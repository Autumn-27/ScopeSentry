# -*- coding:utf-8 -*-　　
# @name: util
# @auth: rainy-autumn@outlook.com
# @version:
import hashlib, random
import re
import string
import sys
from loguru import logger
from core.config import TIMEZONE, APP, SensitiveRuleList, Project_List
from datetime import timezone
from datetime import datetime, timedelta
import json
from urllib.parse import urlparse

from core.db import get_mongo_db


def calculate_md5_from_content(content):
    md5 = hashlib.md5()
    md5.update(content.encode("utf-8"))  # 将内容编码为 utf-8 后更新 MD5
    return md5.hexdigest()


def evaluate_expression(express):
    random_bool = random.choice([True, False])
    return str(random_bool)


def generate_random_string(length):
    # 生成随机字符串，包括大小写字母和数字
    characters = string.ascii_letters + string.digits
    random_string = ''.join(random.choice(characters) for _ in range(length))
    return random_string


def is_valid_string(s):
    # 定义合法字符集
    valid_chars = string.ascii_letters + string.digits
    # 使用正则表达式判断字符串是否仅包含合法字符
    pattern = f"^[{re.escape(valid_chars)}]+$"
    return bool(re.match(pattern, s))

def parse_expression(express, eval_expression):
    parts = []
    part = ""
    operator_flag = False
    parentheses_depth = 0
    for i in range(len(express)):
        if express[i] == '(':
            if i != 0:
                if express[i - 1] != '\\':
                    parentheses_depth += 1
        elif express[i] == ')':
            if i != 0:
                if express[i - 1] != '\\':
                    parentheses_depth -= 1

        if express[i] == '|' and express[i + 1] == '|' and parentheses_depth == 0:
            operator_flag = True
            if part[0] == '(':
                eval_expression += "("
                eval_expression = parse_expression(part.strip("(").strip(")"), eval_expression)
                eval_expression += ") or "
            else:
                eval_expression += evaluate_expression(part) + " or "
            part = ""
        elif express[i] == '&' and express[i + 1] == '&' and parentheses_depth == 0:
            operator_flag = True
            if part[0] == '(':
                eval_expression += "("
                eval_expression = parse_expression(part.strip("(").strip(")"), eval_expression)
                eval_expression += ") and "
            else:
                eval_expression += evaluate_expression(part) + " and "
            part = ""
        else:
            ch = ""
            if operator_flag:
                ch = express[i + 1]
                operator_flag = False
            else:
                ch = express[i]
            part += ch.strip()
    if part[0] == '(':
        eval_expression += "("
        eval_expression = parse_expression(part.strip("(").strip(")"), eval_expression)
        eval_expression += ")"
    else:
        eval_expression += evaluate_expression(part)
    return eval_expression


def get_now_time():
    TZ = timezone(
        timedelta(hours=8),
        name=TIMEZONE,
    )
    utc_now = datetime.utcnow().replace(tzinfo=timezone.utc)
    time_now = utc_now.astimezone(TZ)
    formatted_time = time_now.strftime("%Y-%m-%d %H:%M:%S")
    return formatted_time


def read_json_file(file_path):
    with open(file_path, encoding='utf-8') as f:
        data = json.load(f)
    return data


def transform_db_redis(request_data):
    Subfinder = False
    Ksubdomain = False
    if "Subfinder" in request_data["subdomainConfig"]:
        Subfinder = True
    if "Ksubdomain" in request_data["subdomainConfig"]:
        Ksubdomain = True
    add_redis_task_data = {
        "TaskId": request_data["id"],
        "SubdomainScan": request_data["subdomainScan"],
        "Subfinder": Subfinder,
        "Ksubdomain": Ksubdomain,
        "UrlScan": request_data["urlScan"],
        "Duplicates": request_data["duplicates"],
        "SensitiveInfoScan": request_data["sensitiveInfoScan"],
        "PageMonitoring": request_data["pageMonitoring"],
        "CrawlerScan": request_data["crawlerScan"],
        "VulScan": request_data["vulScan"],
        "VulList": request_data["vulList"],
        "PortScan": request_data["portScan"],
        "Ports": request_data["ports"],
        "Waybackurl": request_data["waybackurl"],
        "DirScan": request_data["dirScan"],
        "type": 'scan'
    }
    return add_redis_task_data


