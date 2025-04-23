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
from pymongo import ASCENDING, UpdateMany, UpdateOne

from core.config import VERSION
from core.default import get_dirDict, get_domainDict, get_sensitive, ModulesConfig, PLUGINS, SCANTEMPLATE
from core.util import print_progress_bar


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

    fs = AsyncIOMotorGridFSBucket(db)
    total_steps = 11
    # 更新目录扫描默认字典
    content = get_dirDict()
    if content != "":
        size = len(content) / (1024 * 1024)
        result = await db["dictionary"].insert_one({"name": "default", "category": "dir", "size": "{:.2f}".format(size)})
        if result.inserted_id:
            await fs.upload_from_stream(
                str(result.inserted_id),  # 使用id作为文件名存储
                content.encode('utf-8')  # 文件内容
            )
    print_progress_bar(1, 11)
    # 更新子域名默认字典
    content = get_domainDict()
    if content != "":
        size = len(content) / (1024 * 1024)
        result = await db["dictionary"].insert_one({"name": "default", "category": "subdomain", "size": "{:.2f}".format(size)})
        if result.inserted_id:
            await fs.upload_from_stream(
                str(result.inserted_id),  # 使用id作为文件名存储
                content.encode('utf-8')  # 文件内容
            )
    print_progress_bar(2, 11)
    # 获取任务名称
    cursor = db['task'].find({"type": {"$ne": "other"}})
    task_list = {}
    result = await cursor.to_list(length=None)
    for item in result:
        task_list[str(item["_id"])] = item['name']

    print_progress_bar(3, 11)

    # 修改资产字段
    cursor = db['asset'].find({})
    result = await cursor.to_list(length=None)
    # 构建批量操作列表
    bulk_operations = []
    for item in result:
        try:
            taskName = ""
            if "timestamp" in item:
                time = item["timestamp"]
            else:
                if "time" in item:
                    time = item["time"]
                else:
                    time = ""
            if item.get("taskId", "") in task_list:
                taskName = task_list[item["taskId"]]

            if item["type"] != "other":
                ip = item["host"]
                parsed_url = urlparse(item['url'])
                host = parsed_url.netloc
                extracted = tldextract.extract(item['url'])
                root_domain = f"{extracted.domain}.{extracted.suffix}"
                bulk_operations.append(UpdateOne(
                    {"_id": item["_id"]},
                    {'$set': {'host': host, "ip": ip, "taskName": taskName, "rootDomain": root_domain, "time": time,
                              "service": item["type"], "tag": []}},
                    upsert=False
                ))
            else:
                service = item["protocol"]
                extracted = tldextract.extract("https://" + item['host'])
                root_domain = f"{extracted.domain}.{extracted.suffix}"
                bulk_operations.append(UpdateOne(
                    {"_id": item["_id"]},
                    {'$set': {"taskName": taskName, "rootDomain": root_domain, "time": time, "service": service,
                              "tag": []}},
                    upsert=False
                ))
        except:
            logger.warning(f"asset 更新失败： {item}")
    # 批量执行更新操作
    if bulk_operations:
        await db['asset'].bulk_write(bulk_operations)
    print_progress_bar(4, 11)


    # # 修改敏感信息字段
    # cursor = db['SensitiveResult'].find({})
    # result = await cursor.to_list(length=None)
    #
    # update_sensitive_result_ops = []
    # update_sensitive_body_ops = []
    # update_list = set()  # 使用集合提高查重效率
    #
    # for item in result:
    #     taskName = task_list[item["taskId"]]
    #     extracted = tldextract.extract(item['url'])
    #     root_domain = f"{extracted.domain}.{extracted.suffix}"
    #
    #     if item["body"]:  # 仅在 body 不为空时处理
    #         if item['md5'] not in update_list:
    #             update_sensitive_body_ops.append(
    #                 UpdateOne({"md5": item['md5']}, {"$set": {"body": item['body']}}, upsert=True)
    #             )
    #             update_list.add(item['md5'])
    #
    #         update_sensitive_result_ops.append(
    #             UpdateOne({"_id": item['_id']},
    #                       {"$set": {"taskName": taskName, "rootDomain": root_domain}, "$unset": {"body": ""}})
    #         )
    #     else:
    #         update_sensitive_result_ops.append(
    #             UpdateOne({"_id": item['_id']},
    #                       {"$set": {"md5": item['md5'].replace("md5==", ""), "taskName": taskName,
    #                                 "rootDomain": root_domain},
    #                        "$unset": {"body": ""}})
    #         )
    #
    # # 批量写入 SensitiveBody
    # if update_sensitive_body_ops:
    #     await db['SensitiveBody'].bulk_write(update_sensitive_body_ops)
    #
    # # 批量写入 SensitiveResult
    # if update_sensitive_result_ops:
    #     await db['SensitiveResult'].bulk_write(update_sensitive_result_ops)
    try:
        await db['SensitiveResult'].rename('SensitiveResult_bak')
    except:
        logger.warning("SensitiveResult not found")
    logger.warning("敏感信息数据由于和新版本不兼容，备份到SensitiveResult_bak中，如果不需要可以连接数据库进行删除")
    logger.warning("Sensitive information data is incompatible with the new version, so it is backed up to SensitiveResult_bak. If it is not needed, you can connect to the database to delete it.")
    print_progress_bar(5, 11)
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
            bulk_operations.append(UpdateMany(query, update_query))

        # 如果有更新操作，执行批量更新
        if bulk_operations:
            await db[doc_name].bulk_write(bulk_operations)
    print_progress_bar(6, 11)

    # 增加默认插件
    await db["plugins"].insert_many(PLUGINS)
    print_progress_bar(7, 11)

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
        await db['PocList'].update_many(
            {"level": value},
            {"$set": {"level": label, "type": 'nuclei'}}
        )
    print_progress_bar(8, 11)
    # 创建默认扫描模板
    collection = db["ScanTemplates"]
    await collection.insert_one(SCANTEMPLATE)
    print_progress_bar(9, 11)
    # 修改页面监控数据
    # 创建页面监控文档，url不重复
    try:
        await db['PageMonitoring'].rename('PageMonitoring_bak')
    except:
        logger.warning("PageMonitoring not found")
    logger.warning("页面监控数据由于和新版本不兼容，备份到PageMonitoring_bak中，如果不需要可以连接数据库进行删除")
    logger.warning("Since the page monitoring data is incompatible with the new version, it is backed up to PageMonitoring_bak. If it is not needed, you can connect to the database to delete it.")
    await db['PageMonitoring'].create_index([('url', ASCENDING)], unique=True)
    await db['PageMonitoringBody'].create_index([('md5', ASCENDING)], unique=True)
    # 修改全局线程配置、节点配置
    await db.config.insert_one(
        {"name": "ModulesConfig", 'value': ModulesConfig, 'type': 'system'})
    print_progress_bar(10, 11)
    # 增加任务运行状态
    # 运行中
    await db.task.update_many({"progress": {"$ne": 100}}, {"$set": {"status": 1}})
    # 完成
    await db.task.update_many({"progress": 100}, {"$set": {"status": 3}})
    print_progress_bar(11, 11)

    logger.success("更新完成")
    logger.success("Update Complete")


