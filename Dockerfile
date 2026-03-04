FROM golang:1.25 AS builder

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -o app ./cmd/main

ENTRYPOINT ["./app"]