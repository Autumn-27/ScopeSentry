# -*- coding:utf-8 -*-　　
# @name: default
# @auth: rainy-autumn@outlook.com
# @version:
import json
import os
from urllib.parse import urlparse

from bson import ObjectId

from loguru import logger
current_directory = os.getcwd()

dict_directory = "dicts"

combined_directory = os.path.join(current_directory, dict_directory)


def read_json_file(file_path):
    with open(file_path, encoding='utf-8') as f:
        data = json.load(f)
    return data


def get_domainDict():
    domainDict = ""
    try:
        # 尝试打开文件并读取内容
        with open(os.path.join(combined_directory, "domainDict"), "r", encoding="utf-8") as file:
            domainDict = file.read()
    except FileNotFoundError:
        logger.error("文件不存在")
    return domainDict


def get_dirDict():
    domainDict = ""
    try:
        # 尝试打开文件并读取内容
        with open(os.path.join(combined_directory, "dirDict"), "r", encoding="utf-8") as file:
            domainDict = file.read()
    except FileNotFoundError:
        logger.error("文件不存在")
    return domainDict


def get_poc():
    pocPath = os.path.join(combined_directory, "ScopeSentry.PocList.json")
    data = read_json_file(pocPath)
    for d in data:
        d.pop('_id', None)
    return data


def get_finger():
    fingerPath = os.path.join(combined_directory, "ScopeSentry.FingerprintRules.json")
    data = read_json_file(fingerPath)
    for d in data:
        d.pop('_id', None)
    return data


def get_project_data():
    project_path = os.path.join(combined_directory, "ScopeSentry.project.json")
    data = read_json_file(project_path)
    target_data = []
    project_data = []
    for d in data:
        project_id = d['_id']['$oid']
        tmp = []
        for t in d['target'].split('\n'):
            root_domain = get_root_domain(t).replace("\n", "").replace("\r", "").strip()
            if root_domain is not None and root_domain != "":
                if root_domain not in tmp:
                    tmp.append(root_domain)
        d["root_domains"] = tmp
        d['_id'] = ObjectId(project_id)
        target_data.append({"id": project_id, "target": d['target']})
        del d["target"]
        project_data.append(d)
    return project_data, target_data


def get_sensitive():
    sensitive_path = os.path.join(combined_directory, "ScopeSentry.SensitiveRule.json")
    data = read_json_file(sensitive_path)
    for d in data:
        d.pop('_id', None)
    return data


SEARCHKEY = {
        'SensitiveResult': {
            'url': 'url',
            'sname': 'sid',
            "body": "body",
            "info": "match",
            'md5': 'md5'
        },
        'DirScanResult': {
            'statuscode': 'status',
            'url': 'url',
            'redirect': 'msg',
            'length': 'length'
        },
        'vulnerability': {
            'url': 'url',
            'vulname': 'vulname',
            'matched': 'matched',
            'request': 'request',
            'response': 'response',
            'level': 'level'
        },
        'subdomain': {
            'domain': 'host',
            'ip': 'ip',
            'type': 'type',
            'value': 'value'
        },
        'asset': {
            'app': 'technologies',
            'body': 'body',
            'header': 'rawheaders',
            'title': 'title',
            'statuscode': 'statuscode',
            'icon': 'faviconmmh3',
            'ip': "ip",
            'domain': "host",
            'port': 'port',
            'service': "service",
            'banner': 'raw',
            'type': 'type',
        },
        'SubdoaminTakerResult': {
            'domain': 'input',
            'value': 'value',
            'type': 'cname',
            'response': 'response',
        },
        'UrlScan': {
            'url': 'output',
            'input': 'input',
            'source': 'source',
            "resultId": 'resultId',
            "type": "outputtype"
        },
        'PageMonitoring': {
            'url': 'url',
            'hash': 'hash',
            'diff': 'diff',
            'response': 'response'
        },
        'crawler': {
            'url': 'url',
            'method': 'method',
            'body': 'body',
            "resultId": 'resultId',
        },
        'RootDomain': {
            'domain': 'domain',
            'icp': 'icp',
            'company': 'company'
        },
        'app': {
            'name': 'name',
            'icp': 'icp',
            'company': 'company',
            'category': 'category',
            'description': 'description',
            'url': 'url',
            'apk': 'apk'
        },
        'mp': {
            'name': 'name',
            'icp': 'icp',
            'company': 'company',
            'category': 'category',
            'description': 'description',
            'url': 'url',
        }
    }

