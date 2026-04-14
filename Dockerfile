FROM golang:1.25 AS builder

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -ldflags="-s -w" -o app ./cmd/main


FROM alpine:3.20

WORKDIR /app

COPY --from=builder /build/app .

CMD ["./app"]
