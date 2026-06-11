FROM alpine


RUN apk update
RUN apk upgrade

ENTRYPOINT [ "echo", "hello", "world" ]
