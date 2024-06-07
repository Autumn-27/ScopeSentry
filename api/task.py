# -*- coding:utf-8 -*-　　
# @name: sensitive
# @auth: rainy-autumn@outlook.com
# @version:
import asyncio
import json
from loguru import logger
from bson import ObjectId
from fastapi import APIRouter, Depends, BackgroundTasks
from pymongo import DESCENDING

from api.users import verify_token
from motor.motor_asyncio import AsyncIOMotorCursor

from core.apscheduler_handler import scheduler
from core.db import get_mongo_db
from core.redis_handler import get_redis_pool
from core.util import *
from api.node import get_redis_online_data
from api.page_monitoring import get_page_monitoring_data
router = APIRouter()


@router.post("/task/data")
async def get_task_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token), background_tasks: BackgroundTasks = BackgroundTasks()):
    try:
        background_tasks.add_task(task_progress)
        search_query = request_data.get("search", "")
        page_index = request_data.get("pageIndex", 1)
        page_size = request_data.get("pageSize", 10)
        # Fuzzy search based on the name field
        query = {"name": {"$regex": search_query, "$options": "i"}}

        # Get the total count of documents matching the search criteria
        total_count = await db.task.count_documents(query)

        # Perform pagination query
        cursor: AsyncIOMotorCursor = db.task.find(query).skip((page_index - 1) * page_size).limit(page_size).sort([("creatTime", DESCENDING)])
        result = await cursor.to_list(length=None)
        if len(result) == 0:
            return {
            "code": 200,
            "data": {
                'list': [],
                'total': 0
            }
        }
        # Process the result as needed
        response_data = [{"id": str(doc["_id"]), "name": doc["name"], "taskNum": doc["taskNum"], "progress": doc["progress"], "creatTime": doc["creatTime"], "endTime": doc["endTime"]} for doc in result]
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


@router.post("/task/add")
async def add_task(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token), redis_con=Depends(get_redis_pool)):
    try:
        name = request_data.get("name")
        target = request_data.get("target", "")
        node = request_data.get("node")
        if name == "" or target == "" or node == []:
            return {"message": "Null", "code": 500}
        scheduledTasks = request_data.get("scheduledTasks", False)
        hour = request_data.get("hour", 1)
        targetList = []
        targetTmp = ""
        for t in target.split("\n"):
            t.replace("http://","").replace("https://","")
            t = t.strip("\n").strip("\r")
            if t != "" and t not in targetList:
                targetList.append(t)
                targetTmp += t + "\n"
        taskNum = len(targetList)
        request_data['taskNum'] = taskNum
        request_data['target'] = targetTmp.strip("\n")
        request_data['progress'] = 0
        request_data["creatTime"] = get_now_time()
        request_data["endTime"] = ""
        if "All Poc" in request_data['vulList']:
            request_data['vulList'] = ["All Poc"]
        result = await db.task.insert_one(request_data)
        # Check if the insertion was successful
        if result.inserted_id:
            if scheduledTasks:
                scheduler.add_job(scheduler_scan_task, 'interval', hours=hour, args=[str(result.inserted_id)],
                                  id=str(result.inserted_id), jobstore='mongo')
                next_time = scheduler.get_job(str(result.inserted_id)).next_run_time
                formatted_time = next_time.strftime("%Y-%m-%d %H:%M:%S")
                db.ScheduledTasks.insert_one(
                    {"id": str(result.inserted_id), "name": name, 'hour': hour, 'type': 'Scan', 'state': True, 'lastTime': get_now_time(), 'nextTime': formatted_time, 'runner_id': str(result.inserted_id)})
            f = await create_scan_task(request_data, result.inserted_id, targetList, redis_con)
            if f:
                return {"code": 200, "message": "Task added successfully"}
            else:
                return {"code": 400, "message": "Failed to add Task"}
        else:
            return {"code": 400, "message": "Failed to add Task"}
    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/task/content")
