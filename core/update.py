# -------------------------------------
# @file      : update.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/6/16 14:14
# -------------------------------------------
from loguru import logger
from motor.motor_asyncio import AsyncIOMotorGridFSBucket

from core.config import VERSION
from core.default import get_dirDict, get_domainDict, get_sensitive, ModulesConfig


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
    # 更新目录扫描字典
    file_docs = await fs.files.find({"filename": "dirdict"}).to_list(1)
    if not file_docs:
        logger.error("File dirdict not found")
        return
    file_doc = file_docs[0]
    new_filename = "dir_default"
    # 更新文件名
    await fs.files.update_one(
        {"_id": file_doc["_id"]},
        {"$set": {"filename": new_filename}}
    )

    # 更新子域名字典
    file_docs = await fs.files.find({"filename": "DomainDic"}).to_list(1)
    if not file_docs:
        logger.error("File DomainDic not found")
    file_doc = file_docs[0]
    new_filename = "domain_default"
    # 更新文件名
    await fs.files.update_one(
        {"_id": file_doc["_id"]},
        {"$set": {"filename": new_filename}}
    )
