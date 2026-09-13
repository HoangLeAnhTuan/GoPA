# GoPA Frontend — Master UI/UX Specification & Engineering Playbook

> **Status:** Living specification — source of truth for all frontend architecture, visual design, and agent implementation.
> **Core Stack:** React 18+ · Vite · TypeScript (strict) · Tailwind CSS · shadcn/ui · TanStack Query · Zustand · React Router v7 · Framer Motion (`motion/react`) · react-i18next · Axios
> **Integrated Agent Skills:** Powered by `.agents/skills/` (`ui-ux-pro-max`, `design-taste-frontend`, `high-end-visual-design`, `ui-styling`, `design-system`, `full-output-enforcement`, `redesign-existing-projects`, `gpt-taste`, `minimalist-ui`, `industrial-brutalist-ui`)
> **Related Documents:** `../MASTER_PROMPT.md` (root spec) · `../AGENTS.md` (agent operations) · `../be/BACKEND_MASTER.md` (backend spec)

---

## 1. Product Mission & Philosophy

GoPA frontend is a **personal operating system** — a calm, exceptionally usable daily workspace for deep work, language learning, personal finance, and journaling. It must feel like an elite $150k+ agency digital build: layered, tactile, alive, and mathematically balanced — never an admin dashboard, a generic CRUD form collection, or a cookie-cutter AI template.

The AI serves as both **Principal Frontend Engineer, Elite UI/UX Architect, and Mentor**: every implementation should be explainable to a backend-oriented developer learning React and TypeScript.

### 1.1 Core Jobs to Be Done

- **Today Dashboard:** See the next meaningful action at a glance on a rich, asymmetric Bento dashboard.
- **Deep Work & Pomodoro:** Start, pause, and finish Pomodoro cycles without losing concentration, with real-time WebSocket synchronization.
- **Linguistics SRS:** Review Japanese (N3 → N2) and English (TOEIC 805 → 900+) vocabulary with zero friction across 4 interactive study modes (3D flashcards, quiz, type-in, audio).
- **Personal Finance:** Track net worth, cashflow, budgets, and savings goals with high data density, tabular figures, and clear visual charts.
- **Markdown Journal:** Capture reflections quickly in an editorial-grade Markdown editor, exploring entries via calendar heatmap, mood badges, and full-text/semantic search.
- **Multi-lingual Mobility:** Switch seamlessly among Vietnamese (`vi`), English (`en`), and Japanese (`ja`) without layout shifts or broken glyphs.

### 1.2 Experience Principles

1. **Calm is productive.** Visual hierarchy, generous macro-whitespace, and tactile surfaces over noisy decoration.
2. **Haptic Depth (Double-Bezel Architecture).** Primary containers feel like machined physical hardware — a frosted glass plate nested inside a CNC aluminum tray.
3. **Rich colors with intent, not flat gray.** Tailored dark backgrounds, tinted colored shadows, and distinct module color identities.
4. **Information earns prominence.** Active timers, priority tasks, review queues, and budget alerts are prominent; configuration is quiet.
5. **Fluid dynamics & physics.** Every interaction acknowledges success, progress, or state change with custom cubic-bezier spring curves.
6. **Accessible by construction.** Keyboard navigation, visible focus rings, WCAG AA contrast (≥ 4.5:1), and 44×44px touch targets are baseline requirements.
7. **Zero-Slop Discipline.** No generic 3-card equal rows, no random AI-purple mesh blobs, no unstyled Inter defaults, no 0ms instant state transitions.

### 1.3 Non-Goals

- Mimicking Apple/macOS or Windows literally without purpose.
- Heavy glass blur on scrolling containers that causes GPU frame drops.
- Generic symmetrical 3-column Bootstrap card grids.
- Duplicating server state in Zustand.
- Permitting `any`, unhandled promises, or truncated placeholder code (`// TODO`, `// ...`).

---

## 2. Design Language: Warm Vibrancy & High-End Visual Architecture

The visual system combines warm dark translucency with vibrant module accents, concentric squircle geometry, and fluid kinetic tension.

### 2.1 The "Absolute Zero" Directive (Anti-Slop Clichés)

