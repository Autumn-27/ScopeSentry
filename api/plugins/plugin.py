# -------------------------------------
# @file      : plugin.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/10/24 19:54
# -------------------------------------------
import json
import os
import zipfile
from io import BytesIO

from bson import ObjectId
from fastapi import APIRouter, Depends, File, UploadFile, Query, Form
from pydantic import BaseModel
from starlette.responses import StreamingResponse

from api.users import verify_token
from motor.motor_asyncio import AsyncIOMotorCursor, AsyncIOMotorGridFSBucket
from core.db import get_mongo_db
from core.default import PLUGINS, PLUGINSMODULES
from core.redis_handler import refresh_config, get_redis_pool
from loguru import logger

from core.util import generate_plugin_hash

router = APIRouter()


@router.post("/list")
async def get_all_plugin(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    page_index = request_data.get("pageIndex", 1)
    page_size = request_data.get("pageSize", 10)
    query = request_data.get("search", "")
    query = {
        "name": {"$regex": query, "$options": "i"}
    }
    total_count = await db['plugins'].count_documents(query)
    cursor = db['plugins'].find(query, {"_id": 0,
                                        "id": {"$toString": "$_id"},
                                        "module": 1,
                                        "name": 1,
                                        "hash": 1,
                                        "parameter": 1,
                                        "help": 1,
                                        "introduction": 1,
                                        "isSystem": 1,
                                        "version": 1,
                                        }).skip((page_index - 1) * page_size).limit(page_size)
    result = await cursor.to_list(length=None)
    return {
        "code": 200,
        "data": {
            'list': result,
            'total': total_count
        }
    }


@router.post("/list/bymodule")
async def get_all_plugin_by_module(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    query = request_data.get("module", "")
    if query == "":
        return {"code": 400, "message": "No query provided"}
    query = {
        "module": query
    }
    cursor = db['plugins'].find(query, {"_id": 0,
                                        "id": {"$toString": "$_id"},
                                        "module": 1,
                                        "name": 1,
                                        "hash": 1,
                                        "parameter": 1,
                                        "help": 1,
                                        "introduction": 1,
                                        })
    result = await cursor.to_list(length=None)
    return {
        "code": 200,
        "data": {
            'list': result
        }
    }


@router.post("/detail")
async def get_plugin_detail(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # Get the ID from the request data
        plugin_id = request_data.get("id")

        # Check if ID is provided
        if not plugin_id:
            return {"message": "ID is missing in the request data", "code": 400}

        # Query the database for content based on ID
        query = {"_id": ObjectId(plugin_id)}
        doc = await db.plugins.find_one(query)
        if not doc:
            return {"message": "Content not found for the provided ID", "code": 404}
        result = {
            "name": doc.get("name", ""),
            "module": doc.get("module", ""),
            "hash": doc.get("hash", ""),
            "parameter": doc.get("parameter", ""),
            "help": doc.get("help", ""),
            "introduction": doc.get("introduction", ""),
            "source": doc.get("source", ""),
            "version": doc.get("version", ""),
            "isSystem": doc.get("isSystem", False)
        }
        return {"code": 200, "data": result}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/save")
async def save_plugin(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    with open("PLUGINKEY", 'r') as file:
        plg_key = file.read()
    key = request_data.get("key", "")
    if key == "":
        return {"message": f"key error", "code": 505}
    if plg_key != key:
        return {"message": f"key error", "code": 505}
    id = request_data.get("id", "")
    request_data.pop("id")
    if id is None or id == "":
        # 新建
        request_data["isSystem"] = False
        request_data["hash"] = generate_plugin_hash()
        result = await db.plugins.insert_one(request_data)

        if result.inserted_id:
            await refresh_config('all', 'install_plugin', str(result.inserted_id))
            return {"code": 200, "message": "plugin added successfully"}
        else:
            return {"code": 400, "message": "Failed to add plugin"}
    else:
        update_query = {"_id": ObjectId(id)}

        # Values to be updated
        update_values = {"$set": request_data}

        # Perform the update
        result = await db.plugins.update_one(update_query, update_values)
        if result:
            await refresh_config('all', 'install_plugin', id)
            return {"code": 200, "message": "plugin updated successfully"}
        else:
            return {"code": 404, "message": "plugin not found"}


@router.post("/delete")
async def delete_plugin(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # Extract the list of IDs from the request_data dictionary
        data = request_data.get("data", [])
        if len(data) == 0:
            return {"code": 404, "message": "data null"}
        # Convert the provided rule_ids to ObjectId
        hash_ids = []
        for plugin_info in data:
            hash = plugin_info["hash"]
            module = plugin_info["module"]
            await refresh_config('all', 'delete_plugin', hash + "_" + module)
            if hash not in str(PLUGINS):
                hash_ids.append(hash)
        # Delete the SensitiveRule documents based on the provided IDs
        result = await db.plugins.delete_many({"hash": {"$in": hash_ids}})

        # Check if the deletion was successful
        if result.deleted_count > 0:
            return {"code": 200, "message": "plugin deleted successfully"}
        else:
            return {"code": 404, "message": "plugin not found"}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/log")
async def get_plugin_logs(request_data: dict, _: dict = Depends(verify_token), redis_con=Depends(get_redis_pool)):
    try:
        module = request_data.get("module")
        hash = request_data.get("hash")
        if module is None or hash is None:
            return {"message": "Node name is required", "code": 400}
        # 构建日志键
        log_key = f"logs:plugins:{module}:{hash}"
        # 从 Redis 中获取日志列表
        logs = await redis_con.smembers(log_key)
        log_data = ""
        for log in logs:
            log_data += log + "\n"
        return {"code": 200, "logs": log_data}
    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "Error retrieving logs", "code": 500}


@router.post("/log/clean")
async def clean_plugin_logs(request_data: dict, _: dict = Depends(verify_token), redis_con=Depends(get_redis_pool)):
    try:
        module = request_data.get("module")
        hash = request_data.get("hash")
        if module is None or hash is None:
            return {"message": "Node name is required", "code": 400}
        # 构建日志键
        log_key = f"logs:plugins:{module}:{hash}"
        # 从 Redis 中获取日志列表
        logs = await redis_con.delete(log_key)
        return {"code": 200, "message": "success"}
    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "Error retrieving logs", "code": 500}


ALLOWED_EXTENSIONS = {'.json', '.go', '.md'}


# 检查文件名是否安全
async def is_safe_filename(filename: str) -> bool:
    # 禁止文件名中出现 ..（路径遍历）和其他特殊字符
    if '..' in filename or filename.startswith('/') or filename.startswith('\\'):
        return False

    # 检查文件名的扩展名是否允许
    ext = os.path.splitext(filename)[1].lower()
    if ext not in ALLOWED_EXTENSIONS:
        return False

    # 可以加更多检查文件名长度等规则
    if len(filename) > 255:
        return False

    return True


@router.post("/import")
async def import_plugin(file: UploadFile = File(...), key: str = Query(...), db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        with open("PLUGINKEY", 'r') as PLUGINKEYFile:
            plg_key = PLUGINKEYFile.read()
        if key == "":
            return {"message": f"key error", "code": 505}
        if plg_key != key:
            return {"message": f"key error", "code": 505}
        content = await file.read()
        plugin_info = {}
        source = ""
        # 使用 BytesIO 将字节内容转换为类似文件对象
        with zipfile.ZipFile(BytesIO(content), 'r') as zip_ref:
            for file_name in zip_ref.namelist():
                if file_name.endswith('/'):
                    continue
                # 检查文件名的安全性
                if not await is_safe_filename(file_name):
                    return {"message": f"文件名不安全: {file_name}", "code": 400}
                if os.path.basename(file_name).lower() == "info.json":
                    # 文件名通过安全性检查后，进行后续处理
                    with zip_ref.open(file_name) as f:
                        file_data = f.read()
                        plugin_info = json.loads(file_data)
                if os.path.basename(file_name).lower() == "plugin.go":
                    with zip_ref.open(file_name) as f:
                        file_data = f.read()
                        source = file_data
        if "hash" not in plugin_info:
            plugin_info["hash"] = generate_plugin_hash()
        if "name" not in plugin_info or "module" not in plugin_info:
            return {"message": "node or moudle not found", "code": 500}
        if plugin_info["hash"] in str(PLUGINS) or plugin_info["module"] not in PLUGINSMODULES:
            return {"message": "plugin is system or module error", "code": 500}
        plugin_info["source"] = source
        plugin_info["isSystem"] = False
        result = await db.plugins.insert_one(plugin_info)
        if result.inserted_id:
            await refresh_config('all', 'install_plugin', str(result.inserted_id))
            return {"code": 200, "message": "plugin added successfully"}
        else:
            return {"code": 400, "message": "Failed to add plugin"}
    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/reinstall")
async def reinstall_plugin(request_data: dict, _: dict = Depends(verify_token)):
    node = request_data.get("node", "")
    if node == "":
        return {"message": "node is null", "code": 500}
    hash = request_data.get("hash", "")
    if hash == "":
        return {"message": "plugin hash is null", "code": 500}
    module = request_data.get("module", "")
    if module == "":
        return {"message": "plugin module is null", "code": 500}
    await refresh_config(node, 're_install_plugin', hash + "_" + module)
    return {"code": 200, "message": "success"}


@router.post("/recheck")
async def recheck_plugin(request_data: dict, _: dict = Depends(verify_token)):
    node = request_data.get("node", "")
    if node == "":
        return {"message": "node is null", "code": 500}
    hash = request_data.get("hash", "")
    if hash == "":
        return {"message": "plugin hash is null", "code": 500}
    module = request_data.get("module", "")
    if module == "":
        return {"message": "plugin module is null", "code": 500}
    await refresh_config(node, 're_check_plugin', hash + "_" + module)
    return {"code": 200, "message": "success"}


@router.post("/uninstall")
async def uninstall_plugin(request_data: dict, _: dict = Depends(verify_token)):
    node = request_data.get("node", "")
    if node == "":
        return {"message": "node is null", "code": 500}
    hash = request_data.get("hash", "")
    if hash == "":
        return {"message": "plugin hash is null", "code": 500}
    module = request_data.get("module", "")
    if module == "":
        return {"message": "plugin module is null", "code": 500}
    await refresh_config(node, 'uninstall_plugin', hash + "_" + module)
    return {"code": 200, "message": "success"}


@router.post("/key/check")
async def check_plugin_key(request_data: dict, _: dict = Depends(verify_token)):
    key = request_data.get("key", "")
    with open("PLUGINKEY", 'r') as PLUGINKEYFile:
        plg_key = PLUGINKEYFile.read()
    if key == "":
        return {"message": f"key error", "code": 505}
    if plg_key != key:
        return {"message": f"key error", "code": 505}
    return {"code": 200, "message": "success"}