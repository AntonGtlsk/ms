FROM golang:1.22.0 AS builder


ENV PROJECT_DIR=/app/auth \
    GO111MODULE=on \
    CGO_ENABLED=0


WORKDIR /app


COPY ./auth ./auth
COPY ./logging ./logging
COPY ./auth/go.mod ./auth/go.sum ./auth
COPY ./auth/config ./auth/config


RUN cd ./auth && \
    go get github.com/gorilla/mux && \
    go get github.com/githubnemo/CompileDaemon && \
    go install github.com/githubnemo/CompileDaemon && \
    go mod download

RUN ls -l /app

WORKDIR /app/auth
ENTRYPOINT CompileDaemon -build="go build -a -installsuffix cgo -o main ./cmd/main.go" -command="./main"
