# -------------------------------------
# @file      : common.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/10/20 21:20
# -------------------------------------------
import traceback

from bson import ObjectId
from fastapi import APIRouter, Depends
from api.users import verify_token
from core.db import get_mongo_db
from core.util import *
from loguru import logger

router = APIRouter()


@router.post("/delete")
async def delete_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        data_ids = request_data.get("ids", [])
        index = request_data.get("index", "")
        key = ["asset", "DirScanResult", "SensitiveResult", "SubdoaminTakerResult", "UrlScan", "crawler", "subdomain", "vulnerability", "PageMonitoring", "app", "RootDomain", "mp"]
        if index not in key:
            return {"code": 404, "message": "Data not found"}
        obj_ids = []
        for data_id in data_ids:
            if data_id is not None and data_id != "" and "http://" not in data_id and "https://" not in data_id:
                if len(data_id) > 6:
                    try:
                        obj_ids.append(ObjectId(data_id))
                    except:
                        continue
        result = await db[index].delete_many({"_id": {"$in": obj_ids}})

        if result.deleted_count > 0:
            return {"code": 200, "message": "Data deleted successfully"}
        else:
            return {"code": 404, "message": "Data not found"}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/add/tag")
async def add_tag(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        tp = request_data.get("tp", "")
        id = request_data.get("id", "")
        tag = request_data.get("tag", "")
        key = ["asset", "DirScanResult", "SensitiveResult", "SubdoaminTakerResult", "UrlScan", "crawler", "subdomain",
               "vulnerability", "PageMonitoring", "RootDomain", "app", "mp"]
        if tp not in key or id == "" or tag == "":
            return {"code": 404, "message": "Data not found"}
        query = {"_id": ObjectId(id)}
        doc = await db[tp].find_one(query, {"tags": 1})

        if not doc:
            return {"message": "Content not found for the provided ID", "code": 404}
        if "tags" not in doc or doc["tags"] is None:
            doc["tags"] = [tag]
        else:
            doc["tags"].append(tag)
        await db[tp].update_one(query, {"$set": {"tags": doc["tags"]}})
        return {"message": "success", "code": 200}
    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/delete/tag")
async def delete_tag(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        tp = request_data.get("tp", "")
        id = request_data.get("id", "")
        tag = request_data.get("tag", "")
        key = ["asset", "DirScanResult", "SensitiveResult", "SubdoaminTakerResult", "UrlScan", "crawler", "subdomain",
               "vulnerability", "PageMonitoring"]
        if tp not in key or id == "" or tag == "":
            return {"code": 404, "message": "Data not found"}
        query = {"_id": ObjectId(id)}
        doc = await db[tp].find_one(query, {"tags": 1})
        if "tags" not in doc or doc["tags"] is None:
            return {"message": "success", "code": 200}
        doc["tags"].remove(tag)
        await db[tp].update_one(query, {"$set": {"tags": doc["tags"]}})
        return {"message": "success", "code": 200}
    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/update/status")
async def delete_tag(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        tp = request_data.get("tp", "")
        id = request_data.get("id", "")
        status = request_data.get("status")
        key = ["SensitiveResult", "vulnerability"]
        if tp not in key or id == "":
            return {"code": 404, "message": "Data not found"}
        query = {"_id": ObjectId(id)}
        await db[tp].update_one(query, {"$set": {"status": status}})
        return {"message": "success", "code": 200}
    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/total")
async def total_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    index = request_data.get("index", "")
    key = ["asset", "DirScanResult", "SensitiveResult", "SubdoaminTakerResult", "UrlScan", "crawler", "subdomain",
           "vulnerability", "PageMonitoring"]
    if index not in key:
        return {"code": 404, "message": "Data not found"}
    query = await get_search_query(index, request_data)
    if query == "":
        return {"message": "Search condition parsing error", "code": 500}
    if len(query) == 0:
        total_count = await db['asset'].count_documents({"_id": {"$exists": True}})
    else:
        total_count = await db['asset'].count_documents(query)
    return {
        "code": 200,
        "data": {
            'total': total_count
        }
    }