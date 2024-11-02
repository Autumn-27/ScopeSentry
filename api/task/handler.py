# -------------------------------------
# @file      : handler.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/10/28 22:09
# -------------------------------------------
import asyncio
import json

from bson import ObjectId

from api.node.handler import get_node_all
from api.page_monitoring import get_page_monitoring_data
from core.apscheduler_handler import scheduler
from api.task.util import get_target_list, parameter_parser
from core.db import get_mongo_db
from core.redis_handler import get_redis_pool, get_redis_online_data
from core.util import get_now_time
from loguru import logger


async def insert_task(request_data, db):
    targetList = await get_target_list(request_data['target'], request_data.get("ignore", ""))
    taskNum = len(targetList)
    request_data['taskNum'] = taskNum
    request_data['target'] = request_data['target'].strip("\n").strip("\r").strip()
    request_data['progress'] = 0
    request_data["creatTime"] = get_now_time()
    request_data["endTime"] = ""
    result = await db.task.insert_one(request_data)
    if result.inserted_id:
        asyncio.create_task(create_scan_task(request_data, str(result.inserted_id)))
        return result.inserted_id


async def create_scan_task(request_data, id):
    logger.info(f"[create_scan_task] begin: {id}")
    try:
        async for db in get_mongo_db():
            async for redis_con in get_redis_pool():
                request_data["id"] = str(id)
                if request_data['allNode']:
                    all_node = await get_node_all(redis_con)
                    for node in all_node:
                        if node not in request_data["node"]:
                            request_data["node"].append(node)
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

                # 获取模板数据
                template_data = await db.ScanTemplates.find_one({"_id": ObjectId(request_data["template"])})
                # 如果选择了poc 将poc参数拼接到nuclei的参数中
                if len(template_data['vullist']) != 0:
                    vul_tmp = ""
                    if "All Poc" in template_data['vullist']:
                        vul_tmp = "*"
                    else:
                        for vul in template_data['vullist']:
                            vul_tmp += vul + ","
                    vul_tmp = vul_tmp.strip(",")
                    if "VulnerabilityScan" in template_data["Parameters"]:
                        if "ed93b8af6b72fe54a60efdb932cf6fbc" in template_data["Parameters"]["VulnerabilityScan"]:
                            template_data["Parameters"]["VulnerabilityScan"]["ed93b8af6b72fe54a60efdb932cf6fbc"] = \
                            template_data["Parameters"]["VulnerabilityScan"][
                                "ed93b8af6b72fe54a60efdb932cf6fbc"] + " -t " + vul_tmp
                # 解析参数，支持{}获取字典
                template_data["Parameters"] = await parameter_parser(template_data["Parameters"], db)
                # 删除原始的vullist
                del template_data["vullist"]
                del template_data["_id"]
                # 设置任务名称
                template_data["TaskName"] = request_data["name"]
                # 设置忽略目标
                template_data["ignore"] = request_data["ignore"]
                # 设置去重
                template_data["duplicates"] = request_data["duplicates"]
                # 任务id
                template_data["ID"] = str(id)
                # 任务类型
                template_data["type"] = "scan"
                # 原始的target生成target list
                target_list = await get_target_list(request_data['target'], request_data.get("ignore", ""))
                # 将任务目标插入redis中
                await redis_con.lpush(f"TaskInfo:{id}", *target_list)
                # 分发任务
                for name in request_data["node"]:
                    await redis_con.rpush(f"NodeTask:{name}", json.dumps(template_data))
                logger.info(f"[create_scan_task] end: {id}")
                return True
    except Exception as e:
        logger.error(str(e))
        logger.info(f"[create_scan_task] error end: {id}")
        # Handle exceptions as needed
        return False


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
        doc["name"] = doc["name"] + f"-{tp}-" + time_now
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
                "type": 'page_monitoring',
                "TaskId": "page_monitoring"
            }
            for name in name_list:
                await redis.rpush(f"NodeTask:{name}", json.dumps(add_redis_task_data))