def string_to_postfix(expression):
    try:
        operands_stack = []
        expression_stack = []
        start_char = 0
        skip_flag = False
        exp_flag = False
        for index, char in enumerate(expression):
            if skip_flag:
                skip_flag = False
                continue
            if char == '|' and expression[index + 1] == '|':
                skip_flag = True
                operands_stack.append("||")
                key = expression[start_char:index]
                if key != "":
                    expression_stack.append(key)
                start_char = index + 2
            elif char == '&' and expression[index + 1] == '&':
                skip_flag = True
                operands_stack.append("&&")
                key = expression[start_char:index]
                if key != "":
                    expression_stack.append(key)
                start_char = index + 2
            elif char == '(' and expression[index - 1] != '\\' and exp_flag != True:
                start_char = index + 1
                operands_stack.append('(')
            elif char == ')' and expression[index - 1] != '\\' and exp_flag != True:
                key = expression[start_char:index]
                if key != "":
                    expression_stack.append(key)
                start_char = index + 1
                popped_value = operands_stack.pop()
                while popped_value != '(':
                    if popped_value != '(':
                        if popped_value != "":
                            expression_stack.append(popped_value)
                    popped_value = operands_stack.pop()
            elif char == " ":
                continue
            elif char == "\"" and expression[index - 1] != "\\":
                if exp_flag == False:
                    exp_flag = True
                else:
                    if index == len(expression):
                        exp_flag = False
                        continue
                    tmp = expression[index:].replace(" ", "")
                    if tmp.startswith("\"||") or (tmp.startswith("\"))") and len(tmp) == 3) or tmp.startswith(
                            "\"&&") or tmp.startswith("\")||") or tmp.startswith("\")&&") or (
                            tmp.startswith("\")") and len(tmp) == 2) or re.findall(r"^\"[)]*(\|\||\&\&)", tmp):
                        exp_flag = False
        if start_char != len(expression):
            key = expression[start_char:]
            if key != "":
                expression_stack.append(key)
        while len(operands_stack) != 0:
            expression_stack.append(operands_stack.pop())
        tmp = []
        for key in expression_stack:
            if key != "" and key != " ":
                tmp.append(
                    key.strip().replace('\(', '(').replace('\)', ')').replace('\|\|', '||').replace('\&\&', '&&'))
        return tmp
    except Exception as e:
        logger.error(f"后缀表达式转换出错：{expression}")
        return ""


async def search_to_mongodb(expression_raw, keyword):
    try:
        keyword["task"] = "taskId"
        if expression_raw == "":
            return [{}]
        if len(APP) == 0:
            logger.error("WebFinger缓存数据为0，请排查~")
        expression = string_to_postfix(expression_raw)
        stack = []
        for expr in expression:
            if expr == "&&":
                right = stack.pop()
                left = stack.pop()
                stack.append({"$and": [left, right]})
            elif expr == "||":
                right = stack.pop()
                left = stack.pop()
                stack.append({"$or": [left, right]})
            elif "!=" in expr:
                key, value = expr.split("!=", 1)
                key = key.strip()
                if key in keyword:
                    value = value.strip("\"")
                    if key == 'statuscode' or key == 'length':
                        value = int(value)
                    if key == 'project':
                        if value.lower() in Project_List:
                            value = Project_List[value.lower()]
                    if key == 'app':
                        finger_id = []
                        for ap_key in APP:
                            if value.lower() in APP[ap_key].lower():
                                finger_id.append(ap_key)
                        tmp_nor = {"$nor": []}
                        for f_i in finger_id:
                            tmp_nor['$nor'].append({"webfinger": {"$in": [f_i]}})
                        tmp_nor['$nor'].append({"technologies": {"$regex": value, "$options": "i"}})
                        stack.append(tmp_nor)
                    if type(keyword[key]) is list:
                        tmp_nor = {"$nor": []}
                        for v in keyword[key]:
                            tmp_nor['$nor'].append({v: {"$regex": value, "$options": "i"}})
                        stack.append(tmp_nor)
                    else:
                        tmp_nor = {"$nor": []}
                        if type(value) is int:
                            tmp_nor['$nor'].append({keyword[key]: {"$eq": value}})
                        else:
                            tmp_nor['$nor'].append({keyword[key]: {"$regex": value, "$options": "i"}})
                        stack.append(tmp_nor)
            elif "==" in expr:
                key, value = expr.split("==", 1)
                key = key.strip()
                if key in keyword:
                    value = value.strip("\"")
                    if key == "task":
                        async for db in get_mongo_db():
                            query = {"name": {"$eq": value}}
                            doc = await db.task.find_one(query)
                            if doc is not None:
                                taskid = str(doc.get("_id"))
                                value = taskid
                    if key == 'statuscode' or key == 'length':
                        value = int(value)
                    if key == 'project':
                        if value.lower() in Project_List:
                            value = Project_List[value.lower()]
                    if key == 'app':
                        finger_id = []
                        for ap_key in APP:
                            if value.lower() == APP[ap_key].lower():
                                finger_id.append(ap_key)
                        tmp_or = {"$or": []}
                        for f_i in finger_id:
                            tmp_or['$or'].append({"webfinger": {"$in": [f_i]}})
                        tmp_or['$or'].append({"technologies": {"$eq": value}})
                        stack.append(tmp_or)
                    if type(keyword[key]) is list:
                        tmp_or = {"$or": []}
                        for v in keyword[key]:
                            tmp_or['$or'].append({v: {"$eq": value}})
                        stack.append(tmp_or)
                    else:
                        tmp_or = {keyword[key]: {"$eq": value}}
                        stack.append(tmp_or)
            elif "=" in expr:
                key, value = expr.split("=", 1)
                key = key.strip()
                if key in keyword:
                    value = value.strip("\"")
                    if key == 'project':
                        if value.lower() in Project_List:
                            value = Project_List[value.lower()]
                    if key == 'app':
                        finger_id = []
                        for ap_key in APP:
                            if value.lower() in APP[ap_key].lower():
                                finger_id.append(ap_key)
                        tmp_or = {"$or": []}
                        for f_i in finger_id:
                            tmp_or['$or'].append({"webfinger": {"$in": [f_i]}})
                        tmp_or['$or'].append({"technologies": {"$regex": value, "$options": "i"}})
                        stack.append(tmp_or)
                    if type(keyword[key]) is list:
                        tmp_or = {"$or": []}
                        for v in keyword[key]:
                            tmp_or['$or'].append({v: {"$regex": value, "$options": "i"}})
                        stack.append(tmp_or)
                    else:
                        stack.append({keyword[key]: {"$regex": value, "$options": "i"}})
        return stack
    except Exception as e:
        logger.error(e)
        return ""