Any frontend implementation containing these patterns is rejected:
- ❌ **Banned Fonts:** Unconfigured browser default fonts or unstyled Inter everywhere. Use `Geist` or `Plus Jakarta Sans` for body/headings, `Geist Mono` for code/data, and `Noto Sans JP` for Japanese text.
- ❌ **Banned Borders & Shadows:** Generic 1px solid gray borders (`border-gray-200`, `border-slate-800`) and harsh black drop shadows (`shadow-md`, `rgba(0,0,0,0.3)`). Shadows must be tinted to the background hue (e.g., `shadow-emerald-950/40`).
- ❌ **Banned Layouts:** Rigid symmetrical 3-column card rows. Use Asymmetric Bento Grids (`col-span-8` + `col-span-4`), split master-detail views, or staggered masonry.
- ❌ **Banned Motion:** Standard `linear` or `ease-in-out` transitions. Instant 0ms state changes without spring interpolation.
- ❌ **Banned Color Schemes:** Pure `#000000` pitch black backgrounds or oversaturated neon AI gradients.

### 2.2 Visual Rules & Haptic Architecture

| Element | Rule & Implementation |
|---|---|
| **Canvas** | Deep OLED charcoal (`#090b10` or `slate-950`) with subtle radial gradients of primary module glow |
| **Surfaces (Double-Bezel)** | **Outer Shell:** `bg-white/5 border border-white/10 p-1.5 rounded-[1.5rem]`<br>**Inner Core:** `bg-slate-900/60 shadow-[inset_0_1px_1px_rgba(255,255,255,0.1)] rounded-[calc(1.5rem-0.375rem)]` |
| **Primary Buttons (Button-in-Button)** | Pill shape `rounded-full px-5 py-2.5 active:scale-[0.98]` with trailing icon encased in inner circular wrapper `w-7 h-7 rounded-full bg-white/10 group-hover:translate-x-1` |
| **Borders** | Subtle tinted hairlines: `border border-white/10` (dark) · module-tinted hover borders |
| **Shadows** | Deep, layered, tinted shadows: `shadow-2xl shadow-black/50` + `shadow-[0_8px_30px_rgb(0,0,0,0.12)]` |
| **Typography** | `Plus Jakarta Sans` / `Geist` (body/headings) + `Geist Mono` (code/numbers) + `Noto Sans JP` (Japanese) |
| **Typography Nuances** | Numeric displays: `tabular-nums` (`font-variant-numeric: tabular-nums`). Headings: `text-wrap: balance`. |
| **Motion Curves** | Custom cubic-bezier physics: `cubic-bezier(0.32, 0.72, 0, 1)` or spring `{ stiffness: 320, damping: 28 }` |
| **Icons** | Lucide React (primary, standardized `strokeWidth={1.75}`) + accessible labels (`aria-label`) |

### 2.3 Layer Hierarchy

```text
Dark canvas with ambient radial module glow (#090b10)
  └─ Floating left navigation rail (glass blur, active module pill indicator)
      └─ Translucent top bar with breadcrumb, timer widget, locale selector
          └─ Content workspace (Double-Bezel Asymmetric Bento Grid)
              └─ Popovers, command palette (higher blur, backdrop-blur-2xl)
                  └─ Modal Dialogs (dimmed overlay + Double-Bezel card)
                      └─ Toasts (fixed corner, spring bounce in, auto-dismiss)
```

### 2.4 Three-Layer Design Tokens (`fe/src/styles/tokens.css`)

Structured according to the `design-system` token specification:

```css
/* fe/src/styles/tokens.css */
:root {
  /* 1. Primitive Tokens (Raw Values) */
  --primitive-slate-950: 222 47% 5%;       /* #090b10 */
  --primitive-slate-900: 222 30% 9%;
  --primitive-slate-800: 217 33% 17%;
  --primitive-emerald-500: 158 64% 52%;
  --primitive-violet-500: 267 83% 69%;
  --primitive-amber-500: 38 92% 55%;
  --primitive-blue-500: 217 91% 60%;
  --primitive-rose-500: 0 84% 60%;

  /* 2. Semantic Tokens (Purpose Aliases) */
  --bg-base: var(--primitive-slate-950);
  --bg-surface: var(--primitive-slate-900);
  --bg-surface-elevated: var(--primitive-slate-800);
  --foreground: 210 40% 96%;
  --muted-foreground: 215 16% 60%;
  --border: rgba(255, 255, 255, 0.08);
  --border-focus: rgba(255, 255, 255, 0.25);

  /* Module Accents */
  --accent-dashboard: 230 85% 65%;
  --accent-tasks: var(--primitive-blue-500);
  --accent-pomodoro: var(--primitive-rose-500);
  --accent-linguistics: var(--primitive-violet-500);
  --accent-finance: var(--primitive-emerald-500);
  --accent-journal: var(--primitive-amber-500);

  /* Status Colors */
  --success: 142 71% 45%;
  --warning: 38 92% 50%;
  --destructive: var(--primitive-rose-500);

  /* 3. Component Geometry & Physics Tokens */
  --radius-outer: 1.5rem;                  /* 24px outer bezel */
  --radius-inner: calc(1.5rem - 0.375rem); /* 18px inner core */
  --radius-control: 0.75rem;               /* 12px buttons & inputs */
  --radius-pill: 9999px;

  --ease-spring: cubic-bezier(0.32, 0.72, 0, 1);
  --ease-out-expo: cubic-bezier(0.16, 1, 0.3, 1);
}
```

