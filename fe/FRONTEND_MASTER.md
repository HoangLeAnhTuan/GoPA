# GoPA Frontend — Master UI/UX Specification & Engineering Playbook

> **Status:** Living specification — update when any design or architectural decision changes.
> **Stack:** React 18+ · Vite · TypeScript (strict) · Tailwind CSS · shadcn/ui · TanStack Query · Zustand · React Router v7 · Framer Motion · react-i18next · Axios
> **Related:** `../MASTER_PROMPT.md` (root spec) · `../be/BACKEND_MASTER.md` (backend spec)

---

## 1. Product mission

GoPA frontend is a **personal operating system** — a calm, exceptionally usable daily workspace for deep work, language learning, personal finance, and journaling. It must feel like a premium, thoughtfully designed application — not an admin dashboard or a generic CRUD form collection.

The AI serves as both Principal Frontend Engineer and mentor: every implementation should be explainable to a backend-oriented developer learning React and TypeScript.

### 1.1 Core jobs to be done

- See the next meaningful action at a glance on a rich Today dashboard.
- Start, pause, and finish a Pomodoro without losing concentration.
- Review Japanese and English vocabulary with minimal interaction friction across 4 study modes.
- Track income, expenses, and savings with a clear, visual financial dashboard.
- Capture a journal entry quickly, then find it later through date, mood, tag, or keyword search.
- Move among Vietnamese, English, and Japanese without breaking the interface.

### 1.2 Experience principles

1. **Calm is productive.** Visual hierarchy and whitespace before panels and decoration.
2. **macOS-inspired aesthetics.** Glassmorphism layers, fluid animations, vibrancy — without compromising readability.
3. **Rich colors, not flat gray.** Curated color palettes, gradient accents, colorful category indicators.
4. **Information earns prominence.** Current task, active timer, review queue, budget alerts are primary; settings are quiet.
5. **Immediate feedback.** Every interaction acknowledges success, progress, or failure visibly.
6. **Accessible by construction.** Keyboard navigation, contrast, focus rings, reduced motion are requirements.
7. **Shared design system.** Feature teams compose primitives; no one-off visual dialects.

### 1.3 Non-goals

- Reproducing Apple/Microsoft interfaces literally.
- Glass blur so aggressive it hurts readability or battery.
- Generic "card grid" dashboards for every screen.
- Duplicating server data in Zustand.
- Shipping untyped API data or `any` as a shortcut.

---

## 2. Design language: Warm Vibrancy

The visual system combines warm macOS-like translucency with vibrant color accents and clear information density. It should feel like a composed personal workspace: layered, rounded, alive, and responsive.

### 2.1 Visual rules

| Element | Rule |
|---|---|
| Canvas | Deep dark background (`#0d0f14` or `slate-950`) with subtle radial gradients of primary color |
| Surfaces | Semi-opaque dark glass: `bg-white/5 backdrop-blur-xl border border-white/10` |
| Accent cards | Strong color identity per module (finance = emerald, linguistics = violet, journal = amber, tasks = blue) |
| Corners | `rounded-xl` controls · `rounded-2xl` cards · `rounded-3xl` modals |
| Borders | `border border-white/10` (dark) · subtle and low-contrast |
| Shadows | `shadow-2xl shadow-black/40` — deep, layered |
| Typography | Inter (body) + Geist (monospace for code) + Noto Sans JP (Japanese glyphs) |
| Motion | Purposeful spring animations; fade + translate; 3D flips for flashcards |
| Color | Semantic per module; gradient accents on primary actions |
| Icons | Lucide React (primary) + Heroicons (secondary) · always with accessible label |

### 2.2 Layer hierarchy

```text
Dark canvas with ambient gradient glow
  └─ floating left navigation rail (glass)
      └─ translucent top bar + global actions
          └─ content surface (semi-opaque dark card)
              └─ popovers, command palette (higher blur, brighter)
                  └─ dialogs (darkened overlay + bright card)
                      └─ toasts (fixed corner, spring animation)
```

### 2.3 Design tokens

```css
/* src/styles/tokens.css */
:root {
  /* Canvas */
  --bg-base: 222 47% 5%;            /* #0d0f14 */
  --bg-surface: 222 30% 9%;          /* dark glass surface */
  --bg-surface-hover: 222 30% 12%;

  /* Glass */
  --glass-bg: rgba(255,255,255,0.04);
  --glass-border: rgba(255,255,255,0.08);
  --glass-blur: 24px;

  /* Module accent colors */
  --accent-finance: 158 64% 52%;     /* emerald */
  --accent-linguistics: 267 83% 69%; /* violet */
  --accent-journal: 38 92% 55%;      /* amber */
  --accent-tasks: 217 91% 60%;       /* blue */
  --accent-pomodoro: 0 84% 60%;      /* red */

  /* Typography */
  --foreground: 210 40% 96%;
  --muted-foreground: 215 16% 60%;
  --border: rgba(255,255,255,0.08);

  /* Semantic */
  --success: 142 71% 45%;
  --warning: 38 92% 50%;
  --destructive: 0 84% 60%;

  /* Radius */
  --radius: 0.75rem;
  --radius-lg: 1rem;
  --radius-xl: 1.5rem;
}
```

### 2.4 Dark mode

Dark mode is the **primary** and default mode. A light mode is optional for later. System preference is respected; user override stored in Zustand UI preferences store. Verify text contrast in both modes where applicable.

### 2.5 Typography

