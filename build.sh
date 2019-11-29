#!/usr/bin/env bash
export GOPROXY=https://goproxy.cn
CGO_ENABLED=0 GOOS=linux go build -o bin/gnproxy
docker build -t gnproxy .

docker rm -f gnproxy
docker run -it -d -p 8888:80 --name gnproxy gnproxy