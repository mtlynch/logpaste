FROM golang:1.13.5-buster

COPY . /app/
WORKDIR /app

RUN go build \
  -o /app/main \
  ./main.go

ENTRYPOINT ["/app/main"]
