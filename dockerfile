FROM python:3.7-slim

ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 更新包列表并安装必要的包，包括 nginx
RUN apt-get update && \
    apt-get install -y git curl ca-certificates libcurl4-openssl-dev vim unzip && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# 设置工作目录
WORKDIR /opt/ScopeSentry/

# 复制 ScopeSentry 项目文件到工作目录
COPY dist/ScopeSentry_linux_amd64_v1/ScopeSentry /opt/ScopeSentry/

RUN chmod +x /opt/ScopeSentry/ScopeSentry

ENTRYPOINT ["/opt/ScopeSentry/ScopeSentry"]

# # 复制 start.sh 脚本到容器
# RUN cp /opt/ScopeSentry/start.sh /usr/local/bin/start.sh
#
# # 给 start.sh 脚本赋予执行权限
# RUN chmod +x /usr/local/bin/start.sh
#
# # 使用 start.sh 启动容器
# CMD ["/usr/local/bin/start.sh"]