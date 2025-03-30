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

from core.db import get_mongo_db
from core.util import *
from pymongo import DESCENDING
from loguru import logger
router = APIRouter()


@router.post("/data")
async def asset_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        page_index = request_data.get("pageIndex", 1)
        page_size = request_data.get("pageSize", 10)
        query = await get_search_query("asset", request_data)
        if query == "":
            return {"message": "Search condition parsing error", "code": 500}
        # total_count = await db['asset'].count_documents(query)
        cursor: AsyncIOMotorCursor = db['asset'].find(query, {"_id": 0,
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
                                          "metadata": 1,
                                          "time": 1,
                                          "iconcontent": 1,
                                          "tags": 1,
                                          "screenshot": 1,
                                          }).skip((page_index - 1) * page_size).limit(page_size).sort(
            [("time", DESCENDING)])
        result_list = []
        async for r in cursor:
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
                    if r['metadata'] is not None:
                        raw_data = json.loads(r['metadata'].decode('utf-8'))
                        for k in raw_data:
                            tmp['banner'] += k + ":" + str(raw_data[k]).strip("\n") + "\n"
                except:
                    try:
                        raw_data = r['metadata'].decode('utf-8')
                        tmp['banner'] = raw_data
                    except:
                        tmp['banner'] = ""
                tmp['products'] = []
            else:
                tmp["screenshot"] = r.get("screenshot", "")
                # if tmp["screenshot"] != "":
                #     tmp["screenshot"] = "true"
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
            }
        }
    except Exception as e:
        logger.error(str(e))
        logger.error(traceback.format_exc())
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/data/card")
async def asset_card_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        page_index = request_data.get("pageIndex", 1)
        page_size = request_data.get("pageSize", 10)
        query = await get_search_query("asset", request_data)
        if query == "":
            return {"message": "Search condition parsing error", "code": 500}
        total_count = await db['asset'].count_documents(query)
        cursor: AsyncIOMotorCursor = db['asset'].find(query, {"_id": 0,
                                          "host": 1,
                                          "url": 1,
                                          "port": 1,
                                          "service": 1,
                                          "type": 1,
                                          "title": 1,
                                          "statuscode": 1,
                                          "screenshot": 1,
                                          }).skip((page_index - 1) * page_size).limit(page_size).sort(
            [("time", DESCENDING)])
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
        logger.error(traceback.format_exc())
        # Handle exceptions as needed
        return {"message": "error", "code": 500}

@router.post("/screenshot")
async def asset_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    id = request_data.get("id", "")
    if id == "":
        return {"message": "not found", "code": 404}
    query = {"_id": ObjectId(id)}
    doc = await db.asset.find_one(query, {"screenshot": 1})
    if doc is None:
        return {"message": "not found", "code": 404}
    screenshot = doc.get('screenshot', "")
    return {
            "code": 200,
            "data": {
                'screenshot': screenshot
            }
        }


# @router.post("/detail")
# async def asset_detail(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
#     try:
#         # Get the ID from the request data
#         asset_id = request_data.get("id")
#
#         # Check if ID is provided
#         if not asset_id:
#             return {"message": "ID is missing in the request data", "code": 400}
#
#         # Query the database for content based on ID
#         query = {"_id": ObjectId(asset_id)}
#         doc = await db.asset.find_one(query)
#
#         if not doc:
#             return {"message": "Content not found for the provided ID", "code": 404}
#         products = []
#         tlsdata = ""
#         hashes = ""
#         banner = ""
#         domain = doc.get('host', "")
#         service = doc.get("service", "")
#         IP = doc.get("ip", "")
#         if doc['type'] == 'other':
#             URL = ""
#             try:
#                 if doc['raw'] is not None:
#                     raw_data = json.loads(doc['raw'].decode('utf-8'))
#                     for k in raw_data:
#                         banner += k + ":" + str(raw_data[k]).strip("\n") + "\n"
#             except:
#                 banner = ""
#         else:
#             URL = doc.get("url", "")
#             products = doc.get('technologies')
#             if products == None:
#                 products = []
#             if doc['tlsdata'] is not None:
#                 for h in doc['tlsdata']:
#                     tlsdata += h + ": " + str(doc['tlsdata'][h]) + '\n'
#             if doc['hashes'] is not None:
#                 for h in doc['hashes']:
#                     hashes += h + ": " + str(doc['hashes'][h]) + '\n'
#             banner = doc.get('rawheaders', "")
#         project_name = ""
#         if doc.get("project", "") != "":
#             query = {"_id": ObjectId(doc.get("project", ""))}
#             project_data = await db.project.find_one(query)
#             project_name = project_data.get("name", "")
#         data = {
#             "host": domain,
#             "IP": IP,
#             "URL": URL,
#             "port": doc.get("port", ""),
#             "service": service,
#             "title": doc.get("title", ""),
#             "status": doc.get("statuscode", ""),
#             "FaviconHash": doc.get("faviconmmh3", ""),
#             "jarm": doc.get("jarm", ""),
#             "time": doc.get("timestamp", ""),
#             "products": products,
#             "TLSData": tlsdata,
#             "hash": hashes,
#             "banner": banner,
#             "ResponseBody": doc.get("responsebody", ""),
#             "project": project_name
#         }
#         return {"code": 200, "data": data}
#     except Exception as e:
#         logger.error(str(e))
#         # Handle exceptions as needed
#         return {"message": "error", "code": 500}


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
        doc["id"] = str(doc["_id"])
        del doc["_id"]
        return {"code": 200, "data": {"json": doc}}
    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/changelog")
async def asset_detail(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # Get the ID from the request data
        asset_id = request_data.get("id")

        # Check if ID is provided
        if not asset_id:
            return {"message": "ID is missing in the request data", "code": 400}

        # Query the database for content based on ID
        query = {"assetid": asset_id}
        cursor: AsyncIOMotorCursor = db.AssetChangeLog.find(query).sort([("time", DESCENDING)])
        results = await cursor.to_list(length=None)
        result_list = []
        for result in results:
            del result["_id"]
            result_list.append(result)
        return {"code": 200, "data": result_list}
    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}