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

sensitiveList = [{'name': 'JSON Web Token',
                  'regular': '(eyJ[A-Za-z0-9_-]{10,}\\.[A-Za-z0-9._-]{10,}|eyJ[A-Za-z0-9_\\/+-]{10,}\\.[A-Za-z0-9._\\/+-]{10,})',
                  'color': 'green', 'state': True}, {'name': 'Swagger UI',
                                                     'regular': '((swagger-ui.html)|(\\"swagger\\":)|(Swagger UI)|(swaggerUi)|(swaggerVersion))',
                                                     'color': 'red', 'state': True},
                 {'name': 'Ueditor', 'regular': '(ueditor\\.(config|all)\\.js)', 'color': 'green', 'state': True},
                 {'name': 'Java Deserialization', 'regular': '(javax\\.faces\\.ViewState)', 'color': 'yellow',
                  'state': True},
                 {'name': 'URL As A Value', 'regular': '(=(https?)(://|%3a%2f%2f))', 'color': 'cyan', 'state': True},
                 {'name': 'Upload Form', 'regular': '(type=\\"file\\")', 'color': 'yellow', 'state': True},
                 {'name': 'Email',
                  'regular': '(([a-z0-9][_|\\.])*[a-z0-9]+@([a-z0-9][-|_|\\.])*[a-z0-9]+\\.((?!js|css|jpg|jpeg|png|ico)[a-z]{2,}))',
                  'color': 'yellow', 'state': True}, {'name': 'Chinese IDCard',
                                                      'regular': "'[^0-9]((\\d{8}(0\\d|10|11|12)([0-2]\\d|30|31)\\d{3}$)|(\\d{6}(18|19|20)\\d{2}(0[1-9]|10|11|12)([0-2]\\d|30|31)\\d{3}(\\d|X|x)))[^0-9]'",
                                                      'color': 'orange', 'state': True},
                 {'name': 'Chinese Mobile Number',
                  'regular': "'[^\\w]((?:(?:\\+|00)86)?1(?:(?:3[\\d])|(?:4[5-79])|(?:5[0-35-9])|(?:6[5-7])|(?:7[0-8])|(?:8[\\d])|(?:9[189]))\\d{8})[^\\w]'",
                  'color': 'orange', 'state': True}, {'name': 'Internal IP Address',
                                                      'regular': "'[^0-9]((127\\.0\\.0\\.1)|(10\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3})|(172\\.((1[6-9])|(2\\d)|(3[01]))\\.\\d{1,3}\\.\\d{1,3})|(192\\.168\\.\\d{1,3}\\.\\d{1,3}))'",
                                                      'color': 'cyan', 'state': True}, {'name': 'MAC Address',
                                                                                        'regular': '(^([a-fA-F0-9]{2}(:[a-fA-F0-9]{2}){5})|[^a-zA-Z0-9]([a-fA-F0-9]{2}(:[a-fA-F0-9]{2}){5}))',
                                                                                        'color': 'green',
                                                                                        'state': True},
                 {'name': 'Chinese Bank Card ID', 'regular': "'[^0-9]([1-9]\\d{12,18})[^0-9]'", 'color': 'orange',
                  'state': True},
                 {'name': 'Cloud Key', 'regular': '((accesskeyid)|(accesskeysecret)|(LTAI[a-z0-9]{12,20}))',
                  'color': 'yellow', 'state': True}, {'name': 'Windows File/Dir Path',
                                                      'regular': "'[^\\w](([a-zA-Z]:\\\\(?:\\w+\\\\?)*)|([a-zA-Z]:\\\\(?:\\w+\\\\)*\\w+\\.\\w+))'",
                                                      'color': 'green', 'state': True}, {'name': 'Password Field',
                                                                                         'regular': '((|\'|")([p](ass|wd|asswd|assword))(|\'|")(:|=)( |)(\'|")(.*?)(\'|")(|,))',
                                                                                         'color': 'yellow',
                                                                                         'state': True},
                 {'name': 'Username Field',
                  'regular': '((|\'|")(([u](ser|name|ame|sername))|(account))(|\'|")(:|=)( |)(\'|")(.*?)(\'|")(|,))',
                  'color': 'green', 'state': True},
                 {'name': 'WeCom Key', 'regular': '([c|C]or[p|P]id|[c|C]orp[s|S]ecret)', 'color': 'green',
                  'state': True},
                 {'name': 'JDBC Connection', 'regular': '(jdbc:[a-z:]+://[a-z0-9\\.\\-_:;=/@?,&]+)', 'color': 'yellow',
                  'state': True}, {'name': 'Authorization Header',
                                   'regular': '((basic [a-z0-9=:_\\+\\/-]{5,100})|(bearer [a-z0-9_.=:_\\+\\/-]{5,100}))',
                                   'color': 'yellow', 'state': True},
                 {'name': 'Github Access Token', 'regular': '([a-z0-9_-]*:[a-z0-9_\\-]+@github\\.com*)',
                  'color': 'green', 'state': True}, {'name': 'Sensitive Field',
                                                     'regular': '((|\'|")([\\w]{0,10})((key)|(secret)|(token)|(config)|(auth)|(access)|(admin))(|\'|")(:|=)(       |)(\'|")(.*?)(\'|")(|,))',
                                                     'color': 'yellow', 'state': True}, {'name': 'Linkfinder',
                                                                                         'regular': '(?:"|\')(((?:[a-zA-Z]{1,10}://|//)[^"\'/]{1,}\\.[a-zA-Z]{2,}[^"\']{0,})|((?:/|\\.\\./|\\./)[^"\'><,;|*()(%%$^/\\\\\\[\\]][^"\'><,;|()]{1,})|([a-zA-Z0-9_\\-/]{1,}/[a-zA-Z0-9_\\-/]{1,}\\.(?:[a-zA-Z]{1,4}|action)(?:[\\?|#][^"|\']{0,}|))|([a-zA-Z0-9_\\-/]{1,}/[a-zA-Z0-9_\\-/]{3,}(?:[\\?|#][^"|\']{0,}|))|([a-zA-Z0-9_\\-]{1,}\\.(?:\\w)(?:[\\?|#][^"|\']{0,}|)))(?:"|\')',
                                                                                         'color': 'gray',
                                                                                         'state': True},
                 {'name': 'Source Map', 'regular': '(\\.js\\.map)', 'color': 'null', 'state': True},
                 {'name': 'HTML Notes', 'regular': '(<!--[\\s\\S]*?-->)', 'color': 'green', 'state': True},
                 {'name': 'Create Script', 'regular': '(createElement\\(\\"script\\"\\))', 'color': 'green',
                  'state': True}, {'name': 'URL Schemes',
                                   'regular': '(?![http]|[https])(([-A-Za-z0-9]{1,20})://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|])',
                                   'color': 'yellow', 'state': True},
                 {'name': 'Potential cryptographic private key', 'regular': '(\\.pem[\'"])', 'color': 'green',
                  'state': True},
                 {'name': 'google_api', 'regular': '(AIza[0-9A-Za-z-_]{35})', 'color': 'red', 'state': True},
                 {'name': 'firebase', 'regular': '(AAAA[A-Za-z0-9_-]{7}:[A-Za-z0-9_-]{140})', 'color': 'red',
                  'state': True},
                 {'name': 'authorization_api', 'regular': '(api[key|_key|\\s+]+[a-zA-Z0-9_\\-]{5,100})', 'color': 'red',
                  'state': True}, {'name': 'Log file', 'regular': '(\\.log[\'"])', 'color': 'green', 'state': True},
                 {'name': 'Potential cryptographic key bundle', 'regular': '(\\.pkcs12[\'"])', 'color': 'yellow',
                  'state': True},
                 {'name': 'Potential cryptographic key bundle', 'regular': '(\\.p12[\'"])', 'color': 'yellow',
                  'state': True},
                 {'name': 'Potential cryptographic key bundle', 'regular': '(\\.pfx[\'"])', 'color': 'yellow',
                  'state': True},
                 {'name': 'Pidgin OTR private key', 'regular': '(otr\\.private_key)', 'color': 'yellow', 'state': True},
                 {'name': 'File',
                  'regular': '(\\.((asc)|(ovpn)|(cscfg)|(rdp)|(mdf)|(sdf)|(sqlite)|(sqlite3)|(bek)|(tpm)|(fve)|(jks)|(psafe3)|(agilekeychain)|(keychain)|(pcap)|(gnucash)|(kwallet)|(tblk)|(dayone)|(exports)|(functions)|(extra)|(proftpdpasswd))[\'"])',
                  'color': 'yellow', 'state': True},
                 {'name': 'Ruby On Rails secret token configuration file', 'regular': '(secret_token\\.rb)',
                  'color': 'yellow', 'state': True},
                 {'name': 'Carrierwave configuration file', 'regular': '(carrierwave\\.rb)', 'color': 'yellow',
                  'state': True},
                 {'name': 'Potential Ruby On Rails database configuration file', 'regular': '(database\\.yml)',
                  'color': 'yellow', 'state': True},
                 {'name': 'OmniAuth configuration file', 'regular': '(omniauth\\.rb)', 'color': 'yellow',
                  'state': True},
                 {'name': 'Django configuration file', 'regular': '(settings\\.py)', 'color': 'yellow', 'state': True},
                 {'name': 'Jenkins publish over SSH plugin file',
                  'regular': '(jenkins.plugins.publish_over_ssh\\.BapSshPublisherPlugin.xml)', 'color': 'yellow',
                  'state': True},
                 {'name': 'Potential Jenkins credentials file', 'regular': '(credentials\\.xml)', 'color': 'yellow',
                  'state': True},
                 {'name': 'Potential MediaWiki configuration file', 'regular': 'LocalSettings\\.php', 'color': 'yellow',
                  'state': True},
                 {'name': 'Sequel Pro MySQL database manager bookmark file', 'regular': '(Favorites\\.plist)',
                  'color': 'yellow', 'state': True},
                 {'name': 'Little Snitch firewall configuration file', 'regular': '(configuration\\.user\\.xpl)',
                  'color': 'yellow', 'state': True},
                 {'name': 'Potential jrnl journal file', 'regular': '(journal\\.txt)', 'color': 'yellow',
                  'state': True},
                 {'name': 'Chef Knife configuration file', 'regular': '(knife\\.rb)', 'color': 'yellow', 'state': True},
                 {'name': 'Robomongo MongoDB manager configuration file', 'regular': '(robomongo\\.json)',
                  'color': 'yellow', 'state': True},
                 {'name': 'FileZilla FTP configuration file', 'regular': '(filezilla\\.xml)', 'color': 'yellow',
                  'state': True},
                 {'name': 'FileZilla FTP recent servers file', 'regular': '(recentservers\\.xml)', 'color': 'yellow',
                  'state': True},
                 {'name': 'Ventrilo server configuration file', 'regular': '(ventrilo_srv\\.ini)', 'color': 'yellow',
                  'state': True},
                 {'name': 'Terraform variable config file', 'regular': '(terraform\\.tfvars)', 'color': 'yellow',
                  'state': True}, {'name': 'Private SSH key', 'regular': '(.*_rsa)', 'color': 'yellow', 'state': True},
                 {'name': 'Private SSH key', 'regular': '(.*_dsa)', 'color': 'yellow', 'state': True},
                 {'name': 'Private SSH key', 'regular': '(.*_ed25519)', 'color': 'yellow', 'state': True},
                 {'name': 'Private SSH key', 'regular': '(.*_ecdsa)', 'color': 'yellow', 'state': True},
                 {'name': 'SSH configuration file', 'regular': '(\\.ssh_config)', 'color': 'yellow', 'state': True},
                 {'name': 'Shell command history file', 'regular': '(\\.?(bash_|zsh_|sh_|z)?history)',
                  'color': 'yellow', 'state': True},
                 {'name': 'MySQL client command history file', 'regular': '(.?mysql_history)', 'color': 'yellow',
                  'state': True},
                 {'name': 'PostgreSQL client command history file', 'regular': '(\\.?psql_history)', 'color': 'yellow',
                  'state': True},
                 {'name': 'PostgreSQL password file', 'regular': '(\\.?pgpass)', 'color': 'yellow', 'state': True},
                 {'name': 'Ruby IRB console history file', 'regular': '(\\.?irb_history)', 'color': 'yellow',
                  'state': True},
                 {'name': 'Pidgin chat client account configuration file', 'regular': '(\\.?purple/accounts\\\\.xml)',
                  'color': 'yellow', 'state': True}, {'name': 'DBeaver SQL database manager configuration file',
                                                      'regular': '(\\.?dbeaver-data-sources.xml)', 'color': 'yellow',
                                                      'state': True},
                 {'name': 'Mutt e-mail client configuration file', 'regular': '(\\.?muttrc)', 'color': 'yellow',
                  'state': True},
                 {'name': 'S3cmd configuration file', 'regular': '(\\.?s3cfg)', 'color': 'yellow', 'state': True},
                 {'name': 'AWS CLI credentials file', 'regular': '(\\.?aws/credentials)', 'color': 'yellow',
                  'state': True},
                 {'name': 'SFTP connection configuration file', 'regular': '(sftp-config(\\.json)?)', 'color': 'yellow',
                  'state': True},
                 {'name': 'T command-line Twitter client configuration file', 'regular': '(\\.?trc)', 'color': 'yellow',
                  'state': True},
                 {'name': 'Shell configuration file', 'regular': '(\\.?(bash|zsh|csh)rc)', 'color': 'yellow',
                  'state': True}, {'name': 'Shell profile configuration file', 'regular': '(\\.?(bash_|zsh_)?profile)',
                                   'color': 'yellow', 'state': True},
                 {'name': 'Shell command alias configuration file', 'regular': '(\\.?(bash_|zsh_)?aliases)',
                  'color': 'yellow', 'state': True},
                 {'name': 'PHP configuration file', 'regular': '(config(\\.inc)?\\.php)', 'color': 'yellow',
                  'state': True},
                 {'name': 'GNOME Keyring database file', 'regular': '(key(store|ring))', 'color': 'yellow',
                  'state': True},
                 {'name': 'KeePass password manager database file', 'regular': '(kdbx?)', 'color': 'yellow',
                  'state': True},
                 {'name': 'SQL dump file', 'regular': '(sql(dump)?)', 'color': 'yellow', 'state': True},
                 {'name': 'Apache htpasswd file', 'regular': '(\\.?htpasswd)', 'color': 'yellow', 'state': True},
                 {'name': 'Configuration file for auto-login process', 'regular': '((\\.|_)?netrc)', 'color': 'yellow',
                  'state': True},
                 {'name': 'Rubygems credentials file', 'regular': '(\\.?gem/credentials)', 'color': 'yellow',
                  'state': True},
                 {'name': 'Tugboat DigitalOcean management tool configuration', 'regular': '(\\.?tugboat)',
                  'color': 'yellow', 'state': True},
                 {'name': 'DigitalOcean doctl command-line client configuration file', 'regular': '(doctl/config.yaml)',
                  'color': 'yellow', 'state': True},
                 {'name': 'git-credential-store helper credentials file', 'regular': '(\\.?git-credentials)',
                  'color': 'yellow', 'state': True},
                 {'name': 'GitHub Hub command-line client configuration file', 'regular': '(config/hub)',
                  'color': 'yellow', 'state': True},
                 {'name': 'Git configuration file', 'regular': '(\\.?gitconfig)', 'color': 'yellow', 'state': True},
                 {'name': 'Chef private key', 'regular': '(\\.?chef/(.*)\\\\.pem)', 'color': 'yellow', 'state': True},
                 {'name': 'Potential Linux shadow file', 'regular': '(etc/shadow)', 'color': 'yellow', 'state': True},
                 {'name': 'Potential Linux passwd file', 'regular': '(etc/passwd)', 'color': 'yellow', 'state': True},
                 {'name': 'Docker configuration file', 'regular': '(\\.?dockercfg)', 'color': 'yellow', 'state': True},
                 {'name': 'NPM configuration file', 'regular': '(\\.?npmrc)', 'color': 'yellow', 'state': True},
                 {'name': 'Environment configuration file', 'regular': '(\\.?env)', 'color': 'yellow', 'state': True},
                 {'name': 'AWS Access Key ID Value',
                  'regular': '((A3T[A-Z0-9]|AKIA|AGPA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16})', 'color': 'red',
                  'state': True}, {'name': 'ak sk',
                                   'regular': '(((access_key|access_token|admin_pass|admin_user|algolia_admin_key|algolia_api_key|alias_pass|alicloud_access_key|amazon_secret_access_key|amazonaws|ansible_vault_password|aos_key|api_key|api_key_secret|api_key_sid|api_secret|api.googlemaps AIza|apidocs|apikey|apiSecret|app_debug|app_id|app_key|app_log_level|app_secret|appkey|appkeysecret|application_key|appsecret|appspot|auth_token|authorizationToken|authsecret|aws_access|aws_access_key_id|aws_bucket|aws_key|aws_secret|aws_secret_key|aws_token|AWSSecretKey|b2_app_key|bashrc password|bintray_apikey|bintray_gpg_password|bintray_key|bintraykey|bluemix_api_key|bluemix_pass|browserstack_access_key|bucket_password|bucketeer_aws_access_key_id|bucketeer_aws_secret_access_key|built_branch_deploy_key|bx_password|cache_driver|cache_s3_secret_key|cattle_access_key|cattle_secret_key|certificate_password|ci_deploy_password|client_secret|client_zpk_secret_key|clojars_password|cloud_api_key|cloud_watch_aws_access_key|cloudant_password|cloudflare_api_key|cloudflare_auth_key|cloudinary_api_secret|cloudinary_name|codecov_token|config|conn.login|connectionstring|consumer_key|consumer_secret|credentials|cypress_record_key|database_password|database_schema_test|datadog_api_key|datadog_app_key|db_password|db_server|db_username|dbpasswd|dbpassword|dbuser|deploy_password|digitalocean_ssh_key_body|digitalocean_ssh_key_ids|docker_hub_password|docker_key|docker_pass|docker_passwd|docker_password|dockerhub_password|dockerhubpassword|dot-files|dotfiles|droplet_travis_password|dynamoaccesskeyid|dynamosecretaccesskey|elastica_host|elastica_port|elasticsearch_password|encryption_key|encryption_password|env.heroku_api_key|env.sonatype_password|eureka.awssecretkey)[a-z0-9_ .\\-,]{0,25})(=|>|:=|\\|\\|:|<=|=>|:).{0,5}[\'\\"]([0-9a-zA-Z\\-_=]{8,64}))\\b`',
                                   'color': 'red', 'state': True}, {'name': 'AWS Access Key ID',
                                                                    'regular': '(("|\'|`)?((?i)aws)?_?((?i)access)_?((?i)key)?_?((?i)id)?("|\'|`)?\\s{0,50}(:|=>|=)\\s{0,50}("|\'|`)?(A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16}("|\'|`)?)',
                                                                    'color': 'red', 'state': True},
                 {'name': 'AWS Account ID',
                  'regular': '(("|\'|`)?((?i)aws)?_?((?i)account)_?((?i)id)?("|\'|`)?\\s{0,50}(:|=>|=)\\s{0,50}("|\'|`)?[0-9]{4}-?[0-9]{4}-?[0-9]{4}("|\'|`)?)',
                  'color': 'red', 'state': True},
                 {'name': 'Artifactory API Token', 'regular': '((?:\\s|=|:|"|^)AKC[a-zA-Z0-9]{10,})', 'color': 'red',
                  'state': True},
                 {'name': 'Artifactory Password', 'regular': '((?:\\s|=|:|"|^)AP[\\dABCDEF][a-zA-Z0-9]{8,})',
                  'color': 'red', 'state': True},
                 {'name': 'Authorization Basic', 'regular': '(basic [a-zA-Z0-9_\\\\-:\\\\.=]+)', 'color': 'red',
                  'state': True},
                 {'name': 'Authorization Authorization Bearer', 'regular': '(bearer [a-zA-Z0-9_\\\\-\\\\.=]+)',
                  'color': 'red', 'state': True}, {'name': 'AWS Client ID',
                                                   'regular': '((A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16})',
                                                   'color': 'red', 'state': True}, {'name': 'AWS MWS Key',
                                                                                    'regular': '(amzn\\.mws\\.[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})',
                                                                                    'color': 'red', 'state': True},
                 {'name': 'AWS MWS Key',
                  'regular': '(amzn\\.mws\\.[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})',
                  'color': 'red', 'state': True},
                 {'name': 'AWS Secret Key', 'regular': '((?i)aws(.{0,20})?(?-i)[\'\\"][0-9a-zA-Z\\/+]{40}[\'"])',
                  'color': 'red', 'state': True}, {'name': 'Base32',
                                                   'regular': '((?:[A-Z2-7]{8})*(?:[A-Z2-7]{2}={6}|[A-Z2-7]{4}={4}|[A-Z2-7]{5}={3}|[A-Z2-7]{7}=)?)',
                                                   'color': 'null', 'state': True},
                 {'name': 'Base64', 'regular': '((eyJ|YTo|Tzo|PD[89]|aHR0cHM6L|aHR0cDo|rO0)[a-zA-Z0-9+/]+={0,2})',
                  'color': 'null', 'state': True}, {'name': 'Basic Auth Credentials',
                                                    'regular': '((?<=:\\/\\/)[a-zA-Z0-9]+:[a-zA-Z0-9]+@[a-zA-Z0-9]+\\.[a-zA-Z]+)',
                                                    'color': 'red', 'state': True},
                 {'name': 'Cloudinary Basic Auth', 'regular': '(cloudinary:\\/\\/[0-9]{15}:[0-9A-Za-z]+@[a-z]+)',
                  'color': 'red', 'state': True},
                 {'name': 'Facebook Access Token', 'regular': '(EAACEdEose0cBA[0-9A-Za-z]+)', 'color': 'red',
                  'state': True},
                 {'name': 'Facebook Client ID', 'regular': '((?i)(facebook|fb)(.{0,20})?[\'\\"][0-9]{13,17})',
                  'color': 'red', 'state': True}, {'name': 'Facebook Oauth',
                                                   'regular': '([f|F][a|A][c|C][e|E][b|B][o|O][o|O][k|K].*[\'|\\"][0-9a-f]{32}[\'|\\"])',
                                                   'color': 'red', 'state': True},
                 {'name': 'Facebook Secret Key', 'regular': '((?i)(facebook|fb)(.{0,20})?(?-i)[\'\\"][0-9a-f]{32})',
                  'color': 'red', 'state': True},
                 {'name': 'Github', 'regular': '((?i)github(.{0,20})?(?-i)[\'\\"][0-9a-zA-Z]{35,40})', 'color': 'red',
                  'state': True},
                 {'name': 'Google API Key', 'regular': '(AIza[0-9A-Za-z\\\\-_]{35})', 'color': 'red', 'state': True},
                 {'name': 'Google Cloud Platform API Key',
                  'regular': '((?i)(google|gcp|youtube|drive|yt)(.{0,20})?[\'\\"][AIza[0-9a-z\\\\-_]{35}][\'\\"])',
                  'color': 'red', 'state': True},
                 {'name': 'Google Oauth', 'regular': '([0-9]+-[0-9A-Za-z_]{32}\\.apps\\.googleusercontent\\.com)',
                  'color': 'red', 'state': True}, {'name': 'Heroku API Key',
                                                   'regular': '([h|H][e|E][r|R][o|O][k|K][u|U].{0,30}[0-9A-F]{8}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{12})',
                                                   'color': 'red', 'state': True},
                 {'name': 'LinkedIn Secret Key', 'regular': '((?i)linkedin(.{0,20})?[\'\\"][0-9a-z]{16}[\'\\"])',
                  'color': 'red', 'state': True},
                 {'name': 'Mailchamp API Key', 'regular': '([0-9a-f]{32}-us[0-9]{1,2})', 'color': 'red', 'state': True},
                 {'name': 'Mailgun API Key', 'regular': '(key-[0-9a-zA-Z]{32})', 'color': 'red', 'state': True},
                 {'name': 'Picatic API Key', 'regular': '(sk_live_[0-9a-z]{32})', 'color': 'red', 'state': True},
                 {'name': 'Slack Token', 'regular': '(xox[baprs]-([0-9a-zA-Z]{10,48})?)', 'color': 'red',
                  'state': True}, {'name': 'Slack Webhook',
                                   'regular': '(https://hooks.slack.com/services/T[a-zA-Z0-9_]{8}/B[a-zA-Z0-9_]{8}/[a-zA-Z0-9_]{24})',
                                   'color': 'red', 'state': True},
                 {'name': 'Stripe API Key', 'regular': '((?:r|s)k_live_[0-9a-zA-Z]{24})', 'color': 'red',
                  'state': True},
                 {'name': 'Square Access Token', 'regular': '(sqOatp-[0-9A-Za-z\\\\-_]{22})', 'color': 'red',
                  'state': True},
                 {'name': 'Square Oauth Secret', 'regular': '(sq0csp-[ 0-9A-Za-z\\\\-_]{43})', 'color': 'red',
                  'state': True},
                 {'name': 'Twilio API Key', 'regular': '(SK[0-9a-fA-F]{32})', 'color': 'red', 'state': True},
                 {'name': 'Twitter Oauth',
                  'regular': '([t|T][w|W][i|I][t|T][t|T][e|E][r|R].{0,30}[\'\\"\\\\s][0-9a-zA-Z]{35,44}[\'\\"\\\\s])',
                  'color': 'red', 'state': True},
                 {'name': 'Twitter Secret Key', 'regular': '(?i)twitter(.{0,20})?[\'\\"][0-9a-z]{35,44}',
                  'color': 'red', 'state': True},
                 {'name': 'google_captcha', 'regular': '(6L[0-9A-Za-z-_]{38}|^6[0-9a-zA-Z_-]{39})', 'color': 'red',
                  'state': True},
                 {'name': 'google_oauth', 'regular': '(ya29\\.[0-9A-Za-z\\-_]+)', 'color': 'red', 'state': True},
                 {'name': 'amazon_aws_access_key_id', 'regular': '(A[SK]IA[0-9A-Z]{16})', 'color': 'red',
                  'state': True},
                 {'name': 'amazon_aws_url', 'regular': 's3\\.amazonaws.com[/]+|[a-zA-Z0-9_-]*\\.s3\\.amazonaws.com',
                  'color': 'red', 'state': True},
                 {'name': 'authorization_api', 'regular': '(api[key|\\s*]+[a-zA-Z0-9_\\-]+)', 'color': 'red',
                  'state': True},
                 {'name': 'twilio_account_sid', 'regular': '(AC[a-zA-Z0-9_\\-]{32})', 'color': 'red', 'state': True},
                 {'name': 'twilio_app_sid', 'regular': '(AP[a-zA-Z0-9_\\-]{32})', 'color': 'red', 'state': True},
                 {'name': 'paypal_braintree_access_token',
                  'regular': '(access_token\\$production\\$[0-9a-z]{16}\\$[0-9a-f]{32})', 'color': 'red',
                  'state': True}, {'name': 'square_oauth_secret',
                                   'regular': '(sq0csp-[ 0-9A-Za-z\\-_]{43}|sq0[a-z]{3}-[0-9A-Za-z\\-_]{22,43})',
                                   'color': 'red', 'state': True},
                 {'name': 'square_access_token', 'regular': '(sqOatp-[0-9A-Za-z\\-_]{22}|EAAA[a-zA-Z0-9]{60})',
                  'color': 'red', 'state': True},
                 {'name': 'rsa_private_key', 'regular': '(-----BEGIN RSA PRIVATE KEY-----)', 'color': 'red',
                  'state': True},
                 {'name': 'ssh_dsa_private_key', 'regular': '(-----BEGIN DSA PRIVATE KEY-----)', 'color': 'red',
                  'state': True},
                 {'name': 'ssh_dc_private_key', 'regular': '(-----BEGIN EC PRIVATE KEY-----)', 'color': 'red',
                  'state': True},
                 {'name': 'pgp_private_block', 'regular': '(-----BEGIN PGP PRIVATE KEY BLOCK-----)', 'color': 'red',
                  'state': True},
                 {'name': 'json_web_token', 'regular': '(ey[A-Za-z0-9-_=]+\\.[A-Za-z0-9-_=]+\\.?[A-Za-z0-9-_.+/=]*)',
                  'color': 'red', 'state': True},
                 {'name': 'Google Cloud', 'regular': '(GOOG[\\w\\W]{10,30})', 'color': 'red', 'state': True},
                 {'name': 'Microsoft Azure', 'regular': '(AZ[A-Za-z0-9]{34,40})', 'color': 'red', 'state': True},
                 {'name': '腾讯云', 'regular': '(AKID[A-Za-z0-9]{13,20})', 'color': 'red', 'state': True},
                 {'name': '亚马逊云', 'regular': '(AKIA[A-Za-z0-9]{16})', 'color': 'red', 'state': True},
                 {'name': 'IBM Cloud', 'regular': '(IBM[A-Za-z0-9]{10,40})', 'color': 'red', 'state': True},
                 {'name': 'Oracle Cloud', 'regular': '(OCID[A-Za-z0-9]{10,40})', 'color': 'red', 'state': True},
                 {'name': '阿里云', 'regular': '(LTAI[A-Za-z0-9]{12,20})', 'color': 'red', 'state': True},
                 {'name': '华为云', 'regular': '(AK[\\w\\W]{10,62})', 'color': 'red', 'state': True},
                 {'name': '百度云', 'regular': '(AK[A-Za-z0-9]{10,40})', 'color': 'red', 'state': True},
                 {'name': '京东云', 'regular': '(AK[A-Za-z0-9]{10,40})', 'color': 'red', 'state': True},
                 {'name': 'UCloud', 'regular': '(UC[A-Za-z0-9]{10,40})', 'color': 'red', 'state': True},
                 {'name': '青云', 'regular': '(QY[A-Za-z0-9]{10,40})', 'color': 'red', 'state': True},
                 {'name': '金山云', 'regular': '(KS3[A-Za-z0-9]{10,40})', 'color': 'red', 'state': True},
                 {'name': '联通云', 'regular': '(LTC[A-Za-z0-9]{10,60})', 'color': 'red', 'state': True},
                 {'name': '移动云', 'regular': '(YD[A-Za-z0-9]{10,60})', 'color': 'red', 'state': True},
                 {'name': '电信云', 'regular': '(CTC[A-Za-z0-9]{10,60})', 'color': 'red', 'state': True},
                 {'name': '一云通', 'regular': '(YYT[A-Za-z0-9]{10,60})', 'color': 'red', 'state': True},
                 {'name': '用友云', 'regular': '(YY[A-Za-z0-9]{10,40})', 'color': 'red', 'state': True},
                 {'name': '南大通用云', 'regular': '(CI[A-Za-z0-9]{10,40})', 'color': 'red', 'state': True},
                 {'name': 'G-Core Labs', 'regular': '(gcore[A-Za-z0-9]{10,30})', 'color': 'red', 'state': True},
                 {'name': 'MailChimp API Key', 'regular': '([0-9a-f]{32}-us[0-9]{12})', 'color': 'red', 'state': True},
                 {'name': 'Outlook team', 'regular': '((https://outlook\\.office.com/webhook/[0-9a-f-]{36}@))',
                  'color': 'red', 'state': True},
                 {'name': 'Sauce Token', 'regular': '(?i)sauce.{0,50}("|\'|`)?[0-9a-f-]{36}("|\'|`)?', 'color': 'red',
                  'state': True},
                 {'name': 'SonarQube Docs API Key', 'regular': '((?i)sonar.{0,50}("|\'|`)?[0-9a-f]{40}("|\'|`)?)',
                  'color': 'red', 'state': True},
                 {'name': 'HockeyApp', 'regular': '(?i)hockey.{0,50}("|\'|`)?[0-9a-f]{32}("|\'|`)?', 'color': 'red',
                  'state': True}, {'name': 'Username and password in URI',
                                   'regular': '(([\\w+]{1,24})(://)([^$<]{1})([^\\s";]{1,}):([^$<]{1})([^\\s";/]{1,})@[-a-zA-Z0-9@:%._\\\\+~#=]{1,256}\\.[a-zA-Z0-9()]{1,24}([^\\s]+))',
                                   'color': 'red', 'state': True},
                 {'name': 'NuGet API Key', 'regular': '(oy2[a-z0-9]{43})', 'color': 'red', 'state': True},
                 {'name': 'StackHawk API Key', 'regular': '(hawk\\.[0-9A-Za-z\\-_]{20}\\.[0-9A-Za-z\\-_]{20})',
                  'color': 'red', 'state': True},
                 {'name': 'Heroku config file', 'regular': '(heroku\\.json)', 'color': 'yellow', 'state': True},
                 {'name': 'jwt_token',
                  'regular': 'eyJ[A-Za-z0-9_\\/+-]{10,}={0,2}\\.[A-Za-z0-9_\\/+\\-]{15,}={0,2}\\\\.[A-Za-z0-9_\\/+\\-]{10,}={0,2}',
                  'color': 'yellow', 'state': True}, {'name': 'INFO-KEY',
                                                      'regular': '(access_key|access_token|admin_pass|admin_user|algolia_admin_key|algolia_api_key|alias_pass|alicloud_access_key|amazon_secret_access_key|amazonaws|ansible_vault_password|aos_key|api_key|api_key_secret|api_key_sid|api_secret|api.googlemaps AIza|apidocs|apikey|apiSecret|app_debug|app_id|app_key|app_log_level|app_secret|appkey|appkeysecret|application_key|appsecret|appspot|auth_token|authorizationToken|authsecret|aws_access|aws_access_key_id|aws_bucket|aws_key|aws_secret|aws_secret_key|aws_token|AWSSecretKey|b2_app_key|bashrc password|bintray_apikey|bintray_gpg_password|bintray_key|bintraykey|bluemix_api_key|bluemix_pass|browserstack_access_key|bucket_password|bucketeer_aws_access_key_id|bucketeer_aws_secret_access_key|built_branch_deploy_key|bx_password|cache_driver|cache_s3_secret_key|cattle_access_key|cattle_secret_key|certificate_password|ci_deploy_password|client_secret|client_zpk_secret_key|clojars_password|cloud_api_key|cloud_watch_aws_access_key|cloudant_password|cloudflare_api_key|cloudflare_auth_key|cloudinary_api_secret|cloudinary_name|codecov_token|config|conn.login|connectionstring|consumer_key|consumer_secret|credentials|cypress_record_key|database_password|database_schema_test|datadog_api_key|datadog_app_key|db_password|db_server|db_username|dbpasswd|dbpassword|dbuser|deploy_password|digitalocean_ssh_key_body|digitalocean_ssh_key_ids|docker_hub_password|docker_key|docker_pass|docker_passwd|docker_password|dockerhub_password|dockerhubpassword|dot-files|dotfiles|droplet_travis_password|dynamoaccesskeyid|dynamosecretaccesskey|elastica_host|elastica_port|elasticsearch_password|encryption_key|encryption_password|env.heroku_api_key|env.sonatype_password|eureka.awssecretkey)',
                                                      'color': 'yellow', 'state': True}]

