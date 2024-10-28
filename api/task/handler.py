# -------------------------------------
# @file      : handler.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/10/28 22:09
# -------------------------------------------
import json

from api.task.util import get_target_list
from core.util import generate_random_string, get_now_time
from loguru import logger
from api.node import get_node_all


async def create_scan_task(request_data, id, targetList, redis_con):
    try:
        request_data["id"] = str(id)
        if request_data['allNode']:
            request_data["node"] = await get_node_all(redis_con)
        # 删除可能存在缓存
        keys_to_delete = [
            f"TaskInfo:tmp:{id}",
            f"TaskInfo:{id}",
            f"TaskInfo:time:{id}",
        ]
        progresskeys = await redis_con.keys(f"TaskInfo:progress:{id}:*")
        keys_to_delete.extend(progresskeys)
        progresskeys = await redis_con.keys(f"duplicates:{id}:")
        keys_to_delete.extend(progresskeys)
        if keys_to_delete:
            await redis_con.delete(*keys_to_delete)

        # 原始的target生成targetlist
        target_list = await get_target_list(request_data['target'])
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


def transform_db_redis(request_data):
    Subfinder = False
    Ksubdomain = False
    if "Subfinder" in request_data["subdomainConfig"]:
        Subfinder = True
    if "Ksubdomain" in request_data["subdomainConfig"]:
        Ksubdomain = True
    add_redis_task_data = {
        "TaskId": request_data["id"],
        "SubdomainScan": request_data["subdomainScan"],
        "Subfinder": Subfinder,
        "Ksubdomain": Ksubdomain,
        "UrlScan": request_data["urlScan"],
        "Duplicates": request_data["duplicates"],
        "SensitiveInfoScan": request_data["sensitiveInfoScan"],
        "PageMonitoring": request_data["pageMonitoring"],
        "CrawlerScan": request_data["crawlerScan"],
        "VulScan": request_data["vulScan"],
        "VulList": request_data["vulList"],
        "PortScan": request_data["portScan"],
        "Ports": request_data["ports"],
        "Waybackurl": request_data["waybackurl"],
        "DirScan": request_data["dirScan"],
        "type": 'scan'
    }
    return add_redis_task_data