# -*- coding:utf-8 -*-　　
# @name: fingerprint
# @auth: rainy-autumn@outlook.com
# @version:
from bson import ObjectId
from fastapi import APIRouter, Depends
from motor.motor_asyncio import AsyncIOMotorCursor
from api.users import verify_token
from core.db import get_mongo_db
from core.redis_handler import refresh_config
from core.util import string_to_postfix
from core.config import APP
from loguru import logger
router = APIRouter()

@router.post("/fingerprint/data")
async def fingerprint_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        search_query = request_data.get("search", "")
        page_index = request_data.get("pageIndex", 1)
        page_size = request_data.get("pageSize", 10)
        query = {"name": {"$regex": search_query, "$options": "i"}}

        # Get the total count of documents matching the search criteria
        total_count = await db.FingerprintRules.count_documents(query)

        # Perform pagination query and sort by time
        cursor: AsyncIOMotorCursor = db.FingerprintRules.find(query, {"_id": 0, "id": {"$toString": "$_id"}, "name": 1, "rule": 1, "category": 1, "parent_category": 1, "amount": 1, "state": 1}).skip((page_index - 1) * page_size).limit(page_size)
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

@router.post("/fingerprint/update")
async def update_fingerprint_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # Get the ID from the request data
        fingerprint_id = request_data.get("id")

        # Check if ID is provided
        if not fingerprint_id:
            return {"message": "ID is missing in the request data", "code": 400}

        # Check if data to update is provided
        if not request_data:
            return {"message": "Data to update is missing in the request", "code": 400}

        # Extract individual fields from the request data
        name = request_data.get("name")
        rule = request_data.get("rule")
        category = request_data.get("category")
        parent_category = request_data.get("parent_category")
        state = request_data.get("state")
        if rule == '':
            return {"code": 500, "message": "rule is null"}
        exp = string_to_postfix(rule)
        if exp == "":
            return {"code": 500, "message": "rule to express error"}
        # Prepare the update document
        update_document = {
            "$set": {
                "name": name,
                "rule": rule,
                "express": exp,
                "category": category,
                "parent_category": parent_category,
                "state": state
            }
        }

        # Remove the ID from the request data to prevent it from being updated
        del request_data["id"]

        # Update data in the database
        result = await db.FingerprintRules.update_one({"_id": ObjectId(fingerprint_id)}, update_document)

        # Check if the update was successful
        if result:
            if fingerprint_id in APP:
                APP[fingerprint_id] = name
            await refresh_config('all', 'finger')
            return {"message": "Data updated successfully", "code": 200}
        else:
            return {"message": "Failed to update data", "code": 404}

    except Exception as e:
        print(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/fingerprint/add")
async def add_fingerprint_rule(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # Extract values from request data
        name = request_data.get("name")
        rule = request_data.get("rule")
        category = request_data.get("category")
        parent_category = request_data.get("parent_category")
        state = request_data.get("state")
        if rule == '':
            return {"code": 500, "message": "rule is null"}
        exp = string_to_postfix(rule)
        if exp == "":
            return {"code": 500, "message": "rule to express error"}
        # Create a new SensitiveRule document
        new_rule = {
            "name": name,
            "rule": rule,
            "category": category,
            "express": exp,
            "parent_category": parent_category,
            'amount': 0,
            "state": state
        }

        # Insert the new document into the SensitiveRule collection
        result = await db.FingerprintRules.insert_one(new_rule)

        # Check if the insertion was successful
        if result.inserted_id:
            if str(result.inserted_id) not in APP:
                APP[str(result.inserted_id)] = name
            await refresh_config('all', 'finger')
            return {"code": 200, "message": "SensitiveRule added successfully"}
        else:
            return {"code": 400, "message": "Failed to add SensitiveRule"}

    except Exception as e:
        print(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}

@router.post("/fingerprint/delete")
async def delete_fingerprint_rules(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # Extract the list of IDs from the request_data dictionary
        fingerprint_ids = request_data.get("ids", [])

        # Convert the provided rule_ids to ObjectId
        obj_ids = [ObjectId(fingerprint_id) for fingerprint_id in fingerprint_ids]

        # Delete the SensitiveRule documents based on the provided IDs
        result = await db.FingerprintRules.delete_many({"_id": {"$in": obj_ids}})

        # Check if the deletion was successful
        if result.deleted_count > 0:
            for fid in fingerprint_ids:
                if fid in APP:
                    del APP[fid]
            await refresh_config('all', 'finger')
            return {"code": 200, "message": "FingerprintRules deleted successfully"}
        else:
            return {"code": 404, "message": "FingerprintRules not found"}

    except Exception as e:
        print(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}