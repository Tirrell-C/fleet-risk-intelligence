.PHONY: help build test lint clean deploy dev

# Variables
DOCKER_REGISTRY ?= your-registry.com
PROJECT_NAME = samsara-risk-mvp
VERSION ?= latest

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: ## Install all dependencies and initialize project
	@echo "🔧 Setting up project dependencies..."
	go mod download
	cd frontend && npm install
	go work sync

build: ## Build all services and frontend
	@echo "🏗️  Building all services..."
	go build -o bin/api ./services/api
	go build -o bin/risk-engine ./services/risk-engine
	go build -o bin/telemetry-ingest ./services/telemetry-ingest
	go build -o bin/websocket ./services/websocket
	cd frontend && npm run build

test: ## Run all tests
	@echo "🧪 Running tests..."
	go test ./...
	cd frontend && npm test

lint: ## Run linters for all code
	@echo "🔍 Running linters..."
	golangci-lint run ./...
	cd frontend && npm run lint

fmt: ## Format all code
	@echo "📝 Formatting code..."
	go fmt ./...
	goimports -w .
	cd frontend && npm run lint:fix

docker-build: ## Build Docker images for all services
	@echo "🐳 Building Docker images..."
	docker build -f docker/api.Dockerfile -t $(PROJECT_NAME)/api:$(VERSION) .
	docker build -f docker/risk-engine.Dockerfile -t $(PROJECT_NAME)/risk-engine:$(VERSION) .
	docker build -f docker/telemetry-ingest.Dockerfile -t $(PROJECT_NAME)/telemetry-ingest:$(VERSION) .
	docker build -f docker/websocket.Dockerfile -t $(PROJECT_NAME)/websocket:$(VERSION) .
	docker build -f docker/frontend.Dockerfile -t $(PROJECT_NAME)/frontend:$(VERSION) .

dev: ## Start development environment
	@echo "🚀 Starting development environment..."
	docker-compose -f docker-compose.dev.yml up -d

dev-down: ## Stop development environment
	@echo "🛑 Stopping development environment..."
	docker-compose -f docker-compose.dev.yml down

deploy-infra: ## Deploy infrastructure with Pulumi
	@echo "☁️  Deploying infrastructure..."
	cd infrastructure && pulumi up

deploy: docker-build deploy-infra ## Build and deploy everything
	@echo "🚀 Deploying application..."

clean: ## Clean build artifacts
	@echo "🧹 Cleaning up..."
	rm -rf bin/
	rm -rf frontend/dist/
	docker system prune -f

generate: ## Generate GraphQL schema and types
	@echo "⚡ Generating GraphQL code..."
	cd services/api && go run github.com/99designs/gqlgen generate