### 2.5 Module Color Identities & Ambient Themes

| Module | Accent Color | Vibe & Visual Persona | Usage & Highlights |
|---|---|---|---|
| **Today / Dashboard** | Blue-Violet (`#6366f1` / `#8b5cf6`) | Asymmetric Bento, Ambient glow | Greeting hero banner, next action priority, KPI summary |
| **Deep Work / Tasks** | Focused Blue (`#3b82f6`) | Kanban flow, clear state chips | Priority badges, drag-and-drop ghost cards |
| **Deep Work / Pomodoro** | Crimson Flame (`#f43f5e`) | Zen Minimalist, high contrast | Circular progress ring, break transition, session dots |
| **Linguistics SRS** | Deep Violet (`#a855f7`) | Interactive, tactile mastery | 3D flashcard flip, furigana accents, Leitner box badges |
| **Personal Finance** | Emerald Wealth (`#10b981`) | Analytical, high-density bento | Tabular amounts, cashflow area chart, budget meters |
| **Markdown Journal** | Warm Amber (`#f59e0b`) | Editorial Luxury, calm paper feel | Markdown toolbar, mood badges, 30-day streak tracker |
| **Settings / IAM** | Slate (`#64748b`) | Quiet, distraction-free | Clean forms, toggle switches, security keys |

---

## 3. Design Intelligence Tooling (`ui-ux-pro-max` Integration)

Agents and developers have access to the local CLI search tool in `.agents/skills/ui-ux-pro-max/scripts/search.py`. Use it before designing or refactoring components to pull authoritative, domain-calibrated guidance:

### 3.1 CLI Execution Recipes

```powershell
# 1. Look up accessibility and form validation UX rules
python .agents/skills/ui-ux-pro-max/scripts/search.py "error summary inline validation" --domain ux

# 2. Look up chart guidelines for Finance or Learning statistics
python .agents/skills/ui-ux-pro-max/scripts/search.py "cashflow budget analytics" --domain chart

# 3. Look up color palettes calibrated for Dark OLED interfaces
python .agents/skills/ui-ux-pro-max/scripts/search.py "personal productivity dark oled" --domain color

# 4. Look up stack-specific guidance for React + Tailwind + shadcn
python .agents/skills/ui-ux-pro-max/scripts/search.py "modal dialog focus trap" --stack shadcn
python .agents/skills/ui-ux-pro-max/scripts/search.py "virtual list performance" --stack react
```

---

## 4. Information Architecture & Routing

### 4.1 Route Map (React Router v7)

```text
/
├── /login                                # Public auth (centered glass card)
├── /register                             # Public registration + strength indicator
├── /app                                  # Protected shell (AppLayout)
│   ├── /today                            # Asymmetric Bento Dashboard
│   ├── /tasks                            # Kanban task board
│   ├── /pomodoro                         # Dedicated fullscreen Pomodoro view
│   ├── /learn                            # Linguistics hub overview
│   │   ├── /review                       # Distraction-free SRS study session (3D flip)
│   │   ├── /vocabulary                   # Full vocabulary management & TTS
│   │   └── /stats                        # 90-day heatmap & Leitner box distribution
│   ├── /finance                          # Finance overview & cashflow analytics
│   │   ├── /transactions                 # Grouped transaction list + quick-add
│   │   ├── /accounts                     # Multi-wallet account cards
│   │   ├── /budgets                      # Monthly budget progress meters
│   │   └── /goals                        # Savings goal rings
│   ├── /journal                          # Calendar heatmap + list view
│   │   ├── /new                          # Markdown editor + live preview
│   │   ├── /:journalId                   # Journal detail & backlinks view
│   │   └── /:journalId/edit              # Edit journal entry
│   └── /settings                         # Preferences, language, profile
└── *                                     # 404 NotFound page with safe return CTA
```

