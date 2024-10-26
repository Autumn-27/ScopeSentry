# -------------------------------------
# @file      : util.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/10/26 22:48
# -------------------------------------------
import json

from bson import ObjectId
from motor.motor_asyncio import AsyncIOMotorCursor

from api.node import get_node_all
from core.apscheduler_handler import scheduler
from core.db import get_mongo_db
from core.redis_handler import get_redis_pool
from core.util import transform_db_redis, generate_random_string, get_now_time
from loguru import logger


async def create_scan_task(request_data, id, targetList, redis_con):
    try:
        request_data["id"] = str(id)
        if request_data['allNode']:
            request_data["node"] = await get_node_all(redis_con)

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


async def delete_asset(task_ids, is_project = False):
    async for db in get_mongo_db():
        key = ["asset", "subdomain", "SubdoaminTakerResult", "UrlScan", "crawler", "SensitiveResult", "DirScanResult", "vulnerability", "PageMonitoring"]
        del_query = {"taskId": {"$in": task_ids}}
        if is_project:
            del_query = {
                            "$or": [
                                {"taskId": {"$in": task_ids}},
                                {"project": {"$in": task_ids}}
                            ]
                        }
        for k in key:
            result = await db[k].delete_many(del_query)
            if result.deleted_count > 0:
                logger.info("Deleted {} {} documents".format(k, result.deleted_count))
            else:
                logger.info("Deleted {} None documents".format(k))