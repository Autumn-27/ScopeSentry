# -------------------------------------
# @file      : page_monitoring.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/4/22 19:46
# -------------------------------------------
import json

from bson import ObjectId
from fastapi import APIRouter, Depends
from motor.motor_asyncio import AsyncIOMotorCursor
from api.users import verify_token
from core.db import get_mongo_db
from pymongo import ASCENDING, DESCENDING
from loguru import logger
from core.redis_handler import refresh_config
from core.util import *

router = APIRouter()


async def get_page_monitoring_data(db, all):
    if all:
        query = {}
    else:
        query = {"state": 1}
    cursor: AsyncIOMotorCursor = db.PageMonitoring.find(query, {"url": 1, "_id": 0, "hash": 1, "statusCode": 1, "md5": 1} )
    result = await cursor.to_list(length=None)
    urls = []
    for url in result:
        urls.append(json.dumps(url))
    return urls


@router.post("/result")
async def page_monitoring_result(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    page_index = request_data.get("pageIndex", 1)
    page_size = request_data.get("pageSize", 10)
    query = await get_search_query("PageMonitoring", request_data)
    if query == "":
        return {"message": "Search condition parsing error", "code": 500}
    query["hash"] = {"$size": 2}
    total_count = await db.PageMonitoring.count_documents(query)
    # Perform pagination query and sort by time
    cursor: AsyncIOMotorCursor = db.PageMonitoring.find(query, {"_id": 0,
                                                                "id": {"$toString": "$_id"},
                                                                "url": 1,
                                                                "hash": 1,
                                                                "time": 1,
                                                                "md5": 1,
                                                                "statusCode": 1,
                                                                "similarity": 1,
                                                                "tags": 1
                                                                }).sort(
        [("time", DESCENDING)]).skip((page_index - 1) * page_size).limit(page_size)
    result = await cursor.to_list(length=None)
    return {
        "code": 200,
        "data": {
            'list': result,
            'total': total_count
        }
    }


@router.post("/response")
async def monitoring_response(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # Get the ID from the request data
        monitoring_id = request_data.get("id")
        flag = request_data.get("flag")  # 1代表获取上一次的响应 2代表获取当前响应
        # Check if ID is provided
        if not monitoring_id:
            return {"message": "ID is missing in the request data", "code": 400}

        # Query the database for content based on ID
        query = {"_id": ObjectId(monitoring_id)}
        doc = await db.PageMonitoring.find_one(query)

        if not doc:
            return {"message": "Content not found for the provided ID", "code": 404}

        # Extract the content
        contents = doc.get("content", [])
        hashes = doc.get("hash", [])
        if flag == "1":
            content = contents[-2]
            c_hash = hashes[-2]
        else:
            content = contents[-1]
            c_hash = hashes[-1]
        return {"code": 200, "data": {"content": content, "hash": c_hash}}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/diff")
async def monitoring_history_diff(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # Get the ID from the request data
        monitoring_id = request_data.get("id")
        # Check if ID is provided
        if not monitoring_id:
            return {"message": "ID is missing in the request data", "code": 400}

        # Query the database for content based on ID
        query = {"md5": monitoring_id}
        doc = await db.PageMonitoringBody.find_one(query)

        if not doc:
            return {"message": "Content not found for the provided ID", "code": 404}

        diff = doc.get("content", ["", ""])
        return {"code": 200, "data": {"diff": diff}}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}