subfinderApiConfig = '''# subfinder can be used right after the installation, however many sources required API keys to work. Learn more here: https://docs.projectdiscovery.io/tools/subfinder/install#post-install-configuration.
bevigil: []
binaryedge: []
bufferover: []
builtwith: []
c99: []
censys: []
certspotter: []
chaos: []
chinaz: []
dnsdb: []
dnsrepo: []
facebook: []
fofa: []
fullhunt: []
github: []
hunter: []
intelx: []
leakix: []
netlas: []
passivetotal: []
quake: []
redhuntlabs: []
robtex: []
securitytrails: []
shodan: []
threatbook: []
virustotal: []
whoisxmlapi: []
zoomeyeapi: []
'''


portDic = [
    {'name': '100个常见端口',
     'value': '21,22,23,25,53,67,68,80,110,111,139,143,161,389,443,445,465,512,513,514,873,993,995,1080,1000,1352,1433,1521,1723,2049,2181,2375,3306,3389,4848,5000,5001,5432,5900,5632,5900,5989,6379,6666,7001,7002,8000,8001,8009,8010,8069,8080,8083,8086,8081,8088,8089,8443,8888,9900,9200,9300,9999,10621,11211,27017,27018,66,81,457,1100,1241,1434,1944,2301,3128,4000,4001,4002,4100,5800,5801,5802,6346,6347,30821,1090,1098,1099,4444,11099,47001,47002,10999,7000-7004,8000-8003,9000-9003,9503,7070,7071,45000,45001,8686,9012,50500,11111,4786,5555,5556,8880,8983,8383,4990,8500,6066'},
    {'name': 'nmap top 1000',
     'value': '1,3-4,6-7,9,13,17,19-26,30,32-33,37,42-43,49,53,70,79-85,88-90,99-100,106,109-111,113,119,125,135,139,143-144,146,161,163,179,199,211-212,222,254-256,259,264,280,301,306,311,340,366,389,406-407,416-417,425,427,443-445,458,464-465,481,497,500,512-515,524,541,543-545,548,554-555,563,587,593,616-617,625,631,636,646,648,666-668,683,687,691,700,705,711,714,720,722,726,749,765,777,783,787,800-801,808,843,873,880,888,898,900-903,911-912,981,987,990,992-993,995,999-1002,1007,1009-1011,1021-1100,1102,1104-1108,1110-1114,1117,1119,1121-1124,1126,1130-1132,1137-1138,1141,1145,1147-1149,1151-1152,1154,1163-1166,1169,1174-1175,1183,1185-1187,1192,1198-1199,1201,1213,1216-1218,1233-1234,1236,1244,1247-1248,1259,1271-1272,1277,1287,1296,1300-1301,1309-1311,1322,1328,1334,1352,1417,1433-1434,1443,1455,1461,1494,1500-1501,1503,1521,1524,1533,1556,1580,1583,1594,1600,1641,1658,1666,1687-1688,1700,1717-1721,1723,1755,1761,1782-1783,1801,1805,1812,1839-1840,1862-1864,1875,1900,1914,1935,1947,1971-1972,1974,1984,1998-2010,2013,2020-2022,2030,2033-2035,2038,2040-2043,2045-2049,2065,2068,2099-2100,2103,2105-2107,2111,2119,2121,2126,2135,2144,2160-2161,2170,2179,2190-2191,2196,2200,2222,2251,2260,2288,2301,2323,2366,2381-2383,2393-2394,2399,2401,2492,2500,2522,2525,2557,2601-2602,2604-2605,2607-2608,2638,2701-2702,2710,2717-2718,2725,2800,2809,2811,2869,2875,2909-2910,2920,2967-2968,2998,3000-3001,3003,3005-3007,3011,3013,3017,3030-3031,3052,3071,3077,3128,3168,3211,3221,3260-3261,3268-3269,3283,3300-3301,3306,3322-3325,3333,3351,3367,3369-3372,3389-3390,3404,3476,3493,3517,3527,3546,3551,3580,3659,3689-3690,3703,3737,3766,3784,3800-3801,3809,3814,3826-3828,3851,3869,3871,3878,3880,3889,3905,3914,3918,3920,3945,3971,3986,3995,3998,4000-4006,4045,4111,4125-4126,4129,4224,4242,4279,4321,4343,4443-4446,4449,4550,4567,4662,4848,4899-4900,4998,5000-5004,5009,5030,5033,5050-5051,5054,5060-5061,5080,5087,5100-5102,5120,5190,5200,5214,5221-5222,5225-5226,5269,5280,5298,5357,5405,5414,5431-5432,5440,5500,5510,5544,5550,5555,5560,5566,5631,5633,5666,5678-5679,5718,5730,5800-5802,5810-5811,5815,5822,5825,5850,5859,5862,5877,5900-5904,5906-5907,5910-5911,5915,5922,5925,5950,5952,5959-5963,5987-5989,5998-6007,6009,6025,6059,6100-6101,6106,6112,6123,6129,6156,6346,6389,6502,6510,6543,6547,6565-6567,6580,6646,6666-6669,6689,6692,6699,6779,6788-6789,6792,6839,6881,6901,6969,7000-7002,7004,7007,7019,7025,7070,7100,7103,7106,7200-7201,7402,7435,7443,7496,7512,7625,7627,7676,7741,7777-7778,7800,7911,7920-7921,7937-7938,7999-8002,8007-8011,8021-8022,8031,8042,8045,8080-8090,8093,8099-8100,8180-8181,8192-8194,8200,8222,8254,8290-8292,8300,8333,8383,8400,8402,8443,8500,8600,8649,8651-8652,8654,8701,8800,8873,8888,8899,8994,9000-9003,9009-9011,9040,9050,9071,9080-9081,9090-9091,9099-9103,9110-9111,9200,9207,9220,9290,9415,9418,9485,9500,9502-9503,9535,9575,9593-9595,9618,9666,9876-9878,9898,9900,9917,9929,9943-9944,9968,9998-10004,10009-10010,10012,10024-10025,10082,10180,10215,10243,10566,10616-10617,10621,10626,10628-10629,10778,11110-11111,11967,12000,12174,12265,12345,13456,13722,13782-13783,14000,14238,14441-14442,15000,15002-15004,15660,15742,16000-16001,16012,16016,16018,16080,16113,16992-16993,17877,17988,18040,18101,18988,19101,19283,19315,19350,19780,19801,19842,20000,20005,20031,20221-20222,20828,21571,22939,23502,24444,24800,25734-25735,26214,27000,27352-27353,27355-27356,27715,28201,30000,30718,30951,31038,31337,32768-32785,33354,33899,34571-34573,35500,38292,40193,40911,41511,42510,44176,44442-44443,44501,45100,48080,49152-49161,49163,49165,49167,49175-49176,49400,49999-50003,50006,50300,50389,50500,50636,50800,51103,51493,52673,52822,52848,52869,54045,54328,55055-55056,55555,55600,56737-56738,57294,57797,58080,60020,60443,61532,61900,62078,63331,64623,64680,65000,65129,65389,280,4567,7001,8008,9080'
     },
    {
        "name": "All Ports",
        "value": "1-65535"
    }
]
radConfig = '''exec_path: ""                     # 启动chrome的路径
disable_headless: false           # 禁用无头模式
subdomain: false                   # 是否自动爬取子域
leakless: true                    # 实验性功能，防止内存泄露，可能造成卡住的现象
force_sandbox: false              # 强制开启sandbox；为 false 时默认开启沙箱，但在容器中会关闭沙箱。为true时强制启用沙箱，可能导致在docker中无法使用。
enable_image: false               # 启用图片显示
parent_path_detect: false          # 是否启用父目录探测功能
proxy: ""                         # 代理配置
user_agent: ""                    # 请求user-agent配置
domain_headers:                   # 请求头配置:[]{domain,map[headerKey]HeaderValue}
- domain: '*'                     # 为哪些域名设置header，glob语法
  headers: {}                     # 请求头，map[key]value
max_depth: 5                     # 最大页面深度限制
navigate_timeout_second: 5       # 访问超时时间，单位秒
load_timeout_second: 5           # 加载超时时间，单位秒
retry: 0                          # 页面访问失败后的重试次数
page_analyze_timeout_second: 10  # 页面分析超时时间，单位秒
max_interactive: 100             # 单个页面最大交互次数
max_interactive_depth: 5         # 页面交互深度限制
max_page_concurrent: 5           # 最大页面并发（不大于10）
max_page_visit: 1000              # 总共允许访问的页面数量
max_page_visit_per_site: 500     # 每个站点最多访问的页面数量
element_filter_strength: 3        # 过滤同站点相似元素强度，1-7取值，强度逐步增大，为0时不进行跨页面元素过滤
new_task_filter_config:           # 检查某个链接是否应该被加入爬取队列
  hostname_allowed: []            # 允许访问的 Hostname，支持格式如 t.com、*.t.com、1.1.1.1、1.1.1.1/24、1.1-4.1.1-8
  hostname_disallowed: []         # 不允许访问的 Hostname，支持格式如 t.com、*.t.com、1.1.1.1、1.1.1.1/24、1.1-4.1.1-8
  port_allowed: []                # 允许访问的端口, 支持的格式如: 80、80-85
  port_disallowed: []             # 不允许访问的端口, 支持的格式如: 80、80-85
  path_allowed: []                # 允许访问的路径，支持的格式如: test、*test*
  path_disallowed: []             # 不允许访问的路径, 支持的格式如: test、*test*
  query_key_allowed: []           # 允许访问的 Query Key，支持的格式如: test、*test*
  query_key_disallowed: []        # 不允许访问的 Query Key, 支持的格式如: test、*test*
  fragment_allowed: []            # 允许访问的 Fragment, 支持的格式如: test、*test*
  fragment_disallowed: []         # 不允许访问的 Fragment, 支持的格式如: test、*test*
  post_key_allowed: []            # 允许访问的 Post Body 中的参数, 支持的格式如: test、*test*
  post_key_disallowed: []         # 不允许访问的 Post Body 中的参数, 支持的格式如: test、*test*
request_send_filter_config:       # 检查某个请求是否应该被发送
  hostname_allowed: []            # 允许访问的 Hostname，支持格式如 t.com、*.t.com、1.1.1.1、1.1.1.1/24、1.1-4.1.1-8
  hostname_disallowed: []         # 不允许访问的 Hostname，支持格式如 t.com、*.t.com、1.1.1.1、1.1.1.1/24、1.1-4.1.1-8
  port_allowed: []                # 允许访问的端口, 支持的格式如: 80、80-85
  port_disallowed: []             # 不允许访问的端口, 支持的格式如: 80、80-85
  path_allowed: []                # 允许访问的路径，支持的格式如: test、*test*
  path_disallowed: []             # 不允许访问的路径, 支持的格式如: test、*test*
  query_key_allowed: []           # 允许访问的 Query Key，支持的格式如: test、*test*
  query_key_disallowed: []        # 不允许访问的 Query Key, 支持的格式如: test、*test*
  fragment_allowed: []            # 允许访问的 Fragment, 支持的格式如: test、*test*
  fragment_disallowed: []         # 不允许访问的 Fragment, 支持的格式如: test、*test*
  post_key_allowed: []            # 允许访问的 Post Body 中的参数, 支持的格式如: test、*test*
  post_key_disallowed: []         # 不允许访问的 Post Body 中的参数, 支持的格式如: test、*test*
request_output_filter_config:     # 检查某个请求是否应该被输出
  hostname_allowed: []            # 允许访问的 Hostname，支持格式如 t.com、*.t.com、1.1.1.1、1.1.1.1/24、1.1-4.1.1-8
  hostname_disallowed: []         # 不允许访问的 Hostname，支持格式如 t.com、*.t.com、1.1.1.1、1.1.1.1/24、1.1-4.1.1-8
  port_allowed: []                # 允许访问的端口, 支持的格式如: 80、80-85
  port_disallowed: []             # 不允许访问的端口, 支持的格式如: 80、80-85
  path_allowed: []                # 允许访问的路径，支持的格式如: test、*test*
  path_disallowed: []             # 不允许访问的路径, 支持的格式如: test、*test*
  query_key_allowed: []           # 允许访问的 Query Key，支持的格式如: test、*test*
  query_key_disallowed: []        # 不允许访问的 Query Key, 支持的格式如: test、*test*
  fragment_allowed: []            # 允许访问的 Fragment, 支持的格式如: test、*test*
  fragment_disallowed: []         # 不允许访问的 Fragment, 支持的格式如: test、*test*
  post_key_allowed: []            # 允许访问的 Post Body 中的参数, 支持的格式如: test、*test*
  post_key_disallowed: []         # 不允许访问的 Post Body 中的参数, 支持的格式如: test、*test*
entrance_retry: 0                 # 入口重试次数
max_similar_request: 0            # 最大相似fetch/XHR请求数（小于等于0时不限制）
'''