async def get_search_query(name, request_data):
    search_query = request_data.get("search", "")
    search_key_v = {
        'sens':{
            'url': 'url',
            'sname': 'sid',
            "body": "body",
            "info": "match",
            'project': 'project',
            'md5': 'md5'
        },
        'dir': {
            'project': 'project',
            'statuscode': 'status',
            'url': 'url',
            'redirect': 'msg',
            'length': 'length'
        },
        'vul': {
            'url': 'url',
            'vulname': 'vulname',
            'project': 'project',
            'matched': 'matched',
            'request': 'request',
            'response': 'response',
            'level': 'level'
        }
    }
    keyword = search_key_v[name]
    query = await search_to_mongodb(search_query, keyword)
    if query == "" or query is None:
        return ""
    query = query[0]
    filter_key = ['color', 'status', 'level']
    filter = request_data.get("filter", {})
    if filter:
        query["$and"] = []
        for f in filter:
            if f in filter_key:
                tmp_or = []
                for v in filter[f]:
                    tmp_or.append({f: v})
                if len(tmp_or) != 0:
                    query["$and"].append({"$or": tmp_or})
    if "$and" in query:
        if len(query["$and"]) == 0:
            query.pop("$and")
    return query


def get_root_domain(url):
    # 如果URL不带协议，添加一个默认的http协议
    global root_domain
    if not url.startswith(('http://', 'https://')):
        url = 'http://' + url

    parsed_url = urlparse(url)

    # 检查是否为IP地址
    try:
        # 使用ip_address来检查
        from ipaddress import ip_address
        ip_address(parsed_url.netloc)
        return parsed_url.netloc  # 如果是IP地址，直接返回
    except ValueError:
        pass

    domain_parts = parsed_url.netloc.split('.')

    # 复合域名列表
    compound_domains = [
    'com.cn', 'net.cn', 'org.cn', 'gov.cn', 'edu.cn', 'ac.cn', 'mil.cn',
    'co.uk', 'org.uk', 'net.uk', 'gov.uk', 'ac.uk', 'sch.uk',
    'co.jp', 'ne.jp', 'or.jp', 'go.jp', 'ac.jp', 'ad.jp',
    'com.de', 'org.de', 'net.de', 'gov.de',
    'com.ca', 'net.ca', 'org.ca', 'gov.ca',
    'com.au', 'net.au', 'org.au', 'gov.au', 'edu.au',
    'com.fr', 'net.fr', 'org.fr', 'gov.fr',
    'com.br', 'com.mx', 'com.ar', 'com.ru',
    'co.in', 'co.za',
    'co.kr', 'com.tw'
]

    # 检查是否为复合域名
    is_compound_domain = False
    for compound_domain in compound_domains:
        if domain_parts[-2:] == compound_domain.split('.'):
            is_compound_domain = True
            root_domain = '.'.join(domain_parts[-3:])
            break

    if not is_compound_domain:
        root_domain = '.'.join(domain_parts[-2:])

    return root_domain