### 4.2 Page Anatomy

1. **Ambient Module Glow:** Subtle top radial gradient reflecting current module accent.
2. **PageHeader:** Double-bezel or translucent strip featuring title, breadcrumbs, contextual description, and nested CTA button (`Button-in-Button`).
3. **Filter & Metric Strip:** Quick KPI chips, date selectors, or search bars.
4. **Primary Workspace Surface:** Asymmetric Bento Grid, Kanban columns, or full-width data tables.
5. **Drawer / Modal Overlays:** Slide-in quick-add drawers or spring-loaded modal dialogs for secondary actions.

---

## 5. Directory Structure & Architecture Standards

```text
fe/src/
├── app/
│   ├── App.tsx                           # Root providers + router mounting
│   ├── providers.tsx                     # QueryClient, Theme, Toast, i18n
│   └── router.tsx                        # React Router v7 route definitions
├── assets/
│   └── icons/                            # Custom SVGs / branding assets
├── components/
│   ├── ui/                               # shadcn/ui primitives (Button, Dialog, Select, etc.)
│   ├── design-system/                    # High-End Design System Primitives
│   │   ├── DoubleBezelCard.tsx           # Machined nested container (outer + inner)
│   │   ├── GlassPanel.tsx                # Translucent glass surface primitive
│   │   ├── ButtonInButton.tsx            # Pill button with nested trailing icon
│   │   ├── PageHeader.tsx                # Module-tinted page hero header
│   │   ├── StatCard.tsx                  # KPI card with tabular-nums & sparkline/trend
│   │   ├── AmountDisplay.tsx             # Locale-aware colored currency (tabular-nums)
│   │   ├── ProgressRing.tsx              # SVG animated circular progress
│   │   ├── DataTable.tsx                 # Accessible, sortable data table
│   │   ├── EmptyState.tsx                # Actionable empty state with CTA
│   │   └── CommandPalette.tsx            # Global search (Ctrl+K)
│   └── feedback/
│       ├── LoadingSkeleton.tsx           # Exact layout-shaped skeleton loader
│       ├── ErrorBoundary.tsx             # Graceful React error boundary
│       └── Toast.tsx                     # Spring-animated accessible notifications
├── constants/
│   └── constants.ts                      # Non-domain constants (timeouts, page sizes, env)
├── domain/
│   └── constants.ts                      # Domain enums & states (as const, mirrors Go domain)
├── features/
│   ├── auth/                             # Login, register, session restoration
│   ├── today/                            # Personal Asymmetric Bento Dashboard
│   ├── tasks/                            # Kanban board & task dialogs
│   ├── pomodoro/                         # Pomodoro widget, timers, WebSocket client
│   ├── linguistics/                      # SRS review, flashcards (3D flip), TTS
│   ├── finance/                          # Wallets, transactions, cashflow charts, budgets
│   └── journal/                          # Markdown editor, calendar heatmap, search drawer
├── hooks/
│   ├── useDebounce.ts
│   ├── useKeyboardShortcut.ts
│   ├── useAudioTTS.ts                   # Audio playback with status & waveform
│   ├── useWebSocket.ts                  # Authenticated WS with backoff reconnection
│   └── useLocalDraft.ts                 # Autosave drafts to localStorage
├── layouts/
│   ├── AppLayout.tsx                    # Main app shell (sidebar + topbar + workspace)
│   ├── AuthLayout.tsx                   # Auth shell with dynamic ambient gradient
│   └── components/
│       ├── Sidebar.tsx                  # Floating glass navigation rail
│       ├── TopBar.tsx                   # Translucent top bar + quick controls
│       └── FloatingPomodoroWidget.tsx   # Persistent minimal floating timer
├── lib/
│   ├── axios.ts                         # Axios client + concurrent refresh queue
│   ├── query-client.ts                  # TanStack Query client & cache configuration
│   ├── cn.ts                            # clsx + tailwind-merge utility
│   └── formatters.ts                    # Currency, date, time formatters
├── locales/
│   ├── en/ vi/ ja/                      # Namespaced JSON translation files
├── stores/
│   ├── ui-store.ts                      # Theme, sidebarCollapsed, locale
│   └── auth-store.ts                    # In-memory access token & user profile
└── styles/
    ├── globals.css                      # Tailwind directives + base reset
    └── tokens.css                       # 3-layer CSS custom properties
```

