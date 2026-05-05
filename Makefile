run:
	go run ./cmd/main

build:
	go build ./cmd/main

test:
	go clean -testcache
	go test -v ./...

coverage:
	go clean -testcache
	-go test ./internal/delivery/... ./internal/re -coverprofile cover.out -covermode=count 
	go tool cover -func cover.out

coverage-win:
	$pkgs = go list ./... | Where-Object { $_ -notmatch '/docs|/utils|/config|/cmd|/generated|/app|/entity|/mocks|/dto|/metrics|/logger|/middleware|/response|/request|/order' }
	go test $pkgs -coverprofile cover.out -covermode=count
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

run_user:
	go run ./cmd/user	

run_poster:
	go run ./cmd/poster

auth_proto:
	mkdir -p internal/delivery/grpc/generated/auth
	protoc --proto_path=proto --go_out=internal/delivery/grpc/generated/auth --go-grpc_out=internal/delivery/grpc/generated/auth --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative ./proto/auth.proto

user_proto:
	mkdir -p internal/delivery/grpc/generated/user
	protoc --proto_path=proto --go_out=internal/delivery/grpc/generated/user --go-grpc_out=internal/delivery/grpc/generated/user --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative ./proto/user.proto

poster_proto:
	mkdir -p internal/delivery/grpc/generated/poster
	protoc --proto_path=proto --go_out=internal/delivery/grpc/generated/poster --go-grpc_out=internal/delivery/grpc/generated/poster --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative ./proto/poster.proto
  
complex_proto:
	mkdir -p internal/delivery/grpc/generated/complex
	protoc --proto_path=proto --go_out=internal/delivery/grpc/generated/complex --go-grpc_out=internal/delivery/grpc/generated/complex --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative ./proto/complex.proto

grpc_mock:
	mockgen -destination internal/mocks/grpc_client/mock_auth_client.go -package grpc_client github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/generated/auth AuthServiceClient
	mockgen -destination internal/mocks/grpc_client/mock_user_client.go -package grpc_client github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/generated/user UserServiceClient
	mockgen -destination internal/mocks/grpc_client/mock_poster_client.go -package grpc_client github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/generated/poster PosterServiceClient
	mockgen -destination internal/mocks/grpc_client/mock_complex_client.go -package grpc_client github.com/go-park-mail-ru/2026_1_TheBugs/internal/delivery/grpc/generated/complex ComplexServiceClient
