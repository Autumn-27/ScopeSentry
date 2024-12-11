# -------------------------------------
# @file      : SubdoaminTaker.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/4/27 15:41
# -------------------------------------------
from fastapi import APIRouter, Depends
from motor.motor_asyncio import AsyncIOMotorCursor
from api.users import verify_token
from core.db import get_mongo_db
from core.util import search_to_mongodb, get_search_query
from loguru import logger
router = APIRouter()


@router.post("/data")
async def get_subdomaintaker_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        search_query = request_data.get("search", "")
        page_index = request_data.get("pageIndex", 1)
        page_size = request_data.get("pageSize", 10)
        # MongoDB collection for SensitiveRule
        # Fuzzy search based on the name field
        query = await get_search_query("SubdoaminTakerResult", request_data)
        if query == "":
            return {"message": "Search condition parsing error", "code": 500}
        # Get the total count of documents matching the search criteria
        total_count = await db.SubdoaminTakerResult.count_documents(query)
        if total_count == 0:
            return {
            "code": 200,
            "data": {
                'list': [],
                'total': 0
            }
        }
        # Perform pagination query
        cursor: AsyncIOMotorCursor = db.SubdoaminTakerResult.find(query).skip((page_index - 1) * page_size).limit(page_size)
        result = await cursor.to_list(length=None)
        # Process the result as needed
        response_data = []
        for doc in result:
            data = {
                "host": doc["input"],
                "value": doc["value"],
                "type": doc["cname"],
                "response": doc["response"],
                "id": str(doc["_id"]),
                "tags": doc.get("tags", [])
            }
            response_data.append(data)
        return {
            "code": 200,
            "data": {
                'list': response_data,
                'total': total_count
            }
        }

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error","code":500}