```css
/* Global font stack */
body {
  font-family: 'Inter', system-ui, sans-serif;
  font-feature-settings: 'cv02', 'cv03', 'cv04', 'cv11';
}

/* Japanese content */
.ja-text {
  font-family: 'Noto Sans JP', 'Hiragino Sans', 'Yu Gothic', sans-serif;
}

/* Monospace (code blocks in journal) */
.mono { font-family: 'Geist Mono', 'Cascadia Code', monospace; }

/* Numeric displays (timer, amounts) */
.tabular { font-variant-numeric: tabular-nums; }
```

### 2.6 Module color identities

Each primary module has a color identity used for: gradient header backgrounds, active nav indicator, accent buttons, chart colors, and empty state illustrations.

| Module | Primary | Usage |
|---|---|---|
| Today / Dashboard | Blue + Violet gradient | Hero gradient, greeting banner |
| Deep Work | `blue-500` | Task cards, Kanban columns, Pomodoro ring |
| Linguistics | `violet-500` | Flashcard back, progress bars, mastery badges |
| Finance | `emerald-500` | Balance numbers, income, savings progress |
| Journal | `amber-500` | Editor toolbar, mood indicators, streak |
| Settings | `slate-500` | Neutral, no accent |

### 2.7 Icons and backgrounds

- **Primary icons:** Lucide React (consistent, clean, 24px default).
- **Decorative/module icons:** Heroicons (solid) for large feature callouts.
- **Background patterns:** Subtle noise texture (`bg-noise`) + radial gradient glows per module.
- **Background images:** Abstract gradient backgrounds for auth pages and empty states (generated via design tool, stored in `src/assets/`).
- **Emoji:** Never use emoji as the only status indicator. Always pair with text or icon.

---

## 3. Information architecture

### 3.1 Primary navigation

| Route | Label (vi/en) | Icon | Color |
|---|---|---|---|
| `/app/today` | Hôm nay / Today | `LayoutDashboard` | Blue-violet gradient |
| `/app/tasks` | Công việc / Tasks | `CheckSquare` | Blue |
| `/app/learn` | Học tập / Learn | `GraduationCap` | Violet |
| `/app/finance` | Tài chính / Finance | `Wallet` | Emerald |
| `/app/journal` | Nhật ký / Journal | `BookOpen` | Amber |
| `/app/settings` | Cài đặt / Settings | `Settings` | Slate |

The active Pomodoro is globally visible in the top bar — not a primary nav item.

### 3.2 Route map

```text
/
├── /login
├── /register
├── /app
│   ├── /today                          # Personal dashboard
│   ├── /tasks                          # Kanban board
│   ├── /pomodoro                       # Detailed timer view
│   ├── /learn                          # Learning hub overview
│   │   ├── /review                     # Active study session
│   │   ├── /vocabulary                 # Vocabulary list/management
│   │   └── /stats                      # Progress, heatmap, streaks
│   ├── /finance                        # Finance dashboard
│   │   ├── /transactions               # Transaction list
│   │   ├── /accounts                   # Account management
│   │   ├── /budgets                    # Budget management
│   │   └── /goals                      # Savings goals
│   ├── /journal                        # Journal list (calendar view)
│   ├── /journal/new                    # New journal entry
│   ├── /journal/:journalId             # Journal detail/edit
│   └── /settings                       # User settings
└── * (not found)
```

Use React Router v7 nested routes. Public routes use `AuthLayout`; protected routes use `AppLayout`. Route guards wait for the current-user query to settle before redirecting — no flash from private to login.

### 3.3 Page anatomy

1. Module-colored page header (gradient strip or accent border).
2. Contextual page title + description + primary action button.
3. Critical summary, filter controls, or quick stats.
4. Main working surface.
5. Secondary/supporting information below or in side panels.

---

## 4. Project structure

