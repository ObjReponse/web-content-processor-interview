.PHONY: help test run build clean lint fmt deps

# Цвета для вывода
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
RESET  := $(shell tput -Txterm sgr0)

## Помощь
help:
	@echo "$(GREEN)🚀 Web Content Processor Interview$(RESET)"
	@echo ""
	@echo "$(YELLOW)Доступные команды:$(RESET)"
	@echo "  make test    - Запустить тесты"
	@echo "  make run     - Запустить пример"
	@echo "  make build   - Собрать бинарный файл"
	@echo "  make lint    - Проверить код"
	@echo "  make fmt     - Форматировать код"
	@echo "  make clean   - Очистить сборки"
	@echo "  make deps    - Установить зависимости"

## Установка зависимостей
deps:
	@echo "$(GREEN)📦 Installing dependencies...$(RESET)"
	go mod download
	go mod tidy

## Запуск тестов
test:
	@echo "$(GREEN)🧪 Running tests...$(RESET)"
	go test -v ./...

## Запуск тестов с покрытием
test-coverage:
	@echo "$(GREEN)🧪 Running tests with coverage...$(RESET)"
	go test -cover ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(YELLOW)📊 Coverage report: coverage.html$(RESET)"

## Запуск примера
run:
	@echo "$(GREEN)🚀 Running example...$(RESET)"
	go run main.go

## Сборка
build:
	@echo "$(GREEN)🔨 Building...$(RESET)"
	go build -o bin/web-processor main.go

## Форматирование кода
fmt:
	@echo "$(GREEN)📝 Formatting code...$(RESET)"
	go fmt ./...

## Линтинг (если установлен golangci-lint)
lint:
	@echo "$(GREEN)🔍 Linting code...$(RESET)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "$(YELLOW)⚠️  golangci-lint not installed, running go vet instead$(RESET)"; \
		go vet ./...; \
	fi

## Очистка
clean:
	@echo "$(GREEN)🧹 Cleaning...$(RESET)"
	rm -rf bin/
	rm -f coverage.out coverage.html

## Запуск в Docker
docker-run:
	@echo "$(GREEN)🐳 Running in Docker...$(RESET)"
	docker-compose up --build

## Остановка Docker
docker-stop:
	@echo "$(GREEN)🛑 Stopping Docker...$(RESET)"
	docker-compose down