async def task_content(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # Get the ID from the request data
        task_id = request_data.get("id")

        # Check if ID is provided
        if not task_id:
            return {"message": "ID is missing in the request data", "code": 400}

        # Query the database for content based on ID
        query = {"_id": ObjectId(task_id)}
        doc = await db.task.find_one(query)
        if not doc:
            return {"message": "Content not found for the provided ID", "code": 404}
        result = {
            "name": doc.get("name", ""),
            "target": doc.get("target", ""),
            "node": doc.get("node", []),
            "subdomainScan": doc.get("subdomainScan", False),
            "subdomainConfig": doc.get("subdomainConfig", []),
            "urlScan": doc.get("urlScan", False),
            "sensitiveInfoScan": doc.get("sensitiveInfoScan", False),
            "pageMonitoring": doc.get("pageMonitoring", ""),
            "crawlerScan": doc.get("crawlerScan", False),
            "vulScan": doc.get("vulScan", False),
            "vulList": doc.get("vulList", []),
            "portScan": doc.get("portScan"),
            "ports": doc.get("ports"),
            "waybackurl": doc.get("waybackurl"),
            "dirScan": doc.get("dirScan"),
            "scheduledTasks": doc.get("scheduledTasks"),
            "hour": doc.get("hour"),
            "duplicates": doc.get("duplicates")
        }

        return {"code": 200, "data": result}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/task/delete")
async def delete_task(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token), redis_con=Depends(get_redis_pool)):
    try:
        # Extract the list of IDs from the request_data dictionary
        task_ids = request_data.get("ids", [])

        # Convert the provided rule_ids to ObjectId
        obj_ids = []
        redis_key = []
        for task_id in task_ids:
            obj_ids.append(ObjectId(task_id))
            redis_key.append("TaskInfo:" + task_id)
        await redis_con.delete(*redis_key)
        # Delete the SensitiveRule documents based on the provided IDs
        result = await db.task.delete_many({"_id": {"$in": obj_ids}})

        # Check if the deletion was successful
        if result.deleted_count > 0:
            return {"code": 200, "message": "Task deleted successfully"}
        else:
            return {"code": 404, "message": "Task not found"}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/task/retest")
async def retest_task(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token), redis_con=Depends(get_redis_pool)):
    try:
        # Get the ID from the request data
        task_id = request_data.get("id")

        # Check if ID is provided
        if not task_id:
            return {"message": "ID is missing in the request data", "code": 400}

        # Query the database for content based on ID
        query = {"_id": ObjectId(task_id)}
        doc = await db.task.find_one(query)
        if not doc:
            return {"message": "Content not found for the provided ID", "code": 404}
        target = doc['target']
        targetList = target.split("\n")
        keys_to_delete = [
            f"TaskInfo:tmp:{task_id}",
            f"TaskInfo:{task_id}",
            f"TaskInfo:time:{task_id}",
            f"duplicates:url:{task_id}",
            f"duplicates:domain:{task_id}",
            f"duplicates:sensresp:{task_id}",
            f"duplicates:craw:{task_id}"
        ]
        progresskeys = await redis_con.keys(f"TaskInfo:progress:{task_id}:*")
        keys_to_delete.extend(progresskeys)

        if keys_to_delete:
            await redis_con.delete(*keys_to_delete)

        f = await create_scan_task(doc, task_id, targetList, redis_con)

        if f:
            update_document = {
                "$set": {
                    "progress": 0,
                    "creatTime": get_now_time(),
                    "endTime": ""
                }
            }
            await db.task.update_one({"_id": ObjectId(task_id)}, update_document)
            return {"code": 200, "message": "Task added successfully"}
        else:
            return {"code": 400, "message": "Failed to add Task"}
    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}



