FROM golang:1.13.5-buster as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . /app/

RUN go build \
  -mod=readonly \
  -v \
  -o /app/server \
  ./main.go

FROM debian:stable-20210208-slim

ARG litestream_version="0.3.2"
ARG litestream_deb_filename="litestream-v${litestream_version}-linux-amd64.deb"

RUN set -x && \
    apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y \
      ca-certificates \
      wget \
      && \
    wget "https://github.com/benbjohnson/litestream/releases/download/v${litestream_version}/${litestream_deb_filename}" && \
    apt-get remove -y wget && \
    rm -rf /var/lib/apt/lists/*

RUN dpkg -i "${litestream_deb_filename}"

COPY --from=builder /app/server /app/server
COPY --from=builder /app/views /app/views
COPY --from=builder /app/static /app/static
COPY ./docker_entrypoint /app/docker_entrypoint

WORKDIR /app

# Frequency that database snapshots are replicated.
ENV DB_SYNC_INTERVAL="10s"

# Should logpaste create a fresh database or pull down the latest replicated
# version?
ENV CREATE_NEW_DB="false"

# Properties that control the site UI.
ENV SITE_TITLE="Log Paster"
ENV SITE_SUBTITLE="A minimalist, open-source debug log upload service"

ENTRYPOINT ["/app/docker_entrypoint"]
