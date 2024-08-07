# -------------------------------------
# @file      : apscheduler_handler.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/4/21 19:36
# -------------------------------------------

from apscheduler.schedulers.asyncio import AsyncIOScheduler
from apscheduler.jobstores.mongodb import MongoDBJobStore
from core.config import *
mongo_config = {
    'host': MONGODB_IP,
    'port': int(MONGODB_PORT),
    'username': DATABASE_USER,
    'password': DATABASE_PASSWORD,
    'database': DATABASE_NAME,
    'collection': 'apscheduler'
}
jobstores = {
    'mongo': MongoDBJobStore(**mongo_config)
}
scheduler = AsyncIOScheduler(jobstores=jobstores)
