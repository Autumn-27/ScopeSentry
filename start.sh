#!/bin/bash

# 判断 UPDATE 文件是否存在
if [ -f "/opt/ScopeSentry/UPDATE" ]; then
    echo "UPDATE 文件存在，开始从 URL 下载并解压 main.zip..."

    # 从 UPDATE 文件中读取 ZIP 文件的 URL
    zip_url=$(cat /opt/ScopeSentry/UPDATE)

    if [ -z "$zip_url" ]; then
        echo "ERROR: UPDATE 文件中没有 URL 地址"
    else
        # 下载 ZIP 文件到 /tmp/main.zip
        curl -o /tmp/main.zip "$zip_url" || echo "ERROR: ZIP 文件下载失败"

        # 解压 /tmp/main.zip 并覆盖原有的 main.zip
        echo "解压并覆盖 main.zip..."
        unzip -o /tmp/main.zip -d /opt/ScopeSentry/ || echo "ERROR: 解压失败"
        pip install -r /opt/ScopeSentry/requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple --no-cache-dir
        # 删除 UPDATE 文件
        rm -f /opt/ScopeSentry/UPDATE || echo "ERROR: 删除 UPDATE 文件失败"
    fi
else
    echo "UPDATE 文件不存在，跳过更新，直接启动 Python 应用..."
fi

# 进入工作目录
cd /opt/ScopeSentry/
# 启动 Python 应用
python main.py
