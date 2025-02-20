FROM golang:1.22.2 as builder

WORKDIR /app
COPY . .

RUN go mod tidy && go build -o server ./cmd/main.go

FROM alpine:latest
WORKDIR /app/

COPY --from=builder /app/server .
COPY --from=builder /app/config ./config

ENV CONFIG_PATH=./config/local.yaml
EXPOSE 8080
CMD ["./server"]