---

## 6. High-End Design System Primitives & Components

All components must strictly adhere to the `ui-styling` and `high-end-visual-design` standards.

### 6.1 `DoubleBezelCard` (Precision Nested Container)

Never place raw content flatly on the background. Use the concentric double-bezel pattern:

```tsx
// components/design-system/DoubleBezelCard.tsx
import React from 'react';
import { cn } from '@/lib/cn';

interface DoubleBezelCardProps extends React.HTMLAttributes<HTMLDivElement> {
  outerClassName?: string;
  glowColor?: string; // e.g. "rgba(16, 185, 129, 0.15)"
}

export const DoubleBezelCard: React.FC<DoubleBezelCardProps> = ({
  children,
  className,
  outerClassName,
  glowColor,
  ...props
}) => {
  return (
    <div
      className={cn(
        "relative p-1.5 rounded-[1.5rem] bg-white/[0.04] border border-white/10 backdrop-blur-xl",
        "transition-all duration-500 ease-[cubic-bezier(0.32,0.72,0,1)]",
        "hover:border-white/20 hover:bg-white/[0.06]",
        outerClassName
      )}
      style={glowColor ? { boxShadow: `0 0 40px -10px ${glowColor}` } : undefined}
      {...props}
    >
      <div
        className={cn(
          "w-full h-full p-6 rounded-[calc(1.5rem-0.375rem)] bg-slate-950/70",
          "shadow-[inset_0_1px_1px_rgba(255,255,255,0.12)] border border-white/5",
          className
        )}
      >
        {children}
      </div>
    </div>
  );
};
```

### 6.2 `ButtonInButton` (Nested Trailing CTA)

```tsx
// components/design-system/ButtonInButton.tsx
import React from 'react';
import { motion, type HTMLMotionProps } from 'framer-motion';
import { cn } from '@/lib/cn';

interface ButtonInButtonProps extends HTMLMotionProps<'button'> {
  icon: React.ElementType;
  variant?: 'primary' | 'secondary' | 'emerald' | 'violet';
}

export const ButtonInButton: React.FC<ButtonInButtonProps> = ({
  children,
  icon: Icon,
  variant = 'primary',
  className,
  ...props
}) => {
  const variantStyles = {
    primary: "bg-blue-600 hover:bg-blue-500 text-white shadow-lg shadow-blue-950/50",
    secondary: "bg-white/10 hover:bg-white/15 text-white border border-white/10",
    emerald: "bg-emerald-600 hover:bg-emerald-500 text-white shadow-lg shadow-emerald-950/50",
    violet: "bg-violet-600 hover:bg-violet-500 text-white shadow-lg shadow-violet-950/50",
  };

  return (
    <motion.button
      whileHover={{ scale: 1.02 }}
      whileTap={{ scale: 0.98 }}
      transition={{ type: 'spring', stiffness: 400, damping: 25 }}
      className={cn(
        "group relative inline-flex items-center justify-between gap-3",
        "rounded-full pl-5 pr-1.5 py-1.5 font-medium text-sm",
        "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-white/40",
        variantStyles[variant],
        className
      )}
      {...props}
    >
      <span>{children}</span>
      <div className="size-7 rounded-full bg-white/15 flex items-center justify-center transition-transform duration-300 group-hover:translate-x-0.5 group-hover:-translate-y-0.5">
        <Icon className="size-3.5" />
      </div>
    </motion.button>
  );
};
```

### 6.3 `AmountDisplay` (Tabular Data-Dense Currency)

```tsx
// components/design-system/AmountDisplay.tsx
import React from 'react';
import { cn } from '@/lib/cn';
import { formatCurrency } from '@/lib/formatters';

interface AmountDisplayProps {
  amount: number;
  currency?: 'VND' | 'USD' | 'JPY';
  locale?: string;
  showSign?: boolean;
  colorize?: boolean;
  className?: string;
}

export const AmountDisplay: React.FC<AmountDisplayProps> = ({
  amount,
  currency = 'VND',
  locale = 'vi-VN',
  showSign = false,
  colorize = true,
  className,
}) => {
  const isPositive = amount > 0;
  const isNegative = amount < 0;

  const colorClass = colorize
    ? isPositive
      ? 'text-emerald-400'
      : isNegative
      ? 'text-rose-400'
      : 'text-muted-foreground'
    : 'text-foreground';

  return (
    <span className={cn("font-mono tabular-nums tracking-tight font-semibold", colorClass, className)}>
      {showSign && isPositive ? '+' : ''}
      {formatCurrency(amount, currency, locale)}
    </span>
  );
};
```

