# GoPA — Golang Personal Assistant

> **Your personal operating system for deep work, language learning, personal finance & journaling.**

GoPA is a private, full-stack web application built as a practical vehicle for learning idiomatic Go while producing a maintainable, production-minded application. It combines a **modular Go monolith** (API + async Worker) with a **React 18 frontend** and a **Python AI microservice** to power five integrated life-management modules.

**Version:** 2.5 · **Branch:** `dev/version2.5` · **Author:** TuanHLA

---

## ✨ Features

### 🎯 Deep Work & Pomodoro
- Kanban task board with `TODO → IN_PROGRESS → DONE` columns and drag reordering
- Priority levels: `LOW`, `MEDIUM`, `HIGH`, `URGENT` (pulsing crimson indicator)
- LeetCode category chip with dedicated icon
- Pomodoro timer with **real-time WebSocket sync** — start, pause, stop, finish
- Pomodoro history log; completed session count linked to tasks

### 🧠 Linguistics — Smart SRS
- Japanese (N3 → N2) and English (TOEIC 805 → 900+) vocabulary management
- **SuperMemo-2 (SM-2) spaced-repetition algorithm** implemented as a pure domain function
- Four interactive study modes: **3D Flashcard flip**, **Multiple Choice**, **Type-in** (fuzzy kana/romaji match), **Audio Quiz**
- 90-day review heatmap (GitHub-style, violet tint) + Leitner box distribution
- CSV/JSON bulk import; AI-assisted TTS audio playback (via AI microservice)
- Mastery score (0–100) and vocabulary statistics dashboard

### 💰 Personal Finance
- Multi-account wallets: `CASH`, `BANK`, `CREDIT_CARD`, `SAVINGS`, `INVESTMENT`, `CRYPTO`
- Full transaction lifecycle: income, expense, and account transfers with multi-currency support
- Category hierarchy (max 2 levels); auto-seeded default categories on first login
- Monthly budget tracker with color-coded progress meters and alert thresholds (80%)
- Savings goals with circular SVG progress rings and target-date countdown
- Analytics: net worth, cashflow summary, spending by category, balance history

### 📓 Markdown Journal
- Rich Markdown editor with autosave draft (browser `localStorage`)
- Mood (`GREAT` → `TERRIBLE`) + energy level (1–5 lightning bolts) tracking
- Pinned notes, bi-directional inter-entry linking
- PostgreSQL full-text search (`search_vector` GIN index + `ts_rank`)
- Writing streak counter, word count statistics, mood trend visualization
- Dedicated reader view with syntax-highlighted code blocks

### 🌐 Multi-language UI
- Full i18n parity across Vietnamese (`vi`), English (`en`), and Japanese (`ja`)
- Runtime locale switching without layout shifts or broken glyphs

---

## 🏗 Architecture

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

### Process Responsibilities

| Process | Language | Role |
|---|---|---|
| **API** | Go | HTTP REST + WebSocket; validation, auth, use-case orchestration |
| **Worker** | Go | Async event consumers; SRS scheduling, Pomodoro history, budget alerts |
| **AI Service** | Python / FastAPI | TTS generation, semantic search embedding, vocabulary enrichment |
| **PostgreSQL** | — | Primary durable source of truth + full-text search vectors |
| **Elasticsearch** | — | Journal & vocabulary semantic search |
| **Redis** | — | Session state, rate limits, live Pomodoro state, Pub/Sub |
| **RabbitMQ** | — | Async decoupling; at-least-once delivery to Worker |

---

## 🛠 Tech Stack

| Layer | Technology | Version |
|---|---|---|
| **Backend language** | Go | 1.22+ |
| **HTTP framework** | `gin-gonic/gin` | v1.10+ |
| **Database access** | `jmoiron/sqlx` + `pgx` | — |
| **Database** | PostgreSQL | 16+ |
| **Search** | Elasticsearch | 8+ |
| **Cache / Pub-Sub** | Redis (`go-redis/v9`) | 7+ |
| **Message broker** | RabbitMQ (`amqp091-go`) | 3.13+ |
| **Auth** | JWT + bcrypt + opaque refresh tokens | — |
| **Logging** | `log/slog` (stdlib) | — |
| **AI microservice** | Python 3.12 + FastAPI | — |
| **TTS** | `edge-tts` (offline, natural voices) | — |
| **Embeddings** | `sentence-transformers` (multilingual) | — |
| **Frontend** | React 18 + TypeScript + Vite | — |
| **Styling** | Tailwind CSS + shadcn/ui | — |
| **State management** | TanStack Query + Zustand | — |
| **Animation** | Framer Motion (`motion/react`) | — |
| **Icons** | Lucide React | — |
| **Containerisation** | Docker Compose (dev) | — |

