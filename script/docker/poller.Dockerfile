FROM golang:1.26-alpine3.23 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -a -installsuffix cgo -ldflags="-w -s" -o ./bin/poller cmd/poller/main.go

FROM alpine:3.23

WORKDIR /app/

COPY --from=builder /app/bin/poller /app/poller

ENTRYPOINT ["./poller", "-config", "./config/config-poller.yaml"]