def get_fingerprint_data():
    try:
        # 尝试打开文件并读取内容
        with open(os.path.join(combined_directory, "fingerprint"), "r", encoding="utf-8") as file:
            fingerprint = file.read()
    except FileNotFoundError:
        logger.error("文件不存在")
    return json.loads(fingerprint)


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


ModulesConfig = '''maxGoroutineCount: 3 # 最大目标并发
subdomainScan:
  goroutineCount: 3  # 设置"子域名扫描"模块最大并发
subdomainSecurity:
  goroutineCount: 10  # 设置"子域名结果处理"模块最大并发
assetMapping:
  goroutineCount: 5  # 设置"资产测绘"模块最大并发
assetHandle:
  goroutineCount: 30  # 设置"资产结果处理"模块最大并发
portScanPreparation:
  goroutineCount: 30  # 设置"端口扫描预处理"模块最大并发
portScan:
  goroutineCount: 2  # 设置"端口扫描"模块最大并发
portFingerprint:
  goroutineCount: 10  # 设置"端口指纹识别"模块最大并发
URLScan:
  goroutineCount: 5  # 设置"URL扫描"模块最大并发
URLSecurity:
  goroutineCount: 15  # 设置"URL扫描结果处理"模块最大并发
webCrawler:
  goroutineCount: 2  # 设置"爬虫扫描"模块最大并发
dirScan:
  goroutineCount: 3  # 设置"目录扫描"模块最大并发
vulnerabilityScan:
  goroutineCount: 2  # 设置"漏洞扫描"模块最大并发
'''

