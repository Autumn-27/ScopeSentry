# -------------------------------------
# @file      : subdomain.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/10/20 21:15
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
async def asset_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        page_index = request_data.get("pageIndex", 1)
        page_size = request_data.get("pageSize", 10)
        query = await get_search_query("subdomain", request_data)
        if query == "":
            return {"message": "Search condition parsing error", "code": 500}
        total_count = await db['subdomain'].count_documents(query)
        cursor: AsyncIOMotorCursor = ((db['subdomain'].find(query, {"_id": 0,
                                                                    "id": {"$toString": "$_id"},
                                                                    "host": 1,
                                                                    "type": 1,
                                                                    "value": 1,
                                                                    "ip": 1,
                                                                    "time": 1,
                                                                    "tags": 1
                                                                    })
                                       .skip((page_index - 1) * page_size)
                                       .limit(page_size))
                                      .sort([("time", DESCENDING)]))
        result = await cursor.to_list(length=None)
        result_list = []
        for r in result:
            if r['value'] is None:
                r['value'] = []
            if r['ip'] is None:
                r['ip'] = []
            result_list.append(r)
        return {
            "code": 200,
            "data": {
                'list': result_list,
                'total': total_count
            }
        }
    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}