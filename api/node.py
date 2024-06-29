# -*- coding:utf-8 -*-　　
# @name: node
# @auth: rainy-autumn@outlook.com
# @version:
from datetime import datetime
from fastapi import WebSocket
from fastapi import APIRouter, Depends
from starlette.websockets import WebSocketDisconnect
from core.config import *
from api.users import verify_token
from core.redis_handler import get_redis_pool
from core.util import get_now_time
from core.redis_handler import refresh_config
import asyncio, json
from loguru import logger

router = APIRouter()


async def update_redis_data(redis, key):
    await redis.hmset(key, {'state': '3'})


@router.get("/node/data")
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
        return {
            "code": 200,
            "data": {
                'list': result
            }
        }


@router.get("/node/data/online")
async def node_data_online(_: dict = Depends(verify_token), redis_con=Depends(get_redis_pool)):
    result = await get_redis_online_data(redis_con)
    return {
        "code": 200,
        "data": {
            'list': result
        }
    }


@router.post("/node/config/update")
async def node_config_update(config_data: dict, _: dict = Depends(verify_token), redis_con=Depends(get_redis_pool)):
    try:
        name = config_data.get("name")
        max_task_num = config_data.get("maxTaskNum")
        state = config_data.get("state")
        if name is None or max_task_num is None or state is None:
            return {"code": 400, "message": "Invalid request, missing required parameters"}

        async with redis_con as redis:
            key = f"node:{name}"
            redis_state = await redis.hget(key, "state")
            if state:
                if redis_state == "2":
                    await redis.hset(key, "state", "1")
            else:
                if redis_state == "1":
                    await redis.hset(key, "state", "2")
            del config_data["name"]
            del config_data["state"]
            for c in config_data:
                await redis.hset(key, c, config_data[c])
            await refresh_config(name, 'nodeConfig')
        return {"code": 200, "message": "Node configuration updated successfully"}

    except Exception as e:
        return {"code": 500, "message": f"Internal server error: {str(e)}"}


@router.post("/node/delete")
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


@router.post("/node/log/data")
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


async def get_redis_online_data(redis_con):
    async with redis_con as redis:
        # 获取所有以 node: 开头的键
        keys = await redis.keys("node:*")
        # 构建结果字典
        result = []
        for key in keys:
            name = key.split(":")[1]
            hash_data = await redis.hgetall(key)
            if hash_data.get('state') == '1':
                update_time_str = hash_data.get('updateTime')
                if update_time_str:
                    update_time = datetime.strptime(update_time_str, '%Y-%m-%d %H:%M:%S')
                    time_difference = (
                                datetime.strptime(get_now_time(), "%Y-%m-%d %H:%M:%S") - update_time).total_seconds()
                    logger.info(f'节点时间差：{time_difference}, {get_now_time()}, {update_time}')
                    if time_difference > NODE_TIMEOUT:
                        await asyncio.create_task(update_redis_data(redis, key))
                        hash_data['state'] = '3'
                    else:
                        result.append(name)
        return result
