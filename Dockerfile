FROM docker.1ms.run/library/golang:latest

WORKDIR /works
COPY src.txt /works
RUN cp /works/src.txt /etc/apt/sources.list \
    && rm /etc/apt/sources.list.d/debian.sources \
    && apt update \
    && apt install -y awscli vim wget curl apt-transport-https ca-certificates

RUN go env -w GO111MODULE='on' && go env -w GOPROXY='https://goproxy.cn,direct' 

ADD go.* *.go /works/
RUN go build -o app .
COPY config.ini /works/

