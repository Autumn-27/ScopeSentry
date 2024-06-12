FROM python:3.7-slim

ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 更新包列表并安装必要的包，包括 nginx
RUN apt-get update && \
    apt-get install -y git curl ca-certificates libcurl4-openssl-dev nginx && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# 设置工作目录
WORKDIR /opt/ScopeSentry/

# 复制 ScopeSentry 项目文件到工作目录
COPY ./ScopeSentry /opt/ScopeSentry/

# 安装 Python 依赖包
RUN pip install -r requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple --no-cache-dir

# 移动 static 目录下的所有文件到 nginx 的默认根目录
RUN cp -r /opt/ScopeSentry/static/* /var/www/html/

# 修改 Nginx 配置文件，添加反向代理规则
RUN echo 'server {\n\
    listen 80;\n\
    server_name localhost;\n\
\n\
    location / {\n\
        root /var/www/html;\n\
        try_files $uri $uri/ =404;\n\
    }\n\
\n\
    location /api/ {\n\
        proxy_pass http://127.0.0.1:8082;\n\
        proxy_set_header Host $host;\n\
        proxy_set_header X-Real-IP $remote_addr;\n\
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;\n\
        proxy_set_header X-Forwarded-Proto $scheme;\n\
    }\n\
}\n' > /etc/nginx/sites-available/default

# 确保 Nginx 默认配置启用
RUN ln -sf /etc/nginx/sites-available/default /etc/nginx/sites-enabled/

# 启动 Nginx 和 Python 应用
CMD service nginx start && python main.py
