# -*- coding:utf-8 -*-　　
# @name: configuration
# @auth: rainy-autumn@outlook.com
# @version:
from bson import ObjectId
from fastapi import APIRouter, Depends
from starlette.background import BackgroundTasks
import datetime
from api.users import verify_token
from core.apscheduler_handler import scheduler
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


@router.get("/deduplication/config")
async def get_deduplication_config(_: dict = Depends(verify_token), db=Depends(get_mongo_db)):
    try:
        job = scheduler.get_job("deduplication")
        next_rune_time = ""
        if job is not None:
            next_rune_time = scheduler.get_job("deduplication").next_run_time.strftime("%Y-%m-%d %H:%M:%S")
        result = await db.config.find_one({"name": "deduplication"})
        result["next_run_time"] = next_rune_time
        result.pop("_id")
        return {
            "code": 200,
            "data": result
        }
    except Exception as e:
        logger.error(str(e))
        # 根据需要处理异常
        return {"message": "error", "code": 500}


@router.post("/deduplication/save")
async def save_deduplication_config(request_data: dict, _: dict = Depends(verify_token), db=Depends(get_mongo_db), background_tasks: BackgroundTasks = BackgroundTasks()):
    try:
        run_now = request_data.get("runNow", False)
        request_data.pop("runNow")
        await db.config.update_one(
            {"name": "deduplication"},
            {"$set": request_data},
            upsert=True
        )
        job = scheduler.get_job("deduplication")
        if job is not None:
            scheduler.remove_job("deduplication")
        if request_data.get('flag', False):
            scheduler.add_job(do_asset_deduplication, 'interval', hours=request_data.get('hour', 3),
                              id='deduplication', jobstore='mongo')
        if run_now:
            background_tasks.add_task(do_asset_deduplication)
        return {"message": "Data saved successfully", "code": 200}
    except Exception as e:
        logger.error(str(e))
        return {"message": "error", "code": 500}


async def do_asset_deduplication():
    async for db in get_mongo_db():
        result = await db.config.find_one({"name": "deduplication"})
        result.pop("_id")
        result.pop("name")
        result.pop("hour")
        result.pop("flag")
        f_g_k = {
              "DirScanResult": {
                  "filters": [],
                  "groups": ["url", "status", "msg"]
              },
              "PageMonitoring": {
                  "filters": [],
                  "groups": ["url"]
              },
              "SensitiveResult": {
                  "filters": [],
                  "groups": ["url"]
              },#############
              "SubdoaminTakerResult": {
                  "filters": [],
                  "groups": ["input", "value"]
              },
              "UrlScan": {
                  "filters": [],
                  "groups": ["output"]
              },
              "asset": {
                  "filters": [],
                  "groups": [""]
              },################
              "crawler": {
                  "filters": [],
                  "groups": ["url", "body"]
              },
              "subdomain": {
                  "filters": [],
                  "groups": ["host", "type", "ip"]
              },
              "vulnerability": {
                  "filters": [],
                  "groups": ["url", "vulnid", "matched"]
              }
            }
        for r in result:
            if result[r]:
                await asset_data_dedup(db, r, )


async def asset_data_dedup(db, collection_name, filters, groups):
    # db[].update_many({}, {'$set': {'process_flag': timestamp}})
    # 去重http资产
    logger.info(f"{collection_name} 开始去重")
    collection = db[collection_name]
    timestamp = datetime.datetime.now()
    collection.update_many({}, {'$set': {'process_flag': timestamp}})
    filter = {
        "process_flag": timestamp
    }
    for f in filter:
        filter[f] = filters[f]
    group = {}
    for g in groups:
        group[g.replace(".")] = "$" + g

    pipeline = [
        {
            "$match": filter
        },
        {
            '$sort': {'_id': -1}
        },
        {
            '$group': {
                '_id': group,
                'latestId': {'$first': '$_id'}
            }
        },
        {
            '$project': {'_id': 0, 'latestId': 1}
        }
    ]
    latest_ids = []
    for doc in collection.aggregate(pipeline):
        latest_ids.append(doc['latestId'])
    collection.update_many({'_id': {'$in': latest_ids}}, {'$set': {'latest': True}})
    collection.delete_many({'process_flag': timestamp, 'latest': {'$ne': True}})
    collection.update_many({'process_flag': timestamp}, {'$unset': {'process_flag': "", 'latest': ""}})
    timestamp2 = datetime.datetime.now()
    time_difference = timestamp2 - timestamp
    time_difference_in_seconds = time_difference.total_seconds()
    logger.info(f"{collection_name} 去重消耗时间: {time_difference_in_seconds}")