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

[x] MILESTONE 4: HTTP Handlers, WebSocket Gateway & Outbox Relay (3 Sessions)
    [x] Session 4.1: Handler Modularization & Response Standardization
    [x] Session 4.2: Real-time Pomodoro WebSocket Gateway (/ws/pomodoro)
    [x] Session 4.3: Transactional Outbox Background Relay & Composition Root

[x] MILESTONE 5: Worker Asynchronous Event Consumers (2 Sessions)
    [x] Session 5.1: RabbitMQ Consumer Framework & Idempotency Engine
    [x] Session 5.2: Async Domain Handlers (Pomodoro, SRS, Budget Alert) & Worker Wiring

[x] MILESTONE 6: Hardening, Integration Testing & Verification (3 Sessions)
    [x] Session 6.1: Unit & Table-Driven Domain Test Suite
    [x] Session 6.2: PostgreSQL Integration Tests & ACID Concurrency Validation
    [x] Session 6.3: Code Quality, Linter, Security Headers & Build Verification

[x] MILESTONE 7: Design System Primitives, Tokens & Layout Foundation (2 Sessions)
    [x] Session 7.1: Double-Bezel, Button-in-Button, AmountDisplay & ProgressRing Primitives
    [x] Session 7.2: Dọn dẹp dead code (features/assets), Domain Constants & Type Reconciliations

[x] MILESTONE 8: Today Asymmetric Bento Dashboard & IAM Settings (2 Sessions)
    [x] Session 8.1: Today Asymmetric Bento Hub (/app/today) với Quick-Action Widgets
    [x] Session 8.2: Settings & User Profile Management (/app/settings) + PATCH /api/v1/me

[x] MILESTONE 9: Deep Work — Kanban Atomic Reordering & Real-time WebSocket Pomodoro (2 Sessions)
    [x] Session 9.1: Kanban Task Board (Urgent pulse, LeetCode chip, Pomodoro estimates & Reorder)
    [x] Session 9.2: Real-time WebSocket Gateway Client (/ws/pomodoro) & History Drawer

[x] MILESTONE 10: Linguistics Hub — 3D Flip Card SRS Engine, Modes & Stats (2 Sessions)
    [x] Session 10.1: Interactive Study Session (/learn/review) với 3D Flip & 4 Study Modes
    [x] Session 10.2: Vocabulary Management (/learn/vocabulary), TTS Audio Playback & 90-day Heatmap

[x] MILESTONE 11: Personal Finance — Budgets Tracker, Savings Goals & Analytics (2 Sessions)
    [x] Session 11.1: Budget Tracker với Alert Thresholds & Category Auto-Seeding (/finance/budgets)
    [x] Session 11.2: Savings Goals với Circular Progress Rings & Net Worth Analytics (/finance/goals)

[x] MILESTONE 12: Journal Enhancements, Command Palette & Polish (2 Sessions)
    [x] Session 12.1: Journal Energy Level, Pinned Notes, Linking UI & Detail View (/journal/:id)
    [x] Session 12.2: Global Command Palette (Ctrl+K), Responsive Mobile Audit & Quality Gates
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