async def create_scan_task(request_data, id, targetList, redis_con):
    try:
        request_data["id"] = str(id)
        if request_data['allNode']:
            request_data["node"] = await get_redis_online_data(redis_con)

        keys_to_delete = [
            f"TaskInfo:tmp:{id}",
            f"TaskInfo:{id}",
            f"TaskInfo:time:{id}",
            f"duplicates:url:{id}",
            f"duplicates:domain:{id}",
            f"duplicates:sensresp:{id}",
            f"duplicates:craw:{id}"
        ]
        progresskeys = await redis_con.keys(f"TaskInfo:progress:{id}:*")
        keys_to_delete.extend(progresskeys)
        if keys_to_delete:
            await redis_con.delete(*keys_to_delete)
        add_redis_task_data = transform_db_redis(request_data)
        async with redis_con as redis:
            await redis.lpush(f"TaskInfo:{id}", *targetList)
            for name in request_data["node"]:
                await redis.rpush(f"NodeTask:{name}", json.dumps(add_redis_task_data))
        return True
    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return False


@router.post("/task/update")
async def update_task_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # Get the ID from the request data
        task_id = request_data.get("id")
        hour = request_data.get("hour")
        # Check if ID is provided
        if not task_id:
            return {"message": "ID is missing in the request data", "code": 400}
        query = {"id": task_id}
        doc = await db.ScheduledTasks.find_one(query)
        oldScheduledTasks = doc["state"]
        old_hour = doc["hour"]
        newScheduledTasks = request_data.get("scheduledTasks")
        if oldScheduledTasks != newScheduledTasks:
            if newScheduledTasks:
                scheduler.add_job(scheduler_scan_task, 'interval', hours=hour, args=[task_id],
                                  id=str(task_id), jobstore='mongo')
                await db.ScheduledTasks.update_one({"id": task_id}, {"$set": {'state': True}})
            else:
                scheduler.remove_job(task_id)
                await db.ScheduledTasks.update_one({"id": task_id}, {"$set": {'state': False}})
        if newScheduledTasks:
            if hour != old_hour:
                job = scheduler.get_job(task_id)
                if job is not None:
                    scheduler.remove_job(job)
                    scheduler.add_job(scheduler_scan_task, 'interval', hours=hour, args=[task_id],
                                      id=str(task_id), jobstore='mongo')

        request_data.pop("id")
        update_document = {
            "$set": request_data
        }
        result = await db.task.update_one({"_id": ObjectId(task_id)}, update_document)
        # Check if the update was successful
        if result:
            return {"message": "Task updated successfully", "code": 200}
        else:
            return {"message": "Failed to update data", "code": 404}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


async def task_progress():
    async for db in get_mongo_db():
        async for redis in get_redis_pool():
            query = {"progress": {"$ne": 100}}
            cursor: AsyncIOMotorCursor = db.task.find(query)
            result = await cursor.to_list(length=None)
            if len(result) == 0:
                return True
            for r in result:
                id = str(r["_id"])
                key = f"TaskInfo:tmp:{id}"
                exists = await redis.exists(key)
                if exists:
                    count = await redis.llen(key)
                    progress_tmp = round(count / r['taskNum'], 2)
                    progress_tmp = round(progress_tmp * 100, 1)
                    if progress_tmp > 100:
                        progress_tmp = 100
                    if progress_tmp == 100:
                        time_key = f"TaskInfo:time:{id}"
                        time_value = await redis.get(time_key)
                        await db.task.update_one({"_id": r["_id"]}, {"$set": {"endTime": time_value}})
                    await db.task.update_one({"_id": r["_id"]}, {"$set": {"progress": progress_tmp}})
                else:
                    await db.task.update_one({"_id": r["_id"]}, {"$set": {"progress": 0}})
            return


