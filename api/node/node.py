# -------------------------------------
# @file      : node.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/11/2 14:51
# -------------------------------------------
from datetime import datetime
from fastapi import APIRouter, Depends

from api.node.handler import update_redis_data, get_redis_online_data
from core.config import *
from api.users import verify_token
from core.db import get_mongo_db
from core.redis_handler import get_redis_pool
from core.util import get_now_time
from core.redis_handler import refresh_config
import asyncio
from loguru import logger

router = APIRouter()


@router.get("/data")
async def node_data(_: dict = Depends(verify_token), redis_con=Depends(get_redis_pool)):
    async with redis_con as redis:
        # 获取所有以 node: 开头的键
        keys = await redis.keys("node:*")
        # 构建结果字典
        result = []
        for key in keys:
            name = key.split(":")[1]
            # 获取哈希中的所有字段和值
            hash_data = await redis.hgetall(key)
            # 在哈希数据中增加键为 name，值为键的名称
            hash_data['name'] = name
            if hash_data.get('state') == '1':
                update_time_str = hash_data.get('updateTime')
                if update_time_str:
                    update_time = datetime.strptime(update_time_str, '%Y-%m-%d %H:%M:%S')
                    time_difference = (datetime.strptime(get_now_time(), "%Y-%m-%d %H:%M:%S") - update_time).total_seconds()
                    logger.info(f'节点时间差：{time_difference}, {get_now_time()}, {update_time}')
                    if time_difference > NODE_TIMEOUT:
                        await asyncio.create_task(update_redis_data(redis, key))
                        hash_data['state'] = '3'

                # 将哈希数据添加到结果数组中
            result.append(hash_data)
            result.sort(key=lambda x: x["name"])
        return {
            "code": 200,
            "data": {
                'list': result
            }
        }


@router.get("/data/online")
async def node_data_online(_: dict = Depends(verify_token), redis_con=Depends(get_redis_pool)):
    result = await get_redis_online_data(redis_con)
    return {
        "code": 200,
        "data": {
            'list': result
        }
    }


@router.post("/config/update")
async def node_config_update(config_data: dict, _: dict = Depends(verify_token), redis_con=Depends(get_redis_pool)):
    try:
        name = config_data.get("name")
        old_name = config_data.get("oldName", "")
        if old_name == "":
            old_name = name
        modulesConfig = config_data.get("ModulesConfig")
        state = config_data.get("state")
        if name is None or modulesConfig is None or state is None:
            return {"code": 400, "message": "Invalid request, missing required parameters"}
        msg = f"{name}[*]{state}[*]{modulesConfig}"
        await refresh_config(old_name, 'nodeConfig', msg)
        return {"code": 200, "message": "Node configuration updated successfully"}
    except Exception as e:
        return {"code": 500, "message": f"Internal server error: {str(e)}"}


@router.post("/delete")
async def delete_node_rules(request_data: dict, _: dict = Depends(verify_token), redis_con=Depends(get_redis_pool)):
    try:
        node_names = request_data.get("names", [])
        for name in node_names:
            logger.info("delete node:" + name)
            await redis_con.delete("node:" + name)
        return {"message": "Node deleted successfully", "code": 200}
    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/log/data")
async def get_node_logs(request_data: dict, _: dict = Depends(verify_token), redis_con=Depends(get_redis_pool)):
    try:
        node_name = request_data.get("name")
        if not node_name:
            return {"message": "Node name is required", "code": 400}
        # 构建日志键
        log_key = f"log:{node_name}"
        # 从 Redis 中获取日志列表
        logs = await redis_con.lrange(log_key, 0, -1)
        log_data = ""
        for log in logs:
            log_data += log
        return {"code": 200, "logs": log_data}
    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "Error retrieving logs", "code": 500}


@router.post("/plugin")
async def get_node_plugin(request_data: dict, _: dict = Depends(verify_token), redis_con=Depends(get_redis_pool), db=Depends(get_mongo_db)):
    try:
        node_name = request_data.get("name")
        if not node_name:
            return {"message": "Node name is required", "code": 400}
        cursor = db['plugins'].find({}, {"_id": 0,
                                         "name": 1,
                                         "hash": 1,
                                         "module": 1
                                         })
        result = await cursor.to_list(length=None)
        plugin_list = {}
        for r in result:
            plugin_list[r["hash"]] = {"name": r["name"], "module": r["module"]}
        key = "NodePlg:" + node_name
        hash_data = await redis_con.hgetall(key)
        result_list = []
        # 0 代表未安装 1代表安装失败 2代表安装成功，未检查 3代表安装成功，检查失败 4代表安装检查都成功
        for plg in plugin_list:
            try:
                install_value = hash_data[plg + "_install"]
                check_value = hash_data[plg + "_check"]
                result_list.append({
                    "name": plugin_list[plg]["name"],
                    "install": install_value,
                    "check": check_value,
                    "hash": plg,
                    "module": plugin_list[plg]["module"]
                })
            except:
                result_list.append({
                    "name": plugin_list[plg]["name"],
                    "install": "0",
                    "check": "0",
                    "hash": plg,
                    "module": plugin_list[plg]["module"]
                })
        return {"code": 200, "data": {"list": result_list}}
    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "Error retrieving logs", "code": 500}


@router.post("/restart")
async def restart_node(request_data: dict, _: dict = Depends(verify_token), redis_con=Depends(get_redis_pool), db=Depends(get_mongo_db)):
    node_name = request_data.get("name")
    if not node_name:
        return {"message": "Node name is required", "code": 400}

    await refresh_config(node_name, 'restart')
    return {"message": "Node restart successfully", "code": 200}
