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

from api.task.handler import scheduler_scan_task, create_page_monitoring_task, insert_scheduled_tasks
from api.users import verify_token
from motor.motor_asyncio import AsyncIOMotorCursor

from core.apscheduler_handler import scheduler
from core.db import get_mongo_db
from core.redis_handler import get_redis_pool
from core.util import *

router = APIRouter()


@router.post("/data")
async def get_scheduled_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
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
        week = doc.get("week", 1)
        day = doc.get("day", 1)
        hour = doc.get("hour", 0)
        minute = doc.get("minute", 0)
        cycle_type = doc.get("cycleType", "nhours")
        cycle = ""
        doc_id = str(doc["_id"])
        id = doc.get("id", doc_id)
        if id == "page_monitoring":
            cycle = f'{hour} hour'
        else:
            if cycle_type == "daily":
                cycle = f'Every day at {hour}:{minute}'
            elif cycle_type == "ndays":
                cycle = f'Every {day} days at {hour}:{minute}'
            elif cycle_type == "nhours":
                cycle = f'Every {hour}h {minute}m'
            elif cycle_type == "weekly":
                cycle = f'Every week on day {week} at {hour}:{minute}'
            elif cycle_type == "monthly":
                cycle = f'Every month on day {day} at {hour}:{minute}'
        tmp = {
            "id": id,
            "name": doc["name"],
            "type": doc.get("type", "scan"),
            "lastTime": doc.get("lastTime", ""),
            "nextTime": doc.get("nextTime", ""),
            "state": doc.get("scheduledTasks", doc.get("state", False)),
            "node": doc.get("node", []),
            "cycle": cycle,
            "allNode": doc.get("allNode", True),
            "runner_id": doc.get("runner_id", ""),
            "project": doc.get("project", []),
            "targetSource": doc.get("targetSource", "general"),
            "day": doc.get("day", 1),
            "minute": doc.get("minute", 1),
            "hour": doc.get("hour", 1),
            "search": doc.get("search", ""),
            "cycleType": cycle_type,
            "scheduledTasks": doc.get("scheduledTasks", True)
        }
        result_list.append(tmp)
    return {
        "code": 200,
        "data": {
            'list': result_list,
            'total': total_count
        }
    }



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


@router.post("/delete")
async def delete_task(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token),
                      redis_con=Depends(get_redis_pool)):
    try:
        # Extract the list of IDs from the request_data dictionary
        task_ids = request_data.get("ids", [])

        # Convert the provided rule_ids to ObjectId
        obj_ids = []
        sched_ids = []
        for task_id in task_ids:
            if task_id != "page_monitoring":
                obj_ids.append(task_id)
                try:
                    sched_ids.append(ObjectId(task_id))
                except:
                    logger.warning(f"Invalid task id {task_id}")
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
        result = await db.ScheduledTasks.delete_many({
            "$or": [
                {"id": {"$in": obj_ids}},
                {"_id": {"$in": sched_ids}}
            ]
        })

        # Check if the deletion was successful
        if result.deleted_count > 0:
            return {"code": 200, "message": "Scheduled Task deleted successfully"}
        else:
            return {"code": 404, "message": "Scheduled Task not found"}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/pagemonit/data")
async def get_scheduled_task_pagemonit_data(request_data: dict, db=Depends(get_mongo_db),
                                            _: dict = Depends(verify_token)):
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
async def update_scheduled_task_pagemonit_data(request_data: dict, db=Depends(get_mongo_db),
                                               _: dict = Depends(verify_token)):
    try:
        if not request_data:
            return {"message": "Data to update is missing in the request", "code": 400}
        state = request_data.get('state')
        formatted_time = ""
        job = scheduler.get_job('page_monitoring')
        if state:
            if job:
                scheduler.remove_job('page_monitoring')
            scheduler.add_job(create_page_monitoring_task, 'interval', hours=request_data.get('hour', 24),
                              id='page_monitoring', jobstore='mongo')
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
async def delete_scheduled_task_pagemonit_data(request_data: dict, db=Depends(get_mongo_db),
                                               _: dict = Depends(verify_token)):
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
async def add_scheduled_task_pagemonit_data(request_data: dict, db=Depends(get_mongo_db),
                                            _: dict = Depends(verify_token)):
    try:
        if not request_data:
            return {"message": "Data to add is missing in the request", "code": 400}
        url = request_data.get("url")
        result = await db.PageMonitoring.insert_one({
            "url": url,
            "hash": [],
            "md5": calculate_md5_from_content(url),
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
async def update_scheduled_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # Get the ID from the request data
        task_id = request_data.get("id")
        if task_id == "":
            return {"message": "id is none", "code": 400}
        scheduled_tasks = request_data.get("scheduledTasks", False)
        if scheduled_tasks:
            job = scheduler.get_job(task_id)
            if job is not None:
                scheduler.remove_job(task_id)
            await insert_scheduled_tasks(request_data, db, True, task_id)
        else:
            job = scheduler.get_job(task_id)
            if job is not None:
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


@router.post("/add")
async def add_scheduled_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        name = request_data.get("name")
        cursor = db.ScheduledTasks.find({"name": {"$eq": name}}, {"_id": 1})
        results = await cursor.to_list(length=None)
        if len(results) != 0:
            return {"code": 400, "message": "name already exists"}
        node = request_data.get("node")
        if name == "" or node == []:
            return {"message": "target is Null", "code": 500}

        scheduled_tasks = request_data.get("scheduledTasks", False)
        if scheduled_tasks:
            await insert_scheduled_tasks(request_data, db)
        else:
            await db.ScheduledTasks.insert_one(request_data)
        return {"message": "Task updated successfully", "code": 200}

    except Exception as e:
        logger.error(str(e))
        logger.error(traceback.format_exc())
        return {"message": "error", "code": 500}


@router.post("/detail")
async def scheduled_detail(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # Get the ID from the request data
        task_id = request_data.get("id")

        # Check if ID is provided
        if not task_id:
            return {"message": "ID is missing in the request data", "code": 400}

        # Query the database for content based on ID
        query = {"$or": [
        {"id": task_id},
        {"_id": ObjectId(task_id)}  # 你可以替换为你想要的条件
    ]}
        doc = await db.ScheduledTasks.find_one(query)
        if not doc:
            return {"message": "Content not found for the provided ID", "code": 404}
        result = {
            "name": doc.get("name", ""),
            "target": doc.get("target", ""),
            "ignore": doc.get("ignore", ""),
            "node": doc.get("node", []),
            "allNode": doc.get("allNode"),
            "hour": doc.get("hour"),
            "duplicates": doc.get("duplicates"),
            "template": doc.get("template", ""),
            "project": doc.get("project", []),
            "targetSource": doc.get("targetSource", "general"),
            "day": doc.get("day", 1),
            "minute": doc.get("minute", 1),
            "search": doc.get("search", ""),
            "cycleType": doc.get("cycleType", "nhours"),
            "scheduledTasks": doc.get("scheduledTasks", True)
        }
        return {"code": 200, "data": result}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}
