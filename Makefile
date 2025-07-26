# Makefile - помещается в корень проекта
.PHONY: build run test clean migrate-up migrate-down migrate-reset migrate-version \
    install-deps install-migrate docker-build docker-run docker-down

# Default variables
APP_NAME ?= chat-service
BINARY_NAME ?= chat-service
MIGRATE_PATH ?= ./migrations
DATABASE_URL ?= postgres://chatuser:chatpass@localhost:5432/chatdb?sslmode=disable

# Build the application
build:
    go build -o $(BINARY_NAME) ./cmd/server

# Run the application
run: build
    ./$(BINARY_NAME)

# Run tests
test:
    go test -v ./...

# Clean build artifacts
clean:
    rm -f $(BINARY_NAME)
    go clean

# Install dependencies
install-deps:
    go mod tidy
    go get -u github.com/golang-migrate/migrate/v4

# Install migrate tool
install-migrate:
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations up
migrate-up:
    go run cmd/migrate/main.go -path=$(MIGRATE_PATH) -database=$(DATABASE_URL) -action=up

# Run migrations down
migrate-down:
    go run cmd/migrate/main.go -path=$(MIGRATE_PATH) -database=$(DATABASE_URL) -action=down -steps=1

# Reset migrations
migrate-reset:
    go run cmd/migrate/main.go -path=$(MIGRATE_PATH) -database=$(DATABASE_URL) -action=reset

# Show current migration version
migrate-version:
    go run cmd/migrate/main.go -path=$(MIGRATE_PATH) -database=$(DATABASE_URL) -action=version

# Docker commands
docker-build:
    docker-compose -f docker/docker-compose.yml build

docker-up:
    docker-compose -f docker/docker-compose.yml up -d

docker-down:
    docker-compose -f docker/docker-compose.yml down

docker-logs:
    docker-compose -f docker/docker-compose.yml logs -f

# Development commands
dev: install-deps build run

# Help
help:
    @echo "Available commands:"
    @echo "  build          - Build the application"
    @echo "  run            - Run the application"
    @echo "  test           - Run tests"
    @echo "  clean          - Clean build artifacts"
    @echo "  install-deps   - Install dependencies"
    @echo "  install-migrate - Install migrate tool"
    @echo "  migrate-up     - Run migrations up"
    @echo "  migrate-down   - Run migrations down"
    @echo "  migrate-reset  - Reset migrations"
    @echo "  migrate-version - Show migration version"
    @echo "  docker-build   - Build Docker images"
    @echo "  docker-up      - Start Docker containers"
    @echo "  docker-down    - Stop Docker containers"
    @echo "  docker-logs    - Show Docker logs"
    @echo "  dev            - Install deps, build and run"
    @echo "  help           - Show this help"


# Docker environment variables
DOCKER_COMPOSE_FILE ?= docker/docker-compose.yml
DOCKER_COMPOSE_PROJECT ?= chat-service

# Docker commands
docker-build:
    docker-compose -f $(DOCKER_COMPOSE_FILE) build

docker-up:
    docker-compose -f $(DOCKER_COMPOSE_FILE) up -d

docker-down:
    docker-compose -f $(DOCKER_COMPOSE_FILE) down

docker-logs:
    docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f

docker-ps:
    docker-compose -f $(DOCKER_COMPOSE_FILE) ps

docker-exec-db:
    docker exec -it chat-postgres psql -U chatuser -d chatdb

docker-exec-app:
    docker exec -it chat-service sh

docker-restart:
    docker-compose -f $(DOCKER_COMPOSE_FILE) restart

docker-stop:
    docker-compose -f $(DOCKER_COMPOSE_FILE) stop

docker-start:
    docker-compose -f $(DOCKER_COMPOSE_FILE) start

# Build and run with Docker
docker-dev: docker-build docker-up

# Clean Docker resources
docker-clean:
    docker-compose -f $(DOCKER_COMPOSE_FILE) down -v --remove-orphans
    docker system prune -f

# Swagger documentation
swagger-init:
    swag init -g cmd/server/main.go -o internal/docs

swagger-fmt:
    swag fmt -d internal/handler

# Install swag
install-swag:
    go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs
generate-docs: swagger-fmt swagger-init