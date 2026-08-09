# GoPA web client

The React application is the Phase 0 shell: feature-oriented routes, a Soft Acrylic Fluent layout, strict TypeScript, TanStack Query, Zustand UI preferences, Axios transport, and Vietnamese/English/Japanese resources.

## Start

```powershell
Copy-Item .env.example .env
npm install
npm run dev
```

Open `http://localhost:5173`. IAM is deliberately deferred to Phase 1, so the application shell is directly reachable during bootstrap.

## Quality checks

```powershell
npm run lint
npm run typecheck
npm run test
npm run build
```
