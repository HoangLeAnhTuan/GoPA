# GoPA — Master Architecture Specification & Engineering Playbook

> **Status:** Living specification — update this document when any architectural decision changes.
> **Primary language:** English (technical), Vietnamese (mentorship explanations where useful)
> **Target stack:** Go 1.22+ · Python 3.12+ (AI microservice) · PostgreSQL 16+ · Elasticsearch 8+ · Redis 7+ · RabbitMQ 3.13+ · React 18 + TypeScript

---

## 1. Project identity

**GoPA** (*Golang Personal Assistant*) is a private, full-stack web application for managing focused work, language learning, personal finance, and a rich Markdown knowledge journal. It is deliberately designed as a practical vehicle for learning idiomatic Go while producing a maintainable, production-minded application.

### 1.1 Product goals

GoPA helps one person to:

- plan and complete deep-work tasks, including daily LeetCode practice, with a synchronized Pomodoro timer;
- learn Japanese (N3 → N2) and English (TOEIC 805 → 900+) through a smart, interactive spaced-repetition engine with AI-assisted pronunciation and context generation;
- track personal finances — income, expenses, budgets, savings goals, and multi-account wallets — with clear analytics;
- write, tag, link, search, and reflect in a daily Markdown journal with full-text and semantic search.

### 1.2 Non-goals for the first release

- Multi-tenant collaboration, billing, or public sharing.
- Native mobile applications.
- Vehicle / home-network asset tracking (removed; replaced by Finance module).
- A generic enterprise workflow engine.
- Premature microservice decomposition. GoPA runs as a modular monolith plus one focused AI microservice.

### 1.3 Engineering principles

1. **Correctness before cleverness.** Prefer readable, explicit code over magic.
2. **Domain-first design.** Business rules must not depend on Gin, SQL, Redis, or React.
3. **Boring technology.** Stable, well-understood libraries and standard patterns.
4. **Explicit boundaries.** Interfaces belong where a core service depends on a replaceable capability.
5. **Observable operation.** Important work is logged, measurable, and safe to retry.
6. **Secure by default.** Secrets stay outside source control; auth and input validation are mandatory.
7. **ACID compliance.** Every business state change that touches multiple records uses explicit transactions.
8. **Incremental delivery.** Build, test, and review one vertical slice at a time.

---

## 2. Developer context

The primary developer is an experienced Java/Spring Boot engineer transitioning to Go. The AI acts as pair programmer, Principal Go Engineer, and mentor.

| Concern | Java / Spring | GoPA / Go |
|---|---|---|
| Dependency injection | Framework container, annotations | Explicit construction in `main.go` |
| Request state | `ThreadLocal`, framework context | `context.Context` passed explicitly |
| Concurrency | Managed thread pools, futures | Goroutines, channels, cancellation contexts |
| ORM / mapping | JPA, Hibernate, MyBatis | SQL-first with `sqlx`; no ORM |
| Transactions | `@Transactional` proxy | Explicit `*sqlx.Tx` within a use case |
| Errors | Checked/unchecked exceptions | Returned `error`, wrapped with context |
| HTTP controllers | Annotation-based methods | Explicit Gin route registration and handlers |

---

## 3. Scope and feature map

| Module | Responsibility | Status |
|---|---|---|
| **IAM** | Identity, session lifecycle, authorization | Phase 1 |
| **Deep Work** | Tasks, Kanban, Pomodoro timer | Phase 2 |
| **Linguistics** | Smart SRS, interactive study modes, AI TTS | Phase 3 |
| **Finance** | Wallets, transactions, budgets, analytics | Phase 4 |
| **Journal** | Rich Markdown, full-text + semantic search, linking | Phase 5 |

### 3.1 High-level system view

```mermaid
flowchart LR
  Browser["React web client"] -->|"HTTPS / REST"| API["Go API · Gin"]
  Browser -->|"WebSocket"| WS["Pomodoro gateway"]
  API --> PostgreSQL[("PostgreSQL")]
  API --> Redis[("Redis")]
  API --> RabbitMQ[("RabbitMQ")]
  API -->|"Internal HTTP"| AISvc["Python AI Service (FastAPI)"]
  WS --> Redis
  RabbitMQ --> Worker["Go worker"]
  Worker --> PostgreSQL
  Worker --> Redis
  AISvc --> Elasticsearch[("Elasticsearch")]
  AISvc --> PostgreSQL
```

### 3.2 Runtime topology

| Process | Language | Role |
|---|---|---|
| **API** | Go | HTTP REST + WebSocket; validation, auth, use-case orchestration |
| **Worker** | Go | Async event consumers; SRS scheduling, Pomodoro history |
| **AI Service** | Python / FastAPI | TTS generation, semantic search embedding, AI context generation |
| **PostgreSQL** | — | Primary durable source of truth |
| **Elasticsearch** | — | Full-text and semantic journal/vocabulary search |
| **Redis** | — | Session state, rate limits, live Pomodoro state |
| **RabbitMQ** | — | Async decoupling; at-least-once delivery to worker |

