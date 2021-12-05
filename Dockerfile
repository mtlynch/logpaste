FROM golang:1.17.3-buster as backend_builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY ./handlers /app/handlers
COPY ./limit /app/limit
COPY ./random /app/random
COPY ./store /app/store
COPY main.go ./

RUN GOOS=linux GOARCH=amd64 \
  go build \
  -tags netgo \
  -ldflags '-w -extldflags "-static"' \
  -mod=readonly \
  -v \
  -o /app/server \
  ./main.go

FROM debian:stable-20211011-slim AS litestream_downloader

ARG litestream_version="v0.3.7"
ARG litestream_binary_tgz_filename="litestream-${litestream_version}-linux-amd64-static.tar.gz"

WORKDIR /litestream

RUN set -x && \
    apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y \
      ca-certificates \
      wget
RUN wget "https://github.com/benbjohnson/litestream/releases/download/${litestream_version}/${litestream_binary_tgz_filename}"
RUN tar -xvzf "${litestream_binary_tgz_filename}"

FROM alpine:3.15

RUN apk add --no-cache bash

COPY --from=backend_builder /app/server /app/server
COPY --from=litestream_downloader /litestream/litestream /app/litestream
COPY ./docker_entrypoint /app/docker_entrypoint
COPY ./litestream.yml /etc/litestream.yml
COPY ./static /app/static
COPY ./views /app/views

WORKDIR /app

# Frequency that database snapshots are replicated.
ENV DB_SYNC_INTERVAL="10s"

ENTRYPOINT ["/app/docker_entrypoint"]