---

## 7. State Management & Data Transport Policy

### 7.1 Server State (TanStack Query) vs UI State (Zustand)

| State Category | Responsible Layer | Rule / Constraint |
|---|---|---|
| **Remote Entities** (Tasks, Vocab, Transactions, Journals) | **TanStack Query** | Never store in Zustand. Use deterministic keys via feature key factories. |
| **Auth Token (Access Token)** | **Zustand (in-memory)** | In-memory only. Never save to `localStorage` (XSS vulnerability). |
| **Theme / Locale / Sidebar collapsed** | **Zustand (`ui-store.ts`)** | UI-only layout preferences, persisted to `localStorage`. |
| **Continuous Pointer/Scroll Physics** | **Framer Motion (`useMotionValue`)** | **NEVER** store mouse or scroll position in `useState` or Zustand (kills 60fps on mobile). |
| **Form Data & Validation** | **React Hook Form + Zod** | Controlled/uncontrolled forms with schema-driven validation. |

### 7.2 Deterministic Query Keys

```ts
// features/linguistics/api/linguistics-keys.ts
export const linguisticsKeys = {
  all: ['linguistics'] as const,
  vocabulary: (filters?: Record<string, unknown>) => [...linguisticsKeys.all, 'vocabulary', filters] as const,
  vocabDetail: (id: string) => [...linguisticsKeys.all, 'vocabulary', id] as const,
  studySessionQueue: (language: string, mode: string) => [...linguisticsKeys.all, 'queue', language, mode] as const,
  stats: (range: string) => [...linguisticsKeys.all, 'stats', range] as const,
};
```

---

## 8. Motion Choreography & Kinetic Standards

### 8.1 Performance & Safety Guardrails
- **GPU-Safe Animations:** Animate **exclusively** via `transform` and `opacity`. Never animate `top`, `left`, `width`, or `height`.
- **Blur Boundaries:** Restrict `backdrop-blur` to fixed/sticky headers and dialog overlays. Never apply blur filters to scrolling content containers.
- **Client Leaf Isolation:** Isolate all Framer Motion components as leaf components to avoid unnecessary parent re-renders.
- **Accessibility:** Fully respect `prefers-reduced-motion` with instant opacity fades.

### 8.2 Standard Animation Presets

```ts
// styles/motion-presets.ts
export const motionPresets = {
  // Fluid page & card reveal
  fadeInUp: {
    initial: { opacity: 0, y: 16 },
    animate: { opacity: 1, y: 0 },
    exit: { opacity: 0, y: -12 },
    transition: { duration: 0.4, ease: [0.32, 0.72, 0, 1] },
  },
  // Double-Bezel card hover
  cardHover: {
    whileHover: { y: -2, transition: { duration: 0.2, ease: [0.32, 0.72, 0, 1] } },
  },
  // 3D Card flip for Flashcards
  flipCard: {
    transition: { type: 'spring', stiffness: 260, damping: 24 },
  },
  // Stagger container
  staggerContainer: {
    animate: { transition: { staggerChildren: 0.06 } },
  },
};
```

---

## 9. Feature Modules Specification

### 9.1 Today — Personal Dashboard (`/app/today`)
- **Layout:** Asymmetric Bento Grid (`col-span-8` focus task hero + `col-span-4` Pomodoro widget, quick SRS queue badge, monthly cashflow mini-stat, journal streak).
- **Vibe:** Ambient blue-violet radial glow, greeting banner with dynamic time-of-day greeting (Vietnamese / English).
- **Interactions:** One-click task start, quick Pomodoro trigger, instant review link.

### 9.2 Deep Work — Tasks & Pomodoro (`/app/tasks`, `/app/pomodoro`)
- **Kanban Flow:** Three responsive columns (`TODO`, `IN_PROGRESS`, `DONE`) with drag-and-drop and full keyboard accessibility (`M` key menu to move).
- **Pomodoro Ring:** Large circular progress ring (`ProgressRing`), pulsing reconnecting state, synchronized with backend WebSocket `/ws/pomodoro`.
- **Audio Feedback:** Optional subtle chime on session completion.

