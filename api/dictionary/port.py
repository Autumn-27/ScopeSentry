# -------------------------------------
# @file      : port.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/10/17 22:44
# -------------------------------------------

from bson import ObjectId
from fastapi import APIRouter, Depends

from api.users import verify_token
from motor.motor_asyncio import AsyncIOMotorCursor
from core.db import get_mongo_db
from core.redis_handler import refresh_config
from loguru import logger
router = APIRouter()


@router.post("/data")
async def get_port_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        page_index = request_data.get("pageIndex", 1)
        page_size = request_data.get("pageSize", 10)
        search = request_data.get("search", None)  # 获取search参数
        # Construct the search query
        search_query = {}
        if search:
            search_regex = {"$regex": search, "$options": "i"}  # Case-insensitive regex
            search_query = {"$or": [{"name": search_regex}, {"value": search_regex}]}
        total_count = await db.PortDict.count_documents(search_query)

        # Perform pagination query
        cursor: AsyncIOMotorCursor = db.PortDict.find(search_query).skip((page_index - 1) * page_size).limit(page_size)
        result = await cursor.to_list(length=None)

        # Process the result as needed
        response_data = [{"id": str(doc["_id"]),"name": doc["name"], "value": doc["value"]} for doc in result]
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


@router.post("/upgrade")
async def upgrade_port_dict(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # Extract values from request data
        port_id = request_data.get("id")
        name = request_data.get("name")
        value = request_data.get("value")

        # Update query based on rule_id
        update_query = {"_id": ObjectId(port_id)}

        # Values to be updated
        update_values = {"$set": {"name": name, "value": value}}

        # Perform the update
        result = await db.PortDict.update_one(update_query, update_values)
        await refresh_config('all', 'port')
        if result:
            return {"code": 200, "message": "SensitiveRule updated successfully"}
        else:
            return {"code": 404, "message": "SensitiveRule not found"}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}

@router.post("/add")
async def add_port_dict(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # Extract values from request data
        name = request_data.get("name")
        value = request_data.get("value",'')
        if value == '':
            return {"code": 400, "message": "value is null"}
        # Create a new SensitiveRule document
        new_port_dict = {
            "name": name,
            "value": value
        }

        # Insert the new document into the SensitiveRule collection
        result = await db.PortDict.insert_one(new_port_dict)
        await refresh_config('all', 'port')
        # Check if the insertion was successful
        if result.inserted_id:
            return {"code": 200, "message": "Port dict added successfully"}
        else:
            return {"code": 400, "message": "Failed to add port dict"}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/delete")
async def delete_port_dict(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # Extract the list of IDs from the request_data dictionary
        port_dict_ids = request_data.get("ids", [])

        # Convert the provided rule_ids to ObjectId
        obj_ids = [ObjectId(port_dict_id) for port_dict_id in port_dict_ids]

        # Delete the SensitiveRule documents based on the provided IDs
        result = await db.PortDict.delete_many({"_id": {"$in": obj_ids}})
        await refresh_config('all', 'port')
        # Check if the deletion was successful
        if result.deleted_count > 0:
            return {"code": 200, "message": "Port dict deleted successfully"}
        else:
            return {"code": 404, "message": "Port dict not found"}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}