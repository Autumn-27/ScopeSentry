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
from core.util import *
from loguru import logger

router = APIRouter()


@router.post("/delete")
async def delete_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        data_ids = request_data.get("ids", [])
        index = request_data.get("index", "")
        key = ["asset", "DirScanResult", "SensitiveResult", "SubdoaminTakerResult", "UrlScan", "crawler", "subdomain", "vulnerability", "PageMonitoring"]
        if index not in key:
            return {"code": 404, "message": "Data not found"}
        obj_ids = []
        for data_id in data_ids:
            if data_id is not None and data_id != "":
                if len(data_id) > 6:
                    obj_ids.append(ObjectId(data_id))
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
        tp = request_data.get("type", "")
        id = request_data.get("id", "")
        tag = request_data.get("tag", "")
        key = ["asset", "DirScanResult", "SensitiveResult", "SubdoaminTakerResult", "UrlScan", "crawler", "subdomain",
               "vulnerability", "PageMonitoring"]
        if tp not in key or id == "" or tag == "":
            return {"code": 404, "message": "Data not found"}
        query = {"_id": ObjectId(id)}
        doc = await db[tp].find_one(query, {"tags": 1})

        if not doc:
            return {"message": "Content not found for the provided ID", "code": 404}
        if doc["tags"] is None:
            doc["tags"] = [tag]
        else:
            doc["tags"] = doc["tags"].append(tag)
        await db[tp].update_one(query, {"$set": {"tags": doc["tags"]}})
    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}

