# -*- coding:utf-8 -*-　　
# @name: sensitive
# @auth: rainy-autumn@outlook.com
# @version:
from datetime import datetime

from bson import ObjectId, SON
from fastapi import APIRouter, Depends
from pymongo import DESCENDING

from api.users import verify_token
from motor.motor_asyncio import AsyncIOMotorCursor

from core.db import get_mongo_db
from core.redis_handler import refresh_config
from loguru import logger

from core.util import search_to_mongodb, get_search_query

router = APIRouter()


@router.post("/data")
async def get_sensitive_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        search_query = request_data.get("search", "")
        page_index = request_data.get("pageIndex", 1)
        page_size = request_data.get("pageSize", 10)
        # MongoDB collection for SensitiveRule
        # Fuzzy search based on the name field
        query = {"name": {"$regex": search_query, "$options": "i"}}
        # Get the total count of documents matching the search criteria
        total_count = await db.SensitiveRule.count_documents(query)

        # Perform pagination query
        cursor: AsyncIOMotorCursor = db.SensitiveRule.find(query).skip((page_index - 1) * page_size).limit(page_size).sort([("_id", DESCENDING)])
        result = await cursor.to_list(length=None)
        if len(result) == 0:
            return {
            "code": 200,
            "data": {
                'list': [],
                'total': 0
            }
        }
        # Process the result as needed
        response_data = [{"id": str(doc["_id"]),"name": doc["name"], "regular": doc["regular"], "state": doc["state"], "color": doc["color"]} for doc in result]
        return {
            "code": 200,
            "data": {
                'list': response_data,
                'total': total_count
            }
        }

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error","code":500}


@router.post("/update")
async def upgrade_sensitive_rule(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # Extract values from request data
        rule_id = request_data.get("id")
        name = request_data.get("name")
        regular = request_data.get("regular")
        color = request_data.get("color")
        state = request_data.get("state")

        # Update query based on rule_id
        update_query = {"_id": ObjectId(rule_id)}

        # Values to be updated
        update_values = {"$set": {"name": name, "regular": regular, "color": color, "state": state}}

        # Perform the update
        result = await db.SensitiveRule.update_one(update_query, update_values)
        if result:
            await refresh_config('all', 'sensitive')
            return {"code": 200, "message": "SensitiveRule updated successfully"}
        else:
            return {"code": 404, "message": "SensitiveRule not found"}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}

