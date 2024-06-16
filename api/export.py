# -------------------------------------
# @file      : export.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/6/16 16:11
# -------------------------------------------
import os

from bson import ObjectId
from fastapi import APIRouter, Depends, BackgroundTasks
from api.users import verify_token
from motor.motor_asyncio import AsyncIOMotorCursor
from core.db import get_mongo_db
import pandas as pd
from core.util import *
from pymongo import ASCENDING, DESCENDING
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
async def export_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token), background_tasks: BackgroundTasks = BackgroundTasks()):
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
        "end_time": ""
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
    cursor = await fetch_data(db, index, query, quantity, Project_List)
    result = await cursor.to_list(length=None)
    relative_path = f'file/{file_name}.xlsx'
    file_path = os.path.join(os.getcwd(), relative_path)
    if index == "asset":
        http_df = None
        other_df = None
        http_columns = ['时间', 'TLS_Data', 'Hash', 'Cdn_Name', '端口', 'url', '标题', '类型', '错误', '响应体', 'IP', '图标Hash', 'faviconpath', '响应头', 'jarm', 'technologies', '响应码', 'contentlength', 'cdn', 'webcheck', '项目', '指纹', '域名']
        other_columns = ['时间', '域名', 'IP', '端口', '协议', "TLS", '版本', 'banner', '项目', '类型']
        for doc in result:
            if doc["type"] == "other":
                other_df = pd.DataFrame(doc, columns=other_columns)
            else:
                if len(doc['webfinger']) != 0:
                    webfinger = []
                    for webfinger_id in doc['webfinger']:
                        webfinger.append(APP[webfinger_id])
                    doc['webfinger'] = webfinger
                http_df = pd.DataFrame(doc, columns=http_columns)
        excel_writer = pd.ExcelWriter(file_path, engine='xlsxwriter')
        http_df.to_excel(excel_writer, sheet_name='HTTP Data', index=False)
        other_df.to_excel(excel_writer, sheet_name='Other Data', index=False)
        excel_writer.save()
    else:
        columns = []
        if index == "subdomain":
            columns = ['域名', '解析类型', '解析值', '解析IP', '项目']
        if index == "SubdoaminTakerResult":
            columns = ['源域名', '解析值', '接管类型', '响应体', '项目']
        if index == "UrlScan":
            columns = ['输入', '来源', '输出', 'statuscode', 'length', '时间', '项目']
        if index == "crawler":
            columns = ['URL', 'Method', 'Body', '项目']
        if index == "SensitiveResult":
            columns = ['URL', '规则名称', '匹配内容', '项目', '响应体', '等级', '时间', '响应体MD5']
        if index == "DirScanResult":
            columns = ['URL', '响应码', '跳转', '项目']
        if index == "vulnerability":
            columns = ['URL', '漏洞', '匹配', '项目', '危害等级', '时间', '请求', '响应']
        if index == "PageMonitoring":
            columns = ['URL', '响应体', '响应体Hash', 'Diff', '状态', '项目', '时间']




