import os, json
from core.util import string_to_postfix
current_directory = os.getcwd()

dict_directory = "..\dicts"
combined_directory = os.path.join(current_directory, dict_directory)
def get_fingerprint():
    try:
        # 尝试打开文件并读取内容
        with open(os.path.join(combined_directory, "fingerprint"), "r", encoding="utf-8") as file:
            fingerprint = file.read()
    except FileNotFoundError:
        print("文件不存在")
    return json.loads(fingerprint)

if __name__ == '__main__':
    fingerprint_rules = get_fingerprint()
    for rule in fingerprint_rules:
        express = string_to_postfix(rule['rule'])
        print(express)