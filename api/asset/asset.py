# -------------------------------------
# @file      : asset.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/10/20 20:52
# -------------------------------------------
import json
import traceback

from bson import ObjectId
from fastapi import APIRouter, Depends
from api.users import verify_token
from motor.motor_asyncio import AsyncIOMotorCursor
from core.util import *
from pymongo import DESCENDING
from loguru import logger

router = APIRouter()


@router.get("/statistics/data")
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


@router.post("/data")
async def asset_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        page_index = request_data.get("pageIndex", 1)
        page_size = request_data.get("pageSize", 10)
        query = await get_search_query("asset", request_data)
        if query == "":
            return {"message": "Search condition parsing error", "code": 500}
        total_count = await db['asset'].count_documents(query)
        cursor = db['asset'].find(query, {"_id": 0,
                                          "id": {"$toString": "$_id"},
                                          "host": 1,
                                          "url": 1,
                                          "ip": 1,
                                          "port": 1,
                                          "service": 1,
                                          "type": 1,
                                          "title": 1,
                                          "statuscode": 1,
                                          "rawheaders": 1,
                                          "technologies": 1,
                                          "raw": 1,
                                          "time": 1,
                                          "iconcontent": 1,
                                          "tags": 1,
                                          "screenshot": 1,
                                          }).skip((page_index - 1) * page_size).limit(page_size).sort(
            [("time", DESCENDING)])
        result = await cursor.to_list(length=None)
        result_list = []
        for r in result:
            tmp = {
                'port': r['port'],
                'time': r['time'],
                'id': r['id'],
                'type': r['type'],
                'domain': r['host'],
                'ip': r['ip'],
                'service': r['service'],
                'tags': r.get("tags", [])
            }
            if tmp['tags'] == None:
                tmp['tags'] = []
            if r['type'] == 'other':
                tmp['title'] = ""
                tmp['status'] = None
                tmp['banner'] = ""
                try:
                    if r['raw'] is not None:
                        raw_data = json.loads(r['raw'].decode('utf-8'))
                        for k in raw_data:
                            tmp['banner'] += k + ":" + str(raw_data[k]).strip("\n") + "\n"
                except:
                    try:
                        raw_data = r['raw'].decode('utf-8')
                        tmp['banner'] = raw_data
                    except:
                        tmp['banner'] = ""
                tmp['products'] = []
            else:
                tmp["screenshot"] = r.get("screenshot", "")
                tmp['title'] = r['title']
                tmp['status'] = r['statuscode']
                tmp['url'] = r['url']
                tmp['banner'] = r['rawheaders']
                if r['technologies'] is None:
                    r['technologies'] = []
                tmp['products'] = r['technologies']
                tmp['icon'] = r['iconcontent']
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
        logger.error(traceback.format_exc())
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/detail")
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
        domain = doc.get('host', "")
        service = doc.get("service", "")
        IP = doc.get("ip", "")
        if doc['type'] == 'other':
            URL = ""
            try:
                if doc['raw'] is not None:
                    raw_data = json.loads(doc['raw'].decode('utf-8'))
                    for k in raw_data:
                        banner += k + ":" + str(raw_data[k]).strip("\n") + "\n"
            except:
                banner = ""
        else:
            URL = doc.get("url", "")
            products = doc.get('technologies')
            if products == None:
                products = []
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


@router.post("/statistics")
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


@router.post("/statistics2")
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


@router.post("/statistics/port")
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


@router.post("/statistics/title")
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


@router.post("/statistics/type")
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


@router.post("/statistics/icon")
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


@router.post("/statistics/app")
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
