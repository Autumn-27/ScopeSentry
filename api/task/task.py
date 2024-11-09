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
from api.task.util import task_progress, delete_asset, get_target_list
from api.users import verify_token
from motor.motor_asyncio import AsyncIOMotorCursor

from core.apscheduler_handler import scheduler
from core.redis_handler import get_redis_pool, refresh_config
from core.util import *

router = APIRouter()


@router.post("/data")
async def get_task_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token),
                        background_tasks: BackgroundTasks = BackgroundTasks(), redis_con=Depends(get_redis_pool)):
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
        cursor: AsyncIOMotorCursor = db.task.find(query).skip((page_index - 1) * page_size).limit(page_size).sort(
            [("creatTime", DESCENDING)])
        result = await cursor.to_list(length=None)
        # Process the result as needed
        response_data = [
            {"id": str(doc["_id"]),"status": doc["status"], "name": doc["name"], "taskNum": doc["taskNum"], "progress": doc["progress"],
             "creatTime": doc["creatTime"], "endTime": doc["endTime"]} for doc in result]

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
        return {"message": "error", "code": 500}


@router.post("/add")
async def add_task(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token),
                   redis_con=Depends(get_redis_pool)):
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
async def delete_task(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token),
                      redis_con=Depends(get_redis_pool), background_tasks: BackgroundTasks = BackgroundTasks()):
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
            await refresh_config("all", "delete_task", ",".join(task_ids))
            return {"code": 200, "message": "Task deleted successfully"}
        else:
            return {"code": 404, "message": "Task not found"}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/retest")
