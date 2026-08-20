# GoPA Backend — Master Architecture Specification & Engineering Playbook

> **Status:** Living specification — update when any architectural decision changes.
> **Stack:** Go 1.22+ · PostgreSQL 16+ · Elasticsearch 8+ · Redis 7+ · RabbitMQ 3.13+ · Docker Compose
> **Related:** `../MASTER_PROMPT.md` (root spec) · `../ai/` (Python AI microservice)

---

## 1. Mission and engineering role

GoPA backend is a **modular monolith**: one HTTP API process and one asynchronous worker process sharing a well-defined domain model. A separate Python AI microservice handles TTS, semantic embedding, and vocabulary enrichment.

The AI acts as an **Elite Principal Go Engineer** and mentor for a developer transitioning from Java/Spring Boot. Every implementation must be idiomatic Go, explained through the meaningful contrast with Java/Spring patterns.

### 1.1 Core principles

1. **Explicit over implicit.** Wiring, transactions, error paths, and concurrency are visible in code.
2. **Domain before transport.** Gin, sqlx, Redis, and RabbitMQ remain at the system edge.
3. **PostgreSQL is the durable truth.** Redis accelerates; it does not replace durable records.
4. **ACID compliance.** Multi-record business changes use explicit `*sqlx.Tx` transactions.
5. **Safe retries.** Every async consumer is idempotent.
6. **Context-aware operations.** Cancellation and deadlines propagate to all I/O calls.
7. **Secure by default.** Ownership enforced in SQL; secrets never in source or logs.
8. **Minimize boilerplate.** Shared helpers for transactions, responses, and error mapping — not duplicated per feature.

---

## 2. Architecture overview

```mermaid
flowchart LR
  Client["React client"] -->|"REST /api/v1"| API["Go API · Gin"]
  Client -->|"WebSocket"| WS["Pomodoro handler"]
  API --> PG[("PostgreSQL")]
  API --> Redis[("Redis")]
  API --> RMQ[("RabbitMQ")]
  API -->|"Internal HTTP"| AISvc["Python AI Service"]
  WS --> Redis
  RMQ --> Worker["Go worker"]
  Worker --> PG
  Worker --> Redis
  AISvc --> ES[("Elasticsearch")]
  AISvc --> PG
```

### 2.1 Process responsibilities

| Process | Owns | Must not do |
|---|---|---|
| API | Validation, auth, use-case execution, transactional writes, event publishing intent | Long-running background processing |
| Worker | Event consumption, idempotent async work (SRS schedule, Pomodoro history, budget alerts) | Serve browser traffic |
| AI Service | TTS audio generation, sentence embeddings, vocabulary enrichment | Business auth/access control |
| PostgreSQL | Business records, outbox, full-text vectors, audit history | Ephemeral/presence state |
| Redis | Sessions, rate limits, live Pomodoro state, rate counters | Sole durable record of business events |
| RabbitMQ | Durable async decoupling | Business source of truth |

### 2.2 Clean Architecture boundary

```text
HTTP handlers / repositories / cache / broker / AI client adapters
                           │
                           ▼
                  ports (interfaces)
                           ▲
                           │
               application services (use cases)
                           │
                           ▼
                 domain entities and pure rules
```

- **Domain:** Pure Go structs, enums, value objects, domain errors, pure scheduling rules.
- **Ports:** Contracts (interfaces) that services need from repos, cache, broker, clock, AI client.
- **Services:** Use cases that coordinate domain behavior through ports.
- **Adapters:** Gin handlers, sqlx repos, Redis cache, RabbitMQ broker, AI HTTP client.
- **Composition roots:** `cmd/api/main.go` and `cmd/worker/main.go` — the only places that wire concrete dependencies.

> **Java/Spring contrast:** Spring's IoC container wires beans through annotations and reflection. In Go, the composition root is plain Go code: readable, debuggable, and refactorable without framework magic.

---

## 3. Repository layout

```text
be/
├── cmd/
│   ├── api/main.go                      # API composition root, graceful shutdown
│   └── worker/main.go                   # Worker composition root, graceful shutdown
├── internal/
│   ├── core/
│   │   ├── domain/                      # Entities, enums, value objects, domain errors
│   │   └── ports/                       # Interfaces required by services
│   ├── adapters/
│   │   ├── handler/http/                # Gin routes, handlers, request/response DTOs
│   │   │   ├── auth/
│   │   │   ├── linguistics/
│   │   │   ├── finance/
│   │   │   ├── deepwork/
│   │   │   ├── journal/
│   │   │   └── middleware/
│   │   ├── repository/                  # PostgreSQL + sqlx implementations
│   │   ├── cache/                       # Redis implementations
│   │   ├── broker/                      # RabbitMQ publisher/consumer
│   │   └── external/
│   │       └── ai_client/               # HTTP client for Python AI service
│   └── services/                        # Application use cases
│       ├── auth/
│       ├── linguistics/
│       ├── finance/
│       ├── deepwork/
│       └── journal/
├── pkg/
│   ├── config/                          # Env loading and validation
│   ├── database/                        # PG pool, WithTx helper, migration runner
│   ├── logger/                          # slog JSON setup
│   └── utils/
│       ├── response/                    # Unified HTTP response/error helpers
│       ├── jwt/                         # Token issuance and validation
│       ├── password/                    # bcrypt helpers
│       └── pagination/                  # Cursor pagination helpers
├── migrations/                          # 000001_init.up.sql, 000001_init.down.sql ...
├── deployments/                         # Docker assets if needed
├── docker-compose.yml
├── Makefile
├── .env.example
├── go.mod
└── README.md
```

