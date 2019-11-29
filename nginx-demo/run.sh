#!/usr/bin/env bash

docker rm -f nginx
docker run -d -it \
-v $(pwd)/nginx.conf:/etc/nginx/nginx.conf \
-v $(pwd)/default.conf:/etc/nginx/conf.d/default.conf \
-p 8088:80 \
--name nginx \
nginx:alpine