# -------------------------------------
# @file      : notification.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/5/12 11:37
# -------------------------------------------
from bson import ObjectId
from fastapi import APIRouter, Depends
from pymongo import DESCENDING

from api.users import verify_token
from motor.motor_asyncio import AsyncIOMotorCursor

from core.db import get_mongo_db
from core.redis_handler import refresh_config
from loguru import logger

router = APIRouter()


@router.get("/notification/data")
async def get_notification_data(db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        cursor: AsyncIOMotorCursor = db.notification.find({}, {"id": {"$toString": "$_id"},  "_id": 0, "name": 1, "method": 1, "url": 1, "contentType": 1,"data": 1, "state": 1})
        result = await cursor.to_list(length=None)
        return {
            "code": 200,
            "data": {
                'list': result
            }
        }
    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/notification/add")
async def add_notification_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        if not request_data:
            return {"message": "Data to add is missing in the request", "code": 400}

        result = await db.notification.insert_one(request_data)

        if result.inserted_id:
            await refresh_config('all', 'notification')
            return {"message": "Data added successfully", "code": 200}
        else:
            return {"message": "Failed to add data", "code": 400}

    except Exception as e:
        logger.error(str(e))
        return {"message": "error", "code": 500}


@router.post("/notification/update")
async def update_notification_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # Get the ID from the request data
        not_id = request_data.get("id")

        # Check if ID is provided
        if not not_id:
            return {"message": "ID is missing in the request data", "code": 400}

        # Check if data to update is provided
        if not request_data:
            await refresh_config('all', 'notification')
            return {"message": "Data to update is missing in the request", "code": 400}
        del request_data["id"]
        update_document = {
            "$set": request_data
        }

        result = await db.notification.update_one({"_id": ObjectId(not_id)}, update_document)
        # Check if the update was successful
        if result:
            await refresh_config('all', 'notification')
            return {"message": "Data updated successfully", "code": 200}
        else:
            return {"message": "Failed to update data", "code": 404}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/notification/delete")
async def delete_notification(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        not_ids = request_data.get("ids", [])

        # Convert the provided rule_ids to ObjectId
        obj_ids = [ObjectId(not_id) for not_id in not_ids]

        # Delete the SensitiveRule documents based on the provided IDs
        result = await db.notification.delete_many({"_id": {"$in": obj_ids}})

        if result.deleted_count > 0:
            await refresh_config('all', 'notification')
            return {"code": 200, "message": "Notification deleted successfully"}
        else:
            return {"code": 404, "message": "Notification not found"}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.get("/notification/config/data")
async def get_notification_config_data(db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        result = await db.config.find_one({"name": "notification"})
        del result['_id']
        del result['type']
        del result['name']
        return {
            "code": 200,
            "data": result
        }
    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/notification/config/update")
async def update_notification_config_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        update_document = {
            "$set": request_data
        }

        result = await db.config.update_one({"name": "notification"}, update_document)
        # Check if the update was successful
        if result:
            await refresh_config('all', 'notification')
            return {"message": "Data updated successfully", "code": 200}
        else:
            return {"message": "Failed to update data", "code": 404}
    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}