PLUGINSMODULES = [
    "TargetHandler",
    "SubdomainScan",
    "SubdomainSecurity",
    "AssetMapping",
    "PortScanPreparation",
    "PortScan",
    "PortFingerprint",
    "AssetHandle",
    "URLScan",
    "URLSecurity",
    "WebCrawler",
    "DirScan",
    "VulnerabilityScan",
    "PassiveScan"
]

PLUGINS = [
    {
        "module": "AssetHandle",
        "name": "WebFingerprint",
        "hash": "80718cc3fcb4827d942e6300184707e2",
        "parameter": "",
        "help": "无需参数",
        "introduction": "web指纹识别",
        "isSystem": True,
        "version": "1.0",
        "source": ""
    },
    {
        "module": "AssetMapping",
        "name": "httpx",
        "hash": "3a0d994a12305cb15a5cb7104d819623",
        "parameter": "-cdncheck true -screenshot false -tlsprobe false",
        "help": "-cdncheck 是否开启cdn检测 -screenshot 是否开启截图，默认关闭,开启需要安装chromium -tlsprobe 从tls信息发送http探测默认true",
        "introduction": "资产测绘",
        "isSystem": True,
        "version": "1.0",
        "source": ""
    },
    {
        "module": "DirScan",
        "name": "SentryDir",
        "hash": "920546788addc6d29ea63e4a314a1b85",
        "parameter": "-d {dict.dir.default} -t 10",
        "help": "-d 目录扫描字典 -t 扫描并发限制",
        "introduction": "目录扫描",
        "isSystem": True,
        "version": "1.0",
        "source": ""
    },
    {
        "module": "PortFingerprint",
        "name": "fingerprintx",
        "hash": "648a6f49eed57b1737ac702e02985b00",
        "parameter": "",
        "help": "无需参数",
        "introduction": "端口指纹识别",
        "isSystem": True,
        "version": "1.0",
        "source": ""
    },
    {
        "module": "PortScan",
        "name": "RustScan",
        "hash": "66b4ddeb983387df2b7ee7726653874d",
        "parameter": "-port {port.nmap top 1000} -b 600 -t 3000",
        "help": "-port 端口扫描范围 -b 端口扫描并发数量  -t 超时时间",
        "introduction": "端口存活扫描",
        "isSystem": True,
        "version": "1.0",
        "source": ""
    },
    {
        "module": "PortScanPreparation",
        "name": "SkipCdn",
        "hash": "9b91e0f18ac9043ec9fe250a39b4a2d9",
        "parameter": "",
        "help": "无需参数",
        "introduction": "检测是否cdn，跳过cdn的端口扫描",
        "isSystem": True,
        "version": "1.0",
        "source": ""
    },
    {
        "module": "SubdomainScan",
        "name": "subfinder",
        "hash": "d60ba73c70aac430a0a54e796e7e19b8",
        "parameter": "-t 10 -timeout 20 -max-time 10",
        "help": "-t 扫描线程  -timeout 超时时间 -max-time 最大等待时间",
        "introduction": "子域名扫描",
        "isSystem": True,
        "version": "1.0",
        "source": ""
    },
    {
        "module": "SubdomainScan",
        "name": "ksubdomain",
        "hash": "e8f55f5e0e9f4af1ca40eb19048b8c82",
        "parameter": "-subfile {dict.subdomain.default} -et 60",
        "help": "-subfile 子域名字典 -et 最长运行时间(分钟)",
        "introduction": "子域名爆破",
        "isSystem": True,
        "version": "1.0",
        "source": ""
    },
    {
        "module": "SubdomainSecurity",
        "name": "SubdomainTakeover",
        "hash": "c0c71c101271f38b8be1767f3626d291",
        "parameter": "",
        "help": "无需参数",
        "introduction": "子域名接管检测",
        "isSystem": True,
        "version": "1.0",
        "source": ""
    },
    {
        "module": "URLScan",
        "name": "wayback",
        "hash": "ef244b3462744dad3040f9dcf3194eb1",
        "parameter": "",
        "help": "无需参数",
        "introduction": "url扫描从Waybackarchive、Alienvault、Commoncrawl获取历史url",
        "isSystem": True,
        "version": "1.0",
        "source": ""
    },
    {
        "module": "URLScan",
        "name": "katana",
        "hash": "9669d0dcc52a5ca6dbbe580ffc99c364",
        "parameter": "-t 10 -timeout 5 -depth 5 -et 20",
        "help": "-t 并发数 -timeout 超时时间 -et 最长运行时间(分钟)",
        "introduction": "url爬取",
        "isSystem": True,
        "version": "1.0",
        "source": ""
    },
    {
        "module": "URLSecurity",
        "name": "sensitive",
        "hash": "2949994c04a4e124b9c98383489510f0",
        "parameter": "",
        "help": "无需参数",
        "introduction": "敏感信息泄露检测",
        "isSystem": True,
        "version": "1.0",
        "source": ""
    },
    {
        "module": "URLSecurity",
        "name": "PageMonitoring",
        "hash": "e52b8b16d49912ca564c22319c495403",
        "parameter": "",
        "help": "无需参数",
        "introduction": "页面监控，将所有url放入页面监控的计划任务中",
        "isSystem": True,
        "version": "1.0",
        "source": ""
    },
    {
        "module": "URLSecurity",
        "name": "trufflehog",
        "hash": "1aa212b9578dc3fb1409ee8de8ed005e",
        "parameter": "-pdf false -verify false",
        "help": "-pdf 开启pdf检测 -exclude 排除提取的规则(name1,name2) -verify 是否进行验证（验证通过再统计结果）",
        "introduction": "trufflehog密钥提取，如果设置了排除规则，需要重新安装才能重新启用被排除的规则",
        "isSystem": True,
        "version": "1.0",
        "source": ""
    },
    {
        "module": "VulnerabilityScan",
        "name": "nuclei",
        "hash": "ed93b8af6b72fe54a60efdb932cf6fbc",
        "parameter": "-s high,critical",
        "help": "参考官方支持t, s, es, tags, etags, rl, rld, bs, c, hbs, headc, jsc, pc, prc参数",
        "introduction": "漏洞扫描",
        "isSystem": True,
        "version": "1.0",
        "source": ""
    },
    {
        "module": "WebCrawler",
        "name": "rad",
        "hash": "4b292861d3228af0e4da8e7ef979497c",
        "parameter": "",
        "help": "无需参数",
        "introduction": "爬虫",
        "isSystem": True,
        "version": "1.0",
        "source": ""
    }
]

