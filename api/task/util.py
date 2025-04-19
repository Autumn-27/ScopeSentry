# -------------------------------------
# @file      : util.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/10/26 22:48
# -------------------------------------------
import ipaddress
import re

from motor.motor_asyncio import AsyncIOMotorCursor

from core.db import get_mongo_db
from core.redis_handler import get_redis_pool
from loguru import logger


async def generate_ip_range(start_ip, end_ip):
    start = ipaddress.ip_address(start_ip)
    end = ipaddress.ip_address(end_ip)

    ip_list = []
    current_ip = start
    while current_ip <= end:
        ip_list.append(str(current_ip))
        current_ip += 1

    return ip_list


# 解析ip段
async def generate_target(target):
    try:
        if "://" in target:
            return [target]
        if '-' in target:
            start_ip, end_ip = target.split('-')
            l = await generate_ip_range(start_ip, end_ip)
            return l
        elif '/' in target:
            network = ipaddress.ip_network(target, strict=False)
            return [str(ip) for ip in network.hosts()]
        else:
            return [target]
    except Exception as e:
        return [target]


# 根据原始target生成目标列表
async def get_target_list(raw_target, ignore):
    ignore_list, regex_list = await generate_ignore(ignore)
    # 使用集合来避免重复
    target_dict = {}
    for t in raw_target.split("\n"):
        t = t.strip("\n").strip("\r").strip()
        result = await generate_target(t)
        for r in result:
            if r.strip("\n").strip("\r").strip() == "":
                continue
            if r not in ignore_list:
                if len(regex_list) != 0:
                    for rege_str in regex_list:
                        rege = re.compile(rege_str)
                        if not rege.search(r):
                            target_dict[r] = None  # 使用字典的键去重
                else:
                    target_dict[r] = None  # 使用字典的键去重
    return list(target_dict.keys())


async def generate_ignore(ignore):
    # 处理'NoneType' object has no attribute 'split'的报错
    if ignore is None:
        ignore = ""
    ignore_list = []
    regex_list = []
    for t in ignore.split("\n"):
        t.replace("http://", "").replace("https://", "")
        t = t.strip("\n").strip("\r").strip()
        if "*" not in t:
            result = await generate_target(t)
            ignore_list.extend(result)
        else:
            t_escaped = re.escape(t)
            regex_list.append(t_escaped.replace(r"\*", ".*"))
    return ignore_list, regex_list


async def task_progress():
    async for db in get_mongo_db():
        async for redis in get_redis_pool():
            query = {"progress": {"$ne": 100}, "status": 1}
            cursor: AsyncIOMotorCursor = db.task.find(query)
            result = await cursor.to_list(length=None)
            if len(result) == 0:
                return True
            for r in result:
                id = str(r["_id"])
                key = f"TaskInfo:tmp:{id}"
                exists = await redis.exists(key)
                if exists:
                    count = await redis.scard(key)
                    progress_tmp = round(count / r['taskNum'], 2)
                    progress_tmp = round(progress_tmp * 100, 1)
                    if progress_tmp > 100:
                        progress_tmp = 100
                    if progress_tmp == 100:
                        time_key = f"TaskInfo:time:{id}"
                        time_value = await redis.get(time_key)
                        await db.task.update_one({"_id": r["_id"]}, {"$set": {"endTime": time_value, "status": 3}})
                        # 任务结束 删除统计信息
                        await redis.delete(key)
                        await redis.delete(time_key)
                    await db.task.update_one({"_id": r["_id"]}, {"$set": {"progress": progress_tmp}})
                else:
                    await db.task.update_one({"_id": r["_id"]}, {"$set": {"progress": 0}})
            return


async def delete_asset(task_ids, is_project=False):
    async for db in get_mongo_db():
        key = ["asset", "subdomain", "SubdoaminTakerResult", "UrlScan", "crawler", "SensitiveResult", "DirScanResult",
               "vulnerability", "PageMonitoring"]
        del_query = {"taskName": {"$in": task_ids}}
        if is_project:
            del_query = {
                "$or": [
                    {"taskName": {"$in": task_ids}},
                    {"project": {"$in": task_ids}}
                ]
            }
        for k in key:
            result = await db[k].delete_many(del_query)
            if result.deleted_count > 0:
                logger.info("Deleted {} {} documents".format(k, result.deleted_count))
            else:
                logger.info("Deleted {} None documents".format(k))
