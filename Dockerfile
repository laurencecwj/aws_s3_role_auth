FROM docker.1ms.run/library/golang:latest

WORKDIR /works
COPY src.txt /works
RUN rm /etc/apt/sources.list.d/* \
    && cp /works/src.txt /etc/apt/sources.list \
    && apt update \
    && apt install -y awscli vim wget curl apt-transport-https ca-certificates python3-pip python3-venv

RUN go env -w GO111MODULE='on' && go env -w GOPROXY='https://goproxy.cn,direct' 

ADD go.* *.go /works/
RUN go build -o app .

RUN python3 -m venv /works/mypy && . /works/mypy/bin/activate \
&& pip3 install minio boto3 --index-url https://mirrors.aliyun.com/pypi/simple 

# COPY main.py /works/
# RUN /works/mypy/bin/python main.py

COPY . /works/

