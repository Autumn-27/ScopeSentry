# -------------------------------------
# @file      : handler.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/10/29 21:00
# -------------------------------------------
import re

from core.db import get_mongo_db
from loguru import logger


async def update_project(root_domain, project_id, change=False):
    asset_collection_list = {
                        'asset': ["url", "host", "ip"],
                        'subdomain': ["host", "ip"],
                        'DirScanResult': ["url"],
                        'vulnerability': ["url"],
                        'SubdoaminTakerResult': ["input"],
                        'PageMonitoring': ["url"],
                        'SensitiveResult': ["url"],
                        'UrlScan': ["input"],
                        'crawler': ["url"],
                        'RootDomain': ['domain', 'icp', 'company'],
                        'app': ['name', 'bundleID', 'company', 'icp'],
                        'mp': ['company', 'icp']
    }
    async for db in get_mongo_db():
        for a in asset_collection_list:
            if change:
                await asset_update_project(root_domain, asset_collection_list[a], a, db, project_id)
            else:
                await asset_add_project(root_domain, a, db, project_id)


async def asset_add_project(root_domain, doc_name, db, project_id):
    # 构建查询条件
    if doc_name == "RootDomain":
        regex_patterns = [f"^{re.escape(item)}" for item in root_domain]
        query = {
                    "$or": [
                        {"domain": {"$in": root_domain}},
                        {"company": {"$in": root_domain}},
                        {"icp": {"$regex": "|".join(regex_patterns), "$options": "i"}}
                    ]
                }
    elif doc_name == "app":
        regex_patterns = [f"^{re.escape(item)}" for item in root_domain]
        query = {
            "$or": [
                {"name": {"$in": root_domain}},
                {"bundleID": {"$in": root_domain}},
                {"company": {"$in": root_domain}},
                {"icp": {"$regex": "|".join(regex_patterns), "$options": "i"}},
            ]
        }
    elif doc_name == "mp":
        regex_patterns = [f"^{re.escape(item)}" for item in root_domain]
        query = {
            "$or": [
                {"company": {"$in": root_domain}},
                {"icp": {"$regex": "|".join(regex_patterns), "$options": "i"}}
            ]
        }
    else:
        query = {
            "rootDomain": {"$in": root_domain}
        }
    update_query = {
        "$set": {
            "project": project_id
        }
    }
    result = await db[doc_name].update_many(query, update_query)
    # 打印更新的文档数量
    logger.info(f"Updated {doc_name} {result.modified_count} documents")


async def asset_update_project(root_domain, db_key, doc_name, db, project_id):
    # 构建查询条件
    if doc_name == "RootDomain":
        regex_patterns = [f"^{re.escape(item)}" for item in root_domain]
        query = {
            "$and": [
                {"project": project_id},
                {"domain": {"$nin": root_domain}},
                {"company": {"$nin": root_domain}},
                {
                    "icp": {
                        "$nin": root_domain,
                        "$not": {"$regex": "|".join(regex_patterns), "$options": "i"}
                    }
                }
            ]
        }
    elif doc_name == "app":
        regex_patterns = [f"^{re.escape(item)}" for item in root_domain]
        query = {
            "$and": [
                {"project": project_id},
                {"company": {"$nin": root_domain}},
                {
                    "icp": {
                        "$nin": root_domain,
                        "$not": {"$regex": "|".join(regex_patterns), "$options": "i"}
                    }
                }
            ]
        }
    elif doc_name == "mp":
        regex_patterns = [f"^{re.escape(item)}" for item in root_domain]
        query = {
            "$and": [
                {"project": project_id},
                {"name": {"$nin": root_domain}},
                {"company": {"$nin": root_domain}},
                {"bundleID": {"$nin": root_domain}},
                {
                    "icp": {
                        "$nin": root_domain,
                        "$not": {"$regex": "|".join(regex_patterns), "$options": "i"}
                    }
                }
            ]
        }
    else:
        query = {
                    "$and": [
                        {"project": project_id},
                        {"rootDomain": {"$nin": root_domain}}
                    ]
                }
    update_query = {
        "$set": {
            "project": ""
        }
    }
    result = await db[doc_name].update_many(query, update_query)
    # 打印更新的文档数量
    logger.info(f"Updated {doc_name} {result.modified_count} documents to null ")
    await asset_add_project(root_domain, doc_name, db, project_id)


async def delete_asset_project(db, collection, project_id):
    try:
        # 直接使用批量更新操作，减少单独更新的次数
        query = {"project": project_id}
        update = {"$set": {"project": ""}}

        result = await db[collection].update_many(query, update)

        logger.info(f"Matched {result.matched_count}, Modified {result.modified_count} documents.")
    except Exception as e:
        logger.error(f"delete_asset_project error: {e}")


async def delete_asset_project_handler(project_id):
    async for db in get_mongo_db():
        asset_collection_list = ['asset', 'subdomain', 'DirScanResult', 'vulnerability', 'SubdoaminTakerResult',
                                 'PageMonitoring', 'SensitiveResult', 'UrlScan', 'crawler']
        for c in asset_collection_list:
            await delete_asset_project(db, c, project_id)



