# -------------------------------------
# @file      : handler.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/10/28 22:09
# -------------------------------------------
import asyncio
import json

from bson import ObjectId
from motor.motor_asyncio import AsyncIOMotorCursor
from pymongo import DESCENDING

from api.asset.page_monitoring import get_page_monitoring_data
from api.node.handler import get_node_all
from core.apscheduler_handler import scheduler
from api.task.util import get_target_list
from core.db import get_mongo_db
from core.handler.task import get_task_data
from core.redis_handler import get_redis_pool, get_redis_online_data
from core.util import get_now_time, get_search_query
from loguru import logger

running_tasks = set()


async def insert_task(request_data, db):
    # 解析多种来源设置target
    targetSource = request_data.get("targetSource", "general")
    if "Source" in targetSource:
        # 如果是从资产处选则数据进行创建任务
        index = targetSource.replace("Source", "")
        targetTp = request_data.get("targetTp")
        if targetTp == "search":
            # 如果是按照当前搜索条件进行搜索
            targetNumber = int(request_data.get("targetNumber", 0))
            if targetNumber == 0:
                return {"code": 400, "message": "targetNumber is 0"}
            query = await get_search_query(index, request_data)
            target = await get_target_search(query, targetNumber, index, db)
            request_data["target"] = target
        else:
            # 按照选择的数据进行创建任务
            targetIds = request_data.get("targetIds", [])
            if len(targetIds) == 0:
                return {"code": 400, "message": "targetIds is null"}
            target = await get_target_ids(targetIds, index, db)
            request_data["target"] = target
    elif targetSource == "general":
        # 普通创建
        target = request_data.get("target", "")
    elif targetSource == "project":
        # 从项目创建
        project_ids = request_data.get("project", [])
        if len(project_ids) == 0:
            return {"code": 400, "message": "project is null"}
        target = await get_target_project(project_ids, db)
        request_data["target"] = target
        request_data["filter"] = {"project": project_ids}
        if len(project_ids) == 1:
            request_data["bindProject"] = project_ids[0]
    else:
        project_ids = request_data.get("project", [])
        request_data["filter"] = {"project": project_ids}
        query = await get_search_query(targetSource, request_data)
        target = await get_target_search(query, 0, targetSource, db)
        request_data["target"] = target
    targetList = await get_target_list(request_data['target'], request_data.get("ignore", ""))
    taskNum = len(targetList)
    if "_id" in request_data:
        del request_data["_id"]
    request_data['taskNum'] = taskNum
    request_data['target'] = "\n".join(targetList)
    request_data['progress'] = 0
    request_data["creatTime"] = get_now_time()
    request_data["endTime"] = ""
    request_data["status"] = 1
    request_data["type"] = request_data.get("targetSource", "scan")
    result = await db.task.insert_one(request_data)
    if result.inserted_id:
        task = asyncio.create_task(create_scan_task(request_data, str(result.inserted_id)))
        running_tasks.add(task)
        task.add_done_callback(lambda t: running_tasks.remove(t))
        return result.inserted_id


async def create_scan_task(request_data, id, stop_to_start=False):
    logger.info(f"[create_scan_task] begin: {id}")
    async for db in get_mongo_db():
        async for redis_con in get_redis_pool():
            request_data["id"] = str(id)
            if request_data['allNode']:
                all_node = await get_node_all(redis_con)
                for node in all_node:
                    if node not in request_data["node"]:
                        request_data["node"].append(node)

            # 如果是暂停之后重新开始的，则不需要删除缓存和填入目标
            if stop_to_start is False:
                # 删除可能存在缓存
                keys_to_delete = [
                    f"TaskInfo:tmp:{id}",
                    f"TaskInfo:{id}",
                    f"TaskInfo:time:{id}",
                ]
                progresskeys = await redis_con.keys(f"TaskInfo:progress:{id}:*")
                keys_to_delete.extend(progresskeys)
                progresskeys = await redis_con.keys(f"duplicates:{id}:*")
                keys_to_delete.extend(progresskeys)
                await redis_con.delete(*keys_to_delete)
                # 原始的target生成target list
                target_list = await get_target_list(request_data['target'], request_data.get("ignore", ""))
                # 将任务目标插入redis中
                await redis_con.lpush(f"TaskInfo:{id}", *target_list)
            # 获取模板数据
            template_data = await get_task_data(db, request_data, id)
            # 分发任务
            for name in request_data["node"]:
                await redis_con.rpush(f"NodeTask:{name}", json.dumps(template_data))
            logger.info(f"[create_scan_task] end: {id}")
            return True


async def get_target_project(ids, db):
    cursor: AsyncIOMotorCursor = db.ProjectTargetData.find({"id": {"$in": ids}})
    targets = ""
    async for doc in cursor:
        targets += doc.get("target", "").strip() + "\n"
    return targets.strip()