### 9.3 Linguistics — Interactive Study Hub (`/app/learn`)
- **Language Selector:** Tabs for 🇯🇵 日本語 and 🇬🇧 English.
- **Study Modes:**
  1. *Flashcard Mode:* 3D Card Flip (`rotateY: 180deg`) with spacebar flip and `1/2/3/4` rating shortcuts.
  2. *Multiple Choice:* 4 glass option buttons with immediate green glow / red shake feedback.
  3. *Type-in Mode:* Inline input with fuzzy matching tolerance.
  4. *Audio Quiz:* Automated TTS playback, waveform indicator, meaning selection.
- **Japanese Typography:** Render Japanese glyphs with `Noto Sans JP`, furigana support, and Romaji toggle for beginners.

### 9.4 Personal Finance (`/app/finance`)
- **High-Density Dashboard:** Net worth KPI card (large emerald font, `tabular-nums`), 30-day cashflow area chart, budget meters with threshold alert badges.
- **Transaction Quick-Add:** Slide-in drawer with numeric keypad-friendly inputs and instant account balance recalculation.
- **Grouped Transaction List:** Grouped by day with colored status badges (+green for income, -rose for expense).

### 9.5 Journal — Knowledge Base (`/app/journal`)
- **Editorial Luxury Aesthetic:** Warm espresso/amber ambient tints, clean monospace typography option for code blocks.
- **Editor & Live Preview:** Markdown toolbar (bold, italic, headings, code, tables), autosave to `localStorage` every 10s, side-by-side or split preview.
- **Contribution Heatmap:** 365-day GitHub-style streak calendar with mood-tinted activity cells.

---

## 10. Quality, Accessibility & Pre-Delivery Audit Checklist

Before delivering any frontend feature, verify against this strict quality matrix:

- [ ] **Full Output:** Zero `// ...`, zero `// TODO`, zero truncated boilerplate code.
- [ ] **No Banned Clichés:** No 3 equal cards, no AI-purple blobs, no unconfigured Inter defaults, no 0ms transitions.
- [ ] **Haptic Containers:** Major cards implement the Double-Bezel nested architecture (`DoubleBezelCard`).
- [ ] **Button-in-Button:** Primary CTAs use the pill format with encased circular trailing icons.
- [ ] **4 Async States:** Skeleton loader, empty state with CTA, error state with retry button, and success state.
- [ ] **Tabular Numerics:** All currency, dates, timers, and metrics use `tabular-nums`.
- [ ] **WCAG AA Compliance:** Contrast ratio ≥ 4.5:1 on dark glass surfaces; visible `:focus-visible` rings.
- [ ] **Touch Targets:** Minimum 44×44px interactive area with ≥8px spacing.
- [ ] **Mobile Safety:** Mobile layouts collapse gracefully into a single vertical column (`min-h-[100dvh]` to avoid Safari viewport jump).
- [ ] **Internationalization:** All strings translated into `vi`, `en`, and `ja` in their respective namespace files.
- [ ] **Lint & Typecheck:** `npm run lint` and `npm run typecheck` pass with zero errors and zero warnings.

---

## 11. Implementation Roadmap & Progress Tracking

### PHASE-0: Foundation & Design System
- [ ] **FE-001** Vite + React 18 + TypeScript strict setup
- [ ] **FE-002** Tailwind CSS + 3-layer tokens (`tokens.css`) + shadcn/ui integration
- [ ] **FE-003** i18n setup (`vi`, `en`, `ja` namespaced files)
- [ ] **FE-004** Axios client + in-memory auth-store + concurrent refresh queue
- [ ] **FE-005** TanStack Query client + global error handler
- [ ] **FE-006** Zustand stores (`ui-store`, `auth-store`)
- [ ] **FE-007** AppLayout (Double-Bezel floating sidebar + glass topbar)
- [ ] **FE-008** Mobile layout (bottom navigation bar + sheet drawer)
- [ ] **FE-009** Design system primitives (`DoubleBezelCard`, `ButtonInButton`, `AmountDisplay`, `ProgressRing`, `StatCard`, `EmptyState`)
- [ ] **FE-010** Accessible toast notification system

