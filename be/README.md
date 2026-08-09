# GoPA backend

The Go API and worker are a modular monolith. The API currently includes bootstrap health checks and the IAM foundation: registration, login, rotating HttpOnly refresh cookies, JWT access tokens, logout, protected `/api/v1/me`, Redis-backed rate limiting, and versioned SQL migrations.

## Start

```powershell
Copy-Item .env.example .env
docker compose up -d
go mod tidy
go run ./cmd/migrate -direction up
go run ./cmd/api
```

In another terminal, run `go run ./cmd/worker`. Check `GET http://localhost:8080/healthz` for process liveness and `GET http://localhost:8080/readyz` for PostgreSQL, Redis, and RabbitMQ connectivity. `make migrate-down` reverts one local migration.

## Quality checks

```powershell
gofmt -w .
go vet ./...
go test ./...
go build ./cmd/api ./cmd/worker
```
