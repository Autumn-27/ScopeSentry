# -------------------------------------
# @file      : dir.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/9/8 15:10
# -------------------------------------------
import re

from bson import ObjectId
from fastapi import APIRouter, Depends, File, UploadFile, Query, Form
from starlette.responses import StreamingResponse

from api.users import verify_token
from motor.motor_asyncio import AsyncIOMotorCursor, AsyncIOMotorGridFSBucket
from core.db import get_mongo_db
from core.redis_handler import refresh_config
from loguru import logger
router = APIRouter()


@router.post("/list")
async def list(db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    fs = AsyncIOMotorGridFSBucket(db)
    cursor = fs.find({"filename": {"$regex": f"^dir_"}})
    files_data = await cursor.to_list(length=None)
    result_list = []
    for file in files_data:
        result_list.append({
            "id": str(file["_id"]),
            "name": file["filename"].replace("dir_", ""),
            "size": "{:.2f}".format(file["length"] / (1024 * 1024)),
        })
    return {
            "code": 200,
            "data": {
                'list': result_list,
            }
        }


@router.post("/save")
async def save_dir_data(file: UploadFile = File(...), id: str = Query(...), db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        content = await file.read()
        fs = AsyncIOMotorGridFSBucket(db)

        # 查找文件，根据传入的 id
        old_file = await fs.find({'_id': ObjectId(id)}).to_list(1)

        if old_file:
            # 如果找到文件，删除旧文件
            await fs.delete(ObjectId(id))
            print(f"File with id {id} deleted.")

            # 将新文件内容存储到 GridFS 中，并保存新的文件 ID
            await fs.upload_from_stream(
                old_file[0]['filename'],  # 使用原文件名存储
                content  # 文件内容
            )
            return {"code": 200, "message": "upload successful"}
        else:
            return {"message": "not found", "code": 500}
    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/add")
async def add_dir_data(file: UploadFile = File(...), name: str = Form(...), db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        if not re.match(r'^[a-zA-Z0-9_-]+$', name):
            return {"message": "Name must contain only letters", "code": 500}
        content = await file.read()
        fs = AsyncIOMotorGridFSBucket(db)
        await fs.upload_from_stream(
            f"dir_{name}",  # 使用原文件名存储
            content  # 文件内容
        )
        return {"code": 200, "message": "upload successful"}
    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.get("/download")
async def get_dir_data(id: str = Query(...), db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        fs = AsyncIOMotorGridFSBucket(db)

        # 查找文件
        file_doc = await fs.find({'_id': ObjectId(id)}).to_list(1)

        if not file_doc:
            return {'code': 404, 'message': 'file is not found'}

        file_id = file_doc[0]['_id']
        grid_out = await fs.open_download_stream(file_id)
        filename = file_doc[0]['filename'].replace("dir_", "")
        # 返回文件流
        return StreamingResponse(grid_out, media_type="application/octet-stream",
                                 headers={"Content-Disposition": f"attachment; filename={filename}"})
    except Exception as e:
        logger.error(str(e))
        return {"message": "error", "code": 500}


@router.post("/delete")
async def delete_dir_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    dir_dict_ids = request_data.get("ids", [])
    fs = AsyncIOMotorGridFSBucket(db)
    for id in dir_dict_ids:
        try:
            file_id = ObjectId(id)
            # 删除文件
            await fs.delete(file_id)
            print(f"文件 {id} 已删除")
        except Exception as e:
            print(f"删除文件 {id} 时发生错误: {e}")
            logger.error(str(e))
            return {"message": "error", "code": 500}
    return {"code": 200, "message": "delete file successful"}
