# -------------------------------------
# @file      : scan.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/10/29 21:01
# -------------------------------------------
from core.apscheduler_handler import scheduler
from core.db import get_mongo_db
from loguru import logger

from core.redis_handler import get_redis_pool


async def scheduler_project(id):
    print("d")
    # logger.info(f"Scheduler project {id}")
    # async for db in get_mongo_db():
    #     async for redis in get_redis_pool():
    #         job = scheduler.get_job(id)
    #         task_id = generate_random_string(15)
    #         if job:
    #             next_time = job.next_run_time
    #             formatted_time = next_time.strftime("%Y-%m-%d %H:%M:%S")
    #             doc = await db.ScheduledTasks.find_one({"id": id})
    #             run_id_last = doc.get("runner_id", "")
    #             if run_id_last != "":
    #                 progresskeys = await redis.keys(f"TaskInfo:progress:{run_id_last}:*")
    #                 for pgk in progresskeys:
    #                     await redis.delete(pgk)
    #             update_document = {
    #                 "$set": {
    #                     "lastTime": get_now_time(),
    #                     "nextTime": formatted_time,
    #                     "runner_id": task_id
    #                 }
    #             }
    #             await db.ScheduledTasks.update_one({"id": id}, update_document)
    #         query = {"_id": ObjectId(id)}
    #         doc = await db.project.find_one(query)
    #         targetList = []
    #         target_data = await db.ProjectTargetData.find_one({"id": id})
    #         for t in target_data.get('target', '').split("\n"):
    #             t.replace("http://", "").replace("https://", "")
    #             t = t.strip("\n").strip("\r").strip()
    #             if t != "" and t not in targetList:
    #                 targetList.append(t)
    #         await create_scan_task(doc, task_id, targetList, redis)