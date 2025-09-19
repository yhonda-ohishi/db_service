.PHONY: proto build test run clean

# Protocol Buffers compilation
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		src/proto/ryohi.proto

# Build the server
build:
	go build -o bin/server cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -cover ./...

# Run integration tests
test-integration:
	go test -v -tags=integration ./tests/integration/...

# Run the server
run:
	go run cmd/server/main.go

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f src/proto/*.pb.go

# Install protoc plugins
install-proto:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Setup development environment
setup: install-proto
	cp .env.example .env
	@echo "Please edit .env file with your database credentials"

# Database migration
migrate:
	go run cmd/migrate/main.go

# All: compile proto, build, and test
all: proto build test