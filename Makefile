.PHONY: build run clean test fmt lint help

APP_NAME := sds-admin
VERSION := 1.0.0
BUILD_DIR := bin
MAIN_PATH := ./cmd/sds-admin

GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.gitCommit=$(GIT_COMMIT) -X main.buildTime=$(BUILD_TIME)"

help:
	@echo "Available targets:"
	@echo "  build     - Build the application"
	@echo "  build-all - Build frontend and backend with static files"
	@echo "  run       - Run the application"
	@echo "  clean     - Clean build artifacts"
	@echo "  test      - Run tests"
	@echo "  fmt       - Format code"
	@echo "  lint      - Run linter"
	@echo "  swagger   - Generate Swagger documentation"
	@echo "  deps      - Install dependencies"
	@echo "  help      - Show this help message"

deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

build-all:
	@echo "Building frontend..."
	@cd fe && npm install && npm run build
	@echo "Copying frontend files to pub directory..."
	@mkdir -p pub
	@cp -r fe/dist/* pub/
	@echo "Packaging frontend files with statik..."
	@rm -rf internal/static
	@statik -src=pub -dest=internal -p=static
	@echo "Building backend..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

run:
	@echo "Running $(APP_NAME)..."
	go run $(MAIN_PATH) --config configs/config.yaml

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	go clean

test:
	@echo "Running tests..."
	go test -v -race -cover ./...

fmt:
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

lint:
	@echo "Running linter..."
	golangci-lint run ./...

swagger:
	@echo "Generating Swagger documentation..."
	swag init -g $(MAIN_PATH)/main.go -o docs --parseDependency

docker-build:
	@echo "Building Docker image..."
	docker build -t $(APP_NAME):$(VERSION) .

docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 $(APP_NAME):$(VERSION)

.DEFAULT_GOAL := help