---

## 4. Technology decisions

| Area | Choice | Rationale |
|---|---|---|
| Language (BE) | Go 1.22+ | Strong stdlib, simple concurrency, single-binary deployment |
| HTTP framework | `gin-gonic/gin` | Productive routing/middleware; isolated to handler adapters |
| Database access | `jmoiron/sqlx` + `pgx` | SQL-first with struct scanning; no ORM magic |
| Database | PostgreSQL 16+ | Relational integrity, JSONB, full-text search, UUID support |
| Search | Elasticsearch 8+ | Journal semantic search + linguistics vocabulary search |
| Cache | `redis/go-redis/v9` | Sessions, rate limits, Pomodoro live state |
| Broker | RabbitMQ + `amqp091-go` | Durable async decoupling |
| Auth | JWT + bcrypt + opaque refresh tokens | Short access lifetime, server-controlled revocation |
| Logging | `log/slog` | Structured JSON; stdlib, no heavy dep |
| AI microservice | Python 3.12 + FastAPI | TTS (gTTS/edge-tts), sentence embedding (sentence-transformers), OpenAI-compatible |
| Frontend | React 18 + Vite + TypeScript | Fast dev cycle, type safety |
| Styling | Tailwind CSS + shadcn/ui | Consistent, accessible components; Glassmorphism design |
| State | TanStack Query + Zustand | Server cache + minimal UI state |
| Icons | Lucide React + Heroicons | Modern, consistent icon set |
| Animation | Framer Motion | Purposeful motion, reduced-motion aware |
| Deployment | Docker Compose (dev) | Reproducible local environment |

---

## 5. Repository architecture

```text
.
├── be/                              # Go backend
│   ├── cmd/
│   │   ├── api/main.go              # API composition root
│   │   └── worker/main.go           # Worker composition root
│   ├── internal/
│   │   ├── core/
│   │   │   ├── domain/              # Entities, value objects, domain errors
│   │   │   └── ports/               # Repository, cache, broker, service contracts
│   │   ├── adapters/
│   │   │   ├── handler/http/        # Gin handlers, route registration, DTOs
│   │   │   ├── repository/          # PostgreSQL/sqlx implementations
│   │   │   ├── cache/               # Redis implementations
│   │   │   └── broker/              # RabbitMQ publisher/consumer
│   │   └── services/                # Use-case orchestration
│   ├── pkg/
│   │   ├── config/                  # Env loading and validation
│   │   ├── database/                # Pool, WithTx helper, migrations
│   │   ├── logger/                  # slog setup
│   │   └── utils/                   # JWT, password, response/error helpers
│   ├── migrations/                  # Ordered raw SQL up/down files
│   ├── docker-compose.yml
│   ├── Makefile
│   └── .env.example
├── ai/                              # Python AI microservice
│   ├── main.py                      # FastAPI application
│   ├── services/
│   │   ├── tts_service.py           # Text-to-speech (gTTS, edge-tts)
│   │   ├── embedding_service.py     # Sentence embeddings
│   │   └── context_service.py       # OpenAI/Ollama context generation
│   ├── routers/
│   ├── requirements.txt
│   └── Dockerfile
└── fe/                              # React frontend
    ├── src/
    │   ├── app/                     # Providers, router, shell
    │   ├── components/              # Design system primitives
    │   ├── features/                # Feature-local code
    │   ├── hooks/                   # Cross-feature hooks
    │   ├── layouts/                 # AppLayout, AuthLayout
    │   ├── lib/                     # Axios, cn(), query keys, formatters
    │   ├── locales/                 # i18n resources
    │   ├── stores/                  # Zustand UI stores
    │   └── styles/                  # Global CSS, design tokens
    └── ...config files
```

---

## 6. Domain model

### 6.1 Shared conventions

- UUIDs are primary keys; generated by PostgreSQL (`gen_random_uuid()`) or the application.
- All timestamps are `TIMESTAMPTZ` stored in UTC.
- Monetary values use `NUMERIC(14,2)`; never `float`.
- Every user-owned record contains `user_id` and every query scopes by it.
- `created_at` and `updated_at` on every mutable entity; `updated_at` set by trigger or explicit update.
- Soft deletion only when an explicit audit/recovery requirement exists.

### 6.2 IAM

**`users`**

| Field | Type | Notes |
|---|---|---|
| `id` | UUID | PK |
| `email` | TEXT (lowercase) | Unique, canonicalized before storage |
| `password_hash` | TEXT | bcrypt; never returned |
| `display_name` | TEXT | User-visible name |
| `role` | TEXT | `USER` or `ADMIN` |
| `avatar_url` | TEXT nullable | Profile image URL |
| `created_at`, `updated_at` | TIMESTAMPTZ | UTC |

**Authentication policy:**

