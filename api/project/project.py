import time
import traceback
import asyncio
from bson import ObjectId
from fastapi import APIRouter, Depends, BackgroundTasks

from api.project.handler import update_project, delete_asset_project_handler
from api.project.scan import scheduler_project
from api.task.handler import scheduler_scan_task, insert_task
from api.task.util import delete_asset, get_target_list
from api.users import verify_token
from core.db import get_mongo_db
from core.redis_handler import refresh_config, get_redis_pool
from core.util import *
from core.apscheduler_handler import scheduler

router = APIRouter()


@router.post("/data")
async def get_projects_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token),
                            background_tasks: BackgroundTasks = BackgroundTasks()):
    background_tasks.add_task(update_project_count)
    search_query = request_data.get("search", "")
    page_index = request_data.get("pageIndex", 1)
    page_size = request_data.get("pageSize", 10)

    query = {
        "$or": [
            {"name": {"$regex": search_query, "$options": "i"}},
            {"root_domain": {"$regex": search_query, "$options": "i"}}
        ]
    } if search_query else {}

    # 获取标签统计信息
    tag_result = await db.project.aggregate([
        {"$group": {"_id": "$tag", "count": {"$sum": 1}}},
        {"$sort": {"count": -1}}
    ]).to_list(None)

    tag_num = {tag["_id"]: tag["count"] for tag in tag_result}
    all_num = sum(tag_num.values())
    tag_num["All"] = all_num

    result_list = {}

    async def fetch_projects(tag, tag_query):
        cursor = db.project.find(tag_query, {
            "_id": 0,
            "id": {"$toString": "$_id"},
            "name": 1,
            "logo": 1,
            "AssetCount": 1,
            "tag": 1
        }).sort("AssetCount", -1).skip((page_index - 1) * page_size).limit(page_size)

        results = await cursor.to_list(length=None)
        for result in results:
            result["AssetCount"] = result.get("AssetCount", 0)
        return results

    fetch_tasks = []
    for tag in tag_num:
        if tag != "All":
            tag_query = {"$and": [query, {"tag": tag}]}
        else:
            tag_query = query

        fetch_tasks.append(fetch_projects(tag, tag_query))

    fetch_results = await asyncio.gather(*fetch_tasks)

    for tag, results in zip(tag_num, fetch_results):
        result_list[tag] = results

    return {
        "code": 200,
        "data": {
            "result": result_list,
            "tag": tag_num
        }
    }


@router.get("/all")
async def get_projects_all(db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    try:
        pipeline = [
            {
                "$group": {
                    "_id": "$tag",  # 根据 tag 字段分组
                    "children": {"$push": {"value": {"$toString": "$_id"}, "label": "$name"}}  # 将每个文档的 _id 和 name 放入 children 集合中
                }
            },
            {
                "$project": {
                    "_id": 0,
                    "label": "$_id",
                    "value": {"$literal": ""},
                    "children": 1
                }
            }
        ]
        result = await db['project'].aggregate(pipeline).to_list(None)
        return {
            "code": 200,
            "data": {
                'list': result
            }
        }
    except Exception as e:
        logger.error(str(e))
        logger.error(traceback.format_exc())
        return {"message": "error","code":500}


async def update_project_count():
    async for db in get_mongo_db():
        cursor = db.project.find({}, {"_id": 0, "id": {"$toString": "$_id"}})
        results = await cursor.to_list(length=None)

        async def update_count(id):
            query = {"project": {"$eq": id}}
            total_count = await db.asset.count_documents(query)
            if total_count != 0:
                update_document = {
                    "$set": {
                        "AssetCount": total_count
                    }
                }
                await db.project.update_one({"_id": ObjectId(id)}, update_document)

        fetch_tasks = [update_count(r['id']) for r in results]

        await asyncio.gather(*fetch_tasks)


@router.post("/content")
async def get_project_content(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token)):
    project_id = request_data.get("id")
    if not project_id:
        return {"message": "ID is missing in the request data", "code": 400}
    query = {"_id": ObjectId(project_id)}
    doc = await db.project.find_one(query)
    if not doc:
        return {"message": "Content not found for the provided ID", "code": 404}
    project_target_data = await db.ProjectTargetData.find_one({"id": project_id})
    result = {
        "name": doc.get("name", ""),
        "tag": doc.get("tag", ""),
        "target": project_target_data.get("target", ""),
        "node": doc.get("node", []),
        "logo": doc.get("logo", ""),
        "scheduledTasks": doc.get("scheduledTasks"),
        "hour": doc.get("hour"),
        "allNode": doc.get("allNode", False),
        "duplicates": doc.get("duplicates"),
        "template": doc.get("template"),
        "ignore": doc.get("ignore"),
    }
    return {"code": 200, "data": result}

