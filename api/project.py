import time
import traceback
import asyncio
from bson import ObjectId
from fastapi import APIRouter, Depends, BackgroundTasks

from api.task import create_scan_task, delete_asset
from api.users import verify_token
from motor.motor_asyncio import AsyncIOMotorCursor

from core.config import Project_List
from core.db import get_mongo_db
from core.redis_handler import refresh_config, get_redis_pool
from loguru import logger
from core.util import *
from core.apscheduler_handler import scheduler

router = APIRouter()


@router.post("/project/data")
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


@router.get("/project/all")
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


@router.post("/project/content")
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
        "subdomainScan": doc.get("subdomainScan", False),
        "subdomainConfig": doc.get("subdomainConfig", []),
        "urlScan": doc.get("urlScan", False),
        "sensitiveInfoScan": doc.get("sensitiveInfoScan", False),
        "pageMonitoring": doc.get("pageMonitoring", ""),
        "crawlerScan": doc.get("crawlerScan", False),
        "vulScan": doc.get("vulScan", False),
        "vulList": doc.get("vulList", []),
        "portScan": doc.get("portScan"),
        "ports": doc.get("ports"),
        "waybackurl": doc.get("waybackurl"),
        "dirScan": doc.get("dirScan"),
        "scheduledTasks": doc.get("scheduledTasks"),
        "hour": doc.get("hour"),
        "allNode": doc.get("allNode", False),
        "duplicates": doc.get("duplicates")
    }
    return {"code": 200, "data": result}


@router.post("/project/add")
async def add_project_rule(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token),
                           background_tasks: BackgroundTasks = BackgroundTasks()):
    try:
        # Extract values from request data
        name = request_data.get("name")
        target = request_data.get("target")
        runNow = request_data.get("runNow")
        del request_data["runNow"]
        scheduledTasks = request_data.get("scheduledTasks", False)
        hour = request_data.get("hour", 1)
        t_list = []
        tmsg = ''
        root_domains = []
        for t in target.split('\n'):
            if t not in t_list:
                targetTmp = t.replace('http://', "").replace('https://', "").strip()
                if targetTmp != "":
                    root_domain = get_root_domain(targetTmp)
                    if root_domain not in root_domains:
                        root_domains.append(root_domain)
                    tmsg += targetTmp + '\n'
                # Create a new SensitiveRule document
        tmsg = tmsg.strip().strip("\n")
        request_data["root_domains"] = root_domains
        del request_data['target']
        if "All Poc" in request_data['vulList']:
            request_data['vulList'] = ["All Poc"]
        cursor = db.project.find({"name": {"$eq": name}}, {"_id": 1})
        results = await cursor.to_list(length=None)
        if len(results) != 0:
            return {"code": 400, "message": "name already exists"}
        # Insert the new document into the SensitiveRule collection
        result = await db.project.insert_one(request_data)
        # Check if the insertion was successful
        if result.inserted_id:
            await db.ProjectTargetData.insert_one({"id": str(result.inserted_id), "target": tmsg})
            if scheduledTasks:
                scheduler.add_job(scheduler_project, 'interval', hours=hour, args=[str(result.inserted_id)],
                                  id=str(result.inserted_id), jobstore='mongo')
                next_time = scheduler.get_job(str(result.inserted_id)).next_run_time
                formatted_time = next_time.strftime("%Y-%m-%d %H:%M:%S")
                db.ScheduledTasks.insert_one(
                    {"id": str(result.inserted_id), "name": name, 'hour': hour, 'type': 'Project', 'state': True,
                     'lastTime': get_now_time(), 'nextTime': formatted_time, 'runner_id': str(result.inserted_id)})
                if runNow:
                    await scheduler_project(str(result.inserted_id))
            background_tasks.add_task(update_project, root_domains, str(result.inserted_id))
            await refresh_config('all', 'project')
            Project_List[name] = str(result.inserted_id)
            return {"code": 200, "message": "Project added successfully"}
        else:
            return {"code": 400, "message": "Failed to add Project"}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/project/delete")
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
                for project_id in Project_List:
                    if pro_id == Project_List[project_id]:
                        del Project_List[project_id]
                        break
            await db.ScheduledTasks.delete_many({"id": {"$in": pro_ids}})
            return {"code": 200, "message": "Project deleted successfully"}
        else:
            return {"code": 404, "message": "Project not found"}

    except Exception as e:
        logger.error(str(e))
        # Handle exceptions as needed
        return {"message": "error", "code": 500}


