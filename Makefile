E2E_COMPOSE_FILE := docker-compose.e2e.yaml

help:
	@echo "                    Доступные команды                         "
	@echo ""
	@echo "  Тестирование:"
	@echo "    test-unit              - Запустить юнит-тесты"
	@echo "    test-load              - Запустить нагрузочный тест"
	@echo "    test-integration       - Запустить интеграционные тесты"
	@echo "    test-e2e               - Запустить E2E тесты"
	@echo "    test-all               - Запустить все тесты (unit, integration, e2e)"
	@echo ""
	@echo "  Покрытие:"
	@echo "    coverage-all           - Запустить все тесты и объединить покрытие"
	@echo ""
	@echo "  Запуск и окружение:"
	@echo "    e2e-up                 - Поднять E2E окружение"
	@echo "    e2e-down               - Остановить E2E окружение"
	@echo "    e2e-up-test-down       - Запустить E2E тесты с поднятием/остановкой"
	@echo "    up                     - Развертывание docker-compose для сервиса"
	@echo "    down                   - Завершение работы контейнеров"
	@echo ""
	@echo "  Разработка:"
	@echo "    lint                   - Пролинтить весь код"
	@echo "    mock                   - Сгенерировать моки"
	@echo "    clean                  - Очистить артефакты"
	@echo "    seed                   - Заполнение локальной базы данных тестовыми данными"
	@echo ""
	@echo "  CI/CD:"
	@echo "    ci                     - Запустить CI пайплайн (clean + lint + mock + test-all)"
	@echo ""

lint: ## Пролинтить весь код
	@echo "Linting..."
	@golangci-lint run ./...

mock: ## Сгенерировать моки
	@echo "Mocking..."
	@go generate ./...

test-unit: ## Запустить юнит-тесты
	@echo "Unit-testing..."
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

e2e-up-test-down: ## Запуск e2e тестов с поднятием и выключением окружения
	$(MAKE) e2e-up
	@echo "Running E2E tests..."
	@E2E_API_URL=http://localhost:8081 go test -v -tags=e2e ./test/e2e/...; \
	EXIT_CODE=$$?; \
	$(MAKE) e2e-down; \
	exit $$EXIT_CODE

test-e2e: e2e-up-test-down

test-load: ## Нагрузочный тест
	$(MAKE) e2e-up
	@echo "Running load test..."
	go test -v -tags=load -run TestLoad ./test/load/... -timeout 10m; \
	EXIT_CODE=$$?; \
	$(MAKE) e2e-down; \
	exit $$EXIT_CODE

test-integration: ## Интеграционные тесты
	$(MAKE) e2e-up
	@go test -v -tags=integration -coverprofile=coverage_integration.out \
    		-coverpkg=./... \
    		./test/integration/...;
	@go tool cover -func=coverage_integration.out
	$(MAKE) e2e-down

test-e2e: e2e-up-test-down

coverage-all: ## Запустить все тесты и объединить покрытие
	@echo "Running all tests with coverage..."
	@echo ""
	@echo "=== Unit tests ==="
	@go test -coverprofile=coverage_unit.out ./...
	@echo ""
	@echo "=== Integration tests ==="
	@$(MAKE) e2e-up
	@go test -tags=integration -coverprofile=coverage_integration.out \
		-coverpkg=./... \
		./test/integration/...
	@echo ""
	@$(MAKE) e2e-down
	@echo "=== E2E tests ==="
	@$(MAKE) e2e-up
	@E2E_API_URL=http://localhost:8081 go test -tags=e2e -coverprofile=coverage_e2e.out ./test/e2e/...
	@$(MAKE) e2e-down
	@echo ""
	@echo "=== Merging coverage ==="
	@go install github.com/wadey/gocovmerge@latest
	@go run github.com/wadey/gocovmerge@latest coverage_unit.out coverage_integration.out coverage_e2e.out > coverage_all.out
	@go tool cover -func=coverage_all.out | grep total
	@echo ""
	@echo "Total coverage report: coverage_all.out"

clean: ## Очистка тестовых артефактов
	@echo "Cleaning..."
	@rm -f coverage.out coverage_unit.out coverage_integration.out coverage_e2e.out coverage_all.out coverage.html
	@rm -rf bin/
	@go clean -testcache

test-all: test-unit test-integration test-e2e

ci: clean lint mock test-all

up:
	@docker-compose up -d --build

down:
	@docker-compose down

seed:
	@echo "Seeding local db..."
	@go run cmd/seed/main.go