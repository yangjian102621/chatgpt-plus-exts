FROM alpine:3.18.2

MAINTAINER yangjian<yangjian102621@163.com>

WORKDIR /var/www/app
COPY ./src/bin/chatgpt-plus-exts-amd64-linux /var/www/app

EXPOSE 9001

# 容器启动时执行的命令
CMD ["./chatgpt-plus-exts-amd64-linux"]
