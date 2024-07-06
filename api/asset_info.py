# -------------------------------------
# @file      : assetinfo.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/4/14 17:14
# -------------------------------------------
import json

from bson import ObjectId
from fastapi import APIRouter, Depends
from api.users import verify_token
from motor.motor_asyncio import AsyncIOMotorCursor
from core.db import get_mongo_db
from core.redis_handler import get_redis_pool
from core.util import *
from pymongo import ASCENDING, DESCENDING
from loguru import logger

router = APIRouter()


@router.get("/asset/statistics/data")
async def asset_statistics_data(db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    aset_count = await db['asset'].count_documents({})
    subdomain_count = await db['subdomain'].count_documents({})
    sensitive_count = await db['SensitiveResult'].count_documents({})
    url_count = await db['UrlScan'].count_documents({})
    vulnerability_count = await db['vulnerability'].count_documents({})
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


@router.post("/asset/data")
async def asset_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        if len(APP) == 0:
            collection = db["FingerprintRules"]
            cursor = await collection.find({}, {"_id": 1, "name": 1})
            async for document in cursor:
                document['id'] = str(document['_id'])
                del document['_id']
                APP[document['id']] = document['name']
        if len(SensitiveRuleList) == 0:
            collection = db["SensitiveRule"]
            cursor = await collection.find({}, {"_id": 1, "name": 1})
            async for document in cursor:
                document['id'] = str(document['_id'])
                del document['_id']
                SensitiveRuleList[document['id']] = {
                    "name": document['name'],
                    "color": document['color']
                }
        page_index = request_data.get("pageIndex", 1)
        page_size = request_data.get("pageSize", 10)
        query = await get_search_query("asset", request_data)
        if query == "":
            return {"message": "Search condition parsing error", "code": 500}
        total_count = await db['asset'].count_documents(query)
        cursor: AsyncIOMotorCursor = ((db['asset'].find(query, {"_id": 0,
                                                                "id": {"$toString": "$_id"},
                                                                "host": 1,
                                                                "url": 1,
                                                                "ip": 1,
                                                                "port": 1,
                                                                "protocol": 1,
                                                                "type": 1,
                                                                "title": 1,
                                                                "statuscode": 1,
                                                                "rawheaders": 1,
                                                                "webfinger": 1,
                                                                "technologies": 1,
                                                                "raw": 1,
                                                                "timestamp": 1,
                                                                "iconcontent": 1
                                                                })
                                       .skip((page_index - 1) * page_size)
                                       .limit(page_size))
                                      .sort([("timestamp", DESCENDING)]))
        result = await cursor.to_list(length=None)
        result_list = []
        for r in result:
            tmp = {}
            tmp['port'] = r['port']
            tmp['time'] = r['timestamp']
            tmp['id'] = r['id']
            tmp['type'] = r['type']
            if r['type'] == 'other':
                tmp['domain'] = r['host']
                tmp['ip'] = r['ip']
                tmp['service'] = r['protocol']
                tmp['title'] = ""
                tmp['status'] = None
                tmp['banner'] = ""
                try:
                    if r['raw'] is not None:
                        raw_data = json.loads(r['raw'].decode('utf-8'))
                        for k in raw_data:
                            tmp['banner'] += k + ":" + str(raw_data[k]).strip("\n") + "\n"
                except:
                    tmp['banner'] = ""
                tmp['products'] = []
            else:
                tmp['domain'] = r['url'].replace(f'{r["type"]}://', '')
                tmp['ip'] = r['host']
                tmp['service'] = r['type']
                tmp['title'] = r['title']
                tmp['status'] = r['statuscode']
                tmp['url'] = r['url']
                tmp['banner'] = r['rawheaders']
                tmp['products'] = []
                tmp['icon'] = r['iconcontent']
                technologies = r['technologies']
                if technologies is not None:
                    tmp['products'] = tmp['products'] + technologies
                if r['webfinger'] is not None:
                    for w in r['webfinger']:
                        tmp['products'].append(APP[w])
            result_list.append(tmp)
        return {
            "code": 200,
            "data": {
                'list': result_list,
                'total': total_count
            }
        }
    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/asset/detail")
async def asset_detail(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # Get the ID from the request data
        asset_id = request_data.get("id")

        # Check if ID is provided
        if not asset_id:
            return {"message": "ID is missing in the request data", "code": 400}

        # Query the database for content based on ID
        query = {"_id": ObjectId(asset_id)}
        doc = await db.asset.find_one(query)

        if not doc:
            return {"message": "Content not found for the provided ID", "code": 404}
        products = []
        tlsdata = ""
        hashes = ""
        banner = ""
        if doc['type'] == 'other':
            domain = doc.get('host', "")
            IP = doc.get("ip", "")
            URL = ""
            service = doc.get("protocol", "")
            try:
                if doc['raw'] is not None:
                    raw_data = json.loads(doc['raw'].decode('utf-8'))
                    for k in raw_data:
                        banner += k + ":" + str(raw_data[k]).strip("\n") + "\n"
            except:
                banner = ""
        else:
            domain = doc.get('url', "").replace("http://", "").replace("https://", "").split(":")[0]
            IP = doc.get("host", "")
            URL = doc.get("url", "")
            service = doc.get("type", "")
            products = doc.get('technologies')
            if products == None:
                products = []
            if doc['webfinger'] is not None:
                for w in doc['webfinger']:
                    products.append(APP[w])
            if doc['tlsdata'] is not None:
                for h in doc['tlsdata']:
                    tlsdata += h + ": " + str(doc['tlsdata'][h]) + '\n'
            if doc['hashes'] is not None:
                for h in doc['hashes']:
                    hashes += h + ": " + str(doc['hashes'][h]) + '\n'
            banner = doc.get('rawheaders', "")
        project_name = ""
        if doc.get("project", "") != "":
            query = {"_id": ObjectId(doc.get("project", ""))}
            project_data = await db.project.find_one(query)
            project_name = project_data.get("name", "")
        data = {
            "host": domain,
            "IP": IP,
            "URL": URL,
            "port": doc.get("port", ""),
            "service": service,
            "title": doc.get("title", ""),
            "status": doc.get("statuscode", ""),
            "FaviconHash": doc.get("faviconmmh3", ""),
            "jarm": doc.get("jarm", ""),
            "time": doc.get("timestamp", ""),
            "products": products,
            "TLSData": tlsdata,
            "hash": hashes,
            "banner": banner,
            "ResponseBody": doc.get("responsebody", ""),
            "project": project_name
        }
        return {"code": 200, "data": data}
    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/asset/statistics")
async def asset_data_statistics(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    query = await get_search_query("asset", request_data)
    if query == "":
        return {"message": "Search condition parsing error", "code": 500}
    cursor: AsyncIOMotorCursor = ((db['asset'].find(query, {
        "port": 1,
        "protocol": 1,
        "type": 1,
        "webfinger": 1,
        "technologies": 1,
        "faviconmmh3": 1,
        "iconcontent": 1
    })))
    result = await cursor.to_list(length=None)
    result_list = {"Port": [], "Service": [], "Product": [], "Icon": []}
    port_list = {}
    service_list = {}
    icon_list = {}
    icon_tmp = {}
    tec_list = {}
    for r in result:
        if r['port'] not in port_list:
            port_list[r['port']] = 1
        else:
            port_list[r['port']] += 1
        if r['type'] == "http" or r['type'] == "https":
            service = r['type']
            icon = r['iconcontent']
            icon_hash = r['faviconmmh3']
            if icon_hash != "":
                icon_tmp[icon_hash] = icon
                if icon_hash not in icon_list:
                    icon_list[icon_hash] = 1
                else:
                    icon_list[icon_hash] += 1
            if r['technologies'] != None:
                for t in r['technologies']:
                    if t != "":
                        if t not in tec_list:
                            tec_list[t] = 1
                        else:
                            tec_list[t] += 1
            if r['webfinger'] != None:
                for wf in r['webfinger']:
                    if wf != None:
                        if APP[wf] not in tec_list:
                            tec_list[APP[wf]] = 1
                        else:
                            tec_list[APP[wf]] += 1
        else:
            service = r['protocol']
        if service != "":
            if service not in service_list:
                service_list[service] = 1
            else:
                service_list[service] += 1
    service_list = dict(sorted(service_list.items(), key=lambda item: -item[1]))
    for service in service_list:
        result_list['Service'].append({"value": service, "number": service_list[service]})
    port_list = dict(sorted(port_list.items(), key=lambda item: -item[1]))
    for port in port_list:
        result_list['Port'].append({"value": port, "number": port_list[port]})
    tec_list = dict(sorted(tec_list.items(), key=lambda item: -item[1]))
    for tec in tec_list:
        result_list['Product'].append({"value": tec, "number": tec_list[tec]})
    icon_list = dict(sorted(icon_list.items(), key=lambda item: -item[1]))
    for ic in icon_list:
        result_list['Icon'].append({"value": icon_tmp[ic], "number": icon_list[ic], "icon_hash": ic})

    return {
        "code": 200,
        "data": result_list
    }


@router.post("/subdomain/data")
async def asset_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        page_index = request_data.get("pageIndex", 1)
        page_size = request_data.get("pageSize", 10)
        query = await get_search_query("subdomain", request_data)
        if query == "":
            return {"message": "Search condition parsing error", "code": 500}
        total_count = await db['subdomain'].count_documents(query)
        cursor: AsyncIOMotorCursor = ((db['subdomain'].find(query, {"_id": 0,
                                                                    "id": {"$toString": "$_id"},
                                                                    "host": 1,
                                                                    "type": 1,
                                                                    "value": 1,
                                                                    "ip": 1,
                                                                    "time": 1,
                                                                    })
                                       .skip((page_index - 1) * page_size)
                                       .limit(page_size))
                                      .sort([("time", DESCENDING)]))
        result = await cursor.to_list(length=None)
        result_list = []
        for r in result:
            if r['value'] is None:
                r['value'] = []
            if r['ip'] is None:
                r['ip'] = []
            result_list.append(r)
        return {
            "code": 200,
            "data": {
                'list': result_list,
                'total': total_count
            }
        }
    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/url/data")
async def url_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        page_index = request_data.get("pageIndex", 1)
        page_size = request_data.get("pageSize", 10)
        query = await get_search_query("url", request_data)
        if query == "":
            return {"message": "Search condition parsing error", "code": 500}
        total_count = await db['UrlScan'].count_documents(query)
        cursor: AsyncIOMotorCursor = ((db['UrlScan'].find(query, {"_id": 0,
                                                                  "id": {"$toString": "$_id"},
                                                                  "input": 1,
                                                                  "source": 1,
                                                                  "type": "$outputtype",
                                                                  "url": "$output",
                                                                  "time": 1,
                                                                  })
                                       .skip((page_index - 1) * page_size)
                                       .limit(page_size))
                                      .sort([("time", DESCENDING)]))
        result = await cursor.to_list(length=None)
        return {
            "code": 200,
            "data": {
                'list': result,
                'total': total_count
            }
        }
    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/crawler/data")
async def crawler_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        page_index = request_data.get("pageIndex", 1)
        page_size = request_data.get("pageSize", 10)
        query = await get_search_query("crawler", request_data)
        if query == "":
            return {"message": "Search condition parsing error", "code": 500}
        total_count = await db['crawler'].count_documents(query)
        cursor: AsyncIOMotorCursor = ((db['crawler'].find(query, {"_id": 0,
                                                                  "id": {"$toString": "$_id"},
                                                                  "method": 1,
                                                                  "body": 1,
                                                                  "url": 1
                                                                  })
                                       .sort([('_id', -1)])
                                       .skip((page_index - 1) * page_size)
                                       .limit(page_size))
        )
        result = await cursor.to_list(length=None)
        return {
            "code": 200,
            "data": {
                'list': result,
                'total': total_count
            }
        }
    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/asset/statistics2")
async def asset_data_statistics2(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    query = await get_search_query("asset", request_data)
    if query == "":
        return {"message": "Search condition parsing error", "code": 500}
    pipeline = [
        {
            "$match": query  # 添加搜索条件
        },
        {
            "$facet": {
                "by_type": [
                    {"$group": {"_id": "$type", "num_tutorial": {"$sum": 1}}},
                    {"$match": {"_id": {"$ne": None}}}
                ],
                "by_port": [
                    {"$group": {"_id": "$port", "num_tutorial": {"$sum": 1}}},
                    {"$match": {"_id": {"$ne": None}}}
                ],
                "by_protocol": [
                    {"$group": {"_id": "$protocol", "num_tutorial": {"$sum": 1}}},
                    {"$match": {"_id": {"$ne": None}}}
                ],
                "by_icon": [
                    {"$group": {"_id": "$faviconmmh3", "num_tutorial": {"$sum": 1},
                                "iconcontent": {"$first": "$iconcontent"}}},
                    {"$match": {"_id": {"$ne": ""}}}
                ],
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
    result_list = {"Port": [], "Service": [], "Product": [], "Icon": []}
    port_list = {}
    service_list = {}
    icon_list = {}
    icon_tmp = {}
    tec_list = {}
    for r in result:
        for port in r['by_port']:
            port_list[port["_id"]] = port["num_tutorial"]
        for icon in r['by_icon']:
            icon_tmp[icon['_id']] = icon['iconcontent']
            icon_list[icon['_id']] = icon['num_tutorial']
        for type in r['by_type']:
            service_list[type['_id']] = type['num_tutorial']
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
    service_list = dict(sorted(service_list.items(), key=lambda item: -item[1]))
    for service in service_list:
        result_list['Service'].append({"value": service, "number": service_list[service]})
    port_list = dict(sorted(port_list.items(), key=lambda item: -item[1]))
    for port in port_list:
        result_list['Port'].append({"value": port, "number": port_list[port]})
    tec_list = dict(sorted(tec_list.items(), key=lambda item: -item[1]))
    for tec in tec_list:
        result_list['Product'].append({"value": tec, "number": tec_list[tec]})
    icon_list = dict(sorted(icon_list.items(), key=lambda item: -item[1]))
    for ic in icon_list:
        result_list['Icon'].append({"value": icon_tmp[ic], "number": icon_list[ic], "icon_hash": ic})

    return {
        "code": 200,
        "data": result_list
    }


@router.post("/asset/statistics/port")
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


@router.post("/asset/statistics/type")
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
                "by_type": [
                    {"$group": {"_id": "$type", "num_tutorial": {"$sum": 1}}},
                    {"$match": {"_id": {"$ne": None}}}
                ],
                "by_protocol": [
                    {"$group": {"_id": "$protocol", "num_tutorial": {"$sum": 1}}},
                    {"$match": {"_id": {"$ne": None}}}
                ]
            }
        }
    ]
    result = await db['asset'].aggregate(pipeline).to_list(None)
    result_list = {"Service": []}
    service_list = {}
    for r in result:
        for t in r['by_type']:
            if t['_id'] != 'other':
                service_list[t['_id']] = t['num_tutorial']
        for p in r['by_protocol']:
            service_list[p['_id']] = p['num_tutorial']
    service_list = dict(sorted(service_list.items(), key=lambda item: -item[1]))
    for service in service_list:
        result_list['Service'].append({"value": service, "number": service_list[service]})
    return {
        "code": 200,
        "data": result_list
    }


@router.post("/asset/statistics/icon")
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


# @router.post("/asset/statistics/icon2")
# async def asset_data_statistics_icon2(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
#     search_query = request_data.get("search", "")
#     keyword = {
#         'app': '',
#         'body': 'responsebody',
#         'header': 'rawheaders',
#         'project': 'project',
#         'title': 'title',
#         'statuscode': 'statuscode',
#         'icon': 'faviconmmh3',
#         'ip': ['host', 'ip'],
#         'domain': ['host', 'url', 'domain'],
#         'port': 'port',
#         'protocol': ['protocol', 'type'],
#         'banner': 'raw',
#     }
#     query = await search_to_mongodb(search_query, keyword)
#     if query == "" or query is None:
#         return {"message": "Search condition parsing error", "code": 500}
#     query = query[0]
#     query["faviconmmh3"] = {"$ne": ""}
#     cursor = db.asset.find(query, {"_id": 0,
#                                    "faviconmmh3": 1,
#                                    "iconcontent": 1
#                                    })
#     results = await cursor.to_list(length=None)
#     result_list = {"Icon": []}
#     icon_list = {}
#     icon_tmp = {}
#     for r in results:
#         if r['faviconmmh3'] not in icon_list:
#             r['faviconmmh3'] = 1
#             icon_tmp[r['faviconmmh3']] = r['iconcontent']
#         else:
#             r['faviconmmh3'] += 1
#     icon_list = dict(sorted(icon_list.items(), key=lambda item: -item[1]))
#     for ic in icon_list:
#         result_list['Icon'].append({"value": icon_tmp[ic], "number": icon_list[ic], "icon_hash": ic})
#
#     return {
#         "code": 200,
#         "data": result_list
#     }


@router.post("/asset/statistics/app")
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


@router.post("/data/delete")
async def delete_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # Extract the list of IDs from the request_data dictionary
        data_ids = request_data.get("ids", [])
        index = request_data.get("index", "")
        # Convert the provided rule_ids to ObjectId
        obj_ids = [ObjectId(data_id) for data_id in data_ids]

        # Delete the SensitiveRule documents based on the provided IDs
        result = await db[index].delete_many({"_id": {"$in": obj_ids}})

        # Check if the deletion was successful
        if result.deleted_count > 0:
            return {"code": 200, "message": "Data deleted successfully"}
        else:
            return {"code": 404, "message": "Data not found"}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}