def get_before_last_dash(s: str) -> str:
    index = s.rfind('-')  # 查找最后一个 '-' 的位置
    if index != -1:
        return s[:index]  # 截取从开头到最后一个 '-' 前的内容
    return s  # 如果没有 '-'，返回原字符串


@router.post("/add")
async def add_project_rule(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token),
                           background_tasks: BackgroundTasks = BackgroundTasks()):
    # Extract values from request data
    name = request_data.get("name")
    cursor = db.project.find({"name": {"$eq": name}}, {"_id": 1})
    results = await cursor.to_list(length=None)
    if len(results) != 0:
        return {"code": 400, "message": "name already exists"}
    target = request_data.get("target").strip("\n").strip("\r").strip()
    runNow = request_data.get("runNow")
    request_data["tp"] = "project"
    del request_data["runNow"]
    scheduledTasks = request_data.get("scheduledTasks", False)
    hour = request_data.get("hour", 24)
    root_domains = []
    target_list = await get_target_list(target, request_data.get("ignore", ""))
    for tg in target_list:
        if "CMP:" in tg or "ICP:" in tg or "APP:" in tg or "APP-ID:" in tg:
            root_domain = tg.replace("CMP:", "").replace("ICP:", "").replace("APP:", "").replace("APP-ID:", "")
            if "ICP:" in tg:
                root_domain = get_before_last_dash(root_domain)
        else:
            root_domain = get_root_domain(tg)
        if root_domain not in root_domains:
            root_domains.append(root_domain)
    request_data["root_domains"] = root_domains
    del request_data['target']
    # Insert the new document into the SensitiveRule collection
    result = await db.project.insert_one(request_data)
    # Check if the insertion was successful6
    if result.inserted_id:
        await db.ProjectTargetData.insert_one({"id": str(result.inserted_id), "target": target})
        if scheduledTasks:
            scheduler.add_job(scheduler_scan_task, 'interval', hours=hour, args=[str(result.inserted_id), "project"],
                              id=str(result.inserted_id), jobstore='mongo')
            next_time = scheduler.get_job(str(result.inserted_id)).next_run_time
            formatted_time = next_time.strftime("%Y-%m-%d %H:%M:%S")
            request_data["type"] = "project"
            request_data["state"] = True
            request_data["lastTime"] = ""
            request_data["nextTime"] = formatted_time
            request_data["id"] = str(result.inserted_id)
            request_data["target"] = target
            del request_data["root_domains"]
            await db.ScheduledTasks.insert_one(request_data)
        if runNow:
            time_now = get_now_time()
            del request_data["scheduledTasks"]
            request_data["target"] = target
            request_data["name"] = request_data["name"] + "-project-" + time_now
            await insert_task(request_data, db)
        background_tasks.add_task(update_project, root_domains, str(result.inserted_id), False)
        await refresh_config('all', 'project', str(result.inserted_id))
        # Project_List[name] = str(result.inserted_id)
        Project_List[str(result.inserted_id)] = name
        return {"code": 200, "message": "Project added successfully"}
    else:
        return {"code": 400, "message": "Failed to add Project"}