async def retest_task(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token),
                      redis_con=Depends(get_redis_pool)):
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
    page_index = request_data.get("pageIndex", 1)
    page_size = request_data.get("pageSize", 10)
    if not task_id:
        return {"message": "ID is missing in the request data", "code": 400}

    query = {"_id": ObjectId(task_id)}

    doc = await db.task.find_one(query)

    if not doc:
        return {"message": "Content not found for the provided ID", "code": 404}

    target = doc['target']
    # 计算起始和结束索引
    start_index = (page_index - 1) * page_size
    end_index = start_index + page_size
    data_array = target.split("\n")
    total = len(data_array)
    # 获取对应的数据
    paged_data = data_array[start_index:end_index]
    tg_key = {}
    redis_key = []
    for tg in paged_data:
        r = await get_target_list(tg, doc.get("ignore", ""))
        if len(r) == 0:
            tg_key[tg + "[ignore]"] = []
        else:
            tg_key[tg] = r
            redis_key.extend(r)

    tasks = []
    for t in redis_key:
        key = f"TaskInfo:progress:{task_id}:{t}"
        tasks.append(redis_con.hgetall(key))

    redis_results = await asyncio.gather(*tasks)
    redis_result_dict = dict(zip(redis_key, redis_results))
    result_list = []
    for t in tg_key:
        progress_result = {
            "_id": generate_random_string(6),
            "TargetHandler": ["", ""],
            "SubdomainScan": ["", ""],
            "SubdomainSecurity": ["", ""],
            "PortScanPreparation": ["", ""],
            "PortScan": ["", ""],
            "PortFingerprint": ["", ""],
            "AssetMapping": ["", ""],
            "AssetHandle": ["", ""],
            "URLScan": ["", ""],
            "WebCrawler": ["", ""],
            "URLSecurity": ["", ""],
            "DirScan": ["", ""],
            "VulnerabilityScan": ["", ""],
            "All": ["", ""],
            "target": t
        }
        if len(tg_key[t]) == 0:
            result_list.append(progress_result)
        else:
            if len(tg_key[t]) == 1:
                if tg_key[t][0] in redis_result_dict:
                    if redis_result_dict[tg_key[t][0]]:
                        progress_result['TargetHandler'][0] = redis_result_dict[tg_key[t][0]].get("TargetHandler_start", "")
                        progress_result['TargetHandler'][1] = redis_result_dict[tg_key[t][0]].get("TargetHandler_start", "")

                        progress_result['SubdomainScan'][0] = redis_result_dict[tg_key[t][0]].get("SubdomainScan_start", "")
                        progress_result['SubdomainScan'][1] = redis_result_dict[tg_key[t][0]].get("SubdomainScan_end", "")

                        progress_result['SubdomainSecurity'][0] = redis_result_dict[tg_key[t][0]].get("SubdomainSecurity_start", "")
                        progress_result['SubdomainSecurity'][1] = redis_result_dict[tg_key[t][0]].get("SubdomainSecurity_end", "")

                        progress_result['PortScanPreparation'][0] = redis_result_dict[tg_key[t][0]].get("PortScanPreparation_start", "")
                        progress_result['PortScanPreparation'][1] = redis_result_dict[tg_key[t][0]].get("PortScanPreparation_end", "")

                        progress_result['PortScan'][0] = redis_result_dict[tg_key[t][0]].get("PortScan_start", "")
                        progress_result['PortScan'][1] = redis_result_dict[tg_key[t][0]].get("PortScan_end", "")

                        progress_result['PortFingerprint'][0] = redis_result_dict[tg_key[t][0]].get("PortFingerprint_start", "")
                        progress_result['PortFingerprint'][1] = redis_result_dict[tg_key[t][0]].get("PortFingerprint_end", "")

                        progress_result['AssetMapping'][0] = redis_result_dict[tg_key[t][0]].get("AssetMapping_start", "")
                        progress_result['AssetMapping'][1] = redis_result_dict[tg_key[t][0]].get("AssetMapping_end", "")

                        progress_result['AssetHandle'][0] = redis_result_dict[tg_key[t][0]].get("AssetHandle_start", "")
                        progress_result['AssetHandle'][1] = redis_result_dict[tg_key[t][0]].get("AssetHandle_end", "")

                        progress_result['URLScan'][0] = redis_result_dict[tg_key[t][0]].get("URLScan_start", "")
                        progress_result['URLScan'][1] = redis_result_dict[tg_key[t][0]].get("URLScan_end", "")

                        progress_result['WebCrawler'][0] = redis_result_dict[tg_key[t][0]].get("WebCrawler_start", "")
                        progress_result['WebCrawler'][1] = redis_result_dict[tg_key[t][0]].get("WebCrawler_end", "")

                        progress_result['URLSecurity'][0] = redis_result_dict[tg_key[t][0]].get("URLSecurity_start", "")
                        progress_result['URLSecurity'][1] = redis_result_dict[tg_key[t][0]].get("URLSecurity_end", "")

                        progress_result['DirScan'][0] = redis_result_dict[tg_key[t][0]].get("DirScan_start", "")
                        progress_result['DirScan'][1] = redis_result_dict[tg_key[t][0]].get("DirScan_end", "")

                        progress_result['VulnerabilityScan'][0] = redis_result_dict[tg_key[t][0]].get("VulnerabilityScan_start", "")
                        progress_result['VulnerabilityScan'][1] = redis_result_dict[tg_key[t][0]].get("VulnerabilityScan_end", "")

                        progress_result['All'][0] = redis_result_dict[tg_key[t][0]].get("scan_start", "")
                        progress_result['All'][1] = redis_result_dict[tg_key[t][0]].get("scan_end", "")
            else:
                progress_result["children"] = []
                for son_target in tg_key[t]:
                    tmp = {
                        "_id": generate_random_string(6),
                        "TargetHandler": ["", ""],
                        "SubdomainScan": ["", ""],
                        "SubdomainSecurity": ["", ""],
                        "PortScanPreparation": ["", ""],
                        "PortScan": ["", ""],
                        "PortFingerprint": ["", ""],
                        "AssetMapping": ["", ""],
                        "AssetHandle": ["", ""],
                        "URLScan": ["", ""],
                        "WebCrawler": ["", ""],
                        "URLSecurity": ["", ""],
                        "DirScan": ["", ""],
                        "VulnerabilityScan": ["", ""],
                        "All": ["", ""],
                        "target": son_target
                    }
                    tmp['TargetHandler'][0] = redis_result_dict[son_target].get("TargetHandler_start", "")
                    tmp['TargetHandler'][1] = redis_result_dict[son_target].get("TargetHandler_start", "")

                    tmp['SubdomainScan'][0] = redis_result_dict[son_target].get("SubdomainScan_start", "")
                    tmp['SubdomainScan'][1] = redis_result_dict[son_target].get("SubdomainScan_end", "")

                    tmp['SubdomainSecurity'][0] = redis_result_dict[son_target].get("SubdomainSecurity_start", "")
                    tmp['SubdomainSecurity'][1] = redis_result_dict[son_target].get("SubdomainSecurity_end", "")

                    tmp['PortScanPreparation'][0] = redis_result_dict[son_target].get("PortScanPreparation_start", "")
                    tmp['PortScanPreparation'][1] = redis_result_dict[son_target].get("PortScanPreparation_end", "")

                    tmp['PortScan'][0] = redis_result_dict[son_target].get("PortScan_start", "")
                    tmp['PortScan'][1] = redis_result_dict[son_target].get("PortScan_end", "")

                    tmp['PortFingerprint'][0] = redis_result_dict[son_target].get("PortFingerprint_start", "")
                    tmp['PortFingerprint'][1] = redis_result_dict[son_target].get("PortFingerprint_end", "")

                    tmp['AssetMapping'][0] = redis_result_dict[son_target].get("AssetMapping_start", "")
                    tmp['AssetMapping'][1] = redis_result_dict[son_target].get("AssetMapping_end", "")

                    tmp['AssetHandle'][0] = redis_result_dict[son_target].get("AssetHandle_start", "")
                    tmp['AssetHandle'][1] = redis_result_dict[son_target].get("AssetHandle_end", "")

                    tmp['URLScan'][0] = redis_result_dict[son_target].get("URLScan_start", "")
                    tmp['URLScan'][1] = redis_result_dict[son_target].get("URLScan_end", "")

                    tmp['WebCrawler'][0] = redis_result_dict[son_target].get("WebCrawler_start", "")
                    tmp['WebCrawler'][1] = redis_result_dict[son_target].get("WebCrawler_end", "")

                    tmp['URLSecurity'][0] = redis_result_dict[son_target].get("URLSecurity_start", "")
                    tmp['URLSecurity'][1] = redis_result_dict[son_target].get("URLSecurity_end", "")

                    tmp['DirScan'][0] = redis_result_dict[son_target].get("DirScan_start", "")
                    tmp['DirScan'][1] = redis_result_dict[son_target].get("DirScan_end", "")

                    tmp['VulnerabilityScan'][0] = redis_result_dict[son_target].get("VulnerabilityScan_start","")
                    tmp['VulnerabilityScan'][1] = redis_result_dict[son_target].get("VulnerabilityScan_end","")

                    tmp['All'][0] = redis_result_dict[son_target].get("scan_start", "")
                    tmp['All'][1] = redis_result_dict[son_target].get("scan_end", "")
                    progress_result["children"].append(tmp)

            result_list.append(progress_result)

    return {
        "code": 200,
        "data": {
            'list': result_list,
            "total": total
        }
    }


@router.post("/stop")
async def stop_task(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # Get the ID from the request data
        task_id = request_data.get("id")
        await refresh_config("all", "stop_task", task_id)
        await db.task.update_one({"_id": ObjectId(task_id)}, {"$set": {"status": 2}})
        return {"message": "success", "code": 200}
    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/start")
async def start_task(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # Get the ID from the request data
        task_id = request_data.get("id")
        query = {"_id": ObjectId(task_id)}
        doc = await db.task.find_one(query)
        if not doc:
            return {"message": "Content not found for the provided ID", "code": 404}
        doc["type"] = "start"
        await create_scan_task(doc, task_id, True)
        await db.task.update_one({"_id": ObjectId(task_id)}, {"$set": {"status": 1}})
        return {"message": "success", "code": 200}
    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}
