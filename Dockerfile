FROM golang:1.13 AS build

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 \
    BINARY_NAME=fluentbit-plugin-natspublisher
RUN apt update 

COPY ./default-entrypoint.sh /entrypoint.sh
RUN chmod u+x /entrypoint.sh
ENTRYPOINT "/entrypoint.sh"