@router.post("/delete")
async def delete_project_rules(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token),
                               background_tasks: BackgroundTasks = BackgroundTasks()):
    try:
        pro_ids = request_data.get("ids", [])
        delA = request_data.get("delA", False)
        if delA:
            background_tasks.add_task(delete_asset, pro_ids, True)
        obj_ids = [ObjectId(poc_id) for poc_id in pro_ids]
        result = await db.project.delete_many({"_id": {"$in": obj_ids}})
        await db.ProjectTargetData.delete_many({"id": {"$in": pro_ids}})
        # Check if the deletion was successful
        if result.deleted_count > 0:
            for pro_id in pro_ids:
                job = scheduler.get_job(pro_id)
                if job:
                    scheduler.remove_job(pro_id)
                background_tasks.add_task(delete_asset_project_handler, pro_id)
                del Project_List[pro_id]
                # for project_id in Project_List:
                #     if pro_id == Project_List[project_id]:
                #         del Project_List[project_id]
                #         break
            await db.ScheduledTasks.delete_many({"id": {"$in": pro_ids}})
            return {"code": 200, "message": "Project deleted successfully"}
        else:
            return {"code": 404, "message": "Project not found"}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/update")
async def update_project_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token),
                              background_tasks: BackgroundTasks = BackgroundTasks()):
    try:
        # Get the ID from the request data
        pro_id = request_data.get("id")
        hour = request_data.get("hour")
        runNow = request_data.get("runNow")
        del request_data["runNow"]
        if not pro_id:
            return {"message": "ID is missing in the request data", "code": 400}
        query = {"id": pro_id}
        # 删除计划任务管理记录
        await db.ScheduledTasks.delete_many(query)
        # 删除计划任务
        job = scheduler.get_job(pro_id)
        if job is not None:
            scheduler.remove_job(pro_id)
        scheduledTasks = request_data.get("scheduledTasks")
        target = request_data.get("target").strip("\n").strip("\r").strip()
        # 更新目标记录
        await db.ProjectTargetData.update_one({"id": pro_id}, {"$set": {"target": target}})
        root_domains = []
        target_list = await get_target_list(target, request_data.get("ignore", ""))
        for tg in target_list:
            if "CMP:" in tg or "ICP:" in tg or "APP:" in tg or "APP-ID:" in tg:
                root_domain = tg.replace("CMP:", "").replace("ICP:", "").replace("APP:", "").replace("APP-ID:", "")
            else:
                root_domain = get_root_domain(tg)
            if root_domain not in root_domains:
                root_domains.append(root_domain)
        request_data["root_domains"] = root_domains
        request_data.pop("id")
        del request_data['target']
        update_document = {
            "$set": request_data
        }
        await db.project.update_one({"_id": ObjectId(pro_id)}, update_document)
        if scheduledTasks:
            scheduler.add_job(scheduler_scan_task, 'interval', hours=hour, args=[pro_id, "project"],
                              id=pro_id, jobstore='mongo')
            next_time = scheduler.get_job(pro_id).next_run_time
            formatted_time = next_time.strftime("%Y-%m-%d %H:%M:%S")
            request_data["state"] = True
            request_data["type"] = "project"
            request_data["lastTime"] = ""
            request_data["nextTime"] = formatted_time
            request_data["id"] = pro_id
            request_data["target"] = target
            del request_data["root_domains"]
            await db.ScheduledTasks.insert_one(request_data)
        if runNow:
            time_now = get_now_time()
            del request_data["scheduledTasks"]
            request_data["target"] = target
            request_data["name"] = request_data["name"] + "-project-" + time_now
            await insert_task(request_data, db)
        background_tasks.add_task(update_project, root_domains, pro_id, True)
        await refresh_config('all', 'project', pro_id)
        # Project_List[request_data.get("name")] = pro_id
        Project_List[pro_id] = request_data.get("name")
        return {"code": 200, "message": "successfully"}
    except Exception as e:
        logger.error(str(e))
        logger.error(traceback.format_exc())
        # Handle exceptions as needed
        return {"message": "error", "code": 500}