### 3.1 File placement rules

- Database struct tags and SQL belong in repository adapters, not handlers.
- SRS scheduling formula belongs in `internal/core/domain/` or a domain service — not the worker adapter.
- HTTP DTOs live beside their handlers and do not leak into repositories.
- Config is loaded once at startup; do not read environment variables inside business code.
- Interfaces live in `internal/core/ports/` — next to the consumer, not the implementation.

### 3.2 Anti-over-abstraction

Do **not** create a `BaseRepository[T]`, a universal generic CRUD service, or an interface for every concrete type. Create abstractions only where:

1. A core use case depends on a replaceable external capability (repo, cache, broker, AI client), OR
2. A clear test seam is needed for isolated unit testing.

### 3.3 Constant Architecture & Conventions

GoPA Backend organizes constants into three clear layers based on visibility and scope:

1. **Public / Global Non-Domain Constants (`be/pkg/constants/constants.go`)**:
   - Contains application environment modes (`Production`, `Development`, `Test`, `Debug`), primitive utility defaults (`ONE_STRING`, `ZERO_STRING`, `ONE_INT32`, `ZERO_INT32`), HTTP headers, cookie names (`RefreshCookieName`), and public configuration default constants (`DefaultBcryptCost`, `DefaultPageSize`).
   - Accessible by all packages across the backend repository.
2. **Domain Business Constants (`be/internal/core/domain/`)**:
   - Contains domain entity types and enums (`TaskStatus`, `TaskPriority`, `UserRole`, `JournalMood`, `TransactionType`, `AccountType`).
   - Pure business domain types live alongside domain entities in `user.go`, `models.go`, `finance.go`.
3. **Internal Non-Domain Constants (`be/internal/constants/constants.go`)**:
   - Contains unexported/internal infrastructure constants (internal cache key prefixes, internal queue/exchange names, internal context keys).
   - Enforces Go package encapsulation so external modules cannot import internal implementation details.

---

## 4. Technology decisions

| Concern | Choice | Reason |
|---|---|---|
| Go | Go 1.22+ | Modern stdlib, simple concurrency, `errors.Is/As` |
| HTTP | `gin-gonic/gin` | Productive routing/middleware; isolated to handler adapters |
| Database | PostgreSQL 16+ | Relational integrity, JSONB, `tsvector` full-text search, `gen_random_uuid()` |
| DB access | `jmoiron/sqlx` + `pgx` driver | SQL-first, struct scanning, explicit queries |
| Search | Elasticsearch 8+ via AI service | Journal semantic search, vocabulary search |
| Cache | `redis/go-redis/v9` | Context-aware client |
| Broker | RabbitMQ + `amqp091-go` | Durable queues, at-least-once delivery |
| Auth | JWT + bcrypt + opaque refresh | Short access lifetime, server-controlled revocation |
| Logging | `log/slog` | Structured JSON; no external dep |
| Config | `joho/godotenv` + custom validator | 12-factor compatible |
| Migrations | `golang-migrate` | Raw SQL, version-controlled, reversible |

---

## 5. Shared backend conventions

### 5.1 Go coding standards

1. `gofmt -w .` and `go vet ./...` must pass on every commit.
2. Never `panic` for expected application errors.
3. Every repository, service, and cache method takes `ctx context.Context` as first argument.
4. Wrap errors with context: `fmt.Errorf("create vocabulary: %w", err)`.
5. Use `errors.Is` / `errors.As`; never string-match error messages for control flow.
6. Prefer small, focused structs and constructor functions (`NewX(deps...) *X`).
7. Both JSON and DB struct tags on persistence/transport structs: `json:"user_id" db:"user_id"`.
8. UTC `time.Time` always; RFC 3339 in JSON.
9. `NUMERIC(14,2)` for money; use a decimal-compatible type (e.g. `shopspring/decimal` or integer cents), never `float64`.
10. Never log: passwords, tokens, auth headers, journal content, PII beyond operational necessity.

### 5.2 Error taxonomy

```go
// internal/core/domain/errors.go
var (
    ErrNotFound     = errors.New("not found")
    ErrConflict     = errors.New("conflict")
    ErrUnauthorized = errors.New("unauthorized")
    ErrForbidden    = errors.New("forbidden")
    ErrValidation   = errors.New("validation failed")
    ErrBadGateway   = errors.New("upstream service error")
)
```

HTTP mapping:

| Domain error | HTTP status | API code |
|---|---|---|
| ErrValidation | 400 / 422 | `VALIDATION_ERROR` |
| ErrUnauthorized | 401 | `UNAUTHORIZED` |
| ErrForbidden | 403 | `FORBIDDEN` |
| ErrNotFound | 404 | `NOT_FOUND` |
| ErrConflict | 409 | `CONFLICT` |
| ErrBadGateway | 502 | `SERVICE_UNAVAILABLE` |
| other | 500 | `INTERNAL_ERROR` |

### 5.3 Unified response helpers

