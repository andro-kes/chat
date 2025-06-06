FROM golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download
COPY ./ ./
RUN go build -o server server.go
ENTRYPOINT [ "/app/server" ]