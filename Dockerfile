FROM golang:1.21.1 as backend_builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY ./cmd /app/cmd
COPY ./dev-scripts /app/dev-scripts
COPY ./handlers /app/handlers
COPY ./limit /app/limit
COPY ./random /app/random
COPY ./store /app/store

RUN TARGETPLATFORM="${TARGETPLATFORM}" ./dev-scripts/build-backend "prod"

FROM litestream/litestream:0.3.13 AS litestream

FROM alpine:3.15

RUN apk add --no-cache bash

COPY --from=backend_builder /app/bin/logpaste /app/logpaste
COPY --from=litestream /usr/local/bin/litestream /app/litestream
COPY ./docker-entrypoint /app/docker-entrypoint
COPY ./litestream.yml /etc/litestream.yml
COPY ./views /app/views

WORKDIR /app

# Frequency that database snapshots are replicated.
ENV DB_SYNC_INTERVAL="10s"

ENTRYPOINT ["/app/docker-entrypoint"]
