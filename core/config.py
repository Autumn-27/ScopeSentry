# -*- coding:utf-8 -*-　　
# @name: config
# @auth: rainy-autumn@outlook.com
# @version:
import os
import random
import string

import yaml

VERSION = "1.0"
UPDATEURL = "http://update.scope-sentry.top"
REMOTE_REPO_URL = "https://github.com/Autumn-27/ScopeSentry.git"
SECRET_KEY = "ScopeSentry-15847412364125411"
MONGODB_IP = ""
MONGODB_PORT = 0
DATABASE_NAME = ""
DATABASE_USER = ''
DATABASE_PASSWORD = ''
REDIS_IP = ""
REDIS_PORT = ""
REDIS_PASSWORD = ""
TIMEZONE = 'Asia/Shanghai'
LOG_INFO = {}
GET_LOG_NAME = []
NODE_TIMEOUT = 50
TOTAL_LOGS = 1000
APP = {}
SensitiveRuleList = {}
POC_LIST = {}
Project_List = {}
def set_timezone(t):
    global TIMEZONE
    TIMEZONE = t


def get_timezone():
    global TIMEZONE
    return TIMEZONE


def generate_random_string(length):
    # 生成随机字符串，包括大小写字母和数字
    characters = string.ascii_letters + string.digits
    random_string = ''.join(random.choice(characters) for _ in range(length))
    return random_string


def set_config():
    global MONGODB_IP, MONGODB_PORT, DATABASE_NAME, DATABASE_USER, DATABASE_PASSWORD, REDIS_IP, REDIS_PORT, REDIS_PASSWORD, SECRET_KEY, TOTAL_LOGS, TIMEZONE
    SECRET_KEY = generate_random_string(16)
    config_file_path = "config.yaml"
    if os.path.exists(config_file_path):
        with open(config_file_path, 'r') as file:
            data = yaml.safe_load(file)
            MONGODB_IP = data['mongodb']['ip']
            MONGODB_PORT = data['mongodb']['port']
            DATABASE_NAME = data['mongodb']['database_name']
            DATABASE_USER = data['mongodb']['username']
            DATABASE_PASSWORD = data['mongodb']['password']
            REDIS_IP = data['redis']['ip']
            REDIS_PORT = data['redis']['port']
            REDIS_PASSWORD = data['redis']['password']
            TOTAL_LOGS = data['logs']['total_logs']
            TIMEZONE = data['system']['timezone']
    else:
        TIMEZONE = os.environ.get("TIMEZONE", default='Asia/Shanghai')
        MONGODB_IP = os.environ.get("MONGODB_IP", default='127.0.0.1')
        MONGODB_PORT = int(os.environ.get("MONGODB_PORT", default=27017))
        DATABASE_NAME = os.environ.get("DATABASE_NAME", default='ScopeSentry')
        DATABASE_USER = os.environ.get("DATABASE_USER", default='root')
        DATABASE_PASSWORD = os.environ.get("DATABASE_PASSWORD", default='QckSdkg5CKvtxfec')
        REDIS_IP = os.environ.get("REDIS_IP", default='127.0.0.1')
        REDIS_PORT = os.environ.get("REDIS_PORT", default="6379")
        REDIS_PASSWORD = os.environ.get("REDIS_PASSWORD", default='ScopeSentry')
        TOTAL_LOGS = 1000
        config_data = {
            'system': {
                'timezone': TIMEZONE
            },
            'mongodb': {
                'ip': MONGODB_IP,
                'port': int(MONGODB_PORT),
                'database_name': DATABASE_NAME,
                'username': DATABASE_USER,
                'password': DATABASE_PASSWORD,
            },
            'redis': {
                'ip': REDIS_IP,
                'port': REDIS_PORT,
                'password': REDIS_PASSWORD,
            },
            'logs': {
                'total_logs': TOTAL_LOGS
            }
        }
        with open(config_file_path, 'w') as file:
            yaml.dump(config_data, file)
