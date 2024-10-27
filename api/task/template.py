# -------------------------------------
# @file      : template.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/10/26 22:54
# -------------------------------------------
import asyncio
import traceback

from bson import ObjectId
from fastapi import APIRouter, Depends, BackgroundTasks
from pymongo import DESCENDING

from api.task.util import create_scan_task, task_progress, scheduler_scan_task, delete_asset
from api.users import verify_token
from motor.motor_asyncio import AsyncIOMotorCursor

from core.apscheduler_handler import scheduler
from core.redis_handler import get_redis_pool
from core.util import *

router = APIRouter()


@router.post("/list")
async def template_list(request_data: dict, _: dict = Depends(verify_token), db=Depends(get_mongo_db)):
    page_index = request_data.get("pageIndex", 1)
    page_size = request_data.get("pageSize", 10)
    query = request_data.get("query", "")
    query = {
        "name": {"$regex": query}
    }
    total_count = await db['ScanTemplates'].count_documents(query)
    cursor = db['ScanTemplates'].find(query, {"_id": 0,
                                              "id": {"$toString": "$_id"},
                                              'name': 1,
                                              # 'SubdomainScan': 1,
                                              # 'SubdomainSecurity': 1,
                                              # 'PortScanPreparation': 1,
                                              # 'PortScan': 1,
                                              # 'PortFingerprint': 1,
                                              # 'AssetMapping': 1,
                                              # 'AssetHandle': 1,
                                              # 'URLScan': 1,
                                              # 'WebCrawler': 1,
                                              # 'URLSecurity': 1,
                                              # 'DirScan': 1,
                                              # 'VulnerabilityScan': 1,
                                              # 'parameters': 1
                                              }).skip((page_index - 1) * page_size).limit(page_size)
    result = await cursor.to_list(length=None)
    return {
        "code": 200,
        "data": {
            'list': result,
            'total': total_count
        }
    }


@router.post("/detail")
async def template_detail(request_data: dict, _: dict = Depends(verify_token), db=Depends(get_mongo_db)):
    id = request_data.get("id", "")
    if id == "":
        return {"code": 400, "message": "No id provided"}
    query = {
        "_id": ObjectId(id)
    }
    result = await db['ScanTemplates'].find_one(query, {"_id": 0,
                                                        "id": {"$toString": "$_id"},
                                                        'name': 1,
                                                        'TargetHandler': 1,
                                                        'SubdomainScan': 1,
                                                        'SubdomainSecurity': 1,
                                                        'PortScanPreparation': 1,
                                                        'PortScan': 1,
                                                        'PortFingerprint': 1,
                                                        'AssetMapping': 1,
                                                        'AssetHandle': 1,
                                                        'URLScan': 1,
                                                        'WebCrawler': 1,
                                                        'URLSecurity': 1,
                                                        'DirScan': 1,
                                                        'VulnerabilityScan': 1,
                                                        'Parameters': 1
                                                        })
    if not result:
        return {"message": "template not found for the provided ID", "code": 400}
    return {
        "code": 200,
        "data": result
    }
