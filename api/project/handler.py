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
    icp_regex = []
    cmp_list = []
    app_name_list = []
    app_id_list = []
    for item in root_domain:
        if "CMP:" in item or "ICP:" in item or "APP:" in item or "APP-ID:" in item:
            if "ICP:" in item:
                icp_regex.append(f"^{re.escape(item.replace('ICP:', ''))}")
            elif "CMP:" in item:
                cmp_list.append(item.replace('CMP:', ''))
            elif "APP:" in item:
                app_name_list.append(re.escape(item.replace("APP:", '')))
            elif "APP-ID:" in item:
                app_id_list.append(item.replace("APP-ID:", ''))
    async for db in get_mongo_db():
        for a in asset_collection_list:
            if change:
                await asset_update_project(root_domain, asset_collection_list[a], a, db, project_id, icp_regex, cmp_list, app_name_list, app_id_list)
            else:
                await asset_add_project(root_domain, a, db, project_id, icp_regex, cmp_list, app_name_list, app_id_list)


async def asset_add_project(root_domain, doc_name, db, project_id, icp_regex, cmp_list, app_name_list, app_id_list):
    # 构建查询条件
    if doc_name == "RootDomain":
        query = {
                    "$or": [
                        {"domain": {"$in": root_domain}},
                        {"company": {"$in": cmp_list}},
                        {"icp": {"$regex": "|".join(icp_regex), "$options": "i"}}
                    ]
                }
    elif doc_name == "app":
        query = {
            "$or": [
                {"name": {"$in": app_name_list}},
                {"bundleID": {"$in": app_id_list}},
                {"company": {"$in": cmp_list}},
                {"icp": {"$regex": "|".join(icp_regex), "$options": "i"}},
            ]
        }
    elif doc_name == "mp":
        query = {
            "$or": [
                {"company": {"$in": cmp_list}},
                {"icp": {"$regex": "|".join(icp_regex), "$options": "i"}}
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


async def asset_update_project(root_domain, db_key, doc_name, db, project_id, icp_regex, cmp_list, app_name_list, app_id_list):
    # 构建查询条件
    if doc_name == "RootDomain":
        query = {
            "$and": [
                {"project": project_id},
                {"domain": {"$nin": root_domain}},
                {"company": {"$nin": cmp_list}},
                {
                    "icp": {
                        "$not": {"$regex": "|".join(icp_regex), "$options": "i"}
                    }
                }
            ]
        }
    elif doc_name == "app":
        query = {
            "$and": [
                {"project": project_id},
                {"company": {"$nin": cmp_list}},
                {"bundleID": {"$nin": app_id_list}},
                {
                    "name": {
                        "$not": {"$regex": "|".join(app_name_list), "$options": "i"}
                    }
                },
                {
                    "icp": {
                        "$not": {"$regex": "|".join(icp_regex), "$options": "i"}
                    }
                }
            ]
        }
    elif doc_name == "mp":
        query = {
            "$and": [
                {"project": project_id},
                {"name": {"$nin": root_domain}},
                {"company": {"$nin": cmp_list}},
                {
                    "icp": {
                        "$not": {"$regex": "|".join(icp_regex), "$options": "i"}
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
    await asset_add_project(root_domain, doc_name, db, project_id, icp_regex, cmp_list, app_name_list, app_id_list)


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



