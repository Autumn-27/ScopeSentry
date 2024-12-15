# -------------------------------------
# @file      : system.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/5/14 21:59
# -------------------------------------------
import subprocess
import sys
import time
import traceback

from fastapi import APIRouter, Depends
import git
import httpx
from loguru import logger
from api.users import verify_token
from core.config import *
from core.redis_handler import get_redis_pool, refresh_config

router = APIRouter()


@router.get("/system/version")
async def get_system_version(redis_con=Depends(get_redis_pool), _: dict = Depends(verify_token)):
    server_lversion = ""
    server_msg = ""
    scan_lversion = ""
    scan_msg = ""

    async with httpx.AsyncClient(verify=False) as client:
        try:
            r = await client.get(f"https://gitee.com/constL/scope-sentry/raw/main/version.json", timeout=10)
            r_json = r.json()
            server_lversion = r_json["server"]
            server_msg = r_json['server_msg']
            scan_lversion = r_json["scan"]
            scan_msg = r_json['scan_msg']
        except:
            try:
                r = await client.get(f"https://raw.githubusercontent.com/Autumn-27/ScopeSentry/main/version.json",
                                     timeout=10)
                r_json = r.json()
                server_lversion = r_json["server"]
                server_msg = r_json['server_msg']
                scan_lversion = r_json["scan"]
                scan_msg = r_json['scan_msg']
            except Exception as e:
                logger.error(f"An unexpected error occurred: {e}")

    result_list = [{"name": "ScopeSentry-Server", "cversion": VERSION, "lversion": server_lversion, "msg": server_msg}]
    try:
        async with redis_con as redis:
            keys = await redis.keys("node:*")
            for key in keys:
                name = key.split(":")[1]
                hash_data = await redis.hgetall(key)
                result_list.append(
                    {"name": name, "cversion": hash_data["version"], "lversion": scan_lversion, "msg": scan_msg})
    except:
        pass
    return {
        "code": 200,
        "data": {
            'list': result_list
        }
    }


@router.post("/system/update")
async def system_update(request_data: dict, _: dict = Depends(verify_token)):
    with open("PLUGINKEY", 'r') as file:
        plg_key = file.read()
    key = request_data.get("key", "")
    if key == "":
        return {"message": f"key error", "code": 505}
    if plg_key != key:
        return {"message": f"key error", "code": 505}
    server_url = request_data.get("server", "")
    scan_url = request_data.get("scan", "")
    if server_url == "" or scan_url == "":
        return {"message": f"url not set", "code": 404}
    await refresh_config("all", 'UpdateSystem', scan_url)

    with open("/opt/ScopeSentry/UPDATE", 'w') as file:
        file.write(server_url)
    time.sleep(3)
    sys.exit("System updated, exiting application.")
    # await update_server()
    # await refresh_config("all", 'UpdateSystem')


async def update_server():
    relative_path = f'requirements.txt'
    file_path = os.path.join(os.getcwd(), relative_path)
    async with httpx.AsyncClient() as client:
        try:
            r = await client.get(f"https://raw.githubusercontent.com/Autumn-27/ScopeSentry/main/requirements.txt",
                                 timeout=5)
            content = r.text
            with open(file_path, "w") as f:
                f.write(content)
            logger.info("requirements.txt write successfully")
            result = subprocess.run(
                ['pip', 'install', '-r', file_path, '-i', 'https://pypi.tuna.tsinghua.edu.cn/simple', '--no-cache-dir'],
                capture_output=True, text=True
            )
            if result.returncode == 0:
                logger.info("Packages installed successfully")
            else:
                logger.error(f"Error installing packages: {result.stderr}")
        except Exception as e:
            logger.error(str(e))
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
