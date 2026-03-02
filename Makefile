run:
	go run ./cmd/main

build:
	go build ./cmd/main

test:
	go test -v ./...

docs:
	swag init -g ./router.go -d ./internal/delivery/restapi -o ./internal/docs