```go
// pkg/utils/response/response.go

type Envelope struct {
    Data  any   `json:"data,omitempty"`
    Error *Err  `json:"error,omitempty"`
    Meta  Meta  `json:"meta"`
}

type Meta struct {
    RequestID  string `json:"request_id"`
    NextCursor string `json:"next_cursor,omitempty"`
    Total      *int   `json:"total,omitempty"`
}

type Err struct {
    Code    string   `json:"code"`
    Message string   `json:"message"`
    Details []ErrDetail `json:"details,omitempty"`
}

func OK(c *gin.Context, data any)
func Created(c *gin.Context, data any)
func Paginated(c *gin.Context, data any, cursor string, total int)
func Fail(c *gin.Context, err error) // maps domain errors → HTTP status + code
```

### 5.4 Handler discipline

A handler does **only** these five things:

1. Bind and validate HTTP DTO (`ShouldBindJSON`, validator tags).
2. Extract authenticated identity from `gin.Context`.
3. Call one service method.
4. Map known domain errors to response via `response.Fail`.
5. Write success response via `response.OK` / `response.Created`.

Handlers must not: build SQL, compute SRS intervals, encode broker retry policy, or duplicate authorization logic.

> **Java/Spring contrast:** A Spring `@RestController` method has similar responsibilities. The difference is explicit error mapping — Go has no `@ExceptionHandler`, so the handler calls `response.Fail(c, err)` which does the mapping centrally.

---

## 6. Common abstractions

### 6.1 Transaction helper

```go
// pkg/database/tx.go
func WithTx(ctx context.Context, db *sqlx.DB, fn func(*sqlx.Tx) error) error {
    tx, err := db.BeginTxx(ctx, nil)
    if err != nil {
        return fmt.Errorf("begin transaction: %w", err)
    }
    defer func() {
        if p := recover(); p != nil {
            _ = tx.Rollback()
            panic(p)
        }
    }()
    if err := fn(tx); err != nil {
        _ = tx.Rollback()
        return err
    }
    return tx.Commit()
}
```

> **Java/Spring contrast:** `@Transactional` uses a proxy to open/commit/rollback automatically. `WithTx` makes transaction scope visible at the call site — the scope is exactly the closure, not "everything until the method returns through all proxy layers."

### 6.2 Cursor pagination helper

```go
// pkg/utils/pagination/cursor.go
type Page struct {
    Limit  int
    Cursor string // opaque base64 encoded timestamp+id
}

func (p Page) Validate() error { ... }
func EncodeCursor(t time.Time, id uuid.UUID) string { ... }
func DecodeCursor(s string) (time.Time, uuid.UUID, error) { ... }
```

### 6.3 AI service client (adapter)

```go
// internal/adapters/external/ai_client/client.go
type AIClient interface {
    GenerateTTS(ctx context.Context, word, language string) (audioURL string, err error)
    EnrichVocabulary(ctx context.Context, word, language string) (*EnrichResult, error)
    SemanticSearch(ctx context.Context, query string, limit int) ([]SearchHit, error)
}
```

Implementation calls the Python AI microservice over internal HTTP. This interface belongs in `ports/` so the service layer depends on the contract, not the HTTP call.

---

## 7. Persistence and migration policy

### 7.1 Universal data conventions

- `gen_random_uuid()` primary keys or application-generated UUIDs.
- Every user-owned row has `user_id UUID NOT NULL REFERENCES users(id)`.
- `TIMESTAMPTZ` everywhere; UTC stored and returned.
- `NUMERIC(14,2)` for monetary values — no `float`.
- `CHECK` constraints and foreign keys for invariants that must survive app bugs.
- Explicit column lists in all queries; no `SELECT *`.
- Parameterized queries always; sqlx does not eliminate injection risk on dynamic column names.

### 7.2 Migration practice

- Raw SQL pairs: `000001_init.up.sql` / `000001_init.down.sql`.
- Never rewrite a migration applied outside a disposable database.
- Every `up` migration has a safe `down` unless destruction is irreversible by design.
- Add indexes for proven query paths (see per-module specs below).
- Test migrations against a blank database and from the prior schema.

### 7.3 Complete table set

```text
-- IAM
users
refresh_sessions

-- Linguistics
vocabularies
vocabulary_reviews
learning_sessions

-- Finance
accounts
categories
transactions
budgets
savings_goals

-- Deep Work
tasks
pomodoro_history

-- Journal
journals
journal_links

-- Infrastructure
outbox_events
processed_events
```

### 7.4 Key indexes

```sql
-- Linguistics: due review queue (most critical for SRS performance)
CREATE INDEX idx_vocabularies_review_queue
  ON vocabularies(user_id, next_review_at ASC)
  WHERE next_review_at IS NOT NULL;

-- Finance: transaction queries
CREATE INDEX idx_transactions_user_occurred
  ON transactions(user_id, occurred_at DESC);
CREATE INDEX idx_transactions_account
  ON transactions(account_id, occurred_at DESC);

-- Journal: full-text search
CREATE INDEX idx_journals_search_vector
  ON journals USING GIN(search_vector);
CREATE INDEX idx_journals_user_published
  ON journals(user_id, published_date DESC);

-- Tasks: board ordering
CREATE INDEX idx_tasks_board
  ON tasks(user_id, status, position);
```

---

## 8. Module specifications

### 8.1 IAM — Identity and access management

**Goal:** Stateless short-lived access tokens with server-controlled refresh revocation.

**User table:** See MASTER_PROMPT.md §6.2.

**Authentication lifecycle:**