# @router.post("/task/progress/info")
# async def progress_info(request_data: dict,  _: dict = Depends(verify_token), redis_con=Depends(get_redis_pool), db=Depends(get_mongo_db)):
#     task_id = request_data.get("id")
#     type = request_data.get("type")
#     runner = request_data.get("runner")
#     # Check if ID is provided
#     if not task_id:
#         return {"message": "ID is missing in the request data", "code": 400}
#     query = {"_id": ObjectId(task_id)}
#     if type == "scan":
#         doc = await db.task.find_one(query)
#     else:
#         doc = await db.project.find_one(query)
#         target_data = await db.ProjectTargetData.find_one({"id": task_id})
#         doc["target"] = target_data["target"]
#     if not doc:
#         return {"message": "Content not found for the provided ID", "code": 404}
#     target = doc['target']
#     result_list = []
#     for t in target.split("\n"):
#         progress_result = {
#             "subdomain": ["", ""],
#             "subdomainTakeover": ["", ""],
#             "portScan": ["", ""],
#             "assetMapping": ["", ""],
#             "urlScan": ["", ""],
#             "sensitive": ["", ""],
#             "crawler": ["", ""],
#             "dirScan": ["", ""],
#             "vulnerability": ["", ""],
#             "all": ["", ""],
#             "target": ""
#         }
#         if not doc['subdomainScan']:
#             progress_result['subdomain'] = ['', '', '']
#         if not doc['portScan']:
#             progress_result['portScan'] = ['', '', '']
#         if not doc['urlScan']:
#             progress_result['urlScan'] = ['', '', '']
#         if not doc['sensitiveInfoScan']:
#             progress_result['sensitive'] = ['', '', '']
#         if not doc['crawlerScan']:
#             progress_result['crawler'] = ['', '', '']
#         if not doc['dirScan']:
#             progress_result['dirScan'] = ['', '', '']
#         if not doc['vulScan']:
#             progress_result['vulnerability'] = ['', '', '']
#         if runner != "":
#             key = "TaskInfo:progress:" + runner + ":" + t
#         else:
#             key = "TaskInfo:progress:" + task_id + ":" + t
#         data = await redis_con.hgetall(key)
#         progress_result["target"] = t
#         if not data:
#             result_list.append(progress_result)
#         else:
#             progress_result['subdomain'][0] = data.get("subdomain_start","")
#             progress_result['subdomain'][1] = data.get("subdomain_end", "")
#             progress_result['subdomainTakeover'][0] = data.get("subdomainTakeover_start", "")
#             progress_result['subdomainTakeover'][1] = data.get("subdomainTakeover_end", "")
#             progress_result['portScan'][0] = data.get("portScan_start", "")
#             progress_result['portScan'][1] = data.get("portScan_end", "")
#             progress_result['assetMapping'][0] = data.get("assetMapping_start", "")
#             progress_result['assetMapping'][1] = data.get("assetMapping_end", "")
#             progress_result['urlScan'][0] = data.get("urlScan_start", "")
#             progress_result['urlScan'][1] = data.get("urlScan_end", "")
#             progress_result['sensitive'][0] = data.get("sensitive_start", "")
#             progress_result['sensitive'][1] = data.get("sensitive_end", "")
#             progress_result['crawler'][0] = data.get("crawler_start", "")
#             progress_result['crawler'][1] = data.get("crawler_end", "")
#             progress_result['dirScan'][0] = data.get("dirScan_start", "")
#             progress_result['dirScan'][1] = data.get("dirScan_end", "")
#             progress_result['vulnerability'][0] = data.get("vulnerability_start", "")
#             progress_result['vulnerability'][1] = data.get("vulnerability_end", "")
#             progress_result['all'][0] = data.get("scan_start", "")
#             progress_result['all'][1] = data.get("scan_end", "")
#             result_list.append(progress_result)
#     return {
#         "code": 200,
#         "data": {
#             'list': result_list,
#             "total": len(result_list)
#         }
#     }
@router.post("/task/progress/info")
async def progress_info(request_data: dict, _: dict = Depends(verify_token), redis_con=Depends(get_redis_pool),
                        db=Depends(get_mongo_db)):
    task_id = request_data.get("id")
    type = request_data.get("type")
    runner = request_data.get("runner")

    if not task_id:
        return {"message": "ID is missing in the request data", "code": 400}

    query = {"_id": ObjectId(task_id)}

    if type == "scan":
        doc = await db.task.find_one(query)
    else:
        doc, target_data = await asyncio.gather(
            db.project.find_one(query),
            db.ProjectTargetData.find_one({"id": task_id})
        )
        if target_data:
            doc["target"] = target_data["target"]

    if not doc:
        return {"message": "Content not found for the provided ID", "code": 404}

    target = doc['target']
    result_list = []

    tasks = []
    for t in target.split("\n"):
        key = f"TaskInfo:progress:{runner or task_id}:{t}"
        tasks.append(redis_con.hgetall(key))

    redis_results = await asyncio.gather(*tasks)

    for t, data in zip(target.split("\n"), redis_results):
        progress_result = {
            "subdomain": ["", ""],
            "subdomainTakeover": ["", ""],
            "portScan": ["", ""],
            "assetMapping": ["", ""],
            "urlScan": ["", ""],
            "sensitive": ["", ""],
            "crawler": ["", ""],
            "dirScan": ["", ""],
            "vulnerability": ["", ""],
            "all": ["", ""],
            "target": t
        }

        if not data:
            result_list.append(progress_result)
            continue

        progress_result['subdomain'][0] = data.get("subdomain_start", "")
        progress_result['subdomain'][1] = data.get("subdomain_end", "")
        progress_result['subdomainTakeover'][0] = data.get("subdomainTakeover_start", "")
        progress_result['subdomainTakeover'][1] = data.get("subdomainTakeover_end", "")
        progress_result['portScan'][0] = data.get("portScan_start", "")
        progress_result['portScan'][1] = data.get("portScan_end", "")
        progress_result['assetMapping'][0] = data.get("assetMapping_start", "")
        progress_result['assetMapping'][1] = data.get("assetMapping_end", "")
        progress_result['urlScan'][0] = data.get("urlScan_start", "")
        progress_result['urlScan'][1] = data.get("urlScan_end", "")
        progress_result['sensitive'][0] = data.get("sensitive_start", "")
        progress_result['sensitive'][1] = data.get("sensitive_end", "")
        progress_result['crawler'][0] = data.get("crawler_start", "")
        progress_result['crawler'][1] = data.get("crawler_end", "")
        progress_result['dirScan'][0] = data.get("dirScan_start", "")
        progress_result['dirScan'][1] = data.get("dirScan_end", "")
        progress_result['vulnerability'][0] = data.get("vulnerability_start", "")
        progress_result['vulnerability'][1] = data.get("vulnerability_end", "")
        progress_result['all'][0] = data.get("scan_start", "")
        progress_result['all'][1] = data.get("scan_end", "")

        result_list.append(progress_result)

    return {
        "code": 200,
        "data": {
            'list': result_list,
            "total": len(result_list)
        }
    }

async def scheduler_scan_task(id):
    logger.info(f"Scheduler scan {id}")
    async for db in get_mongo_db():
        async for redis in get_redis_pool():
            next_time = scheduler.get_job(id).next_run_time
            formatted_time = next_time.strftime("%Y-%m-%d %H:%M:%S")
            doc = await db.ScheduledTasks.find_one({"id": id})
            run_id_last = doc.get("runner_id", "")
            if run_id_last != "" and id != run_id_last:
                progresskeys = await redis.keys(f"TaskInfo:progress:{run_id_last}:*")
                for pgk in progresskeys:
                    await redis.delete(pgk)
            task_id = generate_random_string(15)
            update_document = {
                "$set": {
                    "lastTime": get_now_time(),
                    "nextTime": formatted_time,
                    "runner_id": task_id
                }
            }
            await db.ScheduledTasks.update_one({"id": id}, update_document)
            query = {"_id": ObjectId(id)}
            doc = await db.task.find_one(query)
            targetList = []
            for t in doc['target'].split("\n"):
                t.replace("http://", "").replace("https://", "")
                t = t.strip("\n").strip("\r").strip()
                if t != "" and t not in targetList:
                    targetList.append(t)
            await create_scan_task(doc, task_id, targetList, redis)