```text
fe/src/
├── app/
│   ├── App.tsx                       # Root providers + router
│   ├── providers.tsx                 # Query, i18n, theme, toast providers
│   └── router.tsx                    # Route definitions
├── assets/
│   ├── backgrounds/                  # Module background images/SVGs
│   └── fonts/                        # Local font files if needed
├── components/
│   ├── ui/                           # shadcn primitives (auto-generated)
│   ├── design-system/
│   │   ├── GlassPanel.tsx            # Translucent dark glass surface
│   │   ├── GradientCard.tsx          # Module-colored gradient card
│   │   ├── PageHeader.tsx            # Module header with gradient
│   │   ├── StatCard.tsx              # KPI metric card with trend
│   │   ├── DataTable.tsx             # Sortable/filterable table primitive
│   │   ├── CommandPalette.tsx        # Ctrl+K global search/actions
│   │   ├── EmptyState.tsx            # Absence + CTA
│   │   ├── ConfirmDialog.tsx         # Destructive action dialog
│   │   ├── AmountDisplay.tsx         # Locale-aware currency
│   │   ├── MoodBadge.tsx             # Mood indicator
│   │   └── ProgressRing.tsx          # Circular progress (Pomodoro, goals)
│   └── feedback/
│       ├── LoadingSkeleton.tsx       # Content-shaped skeletons
│       ├── ErrorState.tsx            # Friendly error + retry
│       └── Toast.tsx                 # Accessible notifications
├── features/
│   ├── auth/
│   │   ├── api/ components/ hooks/ pages/ types.ts
│   ├── today/
│   │   ├── components/ hooks/ pages/
│   ├── tasks/
│   │   ├── api/ components/ hooks/ pages/ types.ts
│   ├── pomodoro/
│   │   ├── api/ components/ hooks/ stores/ types.ts
│   ├── linguistics/
│   │   ├── api/ components/ hooks/ pages/ types.ts
│   │   └── components/
│   │       ├── flashcard/           # Flashcard study mode
│   │       ├── multiple-choice/     # Quiz study mode
│   │       ├── type-in/             # Type-in practice mode
│   │       └── audio-quiz/          # Audio quiz mode
│   ├── finance/
│   │   ├── api/ components/ hooks/ pages/ types.ts
│   │   └── components/
│   │       ├── charts/              # Cashflow, category charts
│   │       ├── transaction-form/    # Quick-add drawer
│   │       └── budget-tracker/      # Budget progress bars
│   └── journal/
│       ├── api/ components/ hooks/ pages/ types.ts
│       └── components/
│           ├── editor/              # Markdown editor
│           ├── calendar-view/       # Heatmap calendar
│           └── search/              # Search + filter panel
├── hooks/
│   ├── useDebounce.ts
│   ├── useKeyboardShortcut.ts
│   ├── useAudioTTS.ts               # Vocabulary TTS playback
│   ├── useWebSocket.ts              # Authenticated WS with backoff
│   ├── useInfiniteScroll.ts
│   ├── useLocalDraft.ts             # Local storage draft autosave
│   └── useMediaQuery.ts
├── layouts/
│   ├── AppLayout.tsx                # Main app shell
│   ├── AuthLayout.tsx               # Auth page wrapper
│   └── components/
│       ├── Sidebar.tsx              # Floating left nav rail
│       ├── TopBar.tsx               # Global top bar
│       └── PomodoroWidget.tsx       # Compact global timer
├── lib/
│   ├── axios.ts                     # Axios instance + interceptors
│   ├── query-client.ts              # TanStack Query client config
│   ├── cn.ts                        # clsx + tailwind-merge
│   ├── formatters.ts                # Currency, date, duration formatters
│   └── query-keys.ts                # (optional shared keys)
├── locales/
│   ├── en/ vi/ ja/                  # Namespaced JSON files per locale
├── stores/
│   ├── ui-store.ts                  # theme, sidebar, locale
│   ├── auth-store.ts                # Access token in memory
│   └── command-palette-store.ts
├── styles/
│   ├── globals.css                  # Tailwind directives + global resets
│   └── tokens.css                   # CSS custom properties
└── main.tsx
```

### 4.1 Constant Architecture & Conventions

GoPA Frontend organizes constants into two clear layers based on domain separation:

1. **Non-Domain Constants (`fe/src/constants/constants.ts`)**:
   - Contains general application environment strings (`PRODUCTION`, `DEVELOPMENT`, `DEBUG`, `TEST`), primitive utility defaults (`ONE_STRING`, `ZERO_STRING`, `ONE_INT32`, `ZERO_INT32`), API client timeouts (`DEFAULT_API_TIMEOUT`), and UI pagination defaults (`DEFAULT_PAGE_SIZE`).
   - Imported by shared utilities, libraries, and application configs.
2. **Domain Constants (`fe/src/domain/constants.ts`)**:
   - Contains domain entity states and enums matching `be/internal/core/domain` (`USER_ROLE`, `TASK_STATUS`, `TASK_PRIORITY`, `TASK_CATEGORY`, `VOCABULARY_LANGUAGE`, `JOURNAL_MOOD`, `POMODORO_STATUS`, `ACCOUNT_TYPE`, `CATEGORY_TYPE`, `TRANSACTION_TYPE`).
   - Standardized using `as const` object definitions to preserve strict TypeScript type safety across feature modules.

---

## 5. Technology and responsibility boundaries

| Tool | Use for | Do NOT use for |
|---|---|---|
| React | Rendering, composition, local state | Network cache, global singleton state |
| TypeScript | API contracts, component props, domain UI models | Escaping with `any` |
| TanStack Query | Remote data, loading/error/success, cache invalidation | Theme, sidebar, draft state |
| Zustand | Cross-route UI-only state (theme, sidebar, auth token in memory) | Backend entities, API cache duplicate |
| React Router | URL state, nested layouts, route params | Business logic |
| Axios | HTTP transport, token attachment, refresh coordination | Domain state |
| Framer Motion | Entrances, feedback, layout transitions, 3D card flips | Timing correctness for Pomodoro countdown |
| Tailwind + shadcn/ui | Consistent composition and tokens | Ad-hoc one-off utility strings |
| i18next | All user-visible strings | Data formatting (`Intl` handles that) |

---

## 6. TypeScript contract policy

### 6.1 Strictness

- Enable `strict: true`, `noUncheckedIndexedAccess: true`.
- `any` is prohibited. Use `unknown` at untrusted boundaries, then narrow.
- `interface` for extensible object contracts; `type` for unions and intersections.
- Mirror stable backend DTOs in feature-local types. Future: generate from OpenAPI.

### 6.2 API model split

```text
API response DTO → mapper → UI view model → component props
```

Example for Finance:

```ts
// features/finance/types.ts

// Wire type (mirrors backend)
export interface TransactionDto {
  id: string;
  account_id: string;
  category_id: string | null;
  type: 'INCOME' | 'EXPENSE' | 'TRANSFER';
  amount: string;           // NUMERIC from backend → string
  description: string;
  occurred_at: string;      // ISO 8601
}

// UI view model
export interface Transaction {
  id: string;
  accountId: string;
  categoryId: string | null;
  type: TransactionDto['type'];
  amount: number;           // Parsed decimal
  description: string;
  occurredAt: Date;
}

// features/finance/api/transaction-mappers.ts
export function mapTransaction(dto: TransactionDto): Transaction {
  return {
    ...
    amount: parseFloat(dto.amount),
    occurredAt: new Date(dto.occurred_at),
  };
}
```