1. Registration: canonicalize/validate email → bcrypt hash → INSERT user.
2. Login: verify bcrypt → issue 15-min JWT access token + 7-day opaque refresh token.
3. Store refresh token hash in Redis: `SET auth:refresh:{user_id}:{session_id} {hash} EX 604800`.
4. Refresh: validate token → atomic DELETE old + SET new in Redis → issue new pair.
5. Logout: DELETE refresh session key. Access token expires naturally.
6. Middleware: validate JWT signature + claims → typed identity into `gin.Context`.

**Security rules:**

- bcrypt cost ≥ 12.
- JWT claims: `sub` (user UUID), `role`, `sid` (session ID), `iat`, `exp`.
- Generic error messages for credential failures (no account enumeration).
- Rate limit: 10 login attempts / minute / IP via Redis sorted set or sliding window.

**Routes:**

```text
POST /api/v1/auth/register
POST /api/v1/auth/login
POST /api/v1/auth/refresh
POST /api/v1/auth/logout
GET  /api/v1/me
PATCH /api/v1/me
```

---

### 8.2 Linguistics — Smart Spaced Repetition Engine

**Core domain rule (SM-2 algorithm):** Lives in `internal/core/domain/linguistics/srs.go`. Must be a pure function with no I/O — testable without any infrastructure.

```go
// SRSResult contains the next state after a review
type SRSResult struct {
    NextBox         int
    EaseFactor      float64
    IntervalDays    int
    RepetitionCount int
    NextReviewAt    time.Time
    MasteryDelta    int // +/- change to mastery score
}

// CalculateNextReview is a pure domain function — no I/O
func CalculateNextReview(current VocabularyState, quality int, now time.Time) SRSResult
```

**Review workflow (service layer):**

```go
// In a single transaction:
// 1. Lock the vocabulary row (SELECT ... FOR UPDATE)
// 2. Validate ownership and current state
// 3. INSERT vocabulary_review record
// 4. INSERT outbox_event { event_type: "learning.review.completed", payload: { vocab_id, user_id, quality, ... } }
// 5. COMMIT

// Worker consumes event idempotently:
// 1. Check processed_events for this event_id (INSERT OR IGNORE)
// 2. Apply CalculateNextReview domain function
// 3. UPDATE vocabularies SET ease_factor, interval_days, current_box, next_review_at, mastery_score
```

**AI enrichment flow (on vocabulary creation):**

```go
// After inserting vocabulary, call AI service asynchronously:
// 1. Publish outbox event: "learning.vocabulary.created"
// 2. Worker calls AIClient.EnrichVocabulary(word, language)
// 3. UPDATE vocabularies SET reading, example_sentence, example_translation, audio_url
```

**TTS on demand:**

```
POST /api/v1/vocabularies/{id}/audio
→ Check if audio_url already cached
→ If not: call AIClient.GenerateTTS(word, language)
→ UPDATE vocabularies SET audio_url = ...
→ Return { audio_url }
```

**Study modes (API support):**

Review queue accepts `?mode=FLASHCARD|MULTIPLE_CHOICE|TYPE_IN|AUDIO_QUIZ`:

- `MULTIPLE_CHOICE`: API returns the target word + 3 distractors (same language, similar difficulty).
- `AUDIO_QUIZ`: API returns an audio_url instead of the word text on the front.
- `TYPE_IN`: API accepts the user's typed answer, performs server-side fuzzy match.

**Fuzzy match logic for Type-in (server side):**

```go
// For JP: accept kana or romaji equivalent via pykakasi mapping
// For EN: Levenshtein distance ≤ 1 for words < 6 chars; ≤ 2 for longer
// Return { correct: bool, expected: "...", match_quality: float }
```

**Routes:**

```text
GET    /api/v1/vocabularies              # ?lang=JP&box=1&q=...&limit=&cursor=
POST   /api/v1/vocabularies              # Create + trigger enrichment event
GET    /api/v1/vocabularies/{id}
PATCH  /api/v1/vocabularies/{id}
DELETE /api/v1/vocabularies/{id}
POST   /api/v1/vocabularies/import       # Bulk CSV/JSON
GET    /api/v1/vocabularies/review-queue # ?language=JP&mode=FLASHCARD&limit=20
POST   /api/v1/vocabularies/{id}/reviews
POST   /api/v1/vocabularies/{id}/audio
GET    /api/v1/vocabularies/stats        # Progress, heatmap, streak
POST   /api/v1/learning-sessions         # { language, session_type }
PATCH  /api/v1/learning-sessions/{id}    # End session { duration, correct, reviewed }
```

**Worker events:**

- `learning.review.completed` → Apply SM-2 update idempotently.
- `learning.vocabulary.created` → Enrich via AI service.
- `learning.mastery.upgraded` → (future) push notification or achievement.

---

### 8.3 Personal Finance

**Domain rules:**

- Transaction amount is always positive; direction is determined by `type` (`INCOME`, `EXPENSE`, `TRANSFER`).
- TRANSFER must have both `account_id` and `to_account_id`; category is null.
- Balance update: After INSERT transaction, UPDATE account balance atomically in the same transaction.
- Budget alert: After a transaction exceeds the alert threshold, publish `finance.budget.alert` event.
- Multi-currency: `amount × exchange_rate` converts to the account's base currency for balance updates.

**Balance integrity (critical):**

```sql
-- account.current_balance is maintained via explicit UPDATE in the same transaction as INSERT transaction
-- Never compute balance by summing all transactions at query time for display (too slow)
-- Run a reconciliation job to detect drift periodically
```

**Finance service operations:**

