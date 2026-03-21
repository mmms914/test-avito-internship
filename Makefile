.PHONY: help lint

help: ## Показать справку
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "%-20s %s\n", $$1, $$2}'

lint: ## Пролинтить весь код
	@golangci-lint run ./...

mock: ## Пролинтить весь код
	@go generate ./...