### 6.3 API envelope

```ts
// lib/axios.ts
interface ApiResponse<T> { data: T; meta: { request_id: string; next_cursor?: string; total?: number } }
interface ApiError { code: string; message: string; details?: FieldError[]; status: number; request_id: string }
```

---

## 7. Server state, client state, mutations

### 7.1 TanStack Query rules

1. All REST reads: `useQuery` or `useSuspenseQuery`.
2. All REST writes: `useMutation`.
3. Centralized, deterministic query keys per feature:

```ts
// features/finance/api/finance-keys.ts
export const financeKeys = {
  all: ['finance'] as const,
  accounts: () => [...financeKeys.all, 'accounts'] as const,
  transactions: (filters: TxFilters) => [...financeKeys.all, 'transactions', filters] as const,
  budgets: () => [...financeKeys.all, 'budgets'] as const,
  dashboard: () => [...financeKeys.all, 'dashboard'] as const,
};
```

4. Mutations invalidate exactly the affected keys.
5. UI renders all four states: loading (skeleton), empty, error, success.
6. Never `useEffect` for data fetching.

### 7.2 Zustand policy

Zustand stores only UI-only state:

```ts
// stores/ui-store.ts
interface UIState {
  theme: 'dark' | 'light';
  sidebarCollapsed: boolean;
  locale: 'vi' | 'en' | 'ja';
  setTheme: (t: UIState['theme']) => void;
  toggleSidebar: () => void;
  setLocale: (l: UIState['locale']) => void;
}
```

Must NOT store: users, tasks, vocabulary, journals, transactions, access tokens, or any TanStack Query duplicate.

### 7.3 Forms

React Hook Form + Zod for all non-trivial forms:

```ts
const schema = z.object({
  amount: z.number().positive(),
  description: z.string().min(1).max(200),
  occurred_at: z.string().datetime(),
});

type FormValues = z.infer<typeof schema>;
```

Every form provides: visible labels, inline validation, disabled pending state, server field error mapping, unsaved-change warning for journals, safe cancel/reset, focus to first invalid field.

---

## 8. Authentication and HTTP transport

### 8.1 Credential model

- Access token: in-memory only (Zustand `auth-store.ts`). Never `localStorage`.
- Refresh token: `HttpOnly`, `Secure`, `SameSite=Strict` cookie — JS cannot read it.
- On boot: call `/me` or trigger a silent refresh to restore session before rendering protected routes.
- On logout: call `/auth/logout`, clear auth store, clear query cache, close WS connections, navigate to `/login`.

### 8.2 Axios interceptor behavior

```ts
// lib/axios.ts — interceptor pseudocode

// Request: attach access token from memory
config.headers.Authorization = `Bearer ${getAccessToken()}`;

// Response 401:
// 1. If not a refresh/login call → attempt one silent refresh POST /auth/refresh
// 2. Multiple concurrent 401s → queue behind a single refresh promise
// 3. Success → update access token in memory, retry original request
// 4. Failure → clear auth state, redirect to /login
```

### 8.3 Request standards

- Set a 30s HTTP timeout.
- Surface server `request_id` in UI error messages for support.
- Use `AbortSignal` where TanStack Query supports cancellation.
- Distinguish offline/network errors from 4xx/5xx errors with different UI messages.

---

## 9. Internationalization

| Code | Language | Notes |
|---|---|---|
| `vi` | Vietnamese | Default for primary user |
| `en` | English | Technical fallback |
| `ja` | Japanese | Use natural Japanese; Noto Sans JP for glyphs |

Namespace structure:

```text
locales/
├── en/
│   ├── common.json         # Shared: errors, actions, navigation
│   ├── auth.json
│   ├── tasks.json
│   ├── learning.json
│   ├── finance.json
│   └── journal.json
├── vi/ ...                 # Mirror of en/ structure
└── ja/ ...                 # Mirror of en/ structure
```

- Keys are stable semantic paths: `finance.dashboard.net_worth_label`.
- Never concatenate translated fragments for grammar; use interpolation.
- Use `Intl.DateTimeFormat`, `Intl.NumberFormat` for locale-aware formatting.
- New UI is not done until all 3 locale files are updated.
- For Japanese vocabulary content: provide meaning + Romaji in instructional UI where it assists the learner.

---

## 10. Global layout and shared components

### 10.1 App layout (desktop)

```text
┌─────────────────────────────────────────────────────────────┐
│  Glass sidebar (72px collapsed / 220px expanded)  │ TopBar  │
│  Module icon + label                              │ Ctrl+K  │
│  Active: left accent border + module color glow   │ Lang    │
│                                                   │ Avatar  │
│  Bottom: Theme toggle                             │ Timer   │
├───────────────────────────────────────────────────┴─────────┤
│  Page content area                                          │
│  Module gradient header strip                               │
│  Main content                                               │
│                                              Pomodoro float │
└─────────────────────────────────────────────────────────────┘
```

- Sidebar: collapsible left rail with glass surface, module accent glow on active item.
- Top bar: translucent glass strip with breadcrumbs, command palette trigger, language switcher, theme toggle, compact Pomodoro status, user avatar.
- Floating Pomodoro: bottom-right minimal widget when in other modules.
- Mobile: bottom navigation tab bar (5 items) + sheet-based sidebar.

