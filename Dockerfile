FROM nginx:alpine

ENV APP_PATH /run

COPY bin/gnproxy $APP_PATH/gnproxy

ENTRYPOINT ["./gnproxy"]