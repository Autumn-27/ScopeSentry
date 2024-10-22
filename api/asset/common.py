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