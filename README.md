# Cocktail Recipes

### This is a training project for developing a web application on Go. The goal of the project is to implement a simple service with authorization, registration, and receiving posts with cocktail recipes.

### Chi web framework and sqlite database were used for the project

#### The authorization system is implemented using JSON WEB Token (JWT). The ECDSA algorithm was used to encode tokens.

#### Additionally implemented build system using Dockerfile

```Dockerfile
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
```

