FROM golang:1.19.3-alpine3.15 as builder

WORKDIR /app
COPY ./ ./

RUN go mod tidy && \
    go build -o /bin/gocrawler /app/cmd/gocrawler/

FROM alpine:3.15

ARG UID=1000
ARG GID=1000

RUN apk update && apk add shadow && \
    useradd --create-home --shell /sbin/nologin -u ${UID} gocrawler && \
    mkdir /gocrawler && \
    chown gocrawler:${GID} /gocrawler

COPY --from=builder --chown=gocrawler /bin/gocrawler /gocrawler/gocrawler

USER gocrawler
WORKDIR /gocrawler

ENTRYPOINT ["/gocrawler/gocrawler"]
