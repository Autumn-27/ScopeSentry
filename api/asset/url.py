# -------------------------------------
# @file      : url.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/10/20 21:16
# -------------------------------------------
from fastapi import APIRouter, Depends
from api.users import verify_token
from motor.motor_asyncio import AsyncIOMotorCursor

from core.db import get_mongo_db
from core.util import *
from loguru import logger

router = APIRouter()


@router.post("/data")
async def url_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        page_index = request_data.get("pageIndex", 1)
        page_size = request_data.get("pageSize", 10)
        query = await get_search_query("UrlScan", request_data)
        if query == "":
            return {"message": "Search condition parsing error", "code": 500}
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
        total_count = await db['UrlScan'].count_documents(query)
        cursor: AsyncIOMotorCursor = ((db['UrlScan'].find(query, {"_id": 0,
                                                                  "id": {"$toString": "$_id"},
                                                                  "input": 1,
                                                                  "source": 1,
                                                                  "status": 1,
                                                                  "length": 1,
                                                                  "type": "$outputtype",
                                                                  "url": "$output",
                                                                  "time": 1,
                                                                  "tags": 1,
                                                                  })
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
