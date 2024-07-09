# -------------------------------------
# @file      : project_aggregation.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/7/8 21:02
# -------------------------------------------

import time
import traceback

from bson import ObjectId
from fastapi import APIRouter, Depends, BackgroundTasks
from pymongo import DESCENDING

from api.task import create_scan_task, delete_asset
from api.users import verify_token
from motor.motor_asyncio import AsyncIOMotorCursor

from core.config import Project_List
from core.db import get_mongo_db
from core.redis_handler import refresh_config, get_redis_pool
from loguru import logger
from core.util import *
from core.apscheduler_handler import scheduler

router = APIRouter()


@router.post("/project/info")
async def get_projects_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    id = request_data.get("id", "")
    result = await db.project.find_one({"_id": ObjectId(id)}, {
                                                                    "_id": 0,
                                                                    "tag": 1,
                                                                    "hour": 1,
                                                                    "scheduledTasks": 1,
                                                                    "AssetCount": 1,
                                                                    "root_domains": 1,
                                                                    "name":1
                                                                }
                                       )
    if result['scheduledTasks']:
        job = scheduler.get_job(id)
        if job is not None:
            next_time = job.next_run_time.strftime("%Y-%m-%d %H:%M:%S")
            result['next_time'] = next_time
    return {"code": 200, "data": result}


@router.post("/project/asset/count")
async def get_projects_asset_count(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    id = request_data.get("id", "")
    subdomain_count = await db['subdomain'].count_documents({"project": id})
    vulnerability_count = await db['vulnerability'].count_documents({"project": id})
    return {"code": 200, "data": {
        "subdomainCount": subdomain_count,
        "vulCount": vulnerability_count
    }}


@router.post("/project/vul/statistics")
async def get_projects_vul_statistics(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    id = request_data.get("id", "")
    pipeline = [
        {"$match": {"project": id}},
        {
            "$group": {
                "_id": "$level",
                "count": {"$sum": 1}
            }
        }
    ]
    result = await db['vulnerability'].aggregate(pipeline).to_list(None)
    return {"code": 200, "data": result}


@router.post("/project/vul/data")
async def get_projects_vul_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    id = request_data.get("id", "")
    cursor: AsyncIOMotorCursor = db.vulnerability.find({"project": id}, {"_id": 0, "url": 1, "vulname": 1, "level": 1, "time": 1, "matched": 1}).sort([("time", DESCENDING)])
    result = await cursor.to_list(length=None)
    return {
        "code": 200,
        "data": {
            'list': result
        }
    }