### 10.2 Shared design system primitives

```tsx
// GlassPanel — the core surface primitive
<GlassPanel className="p-6">...</GlassPanel>
// Renders: bg-white/5 backdrop-blur-xl border border-white/10 rounded-2xl

// GradientCard — module-colored card
<GradientCard module="finance" title="Net Worth" value="$12,450">
// Renders: gradient from module accent color, large number, trend indicator

// PageHeader — consistent page hero
<PageHeader
  module="finance"
  title="Tài chính"
  description="Theo dõi thu chi và mục tiêu tài chính"
  action={<Button>+ Giao dịch</Button>}
/>

// StatCard — compact KPI metric
<StatCard label="Tổng chi tiêu" value={850000} currency="VND" trend={-5.2} />

// AmountDisplay — locale-aware currency
<AmountDisplay amount={1250000} currency="VND" colorize />
// Positive: emerald, Negative: red

// ProgressRing — circular progress
<ProgressRing value={65} max={100} size={80} color="var(--accent-finance)" />
```

### 10.3 Dropdown/select improvements

Dropdowns must follow macOS-style aesthetics:

```tsx
// Use shadcn Select + custom styling
<Select>
  <SelectTrigger className="
    bg-white/5 border-white/10 backdrop-blur-sm
    hover:bg-white/10 transition-colors
    rounded-xl text-sm
  ">
    <SelectValue />
  </SelectTrigger>
  <SelectContent className="
    bg-slate-900/95 backdrop-blur-xl
    border border-white/10 rounded-xl
    shadow-2xl shadow-black/50
    p-1
  ">
    {options.map(opt => (
      <SelectItem
        key={opt.value}
        value={opt.value}
        className="rounded-lg hover:bg-white/10 cursor-pointer"
      >
        <div className="flex items-center gap-2">
          {opt.icon && <opt.icon className="size-4 text-muted-foreground" />}
          <span>{opt.label}</span>
        </div>
      </SelectItem>
    ))}
  </SelectContent>
</Select>
```

All dropdowns/selects must have: glass dark surface, rounded items, icon+label pairs where appropriate, smooth animation via Framer Motion.

---

## 11. Feature specifications

### 11.1 Authentication

**Visual:** A centered dark glass card above a deep background with a vibrant multi-color radial gradient mesh (generated abstract art). The GoPA logo / wordmark with a subtle gradient. Language and theme controls accessible before login.

**Flows:**

- Login: email + password, inline validation, loading state with spinner, generic invalid-credential error (no account enumeration hint).
- Register: display_name, email, password (strength indicator), password confirm.
- Session restoration: quiet dark loading shell while `/me` resolves — no flash.
- Session expiry: clear toast notification → redirect to login with `redirect` param.

### 11.2 Today — Personal Dashboard

The most important page. Goal: see what matters right now without scrolling.

**Layout (4-zone):**

```text
┌──────────────────────────────────────────────────────────┐
│  Greeting + Date          │  Pomodoro Active Card        │
│  "Chào buổi sáng, Hoang" │  [====........] 12:34        │
├──────────────────────────┼───────────────────────────────┤
│  Next Up (top 3 tasks)   │  Review Queue Badge           │
│  + Add Task button       │  "15 từ cần ôn hôm nay"      │
├──────────────────────────┼───────────────────────────────┤
│  Finance Summary         │  Journal Streak               │
│  Net Worth + Today Spend │  "🔥 7 ngày liên tiếp"       │
└──────────────────────────┴───────────────────────────────┘
```

All zones are glass cards. Active Pomodoro zone is primary; if no timer, show a "Bắt đầu tập trung" CTA. Each zone links to its full feature page.

### 11.3 Deep Work — Tasks and Pomodoro

**Task board:** Three responsive Kanban columns (TODO, IN_PROGRESS, DONE) with module color header. Cards show: priority badge (colored dot), title, category chip, due date, pomodoro count. Drag-and-drop with keyboard fallback (action menu to move).

**Task priority visual:**

- URGENT: red indicator + pulsing dot animation
- HIGH: orange indicator
- MEDIUM: blue indicator
- LOW: slate indicator

**Pomodoro widget:**

- Global top bar: compact ring + time remaining + task name.
- Detail view (`/pomodoro`): large circular progress ring (blue), session count dots, linked task, timer controls.
- Break mode: ring color changes to green.
- Reconnecting state: pulsing ring border + "Đang kết nối lại..." text.

**Keyboard shortcuts:**

- `Space`: Start/pause Pomodoro (when Pomodoro widget focused).
- `S`: Stop Pomodoro.

### 11.4 Linguistics — Smart Learning Hub

**Overview page (`/learn`):**

- Language selector tabs: 🇯🇵 日本語 | 🇬🇧 English
- Daily goal progress (e.g., "12/20 từ hôm nay")
- Mastery heatmap (last 90 days)
- Quick-start study session button (prominent, violet gradient)
- Study mode selector (4 modes as cards with icons)

**Vocabulary list (`/learn/vocabulary`):**

- Filterable by: language, box (1–5), difficulty, due/all, tags.
- Dense table: word (large), reading (smaller, muted), meaning, box badge, next review date, actions.
- Inline "Play Audio" button (speaker icon) → calls TTS.
- Batch actions: import CSV, delete selected.

**Study session (`/learn/review`):**