- Access JWT: 15 minutes; claims: `sub`, `role`, `sid`, `iat`, `exp`.
- Refresh token: 7 days, opaque random value stored hashed in Redis under `auth:refresh:{user_id}:{session_id}`.
- Rotate refresh tokens on every refresh; delete old before issuing new.
- Logout deletes active refresh key; access tokens naturally expire.
- Rate-limit login and refresh endpoints.

### 6.3 Linguistics — Smart Learning Engine

**`vocabularies`**

| Field | Type | Notes |
|---|---|---|
| `id`, `user_id` | UUID | Ownership |
| `language` | TEXT | `JP` or `EN` |
| `word` | TEXT | Surface form |
| `reading` | TEXT nullable | Furigana/kana for JP; IPA/phonetic for EN |
| `meaning` | TEXT | In learner's preferred language |
| `example_sentence` | TEXT nullable | Contextual usage |
| `example_translation` | TEXT nullable | Translation of example |
| `tags` | TEXT[] | e.g. `["N3", "verb", "JLPT"]` |
| `difficulty_level` | TEXT | `BEGINNER`, `INTERMEDIATE`, `ADVANCED` |
| `audio_url` | TEXT nullable | Cached TTS audio URL |
| `image_url` | TEXT nullable | Optional mnemonic image |
| `ease_factor` | NUMERIC(4,2) | SM-2 ease factor, default 2.5 |
| `interval_days` | INTEGER | SM-2 current interval |
| `repetition_count` | INTEGER | Total successful repetitions |
| `current_box` | SMALLINT | Leitner box 1–5 |
| `next_review_at` | TIMESTAMPTZ | Queue ordering |
| `mastery_score` | SMALLINT | 0–100, calculated from review history |
| `created_at`, `updated_at` | TIMESTAMPTZ | — |

**SM-2 Algorithm (SuperMemo-2):**

For quality score `q` (0–5) after each review:

```
if q >= 3:
    if repetition_count == 0: interval = 1
    elif repetition_count == 1: interval = 6
    else: interval = round(prev_interval * ease_factor)
    repetition_count++
else:
    repetition_count = 0
    interval = 1

ease_factor = ease_factor + (0.1 - (5 - q) * (0.08 + (5 - q) * 0.02))
ease_factor = max(1.3, ease_factor)
next_review_at = now + interval_days
```

Keep this algorithm in a **domain service** (isolated from handlers/repos) so it can later be replaced by FSRS without breaking anything else.

**`vocabulary_reviews`**

| Field | Type | Notes |
|---|---|---|
| `id`, `vocabulary_id`, `user_id` | UUID | — |
| `quality_score` | SMALLINT | 0–5 (SM-2 scale) |
| `study_mode` | TEXT | `FLASHCARD`, `MULTIPLE_CHOICE`, `TYPE_IN`, `AUDIO_QUIZ` |
| `response_time_ms` | INTEGER | Time to answer |
| `was_correct` | BOOLEAN | For quiz modes |
| `reviewed_at` | TIMESTAMPTZ | — |

**`learning_sessions`**

| Field | Type | Notes |
|---|---|---|
| `id`, `user_id` | UUID | — |
| `language` | TEXT | `JP` or `EN` |
| `session_type` | TEXT | `REVIEW`, `LEARN_NEW`, `PRACTICE` |
| `items_reviewed` | INTEGER | — |
| `items_correct` | INTEGER | — |
| `duration_seconds` | INTEGER | — |
| `started_at`, `ended_at` | TIMESTAMPTZ | — |

**Interactive Study Modes:**

1. **Flashcard** — 3D card flip (front: word + reading, back: meaning + example + audio). Quality score 1–5.
2. **Multiple Choice** — 4 options (1 correct, 3 distractors from same language/difficulty). Binary correct/incorrect.
3. **Type-in** — User types the reading or meaning. Fuzzy match with leniency for kana/romaji.
4. **Audio Quiz** — TTS plays the word; user identifies meaning from choices.

**AI Microservice responsibilities (Python):**

- **TTS generation:** `edge-tts` (offline capable, natural voices) for JP and EN. Cache audio files; return URL.
- **Example sentence generation:** OpenAI-compatible API or local Ollama to generate natural example sentences with translation.
- **Semantic search:** `sentence-transformers` (multilingual) to embed vocabulary and journal content into Elasticsearch.
- **Context enrichment:** On vocabulary creation, enrich with reading (for JP, use `pykakasi` or `fugashi`), generate example sentence if missing, generate TTS audio.

**Endpoints for AI Service:**

```
POST /ai/tts                  # { word, language } → { audio_url }
POST /ai/enrich-vocabulary    # { word, language } → { reading, example, translation, audio_url }
POST /ai/semantic-search      # { query, language?, limit } → [{ id, score }]
```

### 6.4 Personal Finance

**`accounts`** (Wallets)