---

## 📁 Repository Structure

```text
GoPA/
├── be/                              # Go backend (modular monolith)
│   ├── cmd/
│   │   ├── api/main.go              # API composition root
│   │   ├── worker/main.go           # Worker composition root
│   │   └── migrate/main.go          # Migration runner
│   ├── internal/
│   │   ├── core/
│   │   │   ├── domain/              # Entities, value objects, domain errors
│   │   │   │   └── linguistics/     # Pure SM-2 SRS algorithm
│   │   │   └── ports/               # Repository, cache, broker contracts
│   │   ├── adapters/
│   │   │   ├── handler/http/        # Gin handlers + WebSocket gateway
│   │   │   ├── repository/          # PostgreSQL / sqlx implementations
│   │   │   ├── cache/               # Redis implementations
│   │   │   ├── broker/              # RabbitMQ publisher / consumer
│   │   │   └── outbox/              # Transactional outbox relay
│   │   └── services/                # Use-case orchestration + worker handlers
│   ├── pkg/
│   │   ├── config/                  # Env loading and validation
│   │   ├── database/                # Pool, WithTx helper, migration runner
│   │   ├── logger/                  # slog structured setup
│   │   └── utils/                   # JWT, password, response/error helpers
│   ├── migrations/                  # Ordered raw SQL up/down files (000001–000005)
│   ├── docker-compose.yml
│   ├── Makefile
│   └── .env.example
├── ai/                              # Python AI microservice (planned / in-progress)
│   ├── main.py                      # FastAPI application (TTS, embedding, enrichment)
│   ├── services/
│   │   ├── tts_service.py
│   │   ├── embedding_service.py
│   │   └── context_service.py
│   └── requirements.txt
└── fe/                              # React 18 frontend
    └── src/
        ├── app/                     # Providers, router, App shell
        ├── components/
        │   ├── design-system/       # DoubleBezelCard, ProgressRing, CommandPalette …
        │   ├── feedback/            # LoadingState, ErrorState
        │   └── ui/                  # shadcn primitives
        ├── features/
        │   ├── auth/                # Login, Register, Auth provider
        │   ├── today/               # Asymmetric Bento dashboard
        │   ├── tasks/               # Kanban board
        │   ├── linguistics/         # SRS review session, vocabulary management
        │   ├── finance/             # Budgets, savings goals, transactions
        │   ├── journal/             # Editor, list, detail view
        │   ├── pomodoro/            # WebSocket Pomodoro hooks
        │   └── settings/            # Profile, appearance, locale
        ├── hooks/                   # useDebounce, usePomodoroWs …
        ├── lib/                     # Axios client, API envelope, cn()
        ├── locales/                 # i18n resources (vi, en, ja)
        ├── stores/                  # Zustand UI preferences
        └── styles/                  # Global CSS, three-layer design tokens
```

---

## 🚀 Quick Start

### Prerequisites

| Tool | Version |
|---|---|
| Go | 1.22+ |
| Node.js | 20+ |
| Docker + Docker Compose | Latest stable |
| `golang-migrate` CLI | v4+ |

