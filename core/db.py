# -*- coding:utf-8 -*-　　
# @name: db
# @auth: rainy-autumn@outlook.com
# @version:
import hashlib
import time
from urllib.parse import quote_plus
from motor.motor_asyncio import AsyncIOMotorGridFSBucket
from motor.motor_asyncio import AsyncIOMotorClient, AsyncIOMotorCursor
from pymongo import ASCENDING

from core.default import *
from core.config import *
from loguru import logger

from core.util import print_progress_bar


async def get_mongo_db():

    client = AsyncIOMotorClient(f"mongodb://{MONGODB_USER}:{quote_plus(str(MONGODB_PASSWORD))}@{MONGODB_IP}:{str(MONGODB_PORT)}",
                                serverSelectionTimeoutMS=10000, unicode_decode_error_handler='ignore')
    db = client[MONGODB_DATABASE]
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
                client = AsyncIOMotorClient(f"mongodb://{quote_plus(MONGODB_USER)}:{quote_plus(str(MONGODB_PASSWORD))}@{MONGODB_IP}:{str(MONGODB_PORT)}",
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
        db = client[MONGODB_DATABASE]
        # 如果数据库不存在，创建数据库
        if MONGODB_DATABASE not in database_names:
            # 在数据库中创建一个集合，比如名为 "user"
            collection = db["user"]
            password = generate_random_string(8)
            print("\n" + "=" * 50)
            print("✨✨✨ IMPORTANT NOTICE: Please review the User/Password below ✨✨✨")
            print("=" * 50)
            print(f"🔑 User/Password: ScopeSentry/{password}")
            print("=" * 50)
            print("✅ Ensure the User/Password is correctly copied!\n")
            print("✅ The initialization password is stored in the file PASSWORD\n")
            with open("PASSWORD", 'w') as file:
                file.write(password)

            total_steps = 16
            # 用户数据
            await collection.insert_one({"username": "ScopeSentry",
                                         'password': hashlib.sha256(password.encode()).hexdigest()})
            logger.info("Project initialization")
            print_progress_bar(1, total_steps, "install")

            collection = db["config"]
            # 扫描模块配置
            await collection.insert_one(
                {"name": "ModulesConfig", 'value': ModulesConfig, 'type': 'system'})
            await collection.insert_one(
                {"name": "timezone", 'value': 'Asia/Shanghai', 'type': 'system'})

            print_progress_bar(2, total_steps, "install")
            # subfinder配置
            collection = db["config"]
            await collection.insert_one(
                {"name": "SubfinderApiConfig", 'value': subfinderApiConfig, 'type': 'subfinder'})
            print_progress_bar(3, total_steps, "install")
            # rad配置
            await collection.insert_one(
                {"name": "RadConfig", 'value': radConfig, 'type': 'rad'})
            print_progress_bar(4, total_steps, "install")
            # 通知配置
            await collection.insert_one(
                {"name": "notification", 'dirScanNotification': True,
                 'portScanNotification': True, 'sensitiveNotification': True,
                 'subdomainTakeoverNotification': True,
                 'pageMonNotification': True,
                 'subdomainNotification': True,
                 'vulNotification': True,
                 'type': 'notification'})

            fs = AsyncIOMotorGridFSBucket(db)
            print_progress_bar(5, total_steps, "install")
            # 更新目录扫描默认字典
            content = get_dirDict()
            size = len(content) / (1024 * 1024)
            result = await db["dictionary"].insert_one(
                {"name": "default", "category": "dir", "size": "{:.2f}".format(size)})
            if result.inserted_id:
                await fs.upload_from_stream(
                    str(result.inserted_id),  # 使用id作为文件名存储
                    content.encode("utf-8")  # 文件内容
                )
            print_progress_bar(6, total_steps, "install")
            # 更新子域名默认字典
            content = get_domainDict()
            size = len(content) / (1024 * 1024)
            result = await db["dictionary"].insert_one(
                {"name": "default", "category": "subdomain", "size": "{:.2f}".format(size)})
            if result.inserted_id:
                await fs.upload_from_stream(
                    str(result.inserted_id),  # 使用id作为文件名存储
                    content.encode("utf-8")  # 文件内容
                )
            print_progress_bar(7, total_steps, "install")
            # 插入敏感信息
            sensitive_data = get_sensitive()
            collection = db["SensitiveRule"]
            if sensitive_data:
                await collection.insert_many(sensitive_data)

            print_progress_bar(8, total_steps, "install")
            # 定时任务
            collection = db["ScheduledTasks"]
            await collection.insert_one(
                {"id": "page_monitoring", "name": "Page Monitoring", 'hour': 24, 'node': [], 'allNode': True, 'type': 'Page Monitoring', 'state': True})
            print_progress_bar(9, total_steps, "install")
            # 通知API
            await db.create_collection("notification")
            print_progress_bar(10, total_steps, "install")
            # 默认端口
            collection = db["PortDict"]
            await collection.insert_many(portDic)
            print_progress_bar(11, total_steps, "install")
            # poc导入
            collection = db["PocList"]
            pocData = get_poc()
            await collection.insert_many(pocData)

            print_progress_bar(12, total_steps, "install")
            # 新版本不内置项目
            # collection = db["project"]
            # project_data, target_data = get_project_data()
            # await collection.insert_many(project_data)
            #
            # collection = db["ProjectTargetData"]
            # await collection.insert_many(target_data)
            print_progress_bar(13, total_steps, "install")
            # 指纹导入
            collection = db["FingerprintRules"]
            fingerprint = get_finger()
            await collection.insert_many(fingerprint)
            print_progress_bar(14, total_steps, "install")
            # 创建默认插件
            collection = db["plugins"]
            await collection.insert_many(PLUGINS)
            print_progress_bar(15, total_steps, "install")
            # 创建默认扫描模板
            collection = db["ScanTemplates"]
            await collection.insert_one(SCANTEMPLATE)
            print_progress_bar(16, total_steps, "install")
            # 创建页面监控文档，url不重复
            db['PageMonitoring'].create_index([('url', ASCENDING)], unique=True)
            db['PageMonitoringBody'].create_index([('md5', ASCENDING)], unique=True)
            # 创建RootDomain
            db['PageMonitoringBody'].create_index([('domain', ASCENDING)], unique=True)
            # 创建asset集合索引
            db['asset'].create_index([('time', ASCENDING)])
            db['asset'].create_index([('url', ASCENDING)])
            db['asset'].create_index([('host', ASCENDING)])
            db['asset'].create_index([('ip', ASCENDING)])
            db['asset'].create_index([('port', ASCENDING)])
            db['asset'].create_index([('host', ASCENDING), ('port', ASCENDING)], unique=True)
            db['asset'].create_index([('project', ASCENDING)])
            db['asset'].create_index([('taskName', ASCENDING)])
            db['asset'].create_index([('rootDomain', ASCENDING)])

            db['subdomain'].create_index([('project', ASCENDING)])
            db['subdomain'].create_index([('taskName', ASCENDING)])
            db['subdomain'].create_index([('rootDomain', ASCENDING)])
            db['subdomain'].create_index([('time', ASCENDING)])

            db['UrlScan'].create_index([('project', ASCENDING)])
            db['UrlScan'].create_index([('taskName', ASCENDING)])
            db['UrlScan'].create_index([('rootDomain', ASCENDING)])

            db['crawler'].create_index([('project', ASCENDING)])
            db['crawler'].create_index([('taskName', ASCENDING)])
            db['crawler'].create_index([('rootDomain', ASCENDING)])

            db['SensitiveResult'].create_index([('project', ASCENDING)])
            db['SensitiveResult'].create_index([('taskName', ASCENDING)])
            db['SensitiveResult'].create_index([('rootDomain', ASCENDING)])

            db['DirScanResult'].create_index([('project', ASCENDING)])
            db['DirScanResult'].create_index([('taskName', ASCENDING)])
            db['DirScanResult'].create_index([('rootDomain', ASCENDING)])

            db['vulnerability'].create_index([('project', ASCENDING)])
            db['vulnerability'].create_index([('taskName', ASCENDING)])
            db['vulnerability'].create_index([('rootDomain', ASCENDING)])

            db['RootDomain'].create_index([('project', ASCENDING)])
            db['RootDomain'].create_index([('taskName', ASCENDING)])
            db['RootDomain'].create_index([('domain', ASCENDING)], unique=True)
            db['RootDomain'].create_index([('time', ASCENDING)])

            db['app'].create_index([('project', ASCENDING)])
            db['app'].create_index([('taskName', ASCENDING)])
            db['app'].create_index([('time', ASCENDING)])
            db['app'].create_index([('name', ASCENDING)])

            db['mp'].create_index([('project', ASCENDING)])
            db['mp'].create_index([('taskName', ASCENDING)])
            db['mp'].create_index([('time', ASCENDING)])
            db['mp'].create_index([('name', ASCENDING)])

            logger.success("Project initialization successful")
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
        # Project_List[document['name'].lower()] = document['id']
        Project_List[document['id']] = document['name']