| Field | Type | Notes |
|---|---|---|
| `id`, `user_id` | UUID | — |
| `name` | TEXT | "VietcomBank Savings", "Cash Wallet" |
| `type` | TEXT | `CASH`, `BANK`, `CREDIT_CARD`, `SAVINGS`, `INVESTMENT`, `CRYPTO` |
| `currency` | TEXT | `VND`, `USD`, `JPY` |
| `initial_balance` | NUMERIC(14,2) | Starting balance |
| `current_balance` | NUMERIC(14,2) | Maintained by transaction triggers or recalculation |
| `color` | TEXT | Hex color for UI |
| `icon` | TEXT | Icon identifier |
| `is_archived` | BOOLEAN | Hide without deleting |
| `created_at`, `updated_at` | TIMESTAMPTZ | — |

**`categories`**

| Field | Type | Notes |
|---|---|---|
| `id`, `user_id` | UUID | — |
| `name` | TEXT | e.g. "Food & Dining", "Transport" |
| `type` | TEXT | `INCOME` or `EXPENSE` |
| `icon` | TEXT | Icon identifier |
| `color` | TEXT | Hex color |
| `parent_id` | UUID nullable | Sub-categories (max 2 levels deep) |
| `monthly_budget` | NUMERIC(14,2) nullable | Optional default budget |

**`transactions`**

| Field | Type | Notes |
|---|---|---|
| `id`, `user_id` | UUID | — |
| `account_id` | UUID | Source account |
| `to_account_id` | UUID nullable | Target for transfers |
| `category_id` | UUID nullable | Null for transfers |
| `type` | TEXT | `INCOME`, `EXPENSE`, `TRANSFER` |
| `amount` | NUMERIC(14,2) | Always positive; direction from type |
| `currency` | TEXT | Transaction currency |
| `exchange_rate` | NUMERIC(10,6) | For multi-currency; 1.0 for same currency |
| `description` | TEXT | User note |
| `merchant` | TEXT nullable | "Grab", "Coop Mart" |
| `tags` | TEXT[] | Free-form tags |
| `occurred_at` | TIMESTAMPTZ | Actual transaction time |
| `is_recurring` | BOOLEAN | — |
| `recurring_rule` | JSONB nullable | `{ frequency, until }` |
| `created_at`, `updated_at` | TIMESTAMPTZ | — |

**`budgets`**

| Field | Type | Notes |
|---|---|---|
| `id`, `user_id` | UUID | — |
| `category_id` | UUID nullable | Null = overall spending budget |
| `name` | TEXT | "Monthly Groceries", "Q1 2026 Savings" |
| `amount` | NUMERIC(14,2) | Budget limit |
| `period` | TEXT | `WEEKLY`, `MONTHLY`, `QUARTERLY`, `YEARLY`, `CUSTOM` |
| `start_date`, `end_date` | DATE | — |
| `alert_threshold` | NUMERIC(4,2) | Notify at e.g. 0.8 (80% spent) |
| `is_active` | BOOLEAN | — |

**`savings_goals`**

| Field | Type | Notes |
|---|---|---|
| `id`, `user_id` | UUID | — |
| `name` | TEXT | "Japan Trip 2027", "Emergency Fund" |
| `target_amount` | NUMERIC(14,2) | — |
| `current_amount` | NUMERIC(14,2) | Linked account balance or manual |
| `linked_account_id` | UUID nullable | — |
| `target_date` | DATE nullable | — |
| `color`, `icon` | TEXT | UI display |

**Finance analytics (service layer):**

- Monthly cashflow summary (total income, expense, net, savings rate).
- Spending by category for a date range.
- Budget utilization with alert status.
- Net worth = sum of all account balances.
- Trend charts (weekly/monthly aggregations from transactions).

### 6.5 Deep Work

**`tasks`:** `id`, `user_id`, `title`, `description`, `status`, `priority`, `category`, `due_date`, `position`, `estimated_pomodoros`, `completed_pomodoros`, `created_at`, `updated_at`.

- Status: `TODO`, `IN_PROGRESS`, `DONE`.
- Priority: `LOW`, `MEDIUM`, `HIGH`, `URGENT`.
- Category: `WORK`, `STUDY`, `LIFE`, `LEETCODE`.
- `position` is a numeric ordering scoped by user + status.

**Pomodoro state** in Redis (`pomodoro:state:{user_id}`):

```json
{
  "status": "running",
  "started_at": "2026-08-03T09:00:00Z",
  "duration_seconds": 1500,
  "break_duration_seconds": 300,
  "task_id": "optional-uuid",
  "session_count": 2,
  "updated_at": "2026-08-03T09:00:00Z"
}
```

On completion, publish `pomodoro.finished`; worker creates durable `pomodoro_history` record idempotently.

**`pomodoro_history`:** `id`, `user_id`, `task_id` (nullable), `duration_seconds`, `type` (`WORK` or `BREAK`), `completed_at`.

### 6.6 Journal — Daily Knowledge Base