```bash
# Install golang-migrate (one-time)
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### 1 — Clone the repository

```bash
git clone https://github.com/TuanHLA/GoPA.git
cd GoPA
```

### 2 — Start infrastructure (PostgreSQL, Redis, RabbitMQ, Elasticsearch)

```bash
cd be
cp .env.example .env          # Edit secrets before use
docker compose up -d
```

Verify all services are healthy:

```bash
docker compose ps
```

### 3 — Run database migrations

```bash
cd be
go run ./cmd/migrate -direction up
# or: migrate -path migrations -database "postgres://gopa:gopa@localhost:5432/gopa?sslmode=disable" up
```

### 4 — Start the Go API

```bash
cd be
go mod tidy
go run ./cmd/api
# API available at http://localhost:8080
# Health: GET http://localhost:8080/healthz
# Ready:  GET http://localhost:8080/readyz
```

### 5 — Start the Go Worker (async consumers)

```bash
# In a new terminal
cd be
go run ./cmd/worker
```

### 6 — Start the React frontend

```bash
cd fe
cp .env.example .env
npm install
npm run dev
# Open http://localhost:5173
```

### 7 — (Optional) Start the Python AI Service

> ⚠️ The AI microservice (`ai/`) is currently **in progress**. TTS audio and semantic search features will degrade gracefully without it.

```bash
cd ai
pip install -r requirements.txt
uvicorn main:app --reload --port 8000
```

---

## ⚙️ Environment Variables

All backend environment variables are defined in `be/.env.example`. Key variables:

| Variable | Default | Description |
|---|---|---|
| `APP_ENV` | `development` | Runtime environment |
| `HTTP_ADDR` | `:8080` | API listen address |
| `POSTGRES_DSN` | `postgres://gopa:gopa@localhost:5432/gopa?sslmode=disable` | PostgreSQL connection string |
| `REDIS_ADDR` | `localhost:6379` | Redis address |
| `RABBITMQ_URL` | `amqp://gopa:gopa@localhost:5672/` | RabbitMQ AMQP URL |
| `ELASTICSEARCH_URL` | `http://localhost:9200` | Elasticsearch URL |
| `AI_SERVICE_URL` | `http://localhost:8000` | Python AI microservice URL |
| `AI_SERVICE_SECRET` | `change-me` | Shared secret for AI service auth |
| `JWT_ACCESS_SECRET` | *(required)* | Secret for signing access tokens |
| `JWT_ISSUER` | `gopa` | JWT issuer claim |
| `ACCESS_TOKEN_TTL` | `15m` | JWT access token lifetime |
| `REFRESH_TOKEN_TTL` | `168h` | Refresh token lifetime (7 days) |
| `WEB_ORIGIN` | `http://localhost:5173` | Allowed CORS origin |
| `LOG_LEVEL` | `debug` | Logging verbosity |
| `SHUTDOWN_TIMEOUT` | `15s` | Graceful shutdown window |

---

## 📡 API Endpoints

Base path: `/api/v1`. All responses use the unified envelope:

```json
// Success
{ "data": {}, "meta": { "request_id": "01J..." } }
// Error
{ "error": { "code": "NOT_FOUND", "message": "..." }, "meta": { "request_id": "01J..." } }
```

### IAM

| Method | Path | Description |
|---|---|---|
| `POST` | `/auth/register` | Create user account |
| `POST` | `/auth/login` | Issue JWT + refresh token |
| `POST` | `/auth/refresh` | Rotate refresh token |
| `POST` | `/auth/logout` | Revoke session |
| `GET` | `/me` | Current user profile |
| `PATCH` | `/me` | Update display name / avatar |

### Deep Work

| Method | Path | Description |
|---|---|---|
| `GET/POST` | `/tasks` | List / create tasks |
| `PATCH/DELETE` | `/tasks/:id` | Update / delete task |
| `PATCH` | `/tasks/:id/status` | Move task status + reorder |
| `GET/PUT` | `/pomodoro` | Read / replace timer state |
| `POST` | `/pomodoro/stop` | Stop timer |
| `GET` | `/pomodoro/history` | Session history |
| `GET` | `/ws/pomodoro` | Authenticated WebSocket (Pub/Sub) |

### Linguistics

| Method | Path | Description |
|---|---|---|
| `GET/POST` | `/vocabularies` | List / create vocabulary |
| `GET/PATCH/DELETE` | `/vocabularies/:id` | Manage vocabulary item |
| `POST` | `/vocabularies/import` | Bulk import (CSV / JSON) |
| `GET` | `/vocabularies/review-queue` | SM-2 due items |
| `POST` | `/vocabularies/:id/reviews` | Submit review event |
| `POST` | `/vocabularies/:id/audio` | Generate / fetch TTS audio |
| `GET` | `/vocabularies/stats` | Learning progress + heatmap |
| `POST` | `/learning-sessions` | Start learning session |
| `PATCH` | `/learning-sessions/:id` | End learning session |

### Personal Finance