async def update16(db):
    # 创建asset索引
    await db['asset'].create_index([('time', ASCENDING)])
    await db['asset'].create_index([('url', ASCENDING)])
    await db['asset'].create_index([('host', ASCENDING)])
    await db['asset'].create_index([('ip', ASCENDING)])
    await db['asset'].create_index([('port', ASCENDING)])
    await db['asset'].create_index([('host', ASCENDING), ('port', ASCENDING)])


async def update17(db):
    await db["plugins"].insert_many([{
        "module": "URLSecurity",
        "name": "trufflehog",
        "hash": "1aa212b9578dc3fb1409ee8de8ed005e",
        "parameter": "-pdf false -verify false",
        "help": "-pdf 开启pdf检测 -exclude 排除提取的规则(name1,name2) -verify 是否进行验证（验证通过再统计结果）",
        "introduction": "trufflehog密钥提取，如果设置了排除规则，需要重新安装才能重新启用被排除的规则",
        "isSystem": True,
        "version": "1.0",
        "source": ""
    }])
    await db['asset'].create_index([('project', ASCENDING)])
    await db['asset'].create_index([('taskName', ASCENDING)])
    await db['asset'].create_index([('rootDomain', ASCENDING)])

    await db['subdomain'].create_index([('project', ASCENDING)])
    await db['subdomain'].create_index([('taskName', ASCENDING)])
    await db['subdomain'].create_index([('rootDomain', ASCENDING)])
    await db['subdomain'].create_index([('time', ASCENDING)])

    await db['UrlScan'].create_index([('project', ASCENDING)])
    await db['UrlScan'].create_index([('taskName', ASCENDING)])
    await db['UrlScan'].create_index([('rootDomain', ASCENDING)])

    await db['crawler'].create_index([('project', ASCENDING)])
    await db['crawler'].create_index([('taskName', ASCENDING)])
    await db['crawler'].create_index([('rootDomain', ASCENDING)])

    await db['SensitiveResult'].create_index([('project', ASCENDING)])
    await db['SensitiveResult'].create_index([('taskName', ASCENDING)])
    await db['SensitiveResult'].create_index([('rootDomain', ASCENDING)])

    await db['DirScanResult'].create_index([('project', ASCENDING)])
    await db['DirScanResult'].create_index([('taskName', ASCENDING)])
    await db['DirScanResult'].create_index([('rootDomain', ASCENDING)])

    await db['vulnerability'].create_index([('project', ASCENDING)])
    await db['vulnerability'].create_index([('taskName', ASCENDING)])
    await db['vulnerability'].create_index([('rootDomain', ASCENDING)])

    await db['RootDomain'].create_index([('project', ASCENDING)])
    await db['RootDomain'].create_index([('taskName', ASCENDING)])
    await db['RootDomain'].create_index([('domain', ASCENDING)], unique=True)
    await db['RootDomain'].create_index([('time', ASCENDING)])

    await db['app'].create_index([('project', ASCENDING)])
    await db['app'].create_index([('taskName', ASCENDING)])
    await db['app'].create_index([('time', ASCENDING)])
    await db['app'].create_index([('name', ASCENDING)])

    await db['mp'].create_index([('project', ASCENDING)])
    await db['mp'].create_index([('taskName', ASCENDING)])
    await db['mp'].create_index([('time', ASCENDING)])
    await db['mp'].create_index([('name', ASCENDING)])