**`journals`**

| Field | Type | Notes |
|---|---|---|
| `id`, `user_id` | UUID | — |
| `title` | TEXT | — |
| `content` | TEXT | Markdown source |
| `mood` | TEXT nullable | `GREAT`, `GOOD`, `OKAY`, `BAD`, `TERRIBLE` |
| `energy_level` | SMALLINT nullable | 1–5 |
| `tags` | TEXT[] | Free-form |
| `linked_journal_ids` | UUID[] | Bi-directional journal links |
| `pinned` | BOOLEAN | — |
| `search_vector` | TSVECTOR | PG full-text, auto-updated by trigger |
| `embedding_vector` | vector(384) nullable | pgvector for semantic search (optional) |
| `word_count` | INTEGER | Computed on save |
| `published_date` | DATE | The date the journal "belongs to" (default: today) |
| `created_at`, `updated_at` | TIMESTAMPTZ | — |

**Journal features:**

- PostgreSQL full-text search via `search_vector` (GIN index).
- Optional semantic search via Elasticsearch (content embedded by AI service).
- Date-range filtering: search by `published_date` range (week, month, custom).
- Mood/energy tracking and trend visualization.
- Bi-directional linking: a journal can reference other journals (`linked_journal_ids`).
- Daily journaling streak tracking.
- Word count and writing time statistics.
- Autosave draft (browser local storage); explicit save to server.
- Markdown rendering with sanitization; code blocks with syntax highlighting.

---

## 7. API conventions

### 7.1 Base rules

- Base path: `/api/v1`.
- JSON only for REST. WebSocket for Pomodoro live sync.
- UUIDs as strings; ISO 8601 / RFC 3339 timestamps in UTC.
- Validate DTOs at the HTTP boundary; never expose database structs directly.
- Pagination: `limit` (default 20, max 100) + cursor-based `cursor` for growing lists.
- `GET /healthz` (liveness) and `GET /readyz` (dependency readiness).

### 7.2 Response envelope

```json
// Success
{ "data": {}, "meta": { "request_id": "01J..." } }

// Paginated
{ "data": [], "meta": { "request_id": "01J...", "next_cursor": "...", "total": 142 } }

// Error
{ "error": { "code": "VALIDATION_ERROR", "message": "...", "details": [] }, "meta": { "request_id": "01J..." } }
```

### 7.3 Endpoint inventory

**IAM**
| Method | Path | Purpose |
|---|---|---|
| POST | `/auth/register` | Create user |
| POST | `/auth/login` | Issue credentials |
| POST | `/auth/refresh` | Rotate refresh token |
| POST | `/auth/logout` | Revoke session |
| GET | `/me` | Current user profile |
| PATCH | `/me` | Update profile |

**Linguistics**
| Method | Path | Purpose |
|---|---|---|
| GET/POST | `/vocabularies` | List/create vocabulary |
| GET/PATCH/DELETE | `/vocabularies/{id}` | Manage a vocabulary item |
| POST | `/vocabularies/import` | Bulk import (CSV/JSON) |
| GET | `/vocabularies/review-queue` | Due items (`?language=JP&mode=FLASHCARD`) |
| POST | `/vocabularies/{id}/reviews` | Submit review event |
| POST | `/vocabularies/{id}/audio` | Generate/fetch TTS audio |
| GET | `/vocabularies/stats` | Learning progress and heatmap |
| POST | `/learning-sessions` | Start learning session |
| PATCH | `/learning-sessions/{id}` | End learning session |

**Finance**
| Method | Path | Purpose |
|---|---|---|
| GET/POST | `/accounts` | List/create accounts |
| GET/PATCH/DELETE | `/accounts/{id}` | Manage account |
| GET | `/accounts/{id}/balance-history` | Balance over time |
| GET/POST | `/transactions` | List/create transactions |
| GET/PATCH/DELETE | `/transactions/{id}` | Manage transaction |
| GET | `/transactions/summary` | Period cashflow summary |
| GET | `/transactions/by-category` | Spending breakdown |
| GET/POST | `/budgets` | List/create budgets |
| GET/PATCH/DELETE | `/budgets/{id}` | Manage budget |
| GET | `/budgets/status` | Utilization with alerts |
| GET/POST | `/savings-goals` | List/create savings goals |
| PATCH | `/savings-goals/{id}` | Update goal progress |
| GET | `/finance/net-worth` | Net worth summary |
| GET | `/finance/dashboard` | Full dashboard data |

**Deep Work**
| Method | Path | Purpose |
|---|---|---|
| GET/POST | `/tasks` | List/create tasks |
| PATCH/DELETE | `/tasks/{id}` | Update/delete task |
| PATCH | `/tasks/{id}/status` | Move task status |
| GET/PUT | `/pomodoro` | Read/replace timer state |
| POST | `/pomodoro/stop` | Stop timer |
| GET | `/ws/pomodoro` | Authenticated WebSocket |
| GET | `/pomodoro/history` | Session history |

