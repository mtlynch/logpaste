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
RUN set -x && \
    apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y \
      ca-certificates \
      wget \
      && \
    rm -rf /var/lib/apt/lists/*

ARG litestream_version="0.3.2"
ARG lisestram_deb_filename="litestream-v${litestream_version}-linux-amd64.deb"
RUN wget "https://github.com/benbjohnson/litestream/releases/download/v${litestream_version}/${lisestram_deb_filename}"
RUN dpkg -i "${lisestram_deb_filename}"

COPY --from=builder /app/server /app/server
COPY --from=builder /app/views /app/views

WORKDIR /app

ENTRYPOINT ["/app/docker_entrypoint"]