```go
// All these operate within WithTx:
CreateTransaction(ctx, tx, input) → validates category ownership, updates account balance, optionally checks budget alert
TransferBetweenAccounts(ctx, tx, input) → two balance updates in one transaction (ACID)
CreateBudget(ctx, tx, input) → validates category if provided
GetCashflowSummary(ctx, userID, from, to) → aggregation query
GetSpendingByCategory(ctx, userID, from, to) → aggregation query
GetNetWorth(ctx, userID) → sum of all non-archived account balances
```

**Routes:**

```text
GET/POST        /api/v1/accounts
GET/PATCH/DELETE /api/v1/accounts/{id}
GET             /api/v1/accounts/{id}/balance-history   # ?period=daily|weekly|monthly

GET/POST        /api/v1/categories
GET/PATCH/DELETE /api/v1/categories/{id}

GET/POST        /api/v1/transactions                    # ?account=&category=&type=&from=&to=&q=&limit=&cursor=
GET/PATCH/DELETE /api/v1/transactions/{id}
GET             /api/v1/transactions/summary             # Cashflow summary for period
GET             /api/v1/transactions/by-category         # Breakdown ?from=&to=

GET/POST        /api/v1/budgets
GET/PATCH/DELETE /api/v1/budgets/{id}
GET             /api/v1/budgets/status                   # All budgets with % used + alert state

GET/POST        /api/v1/savings-goals
PATCH           /api/v1/savings-goals/{id}

GET             /api/v1/finance/net-worth
GET             /api/v1/finance/dashboard                # All summary data in one call
```

**Worker events:**

- `finance.budget.alert` → Log + (future) push notification.
- `finance.recurring.due` → (future) auto-create recurring transactions.

---

### 8.4 Deep Work — Tasks and Pomodoro

**Task domain rules:**

- Status transitions validated in service (even when all transitions initially allowed).
- `position` ordering: scoped by `(user_id, status)`. Reorder in one transaction with an update-all-affected approach.
- `estimated_pomodoros` and `completed_pomodoros` are informational; no derived status from them.

**Pomodoro live state:**

Redis key: `pomodoro:state:{user_id}` with TTL of 4 hours max.

```json
{
  "status": "running|paused|idle|break",
  "started_at": "2026-08-03T09:00:00Z",
  "duration_seconds": 1500,
  "break_duration_seconds": 300,
  "task_id": "optional-uuid",
  "session_count": 2,
  "updated_at": "2026-08-03T09:00:00Z"
}
```

**WebSocket handler:**

- One goroutine per connection; canceled when connection closes.
- Reads authoritative Redis state and broadcasts to the connected browser.
- Also subscribes to Redis Pub/Sub channel `pomodoro:channel:{user_id}` for cross-session broadcast.
- Browser calculates countdown from `started_at`; server publishes periodic corrections every 30s.
- On completion: publish `pomodoro.finished` via outbox → worker persists `pomodoro_history`.

**Routes:**

```text
GET/POST    /api/v1/tasks                # ?status=&category=&priority=&q=
PATCH/DELETE /api/v1/tasks/{id}
PATCH       /api/v1/tasks/{id}/status   # Move between columns
GET/PUT     /api/v1/pomodoro
POST        /api/v1/pomodoro/stop
GET         /api/v1/pomodoro/history    # ?limit=&cursor=
GET         /api/v1/ws/pomodoro         # WebSocket upgrade
```

---

### 8.5 Journal — Rich Markdown Knowledge Base

**Full-text search setup:**

```sql
-- search_vector auto-updated by trigger
CREATE OR REPLACE FUNCTION journals_tsvector_trigger() RETURNS trigger AS $$
BEGIN
  NEW.search_vector :=
    setweight(to_tsvector('english', coalesce(NEW.title, '')), 'A') ||
    setweight(to_tsvector('english', coalesce(NEW.content, '')), 'B');
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER journals_tsvector_update
  BEFORE INSERT OR UPDATE ON journals
  FOR EACH ROW EXECUTE FUNCTION journals_tsvector_trigger();
```

**Semantic search (via AI service + Elasticsearch):**

After journal save, publish `journal.created` event. Worker calls AI service to:

1. Generate sentence embedding for the journal content.
2. Index `{ id, user_id, title, content_preview, embedding, tags, published_date }` into Elasticsearch.

For semantic queries, call AI service `/ai/semantic-search`, which returns ranked journal IDs. API then fetches full records from PostgreSQL.

**Journal linking:**

```sql
-- journal_links table for explicit bi-directional links
CREATE TABLE journal_links (
  journal_id    UUID REFERENCES journals(id) ON DELETE CASCADE,
  linked_id     UUID REFERENCES journals(id) ON DELETE CASCADE,
  user_id       UUID NOT NULL,
  created_at    TIMESTAMPTZ DEFAULT now(),
  PRIMARY KEY (journal_id, linked_id),
  CHECK (journal_id != linked_id)
);
```

Linking is bi-directional: inserting `(A→B)` also inserts `(B→A)` in the same transaction.

**Streak calculation (service):**

A "streak" counts consecutive calendar days where `published_date` has at least one journal. This is a service-layer computation (not DB query), using a sorted list of distinct dates.

**Routes:**

```text
GET/POST        /api/v1/journals          # ?from=&to=&tags=&mood=&q=&limit=&cursor=
GET/PATCH/DELETE /api/v1/journals/{id}
POST            /api/v1/journals/{id}/link   # { linked_id }
DELETE          /api/v1/journals/{id}/link/{linked_id}
GET             /api/v1/journals/search   # Full-text + date range + tags
GET             /api/v1/journals/stats    # Streak, word count, mood trend
```

