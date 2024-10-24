# -*- coding:utf-8 -*-　　
# @name: db
# @auth: rainy-autumn@outlook.com
# @version:
import time
from urllib.parse import quote_plus
from motor.motor_asyncio import AsyncIOMotorGridFSBucket
from motor.motor_asyncio import AsyncIOMotorClient, AsyncIOMotorCursor
from core.default import *
from core.config import *
from loguru import logger


async def get_mongo_db():

    client = AsyncIOMotorClient(f"mongodb://{DATABASE_USER}:{quote_plus(DATABASE_PASSWORD)}@{MONGODB_IP}:{str(MONGODB_PORT)}",
                                serverSelectionTimeoutMS=10000, unicode_decode_error_handler='ignore')
    db = client[DATABASE_NAME]
    try:
        yield db
    finally:
        client.close()


async def create_database():
    client = None
    check_flag = 0
    try:
        while True:
            try:
                # 创建新的 MongoDB 客户端
                client = AsyncIOMotorClient(f"mongodb://{quote_plus(DATABASE_USER)}:{quote_plus(DATABASE_PASSWORD)}@{MONGODB_IP}:{str(MONGODB_PORT)}",
                                            serverSelectionTimeoutMS=2000)
                break
            except Exception as e:
                time.sleep(10)
                check_flag += 1
                if check_flag == 10:
                    logger.error(f"Error re creating database: {e}")
                    exit(1)
        # 获取数据库列表
        database_names = await client.list_database_names()
        db = client[DATABASE_NAME]
        # 如果数据库不存在，创建数据库
        if DATABASE_NAME not in database_names:
            # 在数据库中创建一个集合，比如名为 "user"
            collection = db["user"]

            # 用户数据
            await collection.insert_one({"username": "ScopeSentry",
                                         'password': 'b0ce71fcbed8a6ca579d52800145119cc7d999dc8651b62dfc1ced9a984e6e64'})

            collection = db["config"]
            # 系统配置
            await collection.insert_one(
                {"name": "timezone", 'value': 'Asia/Shanghai', 'type': 'system'})
            await collection.insert_one(
                {"name": "MaxTaskNum", 'value': '7', 'type': 'system'})
            await collection.insert_one(
                {"name": "DirscanThread", 'value': '15', 'type': 'system'})
            await collection.insert_one(
                {"name": "PortscanThread", 'value': '5', 'type': 'system'})
            await collection.insert_one(
                {"name": "CrawlerThread", 'value': '2', 'type': 'system'})
            await collection.insert_one(
                {"name": "UrlMaxNum", 'value': '500', 'type': 'system'})
            await collection.insert_one(
                {"name": "UrlThread", 'value': '5', 'type': 'system'})
            # 设置时区为Asia/Shanghai
            # SHA_TZ = timezone(TIMEZONE)
            # timezone('Asia/Shanghai')
            # utc_now = datetime.utcnow().replace(tzinfo=timezone.utc)
            # time_now = utc_now.astimezone(SHA_TZ)
            # formatted_time = time_now.strftime("%Y-%m-%d %H:%M:%S")
            # subfinder配置
            collection = db["config"]
            # 插入一条数据
            await collection.insert_one(
                {"name": "SubfinderApiConfig", 'value': subfinderApiConfig, 'type': 'subfinder'})
            await collection.insert_one(
                {"name": "RadConfig", 'value': radConfig, 'type': 'rad'})
            # dirDict = get_dirDict()
            # await collection.insert_one(
            #     {"name": "DirDic", 'value': dirDict, 'type': 'dirDict'})
            # 目录扫描字典
            fs = AsyncIOMotorGridFSBucket(db)
            content = get_dirDict()
            if content:
                byte_content = content.encode('utf-8')
                await fs.upload_from_stream('dirdict', byte_content)
            # 子域名字典
            content = get_domainDict()
            if content:
                byte_content = content.encode('utf-8')
                await fs.upload_from_stream('DomainDic', byte_content)
                logger.info("Document DomainDic uploaded to GridFS.")

            await collection.insert_one(
                {"name": "notification", 'dirScanNotification': True,
                 'portScanNotification': True, 'sensitiveNotification': True,
                 'subdomainTakeoverNotification': True,
                 'pageMonNotification': True,
                 'subdomainNotification': True,
                 'vulNotification': True,
                 'type': 'notification'})
            # domainDict = get_domainDict()
            # await collection.insert_one(
            #     {"name": "DomainDic", 'value': domainDict, 'type': 'domainDict'})
            sensitive_data = get_sensitive()
            collection = db["SensitiveRule"]
            if sensitive_data:
                await collection.insert_many(sensitive_data)

            collection = db["ScheduledTasks"]
            await collection.insert_one(
                {"id": "page_monitoring", "name": "Page Monitoring", 'hour': 24, 'node': [], 'allNode': True, 'type': 'Page Monitoring', 'state': True})

            await db.create_collection("notification")

            collection = db["PortDict"]
            await collection.insert_many(portDic)

            collection = db["PocList"]
            pocData = get_poc()
            await collection.insert_many(pocData)

            collection = db["project"]
            project_data, target_data = get_project_data()
            await collection.insert_many(project_data)

            collection = db["ProjectTargetData"]
            await collection.insert_many(target_data)

            collection = db["FingerprintRules"]
            fingerprint = get_finger()
            await collection.insert_many(fingerprint)
        else:
            collection = db["config"]
            result = await collection.find_one({"name": "timezone"})
            set_timezone(result.get('value', 'Asia/Shanghai'))

            collection = db["ScheduledTasks"]
            result = await collection.find_one({"id": "page_monitoring"})
            if not result:
                await collection.insert_one(
                    {"id": "page_monitoring", "name": "Page Monitoring", 'hour': 24, 'type': 'Page Monitoring', 'state': True})
        await get_fingerprint(db)
        # await get_sens_rule(db)
        await get_project(db)
    except Exception as e:
        # 处理异常
        logger.error(f"Error creating database: {e}")
        exit(0)
    finally:
        # 在适当的地方关闭 MongoDB 客户端
        if client:
            client.close()


async def get_fingerprint(client):
    collection = client["FingerprintRules"]
    cursor = collection.find({}, {"_id": 1, "name": 1})
    async for document in cursor:
        document['id'] = str(document['_id'])
        del document['_id']
        APP[document['id']] = document['name']


# async def get_sens_rule(client):
#     collection = client["SensitiveRule"]
#     cursor = collection.find({}, {"_id": 1, "name": 1, "color": 1})
#     async for document in cursor:
#         document['id'] = str(document['_id'])
#         del document['_id']
#         SensitiveRuleList[document['id']] = {
#             "name": document['name'],
#             "color": document['color']
#         }


async def get_project(client):
    collection = client["project"]
    cursor = collection.find({}, {"_id": 1, "name": 1})
    async for document in cursor:
        document['id'] = str(document['_id'])
        Project_List[document['name'].lower()] = document['id']
