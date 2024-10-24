# -------------------------------------
# @file      : plugin.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/10/24 19:54
# -------------------------------------------

from bson import ObjectId
from fastapi import APIRouter, Depends, File, UploadFile, Query, Form
from pydantic import BaseModel
from starlette.responses import StreamingResponse

from api.users import verify_token
from motor.motor_asyncio import AsyncIOMotorCursor, AsyncIOMotorGridFSBucket
from core.db import get_mongo_db
from core.redis_handler import refresh_config
from loguru import logger

router = APIRouter()


@router.post("/list")
async def get_all_plugin(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    page_index = request_data.get("pageIndex", 1)
    page_size = request_data.get("pageSize", 10)
    query = request_data.get("query", "")
    query = {
        "name": {"$regex": query}
    }
    total_count = await db['plugins'].count_documents(query)
    cursor = db['plugins'].find(query, {"_id": 0,
                                        "id": {"$toString": "$_id"},
                                        "module": 1,
                                        "name": 1,
                                        "hash": 1,
                                        "parameter": 1,
                                        "help": 1,
                                        "introduction": 1,
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
async def get_plugin_detail(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # Get the ID from the request data
        plugin_id = request_data.get("id")

        # Check if ID is provided
        if not plugin_id:
            return {"message": "ID is missing in the request data", "code": 400}

        # Query the database for content based on ID
        query = {"_id": ObjectId(plugin_id)}
        doc = await db.plugins.find_one(query)
        if not doc:
            return {"message": "Content not found for the provided ID", "code": 404}
        result = {
            "name": doc.get("name", ""),
            "module": doc.get("module", ""),
            "hash": doc.get("hash", ""),
            "parameter": doc.get("parameter", ""),
            "help": doc.get("help", ""),
            "introduction": doc.get("introduction", ""),
            "source": doc.get("source", "")
        }
        return {"code": 200, "data": result}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}