**Worker events:**

- `journal.created` / `journal.updated` → Trigger AI service embedding + Elasticsearch indexing.
- `journal.deleted` → Remove from Elasticsearch index.

---

## 9. Transactions, outbox, and consistency

**Rule:** Any user-facing operation that must durably trigger async work uses the transactional outbox pattern.

```text
                 ┌─────────────────────────────────────┐
  Request ──────►│  BEGIN TRANSACTION                  │
                 │  INSERT business record             │
                 │  INSERT outbox_event (pending)      │
                 │  COMMIT                             │
                 └─────────────────────────────────────┘
                          │
                          ▼
              Outbox relay (goroutine or scheduler)
                          │
              SELECT pending outbox events
                          │
              Publish to RabbitMQ (exchange: gopa.events)
                          │
              UPDATE outbox SET status = delivered
```

```sql
CREATE TABLE outbox_events (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  event_type     TEXT NOT NULL,
  schema_version INT  NOT NULL DEFAULT 1,
  payload        JSONB NOT NULL,
  occurred_at    TIMESTAMPTZ DEFAULT now(),
  status         TEXT NOT NULL DEFAULT 'pending', -- pending | delivered | failed
  attempts       INT  NOT NULL DEFAULT 0,
  last_error     TEXT,
  delivered_at   TIMESTAMPTZ
);

CREATE INDEX idx_outbox_pending ON outbox_events(status, occurred_at) WHERE status = 'pending';
```

**Idempotency in worker:**

```sql
-- Before applying any event effect:
INSERT INTO processed_events (event_id, processed_at)
VALUES ($1, now())
ON CONFLICT (event_id) DO NOTHING;
-- If 0 rows inserted → already processed; skip safely
```

---

## 10. Redis policy

| Key pattern | Value / TTL | Purpose |
|---|---|---|
| `auth:refresh:{user}:{session}` | Token hash, 7 days | Refresh session validation |
| `rate:login:{ip}` | Counter, 1 min window | Login rate limiting |
| `rate:refresh:{user}` | Counter, 1 min window | Refresh rate limiting |
| `pomodoro:state:{user}` | JSON state, 4 hours | Live timer authority |
| `pomodoro:channel:{user}` | Pub/Sub | State change fan-out |

- Every Redis call has a context with appropriate timeout (≤ 500ms for reads).
- Redis unavailability must not crash the API; handle gracefully per feature risk.
- TTLs are always set; never allow keys to leak indefinitely.

---

## 11. RabbitMQ and worker policy

### 11.1 Event envelope

```json
{
  "event_id": "uuid",
  "event_type": "learning.review.completed",
  "schema_version": 1,
  "occurred_at": "2026-08-03T09:00:00Z",
  "correlation_id": "request-id",
  "user_id": "uuid",
  "payload": {}
}
```

### 11.2 Broker standards

- Durable topic exchange: `gopa.events`.
- Routing keys: `learning.review.completed`, `learning.vocabulary.created`, `pomodoro.finished`, `finance.budget.alert`, `journal.created`, `journal.updated`, `journal.deleted`.
- Messages published persistently; queues durable.
- Acknowledge only after successful, idempotent processing.
- Bounded retry (3 attempts) with exponential backoff; dead-letter queue for failed messages.
- Malformed/unknown schema → send to dead-letter immediately; do not retry.

### 11.3 Worker graceful shutdown

On `SIGINT` / `SIGTERM`:

1. Stop accepting new deliveries (cancel consumer).
2. Wait for active handlers up to `SHUTDOWN_TIMEOUT` (default 15s).
3. Close AMQP channels and connection cleanly.
4. Unacknowledged messages automatically requeued by RabbitMQ.

---

## 12. Python AI Microservice

**Location:** `ai/` directory.
**Framework:** FastAPI + Uvicorn.
**Authentication:** Internal shared secret header `X-AI-Secret` (not exposed externally).

### 12.1 Services

**TTS Service (`edge-tts`):**

```python
# services/tts_service.py
async def generate_tts(word: str, language: str) -> str:
    """
    JP: voice = "ja-JP-NanamiNeural" (natural female voice)
    EN: voice = "en-US-JennyNeural"
    Returns: URL path to cached audio file
    Caches in local filesystem / object storage; returns existing URL if cached.
    """
```

**Vocabulary Enrichment Service:**

```python
# services/enrichment_service.py
async def enrich_vocabulary(word: str, language: str) -> EnrichResult:
    """
    JP: Uses pykakasi/fugashi to get kana reading + romaji
    EN: Gets IPA pronunciation
    Generates example sentence using OpenAI API or local Ollama (llama3)
    Generates TTS audio
    Returns: { reading, romaji, example_sentence, example_translation, audio_url }
    """
```

**Embedding Service (`sentence-transformers`):**

```python
# services/embedding_service.py
# Model: "paraphrase-multilingual-MiniLM-L12-v2" (supports JP + EN + VI)
async def embed_text(text: str) -> list[float]:
    """Returns 384-dim embedding vector"""

async def index_journal(journal_id: str, user_id: str, title: str, content: str, tags: list, published_date: str):
    """Indexes journal into Elasticsearch with embedding"""

async def semantic_search(query: str, user_id: str, limit: int = 10) -> list[SearchHit]:
    """KNN search in Elasticsearch filtered by user_id"""
```

