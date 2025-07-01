APP_NAME     := tfinspector
CMD_PATH     := ./main.go
BIN_DIR      := bin
DOCKER_IMAGE ?= $(APP_NAME)
DOCKER_TAG   ?= latest
DOCKERFILE   ?= Dockerfile
GO_FILES     := $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all fmt vet test build clean docker help install

all: fmt vet test build ## Run all primary steps (format, vet, test, build)

fmt: ## Format Go source files
	@echo "Formatting Go files..."
	@fmt_output=$$(go fmt ./...) ; \
	if [ -n "$$fmt_output" ]; then \
		echo "Formatted files:" ; \
		echo "$$fmt_output" | sed 's/^/ - /' ; \
	else \
		echo "All files are already formatted."; \
	fi

vet: ## Run go vet for static analysis
	@echo "Running go vet..."
	@vet_output=$$(go vet ./... 2>&1) ; \
	if [ -n "$$vet_output" ]; then \
		echo "Issues found:" ; \
		echo "$$vet_output" | sed 's/^/ - /' ; \
		exit 1 ; \
	else \
		echo "No issues found."; \
	fi

test: ## Run unit tests with verbose output
	@echo "Running tests..."
	@output=$$(go test -v ./...) ; \
	status=$$? ; \
	echo "$$output" ; \
	if [ $$status -ne 0 ]; then \
		echo "Tests failed."; \
		exit $$status ; \
	else \
		echo "All tests passed."; \
	fi

build: ## Build a static binary
	@echo "Building static binary..."
	@mkdir -p $(BIN_DIR)
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -trimpath -ldflags="-s -w" \
	-o $(BIN_DIR)/$(APP_NAME) $(CMD_PATH)
	@echo "Binary built: $(BIN_DIR)/$(APP_NAME)"

install: build ## Install the binary into your PATH
	@echo "Installing $(APP_NAME)..."

	@if [ -n "$$GOBIN" ]; then \
		install -m 0755 $(BIN_DIR)/$(APP_NAME) "$$GOBIN/"; \
		echo "Installed to $$GOBIN/$(APP_NAME)"; \
	elif [ -d "$$HOME/go/bin" ]; then \
		install -m 0755 $(BIN_DIR)/$(APP_NAME) "$$HOME/go/bin/"; \
		echo "Installed to $$HOME/go/bin/$(APP_NAME)"; \
	elif [ -d "/usr/local/bin" ]; then \
		echo "Using sudo to install to /usr/local/bin..."; \
		sudo install -m 0755 $(BIN_DIR)/$(APP_NAME) /usr/local/bin/; \
		echo "Installed to /usr/local/bin/$(APP_NAME)"; \
	else \
		echo "Could not determine install location."; \
		echo "Please manually copy $(BIN_DIR)/$(APP_NAME) to a directory in your PATH."; \
		exit 1; \
	fi

	@echo "You can now run '$(APP_NAME)' from anywhere!"

clean: ## Remove build artifacts
	@echo "Cleaning up build artifacts..."
	@if [ -d "$(BIN_DIR)" ]; then \
		rm -rf "$(BIN_DIR)"; \
		echo "Removed $(BIN_DIR)/"; \
	else \
		echo "Nothing to clean."; \
	fi

docker: ## Build Docker image
	@echo "Building Docker image $(DOCKER_IMAGE):$(DOCKER_TAG)..."
	@docker build --file $(DOCKERFILE) --tag $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)"

help: ## Show available Makefile commands
	@echo ""
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'