PROTOC := protoc
GO := go
DOCKER_COMPOSE := docker-compose
SH := bash

PROTO_PATH := api
GEN_PATH := gen
COMPOSE_FILE := docker-compose.yml

proto:
	@echo "Генерация protobuf файлов..."
	$(PROTOC) -I $(PROTO_PATH) -I third_party \
		--go_out=$(GEN_PATH)/pb --go_opt=paths=source_relative \
    	--go-grpc_out=$(GEN_PATH)/pb --go-grpc_opt=paths=source_relative \
    	--grpc-gateway_out=$(GEN_PATH)/pb --grpc-gateway_opt=paths=source_relative \
    	--openapiv2_out=$(GEN_PATH)/openapi --openapiv2_opt=logtostderr=true \
    	$(PROTO_PATH)/calculator.proto

build:
	@echo "Сборка бинарных файлов..."
	$(GO) build -o server ./cmd/server
	$(GO) build -o client ./cmd/client

test:
	@echo "Запуск юнит-тестов..."
	$(GO) test -v ./internal/calculator/...

grpc-test-sh:
	@echo "Запуск grpc тестов..."
	@if [ -f test-grpc.sh ]; then \
    	$(SH) test-grpc.sh; \
  	else \
   		echo "Файл grpc_tests.sh не найден!"; \
    	exit 1; \
  	fi

http-test-sh:
	@echo "Запуск http тестов..."
	@if [ -f test-http.sh ]; then \
    	$(SH) test-http.sh; \
  	else \
   		echo "Файл grpc_tests.sh не найден!"; \
    	exit 1; \
  	fi

docker-up:
	@echo "Запуск Docker-контейнеров..."
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) up -d

docker-down:
	@echo "Остановка Docker-контейнеров..."
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) down --remove-orphans

docker-logs:
	@echo "Логи Docker-контейнеров..."
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) logs

fmt:
	@echo "Форматирование Go-кода..."
	$(GO) fmt ./...