Full-screen, distraction-free mode. Exit button top-left.

**Flashcard mode:**

```
┌─────────────────────────────────────────────────────┐
│  Progress: 5 / 20           [X] Exit                │
│                                                     │
│  ┌─────────────── CARD (3D flip) ───────────────┐  │
│  │  Front: Word + Reading                       │  │
│  │  "食べる" · たべる                            │  │
│  │                                              │  │
│  │  [Speaker icon] Audio                        │  │
│  └──────────────────────────────────────────────┘  │
│                                                     │
│  [Flip Card]  or  Space / Enter                     │
│                                                     │
│  ── After flip ──                                   │
│                                                     │
│  Meaning: "to eat" · Example: 毎日ご飯を食べる。    │
│                                                     │
│  [ 1 Again ]  [ 2 Hard ]  [ 3 Good ]  [ 4 Easy ]   │
│     (red)       (orange)    (green)     (blue)      │
└─────────────────────────────────────────────────────┘
```

3D flip animation (Framer Motion, `rotateY`). Reduced motion: instant content swap. Keyboard: `Space/Enter` to flip; `1/2/3/4` for rating; `→` to skip.

**Multiple choice mode:**

Word + reading shown. 4 option buttons (glass cards). Correct answer reveals green glow + ✓. Wrong answer reveals red shake + correct answer. Brief pause then auto-advance.

**Type-in mode:**

Word shown. Text input for meaning (or reading for JP). Submit on Enter. Fuzzy match tolerance: slight typo acceptance. Result: ✓ (green) or ✗ (red) + expected answer revealed.

**Audio quiz mode:**

Speaker auto-plays. User sees 4 meaning options. Tap to answer. "Replay audio" button. Same scoring as multiple choice.

**Session end screen:**

Accuracy %, time taken, XP/mastery gained, short motivational message. CTA: "Ôn thêm" or "Về tổng quan".

**Progress stats (`/learn/stats`):**

- 90-day review heatmap (GitHub-style, violet shades).
- Box distribution bar chart (how many in each Leitner box).
- Daily/weekly review count line chart.
- Language breakdown (JP vs EN).
- Total words mastered vs in-progress.

### 11.5 Personal Finance

**Dashboard (`/finance`):**

```text
┌─────────────────────────────────────────────────────────┐
│  Net Worth: 125,450,000 ₫        [+2.3% this month]    │  ← Large emerald number
├──────────────────────────────────────────────────────────┤
│  This Month:                                            │
│  Income: +8,500,000 ₫   Expenses: -5,230,000 ₫        │
│  Cashflow chart (area chart, 30 days)                   │
├──────────────────────────────────────────────────────────┤
│  Spending by Category        │  Savings Goals           │
│  (horizontal bar chart)      │  (progress rings)        │
├──────────────────────────────────────────────────────────┤
│  Recent Transactions (last 10, with quick-add button)   │
└─────────────────────────────────────────────────────────┘
```

**Transaction quick-add (drawer from right):**

```
┌─────────────────────────────┐
│  Thêm giao dịch             │
│                             │
│  Type: [Income/Expense/Transfer] (pill toggle)  │
│  Amount: [  1,500,000  ₫ ]  │  ← Large input, numeric keyboard
│  Account: [Dropdown]        │
│  Category: [Dropdown + icon]│
│  Description: [Input]       │
│  Date: [Date picker]        │
│                             │
│  [Hủy]    [Lưu giao dịch]  │
└─────────────────────────────┘
```

Drawer animation: slide-in from right (Framer Motion).

**Transaction list (`/finance/transactions`):**

- Top filters: date range picker, account filter, category filter, type toggle.
- Grouped by day (date separator rows).
- Each row: category colored icon, description, account name, amount (+green / -red), actions.
- Infinite scroll or pagination.
- Export CSV button.

**Budget tracker (`/finance/budgets`):**

Each budget card:

```
┌────────────────────────────────────────┐
│  🍔 Ăn uống             75% used      │
│  ████████████░░░░  3,750k / 5,000k ₫  │
│  12 ngày còn lại                       │
│  [⚠️ Sắp đến giới hạn]                 │  ← Alert badge
└────────────────────────────────────────┘
```

Progress bar color: green < 60%, yellow 60–85%, red > 85%.

**Savings goals (`/finance/goals`):**

Each goal card shows: goal name + icon, target amount, current amount, % ring progress, days remaining, linked account balance if any.

**Account management (`/finance/accounts`):**

Cards per account with: colored icon, name, type badge, balance (large), last activity date. Click to see balance history chart.

### 11.6 Journal — Daily Knowledge Base

**Journal list (`/journal`):**

Default view: **Calendar view** (month grid, dots on days with entries, mood colors).

Toggle to **List view**: reverse-chronological with date separators, title, word count, mood badge, tags, preview snippet.

**Search panel (slide-in from right):**

- Search input (full-text search as you type, debounced 300ms).
- Date range picker (from/to).
- Mood filter (icon buttons: 😄😊😐😔😢).
- Tag filter (multi-select chips).
- Sort by: newest, oldest, word count.

Results show instantly as filters change. Highlight search terms in results.

**Journal editor (`/journal/new`, `/journal/:id`):**

