import time

from loguru import logger
import uvicorn
from starlette.middleware.base import BaseHTTPMiddleware
from starlette.staticfiles import StaticFiles

from core.config import *
from core.update import update14, update15

set_config()

from core.redis_handler import subscribe_log_channel
from core.db import get_mongo_db

from starlette.requests import Request
import asyncio
from fastapi import FastAPI
from fastapi.responses import FileResponse
from fastapi.responses import JSONResponse
from api import dirscan
from core import db
import json
from fastapi import WebSocket
from starlette.exceptions import HTTPException as StarletteHTTPException
from starlette.websockets import WebSocketDisconnect

app = FastAPI(timeout=None)

from core.apscheduler_handler import scheduler


async def update():
    async for db in get_mongo_db():
        # 判断版本
        result = await db.config.find_one({"name": "version"})
        update = False
        if result is not None:
            version = result["version"]
            update = result["update"]
            if version < float(VERSION):
                update = False
        else:
            await db.config.insert_one({"name": "version", "version": float(VERSION), "update": False})
            version = float(VERSION)
        if update is False:
            if version < 1.4:
                await update14(db)
            if version < 1.5:
                await update15(db)


@app.on_event("startup")
async def startup_db_client():
    file_path = os.path.join(os.getcwd(), 'file')
    if not os.path.exists(file_path):
        os.makedirs(file_path)
    await db.create_database()
    await update()
    scheduler.start()
    jobs = scheduler.get_jobs()
    find_page_m = False
    for j in jobs:
        if j.id == 'page_monitoring':
            find_page_m = True
    if not find_page_m:
        from api.task.handler import get_page_monitoring_time, create_page_monitoring_task
        pat, flag = await get_page_monitoring_time()
        if flag:
            scheduler.add_job(create_page_monitoring_task, 'interval', hours=pat, id='page_monitoring',
                              jobstore='mongo')
    asyncio.create_task(subscribe_log_channel())


@app.exception_handler(StarletteHTTPException)
async def http_exception_handler(request, exc):
    if type(exc.detail) == str:
        exc.detail = {'code': 500, 'message': exc.detail}
    return JSONResponse(exc.detail, status_code=exc.status_code)


os.chdir(os.path.dirname(os.path.abspath(__file__)))

from api import users, poc, configuration, fingerprint, node, task, vulnerability, SubdoaminTaker, notification, system, export, project_aggregation
from api.dictionary import router as dictionary_router
from api.asset import router as asset_route
from api.plugins import router as plugin_route
from api.project import router as project_route

app.include_router(plugin_route, prefix='/api')
app.include_router(users.router, prefix='/api')
app.include_router(dictionary_router, prefix='/api/dictionary')
app.include_router(poc.router, prefix='/api')
app.include_router(configuration.router, prefix='/api/configuration')
app.include_router(fingerprint.router, prefix='/api')
app.include_router(node.router, prefix='/api')
app.include_router(project_route, prefix='/api')
app.include_router(task.router, prefix='/api')
app.include_router(asset_route, prefix='/api')
app.include_router(vulnerability.router, prefix='/api')
app.include_router(SubdoaminTaker.router, prefix='/api')
app.include_router(dirscan.router, prefix='/api')
app.include_router(notification.router, prefix='/api')
app.include_router(system.router, prefix='/api')
app.include_router(export.router, prefix='/api')
app.include_router(project_aggregation.router, prefix='/api/project_aggregation')
app.mount("/assets", StaticFiles(directory="static/assets"), name="assets")


@app.get("/logo.png", response_class=FileResponse)
async def get_logo(request: Request):
    return FileResponse("static/logo.png")


@app.get("/favicon.ico", response_class=FileResponse)
async def get_favicon(request: Request):
    return FileResponse("static/favicon.ico")


# @app.middleware("http")
# async def process_http_requests(request, call_next):
#     url = str(request.url)
#     parsed_url = urlparse(url)
#     # 从路径中获取文件名
#     file_name = os.path.basename(parsed_url.path).replace('..', '')
#     # 获取文件后缀名
#     file_extension = os.path.splitext(file_name)[1]
#     if '.html' == file_extension or '.css' == file_extension or '.svg' == file_extension or '.png' == file_extension or '.ico' == file_extension:
#         file_name = file_name.replace('..', '')
#         file_path = os.path.join("static", "assets", file_name)
#         return FileResponse(f"{file_path}")
#     elif '.js' == file_extension:
#         headers = {
#             "Content-Type": "application/javascript; charset=UTF-8"
#         }
#         file_name = file_name.replace('..', '')
#         file_path = os.path.join("static", "assets", file_name)
#         return FileResponse(f"{file_path}", headers=headers)
#     else:
#         response = await call_next(request)
#     return response


@app.get("/")
async def read_root():
    return FileResponse("static/index.html")


# @app.on_event("shutdown")
# async def shutdown_event():
#     global subscriber_task
#     if subscriber_task:
#         subscriber_task.cancel()
#         try:
#             await subscriber_task
#         except asyncio.CancelledError:
#             pass


class MongoDBQueryTimeMiddleware(BaseHTTPMiddleware):
    async def dispatch(self, request: Request, call_next):
        start_time = time.time()
        response = await call_next(request)
        end_time = time.time()
        # 计算查询时间
        query_time = end_time - start_time
        # 获取当前请求的路由信息
        route = request.url.path
        if route.startswith("/api"):
            logger.info(f"MongoDB 查询时间：{query_time} 秒, 路由: {route}")
        return response


SQLTIME = True

if SQLTIME:
    app.add_middleware(MongoDBQueryTimeMiddleware)


@app.websocket("/")
async def websocket_endpoint(websocket: WebSocket):
    await websocket.accept()
    node_name = ""
    try:
        while True:
            data = await websocket.receive_text()
            # 解析收到的消息，假设消息格式为 JSON {"node_name": "example_node"}
            try:
                message = json.loads(data)
                node_name = message.get("node_name")
                if node_name:
                    GET_LOG_NAME.append(node_name)
                    if node_name in LOG_INFO:
                        while LOG_INFO[node_name]:
                            log = LOG_INFO[node_name].pop(0)
                            await websocket.send_text(log)
                else:
                    await websocket.send_text("Invalid message format: missing node_name")
            except json.JSONDecodeError:
                await websocket.send_text("Invalid JSON format")
    except WebSocketDisconnect:
        GET_LOG_NAME.remove(node_name)
        pass


def banner():
    banner = '''   _____                         _____            _              
  / ____|                       / ____|          | |             
 | (___   ___ ___  _ __   ___  | (___   ___ _ __ | |_ _ __ _   _ 
  \___ \ / __/ _ \| '_ \ / _ \  \___ \ / _ \ '_ \| __| '__| | | |
  ____) | (_| (_) | |_) |  __/  ____) |  __/ | | | |_| |  | |_| |
 |_____/ \___\___/| .__/ \___| |_____/ \___|_| |_|\__|_|   \__, |
                  | |                                       __/ |
                  |_|                                      |___/ '''
    print(banner)
    print("Server Version:", VERSION)


if __name__ == "__main__":
    banner()
    file_path = os.path.join(os.getcwd(), "file")
    uvicorn.run("main:app", host="0.0.0.0", port=8082, reload=False, reload_excludes=[file_path])