### PHASE-1: IAM & Authentication
- [ ] **FE-101** Auth layout with dynamic radial gradient mesh
- [ ] **FE-102** Login page + form + interceptor hook-up
- [ ] **FE-103** Register page + password strength indicator
- [ ] **FE-104** Protected route guards
- [ ] **FE-105** Silent session restoration on boot

### PHASE-2: Deep Work & Pomodoro
- [ ] **FE-201** Today Asymmetric Bento Dashboard
- [ ] **FE-202** Tasks Kanban board (drag-and-drop + keyboard menu)
- [ ] **FE-203** Task create/edit drawer
- [ ] **FE-204** Pomodoro fullscreen view + global floating widget
- [ ] **FE-205** WebSocket client hook (`useWebSocket` with backoff reconnection)
- [ ] **FE-206** Pomodoro session history & statistics

### PHASE-3: Linguistics SRS
- [ ] **FE-301** Learning hub overview + language tabs (🇯🇵 / 🇬🇧)
- [ ] **FE-302** Vocabulary management table + TTS audio playback
- [ ] **FE-303** Flashcard study mode (3D flip + keyboard rating)
- [ ] **FE-304** Multiple choice study mode (kinetic feedback)
- [ ] **FE-305** Type-in practice mode (fuzzy matching)
- [ ] **FE-306** Audio quiz mode (TTS listener)
- [ ] **FE-307** Session summary screen with XP & mastery badges
- [ ] **FE-308** 90-day learning heatmap & Leitner box distribution
- [ ] **FE-309** Bulk vocabulary CSV import UI

### PHASE-4: Personal Finance
- [ ] **FE-401** Finance dashboard (net worth KPI, 30-day cashflow area chart)
- [ ] **FE-402** Transaction list grouped by day with filter controls
- [ ] **FE-403** Transaction quick-add slide-in drawer
- [ ] **FE-404** Budget tracker with threshold alert badges
- [ ] **FE-405** Savings goals circular progress rings
- [ ] **FE-406** Multi-wallet account management
- [ ] **FE-407** Category & tag manager

### PHASE-5: Knowledge Journal
- [ ] **FE-501** Journal list with 365-day calendar heatmap toggle
- [ ] **FE-502** Journal list view + slide-in search drawer
- [ ] **FE-503** Markdown editor (toolbar + live preview + 10s autosave)
- [ ] **FE-504** Journal read-only detail view with backlinks
- [ ] **FE-505** Mood & energy level picker component
- [ ] **FE-506** Journal linking & reference UI
- [ ] **FE-507** Writing streak widget
- [ ] **FE-508** Full-text & semantic search result display

### PHASE-6: Polish & Performance
- [ ] **FE-601** Command palette (`Ctrl+K` global spotlight)
- [ ] **FE-602** Global keyboard shortcuts registration
- [ ] **FE-603** Mobile responsive review (`min-h-[100dvh]`, no horizontal overflow)
- [ ] **FE-604** Playwright E2E tests for critical user journeys
- [ ] **FE-605** Performance audit (bundle analysis, image lazy-loading, LCP < 1.2s)
- [ ] **FE-606** WCAG AA accessibility audit

---

## 12. Required AI Execution Protocol

Build GoPA frontend one vertical slice at a time. **Do not generate the entire application in one unreviewable response.**

For every task, follow this 4-phase protocol:

### Phase A — Design & Architectural Blueprint
- State the deliverable and exact scope.
- Explain how the UI applies the **Warm Vibrancy & Double-Bezel architecture**.
- Cite any relevant guidance queried from `ui-ux-pro-max` search CLI.
- Define component props, TypeScript view models, TanStack Query keys, and routing boundaries.
- State accessibility considerations (WCAG contrast, touch targets, focus rings) and test strategy.

### Phase B — Execution Commands
- Provide complete, copyable commands for dependencies, tests, and builds.

### Phase C — Complete Code Delivery (`full-output-enforcement`)
- Full, untruncated files with exact project-relative paths.
- All types, translations (`vi`, `en`, `ja`), styles, and wiring included.
- **ZERO** `any`, zero `// ...`, zero `// TODO`.

### Phase D — Verification & Handoff
- Run `npm run typecheck`, `npm run lint`, and tests.
- Check off completed items in Section 11.
- Outline the next smallest slice and pause for user feedback.

---

*This document is the frontend single source of truth. Any visual, technical, or product decision must be documented here before implementation begins.*
