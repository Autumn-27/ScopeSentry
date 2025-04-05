# -------------------------------------
# @file      : app.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2025/4/5 22:31
# -------------------------------------------
from fastapi import APIRouter, Depends
from api.users import verify_token
from motor.motor_asyncio import AsyncIOMotorCursor

from core.db import get_mongo_db
from core.util import *
from pymongo import DESCENDING
from loguru import logger

router = APIRouter()


@router.post("/data")
async def app_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    page_index = request_data.get("pageIndex", 1)
    page_size = request_data.get("pageSize", 10)
    query = await get_search_query("app", request_data)
    if query == "":
        return {"message": "Search condition parsing error", "code": 500}
    if len(query) == 0:
        total_count = await db['app'].count_documents({"_id": {"$exists": True}})
    else:
        total_count = await db['app'].count_documents(query)
    cursor: AsyncIOMotorCursor = ((db['app'].find(query, {"_id": 0,
                                                                 "id": {"$toString": "$_id"},
                                                                 "name": 1,
                                                                 "icp": 1,
                                                                 "company": 1,
                                                                 "project": 1,
                                                                 "time": 1,
                                                                 "tags": 1,
                                                                 "category": 1,
                                                                 "description": 1,
                                                                 "bundleID": 1,
                                                                 "apk": 1,
                                                                 "url": 1
                                                                 })
                                   .skip((page_index - 1) * page_size)
                                   .limit(page_size))
                                  .sort([("time", DESCENDING)]))
    result = await cursor.to_list(length=None)
    result_list = []
    for i in result:
        if i["project"] in Project_List:
            i["project"] = Project_List[i["project"]]
        else:
            i["project"] = ""
        result_list.append(i)
    return {
        "code": 200,
        "data": {
            'list': result_list,
            'total': total_count
        }
    }