```
┌──────────────────────────────────────────────────────────┐
│  ← Back    [Title input]                  [Lưu] [Preview]│
├──────────────────────────────────────────────────────────┤
│  Date: [2026-08-09]   Mood: [😊]  Energy: [●●●○○]       │
│  Tags: [#golang] [#productivity] [+ Add tag]             │
├──────────────────────────────────────────────────────────┤
│  Markdown toolbar: B I ~ H1 H2 H3 — link img code table  │
│                                                          │
│  Markdown textarea (full height, monospace font)        │
│                                                          │
│  Word count: 234 words · Auto-saved 2 min ago           │
├──────────────────────────────────────────────────────────┤
│  Linked journals: [Journal A] [+ Link another]          │
└──────────────────────────────────────────────────────────┘
```

Preview mode: side-by-side (desktop) or full-screen toggle (mobile). Rendered Markdown: sanitized HTML, code syntax highlighting (highlight.js or Prism), readable typography.

Autosave: saves to `localStorage` draft every 10 seconds. "Draft not saved to server" indicator. Explicit "Lưu" button saves to server.

**Journal detail (`/journal/:id`):**

Read-only rendered view. Top: title, date, mood, tags, word count, linked journals. Bottom: "Chỉnh sửa" button. Backlinks: shows other journals that link to this one.

**Writing streak widget (Today page + Journal):**

- Current streak: "🔥 7 ngày" in amber.
- Calendar dot view (last 30 days, amber for journaled days).
- CTA if no journal today: "Viết nhật ký hôm nay".

---

## 12. Motion and interaction standards

| Interaction | Motion | Duration |
|---|---|---|
| Route transition | Fade + 6px vertical settle | 200ms |
| Dialog open | Opacity + spring scale 0.95→1 | 150ms |
| Drawer slide | Slide from right, spring | 250ms |
| Card hover | Subtle brightness + 2px Y lift | 100ms |
| Toast | Slide from bottom-right + fade out | 200ms / 3s auto |
| Flashcard flip | 3D rotateY 180°, spring | 400ms |
| Budget bar fill | Width tween on mount | 600ms, ease-out |
| Kanban move | Layout animation | 200ms |
| Number count-up | Tween on mount (finance values) | 800ms, ease-out |

- Default spring: `{ type: 'spring', stiffness: 300, damping: 30 }`.
- Respect `prefers-reduced-motion`: use opacity-only or instant transitions.
- Never make timer logic depend on animation frame accuracy.
- Micro-interactions: button press scale `0.97`, icon hover rotate, active nav item glow pulse.

---

## 13. Responsiveness and accessibility

### 13.1 Breakpoints

- Desktop: ≥1024px — sidebar + top bar + full content.
- Tablet: 768–1023px — collapsed sidebar + reduced padding.
- Mobile: <768px — bottom tab nav + drawer sidebar.

Touch targets: minimum 44×44px. No hover-only essential actions.

### 13.2 Accessibility checklist

1. Semantic landmarks: `header`, `nav`, `main`, proper heading order.
2. Keyboard accessible: Tab, Enter, Space, Arrow keys for all interactive elements.
3. Visible `:focus-visible` ring (white, 2px, offset).
4. Dialogs/drawers: focus trap, restore focus on close, `aria-modal`.
5. Form fields: visible labels, `aria-describedby` for error messages.
6. Status changes: `aria-live="polite"` for dynamic content.
7. WCAG AA contrast: test text on glass surfaces in both dark context.
8. Color never the only status indicator (icon + text + color).
9. Reduced motion fallbacks for all Framer Motion animations.
10. Flashcard keyboard equivalents: always work without animation.

---

## 14. Performance

- Route-level `React.lazy` + `Suspense` for all feature bundles.
- Virtualized lists only when actual data warrants (>100 items visible).
- Images: sized, lazy-loaded, WebP format, `srcset`.
- Fonts: `font-display: swap`, preload critical font files.
- Debounce: 300ms for search inputs.
- TanStack Query `staleTime`: calibrate per endpoint (dashboard: 30s, vocabulary list: 60s, transactions: 10s).
- Measure before memoizing: `useMemo`/`useCallback` only when profiler shows real cost.

---

## 15. Error, loading, empty, and offline states

Every async surface must have all four:

| State | Required behavior |
|---|---|
| Loading | Skeleton shaped exactly like the final content |
| Empty | Explain why + primary action CTA with friendly illustration |
| Error | Human message + retry button + request_id for support |
| Success | Updated UI + brief toast only when action not self-evident |

Offline: stale data labeled as "Đã lưu trong bộ nhớ" with timestamp. Never claim a mutation succeeded until server confirms.

---

## 16. Testing

| Layer | Tooling | Focus |
|---|---|---|
| Unit | Vitest | Formatters, mappers, reducers, SM-2 display logic |
| Component | React Testing Library | Semantic rendering, keyboard behavior, state variants |
| API mocking | MSW | Loading, error, success, mutation cases |
| E2E | Playwright | Login, task create, vocabulary review, transaction add, journal search |

Test by observable behavior. Prefer `getByRole`, `getByLabelText`. `data-testid` only as last resort.

Quality gates before merge:

```bash
npm run lint
npm run typecheck
npm run test
npm run build
```

---

## 17. Security and privacy (client)

1. Never log: tokens, passwords, journal content, personal data.
2. Treat all API/WS input as untrusted; validate/narrow before use.
3. Render Markdown only through DOMPurify or a trusted sanitizer.
4. Avoid `dangerouslySetInnerHTML`; when essential, wrap in a component that always sanitizes.
5. Client-side auth hides unavailable actions (UX only); backend enforces.
6. `VITE_` variables are browser-visible; never put secrets there.
7. Redact sensitive data from error trackers (Sentry etc.).

