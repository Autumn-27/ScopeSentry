# -*- coding:utf-8 -*-　　
# @name: poc_manage
# @auth: rainy-autumn@outlook.com
# @version:
import os

import yaml
from bson import ObjectId
from fastapi import APIRouter, Depends, File, UploadFile
from motor.motor_asyncio import AsyncIOMotorCursor
from api.users import verify_token
from core.db import get_mongo_db
from pymongo import ASCENDING, DESCENDING
from loguru import logger
from core.redis_handler import refresh_config
from core.util import *
import zipfile
router = APIRouter()


@router.post("/poc/data")
async def poc_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        search_query = request_data.get("search", "")
        page_index = request_data.get("pageIndex", 1)
        page_size = request_data.get("pageSize", 10)
        query = {"name": {"$regex": search_query, "$options": "i"}}

        # Get the total count of documents matching the search criteria
        total_count = await db.PocList.count_documents(query)
        # Perform pagination query and sort by time
        cursor: AsyncIOMotorCursor = db.PocList.find(query, {"_id": 0, "id": {"$toString": "$_id"}, "name": 1, "level": 1, "time": 1}).sort([("time", DESCENDING)]).skip((page_index - 1) * page_size).limit(page_size)
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


def is_safe_path(base_path, target_path):
    # 计算规范化路径
    abs_base_path = os.path.abspath(base_path)
    abs_target_path = os.path.abspath(target_path)
    return abs_target_path.startswith(abs_base_path)


@router.post("/poc/data/import")
async def poc_import(file: UploadFile = File(...), db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    if not file.filename.endswith('.zip'):
        return {"message": "not zip", "code": 500}
    file_name = generate_random_string(5)
    relative_path = f'file\\{file_name}.zip'
    zip_file_path = os.path.join(os.getcwd(), relative_path)
    with open(zip_file_path, "wb") as f:
        f.write(await file.read())

    yaml_files = []
    unzip_path = f'file\\{file_name}'
    file_path = os.path.join(os.getcwd(), unzip_path)
    extract_path = file_path
    os.makedirs(extract_path, exist_ok=True)
    with zipfile.ZipFile(zip_file_path, 'r') as zip_ref:
        for member in zip_ref.namelist():
            member_path = os.path.join(extract_path, member)
            if not is_safe_path(extract_path, member_path):
                return {"message": "Unsafe file path detected in ZIP file", "code": 500}
        zip_ref.extractall(extract_path)

        for root, dirs, files in os.walk(extract_path):
            for filename in files:
                if filename.endswith('.yaml') or filename.endswith('.yml'):
                    file_path = os.path.join(root, filename)
                    yaml_files.append(file_path)
    for yaml_file in yaml_files:
        with open(yaml_file, 'r') as stream:
            try:
                data = yaml.safe_load(stream)
                print(data["id"])
            except:
                pass


@router.get("/poc/data/all")
async def poc_data(db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        cursor: AsyncIOMotorCursor = db.PocList.find({}, {"id": {"$toString": "$_id"}, "name": 1, "time": -1, "_id": 0}).sort([("time", DESCENDING)])
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


@router.post("/poc/content")
async def poc_content(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # Get the ID from the request data
        poc_id = request_data.get("id")

        # Check if ID is provided
        if not poc_id:
            return {"message": "ID is missing in the request data", "code": 400}

        # Query the database for content based on ID
        query = {"_id": ObjectId(poc_id)}
        doc = await db.PocList.find_one(query)

        if not doc:
            return {"message": "Content not found for the provided ID", "code": 404}

        # Extract the content
        content = doc.get("content", "")

        return {"code": 200, "data": {"content": content}}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/poc/update")
async def update_poc_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # Get the ID from the request data
        poc_id = request_data.get("id")

        # Check if ID is provided
        if not poc_id:
            return {"message": "ID is missing in the request data", "code": 400}

        # Check if data to update is provided
        if not request_data:
            return {"message": "Data to update is missing in the request", "code": 400}

        # Extract individual fields from the request data
        name = request_data.get("name")
        content = request_data.get("content")
        hash_value = calculate_md5_from_content(content)
        level = request_data.get("level")

        # Prepare the update document
        update_document = {
            "$set": {
                "name": name,
                "content": content,
                "hash": hash_value,
                "level": level
            }
        }

        # Remove the ID from the request data to prevent it from being updated
        del request_data["id"]

        # Update data in the database
        result = await db.PocList.update_one({"_id": ObjectId(poc_id)}, update_document)
        # Check if the update was successful
        if result:
            await refresh_config('all', 'poc')
            return {"message": "Data updated successfully", "code": 200}
        else:
            return {"message": "Failed to update data", "code": 404}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/poc/add")
async def add_poc_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # Check if data to add is provided
        if not request_data:
            return {"message": "Data to add is missing in the request", "code": 400}

        # Extract individual fields from the request data
        name = request_data.get("name")
        content = request_data.get("content")
        hash_value = calculate_md5_from_content(content)
        level = request_data.get("level")
        formatted_time = get_now_time()
        # Insert data into the database
        result = await db.PocList.insert_one({
            "name": name,
            "content": content,
            "hash": hash_value,
            "level": level,
            "time": formatted_time
        })

        # Check if the insertion was successful
        if result.inserted_id:
            await refresh_config('all', 'poc')
            return {"message": "Data added successfully", "code": 200}
        else:
            return {"message": "Failed to add data", "code": 400}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/poc/delete")
async def delete_poc_rules(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # Extract the list of IDs from the request_data dictionary
        poc_ids = request_data.get("ids", [])

        # Convert the provided rule_ids to ObjectId
        obj_ids = [ObjectId(poc_id) for poc_id in poc_ids]

        # Delete the SensitiveRule documents based on the provided IDs
        result = await db.PocList.delete_many({"_id": {"$in": obj_ids}})

        # Check if the deletion was successful
        if result.deleted_count > 0:
            return {"code": 200, "message": "Poc deleted successfully"}
        else:
            return {"code": 404, "message": "Poc not found"}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}
