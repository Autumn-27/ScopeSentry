# -------------------------------------
# @file      : scheduled_tasks.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/4/28 20:58
# -------------------------------------------
import traceback

from apscheduler.events import JobSubmissionEvent, EVENT_JOB_MAX_INSTANCES, EVENT_JOB_SUBMITTED
from apscheduler.executors.base import MaxInstancesReachedError
from bson import ObjectId
from fastapi import APIRouter, Depends
from pytz import utc

from api.task.handler import scheduler_scan_task
from api.users import verify_token
from motor.motor_asyncio import AsyncIOMotorCursor

from core.apscheduler_handler import scheduler
from core.db import get_mongo_db
from core.redis_handler import get_redis_pool
from core.util import *
from api.node import get_redis_online_data
from api.page_monitoring import get_page_monitoring_data

router = APIRouter()


@router.post("/scheduled/task/data")
async def get_scheduled_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        search_query = request_data.get("search", "")
        page_index = request_data.get("pageIndex", 1)
        page_size = request_data.get("pageSize", 10)
        # Fuzzy search based on the name field
        query = {"name": {"$regex": search_query, "$options": "i"}}

        # Get the total count of documents matching the search criteria
        total_count = await db.ScheduledTasks.count_documents(query)

        # Perform pagination query
        cursor: AsyncIOMotorCursor = db.ScheduledTasks.find(query).skip((page_index - 1) * page_size).limit(page_size)
        result = await cursor.to_list(length=None)
        if len(result) == 0:
            return {
                "code": 200,
                "data": {
                    'list': [],
                    'total': 0
                }
            }
        result_list = []
        for doc in result:
            tmp = {
                "id": doc["id"],
                "name": doc["name"],
                "type": doc["type"],
                "lastTime": doc.get("lastTime", ""),
                "nextTime": doc.get("nextTime", ""),
                "state": doc.get("state"),
                "cycle": doc.get("hour"),
                "node": doc.get("node", []),
                "allNode": doc.get("allNode", True),
                "runner_id": doc.get("runner_id", "")
            }
            result_list.append(tmp)
        return {
            "code": 200,
            "data": {
                'list': result_list,
                'total': total_count
            }
        }

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}

# @router.post("/scheduled/task/run")
# async def scheduled_run(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token),
#                         jobstore_alias=None):
#     try:
#         id = request_data.get("id", "")
#         job = scheduler.get_job(id)
#         if job:
#             executor = scheduler._lookup_executor(job.executor)
#             run_times = [datetime.now(utc)]
#             try:
#                 executor.submit_job(job, run_times)
#             except MaxInstancesReachedError:
#                 scheduler._logger.warning(
#                     'Execution of job "%s" skipped: maximum number of running '
#                     'instances reached (%d)', job, job.max_instances)
#                 event = JobSubmissionEvent(EVENT_JOB_MAX_INSTANCES, job.id,
#                                            jobstore_alias, run_times)
#                 scheduler._dispatch_event(event)
#             except BaseException:
#                 scheduler._logger.exception('Error submitting job "%s" to executor "%s"',
#                                             job, job.executor)
#             else:
#                 event = JobSubmissionEvent(EVENT_JOB_SUBMITTED, job.id, jobstore_alias,
#                                            run_times)
#                 scheduler._dispatch_event(event)
#             return {"message": "task run success", "code": 200}
#         else:
#             return {"message": "Not Found Task", "code": 500}
#     except:
#         return {"message": "error", "code": 500}
@router.post("/scheduled/task/delete")
async def delete_task(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token),
                      redis_con=Depends(get_redis_pool)):
    try:
        # Extract the list of IDs from the request_data dictionary
        task_ids = request_data.get("ids", [])

        # Convert the provided rule_ids to ObjectId
        obj_ids = []
        for task_id in task_ids:
            if task_id != "page_monitoring":
                obj_ids.append(task_id)
                job = scheduler.get_job(task_id)
                if job:
                    function_name = job.func.__name__ if hasattr(job.func, '__name__') else job.func
                    update_document = {
                        "$set": {
                            "scheduledTasks": False
                        }
                    }
                    if function_name == "scheduler_scan_task":
                        await db.task.update_one({"_id": ObjectId(task_id)}, update_document)
                    else:
                        await db.project.update_one({"_id": ObjectId(task_id)}, update_document)
                    scheduler.remove_job(task_id)
        result = await db.ScheduledTasks.delete_many({"id": {"$in": obj_ids}})

        # Check if the deletion was successful
        if result.deleted_count > 0:
            return {"code": 200, "message": "Scheduled Task deleted successfully"}
        else:
            return {"code": 404, "message": "Scheduled Task not found"}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


async def get_page_monitoring_time():
    async for db in get_mongo_db():
        result = await db.ScheduledTasks.find_one({"id": "page_monitoring"})
        time = result['hour']
        flag = result['state']
        return time, flag


async def create_page_monitoring_task():
    logger.info("create_page_monitoring_task")
    async for db in get_mongo_db():
        async for redis in get_redis_pool():
            name_list = []
            result = await db.ScheduledTasks.find_one({"id": "page_monitoring"})
            next_time = scheduler.get_job("page_monitoring").next_run_time
            formatted_time = next_time.strftime("%Y-%m-%d %H:%M:%S")
            update_document = {
                "$set": {
                    "lastTime": get_now_time(),
                    "nextTime": formatted_time
                }
            }
            await db.ScheduledTasks.update_one({"_id": result['_id']}, update_document)
            if result['allNode']:
                tmp = await get_redis_online_data(redis)
                name_list += tmp
            else:
                name_list += result['node']
            targetList = await get_page_monitoring_data(db, False)
            if len(targetList) == 0:
                return
            await redis.delete(f"TaskInfo:page_monitoring")
            await redis.lpush(f"TaskInfo:page_monitoring", *targetList)
            add_redis_task_data = {
                "type": 'page_monitoring',
                "TaskId": "page_monitoring"
            }
            for name in name_list:
                await redis.rpush(f"NodeTask:{name}", json.dumps(add_redis_task_data))


