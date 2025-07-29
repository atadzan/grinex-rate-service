.PHONY: build test docker-build run lint clean proto

# Build the application
build:
	go build -o bin/grinex-rate-service ./cmd

# Run unit tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Build Docker image
docker-build:
	docker build -t grinex-rate-service .

# Run the application locally
run:
	go run ./cmd

# Run with Docker Compose
run-docker:
	docker-compose up --build

# Stop Docker Compose
stop-docker:
	docker-compose down

# Run linter
lint:
	golangci-lint run

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Generate protobuf files
proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/v1/rate-service.proto

# Install dependencies
deps:
	docker compose up -d

# Run database migrations
migrate-up:
	migrate -path migrations -database "postgres://postgres:password@localhost:5432/grinex_rates?sslmode=disable" up

# Rollback database migrations
migrate-down:
	migrate -path migrations -database "postgres://postgres:password@localhost:5432/grinex_rates?sslmode=disable" down

# Show help
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  test          - Run unit tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  docker-build  - Build Docker image"
	@echo "  run           - Run the application locally"
	@echo "  run-docker    - Run with Docker Compose"
	@echo "  stop-docker   - Stop Docker Compose"
	@echo "  lint          - Run linter"
	@echo "  clean         - Clean build artifacts"
	@echo "  proto         - Generate protobuf files"
	@echo "  deps          - Install dependencies"
	@echo "  migrate-up    - Run database migrations"
	@echo "  migrate-down  - Rollback database migrations"