# -------------------------------------
# @file      : page_monitoring.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/4/22 19:46
# -------------------------------------------
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
    cursor: AsyncIOMotorCursor = db.PageMonitoring.find(query, {"url": 1, "_id": 0})
    result = await cursor.to_list(length=None)
    urls = [item['url'] for item in result]

    return urls


@router.post("/page/monitoring/result")
async def page_monitoring_result(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    search_query = request_data.get("search", "")
    page_index = request_data.get("pageIndex", 1)
    page_size = request_data.get("pageSize", 10)
    keyword = {
        'url': 'url',
        'project': 'project',
        'hash': 'hash',
        'diff': 'diff',
        'response': 'response'
    }
    query = await search_to_mongodb(search_query, keyword)
    if query == "" or query is None:
        return {"message": "Search condition parsing error", "code": 500}
    query = query[0]
    # Get the total count of documents matching the search criteria
    query["diff"] = {"$ne": []}
    total_count = await db.PageMonitoring.count_documents(query)
    # Perform pagination query and sort by time
    cursor: AsyncIOMotorCursor = db.PageMonitoring.find(query, {"_id": 0,
                                                                "id": {"$toString": "$_id"},
                                                                "url": 1,
                                                                "content": 1,
                                                                "hash": 1,
                                                                "diff": 1}).sort(
        [("time", DESCENDING)]).skip((page_index - 1) * page_size).limit(page_size)
    result = await cursor.to_list(length=None)
    result_list = []
    for r in result:
        if len(r['content']) < 2:
            continue
        result_list.append({
            "url": r['url'],
            "response1": r['content'][-2],
            "response2": r['content'][-1],
            "hash1": r['hash'][-2],
            "hash2": r['hash'][-1],
            "diff": r['diff'][-1],
            "history_diff": r['diff'][::1]
        })
    return {
        "code": 200,
        "data": {
            'list': result_list,
            'total': total_count
        }
    }