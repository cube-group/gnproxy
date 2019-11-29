FROM nginx:alpine

ENV APP_PATH /run

COPY conf/nginx.conf /etc/nginx/nginx.conf
COPY bin/gnproxy $APP_PATH/gnproxy

RUN mkdir -p /etc/nginx/conf.d && \
    mkdir -p /etc/nginx/stream.d && \
    mkdir -p /etc/nginx/basic_auth


WORKDIR $APP_PATH
ENTRYPOINT ["./gnproxy"]