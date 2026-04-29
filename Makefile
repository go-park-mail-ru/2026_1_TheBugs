run:
	go run ./cmd/main

build:
	go build ./cmd/main

test:
	go clean -testcache
	go test -v ./...

coverage:
	go clean -testcache
	-go test ./... -coverprofile cover.out -covermode=count 
	go tool cover -func cover.out

swag:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init --parseDependency --parseInternal -g ./register.go -d ./internal/delivery/restapi -o ./docs

keys:
	openssl genrsa -out private.pem 2048
	openssl rsa -in private.pem -pubout -out public.pem

install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint:
	golangci-lint run ./... --fix

format: 
	gofmt -w .

proto_install:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

run_auth:
	go run ./cmd/auth

auth_proto:
	mkdir -p internal/delivery/grpc/generated/auth
	protoc --proto_path=proto --go_out=internal/delivery/grpc/generated/auth --go-grpc_out=internal/delivery/grpc/generated/auth --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative ./proto/auth.proto