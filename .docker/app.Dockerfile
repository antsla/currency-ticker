FROM golang:1.17-alpine as builder

ARG SERVICE_NAME

COPY ./${SERVICE_NAME} /app/${SERVICE_NAME}
WORKDIR /app/${SERVICE_NAME}

RUN apk add --no-cache nano bash shadow build-base

RUN mkdir /.cache
RUN chown nobody:nobody -R /.cache

USER nobody