@router.post("/pagemonit/data")
async def get_scheduled_task_pagemonit_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        search_query = request_data.get("search", "")
        page_index = request_data.get("pageIndex", 1)
        page_size = request_data.get("pageSize", 10)
        # Fuzzy search based on the name field
        query = {"url": {"$regex": search_query, "$options": "i"}}

        # Get the total count of documents matching the search criteria
        total_count = await db.PageMonitoring.count_documents(query)

        # Perform pagination query
        cursor: AsyncIOMotorCursor = db.PageMonitoring.find(query).skip((page_index - 1) * page_size).limit(page_size)
        result = await cursor.to_list(length=None)
        result_list = []
        for doc in result:
            tmp = {
                "id": str(doc["_id"]),
                "url": doc["url"]
            }
            result_list.append(tmp)
        return {
            "code": 200,
            "data": {
                'list': result_list,
                'total': total_count
            }
        }

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/pagemonit/update")
async def update_scheduled_task_pagemonit_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        if not request_data:
            return {"message": "Data to update is missing in the request", "code": 400}
        state = request_data.get('state')
        formatted_time = ""
        job = scheduler.get_job('page_monitoring')
        if state:
            if job is None:
                scheduler.add_job(create_page_monitoring_task, 'interval', hours=request_data.get('hour', 24), id='page_monitoring', jobstore='mongo')
                next_time = scheduler.get_job('page_monitoring').next_run_time
                formatted_time = next_time.strftime("%Y-%m-%d %H:%M:%S")
        else:
            if job:
                scheduler.remove_job('page_monitoring')
        update_document = {
            "$set": {
                "hour": request_data.get('hour', 24),
                "node": request_data.get('node', []),
                "allNode": request_data.get('allNode', True),
                "nextTime": formatted_time,
                "state": request_data.get('state'),
            }
        }
        result = await db.ScheduledTasks.update_one({"id": 'page_monitoring'}, update_document)
        if result:
            return {"message": "Data updated successfully", "code": 200}
        else:
            return {"message": "Failed to update data", "code": 404}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/pagemonit/delete")
async def delete_scheduled_task_pagemonit_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # Extract the list of IDs from the request_data dictionary
        url_ids = request_data.get("ids", [])

        # Convert the provided rule_ids to ObjectId
        obj_ids = [ObjectId(url_id) for url_id in url_ids]

        # Delete the SensitiveRule documents based on the provided IDs
        result = await db.PageMonitoring.delete_many({"_id": {"$in": obj_ids}})

        # Check if the deletion was successful
        if result.deleted_count > 0:
            return {"code": 200, "message": "URL deleted successfully"}
        else:
            return {"code": 404, "message": "URL not found"}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/pagemonit/add")
async def add_scheduled_task_pagemonit_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        if not request_data:
            return {"message": "Data to add is missing in the request", "code": 400}
        url = request_data.get("url")
        result = await db.PageMonitoring.insert_one({
            "url": url,
            "content": [],
            "hash": [],
            "diff": [],
            "state": 1,
            "project": '',
            "time": ''
        })

        if result.inserted_id:
            return {"message": "Data added successfully", "code": 200}
        else:
            return {"message": "Failed to add data", "code": 400}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/update")
async def update_task_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # Get the ID from the request data
        task_id = request_data.get("id")
        hour = request_data.get("hour")
        # Check if ID is provided
        if not task_id:
            return {"message": "ID is missing in the request data", "code": 400}
        query = {"id": ObjectId(task_id)}
        doc = await db.ScheduledTasks.find_one(query)
        oldScheduledTasks = doc["state"]
        old_hour = doc["hour"]
        newScheduledTasks = request_data.get("scheduledTasks")
        if oldScheduledTasks is False:
            # 原本是关闭的
            # 如果重新开启，则开启计划任务
            if newScheduledTasks:
                scheduler.add_job(scheduler_scan_task, 'interval', hours=hour, args=[task_id], id=str(task_id),
                                  jobstore='mongo')
                await db.ScheduledTasks.update_one({"id": task_id}, {"$set": {'state': True}})
            # 如果新的还是关闭，则不用管
        else:
            # 原本是开启的
            # 如果重新提交的是开启的
            if newScheduledTasks:
                # 如果旧的时间间隔和新的时间剑客不同删除旧的然后按照新的时间重新提交，否则无需改变
                if old_hour != hour:
                    scheduler.remove_job(task_id)
                    scheduler.add_job(scheduler_scan_task, 'interval', hours=hour, args=[task_id], id=str(task_id),
                                      jobstore='mongo')
            else:
                # 重新提交的是关闭的，则移除计划任务
                scheduler.remove_job(task_id)
        request_data.pop("id")
        update_document = {
            "$set": request_data
        }
        result = await db.ScheduledTasks.update_one({"id": task_id}, update_document)
        if result:
            return {"message": "Task updated successfully", "code": 200}
        else:
            return {"message": "Failed to update data", "code": 404}

    except Exception as e:
        logger.error(str(e))
        logger.error(traceback.format_exc())
        return {"message": "error", "code": 500}