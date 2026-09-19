# Changelog

## [2.5.0] - 2026-09-20

### Summary
Ponytail audit (lazy senior dev pass) + Code Review + Security Audit + API Design review
across the full stack (Go backend + React/TS frontend). Applied all Critical and Required findings.

---

### 🔒 Security Fixes

#### Backend
- **`middleware.go` — CORS wildcard fix (Critical)**
  `WEB_ORIGIN=*` now correctly sends `Access-Control-Allow-Origin: *`.
  Previously, wildcard mode matched `requestOrigin == "*"` which never matches a real browser
  Origin header — all CORS preflight requests silently failed in dev.
  Note: wildcard cannot be combined with `Allow-Credentials` per the CORS spec (documented in code).

- **`ws/pomodoro_ws.go` — WebSocket `checkOrigin` rejects empty Origin (High)**
  Non-browser clients (curl, server-side scripts) can omit the `Origin` header entirely,
  bypassing origin checking. Now returns `false` for empty Origin in non-wildcard mode.

#### Frontend
- **JWT in WebSocket URL** — tracked as known risk (backend coordination required).
  The `?token=<jwt>` query param exposes the access token in proxy logs and browser history.
  Fix deferred: requires backend to support first-message auth pattern before FE can remove it.

---

### 🐛 Correctness Fixes

#### Backend
- **`auth_middleware.go` — Duplicate WS query-param block (Ponytail / Required)**
  Two consecutive `if tokenStr == "" && path == "/ws/pomodoro"` blocks merged into one.
  Second block was dead code when the first already found a token.

- **`auth_middleware.go` — Removed duplicate identity context key (Architecture / Required)**
  `c.Set("user_id", userID)` was a duplicate of `c.Set(constants.IdentityKey, Identity{...})`.
  Both carried the same UUID. Deleted the raw `"user_id"` key.

- **`ws/pomodoro_ws.go` — Identity lookup via `constants.IdentityKey` (Architecture / Required)**
  WS handler was reading the now-deleted raw `"user_id"` key. Updated to read
  `constants.IdentityKey` and type-assert through JSON round-trip (avoids circular import).

- **`ws/pomodoro_ws.go` — Real WebSocket ping replaces JSON heartbeat (Correctness / Required)**
  The ticker was sending a JSON `{"type":"heartbeat"}` message — not a WebSocket protocol ping.
  The `SetPongHandler` (which resets the read deadline) only fires on protocol-level pong frames,
  so the 60s read deadline was never being reset by the heartbeat. Now sends `PingMessage`.

- **`config.go` — `optionalInt` parse error returns fallback not 0 (Ponytail)**
  Bad env var value (e.g. `BCRYPT_COST=abc`) previously returned `0`, which then failed
  `Validate()` with a cryptic range error instead of using the documented default.

#### Frontend
- **`review-session-page.tsx` — Stale closure in keyboard handler (Critical)**
  `handleGrade` and `handlePlayAudio` were plain `async` functions captured by the `keydown`
  `useEffect` without being in the deps array (suppressed with `eslint-disable`).
  If `currentCard` or `reviewMutation` changed mid-session, the handler graded the wrong card.
  Fixed: both functions wrapped in `useCallback` with correct deps; `eslint-disable` removed;
  both added to the effect deps array.

---

### 🌐 API Design Fixes

#### Backend
- **`response.go` — `ErrBadGateway` error code mismatch (Required)**
  Error code was `"SERVICE_UNAVAILABLE"` (implies 503) but HTTP status was 502.
  Fixed: renamed code to `"BAD_GATEWAY"` to match the HTTP status.

---

### ♿ Accessibility / WCAG AA Fixes

#### Frontend
- **`review-session-page.tsx` — Rating buttons below 44px touch target**
  Four SM-2 recall rating buttons had `py-2.5` with no minimum height.
  Added `min-h-11` (44px) to all four buttons.

- **`goals-page.tsx` — Edit/Delete action buttons below 44px touch target**
  Buttons were `h-7 w-7` (28px square), well below the WCAG AA 44×44px minimum.
  Fixed to `size-11` (44px square) with `rounded-xl`.

- **`journal-page.tsx` — "New Entry" button below 44px touch target**
  Button had `h-10` (40px). Changed to `min-h-11` (44px minimum).

---

### 🛠 UX Fixes

#### Frontend
- **`review-session-page.tsx` — `<ErrorState>` without `onRetry`**
  Users had no in-app recovery from a failed queue fetch (required page reload).
  Added `onRetry={() => void queueQuery.refetch()}`.

- **`journal-page.tsx` — `<ErrorState>` without `onRetry`**
  Added `onRetry={() => void journalsQ.refetch()}`.

- **`goals-page.tsx` — `<ErrorState>` without `onRetry`**
  Added `onRetry` that refetches both `goalsQuery` and `accountsQuery`.

---

### 🐴 Ponytail (Code Simplification)

#### Backend
- `go mod tidy` — promotes `gorilla/websocket` from `// indirect` to direct dependency.

#### Frontend
- **`tasks-page.tsx`** — Arrow function parameter `t` (in `.filter()`) shadowed the outer
  `useTranslation` `t` function. Renamed to `task` to eliminate the silent shadowing.

---

### 📋 Known Remaining Issues (Not Fixed — Require Further Decisions)

| Issue | File | Reason Deferred |
|-------|------|-----------------|
| JWT in WebSocket URL query param | `use-pomodoro-ws.ts:27` | Backend coordination required (first-message auth or httpOnly cookie) |
| CommandPalette always-active queries | `command-palette.tsx:52-55` | Requires component split — 3 live queries regardless of palette open state |
| `confirm()` in goals delete | `goals-page.tsx:228` | Should be replaced with inline confirmation state (design decision) |
| `auth-api.ts` module-level singleton | `auth-api.ts` | Interceptor setup should move to explicit `setupAuthInterceptors()` call in `main.tsx` |
| `review-session-page.tsx` `Math.random()` in `useMemo` | `review-session-page.tsx:80-84` | Impure memo — use seeded RNG or compute outside React state |
