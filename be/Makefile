.PHONY: up down api worker test test-integration lint migrate-up migrate-down build

up:
	docker compose up -d

down:
	docker compose down

api:
	go run ./cmd/api

worker:
	go run ./cmd/worker

test:
	go test ./...

test-integration:
	go test -tags=integration ./...

lint:
	gofmt -w .
	go vet ./...

migrate-up:
	go run ./cmd/migrate -direction up

migrate-down:
	go run ./cmd/migrate -direction down -steps 1

build:
	go build -o bin/api ./cmd/api
	go build -o bin/worker ./cmd/worker