- [x] **[NEW] [`be/internal/adapters/handler/http/ws/pomodoro_ws.go`](file:///g:/GoPA/be/internal/adapters/handler/http/ws/pomodoro_ws.go)**
  - WebSocket upgrader configuration with safe origin checks.
  - Connection handler for `GET /api/v1/ws/pomodoro`.
  - Goroutine per connection reading current Redis state on connection and subscribing to Redis Pub/Sub channel `pomodoro:channel:{user_id}`.
  - Periodic heartbeat ticker (every 30s) and clean connection termination on context cancellation.

---

### Session 4.3: Transactional Outbox Background Relay & Composition Root
**Objective:** Implement the background publisher that relays durable outbox events from PostgreSQL to RabbitMQ exchange `gopa.events`.

- [x] **[NEW] [`be/internal/adapters/outbox/relay.go`](file:///g:/GoPA/be/internal/adapters/outbox/relay.go)**
  - OutboxRelay background daemon with context cancellation.
  - Polls `outbox_events WHERE published_at IS NULL ORDER BY created_at ASC LIMIT 50`.
  - Publishes messages persistently to RabbitMQ topic exchange `gopa.events` with routing key.
  - Atomically marks `published_at = now()` upon successful broker confirmation.
- [x] **[MODIFY] [`be/cmd/api/main.go`](file:///g:/GoPA/be/cmd/api/main.go)**
  - Wire all new handlers, start Outbox Relay in background goroutine, remove legacy `AssetService`, and ensure clean graceful shutdown.

---

## Milestone 5: Worker Asynchronous Event Consumers

### Session 5.1: RabbitMQ Consumer Framework & Idempotency Engine
**Objective:** Establish the consumer topology with dead-letter exchange, bounded retries, and deduplication via `processed_events`.

- [x] **[NEW] [`be/internal/adapters/broker/consumer.go`](file:///g:/GoPA/be/internal/adapters/broker/consumer.go)**
  - Declare durable topic exchange `gopa.events`.
  - Declare queues: `gopa.queue.srs`, `gopa.queue.pomodoro`, `gopa.queue.finance`.
  - Bind queues with matching routing keys (`learning.#`, `pomodoro.#`, `finance.#`).
  - Implement idempotency wrapper checking `processed_events (event_id, processed_at)` before invoking handlers.

---

### Session 5.2: Async Domain Handlers & Worker Wiring
**Objective:** Implement background consumers for spaced repetition updates, pomodoro history archiving, and budget threshold alerts.

- [x] **[NEW] [`be/internal/services/worker_handlers.go`](file:///g:/GoPA/be/internal/services/worker_handlers.go)**
  - `HandlePomodoroFinished`: persists `pomodoro_history` and increments `completed_pomodoros` on linked task.
  - `HandleReviewCompleted`: executes async SRS schedule updates and mastery adjustments.
  - `HandleBudgetAlert`: logs budget utilization warning via `GetBudgetStatus`.
- [x] **[NEW] [`be/internal/services/worker_handlers_test.go`](file:///g:/GoPA/be/internal/services/worker_handlers_test.go)**
  - 9 unit tests across all 3 handlers (happy paths, error cases, input validation).
- [x] **[MODIFY] [`be/cmd/worker/main.go`](file:///g:/GoPA/be/cmd/worker/main.go)**
  - Connects PostgreSQL, Redis, RabbitMQ.
  - Registers all 3 event handlers and starts consumer goroutines.
  - Handles `SIGINT`/`SIGTERM` with graceful `Consumer.Wait()` drain.

---

## Milestone 6: Hardening, Integration Testing & Verification

### Session 6.1: Unit & Table-Driven Domain Test Suite
**Objective:** Maximize unit test coverage on pure domain rules, finance math, and input validation.

- [x] Table-driven tests for SM-2 edge cases (ease factor minimums, repetition resets).
- [x] Unit tests for multi-currency conversion and cashflow calculations.
- [x] Unit tests for token issuance, expiration, and password hashing.
- [x] Run `go test -v -cover ./internal/core/... ./pkg/...`

---

### Session 6.2: PostgreSQL Integration Tests & ACID Concurrency Validation
**Objective:** Validate transaction isolation and database constraints under concurrent load.

- [x] Integration test for `CreateTransaction` verifying atomic account balance updates and rollbacks.
- [x] Integration test for `TransferBetweenAccounts` ensuring atomic double-entry balance consistency.
- [x] Integration test for PostgreSQL full-text search verifying `search_vector` GIN index behavior and `ts_rank`.
- [x] Integration test for Outbox Relay and worker consumer idempotency.

---

### Session 6.3: Code Quality, Linter, Security Headers & Build Verification
**Objective:** Achieve clean static analysis, verify security constraints, and produce production binaries.

- [x] Run `gofmt -w .` across entire repository.
- [x] Run `go vet ./...` (0 warnings).
- [x] Verify CORS headers, token lifetimes, and SQL injection prevention (no unparameterized dynamic queries).
- [x] Build all binaries:
  - `go build -o api.exe ./cmd/api`
  - `go build -o worker.exe ./cmd/worker`
  - `go build -o migrate.exe ./cmd/migrate`
- [x] Verify API health checks: `GET /healthz` and `GET /readyz`.

---

## Milestone 7: Design System Primitives, Tokens & Layout Foundation

### Session 7.1: Double-Bezel, Button-in-Button, AmountDisplay & ProgressRing Primitives
**Objective:** Establish core Agency-Tier UI building blocks in `fe/src/components/design-system/` following `high-end-visual-design` and `ui-styling`.

- [x] **[NEW] `fe/src/components/design-system/double-bezel-card.tsx`**
  - Implement nested hardware enclosure: outer wrapper `rounded-[1.5rem] bg-white/5 border border-white/10 p-1.5 backdrop-blur-xl` wrapping inner core `rounded-[calc(1.5rem-0.375rem)] bg-slate-950/70 border border-white/5`.
  - Add optional dynamic glow shadow property (`glowColor`).
- [x] **[NEW] `fe/src/components/design-system/button-in-button.tsx`**
  - Pill button (`rounded-full px-5 py-2 active:scale-[0.98]`) with trailing icon encased in an inner circular pill wrapper (`w-7 h-7 rounded-full bg-white/15`).
  - Add Framer Motion spring hover translation (`group-hover:translate-x-0.5 group-hover:-translate-y-0.5`).
- [x] **[NEW] `fe/src/components/design-system/progress-ring.tsx`**
  - Circular SVG progress ring with smooth stroke-dashoffset transition, configurable stroke width, track color, and indicator color.
- [x] **[MODIFY] `fe/src/components/design-system/amount-display.tsx`**
  - Ensure strict `tabular-nums` formatting, locale awareness (`VND`, `USD`, `JPY`), and semantic color coding (emerald for positive, rose for negative).

### Session 7.2: Dead Code Cleanup & Domain Constants Reconciliation
**Objective:** Remove legacy asset tracking artifacts and synchronize frontend domain constants with backend database models.

- [x] **[DELETE] `fe/src/features/assets/`**
  - Remove dead directory and obsolete types from vehicle/network tracking.
- [x] **[MODIFY] `fe/src/domain/constants.ts`**
  - Add `TASK_PRIORITY.URGENT = "URGENT"`
  - Add `TASK_CATEGORY.LEETCODE = "LEETCODE"`
  - Add `STUDY_MODE = { FLASHCARD: "FLASHCARD", MULTIPLE_CHOICE: "MULTIPLE_CHOICE", TYPE_IN: "TYPE_IN", AUDIO_QUIZ: "AUDIO_QUIZ" }`
  - Add `BUDGET_PERIOD = { MONTHLY: "MONTHLY", QUARTERLY: "QUARTERLY", YEARLY: "YEARLY" }`
- [x] **[MODIFY] `fe/src/styles/tokens.css`**
  - Align custom properties with the 3-layer architecture (`--radius-outer`, `--radius-inner`, `--ease-spring`).

---

## Milestone 8: Today Asymmetric Bento Dashboard & IAM Settings

### Session 8.1: Today Asymmetric Bento Hub (`/app/today`)
**Objective:** Replace the placeholder Today route with an Asymmetric Bento Grid dashboard aggregating live tasks, Pomodoro status, review queue, and financial KPIs.

- [x] **[NEW] `fe/src/features/today/pages/today-page.tsx`**
  - Zone 1 (Hero Left, `col-span-8`): Dynamic greeting banner ("Chào buổi sáng, Hoang"), top 3 next priority tasks (with priority badges), and "+ Thêm công việc" action.
  - Zone 2 (Hero Right, `col-span-4`): Live Pomodoro widget card with circular progress ring and Start/Pause CTA.
  - Zone 3 (Bottom Left): Daily SRS review queue card ("15 từ cần ôn hôm nay") with Violet gradient button to launch review.
  - Zone 4 (Bottom Center): Personal finance snapshot (Net Worth & Today's spend) with quick "+ Giao dịch" button.
  - Zone 5 (Bottom Right): Journal writing streak counter ("🔥 7 ngày liên tiếp") with link to new entry.
- [x] **[MODIFY] `fe/src/app/router.tsx`**
  - Mount `TodayPage` to `/app/today`.

### Session 8.2: Settings & User Profile Management (`/app/settings`)
**Objective:** Deliver complete IAM profile management hooked into backend `PATCH /api/v1/me`.

- [x] **[MODIFY] `fe/src/features/auth/types.ts` & `fe/src/features/auth/api/auth-api.ts`**
  - Add `display_name`, `avatar_url` to `UserDto`, `User`, and `RegisterInput`.
  - Implement `updateProfile(input: { display_name?: string; avatar_url?: string }): Promise<User>`.
- [x] **[NEW] `fe/src/features/settings/pages/settings-page.tsx`**
  - Profile settings card: avatar display, display name input, email (read-only), role badge.
  - Appearance & Locale card: Theme toggle (Dark/Light), Language switcher (`vi`, `en`, `ja`).
  - Keyboard shortcuts reference drawer / cheatsheet (`Space`, `1-4`, `Ctrl+K`).
- [x] **[MODIFY] `fe/src/app/router.tsx`**
  - Mount `SettingsPage` to `/app/settings`.

---

## Milestone 9: Deep Work — Kanban Atomic Reordering & Real-time WebSocket Pomodoro

### Session 9.1: Kanban Task Board (Urgent Pulse, LeetCode Chip & Reordering)
**Objective:** Upgrade Kanban board with new backend domain fields and status reordering.

- [x] **[MODIFY] `fe/src/features/tasks/types.ts` & `tasks-api.ts`**
  - Add `estimated_pomodoros`, `completed_pomodoros`, `URGENT` priority, `LEETCODE` category.
  - Implement `updateTaskStatus(id: string, status: TaskStatus, targetPosition?: number)`.
- [x] **[MODIFY] `fe/src/features/tasks/pages/tasks-page.tsx`**
  - Upgrade Kanban task cards with `DoubleBezelCard` nested architecture.
  - Add pulsing crimson dot indicator for `URGENT` priority tasks.
  - Add LeetCode category badge with specialized icon (`Code2`).
  - Add Pomodoro count chips (`completed / estimated`).
  - Implement drag-and-drop or move action menu invoking `updateTaskStatus`.

### Session 9.2: Real-time WebSocket Gateway Client (`/ws/pomodoro`) & History
**Objective:** Replace 30-second HTTP polling with persistent WebSocket state synchronization and history tracking.

- [x] **[NEW] `fe/src/features/pomodoro/hooks/use-pomodoro-ws.ts`**
  - Connect to `GET /api/v1/ws/pomodoro` with Bearer auth token query param.
  - Listen for Redis Pub/Sub events (`pomodoro.started`, `pomodoro.paused`, `pomodoro.stopped`, `pomodoro.finished`).
  - Automatic reconnection with exponential backoff and connection status indicator (`connected`, `connecting`, `offline`).
- [x] **[MODIFY] `fe/src/pages/pomodoro-page.tsx`**
  - Integrate `usePomodoroWs` for instant 0ms state changes.
  - Add session history drawer querying `GET /api/v1/pomodoro/history`.
  - Add linked task selector linking the active Pomodoro to a Kanban task.

---

## Milestone 10: Linguistics Hub — 3D Flip Card SRS Engine, Modes & Stats

### Session 10.1: Interactive Study Session (`/learn/review`) với 3D Flip & 4 Modes
**Objective:** Deliver distraction-free SM-2 spaced repetition study session.

- [x] **[MODIFY] `fe/src/features/linguistics/types.ts` & `linguistics-api.ts`**
  - Update `Vocabulary` model with all SM-2 fields (`ease_factor`, `interval_days`, `tags`, `audio_url`, `difficulty_level`, `example_translation`).
  - Implement `getReviewQueue(language: string): Promise<Vocabulary[]>`.
  - Implement `recordReview(id: string, quality: number, studyMode: string, responseTimeMs: number): Promise<void>`.
  - Implement `startLearningSession` and `endLearningSession`.
- [x] **[NEW] `fe/src/features/linguistics/pages/review-session-page.tsx`**
  - Fullscreen distraction-free layout.
  - 3D Card Flip (`rotateY: 180deg`) via Framer Motion with spacebar shortcut.
  - SM-2 Rating buttons: `1` (Again - Red), `2` (Hard - Orange), `3` (Good - Green), `4` (Easy - Blue).
  - Multi-mode support: Flashcard, Multiple Choice (green glow/red shake), Type-in (fuzzy matching), Audio Quiz.
  - Session Summary screen showing accuracy %, XP gained, and mastery score.
- [x] **[MODIFY] `fe/src/app/router.tsx`**
  - Map `/app/learn/review` to `ReviewSessionPage`.

### Session 10.2: Vocabulary Management (`/learn/vocabulary`), TTS & 90-Day Heatmap
**Objective:** Complete vocabulary table view, audio playback, CSV import, and learning statistics.

- [x] **[NEW] `fe/src/features/linguistics/pages/vocabulary-manage-page.tsx`**
  - High-density data table with furigana rendering (`Noto Sans JP`), box indicators (1-5), and tags.
  - Audio TTS playback button calling `POST /api/v1/vocabularies/:id/audio`.
  - CSV Bulk Import modal calling `POST /api/v1/vocabularies/import`.
- [x] **[MODIFY] `fe/src/features/linguistics/pages/learn-page.tsx`**
  - 90-day review heatmap (GitHub-style, violet tint).
  - Leitner box distribution chart.
  - Quick-start study session button (prominent violet gradient pill).
- [x] **[MODIFY] `fe/src/app/router.tsx`**
  - Map `/app/learn/vocabulary` to `VocabularyManagePage`.

---

## Milestone 11: Personal Finance — Budgets Tracker, Savings Goals & Analytics

### Session 11.1: Budget Tracker với Alert Thresholds & Category Seeding (`/finance/budgets`)
**Objective:** Deliver monthly budget management with dynamic progress bars and alert threshold badges.

- [x] **[MODIFY] `fe/src/features/finance/types.ts` & `finance-api.ts`**
  - Define `Budget`, `BudgetStatus`, `BudgetInput`.
  - Implement `listBudgets()`, `createBudget()`, `updateBudget()`, `deleteBudget()`, `getBudgetStatus()`.
- [x] **[NEW] `fe/src/features/finance/pages/budgets-page.tsx`**
  - Budget cards (`DoubleBezelCard`): category icon, budget limit, current expense sum, remaining amount.
  - Color-coded progress meter: Green (<60%), Amber (60-85%), Crimson (>85%).
  - Over-budget warning badge when actual spend exceeds `alert_threshold`.
  - Create/Edit budget modal dialog.
- [x] **[MODIFY] `fe/src/app/router.tsx`**
  - Map `/app/finance/budgets` to `BudgetsPage`.

### Session 11.2: Savings Goals với Circular Progress Rings & Analytics (`/finance/goals`)
**Objective:** Implement savings goal tracking with circular SVG progress indicators and net worth analytics.

- [x] **[MODIFY] `fe/src/features/finance/types.ts` & `finance-api.ts`**
  - Define `SavingsGoal`, `SavingsGoalInput`.
  - Implement `listSavingsGoals()`, `createSavingsGoal()`, `updateSavingsGoal()`, `deleteSavingsGoal()`, `updateGoalProgress()`.
  - Implement `getNetWorthHistory()`.
- [x] **[NEW] `fe/src/features/finance/pages/goals-page.tsx`**
  - Savings goal cards with `ProgressRing`, target date countdown, and linked account balance.
  - Quick deposit / progress adjustment dialog.
- [x] **[MODIFY] `fe/src/app/router.tsx`**
  - Map `/app/finance/goals` to `GoalsPage`.

---

## Milestone 12: Journal Enhancements, Command Palette & Polish

### Session 12.1: Journal Energy Level, Pinned Notes, Linking UI & Detail View (`/journal/:id`)
**Objective:** Upgrade Markdown journal with energy level indicators, note pinning, inter-entry linking, and dedicated reader view.

- [x] **[MODIFY] `fe/src/features/journal/types.ts` & `journal-entry-editor.tsx`**
  - Add `energy_level` (1-5) picker component (lightning bolt icons).
  - Add `pinned` toggle button.
  - Add note linking/unlinking controls (`POST /api/v1/journals/:id/link`).
- [x] **[NEW] `fe/src/features/journal/pages/journal-detail-page.tsx`**
  - Full-page rendered Markdown view with code block syntax highlighting.
  - Backlinks drawer/section displaying connected entries.
  - Edit entry button and quick export.
- [x] **[MODIFY] `fe/src/app/router.tsx`**
  - Map `/app/journal/:journalId` to `JournalDetailPage`.

### Session 12.2: Global Command Palette (`Ctrl+K`), Mobile Audit & Final Quality Gates
**Objective:** Implement system-wide spotlight search, audit mobile responsiveness, and pass all quality gates.

- [x] **[NEW] `fe/src/components/design-system/command-palette.tsx`**
  - Global `Ctrl+K` dialog with instant keyboard search across Tasks, Vocabulary, Journals, and Actions.
- [x] **[TEST] Verification & Build Quality Gates**
  - Run `npm run typecheck` (0 errors).
  - Run `npm run lint` (0 warnings).
  - Run `npm run test` (all test suites pass).
  - Run `npm run build` (clean Vite production bundle).

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
