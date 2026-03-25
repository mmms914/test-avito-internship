.PHONY: help lint mock test e2e-down e2e-up e2e-test e2e-test-coverage

help: ## Показать справку
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "%-20s %s\n", $$1, $$2}'

lint: ## Пролинтить весь код
	@golangci-lint run ./...

mock: ## Сгенерировать моки
	@go generate ./...

test: ## Запустить юнит-тесты
	@go test ./...

E2E_COMPOSE_FILE := docker-compose.e2e.yaml

e2e-up: ## Поднятие e2e среды
	docker-compose -f $(E2E_COMPOSE_FILE) up -d
	@echo "Waiting for services to be ready..."
	@sleep 10

e2e-down: ## Остановка e2e среды
	@echo "Stopping E2E environment..."
	docker-compose -f $(E2E_COMPOSE_FILE) down

e2e-test: ## Запуск e2e тестов
	@echo "Running E2E tests..."
	E2E_API_URL=http://localhost:8081 go test -v -tags=e2e ./test/e2e/...

e2e-test-coverage: ## Запуск e2e тестов с проверкой покрытия
	@echo "Running E2E tests with coverage..."
	E2E_API_URL=http://localhost:8081 go test -v -tags=e2e -coverprofile=e2e_coverage.out ./test/e2e/...
	go tool cover -html=e2e_coverage.out -o e2e_coverage.html
	@echo "Coverage report: e2e_coverage.html"