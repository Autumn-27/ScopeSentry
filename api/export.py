# -------------------------------------
# @file      : export.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/6/16 16:11
# -------------------------------------------
import os

from bson import ObjectId
from fastapi import APIRouter, Depends, BackgroundTasks
from openpyxl.utils.exceptions import IllegalCharacterError
from starlette.responses import FileResponse

from api.users import verify_token
from motor.motor_asyncio import AsyncIOMotorCursor
from core.db import get_mongo_db
import pandas as pd
from core.util import *
from pymongo import ASCENDING, DESCENDING, results
from loguru import logger

router = APIRouter()

keywords = {
    "asset": {
        'app': '',
        'body': 'responsebody',
        'header': 'rawheaders',
        'project': 'project',
        'title': 'title',
        'statuscode': 'statuscode',
        'icon': 'faviconmmh3',
        'ip': ['host', 'ip'],
        'domain': ['host', 'url', 'domain'],
        'port': 'port',
        'protocol': ['protocol', 'type'],
        'banner': 'raw',
    },
    "subdomain": {
        'domain': 'host',
        'ip': 'ip',
        'type': 'type',
        'project': 'project',
        'value': 'value'
    },
    "SubdoaminTakerResult": {
        'domain': 'input',
        'value': 'value',
        'type': 'cname',
        'response': 'response',
        'project': 'project',
    },
    "UrlScan": {
        'url': 'output',
        'project': 'project',
        'input': 'input',
        'source': 'source',
        "type": "outputtype"
    },
    "crawler": {
        'url': 'url',
        'method': 'method',
        'body': 'body',
        'project': 'project'
    },
    "SensitiveResult": {
        'url': 'url',
        'sname': 'sid',
        "body": "body",
        "info": "match",
        'project': 'project',
        'md5': 'md5'
    },
    "DirScanResult": {
        'project': 'project',
        'statuscode': 'status',
        'url': 'url',
        'redirect': 'msg'
    },
    "vulnerability": {
        'url': 'url',
        'vulname': 'vulname',
        'project': 'project',
        'matched': 'matched',
        'request': 'request',
        'response': 'response',
        'level': 'level'
    },
    "PageMonitoring": {
        'url': 'url',
        'project': 'project',
        'hash': 'hash',
        'diff': 'diff',
        'response': 'response'
    }
}


@router.post("/export")
async def export_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token),
                      background_tasks: BackgroundTasks = BackgroundTasks()):
    index = request_data.get("index", "")
    quantity = int(request_data.get("quantity", 0))
    export_type = request_data.get("type", "")
    search_query = request_data.get("search", "")
    if index == "" or quantity == 0 or export_type == "":
        return {"code": 500, "message": f"get index, quantity, export_type null"}
    query = {}
    if export_type == "search":
        query = await search_to_mongodb(search_query, keywords[index])
        if query == "" or query is None:
            return {"message": "Search condition parsing error", "code": 500}
        query = query[0]
    if index == "PageMonitoring":
        query["diff"] = {"$ne": []}
    file_name = generate_random_string(16)
    result = await db.export.insert_one({
        "file_name": file_name,
        "create_time": get_now_time(),
        "quantity": quantity,
        "data_type": index,
        "state": 0,
        "end_time": "",
        "file_size": ""
    })
    if result.inserted_id:
        background_tasks.add_task(export_data_from_mongodb, quantity, query, file_name, index, db)
        return {"message": "Successfully added data export task", "code": 200}
    else:
        return {"message": "Failed to export data", "code": 500}


async def fetch_data(db, collection, query, quantity, project_list):
    # 构造替换字段值的pipeline
    branches = []
    for new_value, original_value in project_list.items():
        branches.append({"case": {"$eq": ["$project", original_value]}, "then": new_value})

    pipeline = [
        {"$match": query},
        {"$limit": quantity},
        {"$addFields": {
            "project": {
                "$switch": {
                    "branches": branches,
                    "default": "$project"
                }
            }
        }},
        {"$project": {"_id": 0, "vulnid": 0}}
    ]

    cursor = db[collection].aggregate(pipeline)
    return cursor