@router.post("/project/update")
async def update_project_data(request_data: dict, db=Depends(get_mongo_db), _: dict = Depends(verify_token),
                              background_tasks: BackgroundTasks = BackgroundTasks()):
    try:
        # Get the ID from the request data
        pro_id = request_data.get("id")
        hour = request_data.get("hour")
        runNow = request_data.get("runNow")
        del request_data["runNow"]
        # Check if ID is provided
        if not pro_id:
            return {"message": "ID is missing in the request data", "code": 400}
        query = {"id": pro_id}
        doc = await db.ScheduledTasks.find_one(query)

        newScheduledTasks = request_data.get("scheduledTasks")
        if doc:
            oldScheduledTasks = doc["state"]
            old_hour = doc["hour"]
            if oldScheduledTasks != newScheduledTasks:
                if newScheduledTasks:
                    scheduler.add_job(scheduler_project, 'interval', hours=hour, args=[pro_id],
                                      id=str(pro_id), jobstore='mongo')
                    await db.ScheduledTasks.update_one({"id": pro_id}, {"$set": {'state': True}})
                else:
                    job = scheduler.get_job(pro_id)
                    if job is not None:
                        scheduler.remove_job(pro_id)
                    await db.ScheduledTasks.update_one({"id": pro_id}, {"$set": {'state': False}})
            else:
                if newScheduledTasks:
                    if hour != old_hour:
                        job = scheduler.get_job(pro_id)
                        if job is not None:
                            scheduler.remove_job(pro_id)
                        scheduler.add_job(scheduler_project, 'interval', hours=hour, args=[pro_id],
                                          id=str(pro_id), jobstore='mongo')
        else:
            if newScheduledTasks:
                scheduler.add_job(scheduler_project, 'interval', hours=hour, args=[str(pro_id)],
                                  id=str(pro_id), jobstore='mongo')
                next_time = scheduler.get_job(str(pro_id)).next_run_time
                formatted_time = next_time.strftime("%Y-%m-%d %H:%M:%S")
                db.ScheduledTasks.insert_one(
                    {"id": str(pro_id), "name": request_data['name'], 'hour': hour, 'type': 'Project', 'state': True,
                     'lastTime': get_now_time(), 'nextTime': formatted_time})

        result = await db.project.find_one({"_id": ObjectId(pro_id)})
        project_target_data = await db.ProjectTargetData.find_one({"id": pro_id})
        old_targets = project_target_data['target']
        old_name = result['name']
        new_name = request_data['name']
        new_targets = request_data['target']
        if old_targets != new_targets.strip().strip('\n'):
            new_root_domain = []
            update_document = {
                "$set": {
                    "target": new_targets.strip().strip('\n')
                }
            }
            for n_t in new_targets.strip().strip('\n').split('\n'):
                t_root_domain = get_root_domain(n_t)
                if t_root_domain not in new_root_domain:
                    new_root_domain.append(t_root_domain)
            request_data["root_domains"] = new_root_domain
            await db.ProjectTargetData.update_one({"id": pro_id}, update_document)
            background_tasks.add_task(update_project, new_root_domain, pro_id, True)
        if old_name != new_name:
            del Project_List[old_name]
            Project_List[new_name] = pro_id
            await db.ScheduledTasks.update_one({"id": pro_id}, {"$set": {'name': new_name}})
        request_data.pop("id")
        del request_data['target']
        update_document = {
            "$set": request_data
        }
        result = await db.project.update_one({"_id": ObjectId(pro_id)}, update_document)
        # Check if the update was successful
        if result:
            if runNow:
                await scheduler_project(str(pro_id))
            return {"message": "Task updated successfully", "code": 200}
        else:
            return {"message": "Failed to update data", "code": 404}

    except Exception as e:
        logger.error(str(e))
        logger.error(traceback.format_exc())
        # Handle exceptions as needed
        return {"message": "error", "code": 500}

