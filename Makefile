run:
	go run ./cmd/main

build:
	go build ./cmd/main

test:
	go clean -testcache
	go test -v ./... --cover

docs:
	swag init --parseDependency --parseInternal -g ./register.go -d ./internal/delivery/restapi -o ./internal/docs

keys:
	openssl genrsa -out private.pem 2048
	openssl rsa -in private.pem -pubout -out public.pem