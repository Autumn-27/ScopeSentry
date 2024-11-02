# -------------------------------------
# @file      : handler.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/11/2 14:52
# -------------------------------------------


async def update_redis_data(redis, key):
    await redis.hmset(key, {'state': '3'})


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
                # update_time_str = hash_data.get('updateTime')
                # if update_time_str:
                #     update_time = datetime.strptime(update_time_str, '%Y-%m-%d %H:%M:%S')
                #     time_difference = (
                #                 datetime.strptime(get_now_time(), "%Y-%m-%d %H:%M:%S") - update_time).total_seconds()
                #     logger.info(f'节点时间差：{time_difference}, {get_now_time()}, {update_time}')
                #     if time_difference > NODE_TIMEOUT:
                #         await asyncio.create_task(update_redis_data(redis, key))
                #         hash_data['state'] = '3'
                #     else:
                result.append(name)
        return result


async def get_node_all(redis_con):
    try:
        result = []
        async with redis_con as redis:
            # 获取所有以 node: 开头的键
            keys = await redis.keys("node:*")
            for key in keys:
                name = key.split(":")[1]
                hash_data = await redis.hgetall(key)
                if hash_data.get('state') != '2':
                    result.append(name)

        return result
    except:
        return []