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
from core.util import search_to_mongodb, get_search_query
from loguru import logger

router = APIRouter()


@router.post("/data")
async def dirscan_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        page_index = request_data.get("pageIndex", 1)
        page_size = request_data.get("pageSize", 10)
        query = await get_search_query("DirScanResult", request_data)
        if query == "":
            return {"message": "Search condition parsing error", "code": 500}
        total_count = await db['DirScanResult'].count_documents(query)
        sort = request_data.get("sort", {})
        sort_by = [('_id', -1)]
        if sort != {}:
            if 'length' in sort:
                sort_value = sort['length']
                if sort_value is not None:
                    if sort_value == "ascending":
                        sort_value = 1
                    else:
                        sort_value = -1
                    sort_by = [('length', sort_value)]
        cursor: AsyncIOMotorCursor = ((db['DirScanResult'].find(query, {"_id": 0, "id": {"$toString": "$_id"}, "url": 1,
                                                                        "status": 1, "msg": 1, "length": 1, "tags": 1})
                                       .sort(sort_by)
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
