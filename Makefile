# Go parameters
MAIN_PATH=cmd/http-ws-server

.PHONY: setup
setup: ## Setup Project
	go mod tidy

.PHONY: start-dev
start-dev: ## Start the local environment
	docker-compose --project-directory local-environment up -d

.PHONY: stop-dev
stop-dev: ## Stop the local environment
	docker-compose --project-directory local-environment down --remove-orphans

.PHONY: test
test: ## Run tests
	go test ./...

.PHONY: test-race
test-race: ## Run tests
	go test ./... -race

.PHONY: build
build: ## Build the application
	go build -o build/http-ws-server cmd/http-ws-server/main.go
	echo "Build completed: ./build/http-ws-server"

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
