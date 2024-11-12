# -------------------------------------
# @file      : statistics.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/11/12 18:58
# -------------------------------------------
import asyncio
import traceback

from bson import ObjectId
from fastapi import APIRouter, Depends
from api.users import verify_token
from motor.motor_asyncio import AsyncIOMotorCursor
from core.util import *
from pymongo import DESCENDING
from loguru import logger

router = APIRouter()


@router.get("/data")
async def asset_statistics_data(db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    # 使用 asyncio.gather 并行执行估计计数查询
    aset_count_task = db['asset'].estimated_document_count()
    subdomain_count_task = db['subdomain'].estimated_document_count()
    sensitive_count_task = db['SensitiveResult'].estimated_document_count()
    url_count_task = db['UrlScan'].estimated_document_count()
    vulnerability_count_task = db['vulnerability'].estimated_document_count()

    # 等待所有计数任务完成
    aset_count, subdomain_count, sensitive_count, url_count, vulnerability_count = await asyncio.gather(
        aset_count_task, subdomain_count_task, sensitive_count_task, url_count_task, vulnerability_count_task
    )

    return {
        "code": 200,
        "data": {
            "assetCount": aset_count,
            "subdomainCount": subdomain_count,
            "sensitiveCount": sensitive_count,
            "urlCount": url_count,
            "vulnerabilityCount": vulnerability_count
        }
    }


@router.post("/port")
async def asset_data_statistics_port(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    query = await get_search_query("asset", request_data)
    if query == "":
        return {"message": "Search condition parsing error", "code": 500}
    pipeline = [
        {
            "$match": query  # 添加搜索条件
        },
        {
            "$facet": {
                "by_port": [
                    {"$group": {"_id": "$port", "num_tutorial": {"$sum": 1}}},
                    {"$match": {"_id": {"$ne": None}}}
                ]
            }
        }
    ]
    result = await db['asset'].aggregate(pipeline).to_list(None)
    result_list = {"Port": []}
    port_list = {}

    for r in result:
        for port in r['by_port']:
            port_list[port["_id"]] = port["num_tutorial"]

    port_list = dict(sorted(port_list.items(), key=lambda item: -item[1]))
    for port in port_list:
        result_list['Port'].append({"value": port, "number": port_list[port]})
    return {
        "code": 200,
        "data": result_list
    }


@router.post("/title")
async def asset_data_statistics_title(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    request_data['filter']['type'] = ['https', 'http']
    query = await get_search_query("asset", request_data)
    if query == "":
        return {"message": "Search condition parsing error", "code": 500}
    pipeline = [
        {
            "$match": query  # 添加搜索条件
        },
        {
            "$facet": {
                "by_title": [
                    {"$group": {"_id": "$title", "num_tutorial": {"$sum": 1}}},
                    {"$match": {"_id": {"$ne": ""}}}
                ]
            }
        }
    ]
    result = await db['asset'].aggregate(pipeline).to_list(None)
    result_list = {"Title": []}
    title_list = {}

    for r in result:
        for port in r['by_title']:
            title_list[port["_id"]] = port["num_tutorial"]

    title_list = dict(sorted(title_list.items(), key=lambda item: -item[1]))
    for title in title_list:
        result_list['Title'].append({"value": title, "number": title_list[title]})
    return {
        "code": 200,
        "data": result_list
    }


@router.post("/type")
async def asset_data_statistics_type(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    query = await get_search_query("asset", request_data)
    if query == "":
        return {"message": "Search condition parsing error", "code": 500}
    pipeline = [
        {
            "$match": query  # 添加搜索条件
        },
        {
            "$facet": {
                "by_service": [
                    {"$group": {"_id": "$service", "num_tutorial": {"$sum": 1}}},
                    {"$match": {"_id": {"$ne": None}}}
                ]
            }
        }
    ]
    result = await db['asset'].aggregate(pipeline).to_list(None)
    result_list = {"Service": []}
    service_list = {}
    for r in result:
        for p in r['by_service']:
            service_list[p['_id']] = p['num_tutorial']
    service_list = dict(sorted(service_list.items(), key=lambda item: -item[1]))
    for service in service_list:
        result_list['Service'].append({"value": service, "number": service_list[service]})
    return {
        "code": 200,
        "data": result_list
    }


@router.post("/icon")
async def asset_data_statistics_icon(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    query = await get_search_query("asset", request_data)
    if query == "":
        return {"message": "Search condition parsing error", "code": 500}
    pipeline = [
        {
            "$match": query  # 添加搜索条件
        },
        {
            "$project": {
                "faviconmmh3": 1,
                "iconcontent": 1
            }
        },
        {
            "$facet": {
                "by_icon": [
                    {"$group": {"_id": "$faviconmmh3",
                                "num_tutorial": {"$sum": 1},
                                "iconcontent": {"$first": "$iconcontent"}
                                }
                     },
                    {"$match": {"_id": {"$ne": None}}}
                ]
            }
        }
    ]
    result = await db['asset'].aggregate(pipeline).to_list(None)
    result_list = {"Icon": []}
    icon_list = {}
    icon_tmp = {}
    for r in result:
        for icon in r['by_icon']:
            if icon['_id'] != "":
                icon_tmp[icon['_id']] = icon['iconcontent']
                icon_list[icon['_id']] = icon['num_tutorial']
    icon_list = dict(sorted(icon_list.items(), key=lambda item: -item[1]))
    for ic in icon_list:
        result_list['Icon'].append({"value": icon_tmp[ic], "number": icon_list[ic], "icon_hash": ic})

    return {
        "code": 200,
        "data": result_list
    }


@router.post("/app")
async def asset_data_statistics_app(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    query = await get_search_query("asset", request_data)
    if query == "":
        return {"message": "Search condition parsing error", "code": 500}
    pipeline = [
        {
            "$match": query  # 添加搜索条件
        },
        {
            "$facet": {
                "by_webfinger": [
                    {"$unwind": "$webfinger"},
                    {"$group": {"_id": "$webfinger", "num_tutorial": {"$sum": 1}}},
                    {"$match": {"_id": {"$ne": None}}}
                ],
                "by_technologies": [
                    {"$unwind": "$technologies"},
                    {"$group": {"_id": "$technologies", "num_tutorial": {"$sum": 1}}},
                    {"$match": {"_id": {"$ne": None}}}
                ]
            }
        }
    ]
    result = await db['asset'].aggregate(pipeline).to_list(None)
    result_list = {"Product": []}
    tec_list = {}
    for r in result:
        for technologie in r['by_technologies']:
            tec_list[technologie['_id']] = technologie['num_tutorial']
        for webfinger in r['by_webfinger']:
            try:
                if APP[webfinger['_id']] not in tec_list:
                    tec_list[APP[webfinger['_id']]] = webfinger['num_tutorial']
                else:
                    tec_list[APP[webfinger['_id']]] += webfinger['num_tutorial']
            except:
                pass
    tec_list = dict(sorted(tec_list.items(), key=lambda item: -item[1]))
    for tec in tec_list:
        result_list['Product'].append({"value": tec, "number": tec_list[tec]})
    return {
        "code": 200,
        "data": result_list
    }
