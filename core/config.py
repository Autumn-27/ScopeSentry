# -*- coding:utf-8 -*-　　
# @name: config
# @auth: rainy-autumn@outlook.com
# @version:
import os
import random
import string

import yaml

VERSION = "1.7"
UPDATEURL = "http://update.scope-sentry.top"
REMOTE_REPO_URL = "https://github.com/Autumn-27/ScopeSentry.git"
SECRET_KEY = "ScopeSentry-15847412364125411"
MONGODB_IP = ""
MONGODB_PORT = 0
MONGODB_DATABASE = ""
MONGODB_USER = ''
MONGODB_PASSWORD = ''
REDIS_IP = ""
REDIS_PORT = ""
REDIS_PASSWORD = ""
TIMEZONE = 'Asia/Shanghai'
LOG_INFO = {}
GET_LOG_NAME = []
NODE_TIMEOUT = 50
TOTAL_LOGS = 1000
APP = {}
Project_List = {}
PLUGINKEY = ""

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
    global MONGODB_IP, MONGODB_PORT, MONGODB_DATABASE, MONGODB_USER, MONGODB_PASSWORD, REDIS_IP, REDIS_PORT, REDIS_PASSWORD, SECRET_KEY, TOTAL_LOGS, TIMEZONE
    SECRET_KEY = generate_random_string(16)
    global PLUGINKEY
    if os.path.exists("PLUGINKEY"):
        with open("PLUGINKEY", 'r') as file:
            PLUGINKEY = file.read()
    else:
        PLUGINKEY = generate_random_string(6)
        with open("PLUGINKEY", 'w') as file:
            file.write(PLUGINKEY)
    config_file_path = "config.yaml"
    if os.path.exists(config_file_path):
        with open(config_file_path, 'r') as file:
            data = yaml.safe_load(file)
            MONGODB_IP = data['mongodb']['ip']
            MONGODB_PORT = data['mongodb']['port']
            MONGODB_DATABASE = data['mongodb']['mongodb_database']
            MONGODB_USER = data['mongodb']['username']
            MONGODB_PASSWORD = data['mongodb']['password']
            REDIS_IP = data['redis']['ip']
            REDIS_PORT = data['redis']['port']
            REDIS_PASSWORD = data['redis']['password']
            TOTAL_LOGS = data['logs']['total_logs']
            TIMEZONE = data['system']['timezone']
        env_db_user = os.environ.get("MONGODB_USER", default='')
        if env_db_user != '' and env_db_user != MONGODB_USER:
            MONGODB_USER = env_db_user
        env_db_password = os.environ.get("MONGODB_PASSWORD", default='')
        if env_db_password != '' and env_db_password != MONGODB_PASSWORD:
            MONGODB_PASSWORD = env_db_password
        env_redis_password = os.environ.get("REDIS_PASSWORD", default='')
        if env_redis_password != '' and env_redis_password != REDIS_PASSWORD:
            REDIS_PASSWORD = env_redis_password
    else:
        TIMEZONE = os.environ.get("TIMEZONE", default='Asia/Shanghai')
        MONGODB_IP = os.environ.get("MONGODB_IP", default='127.0.0.1')
        MONGODB_PORT = int(os.environ.get("MONGODB_PORT", default=27017))
        MONGODB_DATABASE = os.environ.get("MONGODB_DATABASE", default='ScopeSentry')
        MONGODB_USER = os.environ.get("MONGODB_USER", default='root')
        MONGODB_PASSWORD = os.environ.get("MONGODB_PASSWORD", default='QckSdkg5CKvtxfec')
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
                'mongodb_database': MONGODB_DATABASE,
                'username': MONGODB_USER,
                'password': MONGODB_PASSWORD,
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
