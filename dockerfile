FROM debian:buster-slim AS git_installer

RUN sed -i 's/deb.debian.org/mirrors.aliyun.com/g' /etc/apt/sources.list && \
    sed -i 's/security.debian.org/mirrors.aliyun.com\/debian-security/g' /etc/apt/sources.list

RUN apt-get update && \
    apt-get install -y git curl && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

FROM python:3.7-slim

ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

COPY --from=git_installer /usr/bin/git /usr/bin/git
COPY --from=git_installer /usr/bin/curl /usr/bin/curl

RUN apt-get update && \
    apt-get install -y gcc  && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /opt/ScopeSentry/

COPY ./ScopeSentry /opt/ScopeSentry/

RUN pip install -r requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple --no-cache-dir

CMD ["python", "main.py"]
