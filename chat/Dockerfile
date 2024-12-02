FROM golang:1.22.0 AS builder


ENV PROJECT_DIR=/app/chat \
    GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /app


COPY ./chat ./chat
COPY ./logging ./logging
COPY ./chat/go.mod ./chat/go.sum ./chat
COPY ./chat/config ./chat/config


RUN cd ./chat && \
    go get github.com/githubnemo/CompileDaemon && \
    go install github.com/githubnemo/CompileDaemon && \
    go mod download

RUN ls -l /app

    
WORKDIR /app/chat
ENTRYPOINT CompileDaemon -build="go build -a -installsuffix cgo -o main ./cmd/main.go" -command="./main"
