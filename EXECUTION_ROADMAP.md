# GoPA — Master Engineering Execution Roadmap & Session Checklist

> **Target Stack:** Go 1.22+ · PostgreSQL 16+ · Redis 7+ · RabbitMQ 3.13+ · Python 3.12+ (FastAPI) · React 18 (TS)  
> **Source Documents:** [`MASTER_PROMPT.md`](file:///g:/GoPA/MASTER_PROMPT.md) · [`be/BACKEND_MASTER.md`](file:///g:/GoPA/be/BACKEND_MASTER.md) · [`fe/FRONTEND_MASTER.md`](file:///g:/GoPA/fe/FRONTEND_MASTER.md)  
> **Role:** Lead Go Architect & Principal Engineer  

---

## Roadmap Overview & Progress Dashboard

```text
[x] MILESTONE 1: Core Domain, Unified Response, Error Taxonomy & Schema Migrations (3 Sessions)
    [x] Session 1.1: Unified Response Utilities & Domain Error Taxonomy
    [x] Session 1.2: Pure SM-2 Algorithm & Domain Model Reconciliations
    [x] Session 1.3: PostgreSQL Migration 000005 & Legacy Schema Cleanup

[x] MILESTONE 2: Repositories & Data Access Layer (3 Sessions)
    [x] Session 2.1: Ports Modernization & Finance Persistence (Budgets & Goals)
    [x] Session 2.2: Task Kanban Ordering & Journal PostgreSQL Full-Text Search
    [x] Session 2.3: Linguistics SM-2 & Learning Session Repository Implementation

[x] MILESTONE 3: Application Services & Use Cases (3 Sessions)
    [x] Session 3.1: IAM Profile Management & Task Kanban Service Workflows
    [x] Session 3.2: Finance Budgets, Savings Goals & Category Auto-Seeding
    [x] Session 3.3: Linguistics SM-2 Engine, Session Orchestration & Pomodoro Outbox

[ ] MILESTONE 4: HTTP Handlers, WebSocket Gateway & Outbox Relay (3 Sessions)
    [x] Session 4.1: Handler Modularization & Response Standardization
    [ ] Session 4.2: Real-time Pomodoro WebSocket Gateway (/ws/pomodoro)
    [ ] Session 4.3: Transactional Outbox Background Relay & Composition Root

[ ] MILESTONE 5: Worker Asynchronous Event Consumers (2 Sessions)
    [ ] Session 5.1: RabbitMQ Consumer Framework & Idempotency Engine
    [ ] Session 5.2: Async Domain Handlers (Pomodoro, SRS, Budget Alert) & Worker Wiring

[ ] MILESTONE 6: Hardening, Integration Testing & Verification (3 Sessions)
    [ ] Session 6.1: Unit & Table-Driven Domain Test Suite
    [ ] Session 6.2: PostgreSQL Integration Tests & ACID Concurrency Validation
    [ ] Session 6.3: Code Quality, Linter, Security Headers & Build Verification
```

---

## Milestone 1: Core Domain, Unified Response, Error Taxonomy & Schema Migrations

### Session 1.1: Unified Response Utilities & Domain Error Taxonomy
**Objective:** Eliminate ad-hoc response helpers across handlers and establish a single source of truth for JSON envelopes and domain error mapping.

- [x] **[NEW] [`be/pkg/utils/response/response.go`](file:///g:/GoPA/be/pkg/utils/response/response.go)**
  - Define `Envelope`, `Meta`, `Err`, `ErrDetail` structs.
  - Implement `OK(c *gin.Context, data any)` → HTTP 200 `{ "data": ..., "meta": { "request_id": ... } }`.
  - Implement `Created(c *gin.Context, data any)` → HTTP 201.
  - Implement `Paginated(c *gin.Context, data any, cursor string, total *int)` → HTTP 200 with pagination metadata.
  - Implement `Fail(c *gin.Context, err error)` with centralized domain error matching (`errors.Is`):
    - `ErrValidation` → HTTP 400 / 422 (`VALIDATION_ERROR`)
    - `ErrUnauthorized` → HTTP 401 (`UNAUTHORIZED`)
    - `ErrForbidden` → HTTP 403 (`FORBIDDEN`)
    - `ErrNotFound` → HTTP 404 (`NOT_FOUND`)
    - `ErrConflict` → HTTP 409 (`CONFLICT`)
    - `ErrBadGateway` → HTTP 502 (`SERVICE_UNAVAILABLE`)
    - Default fallback → HTTP 500 (`INTERNAL_ERROR`)
- [x] **[MODIFY] [`be/internal/core/domain/errors.go`](file:///g:/GoPA/be/internal/core/domain/errors.go)**
  - Add `ErrForbidden = errors.New("forbidden")`.
  - Add `ErrBadGateway = errors.New("upstream service error")`.
- [x] **[NEW] [`be/pkg/utils/response/response_test.go`](file:///g:/GoPA/be/pkg/utils/response/response_test.go)**
  - Table-driven unit tests verifying HTTP status codes and envelope formatting.

---

### Session 1.2: Pure SM-2 Algorithm & Domain Model Reconciliations
**Objective:** Replace naive review ladders with the pure SuperMemo-2 spaced-repetition algorithm and reconcile domain entities with the Master Spec.

- [x] **[NEW] [`be/internal/core/domain/linguistics/srs.go`](file:///g:/GoPA/be/internal/core/domain/linguistics/srs.go)**
  - Implement `VocabularyState` struct (`CurrentBox`, `EaseFactor`, `IntervalDays`, `RepetitionCount`).
  - Implement `SRSResult` struct (`NextBox`, `EaseFactor`, `IntervalDays`, `RepetitionCount`, `NextReviewAt`, `MasteryDelta`).
  - Implement pure domain function `CalculateNextReview(current VocabularyState, quality int16, now time.Time) SRSResult`:
    - Quality score $q \ge 3$: increment repetitions, recalculate interval ($I_1=1, I_2=6, I_n=\text{round}(I_{n-1} \times EF)$).
    - Quality score $q < 3$: reset repetitions to 0, interval to 1 day.
    - Update ease factor: $EF' = \max(1.3, EF + (0.1 - (5 - q) \times (0.08 + (5 - q) \times 0.02)))$.
- [x] **[NEW] [`be/internal/core/domain/linguistics/srs_test.go`](file:///g:/GoPA/be/internal/core/domain/linguistics/srs_test.go)**
  - Table-driven tests validating all quality scores (0 to 5) and ease factor floor ($1.3$).
- [x] **[MODIFY] [`be/internal/core/domain/user.go`](file:///g:/GoPA/be/internal/core/domain/user.go)**
  - Add `DisplayName string`, `AvatarURL *string` to `User`.
- [x] **[MODIFY] [`be/internal/core/domain/models.go`](file:///g:/GoPA/be/internal/core/domain/models.go)**
  - Add `TaskPriorityUrgent = "URGENT"` and `TaskCategoryLeetcode = "LEETCODE"`.
  - Add `EstimatedPomodoros int`, `CompletedPomodoros int` to `Task`.
  - Reconcile `Vocabulary` struct with `EaseFactor`, `IntervalDays`, `RepetitionCount`, `MasteryScore`, `Tags`, `DifficultyLevel`, `AudioURL`, `ImageURL`, `ExampleTranslation`.
  - Reconcile `VocabularyReview` struct with `StudyMode`, `ResponseTimeMs`, `WasCorrect`.
  - Add `LearningSession` struct (`SessionType`, `ItemsReviewed`, `ItemsCorrect`, `DurationSeconds`, `StartedAt`, `EndedAt`).
  - Reconcile `Journal` struct with `EnergyLevel`, `Pinned`, `WordCount`.
- [x] **[MODIFY] [`be/internal/core/domain/finance.go`](file:///g:/GoPA/be/internal/core/domain/finance.go)**
  - Define domain entities for `Budget` (`Period`, `Amount`, `AlertThreshold`, `IsActive`) and `SavingsGoal` (`TargetAmount`, `CurrentAmount`, `LinkedAccountID`, `TargetDate`).

---

### Session 1.3: PostgreSQL Migration 000005 & Legacy Schema Cleanup
**Objective:** Deliver a reversible, raw SQL migration that aligns database tables with domain definitions and removes dead legacy tables.

- [x] **[NEW] [`be/migrations/000005_spec_reconciliation.up.sql`](file:///g:/GoPA/be/migrations/000005_spec_reconciliation.up.sql)**
  - `ALTER TABLE users ADD COLUMN display_name TEXT NOT NULL DEFAULT '', ADD COLUMN avatar_url TEXT;`
  - `ALTER TABLE tasks` update CHECK constraints to include `'URGENT'` and `'LEETCODE'`; add `estimated_pomodoros INT DEFAULT 0`, `completed_pomodoros INT DEFAULT 0`.
  - `ALTER TABLE vocabularies` add `ease_factor NUMERIC(4,2) DEFAULT 2.50`, `interval_days INT DEFAULT 0`, `repetition_count INT DEFAULT 0`, `mastery_score SMALLINT DEFAULT 0`, `tags TEXT[] DEFAULT '{}'`, `difficulty_level TEXT DEFAULT 'BEGINNER'`, `audio_url TEXT`, `image_url TEXT`, `example_translation TEXT`.
  - `ALTER TABLE vocabulary_reviews` add `study_mode TEXT DEFAULT 'FLASHCARD'`, `response_time_ms INT DEFAULT 0`, `was_correct BOOLEAN DEFAULT TRUE`.
  - `CREATE TABLE learning_sessions (...)` with primary key and foreign keys.
  - `ALTER TABLE journals` add `energy_level SMALLINT CHECK (energy_level BETWEEN 1 AND 5)`, `pinned BOOLEAN DEFAULT FALSE`, `word_count INT DEFAULT 0`.
  - Add indexes:
    - `CREATE INDEX idx_vocabularies_review_queue ON vocabularies(user_id, next_review_at ASC);`
    - `CREATE INDEX idx_learning_sessions_user ON learning_sessions(user_id, started_at DESC);`
    - `CREATE INDEX idx_outbox_pending ON outbox_events(created_at) WHERE published_at IS NULL;`
  - Cleanly drop legacy tables: `DROP TABLE IF EXISTS vehicle_logs, vehicles, network_nodes CASCADE;`
- [x] **[NEW] [`be/migrations/000005_spec_reconciliation.down.sql`](file:///g:/GoPA/be/migrations/000005_spec_reconciliation.down.sql)**
  - Provide safe rollback script.

---

## Milestone 2: Repositories & Data Access Layer

### Session 2.1: Ports Modernization & Finance Persistence (Budgets & Goals)
**Objective:** Define clean port interfaces for all modules and implement repository persistence for Budgets and Savings Goals.

- [x] **[MODIFY] [`be/internal/core/ports/auth.go`](file:///g:/GoPA/be/internal/core/ports/auth.go)**
  - Add `Update(context.Context, domain.User) (domain.User, error)` to `UserRepository`.
- [x] **[MODIFY] [`be/internal/core/ports/finance.go`](file:///g:/GoPA/be/internal/core/ports/finance.go)**
  - Add `BudgetRepository` interface: `ListBudgets`, `GetBudget`, `CreateBudget`, `UpdateBudget`, `DeleteBudget`, `GetBudgetStatus`.
  - Add `SavingsGoalRepository` interface: `ListSavingsGoals`, `GetSavingsGoal`, `CreateSavingsGoal`, `UpdateSavingsGoal`, `DeleteSavingsGoal`, `UpdateGoalProgress`.
- [x] **[MODIFY] [`be/internal/adapters/repository/finance_repository.go`](file:///g:/GoPA/be/internal/adapters/repository/finance_repository.go)**
  - Implement all `BudgetRepository` methods with SQL parameterization and `user_id` scoping.
  - Implement `GetBudgetStatus` calculating actual expense sum vs budget amount and checking `alert_threshold`.
  - Implement all `SavingsGoalRepository` methods.
- [x] **[MODIFY] [`be/internal/adapters/repository/user_repository.go`](file:///g:/GoPA/be/internal/adapters/repository/user_repository.go)**
  - Update `Create`, `FindByID`, `FindByEmail` to scan `display_name` and `avatar_url`.
  - Implement `Update` for user profile changes.

---

### Session 2.2: Task Kanban Ordering & Journal PostgreSQL Full-Text Search
**Objective:** Upgrade Task repository with column reordering and Journal repository with ranked full-text search.

- [x] **[MODIFY] [`be/internal/adapters/repository/task_repository.go`](file:///g:/GoPA/be/internal/adapters/repository/task_repository.go)**
  - Update query columns for `estimated_pomodoros`, `completed_pomodoros`.
  - Implement `UpdateStatus(ctx, userID, id, status, position)` with atomic reordering in `database.WithTx`.
- [x] **[MODIFY] [`be/internal/adapters/repository/journal_repository.go`](file:///g:/GoPA/be/internal/adapters/repository/journal_repository.go)**
  - Update scan for `energy_level`, `pinned`, `word_count`.
  - Implement `Search(ctx context.Context, userID uuid.UUID, query string, filter domain.JournalFilter) ([]domain.Journal, error)` using `search_vector @@ plainto_tsquery('english', $1)` with `ts_rank(search_vector, plainto_tsquery('english', $1)) DESC` ordering.
- [x] **[DELETE] [`be/internal/adapters/repository/vehicle_repository.go`](file:///g:/GoPA/be/internal/adapters/repository/vehicle_repository.go)**
- [x] **[DELETE] [`be/internal/adapters/repository/network_repository.go`](file:///g:/GoPA/be/internal/adapters/repository/network_repository.go)**

---

### Session 2.3: Linguistics SM-2 & Learning Session Repository Implementation
**Objective:** Implement full SM-2 data persistence, bulk import, learning session recording, and mastery statistics.

- [x] **[MODIFY] [`be/internal/adapters/repository/vocabulary_repository.go`](file:///g:/GoPA/be/internal/adapters/repository/vocabulary_repository.go)**
  - Update `Create`, `Update`, `Get`, `List` to scan all SM-2 fields (`ease_factor`, `interval_days`, `repetition_count`, `mastery_score`, `tags`, etc.).
  - Implement `BulkCreate(ctx context.Context, userID uuid.UUID, items []domain.Vocabulary) error` using batch SQL insertion.
  - Refactor `RecordReview` to atomically insert `vocabulary_reviews` and `outbox_events` (`learning.review.completed`) in `database.WithTx`.
  - Implement `CreateLearningSession` & `UpdateLearningSession`.
  - Implement `GetVocabularyStats(ctx context.Context, userID uuid.UUID) (domain.VocabularyStats, error)`.

---

## Milestone 3: Application Services & Use Cases

### Session 3.1: IAM Profile Management & Task Kanban Service Workflows
**Objective:** Build clean use-case layers for User Profile updates and Kanban task status transitions.

- [x] **[MODIFY] [`be/internal/services/auth_service.go`](file:///g:/GoPA/be/internal/services/auth_service.go)**
  - Update `Register` to accept `display_name`.
  - Implement `UpdateProfile(ctx context.Context, userID uuid.UUID, input UpdateProfileInput) (domain.User, error)`.
- [x] **[MODIFY] [`be/internal/services/task_service.go`](file:///g:/GoPA/be/internal/services/task_service.go)**
  - Implement `UpdateStatus(ctx context.Context, userID, taskID uuid.UUID, status domain.TaskStatus, targetPosition *int64) (domain.Task, error)`.
  - Update task validation and default assignments.

---

### Session 3.2: Finance Budgets, Savings Goals & Category Auto-Seeding
**Objective:** Complete finance business logic including budget alerts and default category initialization.

- [x] **[MODIFY] [`be/internal/services/finance_service.go`](file:///g:/GoPA/be/internal/services/finance_service.go)**
  - Implement `CreateBudget`, `UpdateBudget`, `DeleteBudget`, `ListBudgets`, `GetBudgetStatus`.
  - Implement `CreateSavingsGoal`, `UpdateSavingsGoal`, `DeleteSavingsGoal`, `ListSavingsGoals`, `UpdateGoalProgress`.
  - Implement `SeedDefaultCategories(ctx context.Context, userID uuid.UUID) error` to initialize standard categories (Food, Transport, Salary, etc.) on first user setup.
  - In `CreateTransaction`: check active budgets for threshold violations and insert `finance.budget.alert` outbox event if exceeded.

---

### Session 3.3: Linguistics SM-2 Engine, Session Orchestration & Pomodoro Outbox
**Objective:** Wire SM-2 review calculations, learning sessions, and Pomodoro session completion events.

- [x] **[MODIFY] [`be/internal/services/vocabulary_service.go`](file:///g:/GoPA/be/internal/services/vocabulary_service.go)**
  - In `Review`: call pure domain function `linguistics.CalculateNextReview`, persist review record, and emit `learning.review.completed` outbox event.
  - Implement `ImportVocabularies(ctx context.Context, userID uuid.UUID, data []byte, format string) (int, error)`.
  - Implement `StartSession` and `EndSession`.
  - Implement `GetStats` (retrieving review heatmaps and box distributions).
- [x] **[MODIFY] [`be/internal/services/pomodoro_service.go`](file:///g:/GoPA/be/internal/services/pomodoro_service.go)**
  - On `Stop` / timer completion: insert `pomodoro_history` and emit `pomodoro.finished` event atomically.
  - Implement `ListHistory(ctx context.Context, userID uuid.UUID, limit int, cursor string) ([]domain.PomodoroHistory, error)`.
- [x] **[DELETE] [`be/internal/services/asset_service.go`](file:///g:/GoPA/be/internal/services/asset_service.go)**

---

## Milestone 4: HTTP Handlers, WebSocket Gateway & Outbox Relay

### Session 4.1: Handler Modularization & Response Standardization
**Objective:** Split the monolithic `modules.go` into dedicated domain controllers, adopt `pkg/utils/response`, and register all remaining REST endpoints.

- [x] **[NEW] [`be/internal/adapters/handler/http/tasks_handler.go`](file:///g:/GoPA/be/internal/adapters/handler/http/tasks_handler.go)**
  - Endpoints: `GET/POST /api/v1/tasks`, `PATCH/DELETE /api/v1/tasks/:id`, `PATCH /api/v1/tasks/:id/status`.
- [x] **[NEW] [`be/internal/adapters/handler/http/linguistics_handler.go`](file:///g:/GoPA/be/internal/adapters/handler/http/linguistics_handler.go)**
  - Endpoints: `GET/POST /api/v1/vocabularies`, `GET/PATCH/DELETE /api/v1/vocabularies/:id`, `POST /api/v1/vocabularies/import`, `GET /api/v1/vocabularies/review-queue`, `POST /api/v1/vocabularies/:id/reviews`, `POST /api/v1/vocabularies/:id/audio`, `GET /api/v1/vocabularies/stats`, `POST /api/v1/learning-sessions`, `PATCH /api/v1/learning-sessions/:id`.
- [x] **[NEW] [`be/internal/adapters/handler/http/finance_handler.go`](file:///g:/GoPA/be/internal/adapters/handler/http/finance_handler.go)**
  - Endpoints: Accounts, Categories, Transactions, `GET /api/v1/finance/dashboard`, `GET /api/v1/finance/net-worth`, `GET/POST /api/v1/budgets`, `GET/PATCH/DELETE /api/v1/budgets/:id`, `GET /api/v1/budgets/status`, `GET/POST /api/v1/savings-goals`, `PATCH/DELETE /api/v1/savings-goals/:id`.
- [x] **[NEW] [`be/internal/adapters/handler/http/journal_handler.go`](file:///g:/GoPA/be/internal/adapters/handler/http/journal_handler.go)**
  - Endpoints: `GET/POST /api/v1/journals`, `GET/PATCH/DELETE /api/v1/journals/:id`, `POST /api/v1/journals/:id/link`, `DELETE /api/v1/journals/:id/link/:linked_id`, `GET /api/v1/journals/search`, `GET /api/v1/journals/stats`.
- [x] **[NEW] [`be/internal/adapters/handler/http/pomodoro_handler.go`](file:///g:/GoPA/be/internal/adapters/handler/http/pomodoro_handler.go)**
  - Endpoints: `GET/PUT /api/v1/pomodoro`, `POST /api/v1/pomodoro/stop`, `GET /api/v1/pomodoro/history`.
- [x] **[MODIFY] [`be/internal/adapters/handler/http/auth.go`](file:///g:/GoPA/be/internal/adapters/handler/http/auth.go)**
  - Add `PATCH /api/v1/me` profile update endpoint.
  - Refactor all responses to use `response.OK`, `response.Created`, `response.Fail`.
- [x] **[DELETE] [`be/internal/adapters/handler/http/modules.go`](file:///g:/GoPA/be/internal/adapters/handler/http/modules.go)**

---

### Session 4.2: Real-time Pomodoro WebSocket Gateway (`/ws/pomodoro`)
**Objective:** Deliver real-time synchronization between browser clients and Redis timer state over an authenticated WebSocket connection.

- [ ] **[NEW] [`be/internal/adapters/handler/http/ws/pomodoro_ws.go`](file:///g:/GoPA/be/internal/adapters/handler/http/ws/pomodoro_ws.go)**
  - WebSocket upgrader configuration with safe origin checks.
  - Connection handler for `GET /api/v1/ws/pomodoro`.
  - Goroutine per connection reading current Redis state on connection and subscribing to Redis Pub/Sub channel `pomodoro:channel:{user_id}`.
  - Periodic heartbeat ticker (every 30s) and clean connection termination on context cancellation.

---

### Session 4.3: Transactional Outbox Background Relay & Composition Root
**Objective:** Implement the background publisher that relays durable outbox events from PostgreSQL to RabbitMQ exchange `gopa.events`.

- [ ] **[NEW] [`be/internal/adapters/outbox/relay.go`](file:///g:/GoPA/be/internal/adapters/outbox/relay.go)**
  - OutboxRelay background daemon with context cancellation.
  - Polls `outbox_events WHERE published_at IS NULL ORDER BY created_at ASC LIMIT 50`.
  - Publishes messages persistently to RabbitMQ topic exchange `gopa.events` with routing key.
  - Atomically marks `published_at = now()` upon successful broker confirmation.
- [ ] **[MODIFY] [`be/cmd/api/main.go`](file:///g:/GoPA/be/cmd/api/main.go)**
  - Wire all new handlers, start Outbox Relay in background goroutine, remove legacy `AssetService`, and ensure clean graceful shutdown.

---

## Milestone 5: Worker Asynchronous Event Consumers

### Session 5.1: RabbitMQ Consumer Framework & Idempotency Engine
**Objective:** Establish the consumer topology with dead-letter exchange, bounded retries, and deduplication via `processed_events`.

- [ ] **[NEW] [`be/internal/adapters/broker/consumer.go`](file:///g:/GoPA/be/internal/adapters/broker/consumer.go)**
  - Declare durable topic exchange `gopa.events`.
  - Declare queues: `gopa.queue.srs`, `gopa.queue.pomodoro`, `gopa.queue.finance`.
  - Bind queues with matching routing keys (`learning.#`, `pomodoro.#`, `finance.#`).
  - Implement idempotency wrapper checking `processed_events (event_id, processed_at)` before invoking handlers.

---

### Session 5.2: Async Domain Handlers & Worker Wiring
**Objective:** Implement background consumers for spaced repetition updates, pomodoro history archiving, and budget threshold alerts.

- [ ] **[NEW] [`be/internal/services/worker_handlers.go`](file:///g:/GoPA/be/internal/services/worker_handlers.go)**
  - `HandlePomodoroFinished`: persists `pomodoro_history` and increments `completed_pomodoros` on linked task.
  - `HandleReviewCompleted`: executes async SRS schedule updates and mastery adjustments.
  - `HandleBudgetAlert`: logs budget utilization warning and stores notification state.
- [ ] **[MODIFY] [`be/cmd/worker/main.go`](file:///g:/GoPA/be/cmd/worker/main.go)**
  - Connect to PostgreSQL, Redis, RabbitMQ.
  - Register event consumers and handle `SIGINT`/`SIGTERM` with graceful drain.

---

## Milestone 6: Hardening, Integration Testing & Verification

### Session 6.1: Unit & Table-Driven Domain Test Suite
**Objective:** Maximize unit test coverage on pure domain rules, finance math, and input validation.

- [ ] Table-driven tests for SM-2 edge cases (ease factor minimums, repetition resets).
- [ ] Unit tests for multi-currency conversion and cashflow calculations.
- [ ] Unit tests for token issuance, expiration, and password hashing.
- [ ] Run `go test -v -cover ./internal/core/... ./pkg/...`

---

### Session 6.2: PostgreSQL Integration Tests & ACID Concurrency Validation
**Objective:** Validate transaction isolation and database constraints under concurrent load.

- [ ] Integration test for `CreateTransaction` verifying atomic account balance updates and rollbacks.
- [ ] Integration test for `TransferBetweenAccounts` ensuring atomic double-entry balance consistency.
- [ ] Integration test for PostgreSQL full-text search verifying `search_vector` GIN index behavior and `ts_rank`.
- [ ] Integration test for Outbox Relay and worker consumer idempotency.

---

### Session 6.3: Code Quality, Linter, Security Headers & Build Verification
**Objective:** Achieve clean static analysis, verify security constraints, and produce production binaries.

- [ ] Run `gofmt -w .` across entire repository.
- [ ] Run `go vet ./...` (0 warnings).
- [ ] Verify CORS headers, token lifetimes, and SQL injection prevention (no unparameterized dynamic queries).
- [ ] Build all binaries:
  - `go build -o api.exe ./cmd/api`
  - `go build -o worker.exe ./cmd/worker`
  - `go build -o migrate.exe ./cmd/migrate`
- [ ] Verify API health checks: `GET /healthz` and `GET /readyz`.

---

## Quick Reference: Commands & Make Targets

```bash
# Apply migrations
migrate -path migrations -database "postgres://gopa:gopa@localhost:5432/gopa?sslmode=disable" up

# Run all unit and integration tests
go test -v -race ./...

# Run the API server locally
go run ./cmd/api

# Run the async worker process
go run ./cmd/worker
```
