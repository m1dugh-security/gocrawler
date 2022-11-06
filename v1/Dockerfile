FROM golang:1.19.3-alpine3.15 as builder

WORKDIR /app
COPY . /app

RUN go build -o crawler .

FROM alpine:3.15

ARG UID=1000
ARG GID=1000

RUN apk update && apk add shadow && \
    useradd --create-home --shell /sbin/nologin -u ${UID} crawler && \
    mkdir /crawler && \
    chown crawler:${GID} /crawler

COPY --from=builder --chown=crawler /app/crawler /crawler/crawler

USER crawler
WORKDIR /crawler

ENTRYPOINT ["/crawler/crawler"]
