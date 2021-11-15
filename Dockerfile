FROM golang:1.17.3-buster as backend_builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . /app/

RUN go build \
  -mod=readonly \
  -v \
  -o /app/server \
  ./main.go

FROM golang:1.17.3-buster AS litestream_builder

ARG litestream_version="v0.3.6"

RUN set -x && \
    apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y git

RUN set -x && \
    git clone --branch "${litestream_version}" --single-branch https://github.com/benbjohnson/litestream.git

RUN set -x && \
    cd litestream && \
    go install ./cmd/litestream && \
    echo "litestream installed to ${GOPATH}/bin/litestream"

FROM debian:stable-20211011-slim

COPY --from=backend_builder /app/server /app/server
COPY --from=backend_builder /app/views /app/views
COPY --from=backend_builder /app/static /app/static
COPY --from=litestream_builder /go/bin/litestream /app/litestream
COPY ./litestream.yml /etc/litestream.yml
COPY ./docker_entrypoint /app/docker_entrypoint
COPY ./litestream.yml /etc/litestream.yml
COPY ./docker_entrypoint /app/docker_entrypoint

WORKDIR /app

# Frequency that database snapshots are replicated.
ENV DB_SYNC_INTERVAL="10s"

ENTRYPOINT ["/app/docker_entrypoint"]
