FROM golang:1.16 as litestream-builder

#TODO: Replace this with normal litestream.

WORKDIR /src/litestream

RUN git clone https://github.com/benbjohnson/litestream.git .

RUN --mount=type=cache,target=/root/.cache/go-build \
	--mount=type=cache,target=/go/pkg \
	go build -ldflags '-s -w -extldflags "-static"' -tags osusergo,netgo,sqlite_omit_load_extension -o /usr/local/bin/litestream ./cmd/litestream

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

COPY --from=litestream-builder /usr/local/bin/litestream /usr/local/bin/litestream

COPY --from=builder /app/server /app/server
COPY --from=builder /app/views /app/views
COPY --from=builder /app/static /app/static
COPY ./docker_entrypoint /app/docker_entrypoint

WORKDIR /app

# Frequency that database snapshots are replicated.
ENV DB_SYNC_INTERVAL="10s"

ENTRYPOINT ["/app/docker_entrypoint"]
