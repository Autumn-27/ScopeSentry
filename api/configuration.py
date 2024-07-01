# -*- coding:utf-8 -*-　　
# @name: configuration
# @auth: rainy-autumn@outlook.com
# @version:
from bson import ObjectId
from fastapi import APIRouter, Depends
from api.users import verify_token
from core.db import get_mongo_db
from core.redis_handler import refresh_config
from core.config import set_timezone
from loguru import logger
router = APIRouter()

@router.get("/subfinder/data")
async def get_subfinder_data(db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # Find document with name equal to "DomainDic"
        result = await db.config.find_one({"name": "SubfinderApiConfig"})
        return {
            "code": 200,
            "data": {
                "content": result.get("value", '')
            }
        }

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error","code":500}

@router.post("/subfinder/save")
async def save_subfinder_data(data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # Update the document with name equal to "DomainDic"
        result = await db.config.update_one({"name": "SubfinderApiConfig"}, {"$set": {"value": data.get('content','')}}, upsert=True)
        if result:
            await refresh_config('all', 'subfinder')
            return {"code": 200, "message": "Successfully updated SubfinderApiConfig value"}
        else:
            return {"code": 404, "message": "SubfinderApiConfig not found"}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}

@router.get("/rad/data")
async def get_rad_data(db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # Find document with name equal to "DomainDic"
        result = await db.config.find_one({"name": "RadConfig"})
        return {
            "code": 200,
            "data": {
                "content": result.get("value", '')
            }
        }

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error","code":500}

@router.post("/rad/save")
async def save_rad_data(data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # Update the document with name equal to "DomainDic"
        result = await db.config.update_one({"name": "RadConfig"}, {"$set": {"value": data.get('content','')}}, upsert=True)
        if result:
            await refresh_config('all', 'rad')
            return {"code": 200, "message": "Successfully updated RadConfig value"}
        else:
            return {"code": 404, "message": "SubfinderApiConfig not found"}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.get("/system/data")
async def get_system_data(db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # 查询所有 type 为 "system" 的文档
        cursor = db.config.find({"type": "system"})
        system_data = {}

        async for document in cursor:
            # 提取 name 和 value 字段，并添加到 system_data 中
            system_data[document["name"]] = document["value"]

        return {
            "code": 200,
            "data": system_data
        }

    except Exception as e:
        logger.error(str(e))
        # 根据需要处理异常
        return {"message": "error", "code": 500}

@router.post("/system/save")
async def save_system_data(data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        for key, value in data.items():
            if key == 'timezone':
                set_timezone(value)
            # 使用键来查找并更新相应的文档
            await db.config.update_one(
                {"type": "system", "name": key},
                {"$set": {"value": value}},
                upsert=True
            )
        await refresh_config('all', 'system')
        return {"message": "Data saved successfully", "code": 200}

    except Exception as e:
        return {"message": "error", "code": 500}


@router.get("/system/deduplication/config")
async def get_deduplication_config(_: dict = Depends(verify_token), db=Depends(get_mongo_db)):
    try:
        # 查询所有 type 为 "system" 的文档
        cursor = await db.config.find_one({"name": "deduplication"})
        deduplication_data = {}

        async for document in cursor:
            deduplication_data[document["name"]] = document["value"]
        return {
            "code": 200,
            "data": deduplication_data
        }
    except Exception as e:
        logger.error(str(e))
        # 根据需要处理异常
        return {"message": "error", "code": 500}