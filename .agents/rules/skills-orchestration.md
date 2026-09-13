---
description: Comprehensive rules for orchestrating and applying the 20 specialized design and frontend skills in .agents/skills/ for GoPA.
globs: ["**/*"]
alwaysApply: true
---

# GoPA Agent Skills Orchestration Rules

Every AI Agent executing in the GoPA repository must automatically conform to the following behavioral and engineering rules:

## 1. Skill Discovery & Execution Rules
- The repository contains 20 specialized skills located in `.agents/skills/`.
- **Search CLI Execution**: Whenever you make visual, UX, layout, or chart decisions, consult `.agents/skills/ui-ux-pro-max/scripts/search.py` using:
  ```powershell
  python .agents/skills/ui-ux-pro-max/scripts/search.py "<query>" --domain <domain> --stack react
  ```
- Available domains: `style`, `color`, `chart`, `landing`, `product`, `ux`, `typography`, `icons`, `gsap`, `react`.

## 2. Mandatory Full-Output Enforcement
- Apply the `full-output-enforcement` skill to every coding task.
- **NEVER** use `// ...`, `// rest of code`, `// implement here`, `// TODO`, or skeletal placeholders.
- Always output complete, self-contained files with all TypeScript types, error handling, and imports.

## 3. High-End Visual Standards (Anti-Slop Discipline)
- Adhere to `high-end-visual-design` and `design-taste-frontend`:
  - **Double-Bezel Container Architecture**: Cards and panels must use nested shells (`rounded-[1.5rem] bg-white/5 border border-white/10 p-1.5` wrapping an inner core with concentric radius `rounded-[calc(1.5rem-0.375rem)]`).
  - **Nested CTA & Button-in-Button**: Trailing icons inside primary action buttons must be encased in their own circular container (`w-7 h-7 rounded-full bg-white/10`).
  - **Fluid Motion**: Transitions must use cubic-bezier curves (e.g., `cubic-bezier(0.32, 0.72, 0, 1)`) and GPU-safe properties (`transform`, `opacity`). Never animate layout dimensions (`width`, `height`, `top`, `left`).
  - **No Cliché Defaults**: Ban generic 3-card equal grids, ban pure `#000000` backgrounds, ban raw default Inter typography, and ban generic untinted black shadows.

## 4. Design System & Component Architecture
- Use `ui-styling` and `design-system`:
  - Rely on shadcn/ui components powered by Radix UI primitives and Tailwind CSS.
  - Structure tokens into 3 layers: Primitive → Semantic → Component.
  - Implement all 4 mandatory UI states for every data-fetching view: Loading Skeleton, Empty State with CTA, Error State with Retry, and Success State.

## 5. Accessibility & Mobile Standards
- Maintain WCAG AA compliance (contrast ≥ 4.5:1, keyboard focus-visible rings, screen-reader labels).
- Minimum interactive touch target is 44×44px with ≥8px spacing.
- Mobile layouts must use `min-h-[100dvh]` (not `100vh`) to prevent iOS Safari viewport jump, collapsing multi-column bento grids into single-column vertical flows below `768px`.
