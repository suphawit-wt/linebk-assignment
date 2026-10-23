SHELL := /bin/bash

CYAN    := \033[36m
GREEN   := \033[32m
YELLOW  := \033[33m
RED     := \033[31m
BOLD    := \033[1m
RESET   := \033[0m

.DEFAULT_GOAL := help

.PHONY: init
init: ## Initialize project (install dependencies, config .env files, generate Vault JWT keys, and check/install k6)
	@echo -e "$(CYAN)[INFO]$(RESET) Installing Node modules and setting up Lefthook hooks..."
	@pnpm install && lefthook install
	
	@echo -e "$(CYAN)[INFO]$(RESET) Configuring environment files..."
	@if [ -f .env.example ] && [ ! -f .env ]; then \
		cp .env.example .env; \
		echo -e "$(GREEN)[SUCCESS]$(RESET) Created root .env from .env.example"; \
	else \
		echo -e "$(YELLOW)[SKIP]$(RESET) Root .env already exists or .env.example missing"; \
	fi
	
	@if [ -f linebk-api/.env.example ] && [ ! -f linebk-api/.env ]; then \
		cp linebk-api/.env.example linebk-api/.env; \
		echo -e "$(GREEN)[SUCCESS]$(RESET) Created linebk-api/.env from .env.example"; \
	else \
		echo -e "$(YELLOW)[SKIP]$(RESET) linebk-api/.env already exists or .env.example missing"; \
	fi
	
	@echo -e "$(CYAN)[INFO]$(RESET) Preparing Vault JWT secrets directory..."
	@if [ ! -f vault/secrets/jwt/user/access_private.pem ]; then \
		echo -e "$(CYAN)[INFO]$(RESET) Generating 2048-bit RSA key pairs for JWT authentication..."; \
		openssl genrsa -out vault/secrets/jwt/user/access_private.pem 2048 2>/dev/null; \
		openssl rsa -in vault/secrets/jwt/user/access_private.pem -pubout -out vault/secrets/jwt/user/access_public.pem 2>/dev/null; \
		openssl genrsa -out vault/secrets/jwt/user/refresh_private.pem 2048 2>/dev/null; \
		openssl rsa -in vault/secrets/jwt/user/refresh_private.pem -pubout -out vault/secrets/jwt/user/refresh_public.pem 2>/dev/null; \
		chmod +x vault/vault-init.sh; \
		echo -e "$(GREEN)[SUCCESS]$(RESET) JWT RSA keys generated successfully under vault/secrets/jwt/user/"; \
	else \
		echo -e "$(YELLOW)[SKIP]$(RESET) JWT keys already exist, skipping generation."; \
	fi

	@echo -e "$(CYAN)[INFO]$(RESET) Checking for k6 performance testing tool..."
	@if ! command -v k6 &> /dev/null; then \
		if command -v brew &> /dev/null; then \
			echo -e "$(CYAN)[INFO]$(RESET) k6 not found. Installing k6 via Homebrew..."; \
			brew install k6; \
			echo -e "$(GREEN)[SUCCESS]$(RESET) k6 installed successfully!"; \
		else \
			echo -e "$(RED)[WARNING]$(RESET) k6 is not installed and Homebrew (brew) is missing. Please install k6 manually (https://k6.io)."; \
		fi \
	else \
		echo -e "$(YELLOW)[SKIP]$(RESET) k6 is already installed."; \
	fi

	@echo -e "$(GREEN)[SUCCESS]$(RESET) Project initialization completed successfully!"


.PHONY: dev-up
dev-up: ## Start development infrastructure containers in background
	@echo -e "$(CYAN)[INFO]$(RESET) Starting development Docker containers..."
	@docker compose -f docker-compose.dev.yml up -d
	@echo -e "$(GREEN)[SUCCESS]$(RESET) Development containers are up and running."

.PHONY: dev-down
dev-down: ## Stop development infrastructure containers
	@echo -e "$(CYAN)[INFO]$(RESET) Stopping development Docker containers..."
	@docker compose -f docker-compose.dev.yml down
	@echo -e "$(GREEN)[SUCCESS]$(RESET) Development containers stopped."

.PHONY: dev-clean
dev-clean: ## Stop dev containers and remove volumes, orphan containers, and local dev images
	@echo -e "$(YELLOW)[WARNING]$(RESET) Cleaning up development containers, volumes, and local images..."
	@docker compose -f docker-compose.dev.yml down -v --rmi local --remove-orphans
	@echo -e "$(GREEN)[SUCCESS]$(RESET) Development environment fully cleaned."

.PHONY: up
up: ## Start production-grade stack containers in background
	@echo -e "$(CYAN)[INFO]$(RESET) Starting full stack containers..."
	@docker compose up -d
	@echo -e "$(GREEN)[SUCCESS]$(RESET) Full stack containers are up and running."

.PHONY: down
down: ## Stop production-grade stack containers
	@echo -e "$(CYAN)[INFO]$(RESET) Stopping full stack containers..."
	@docker compose down
	@echo -e "$(GREEN)[SUCCESS]$(RESET) Full stack containers stopped."

.PHONY: clean
clean: ## Stop containers and remove volumes, orphan containers, and local build images
	@echo -e "$(YELLOW)[WARNING]$(RESET) Cleaning up containers, volumes, and local images..."
	@docker compose down -v --rmi local --remove-orphans
	@echo -e "$(GREEN)[SUCCESS]$(RESET) Environment fully cleaned."

.PHONY: lint
lint: ## Run golangci-lint with auto-fix enabled
	@echo -e "$(CYAN)[INFO]$(RESET) Running golangci-lint..."
	@cd linebk-api && golangci-lint run --fix ./...
	@echo -e "$(GREEN)[SUCCESS]$(RESET) Linting check passed."

.PHONY: test
test: ## Run Go unit tests with coverage reporting
	@echo -e "$(CYAN)[INFO]$(RESET) Running Go unit tests..."
	@cd linebk-api && go test -cover ./...
	@echo -e "$(GREEN)[SUCCESS]$(RESET) Tests executed successfully."

.PHONY: swag
swag: ## Generate Swagger API documentation specs
	@echo -e "$(CYAN)[INFO]$(RESET) Generating Swagger documentation..."
	@swag init -d linebk-api -g cmd/http/main.go -o linebk-api/docs
	@echo -e "$(GREEN)[SUCCESS]$(RESET) Swagger docs generated successfully."

.PHONY: help
help: ## Display available Makefile commands and descriptions
	@echo -e "$(BOLD)Available Commands for LINE BK Development:$(RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(CYAN)%-12s$(RESET) %s\n", $$1, $$2}'