@router.post("/add")
async def add_sensitive_rule(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # Extract values from request data
        name = request_data.get("name")
        regular = request_data.get("regular",'')
        color = request_data.get("color")
        state = request_data.get("state")
        if regular == '':
            return {"code": 500, "message": "regular is null"}
        # Create a new SensitiveRule document
        new_rule = {
            "name": name,
            "regular": regular,
            "color": color,
            "state": state
        }

        # Insert the new document into the SensitiveRule collection
        result = await db.SensitiveRule.insert_one(new_rule)

        # Check if the insertion was successful
        if result.inserted_id:
            await refresh_config('all', 'sensitive')
            return {"code": 200, "message": "SensitiveRule added successfully"}
        else:
            return {"code": 400, "message": "Failed to add SensitiveRule"}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/update/state")
async def update_state_sensitive_rule(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        rule_ids = request_data.get("ids", [])
        state = request_data.get("state")
        if state is None:
            return {"code": 500, "message": "state not found"}
        obj_ids = []
        for rule_id in rule_ids:
            if rule_id != None and rule_id != "":
                obj_ids.append(ObjectId(rule_id))
        result = await db.SensitiveRule.update_many({'_id': {'$in': obj_ids}}, {'$set': {'state': state}})
        # Check if the deletion was successful
        if result.modified_count > 0:
            await refresh_config('all', 'sensitive')
            return {"code": 200, "message": "SensitiveRules update successfully"}
        else:
            return {"code": 404, "message": "SensitiveRules not found"}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/delete")
async def delete_sensitive_rules(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # Extract the list of IDs from the request_data dictionary
        rule_ids = request_data.get("ids", [])

        # Convert the provided rule_ids to ObjectId
        obj_ids = []
        for rule_id in rule_ids:
            if rule_id is not None and rule_id != "":
                obj_ids.append(ObjectId(rule_id))
        # Delete the SensitiveRule documents based on the provided IDs
        result = await db.SensitiveRule.delete_many({"_id": {"$in": obj_ids}})

        # Check if the deletion was successful
        if result.deleted_count > 0:
            await refresh_config('all', 'sensitive')
            return {"code": 200, "message": "SensitiveRules deleted successfully"}
        else:
            return {"code": 404, "message": "SensitiveRules not found"}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/result/data")
async def get_sensitive_result_rules(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        search_query = request_data.get("search", "")
        page_index = request_data.get("pageIndex", 1)
        page_size = request_data.get("pageSize", 10)
        keyword = {
            'url': 'url',
            'sname': 'sid',
            "body": "body",
            "info": "match",
            'project': 'project',
            'md5': 'md5'
        }
        query = await search_to_mongodb(search_query, keyword)
        if query == "" or query is None:
            return {"message": "Search condition parsing error", "code": 500}
        query = query[0]
        filter = request_data.get("filter", {})
        if filter:
            query["$and"] = []
            for f in filter:
                tmp_or = []
                for v in filter[f]:
                    tmp_or.append({f: v})
                query["$and"].append({"$or": tmp_or})
        total_count = await db['SensitiveResult'].count_documents(query)
        cursor: AsyncIOMotorCursor = ((db['SensitiveResult'].find(query, {"_id": 0,
                                                                "id": {"$toString": "$_id"},
                                                                "url": 1,
                                                                "sid": 1,
                                                                "match": 1,
                                                                "time": 1,
                                                                "color": 1
                                                                })
                                       .skip((page_index - 1) * page_size)
                                       .limit(page_size))
                                      .sort([("time", DESCENDING)]))
        result = await cursor.to_list(length=None)
        result_list = []
        for r in result:
            tmp = {
                'id': r['id'],
                'url': r['url'],
                'name': r['sid'],
                'color': r['color'],
                'match': r['match'],
                'time': r['time']
            }
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
        return {"message": "error","code":500}


@router.post("/result/data2")
async def get_sensitive_result_data2(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        page_index = request_data.get("pageIndex", 1)
        page_size = request_data.get("pageSize", 10)
        query = await get_search_query("SensitiveResult", request_data)
        if query == "":
            return {"message": "Search condition parsing error", "code": 500}
        total_count = await db['SensitiveResult'].count_documents(query)
        pipeline = [
            {
                "$match": query  # 增加搜索条件
            },
            {
                "$project": {
                    "_id": 1,
                    "url": 1,
                    "time": 1,
                    "sid": 1,
                    "match": 1,
                    "color": 1,
                    "md5": 1,
                    "tags": 1,
                    "status": 1,
                }
            },
            {
                "$sort": {"_id": DESCENDING}  # 按时间降序排序
            },
            {
                "$group": {
                    "_id": "$url",
                    "time": {"$first": "$time"},  # 记录相同url下最早插入数据的时间
                    "url": {"$first": "$url"},
                    "body_id": {"$last": {"$toString": "$md5"}},  # 记录相同url下最早插入数据的_id
                    "children": {
                        "$push": {
                            "id": {"$toString": "$_id"},
                            "name": "$sid",
                            "color": "$color",
                            "match": "$match",
                            "time": "$time",
                            "tags": "$tags",
                            "status": "$status"
                        }
                    }
                }
            },
            {
                "$sort": {"time": DESCENDING}  # 按每组的最新时间降序排序
            },
            {
                "$skip": (page_index - 1) * page_size  # 跳过前面的URL，用于分页
            },
            {
                "$limit": page_size  # 获取当前页的URL
            },
            {
                "$project": {
                    "id": "$_id",  # 这里将最终的 _id 输出为 id
                    "url": 1,
                    "time": 1,
                    "body_id": 1,
                    "children": 1,
                    "tags": 1,
                    "status": 1
                }
            }
        ]
        # 执行聚合查询
        result = await db['SensitiveResult'].aggregate(pipeline).to_list(None)
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
        return {"message": "error","code":500}


@router.post("/result/names")
async def get_sensitive_result_names(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    query = await get_search_query("SensitiveResult", request_data)
    if query == "":
        return {"message": "Search condition parsing error", "code": 500}

    pipeline = [
        {
            "$match": query
        },
        {"$group": {"_id": "$sid", "count": {"$sum": 1}, "color": {"$first": "$color"}}},
        {"$sort": SON([("count", -1)])},
        {
            "$project": {
                "name": "$_id",
                "count": 1,
                "_id": 0,
                "color": 1
            }
        }
    ]
    result = await db['SensitiveResult'].aggregate(pipeline).to_list(None)
    return {
            "code": 200,
            "data": {
                'list': result
            }
        }


@router.post("/result/info")
async def get_sensitive_result_info(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    sid = request_data.get("sid", "")
    if sid == "":
        return {"message": "sid is null", "code": 500}
    query = await get_search_query("SensitiveResult", request_data)
    if query == "":
        return {"message": "Search condition parsing error", "code": 500}
    query["sid"] = sid
    pipeline = [
        {
            "$match": query
        },
        {
            "$project": {
                "match": 1  # 只保留 match 字段
            }
        },
        {
            "$unwind": "$match"  # 将 match 数组展开为多个文档
        },
        {
            "$group": {
                "_id": None,
                "unique_matches": {"$addToSet": "$match"}  # 使用 $addToSet 去重
            }
        }
    ]
    result = await db['SensitiveResult'].aggregate(pipeline).to_list(None)
    return {
        "code": 200,
        "data": {
            'list': result[0]["unique_matches"]
        }
    }


@router.post("/result/body")
async def get_sensitive_result_body_rules(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        # Get the ID from the request data
        sensitive_result_id = request_data.get("id")

        # Check if ID is provided
        if not sensitive_result_id:
            return {"message": "ID is missing in the request data", "code": 400}

        # Query the database for content based on ID
        query = {"md5": sensitive_result_id}
        doc = await db.SensitiveBody.find_one(query)

        if not doc:
            return {"message": "Content not found for the provided ID", "code": 404}

        # Extract the content
        content = doc.get("body", "")

        return {"code": 200, "data": {"body": content}}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}