async def export_data_from_mongodb(quantity, query, file_name, index, db):
    try:
        cursor = await fetch_data(db, index, query, quantity, Project_List)
        result = await cursor.to_list(length=None)
        relative_path = f'file/{file_name}.xlsx'
        file_path = os.path.join(os.getcwd(), relative_path)
        if index == "asset":
            http_columns = {
                "timestamp": "时间",
                "tlsdata": "TLS_Data",
                "hashes": "Hash",
                "cdnname": "Cdn_Name",
                "port": "端口",
                "url": "url",
                "title": "标题",
                "type": "类型",
                "error": "错误",
                "responsebody": "响应体",
                "host": "IP",
                "faviconmmh3": "图标Hash",
                "faviconpath": "faviconpath",
                "rawheaders": "响应头",
                "jarm": "jarm",
                "technologies": "technologies",
                "statuscode": "响应码",
                "contentlength": "contentlength",
                "cdn": "cdn",
                "webcheck": "webcheck",
                "project": "项目",
                "webfinger": "指纹",
                "iconcontent": "图标",
                "domain": "域名"
            }
            other_columns = {
                "timestamp": "时间",
                "host": "域名",
                "ip": "IP",
                "port": "端口",
                "protocol": "协议",
                "tls": "TLS",
                "transport": "transport",
                "version": "版本",
                "raw": "banner",
                "project": "项目",
                "type": "类型"
            }
            other_df = pd.DataFrame()
            http_df = pd.DataFrame()
            for doc in result:
                if doc["type"] == "other":
                    other_df = pd.concat([other_df, pd.DataFrame([doc])], ignore_index=True)
                else:
                    if doc['webfinger'] is not None:
                        webfinger = []
                        for webfinger_id in doc['webfinger']:
                            webfinger.append(APP[webfinger_id])
                        doc['webfinger'] = webfinger
                    http_df = pd.concat([http_df, pd.DataFrame([doc])], ignore_index=True)
            try:
                excel_writer = pd.ExcelWriter(file_path, engine='xlsxwriter')
                http_df.rename(columns=http_columns, inplace=True)
                http_df.to_excel(excel_writer, sheet_name='HTTP Data', index=False)
                other_df.rename(columns=other_columns, inplace=True)
                other_df.to_excel(excel_writer, sheet_name='Other Data', index=False)
                excel_writer.close()
            except IllegalCharacterError as e:
                logger.error("导出内容有不可见字符，忽略此错误")
        else:
            columns = {}
            if index == "subdomain":
                columns = {'host': '域名', 'type': '解析类型', 'value': '解析值', 'ip': '解析IP', 'project': '项目',
                           'time': '时间'}
            if index == "SubdoaminTakerResult":
                columns = {
                    'input': '源域名', 'value': '解析值', 'cname': '接管类型', 'response': '响应体', 'project': '项目'
                }
            if index == "UrlScan":
                columns = {
                    'input': '输入', 'source': '来源', 'outputtype': '输出类型', 'output': '输出',
                    'statuscode': 'statuscode', 'length': 'length', 'time': '时间', 'project': '项目'
                }
            if index == "crawler":
                columns = {
                    'url': 'URL', 'method': 'Method', 'body': 'Body', 'project': '项目'
                }
            if index == "SensitiveResult":
                columns = {
                    'url': 'URL', 'sid': '规则名称', 'match': '匹配内容', 'project': '项目', 'body': '响应体',
                    'color': '等级', 'time': '时间', 'md5': '响应体MD5'
                }
            if index == "DirScanResult":
                columns = {
                    'url': 'URL', 'status': '响应码', 'msg': '跳转', 'project': '项目'
                }
            if index == "vulnerability":
                columns = {
                    'url': 'URL', 'vulname': '漏洞', 'matched': '匹配', 'project': '项目', 'level': '危害等级',
                    'time': '时间', 'request': '请求', 'response': '响应'
                }
            if index == "PageMonitoring":
                columns = {
                    'url': 'URL', 'content': '响应体', 'hash': '响应体Hash', 'diff': 'Diff',
                    'state': '状态', 'project': '项目', 'time': '时间'
                }
            try:
                df = pd.DataFrame(result)
                df.rename(columns=columns, inplace=True)
                df.to_excel(file_path, index=False)
            except IllegalCharacterError as e:
                logger.error("导出内容有不可见字符，忽略此错误")
        file_size = os.path.getsize(file_path) / (1024 * 1024)  # kb
        update_document = {
            "$set": {
                "state": 1,
                "end_time": get_now_time(),
                "file_size": str(round(file_size, 2))
            }
        }
        await db.export.update_one({"file_name": file_name}, update_document)
    except Exception as e:
        logger.error(str(e))
        update_document = {
            "$set": {
                "state": 2,
            }
        }
        await db.export.update_one({"file_name": file_name}, update_document)


@router.get("/export/record")
async def get_export_record(db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    cursor: AsyncIOMotorCursor = db.export.find({},
                                                {"_id": 0, "id": {"$toString": "$_id"}, "file_name": 1, "end_time": 1,
                                                 "create_time": 1, "data_type": 1, "state": 1, 'file_size': 1}).sort([("create_time", DESCENDING)])
    result = await cursor.to_list(length=None)
    return {
        "code": 200,
        "data": {
            'list': result
        }
    }


@router.post("/export/delete")
async def delete_export(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        export_ids = request_data.get("ids", [])
        delete_filename = []
        for id in export_ids:
            flag = is_valid_string(id)
            if flag and len(id) == 16:
                relative_path = f'file/{id}.xlsx'
                file_path = os.path.join(os.getcwd(), relative_path)
                if os.path.exists(file_path):
                    os.remove(file_path)
                delete_filename.append(id)
        if len(delete_filename) == 0:
            return {"code": 404, "message": "Export file not found"}
        result = await db.export.delete_many({"file_name": {"$in": delete_filename}})

        if result.deleted_count > 0:
            return {"code": 200, "message": "Export file deleted successfully"}
        else:
            return {"code": 404, "message": "Export file not found"}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.get("/export/download")
async def download_export(file_name: str):
    if len(file_name) == 16 and is_valid_string(file_name):
        relative_path = f'file/{file_name}.xlsx'
        file_path = os.path.join(os.getcwd(), relative_path)
        if os.path.exists(file_path):
            return FileResponse(path=file_path, filename=file_name + '.xlsx')
        else:
            return {"message": "file not found", "code": 500}
    else:
        return {"message": "file not found", "code": 500}