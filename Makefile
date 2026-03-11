run:
	go run ./cmd/main

build:
	go build ./cmd/main

test:
	go clean -testcache
	go test -v ./... --cover

swag:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init --parseDependency --parseInternal -g ./register.go -d ./internal/delivery/restapi -o ./internal/docs

keys:
	openssl genrsa -out private.pem 2048
	openssl rsa -in private.pem -pubout -out public.pem

install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint:
	golangci-lint run ./... --fix

format: 
	gofmt -w .