**Journal**
| Method | Path | Purpose |
|---|---|---|
| GET/POST | `/journals` | List/create (`?from=&to=&tags=&mood=&q=`) |
| GET/PATCH/DELETE | `/journals/{id}` | Manage journal |
| POST | `/journals/{id}/link` | Link to another journal |
| GET | `/journals/stats` | Writing streak, word count, mood trends |
| GET | `/journals/search` | Full-text + date-range search |

---

## 8. Common abstractions & anti-boilerplate standards

### 8.1 Backend — Reduce Go boilerplate

**Generic base repository helper (use sparingly; only for truly repetitive patterns):**

```go
// pkg/database/tx.go
func WithTx(ctx context.Context, db *sqlx.DB, fn func(*sqlx.Tx) error) error {
    tx, err := db.BeginTxx(ctx, nil)
    if err != nil { return fmt.Errorf("begin tx: %w", err) }
    defer func() {
        if p := recover(); p != nil { _ = tx.Rollback(); panic(p) }
    }()
    if err := fn(tx); err != nil {
        _ = tx.Rollback()
        return err
    }
    return tx.Commit()
}
```

**Unified HTTP response builder:**

```go
// pkg/utils/response.go
func OK(c *gin.Context, data any)           { c.JSON(200, envelope{Data: data, Meta: meta(c)}) }
func Created(c *gin.Context, data any)      { c.JSON(201, envelope{Data: data, Meta: meta(c)}) }
func Paginated(c *gin.Context, data any, cursor string, total int) { ... }
func Fail(c *gin.Context, err error)        { /* map domain errors → HTTP status */ }
```

**Domain errors:**

```go
// internal/core/domain/errors.go
var (
    ErrNotFound     = errors.New("not found")
    ErrConflict     = errors.New("conflict")
    ErrUnauthorized = errors.New("unauthorized")
    ErrForbidden    = errors.New("forbidden")
    ErrValidation   = errors.New("validation failed")
)
```

### 8.2 Frontend — Design system primitives

Build once, use everywhere. No one-off styles inside feature components.

```text
components/
├── ui/                        # shadcn primitives (auto-generated)
├── design-system/
│   ├── GlassPanel.tsx         # Translucent surface with blur variants
│   ├── PageHeader.tsx         # Title + description + actions
│   ├── StatCard.tsx           # Metric card with trend indicator
│   ├── DataTable.tsx          # Generic sortable/filterable table
│   ├── CommandPalette.tsx     # Ctrl+K global search/action
│   ├── EmptyState.tsx         # Absence explanation + CTA
│   ├── ConfirmDialog.tsx      # Destructive action dialog
│   └── AmountDisplay.tsx      # Locale-aware currency renderer
└── feedback/
    ├── LoadingSkeleton.tsx    # Content-shaped loading states
    ├── ErrorState.tsx         # Friendly error + retry
    └── Toast.tsx              # Accessible notifications
```

**Cross-feature custom hooks:**

```text
hooks/
├── useDebounce.ts             # Debounced value hook
├── useKeyboardShortcut.ts     # Global keyboard shortcut registration
├── useAudioTTS.ts             # Play vocabulary TTS audio
├── useWebSocket.ts            # Authenticated WS with exponential backoff
├── useInfiniteScroll.ts       # Intersection observer pagination
└── useLocalDraft.ts           # Local storage draft persistence
```

---

## 9. Security requirements

1. bcrypt passwords (cost ≥ 12); never store/log plaintext.
2. Validate all untrusted fields: JSON shape, UUIDs, enums, size limits, email, dates, money ranges.
3. Scope every SQL query by `user_id`; a guessed UUID reveals nothing.
4. Rate-limit login, refresh, and AI endpoints via Redis.
5. Explicit CORS origins; never wildcard with credentials.
6. TLS in production; `Secure`, `HttpOnly`, `SameSite` cookie flags.
7. Never commit `.env`; use `.env.example` with placeholders.
8. Sanitize rendered Markdown (server-side or DOMPurify on frontend).
9. AI service endpoints authenticated internally (shared secret or mTLS).
10. Audit dependency and container-image vulnerabilities in CI.

---

## 10. Observability and reliability

- **Logging:** `slog` JSON in production; request ID, method, route, status, latency, authenticated user ID (safe).
- **Health:** `/healthz` (process alive), `/readyz` (checks PG + Redis + RabbitMQ + Elasticsearch + AI service).
- **Metrics (roadmap):** OpenTelemetry after core flows stable — HTTP latency/errors, DB/Redis latency, queue depth, business counters.
- **Tracing (roadmap):** Distributed trace IDs linking API → Worker → AI Service.
- **Backup:** Automated PostgreSQL backups; practice restore; treat Redis as reconstructable.

---