async def get_target_search(query, number, index, db):
    displayKey = {
        'subdomain': {
            'host': 1,
        },
        'asset': {
            'url': 1,
            'host': 1,
            'port': 1,
            'service': 1,
            'type': 1,
        },
        'UrlScan': {
            'output': 1,
        },
        'RootDomain': {
            'domain': 1
        }
    }
    if index not in displayKey:
        return ""
    if number == 0:
        cursor: AsyncIOMotorCursor = db[index].find(query, displayKey[index])
    else:
        cursor: AsyncIOMotorCursor = db[index].find(query, displayKey[index]).limit(number).sort([("time", DESCENDING)])
    target = ""
    async for doc in cursor:
        if index == "asset":
            if doc["type"] == "http":
                target += doc.get("url", "") + "\n"
            else:
                target += doc.get("service", "http") + "://" + doc["host"] + ":" + str(doc["port"]) + "\n"
        elif index == "subdomain":
            target += doc.get("host", "") + "\n"
        elif index == "UrlScan":
            target += doc.get("output", "") + "\n"
        elif index == "RootDomain":
            target += doc.get("domain", "") + "\n"
    return target


async def get_target_ids(ids, index, db):
    key = ["asset", "UrlScan", "subdomain"]
    if index not in key:
        return {"code": 404, "message": "Data not found"}
    obj_ids = []
    for data_id in ids:
        obj_ids.append(ObjectId(data_id))
    cursor = db[index].find({"_id": {"$in": obj_ids}})
    target = ''
    async for doc in cursor:
        if index == "asset":
            if doc["type"] == "http":
                target += doc.get("url", "") + "\n"
            else:
                target += doc.get("service", "http") + "://" + doc["host"] + ":" + str(doc["port"]) + "\n"
        elif index == "subdomain":
            target += doc.get("host", "") + "\n"
        elif index == "UrlScan":
            target += doc.get("output", "") + "\n"
    return target


async def scheduler_scan_task(id, tp):
    logger.info(f"Scheduler scan {id}")
    async for db in get_mongo_db():
        next_time = scheduler.get_job(id).next_run_time
        formatted_time = next_time.strftime("%Y-%m-%d %H:%M:%S")
        time_now = get_now_time()
        update_document = {
            "$set": {
                "lastTime": time_now,
                "nextTime": formatted_time
            }
        }
        await db.ScheduledTasks.update_one({"id": id}, update_document)
        doc = await db.ScheduledTasks.find_one({"id": id})
        doc["name"] = doc["name"] + f"-{doc.get('targetSource', 'None')}-" + time_now
        await insert_task(doc, db)


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
                "ID": 'page_monitoring',
                "type": "page_monitoring"
            }
            for name in name_list:
                await redis.rpush(f"NodeTask:{name}", json.dumps(add_redis_task_data))


async def insert_scheduled_tasks(request_data, db, update=False, id=""):
    cycle_type = request_data['cycleType']
    if cycle_type == "":
        return
    task_id = ""
    if update is False:
        result = await db.ScheduledTasks.insert_one(request_data)
        if result.inserted_id:
            task_id = str(result.inserted_id)
        else:
            return
    else:
        task_id = id
    week = request_data.get("week", 1)
    day = int(request_data.get("day", 1))
    hour = int(request_data.get("hour", 0))
    minute = int(request_data.get("minute", 0))
    if cycle_type == "daily":
        # 每天固定时间执行
        scheduler.add_job(
            scheduler_scan_task, 'cron',
            hour=hour, minute=minute,
            args=[str(task_id), "scan"],
            id=task_id, jobstore='mongo'
        )
    elif cycle_type == "ndays":
        # 每 N 天执行一次
        scheduler.add_job(
            scheduler_scan_task, 'interval',
            days=day, hours=hour, minutes=minute,
            args=[str(task_id), "scan"],
            id=task_id, jobstore='mongo'
        )

    elif cycle_type == "nhours":
        # 每 N 小时执行一次
        scheduler.add_job(
            scheduler_scan_task, 'interval',
            hours=hour, minutes=minute,
            args=[str(task_id), "scan"],
            id=task_id, jobstore='mongo'
        )

    elif cycle_type == "weekly":
        # 每星期几执行一次
        scheduler.add_job(
            scheduler_scan_task, 'cron',
            day_of_week=week,
            hour=hour, minute=minute,
            args=[str(task_id), "scan"],
            id=task_id, jobstore='mongo'
        )

    elif cycle_type == "monthly":
        # 每月第几天固定时间执行
        scheduler.add_job(
            scheduler_scan_task, 'cron',
            day=day, hour=hour, minute=minute,
            args=[str(task_id), "scan"],
            id=task_id, jobstore='mongo'
        )
    next_time = scheduler.get_job(str(task_id)).next_run_time
    formatted_time = next_time.strftime("%Y-%m-%d %H:%M:%S")
    update_document = {
        "$set": {
            "lastTime": "",
            "nextTime": formatted_time,
            "id": str(task_id)
        }
    }
    await db.ScheduledTasks.update_one({"_id": ObjectId(task_id)}, update_document)
    return
