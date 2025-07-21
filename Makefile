.PHONY: build run test clean docker-build docker-run migrate-up migrate-down deps

# Build the application
build:
	go build -o bin/gocbt cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -f *.db

# Install dependencies
deps:
	go mod tidy
	go mod download

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Build Docker image
docker-build:
	docker build -t gocbt:latest .

# Run with Docker Compose (development)
docker-run:
	docker-compose up --build

# Run with Docker Compose (production)
docker-run-prod:
	docker-compose --profile production up --build

# Stop Docker Compose
docker-stop:
	docker-compose down

# Database migrations (will implement later)
migrate-up:
	@echo "Running database migrations..."
	# go run migrations/migrate.go up

migrate-down:
	@echo "Rolling back database migrations..."
	# go run migrations/migrate.go down

# Create new migration
migrate-create:
	@echo "Creating new migration: $(name)"
	# go run migrations/migrate.go create $(name)

# Development setup
dev-setup: deps
	@echo "Setting up development environment..."
	mkdir -p data logs
	@echo "Development environment ready!"

# Help
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Install dependencies"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run with Docker Compose (dev)"
	@echo "  docker-run-prod - Run with Docker Compose (prod)"
	@echo "  docker-stop   - Stop Docker Compose"
	@echo "  migrate-up    - Run database migrations"
	@echo "  migrate-down  - Rollback database migrations"
	@echo "  dev-setup     - Setup development environment"
	@echo "  help          - Show this help message"