portDic = [
    {'name': '100个常见端口',
     'value': '21,22,23,25,53,67,68,80,110,111,139,143,161,389,443,445,465,512,513,514,873,993,995,1080,1000,1352,1433,1521,1723,2049,2181,2375,3306,3389,4848,5000,5001,5432,5900,5632,5900,5989,6379,6666,7001,7002,8000,8001,8009,8010,8069,8080,8083,8086,8081,8088,8089,8443,8888,9900,9200,9300,9999,10621,11211,27017,27018,66,81,457,1100,1241,1434,1944,2301,3128,4000,4001,4002,4100,5800,5801,5802,6346,6347,30821,1090,1098,1099,4444,11099,47001,47002,10999,7000-7004,8000-8003,9000-9003,9503,7070,7071,45000,45001,8686,9012,50500,11111,4786,5555,5556,8880,8983,8383,4990,8500,6066'},
    {'name': 'nmap top 1000',
     'value': '1,3-4,6-7,9,13,17,19-26,30,32-33,37,42-43,49,53,70,79-85,88-90,99-100,106,109-111,113,119,125,135,139,143-144,146,161,163,179,199,211-212,222,254-256,259,264,280,301,306,311,340,366,389,406-407,416-417,425,427,443-445,458,464-465,481,497,500,512-515,524,541,543-545,548,554-555,563,587,593,616-617,625,631,636,646,648,666-668,683,687,691,700,705,711,714,720,722,726,749,765,777,783,787,800-801,808,843,873,880,888,898,900-903,911-912,981,987,990,992-993,995,999-1002,1007,1009-1011,1021-1100,1102,1104-1108,1110-1114,1117,1119,1121-1124,1126,1130-1132,1137-1138,1141,1145,1147-1149,1151-1152,1154,1163-1166,1169,1174-1175,1183,1185-1187,1192,1198-1199,1201,1213,1216-1218,1233-1234,1236,1244,1247-1248,1259,1271-1272,1277,1287,1296,1300-1301,1309-1311,1322,1328,1334,1352,1417,1433-1434,1443,1455,1461,1494,1500-1501,1503,1521,1524,1533,1556,1580,1583,1594,1600,1641,1658,1666,1687-1688,1700,1717-1721,1723,1755,1761,1782-1783,1801,1805,1812,1839-1840,1862-1864,1875,1900,1914,1935,1947,1971-1972,1974,1984,1998-2010,2013,2020-2022,2030,2033-2035,2038,2040-2043,2045-2049,2065,2068,2099-2100,2103,2105-2107,2111,2119,2121,2126,2135,2144,2160-2161,2170,2179,2190-2191,2196,2200,2222,2251,2260,2288,2301,2323,2366,2381-2383,2393-2394,2399,2401,2492,2500,2522,2525,2557,2601-2602,2604-2605,2607-2608,2638,2701-2702,2710,2717-2718,2725,2800,2809,2811,2869,2875,2909-2910,2920,2967-2968,2998,3000-3001,3003,3005-3007,3011,3013,3017,3030-3031,3052,3071,3077,3128,3168,3211,3221,3260-3261,3268-3269,3283,3300-3301,3306,3322-3325,3333,3351,3367,3369-3372,3389-3390,3404,3476,3493,3517,3527,3546,3551,3580,3659,3689-3690,3703,3737,3766,3784,3800-3801,3809,3814,3826-3828,3851,3869,3871,3878,3880,3889,3905,3914,3918,3920,3945,3971,3986,3995,3998,4000-4006,4045,4111,4125-4126,4129,4224,4242,4279,4321,4343,4443-4446,4449,4550,4567,4662,4848,4899-4900,4998,5000-5004,5009,5030,5033,5050-5051,5054,5060-5061,5080,5087,5100-5102,5120,5190,5200,5214,5221-5222,5225-5226,5269,5280,5298,5357,5405,5414,5431-5432,5440,5500,5510,5544,5550,5555,5560,5566,5631,5633,5666,5678-5679,5718,5730,5800-5802,5810-5811,5815,5822,5825,5850,5859,5862,5877,5900-5904,5906-5907,5910-5911,5915,5922,5925,5950,5952,5959-5963,5987-5989,5998-6007,6009,6025,6059,6100-6101,6106,6112,6123,6129,6156,6346,6389,6502,6510,6543,6547,6565-6567,6580,6646,6666-6669,6689,6692,6699,6779,6788-6789,6792,6839,6881,6901,6969,7000-7002,7004,7007,7019,7025,7070,7100,7103,7106,7200-7201,7402,7435,7443,7496,7512,7625,7627,7676,7741,7777-7778,7800,7911,7920-7921,7937-7938,7999-8002,8007-8011,8021-8022,8031,8042,8045,8080-8090,8093,8099-8100,8180-8181,8192-8194,8200,8222,8254,8290-8292,8300,8333,8383,8400,8402,8443,8500,8600,8649,8651-8652,8654,8701,8800,8873,8888,8899,8994,9000-9003,9009-9011,9040,9050,9071,9080-9081,9090-9091,9099-9103,9110-9111,9200,9207,9220,9290,9415,9418,9485,9500,9502-9503,9535,9575,9593-9595,9618,9666,9876-9878,9898,9900,9917,9929,9943-9944,9968,9998-10004,10009-10010,10012,10024-10025,10082,10180,10215,10243,10566,10616-10617,10621,10626,10628-10629,10778,11110-11111,11967,12000,12174,12265,12345,13456,13722,13782-13783,14000,14238,14441-14442,15000,15002-15004,15660,15742,16000-16001,16012,16016,16018,16080,16113,16992-16993,17877,17988,18040,18101,18988,19101,19283,19315,19350,19780,19801,19842,20000,20005,20031,20221-20222,20828,21571,22939,23502,24444,24800,25734-25735,26214,27000,27352-27353,27355-27356,27715,28201,30000,30718,30951,31038,31337,32768-32785,33354,33899,34571-34573,35500,38292,40193,40911,41511,42510,44176,44442-44443,44501,45100,48080,49152-49161,49163,49165,49167,49175-49176,49400,49999-50003,50006,50300,50389,50500,50636,50800,51103,51493,52673,52822,52848,52869,54045,54328,55055-55056,55555,55600,56737-56738,57294,57797,58080,60020,60443,61532,61900,62078,63331,64623,64680,65000,65129,65389,280,4567,7001,8008,9080'}
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
max_depth: 10                     # 最大页面深度限制
navigate_timeout_second: 10       # 访问超时时间，单位秒
load_timeout_second: 10           # 加载超时时间，单位秒
retry: 0                          # 页面访问失败后的重试次数
page_analyze_timeout_second: 100  # 页面分析超时时间，单位秒
max_interactive: 500             # 单个页面最大交互次数
max_interactive_depth: 10         # 页面交互深度限制
max_page_concurrent: 5           # 最大页面并发（不大于10）
max_page_visit: 1000              # 总共允许访问的页面数量
max_page_visit_per_site: 1000     # 每个站点最多访问的页面数量
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