### 12.2 Elasticsearch journal index mapping

```json
{
  "mappings": {
    "properties": {
      "journal_id":      { "type": "keyword" },
      "user_id":         { "type": "keyword" },
      "title":           { "type": "text", "analyzer": "standard" },
      "content_preview": { "type": "text", "analyzer": "standard" },
      "tags":            { "type": "keyword" },
      "published_date":  { "type": "date" },
      "embedding":       { "type": "dense_vector", "dims": 384, "index": true, "similarity": "cosine" }
    }
  }
}
```

### 12.3 AI service routes

```text
GET  /health
POST /tts                 { word, language }         → { audio_url }
POST /enrich-vocabulary   { word, language }         → { reading, romaji, example, translation, audio_url }
POST /embed               { text }                   → { embedding: [float] }
POST /index-journal       { journal_id, user_id, title, content, tags, published_date }
DELETE /index-journal/{id}
POST /semantic-search     { query, user_id, limit }  → [{ journal_id, score }]
```

---

## 13. Configuration and environment

```dotenv
# be/.env.example
APP_ENV=development
HTTP_ADDR=:8080
POSTGRES_DSN=postgres://gopa:gopa@localhost:5432/gopa?sslmode=disable
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
RABBITMQ_URL=amqp://gopa:gopa@localhost:5672/
ELASTICSEARCH_URL=http://localhost:9200
AI_SERVICE_URL=http://localhost:8000
AI_SERVICE_SECRET=change-me-in-production
JWT_ISSUER=gopa
JWT_ACCESS_SECRET=replace-with-a-long-random-secret-min-32-chars
ACCESS_TOKEN_TTL=15m
REFRESH_TOKEN_TTL=168h
LOG_LEVEL=debug
WEB_ORIGIN=http://localhost:5173
SHUTDOWN_TIMEOUT=15s
OUTBOX_POLL_INTERVAL=5s
RATE_LIMIT_WINDOW=1m
RATE_LIMIT_MAX_LOGIN=10
```

Config loader must fail fast with a specific error when any required production value is absent or uses a sample default value (e.g., `"replace-with-a-long-random-secret"`).

---

## 14. Security requirements

1. bcrypt cost ≥ 12; never return/store/log plaintext passwords.
2. Validate all untrusted input: JSON shape, UUIDs, enums, size limits, email, dates, numeric ranges.
3. Scope every query by `user_id` in SQL — application-level filtering is not sufficient.
4. Rate-limit login (10/min/IP), refresh (20/min/user), and AI TTS endpoints.
5. Explicit CORS origins; never `Access-Control-Allow-Origin: *` with credentials.
6. TLS in production; `Secure`, `HttpOnly`, `SameSite=Strict` cookie flags.
7. AI service authenticated by internal shared secret; never exposed to browser directly.
8. Sanitize Markdown HTML at render time (frontend) and store raw Markdown only.
9. Audit all dependencies and container images in CI once delivery pipeline exists.
10. JWT signing key rotation plan documented before production deployment.

---

## 15. Observability

### 15.1 Structured logging

JSON `slog` in production. Every request log includes: `timestamp`, `level`, `service`, `request_id`, `method`, `route`, `status`, `latency_ms`, `user_id` (safe). Error logs include wrapped context chain. Never log: tokens, passwords, raw journal content.

### 15.2 Health checks

- `GET /healthz` — process alive; no downstream calls; always 200 if running.
- `GET /readyz` — checks PG connectivity, Redis PING, RabbitMQ heartbeat, AI service `/health`; returns 503 if any required dep fails. Fail with specific reason.

### 15.3 Metrics roadmap (post-stabilization)

OpenTelemetry instruments: HTTP request rate/latency/error by route, DB/Redis call latency, queue depth/consumer failures, SRS review throughput, Pomodoro completion rate, budget alert frequency.

---

## 16. Testing strategy

| Layer | What to test | Approach |
|---|---|---|
| Domain/service | SM-2 algorithm, finance calculations, schedule rules | Fast table-driven unit tests; pure functions |
| Repository | SQL correctness, constraint behavior, transaction rollback | Real PostgreSQL via `testcontainers-go` |
| Handler | DTO binding, auth middleware, status codes, envelope shape | `httptest.NewRecorder` with test router |
| Broker/worker | Idempotency, ack behavior, retry | Integration with disposable RabbitMQ |
| AI service | TTS response, enrichment, embedding | pytest + mocked OpenAI API |

Test naming: `TestVocabularyService_Review_ResetsBox_WhenQualityBelowThree`.
Every bug fix includes a regression test where practical.

**Required quality gates:**

```bash
gofmt -w .
go vet ./...
go test ./...
go test -race ./...          # for concurrency-sensitive code
pytest ai/ -v                # AI service tests
```

---

## 17. Local development

```bash
make up              # docker compose up -d (PG, Redis, RMQ, Elasticsearch)
make down            # docker compose down
make api             # go run ./cmd/api
make worker          # go run ./cmd/worker
make ai              # uvicorn main:app --reload --port 8000
make migrate-up      # golang-migrate up
make migrate-down    # golang-migrate down 1
make test            # go test ./...
make test-integration
make lint            # gofmt + go vet + golangci-lint
make build           # go build ./cmd/api ./cmd/worker
```

---

## 18. Implementation roadmap & progress tracking

