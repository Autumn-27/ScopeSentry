# -------------------------------------
# @file      : manage.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/10/17 21:09
# -------------------------------------------
import re
from typing import List

from bson import ObjectId
from fastapi import APIRouter, Depends, File, UploadFile, Query, Form
from pydantic import BaseModel
from starlette.responses import StreamingResponse

from api.users import verify_token
from motor.motor_asyncio import AsyncIOMotorCursor, AsyncIOMotorGridFSBucket
from core.db import get_mongo_db
from core.redis_handler import refresh_config
from loguru import logger
router = APIRouter()


@router.get("/list")
async def get_all_files(db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        cursor: AsyncIOMotorCursor = db["dictionary"].find({})
        result = await cursor.to_list(length=None)
        response_data = []
        for doc in result:
            data = {
                "name": doc["name"],
                "category": doc["category"],
                "size": doc["size"],
                "id": str(doc["_id"])
            }
            response_data.append(data)

        return {
            "code": 200,
            "data": {
                'list': response_data,
            }
        }
    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/create")
async def add_dictionary_data(file: UploadFile = File(...), name: str = Form(...), category: str = Form(...), db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        result = await db["dictionary"].find_one({"name": name, "category": category})
        if result:
            return {
                "code": 500,
                "message": "duplication file name"
            }
        content = await file.read()
        size = len(content) / (1024 * 1024)
        result = await db["dictionary"].insert_one({"name": name, "category": category, "size": "{:.2f}".format(size)})
        if result.inserted_id:
            fs = AsyncIOMotorGridFSBucket(db)
            await fs.upload_from_stream(
                str(result.inserted_id),  # 使用id作为文件名存储
                content  # 文件内容
            )
            await refresh_config('all', f'dictionary', f"add:{str(result.inserted_id)}")
            return {"message": "file added successfully", "code": 200}
        else:
            return {"message": "Failed to add file", "code": 500}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.get("/download")
async def get_dictionary_file(id: str = Query(...), db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        fs = AsyncIOMotorGridFSBucket(db)

        # 查找文件
        file_doc = await fs.find({'filename': id}).to_list(1)

        if not file_doc:
            return {'code': 404, 'message': 'file is not found'}

        file_id = file_doc[0]['_id']
        grid_out = await fs.open_download_stream(file_id)
        filename = file_doc[0]['filename']
        # 返回文件流
        return StreamingResponse(grid_out, media_type="application/octet-stream",
                                 headers={"Content-Disposition": f"attachment; filename={filename}"})
    except Exception as e:
        logger.error(str(e))
        return {"message": "error", "code": 500}


@router.post("/delete")
async def delete_dictionary(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        dir_dict_ids = request_data.get("ids", [])
        obj_ids = []
        for id in dir_dict_ids:
            obj_ids.append(ObjectId(id))
            await refresh_config('all', 'dictionary', f"delete:{id}")
        await db["dictionary"].delete_many({"_id": {"$in": obj_ids}})

        fs = AsyncIOMotorGridFSBucket(db)
        for id in dir_dict_ids:
            try:
                # 查找文件
                file_doc = await fs.find({'filename': id}).to_list(1)

                if not file_doc:
                    logger.error(f"file {id} is not found")
                    continue
                file_id = ObjectId(file_doc[0]['_id'])

                # 删除文件
                await fs.delete(file_id)
            except Exception as e:
                print(f"删除文件 {id} 时发生错误: {e}")
                logger.error(str(e))
                return {"message": "error", "code": 500}
        return {"code": 200, "message": "delete file successful"}
    except Exception as e:
        logger.error(str(e))
        return {"message": "error", "code": 500}


@router.post("/save")
async def save_dir_data(file: UploadFile = File(...), id: str = Query(...), db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        content = await file.read()
        fs = AsyncIOMotorGridFSBucket(db)

        # 查找文件，根据传入的 id
        old_file = await fs.find({'filename': id}).to_list(1)

        if old_file:
            # 如果找到文件，删除旧文件
            await fs.delete(ObjectId(old_file[0]['_id']))
            logger.info(f"File with id {id} deleted.")

            # 将新文件内容存储到 GridFS 中，并保存新的文件 ID
            await fs.upload_from_stream(
                id,  # 使用原文件名存储
                content  # 文件内容
            )
            await db["dictionary"].update_one({"_id": ObjectId(id)}, {"$set": {"size": "{:.2f}".format(len(content) / (1024 * 1024))}})
            await refresh_config('all', f'dictionary', f"add:{id}")
            return {"code": 200, "message": "upload successful"}
        else:
            logger.error(f"File not found {id}")
            return {"message": "not found", "code": 500}
    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}