## 11. Local development

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
AI_SERVICE_SECRET=change-me
JWT_ISSUER=gopa
JWT_ACCESS_SECRET=replace-with-a-long-random-secret
ACCESS_TOKEN_TTL=15m
REFRESH_TOKEN_TTL=168h
LOG_LEVEL=debug
WEB_ORIGIN=http://localhost:5173
SHUTDOWN_TIMEOUT=15s
```

**Make targets:**

```text
make up                  # Start all Docker dependencies
make down                # Stop all Docker dependencies
make api                 # Run API
make worker              # Run worker
make ai                  # Run Python AI service
make web                 # Run Vite dev server
make test                # Unit tests
make test-integration    # Integration tests
make lint                # Format + vet + eslint
make migrate-up          # Apply migrations
make migrate-down        # Revert one migration
make build               # Build all binaries
```

---

## 12. Testing strategy

| Layer | Tests | Verify |
|---|---|---|
| Domain/service | Fast table-driven unit | Rules, SRS schedule, finance calculations |
| Repository | PG integration | SQL, scanning, constraints, transactions |
| Handler | HTTP with test router | Binding, auth, status, envelope |
| Worker/broker | Integration | Idempotency, retries, ack |
| AI service | pytest | TTS, embedding, enrichment |
| Frontend | Vitest + RTL | Rendering, keyboard, mutation states |
| E2E | Playwright | Login, task, review, finance, journal search |

---

## 13. Implementation roadmap & progress tracking

> **For AI Agents:** Update the status below as each item is completed. This is the source of truth for project progress.

### PHASE-0: Infrastructure & Bootstrap
- [ ] **FEAT-001** Go module init, directory skeleton, `go.mod`
- [ ] **FEAT-002** Docker Compose (PG, Redis, RabbitMQ, Elasticsearch)
- [ ] **FEAT-003** Python AI service Dockerfile + docker-compose service
- [ ] **FEAT-004** Config loader with validation (fail-fast on missing/unsafe values)
- [ ] **FEAT-005** Structured slog logger setup
- [ ] **FEAT-006** Database pool + `WithTx` helper + migration runner
- [ ] **FEAT-007** Redis client setup
- [ ] **FEAT-008** RabbitMQ connection + graceful shutdown
- [ ] **FEAT-009** `/healthz` and `/readyz` endpoints
- [ ] **FEAT-010** Makefile with all targets
- [ ] **FEAT-011** Vite + React + TypeScript project init
- [ ] **FEAT-012** Tailwind CSS + shadcn/ui + design tokens
- [ ] **FEAT-013** i18n setup (vi, en, ja)
- [ ] **FEAT-014** Base AppLayout (Glassmorphism, sidebar, topbar)
- [ ] **FEAT-015** Common response helpers + domain errors

### PHASE-1: IAM
- [ ] **FEAT-101** `users` migration + repository
- [ ] **FEAT-102** Registration endpoint
- [ ] **FEAT-103** Login endpoint + JWT issuance
- [ ] **FEAT-104** Refresh token rotation
- [ ] **FEAT-105** Logout endpoint
- [ ] **FEAT-106** JWT middleware
- [ ] **FEAT-107** `/me` GET + PATCH
- [ ] **FEAT-108** Rate limiting (Redis)
- [ ] **FEAT-109** Auth UI (Login/Register pages, Axios interceptor, protected routes)

### PHASE-2: Deep Work
- [ ] **FEAT-201** `tasks` migration + repository
- [ ] **FEAT-202** Task CRUD endpoints
- [ ] **FEAT-203** Kanban status/order management
- [ ] **FEAT-204** Pomodoro Redis state management
- [ ] **FEAT-205** Authenticated WebSocket `/ws/pomodoro`
- [ ] **FEAT-206** `PomodoroFinished` event + worker history persistence
- [ ] **FEAT-207** Task board UI (Kanban columns, drag support)
- [ ] **FEAT-208** Pomodoro widget (global floating + detail view)
- [ ] **FEAT-209** Today page (dashboard: active timer, next task, streak)

### PHASE-3: Linguistics
- [ ] **FEAT-301** Python AI service setup (FastAPI, health check)
- [ ] **FEAT-302** TTS service (edge-tts, JP + EN voices, audio caching)
- [ ] **FEAT-303** Vocabulary enrichment service (reading, example, translation)
- [ ] **FEAT-304** `vocabularies` + `vocabulary_reviews` + `learning_sessions` migrations
- [ ] **FEAT-305** Vocabulary CRUD endpoints
- [ ] **FEAT-306** Bulk import (CSV/JSON)
- [ ] **FEAT-307** Review queue endpoint (SM-2 due items)
- [ ] **FEAT-308** Review submission + outbox event
- [ ] **FEAT-309** Worker: idempotent SM-2 schedule update
- [ ] **FEAT-310** Learning session start/end
- [ ] **FEAT-311** Flashcard study mode UI (3D flip, keyboard)
- [ ] **FEAT-312** Multiple choice quiz UI
- [ ] **FEAT-313** Type-in practice UI (fuzzy match for kana/romaji)
- [ ] **FEAT-314** Audio quiz UI (TTS playback)
- [ ] **FEAT-315** Vocabulary list/management UI
- [ ] **FEAT-316** Learning stats + mastery heatmap

### PHASE-4: Personal Finance
- [x] **FEAT-401** `accounts`, `categories`, `transactions`, `budgets`, `savings_goals` migrations
- [x] **FEAT-402** Accounts CRUD endpoints
- [ ] **FEAT-403** Categories CRUD endpoints (with default seed data)
- [x] **FEAT-404** Transactions CRUD + filtering
- [ ] **FEAT-405** Budgets CRUD + utilization endpoint
- [ ] **FEAT-406** Savings goals CRUD
- [x] **FEAT-407** Finance analytics endpoints (dashboard, cashflow, by-category, net-worth)
- [ ] **FEAT-408** Finance dashboard UI (net worth, cashflow chart, recent transactions)
- [ ] **FEAT-409** Transaction list UI (filterable, searchable)
- [ ] **FEAT-410** Transaction quick-add form (drawer/modal)
- [ ] **FEAT-411** Budget management UI (progress bars, alerts)
- [ ] **FEAT-412** Savings goals UI
- [ ] **FEAT-413** Account management UI

### PHASE-5: Journal
- [ ] **FEAT-501** Elasticsearch setup + journal index mapping
- [ ] **FEAT-502** `journals` migration (incl. `search_vector` trigger, `linked_journal_ids`)
- [ ] **FEAT-503** Journal CRUD endpoints
- [ ] **FEAT-504** Full-text search endpoint (PG `search_vector` + date range + tags)
- [ ] **FEAT-505** Elasticsearch semantic search integration (AI service embedding)
- [ ] **FEAT-506** Journal linking endpoint
- [ ] **FEAT-507** Journal stats endpoint (streak, word count, mood trend)
- [ ] **FEAT-508** Journal editor UI (Markdown, autosave draft, preview)
- [ ] **FEAT-509** Journal list UI (calendar heatmap, search, filters)
- [ ] **FEAT-510** Journal detail/view UI (rendered Markdown, linked journals)
- [ ] **FEAT-511** Mood + energy tracking UI
- [ ] **FEAT-512** Writing streak widget

### PHASE-6: Hardening
- [ ] **FEAT-601** OpenTelemetry instrumentation
- [ ] **FEAT-602** Integration test coverage for all critical paths
- [ ] **FEAT-603** Playwright E2E tests
- [ ] **FEAT-604** Production Docker build + deployment guide
- [ ] **FEAT-605** Security review (CORS, rate limits, CSP headers)
- [ ] **FEAT-606** Performance review (bundle size, DB indexes, query plans)

---

## 14. Required AI execution protocol

Build GoPA one vertical slice at a time. **Do not generate the entire application in one response.**

For every implementation task, respond in this sequence:

### Phase A — Architecture + Go vs Java explanation

- State the deliverable and narrow scope.
- Explain component placement, dependency direction, and design decisions.
- Compare the meaningful Go choice to the closest Java/Spring equivalent.
- Identify transactions, context handling, security implications, and trade-offs.
- Reference the relevant `[FEAT-XXX]` items from Section 13.

### Phase B — Commands

Provide exact, copyable commands: `go get`, `pip install`, `npm install`, Docker, migrations, tests. State expected outcome.

### Phase C — Complete code

- Full, untruncated file contents.
- Every file labeled with its exact project-relative path.
- All imports, DTOs, migrations, tests, i18n keys, and wiring included.
- No pseudo-code, "implementation omitted," or unexplained `any` types.
- Follow directory structure and standards in this document.

### Phase D — Verification + handoff

- Commands: format, lint, typecheck, test, manual acceptance.
- Known limitations and next slice.
- Mark completed `[FEAT-XXX]` items in Section 13.
- Stop and wait for review before the next module.

---

## 15. Definition of done

A feature is done only when:

- [ ] Domain rules and ownership authorization are enforced in service layer
- [ ] Migrations, indexes, and constraints exist
- [ ] ACID compliance verified for multi-record operations
- [ ] API errors are safe, consistent, and mapped correctly
- [ ] Unit tests cover core business rules
- [ ] Integration tests cover external boundaries
- [ ] Structured logs and cancellation behavior are correct
- [ ] Frontend has loading, empty, success, and error states
- [ ] All user-visible strings translated (`vi`, `en`, `ja`)
- [ ] Motion is optional for function (reduced-motion respected)
- [ ] `gofmt`, `go vet`, `npm run lint`, `npm run typecheck`, tests pass
- [ ] Relevant `[FEAT-XXX]` items in Section 13 marked complete
- [ ] README/API docs updated

---

*This specification is the source of truth. Any architectural decision that supersedes a section must be documented here before implementation begins.*
