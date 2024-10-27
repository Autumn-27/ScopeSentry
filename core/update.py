# -------------------------------------
# @file      : update.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/6/16 14:14
# -------------------------------------------
from urllib.parse import urlparse

from loguru import logger
from motor.motor_asyncio import AsyncIOMotorGridFSBucket
import tldextract
from pymongo import UpdateOne, ASCENDING

from core.config import VERSION
from core.default import get_dirDict, get_domainDict, get_sensitive, ModulesConfig, PLUGINS, SCANTEMPLATE


async def update14(db):
    # 默认项目有个root_domain为空导致匹配上所有资产
    cursor = db.project.find({"root_domains": ""}, {"_id": 1, "root_domains": 1})
    async for document in cursor:
        logger.info("Update found empty root_domains")
        root_domain = []
        for root in document["root_domains"]:
            if root != "":
                root_domain.append(root)
        update_document = {
            "$set": {
                "root_domains": root_domain,
            }
        }
        await db.project.update_one({"_id": document['_id']}, update_document)
    # 修改目录字典存储方式
    fs = AsyncIOMotorGridFSBucket(db)
    result = await db.config.find_one({"name": "DirDic"})
    if result:
        await db.config.delete_one({"name": "DirDic"})
        content = get_dirDict()
        if content:
            byte_content = content.encode('utf-8')
            await fs.upload_from_stream('dirdict', byte_content)
            logger.info("Document DirDict uploaded to GridFS.")
        else:
            logger.error("No dirdict content to upload.")
    # 修改子域名字典存储方式
    result = await db.config.find_one({"name": "DomainDic"})
    if result:
        await db.config.delete_one({"name": "DomainDic"})
        content = get_domainDict()
        if content:
            byte_content = content.encode('utf-8')
            await fs.upload_from_stream('DomainDic', byte_content)
            logger.info("Document DomainDic uploaded to GridFS.")
        else:
            logger.error("No DomainDic content to upload.")

    # 更新敏感信息
    sensitive_data = get_sensitive()
    collection = db["SensitiveRule"]
    if sensitive_data:
        for s in sensitive_data:
            await collection.update_one(
                {"name": s['name']},
                {"$set": s},
                upsert=True
            )
    await db.config.update_one({"name": "version"}, {"$set": {"update": True, "version": float(VERSION)}})


async def update15(db):
    await db.config.insert_one(
        {"name": "ModulesConfig", 'value': ModulesConfig, 'type': 'system'})
    fs = AsyncIOMotorGridFSBucket(db)

    # 更新目录扫描默认字典
    content = get_dirDict()
    size = len(content) / (1024 * 1024)
    result = await db["dictionary"].insert_one({"name": "default", "category": "dir", "size": "{:.2f}".format(size)})
    if result.inserted_id:
        await fs.upload_from_stream(
            str(result.inserted_id),  # 使用id作为文件名存储
            content  # 文件内容
        )

    # 更新子域名默认字典
    content = get_domainDict()
    size = len(content) / (1024 * 1024)
    result = await db["dictionary"].insert_one({"name": "default", "category": "subdomain", "size": "{:.2f}".format(size)})
    if result.inserted_id:
        await fs.upload_from_stream(
            str(result.inserted_id),  # 使用id作为文件名存储
            content  # 文件内容
        )
    # 获取任务名称
    cursor = db['task'].find({"type": {"$ne": "other"}})
    task_list = {}
    result = await cursor.to_list(length=None)
    for item in result:
        task_list[str(item["id"])] = item['name']
    # 修改资产字段
    cursor = db['asset'].find({})
    result = await cursor.to_list(length=None)
    for item in result:
        taskName = ""
        time = item["timestamp"]
        if item["taskId"] in task_list:
            taskName = task_list[item["taskId"]]
        if item["type"] != "other":
            ip = item["host"]
            parsed_url = urlparse(item['url'])
            host = parsed_url.netloc
            extracted = tldextract.extract(item['url'])
            root_domain = f"{extracted.domain}.{extracted.suffix}"
            await db['asset'].update_many({"_id": item["_id"]}, [{'$set': {'host': host, "ip": ip, "taskName": taskName, "rootDomain": root_domain, "time": time, "service": item["type"], "tag": []}}, {'$unset': 'timestamp'}])
        else:
            service = item["protocol"]
            extracted = tldextract.extract("https://" + item['host'])
            root_domain = f"{extracted.domain}.{extracted.suffix}"
            await db['asset'].update_many({"_id": item["_id"]}, [{'$set': {"taskName": taskName, "rootDomain": root_domain, "time": time, "service": service, "tag": []}}, {'$unset': 'timestamp'}, {'$unset': 'protocol'}])

    # 修改敏感信息字段
    cursor = db['SensitiveResult'].find({})
    update_list =[]
    result = await cursor.to_list(length=None)
    for item in result:
        taskName = task_list[item["taskId"]]
        extracted = tldextract.extract(item['url'])
        root_domain = f"{extracted.domain}.{extracted.suffix}"
        if item["body"] is not None or item["body"] != "":
            if item['md5'] not in update_list:
                await db['SensitiveBody'].update_one({"md5": item['md5']}, {"body": item['body']}, upsert=True)
                update_list.append(item['md5'])
            await db['SensitiveResult'].update_one({"_id": item['_id']}, {"taskName": taskName, "rootDomain": root_domain}, {'$unset': 'body'})
        else:
            await db['SensitiveResult'].update_one({"_id": item['_id']}, {'$set': {"md5": item['md5'].replace("md5==", ""), "taskName": taskName, "rootDomain": root_domain}}, {'$unset': 'body'})
    # 修改任务字段
    doc_names = ["DirScanResult", "PageMonitoring",
                 "SubdoaminTakerResult",
                 "UrlScan",
                 "crawler",
                 "subdomain"]
    for doc_name in doc_names:
        bulk_operations = []

        # 为每个 taskId 创建相应的更新操作
        for task_id, task_name in task_list.items():
            query = {"taskId": task_id}
            update_query = {"$set": {"taskName": task_name}}
            bulk_operations.append(UpdateOne(query, update_query))

        # 如果有更新操作，执行批量更新
        if bulk_operations:
            await db[doc_name].bulk_write(bulk_operations)

    # 增加默认插件
    await db["plugins"].insert_many(PLUGINS)

    # 更新POC严重等级
    level_map = {
        6: 'critical',
        5: 'high',
        4: 'medium',
        3: 'low',
        2: 'info',
        1: 'unknown'
    }
    for value, label in level_map.items():
        db['PocList'].update_many(
            {"level": value},
            {"$set": {"level": label, "type": 'nuclei'}}
        )

    # 创建默认扫描模板
    collection = db["ScanTemplates"]
    await collection.insert_one(SCANTEMPLATE)

    # 修改页面监控数据
    # 创建页面监控文档，url不重复
    await db['PageMonitoring'].rename('PageMonitoring_bak')
    await db['PageMonitoring'].create_index([('url', ASCENDING)], unique=True)
    await db['PageMonitoringBody'].create_index([('md5', ASCENDING)], unique=True)


    # 修改全局线程配置、节点配置
