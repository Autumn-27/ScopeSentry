# -------------------------------------
# @file      : dirscan.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/5/9 20:22
# -------------------------------------------

from bson import ObjectId
from fastapi import APIRouter, Depends
from motor.motor_asyncio import AsyncIOMotorCursor
from api.users import verify_token
from core.db import get_mongo_db
from core.util import search_to_mongodb
from loguru import logger
router = APIRouter()


@router.post("/dirscan/result/data")
async def dirscan_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        search_query = request_data.get("search", "")
        page_index = request_data.get("pageIndex", 1)
        page_size = request_data.get("pageSize", 10)
        keyword = {
            'project': 'project',
            'statuscode': 'status',
            'url': 'url',
            'redirect': 'msg',
            'length': 'length'
        }
        query = await search_to_mongodb(search_query, keyword)
        if query == "" or query is None:
            return {"message": "Search condition parsing error", "code": 500}
        query = query[0]
        total_count = await db['DirScanResult'].count_documents(query)
        cursor: AsyncIOMotorCursor = ((db['DirScanResult'].find(query, {"_id": 0, "id": {"$toString": "$_id"}, "url": 1, "status": 1, "msg":1, "length": 1})
                                       .sort([('_id', -1)])
                                       .skip((page_index - 1) * page_size)
                                       .limit(page_size)))
        result = await cursor.to_list(length=None)
        return {
            "code": 200,
            "data": {
                'list': result,
                'total': total_count
            }
        }
    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}