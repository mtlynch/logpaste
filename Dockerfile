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

FROM debian:buster-slim
RUN set -x && \
    apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/server /app/server
COPY --from=builder /app/views /app/views

WORKDIR /app

ENTRYPOINT ["/app/server"]
