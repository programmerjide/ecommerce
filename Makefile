.PHONY: run test build migrate-up migrate-down lint docker-dev docker-build clean help

# Variables
APP_NAME=ecommerce
BINARY_DIR=bin
MAIN_PATH=cmd/api/main.go
WORKER_PATH=cmd/worker/main.go
MIGRATE_PATH=migrate/main.go

# Database
DB_USER=ecommerce
DB_PASSWORD=ecommerce
DB_NAME=ecommerce
DB_HOST=localhost
DB_PORT=5432
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

run: ## Run the API server
	@echo "🚀 Starting API server..."
	@go run $(MAIN_PATH)

run-worker: ## Run background worker
	@echo "⚙️  Starting worker..."
	@go run $(WORKER_PATH)

test: ## Run tests
	@echo "🧪 Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

test-integration: ## Run integration tests
	@echo "🧪 Running integration tests..."
	@go test -v -race ./tests/integration/...

build: ## Build the application
	@echo "🔨 Building application..."
	@mkdir -p $(BINARY_DIR)
	@go build -o $(BINARY_DIR)/$(APP_NAME)-api $(MAIN_PATH)
	@go build -o $(BINARY_DIR)/$(APP_NAME)-worker $(WORKER_PATH)
	@go build -o $(BINARY_DIR)/$(APP_NAME)-migrate $(MIGRATE_PATH)
	@echo "✅ Build complete: $(BINARY_DIR)/"

clean: ## Clean build artifacts
	@echo "🧹 Cleaning..."
	@rm -rf $(BINARY_DIR)
	@rm -f coverage.out coverage.html
	@go clean

migrate-up: ## Run database migrations up
	@echo "⬆️  Running migrations..."
	@migrate -path migrations -database "$(DB_URL)" up

migrate-down: ## Rollback last migration
	@echo "⬇️  Rolling back migration..."
	@migrate -path migrations -database "$(DB_URL)" down 1

migrate-create: ## Create a new migration (usage: make migrate-create name=create_users)
	@echo "📝 Creating migration: $(name)"
	@migrate create -ext sql -dir migrations -seq $(name)

lint: ## Run linter
	@echo "🔍 Running linter..."
	@golangci-lint run --timeout 5m

format: ## Format code
	@echo "💅 Formatting code..."
	@go fmt ./...
	@goimports -w .

docker-dev: ## Start development environment
	@echo "🐳 Starting development environment..."
	@docker-compose -f docker/docker-compose.yml up -d

docker-down: ## Stop development environment
	@echo "🛑 Stopping development environment..."
	@docker-compose -f docker/docker-compose.yml down

docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	@docker build -f docker/Dockerfile -t $(APP_NAME):latest .

seed: ## Seed database with sample data
	@echo "🌱 Seeding database..."
	@go run scripts/seed.go

deps: ## Download dependencies
	@echo "📦 Downloading dependencies..."
	@go mod download
	@go mod tidy

install-tools: ## Install development tools
	@echo "🔧 Installing tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install github.com/cosmtrek/air@latest
	@go install golang.org/x/tools/cmd/goimports@latest

dev: ## Run with hot reload
	@echo "🔥 Starting with hot reload..."
	@air

.DEFAULT_GOAL := help
