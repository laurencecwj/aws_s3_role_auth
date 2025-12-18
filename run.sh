#!/bin/bash

# resolvectl flush-caches

S3NAME=`python3 -c "import configparser; config = configparser.ConfigParser(); config.read('./config.ini'); print(config['aws']['s3_bucket'])"`

docker build -t cwj/test1:latest .

echo "try to s3 ls s3://$S3NAME"
docker run --rm -it cwj/test1:latest aws s3 ls s3://$S3NAME
echo ""

echo "try to execute golang app inside docker"
docker run --rm -it cwj/test1:latest /works/app
echo ""

echo "try to execute python app inside docker"
docker run --rm -it cwj/test1:latest /works/mypy/bin/python main.py
echo ""