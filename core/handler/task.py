# -------------------------------------
# @file      : task.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/11/6 22:02
# -------------------------------------------
import re

from bson import ObjectId
from motor.motor_asyncio import AsyncIOMotorCursor
from loguru import logger


async def get_task_data(db, request_data, id):
    # 获取模板数据
    template_data = await db.ScanTemplates.find_one({"_id": ObjectId(request_data["template"])})
    # 如果选择了poc 将poc参数拼接到nuclei的参数中
    if len(template_data['vullist']) != 0:
        vul_tmp = ""
        if "All Poc" in template_data['vullist']:
            vul_tmp = "*"
        else:
            for vul in template_data['vullist']:
                vul_tmp += vul + ".yaml" + ","
        vul_tmp = vul_tmp.strip(",")

        if "VulnerabilityScan" not in template_data["Parameters"]:
            template_data["Parameters"]["VulnerabilityScan"] = {"ed93b8af6b72fe54a60efdb932cf6fbc": ""}
        if "ed93b8af6b72fe54a60efdb932cf6fbc" not in template_data["Parameters"]["VulnerabilityScan"]:
            template_data["Parameters"]["VulnerabilityScan"]["ed93b8af6b72fe54a60efdb932cf6fbc"] = ""

        if "ed93b8af6b72fe54a60efdb932cf6fbc" in template_data["VulnerabilityScan"]:
            template_data["Parameters"]["VulnerabilityScan"]["ed93b8af6b72fe54a60efdb932cf6fbc"] = \
                template_data["Parameters"]["VulnerabilityScan"][
                    "ed93b8af6b72fe54a60efdb932cf6fbc"] + " -t " + vul_tmp
    # 解析参数，支持{}获取字典
    template_data["Parameters"] = await parameter_parser(template_data["Parameters"], db)
    # 删除原始的vullist
    del template_data["vullist"]
    del template_data["_id"]
    # 设置任务名称
    template_data["TaskName"] = request_data["name"]
    # 设置忽略目标
    template_data["ignore"] = request_data["ignore"]
    # 设置去重
    template_data["duplicates"] = request_data["duplicates"]
    # 任务id
    template_data["ID"] = str(id)
    # 任务类型
    template_data["type"] = request_data.get("type", "scan")
    # 是否暂停后开启
    template_data["IsStart"] = request_data.get("IsStart", False)
    return template_data


async def parameter_parser(parameter, db):
    dict_list = {}
    port_list = {}
    # 获取字典
    cursor: AsyncIOMotorCursor = db["dictionary"].find({})
    result = await cursor.to_list(length=None)
    for doc in result:
        dict_list[f'{doc["category"].lower()}.{doc["name"].lower()}'] = str(doc['_id'])
    # 获取端口
    cursor: AsyncIOMotorCursor = db.PortDict.find({})
    result = await cursor.to_list(length=None)
    for doc in result:
        port_list[f'{doc["name"].lower()}'] = doc["value"]

    for module_name in parameter:
        for plugin in parameter[module_name]:
            matches = re.findall(r'\{(.*?)\}', parameter[module_name][plugin])
            for match in matches:
                tp, value = match.split(".", 1)
                if tp == "dict":
                    if value.lower() in dict_list:
                        real_param = dict_list[value.lower()]
                    else:
                        real_param = match
                        logger.error(f"parameter error:module {module_name} plugin {plugin}  parameter {parameter[module_name][plugin]}")
                    parameter[module_name][plugin] = parameter[module_name][plugin].replace("{" + match + "}", real_param)
                elif tp == "port":
                    if value.lower() in port_list:
                        real_param = port_list[value.lower()]
                    else:
                        real_param = match
                        logger.error(
                            f"parameter error:module {module_name} plugin {plugin}  parameter {parameter[module_name][plugin]}")
                    parameter[module_name][plugin] = parameter[module_name][plugin].replace("{" + match + "}", real_param)
    return parameter