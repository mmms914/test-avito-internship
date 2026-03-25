.PHONY: help lint mock test e2e-down e2e-up e2e-test e2e-test-coverage all-tests

E2E_COMPOSE_FILE := docker-compose.e2e.yaml

help: ## Показать справку
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "%-20s %s\n", $$1, $$2}'

lint: ## Пролинтить весь код
	@golangci-lint run ./...

mock: ## Сгенерировать моки
	@go generate ./...

unit-test: ## Запустить юнит-тесты
	@go test ./...

e2e-up: ## Поднятие e2e среды
	docker-compose -f $(E2E_COMPOSE_FILE) up -d
	@echo "Waiting for services to be ready..."
	@sleep 10

e2e-down: ## Остановка e2e среды
	@echo "Stopping E2E environment..."
	docker-compose -f $(E2E_COMPOSE_FILE) down

e2e-test: ## Запуск e2e тестов
	@echo "Running E2E tests..."
	@E2E_API_URL=http://localhost:8081 go test -v -tags=e2e ./test/e2e/...;

e2e-up-test-down: ## Запуск e2e тестов
	$(MAKE) e2e-up
	@echo "Running E2E tests..."
	@E2E_API_URL=http://localhost:8081 go test -v -tags=e2e ./test/e2e/...; \
	EXIT_CODE=$$?; \
	$(MAKE) e2e-down; \
	exit $$EXIT_CODE

all-tests: ## Запустить все тесты (юнит и e2e)
	@echo "Running unit tests..."
	$(MAKE) unit-test
	@echo "Running E2E tests..."
	$(MAKE) e2e-up-test-down
