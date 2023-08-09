#!/bin/bash

version=$1
if [ -z $version ];then
	echo "Please input version"
	exit 1
fi

cd ../api
make clean linux

cd ../web
npm run build

cd ../docker

# remove docker image if exists
docker rmi -f registry.cn-hangzhou.aliyuncs.com/geekmaster/chatgpt-plus-exts:$version
# build docker image
docker build -t registry.cn-hangzhou.aliyuncs.com/geekmaster/chatgpt-plus-exts:$version -f Dockerfile ../

if [ "$2" = "push" ];then
  docker push registry.cn-hangzhou.aliyuncs.com/geekmaster/chatgpt-plus-exts:$version
fi