---

## 18. Implementation roadmap & progress tracking

> **For AI Agents:** Update status as each item is completed. This is the source of truth for frontend progress.

### PHASE-0: Foundation
- [ ] **FE-001** Vite + React + TypeScript strict init
- [ ] **FE-002** Tailwind CSS + shadcn/ui + token system
- [ ] **FE-003** i18n setup (vi/en/ja + namespaced files)
- [ ] **FE-004** Axios instance + interceptors + auth-store
- [ ] **FE-005** TanStack Query client + global error handling
- [ ] **FE-006** Zustand stores (ui-store, auth-store)
- [ ] **FE-007** AppLayout (sidebar + topbar + glass surfaces)
- [ ] **FE-008** Mobile layout (bottom nav + drawer)
- [ ] **FE-009** Design system components (GlassPanel, PageHeader, StatCard, EmptyState, etc.)
- [ ] **FE-010** Toast notification system

### PHASE-1: Auth
- [ ] **FE-101** Auth layout + background design
- [ ] **FE-102** Login page + form + interceptor hook-up
- [ ] **FE-103** Register page + password strength
- [ ] **FE-104** Protected route guard
- [ ] **FE-105** Session restoration on boot

### PHASE-2: Deep Work
- [ ] **FE-201** Today dashboard page (4-zone layout)
- [ ] **FE-202** Tasks page (Kanban board)
- [ ] **FE-203** Task form (create/edit)
- [ ] **FE-204** Pomodoro detail page + global widget
- [ ] **FE-205** WebSocket hook (useWebSocket with backoff)
- [ ] **FE-206** Pomodoro session history

### PHASE-3: Linguistics
- [ ] **FE-301** Learn hub overview page + language tabs
- [ ] **FE-302** Vocabulary list/management page
- [ ] **FE-303** Flashcard study mode (3D flip + keyboard)
- [ ] **FE-304** Multiple choice study mode
- [ ] **FE-305** Type-in practice mode (fuzzy match display)
- [ ] **FE-306** Audio quiz mode (TTS hook)
- [ ] **FE-307** Session end summary screen
- [ ] **FE-308** Learning stats page (heatmap, charts)
- [ ] **FE-309** Bulk import CSV UI

### PHASE-4: Finance
- [ ] **FE-401** Finance dashboard page (net worth, cashflow chart)
- [ ] **FE-402** Transaction list page (grouped, filterable)
- [ ] **FE-403** Transaction quick-add drawer
- [ ] **FE-404** Budget tracker page (progress bars)
- [ ] **FE-405** Savings goals page (rings)
- [ ] **FE-406** Account management page
- [ ] **FE-407** Category management

### PHASE-5: Journal
- [ ] **FE-501** Journal list page (calendar heatmap view)
- [ ] **FE-502** Journal list view toggle + search panel
- [ ] **FE-503** Journal editor (Markdown + toolbar + autosave)
- [ ] **FE-504** Journal detail/read view
- [ ] **FE-505** Mood/energy picker component
- [ ] **FE-506** Journal linking UI
- [ ] **FE-507** Writing streak widget
- [ ] **FE-508** Journal search (full-text + semantic result display)

### PHASE-6: Polish
- [ ] **FE-601** Command palette (Ctrl+K)
- [ ] **FE-602** Keyboard shortcuts global registration
- [ ] **FE-603** Responsive review (all pages on mobile)
- [ ] **FE-604** Playwright E2E tests for critical paths
- [ ] **FE-605** Performance audit (bundle size, LCP)
- [ ] **FE-606** Accessibility audit

---

## 19. Required AI execution protocol

Build GoPA frontend one vertical slice at a time. **Do not generate the entire application in one response.**

For each task, respond using this format:

### Phase A — Design and logic explanation

- State the deliverable and exact scope.
- Explain how the UI applies the Warm Vibrancy design system.
- Explain component, type, query/state, and routing responsibilities.
- Compare the React pattern to a Java/Spring concept where useful.
- State assumptions, accessibility, and security considerations.
- Reference relevant `[FE-XXX]` items.

### Phase B — Commands

Complete, copyable commands: npm install, shadcn component generation, tests, build check.

### Phase C — Complete code

- Full, untruncated files with exact project-relative paths.
- All types, translations, tests, styles, and wiring for the slice.
- No `any`, unexplained placeholders, or omitted imports.
- Follow the project structure defined in this document.

### Phase D — Verification + handoff

- `lint`, `typecheck`, `test`, `build`, manual verification commands.
- Mark completed `[FE-XXX]` items.
- State next smallest slice.
- Stop for review before next module.

---

## 20. Definition of done (frontend)

A frontend feature is done only when:

- [ ] Clear route/component ownership with strict TypeScript contracts
- [ ] Server data managed by TanStack Query; UI state by Zustand only where justified
- [ ] All four states implemented: loading (skeleton), empty, error, success
- [ ] Keyboard accessible; visible focus; WCAG AA contrast
- [ ] Responsive: desktop + mobile layouts both work
- [ ] Framer Motion reduced-motion fallback exists
- [ ] All user-visible strings translated in `vi`, `en`, `ja`
- [ ] Expired session handled safely
- [ ] Aligns with Warm Vibrancy design system (no one-off styles)
- [ ] `lint`, `typecheck`, `test`, `build` all pass
- [ ] Relevant `[FE-XXX]` items marked complete

---

*This document is the frontend source of truth. Any design or product decision that supersedes a section must be documented here before implementation begins.*
