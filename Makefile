TASK := go tool task

.DEFAULT_GOAL := help

.PHONY: help build generate start dev run debug test download kill \
        lint fmt vet tidy check clean

help: ## Show this help
	@echo "Targets:"
	@grep -hE '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "} {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

build: ## Production build to bin/main (generates templ first)
	$(TASK) build

generate: ## Generate templ *_templ.go files
	$(TASK) build:templ

start: ## Dev server with hot reload (templ watch + air)
	$(TASK) start

dev: start ## Alias for `start`

run: ## Build then run bin/main
	$(TASK) run

debug: ## Build then launch the delve debugger
	$(TASK) debug

test: ## Refresh schema.sql and run all tests
	$(TASK) test

download: ## Download local copies of prod DB + event images
	$(TASK) download

kill: ## Kill stray main/templ/air/task processes
	$(TASK) kill

lint: ## Run golangci-lint (external binary)
	@command -v golangci-lint >/dev/null 2>&1 \
		|| { echo "golangci-lint not found; see https://golangci-lint.run/welcome/install/"; exit 1; }
	golangci-lint run

fmt: ## Format all Go code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

tidy: ## Tidy go.mod / go.sum
	go mod tidy

check: generate vet test lint ## Generate, vet, test and lint (pre-commit gate)

clean: ## Remove build artifacts
	rm -rf bin tmp/bin