| Method | Path | Description |
|---|---|---|
| `GET/POST` | `/accounts` | List / create accounts |
| `GET/PATCH/DELETE` | `/accounts/:id` | Manage account |
| `GET` | `/accounts/:id/balance-history` | Balance over time |
| `GET/POST` | `/transactions` | List / create transactions |
| `GET/PATCH/DELETE` | `/transactions/:id` | Manage transaction |
| `GET` | `/transactions/summary` | Period cashflow summary |
| `GET` | `/transactions/by-category` | Spending breakdown |
| `GET/POST` | `/categories` | List / create categories |
| `GET/POST` | `/budgets` | List / create budgets |
| `GET/PATCH/DELETE` | `/budgets/:id` | Manage budget |
| `GET` | `/budgets/status` | Utilization with alerts |
| `GET/POST` | `/savings-goals` | List / create savings goals |
| `PATCH/DELETE` | `/savings-goals/:id` | Manage savings goal |
| `GET` | `/finance/net-worth` | Net worth summary |
| `GET` | `/finance/dashboard` | Full dashboard data |

### Journal

| Method | Path | Description |
|---|---|---|
| `GET/POST` | `/journals` | List / create journals |
| `GET/PATCH/DELETE` | `/journals/:id` | Manage journal entry |
| `POST` | `/journals/:id/link` | Link to another journal |
| `DELETE` | `/journals/:id/link/:linked_id` | Remove link |
| `GET` | `/journals/search` | Full-text + date-range search |
| `GET` | `/journals/stats` | Streak, word count, mood trend |

### Health

| Method | Path | Description |
|---|---|---|
| `GET` | `/healthz` | Process liveness |
| `GET` | `/readyz` | Dependency readiness (PG, Redis, RabbitMQ) |

---

## 📊 Implementation Status

| Module | Backend | Frontend | AI Integration |
|---|---|---|---|
| **IAM** (auth, profile, rate-limiting) | ✅ Complete | ✅ Complete | — |
| **Deep Work** (Kanban, Pomodoro WS) | ✅ Complete | ✅ Complete | — |
| **Linguistics** (SM-2, 4 study modes) | ✅ Complete | ✅ Complete | ⏳ TTS (AI service pending) |
| **Finance** (accounts, txns, budgets, goals) | ✅ Complete | ✅ Complete | — |
| **Journal** (editor, search, linking) | ✅ Complete | ✅ Complete | ⏳ Semantic search (pending) |
| **Python AI Microservice** | ⏳ In Progress | — | — |
| **Elasticsearch** (semantic journal search) | ⏳ In Progress | — | — |
| **OpenTelemetry** (metrics + tracing) | 📋 Planned | — | — |
| **Playwright E2E tests** | 📋 Planned | — | — |
| **Production Docker build** | 📋 Planned | — | — |

**Legend:** ✅ Complete · ⏳ In Progress · 📋 Planned

---

## 🧪 Development Commands

### Backend (Go)

```bash
cd be

# Infrastructure
docker compose up -d          # Start PostgreSQL, Redis, RabbitMQ, Elasticsearch
docker compose down           # Stop all services

# Database
go run ./cmd/migrate -direction up    # Apply all pending migrations
go run ./cmd/migrate -direction down  # Revert one migration

# Run processes
go run ./cmd/api              # Start HTTP API (port 8080)
go run ./cmd/worker           # Start async worker

# Quality gates
gofmt -w .                    # Format all Go files
go vet ./...                  # Static analysis (target: 0 warnings)
go test -v -race ./...        # Unit + integration tests with race detector
go build ./cmd/api ./cmd/worker ./cmd/migrate   # Build all binaries
```

### Makefile Targets

```bash
make up                       # docker compose up -d
make down                     # docker compose down
make api                      # go run ./cmd/api
make worker                   # go run ./cmd/worker
make migrate-up               # Apply migrations
make migrate-down             # Revert one migration
make test                     # Unit tests
make test-integration         # Integration tests (requires Docker)
make lint                     # gofmt + go vet + golangci-lint
make build                    # Build all binaries
```

### Frontend (React)

```bash
cd fe

npm run dev                   # Vite dev server (port 5173, HMR)
npm run build                 # Production bundle (dist/)
npm run preview               # Preview production build locally
npm run typecheck             # tsc --noEmit (target: 0 errors)
npm run lint                  # ESLint (target: 0 warnings)
npm run test                  # Vitest unit tests
npm run test:ui               # Vitest UI mode
```

