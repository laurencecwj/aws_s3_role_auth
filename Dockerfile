FROM docker.1ms.run/library/golang:latest

WORKDIR /works
COPY . /works
RUN cp /works/src.txt /etc/apt/sources.list \
    && rm /etc/apt/sources.list.d/debian.sources \
    && apt update \
    && apt install -y awscli vim wget curl apt-transport-https ca-certificates

RUN go build -o app .