> **For AI Agents:** This tracks what has been implemented in the backend. Update status when completing items. Cross-reference with root `MASTER_PROMPT.md` Section 13 for the full project matrix.

### PHASE-0: Bootstrap
- [ ] **BE-001** Go module, directory skeleton
- [ ] **BE-002** Docker Compose (PG + Redis + RMQ + Elasticsearch)
- [ ] **BE-003** Config loader with validation
- [ ] **BE-004** slog structured logger
- [ ] **BE-005** PG pool + `WithTx` helper
- [ ] **BE-006** Redis client
- [ ] **BE-007** RabbitMQ connection + graceful shutdown
- [ ] **BE-008** `/healthz` + `/readyz`
- [ ] **BE-009** Domain errors + response helpers (`pkg/utils/response`)
- [ ] **BE-010** Outbox table migration + relay goroutine
- [ ] **BE-011** Makefile

### PHASE-1: IAM
- [ ] **BE-101** `users` + `refresh_sessions` migration
- [ ] **BE-102** User repository
- [ ] **BE-103** Auth service (register, login, refresh, logout)
- [ ] **BE-104** JWT middleware
- [ ] **BE-105** Auth handlers + routes
- [ ] **BE-106** Redis rate limiting middleware
- [ ] **BE-107** `/me` GET + PATCH

### PHASE-2: Deep Work
- [ ] **BE-201** `tasks` + `pomodoro_history` migration
- [ ] **BE-202** Task repository + service
- [ ] **BE-203** Task CRUD handlers
- [ ] **BE-204** Pomodoro Redis state + service
- [ ] **BE-205** WebSocket handler (`/ws/pomodoro`)
- [ ] **BE-206** `pomodoro.finished` event + worker consumer

### PHASE-3: Linguistics
- [ ] **BE-301** Python AI service (FastAPI, health, TTS, enrichment, embedding)
- [ ] **BE-302** `vocabularies` + `vocabulary_reviews` + `learning_sessions` migration
- [ ] **BE-303** Vocabulary repository + SM-2 domain function
- [ ] **BE-304** Vocabulary CRUD + bulk import service
- [ ] **BE-305** Review queue endpoint (SM-2 due items)
- [ ] **BE-306** Review submission + outbox event
- [ ] **BE-307** Worker: idempotent SM-2 update consumer
- [ ] **BE-308** Worker: vocabulary enrichment via AI service
- [ ] **BE-309** Audio TTS endpoint
- [ ] **BE-310** Learning session endpoints
- [ ] **BE-311** Stats endpoint

### PHASE-4: Finance
- [x] **BE-401** Finance migrations (accounts, categories, transactions, budgets, savings_goals)
- [x] **BE-402** Account repository + service
- [ ] **BE-403** Category repository + service + seed data
- [x] **BE-404** Transaction repository + service (balance update in tx)
- [ ] **BE-405** Budget service (utilization + alert event)
- [ ] **BE-406** Savings goal service
- [x] **BE-407** Analytics service (cashflow, by-category, net-worth)
- [ ] **BE-408** Finance handlers + routes

### PHASE-5: Journal
- [ ] **BE-501** Elasticsearch index setup + `journals` + `journal_links` migration
- [x] **BE-502** `search_vector` trigger
- [x] **BE-503** Journal repository + service
- [x] **BE-504** Journal CRUD + full-text search handlers
- [x] **BE-505** Journal linking endpoint
- [x] **BE-506** Journal stats (streak, word count, mood)
- [ ] **BE-507** Worker: journal indexing to Elasticsearch via AI service

---

## 19. Required AI execution protocol

Do **not** implement the entire backend in one response. Build one vertical slice at a time.

For each task, respond using this structure:

### Phase A — Architecture + Go vs Java

- State the deliverable and narrow scope.
- Explain component placement and Clean Architecture layer.
- Compare the significant Go decision to its Java/Spring equivalent.
- Identify transactions, context, concurrency, security, and trade-off decisions.
- Reference relevant `[BE-XXX]` items.

### Phase B — Commands

Exact, copyable commands: `go get`, migration SQL, Docker, tests. State expected output.

### Phase C — Complete code

- Full, untruncated files with exact project-relative paths.
- All imports, constructors, migrations, tests, and wiring needed for the slice.
- No pseudo-code, "omitted for brevity," or unexplained `any` types.
- Respect this architecture and all code standards.

### Phase D — Verification + handoff

- Format, test, manual verification commands.
- Mark completed `[BE-XXX]` items.
- State next smallest slice.
- Stop for review before next module.

---

## 20. Definition of done (backend)

A backend feature is done when:

- [ ] Domain behavior and HTTP/event contract are explicit and documented
- [ ] User ownership enforced in SQL
- [ ] Migrations, constraints, and indexes support behavior
- [ ] Multi-record operations use explicit transactions (ACID)
- [ ] API errors are safe, consistent, and tested
- [ ] Outbox pattern used where async reliability matters
- [ ] Contexts, timeouts, and graceful cleanup are correct
- [ ] Unit tests cover domain rules; integration tests cover DB + broker paths
- [ ] Structured logs appropriate; no sensitive data logged
- [ ] `gofmt`, `go vet`, `go test ./...` pass
- [ ] Relevant `[BE-XXX]` items marked complete
- [ ] README / API docs updated

---

*This document is the backend source of truth. Any architectural decision that supersedes a section must be documented here before implementation begins.*
