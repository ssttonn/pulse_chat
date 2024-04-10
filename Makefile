.PHONY: help build test lint

help: ## Hiển thị danh sách các lệnh
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

build: ## Build entire project
	@echo "Building services..."
	# ví dụ: go build -o bin/api-core ./src/api-core/cmd/api
	@echo "Build complete."

test: ## Run unit tests with race detector
	@echo "Running tests..."
	go test -v -race ./...

lint: ## Run golangci-lint to check code quality
	@echo "Running linter..."
	golangci-lint run