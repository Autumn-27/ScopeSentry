# -------------------------------------
# @file      : system.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/5/14 21:59
# -------------------------------------------

from fastapi import APIRouter, Depends
import git
import httpx
from loguru import logger
from api.users import verify_token
from core.db import get_mongo_db
from core.config import *
import requests

from core.redis_handler import get_redis_pool, refresh_config

router = APIRouter()


@router.get("/system/version")
async def get_system_version(redis_con=Depends(get_redis_pool), _: dict = Depends(verify_token)):
    server_lversion = ""
    server_msg = ""
    scan_lversion = ""
    scan_msg = ""

    async with httpx.AsyncClient() as client:
        try:
            r = await client.get(f"{UPDATEURL}/get/version?name=server", timeout=5)
            r_json = r.json()
            server_lversion = r_json["value"]
            server_msg = r_json['msg']

            r = await client.get(f"{UPDATEURL}/get/version?name=scan", timeout=5)
            r_json = r.json()
            scan_lversion = r_json["value"]
            scan_msg = r_json['msg']
        except Exception as e:
            # 这里可以添加一些日志记录错误信息
            logger.error(str(e))

    result_list = [{"name": "ScopeSentry-Server", "cversion": VERSION, "lversion": server_lversion, "msg": server_msg}]

    async with redis_con as redis:
        keys = await redis.keys("node:*")
        for key in keys:
            name = key.split(":")[1]
            hash_data = await redis.hgetall(key)
            result_list.append({"name": name, "cversion": hash_data["version"], "lversion": scan_lversion, "msg": scan_msg})
    return {
            "code": 200,
            "data": {
                'list': result_list
            }
        }


@router.get("/system/update")
async def system_update():
    update_server()
    await refresh_config("all", 'UpdateSystem')



def update_server():
    repo_path = os.getcwd()
    if not os.path.isdir('.git'):
        repo = git.Repo.init(repo_path)
    else:
        repo = git.Repo.init(repo_path, bare=False)
    if not repo.remotes:
        # 添加远程仓库地址
        repo.create_remote('origin', REMOTE_REPO_URL)
    else:
        # 获取远程地址
        remote_url = repo.remotes.origin.url
        # 检查远程地址是否为预期地址
        if remote_url != REMOTE_REPO_URL:
            repo.remotes.origin.set_url(REMOTE_REPO_URL)
    result = repo.remotes.origin.pull()
    for info in result:
        print(info)