---

## 🤖 Python AI Microservice

The AI microservice (`ai/`) is **currently in progress**. When complete it will provide:

| Endpoint | Description |
|---|---|
| `POST /ai/tts` | Generate TTS audio for a vocabulary word (edge-tts, JP + EN voices) |
| `POST /ai/enrich-vocabulary` | Auto-generate reading, example sentence, translation, and audio |
| `POST /ai/semantic-search` | Embed query with `sentence-transformers` → Elasticsearch KNN search |

**Technology:** Python 3.12 · FastAPI · edge-tts · sentence-transformers (multilingual) · pykakasi (JP furigana) · OpenAI-compatible context generation

All AI service endpoints are authenticated with a shared secret (`AI_SERVICE_SECRET`). The Go API degrades gracefully when the service is unavailable — vocabulary entries work without audio, and search falls back to PostgreSQL full-text.

---

## 🎨 Design System

GoPA's frontend is built on a **Double-Bezel Haptic Architecture** — surfaces feel like machined physical hardware, not flat admin templates:

```text
Dark canvas (OLED charcoal #090b10 + ambient module radial glow)
  └─ Floating glass navigation rail (backdrop-blur-xl)
      └─ Translucent top bar (breadcrumb, Pomodoro widget, locale)
          └─ Double-Bezel Bento Grid (col-span-8 / col-span-4 asymmetric)
              └─ Command Palette (Ctrl+K) · Modal Dialogs · Toasts
```

**Key primitives** (all in `fe/src/components/design-system/`):

| Component | Purpose |
|---|---|
| `DoubleBezelCard` | Nested hardware-enclosure surface (outer glass + inner core) |
| `ProgressRing` | SVG circular progress with smooth stroke-dashoffset animation |
| `ButtonInButton` | Pill CTA with trailing icon in nested circular wrapper |
| `AmountDisplay` | `tabular-nums` locale-aware currency renderer (VND / USD / JPY) |
| `CommandPalette` | `Ctrl+K` global search across tasks, vocabulary, journals |
| `GlassPanel` | Reusable frosted translucent container |

**Motion:** All transitions use physics springs (`stiffness: 320, damping: 28`) via Framer Motion. Instant 0ms state changes are rejected — every interaction acknowledges itself.

---

## 🔒 Security

- Passwords hashed with bcrypt (cost ≥ 12); plaintext never stored or logged.
- Access JWT: 15-minute lifetime; refresh token: 7-day opaque value stored hashed in Redis.
- Refresh tokens rotate on every use; logout immediately revokes the active token.
- Every SQL query scopes by `user_id` — a guessed UUID reveals nothing.
- Redis-backed rate limiting on `/auth/login` and `/auth/refresh`.
- CORS: explicit origins only; no wildcard with credentials.
- Markdown rendered with sanitization (DOMPurify on frontend).
- Secrets in `.env` (git-ignored); `.env.example` contains placeholders only.

---

## 📐 Engineering Principles

1. **Correctness before cleverness** — readable, explicit code over magic.
2. **Domain-first design** — business rules never depend on Gin, sqlx, Redis, or React.
3. **ACID compliance** — every multi-record business change uses explicit `*sqlx.Tx`.
4. **Explicit over implicit** — wiring, transactions, error paths, and concurrency are visible.
5. **Secure by default** — auth and input validation are mandatory, not optional.
6. **Observable operation** — structured `slog` JSON logging with request ID on every line.
7. **Boring technology** — stable, well-understood libraries and patterns.

---

## 📚 Documentation

| Document | Purpose |
|---|---|
| [`MASTER_PROMPT.md`](./MASTER_PROMPT.md) | Root architecture specification & engineering playbook |
| [`EXECUTION_ROADMAP.md`](./EXECUTION_ROADMAP.md) | Milestone-by-milestone implementation checklist |
| [`be/BACKEND_MASTER.md`](./be/BACKEND_MASTER.md) | Backend architecture deep-dive |
| [`fe/FRONTEND_MASTER.md`](./fe/FRONTEND_MASTER.md) | Frontend UI/UX specification |

---

## 👤 Author

**TuanHLA** — A Java/Spring Boot engineer learning idiomatic Go through building a production-quality personal application.

---

*GoPA — Built with Go, React, and a lot of Pomodoro sessions.*
