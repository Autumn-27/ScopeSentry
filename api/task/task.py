# -*- coding:utf-8 -*-　　
# @name: sensitive
# @auth: rainy-autumn@outlook.com
# @version:
import asyncio
import traceback

from bson import ObjectId
from fastapi import APIRouter, Depends, BackgroundTasks
from pymongo import DESCENDING

from api.task.handler import insert_task, scheduler_scan_task, create_scan_task
from api.task.util import task_progress, delete_asset
from api.users import verify_token
from motor.motor_asyncio import AsyncIOMotorCursor

from core.apscheduler_handler import scheduler
from core.redis_handler import get_redis_pool
from core.util import *
router = APIRouter()


@router.post("/data")
async def get_task_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token), background_tasks: BackgroundTasks = BackgroundTasks(), redis_con=Depends(get_redis_pool)):
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


@router.post("/add")
async def add_task(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token), redis_con=Depends(get_redis_pool)):
    try:
        name = request_data.get("name")
        cursor = db.task.find({"name": {"$eq": name}}, {"_id": 1})
        results = await cursor.to_list(length=None)
        if len(results) != 0:
            return {"code": 400, "message": "name already exists"}
        target = request_data.get("target", "")
        node = request_data.get("node")
        if name == "" or target == "" or node == []:
            return {"message": "Null", "code": 500}
        scheduledTasks = request_data.get("scheduledTasks", False)
        hour = request_data.get("hour", 24)
        task_id = await insert_task(request_data, db)
        if task_id:
            if scheduledTasks:
                scheduler.add_job(scheduler_scan_task, 'interval', hours=hour, args=[str(task_id), "scan"],
                                  id=str(task_id), jobstore='mongo')
                next_time = scheduler.get_job(str(task_id)).next_run_time
                formatted_time = next_time.strftime("%Y-%m-%d %H:%M:%S")
                # 插入计划任务管理
                request_data["type"] = "scan"
                request_data["state"] = True
                request_data["lastTime"] = ""
                request_data["nextTime"] = formatted_time
                request_data["id"] = str(task_id)
                await db.ScheduledTasks.insert_one(request_data)
            return {"code": 200, "message": "Task added successfully"}
        else:
            return {"code": 400, "message": "Failed to add Task"}
    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/detail")
async def task_detail(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
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
            "ignore": doc.get("ignore", ""),
            "node": doc.get("node", []),
            "allNode": doc.get("allNode"),
            "scheduledTasks": doc.get("scheduledTasks"),
            "hour": doc.get("hour"),
            "duplicates": doc.get("duplicates"),
            "template": doc.get("template", "")
        }
        return {"code": 200, "data": result}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/delete")
async def delete_task(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token), redis_con=Depends(get_redis_pool), background_tasks: BackgroundTasks = BackgroundTasks()):
    try:
        # Extract the list of IDs from the request_data dictionary
        task_ids = request_data.get("ids", [])
        delA = request_data.get("delA", False)

        # Convert the provided rule_ids to ObjectId
        obj_ids = []
        redis_key = []
        for task_id in task_ids:
            obj_ids.append(ObjectId(task_id))
            redis_key.append("TaskInfo:" + task_id)
            job = scheduler.get_job(task_id)
            if job:
                scheduler.remove_job(task_id)
        # 删除redis中的任务
        await redis_con.delete(*redis_key)

        # 删除计划任务
        await db.ScheduledTasks.delete_many({"id": {"$in": task_ids}})
        if delA:
            # 如果选择了删除资产，则删除资产
            task_name = await db.task.find({"_id": {"$in": obj_ids}}, {"name": 1})
            background_tasks.add_task(delete_asset, task_name, False)

        # 删除mongdob中的任务
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


@router.post("/retest")
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

        f = await create_scan_task(doc, task_id)

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


@router.post("/progress/info")
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