SCANTEMPLATE = {
  "TargetHandler": [],
  "Parameters": {
    "TargetHandler": {},
    "SubdomainScan": {
      "d60ba73c70aac430a0a54e796e7e19b8": "-t 10 -timeout 20 -max-time 10",
      "e8f55f5e0e9f4af1ca40eb19048b8c82": "-subfile {dict.subdomain.default} -et 60"
    },
    "SubdomainSecurity": {},
    "PortScanPreparation": {},
    "PortScan": {
      "66b4ddeb983387df2b7ee7726653874d": "-port {port.nmap top 1000} -b 600 -t 3000"
    },
    "PortFingerprint": {},
    "AssetMapping": {
      "3a0d994a12305cb15a5cb7104d819623": "-cdncheck true -screenshot false"
    },
    "AssetHandle": {},
    "URLScan": {
      "9669d0dcc52a5ca6dbbe580ffc99c364": "-t 10 -timeout 5 -depth 5 -et 20"
    },
    "WebCrawler": {},
    "URLSecurity": {},
    "DirScan": {
      "920546788addc6d29ea63e4a314a1b85": "-d {dict.dir.default} -t 10"
    },
    "VulnerabilityScan": {
      "ed93b8af6b72fe54a60efdb932cf6fbc": "-s high,critical"
    }
  },
  "SubdomainScan": [
    "d60ba73c70aac430a0a54e796e7e19b8",
    "e8f55f5e0e9f4af1ca40eb19048b8c82"
  ],
  "SubdomainSecurity": [
    "c0c71c101271f38b8be1767f3626d291"
  ],
  "PortScanPreparation": [],
  "PortScan": [
    "66b4ddeb983387df2b7ee7726653874d"
  ],
  "PortFingerprint": [
    "648a6f49eed57b1737ac702e02985b00"
  ],
  "AssetMapping": [
    "3a0d994a12305cb15a5cb7104d819623"
  ],
  "AssetHandle": [
    "80718cc3fcb4827d942e6300184707e2"
  ],
  "URLScan": [
    "ef244b3462744dad3040f9dcf3194eb1",
    "9669d0dcc52a5ca6dbbe580ffc99c364"
  ],
  "WebCrawler": [],
  "URLSecurity": [
    "1aa212b9578dc3fb1409ee8de8ed005e"
  ],
  "DirScan": [
  ],
  "VulnerabilityScan": [],
  "name": "default",
  "vullist": []
}