async def update_project(root_domain, project_id, change=False):
    asset_collection_list = {
                        'asset': ["url", "host", "ip"],
                        'subdomain': ["host", "ip"],
                        'DirScanResult': ["url"],
                        'vulnerability': ["url"],
                        'SubdoaminTakerResult': ["input"],
                        'PageMonitoring': ["url"],
                        'SensitiveResult': ["url"],
                        'UrlScan': ["input"],
                        'crawler': ["url"]}
    async for db in get_mongo_db():
        for a in asset_collection_list:
            if change:
                await asset_update_project(root_domain, asset_collection_list[a], a, db, project_id)
            else:
                await asset_add_project(root_domain, asset_collection_list[a], a, db, project_id)


async def asset_add_project(root_domain, db_key, doc_name, db, project_id):
    regex_patterns = [f".*{domain}.*" for domain in root_domain]
    pattern = "|".join(regex_patterns)
    # 构建查询条件
    query = {
        "$and": [
            {
                "$or": [
                    {key: {"$regex": pattern, "$options": "i"}} for key in db_key
                ]
            },
            {"project": {"$exists": True, "$eq": ""}}
        ]
    }
    update_query = {
        "$set": {
            "project": project_id
        }
    }
    result = await db[doc_name].update_many(query, update_query)
    # 打印更新的文档数量
    logger.info(f"Updated {doc_name} {result.modified_count} documents")


async def asset_update_project(root_domain, db_key, doc_name, db, project_id):
    regex_patterns = [f".*{domain}.*" for domain in root_domain]
    pattern = "|".join(regex_patterns)
    # 构建查询条件
    query = {
        "$and": [
            {"project": project_id},
            {
                "$nor": [
                    {key: {"$regex": pattern, "$options": "i"}} for key in db_key
                ]
            }
        ]
    }
    update_query = {
        "$set": {
            "project": ""
        }
    }
    result = await db[doc_name].update_many(query, update_query)
    # 打印更新的文档数量
    logger.info(f"Updated {doc_name} {result.modified_count} documents to null ")
    await asset_add_project(root_domain, db_key, doc_name, db, project_id)


async def delete_asset_project(db, collection, project_id):
    try:
        query = {"project": project_id}

        cursor = db[collection].find(query)

        async for document in cursor:
            await db[collection].update_one({"_id": document["_id"]}, {"$set": {"project": ""}})

    except Exception as e:
        logger.error(f"delete_asset_project error:{e}")


async def delete_asset_project_handler(project_id):
    async for db in get_mongo_db():
        asset_collection_list = ['asset', 'subdomain', 'DirScanResult', 'vulnerability', 'SubdoaminTakerResult',
                                 'PageMonitoring', 'SensitiveResult', 'UrlScan', 'crawler']
        for c in asset_collection_list:
            await delete_asset_project(db, c, project_id)


async def scheduler_project(id):
    logger.info(f"Scheduler project {id}")
    async for db in get_mongo_db():
        async for redis in get_redis_pool():
            job = scheduler.get_job(id)
            task_id = generate_random_string(15)
            if job:
                next_time = job.next_run_time
                formatted_time = next_time.strftime("%Y-%m-%d %H:%M:%S")
                doc = await db.ScheduledTasks.find_one({"id": id})
                run_id_last = doc.get("runner_id", "")
                if run_id_last != "":
                    progresskeys = await redis.keys(f"TaskInfo:progress:{run_id_last}:*")
                    for pgk in progresskeys:
                        await redis.delete(pgk)
                update_document = {
                    "$set": {
                        "lastTime": get_now_time(),
                        "nextTime": formatted_time,
                        "runner_id": task_id
                    }
                }
                await db.ScheduledTasks.update_one({"id": id}, update_document)
            query = {"_id": ObjectId(id)}
            doc = await db.project.find_one(query)
            targetList = []
            target_data = await db.ProjectTargetData.find_one({"id": id})
            for t in target_data.get('target', '').split("\n"):
                t.replace("http://", "").replace("https://", "")
                t = t.strip("\n").strip("\r").strip()
                if t != "" and t not in targetList:
                    targetList.append(t)
            await create_scan_task(doc, task_id, targetList, redis)
