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
            # 扫描模块配置
            await collection.insert_one(
                {"name": "ModulesConfig", 'value': ModulesConfig, 'type': 'system'})
            await collection.insert_one(
                {"name": "timezone", 'value': 'Asia/Shanghai', 'type': 'system'})
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
            fs = AsyncIOMotorGridFSBucket(db)
            # 更新目录扫描默认字典
            content = get_dirDict()
            size = len(content) / (1024 * 1024)
            result = await db["dictionary"].insert_one(
                {"name": "default", "category": "dir", "size": "{:.2f}".format(size)})
            if result.inserted_id:
                await fs.upload_from_stream(
                    str(result.inserted_id),  # 使用id作为文件名存储
                    content  # 文件内容
                )

            # 更新子域名默认字典
            content = get_domainDict()
            size = len(content) / (1024 * 1024)
            result = await db["dictionary"].insert_one(
                {"name": "default", "category": "subdomain", "size": "{:.2f}".format(size)})
            if result.inserted_id:
                await fs.upload_from_stream(
                    str(result.inserted_id),  # 使用id作为文件名存储
                    content  # 文件内容
                )

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
            # 创建默认插件
            collection = db["plugins"]
            await collection.insert_many(PLUGINS)

            # 创建默认扫描模板
            collection = db["ScanTemplates"]
            await collection.insert_one(SCANTEMPLATE)
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