FIELD = {
    "asset": ["time", "lastScanTime", "host", "ip", "port", "service", "tls", "transport", "version", "metadata", "project", "type", "tags", "taskName", "rootDomain", "urlPath", "hash", "cdnname", "url", "title", "error", "body", "screenshot", "faviconmmh3", "faviconpath", "rawheaders", "jarm", "technologies", "statuscode", "contentlength", "cdn", "webcheck", "iconcontent", "domain", "webServer"],
    "crawler": ['url', 'method', 'body', 'project', 'taskName', 'resultId', 'rootDomain', 'time', 'tags'],
    "DirScanResult": ['url', 'status', 'msg', 'project', 'length', 'taskName', 'rootDomain', 'tags'],
    "SensitiveResult": ['url', 'urlid', 'sid', 'match', 'project', 'color', 'time', 'md5', 'taskName', 'rootDomain', 'tags', 'status'],
    "subdomain": ['host', 'type', 'value', 'ip', 'time', 'tags', 'project', 'taskName', 'rootDomain'],
    "SubdoaminTakerResult": ['input', 'value', 'cname', 'response', 'project', 'taskName', 'rootDomain', 'tags'],
    "UrlScan": ['input', 'source', 'outputtype', 'output', 'status', 'length', 'time', 'body', 'project', 'taskName', 'resultId', 'rootDomain', 'tags'],
    "vulnerability": ['url', 'vulnid', 'vulname', 'matched', 'project', 'level', 'time', 'request', 'response', 'taskName', 'rootdomain', 'tags', 'status'],
    "httpAsset": ["time", "lastScanTime", "tls", "hash", "cdnname", "port", "url", "title", "type", "error", "body", "host", "ip", "screenshot", "faviconmmh3", "faviconpath", "rawheaders", "jarm", "technologies", "statuscode", "contentlength", "cdn", "webcheck", "project", "iconcontent", "domain", "taskName", "webServer", "service", "rootDomain", "tags"],
    "otherAsset": ["time", "lastScanTime", "host", "ip", "port", "service", "tls", "transport", "version", "metadata", "project", "type", "tags", "taskName", "rootDomain", "urlPath"]
}