# -------------------------------------
# @file      : project_aggregation.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/7/8 21:02
# -------------------------------------------
import asyncio
import copy
import time
import traceback

from bson import ObjectId
from fastapi import APIRouter, Depends, BackgroundTasks
from pymongo import DESCENDING

from api.users import verify_token
from motor.motor_asyncio import AsyncIOMotorCursor

from core.config import Project_List
from core.db import get_mongo_db
from core.redis_handler import refresh_config, get_redis_pool
from loguru import logger
from core.util import *
from core.apscheduler_handler import scheduler

router = APIRouter()


@router.post("/project/info")
async def get_projects_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    id = request_data.get("id", "")
    result = await db.project.find_one({"_id": ObjectId(id)}, {
        "_id": 0,
        "tag": 1,
        "hour": 1,
        "scheduledTasks": 1,
        "AssetCount": 1,
        "root_domains": 1,
        "name": 1
    }
                                       )
    if result['scheduledTasks']:
        job = scheduler.get_job(id)
        if job is not None:
            next_time = job.next_run_time.strftime("%Y-%m-%d %H:%M:%S")
            result['next_time'] = next_time
    return {"code": 200, "data": result}


@router.post("/project/asset/count")
async def get_projects_asset_count(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    id = request_data.get("id", "")
    subdomain_count = await db['subdomain'].count_documents({"project": id})
    vulnerability_count = await db['vulnerability'].count_documents({"project": id})
    return {"code": 200, "data": {
        "subdomainCount": subdomain_count,
        "vulCount": vulnerability_count
    }}


@router.post("/project/vul/statistics")
async def get_projects_vul_statistics(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    id = request_data.get("id", "")
    pipeline = [
        {"$match": {"project": id}},
        {
            "$group": {
                "_id": "$level",
                "count": {"$sum": 1}
            }
        }
    ]
    result = await db['vulnerability'].aggregate(pipeline).to_list(None)
    return {"code": 200, "data": result}


@router.post("/project/vul/data")
async def get_projects_vul_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    id = request_data.get("id", "")
    cursor: AsyncIOMotorCursor = db.vulnerability.find({"project": id},
                                                       {"_id": 0, "url": 1, "vulname": 1, "level": 1, "time": 1,
                                                        "matched": 1}).sort([("time", DESCENDING)])
    result = await cursor.to_list(length=None)
    return {
        "code": 200,
        "data": {
            'list': result
        }
    }


async def process_domains(root_domains, query, db):
    cursor: AsyncIOMotorCursor = db['subdomain'].find(query, {
        "_id": 0, "id": {"$toString": "$_id"}, "host": 1, "type": 1, "value": 1, "ip": 1, "time": 1
    }).sort([("time", -1)])
    result = await cursor.to_list(length=None)

    domain_results = {root_domain: [] for root_domain in root_domains}
    for r in result:
        host = r.get('host', '')
        for root_domain in root_domains:
            if host.endswith(root_domain):
                domain_results[root_domain].append(r)
                break

    processed_results = []
    for root_domain, result_list in domain_results.items():
        for r in result_list:
            if r['value'] is None:
                r['value'] = []
            if r['ip'] is None:
                r['ip'] = []
        processed_results.append({
            "host": root_domain,
            "type": "",
            "value": [],
            "ip": [],
            "id": generate_random_string(5),
            "children": result_list,
            "count": len(result_list)
        })

    return processed_results


@router.post("/project/subdomain/data")
async def get_projects_subdomain_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    filter = request_data.get("filter", {})
    project_id = filter["project"][0]
    project_query = {}
    project_query["_id"] = ObjectId(project_id)
    doc = await db.project.find_one(project_query, {"_id": 0, "root_domains": 1})
    if not doc or "root_domains" not in doc:
        return {"code": 404, "message": "domain is null"}
    query = await get_search_query("subdomain", request_data)
    if query == "":
        return {"message": "Search condition parsing error", "code": 500}
    root_domains = doc["root_domains"]
    results = await process_domains(root_domains, query, db)
    return {
        "code": 200,
        "data": {
            'list': results
        }
    }


@router.post("/project/port/data")
async def get_projects_vul_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    query = await get_search_query("asset", request_data)
    if query == "":
        return {"message": "Search condition parsing error", "code": 500}
    pipeline = [
        {
            "$match": query  # 添加搜索条件
        },
        {
            "$group":
                {
                    "_id": "$port",
                    "count":
                        {
                            "$sum": 1
                        }
                }
        },
        {
            "$sort": {"count": -1}
        }
    ]
    result = await db['asset'].aggregate(pipeline).to_list(None)

    async def fetch_asset_data(r):
        tmp_result = {
            "port": r['_id'],
            "id": generate_random_string(5),
            "count": r['count']
        }

        query_copy = query.copy()
        query_copy["port"] = r['_id']

        cursor: AsyncIOMotorCursor = db['asset'].find(query_copy, {
            "_id": 0,
            "id": {"$toString": "$_id"},
            "host": 1,
            "ip": 1,
            "type": 1,
            "time": 1,
            "service": 1
        }).sort([("time", DESCENDING)])

        asset_result = await cursor.to_list(length=None)
        children_list = []

        for asset in asset_result:
            children_data = {}
            children_data['service'] = asset['service']
            children_data['host'] = asset['host']
            children_data['time'] = asset['time']
            children_data['ip'] = asset['ip']
            children_data['id'] = asset['id']
            children_list.append(children_data)

        if len(children_list) != 0:
            tmp_result["children"] = children_list

        return tmp_result

    tasks = [fetch_asset_data(r) for r in result]
    result_list = await asyncio.gather(*tasks)

    return {
        "code": 200,
        "data": {
            'list': result_list
        }
    }


@router.post("/project/service/data")
async def get_projects_service_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    query = await get_search_query("asset", request_data)
    if query == "":
        return {"message": "Search condition parsing error", "code": 500}

    pipeline = [
        {
            "$match": query
        },
        {
            "$group": {
                "_id": "$service",
                "count": {"$sum": 1}
            }
        },
        {
            "$sort": {"count": -1}
        }
    ]

    result = await db['asset'].aggregate(pipeline).to_list(None)

    async def fetch_asset_data(r):
        tmp_result = {
            "service": r['_id'],
            "count": r['count'],
            "id": generate_random_string(5),
            "host": "",
            "ip": "",
            "time": "",
            "port": ""
        }

        query_copy = query.copy()
        query_copy["service"] = r['_id']
        if r['_id'] == "":
            tmp_result['service'] = 'unknown'

        cursor: AsyncIOMotorCursor = db['asset'].find(query_copy, {
            "_id": 0,
            "id": {"$toString": "$_id"},
            "host": 1,
            "ip": 1,
            "service": 1,
            "time": 1,
            "webServer": 1,
            "port": 1,
            "type": 1
        }).sort([("time", DESCENDING)])

        asset_result = await cursor.to_list(length=None)
        children_list = []

        for asset in asset_result:
            children_data = {}

            if asset['type'] == "other":
                children_data['service'] = ''
            else:
                if 'webServer' in asset:
                    children_data['service'] = asset['webServer']
                else:
                    children_data['service'] = ''

            children_data['host'] = asset['host']
            children_data['ip'] = asset['ip']
            children_data['time'] = asset['time']
            children_data['port'] = asset['port']
            children_data['id'] = asset['id']
            children_list.append(children_data)

        if len(children_list) != 0:
            tmp_result["children"] = children_list

        return tmp_result

    tasks = [fetch_asset_data(r) for r in result]
    result_list = await asyncio.gather(*tasks)
    return {
        "code": 200,
        "data": {
            'list': result_list
        }
    }
