# Arthur Frontend: Build Plan

Arthur is a Vercel-style deployment platform for static sites. The backend (Go API + build worker) already exists. This plan covers the **frontend dashboard and landing page**, to be built by an AI coding agent phase by phase.

Backend repo: https://github.com/ShreyanshK1103/Project-Arthur/tree/main/backend

---

## 0. Read this first (gaps to settle)

1. **Missing endpoints.** The backend README lists no `GET /projects`, `GET /projects/{id}`, or `DELETE /projects/{id}`. A dashboard needs all three. Either add them to the Go backend first, or the agent must flag them and stop. Do not fake them.
2. **Unknown shapes.** The README gives no request/response bodies and no deployment status values. The agent must read these backend paths and derive TypeScript types from them:
   - `backend/internal/handlers`
   - `backend/internal/models`
   - `backend/sql/schema`
   - `backend/sql/queries`
3. **Never invent** endpoints, fields, or status values. If something is missing, list it under "Blocked by backend" in the phase summary.

---

## 1. Stack

| Concern | Choice | Why |
| --- | --- | --- |
| Build tool | Vite + React + TypeScript | Backend is a pure API, no SSR needed |
| Styling | Tailwind CSS + shadcn/ui primitives | Componentry and React Bits are built for this workflow |
| Routing | React Router | Simple SPA routing |
| Data | TanStack Query | Polling, caching, and mutations suit logs and deployments |
| Forms | React Hook Form + Zod | Validates repo URLs and env var keys |
| Animation | `motion` | Used by most React Bits components |
| Icons | lucide-react | Consistent icon set |
| Toasts | sonner | Lightweight notifications |
| ANSI logs | `anser` or `ansi-to-react` (state why when adding) | Build logs may contain color codes |

### Auth and CORS

- The backend uses **HttpOnly cookie sessions**, and the refresh cookie is scoped to `/v1/auth`.
- **Dev:** use a Vite proxy (`/v1` -> `http://localhost:8080`) so everything is same-origin.
- **Production:** the frontend and API must share a registrable domain (cookies are `SameSite=Lax`). Set `FRONTEND_URL` on the backend to match.
- Every fetch uses `credentials: "include"`.
- The `/v1/auth/github/login` flow is a full-page redirect, not a fetch.

---

## 2. Component sources

The agent must **read each site's docs first** and use the install method currently documented there. Do not guess commands or props.

| Source | URL | Use for |
| --- | --- | --- |
| Componentry | https://componentry.dev/docs | Core UI kit: buttons, cards, tabs, dialogs, tables, badges, inputs. Main source for the dashboard. |
| React Bits | https://reactbits.dev/get-started/index | Animated, expressive pieces for landing and auth pages: animated backgrounds, text effects, spotlight cards, counters. Pick the **Tailwind + TypeScript** variants. |
| Uiverse | https://uiverse.io/ | Small flourishes only: a loader, a toggle, a special button. These are HTML/CSS, so port to React + Tailwind, scope the CSS, and credit the source in a comment. |

### Rules for mixing them

1. **Design tokens first.** Define colors, radius, fonts, and spacing as CSS variables and in the Tailwind config before pulling in any component. Restyle every pulled-in component to use the tokens so the app does not look like a patchwork.
2. **Dark theme by default**, with a light toggle if time permits.
3. **Heavy effects** (WebGL or canvas backgrounds) appear only on landing and auth pages, are lazy-loaded, and are disabled under `prefers-reduced-motion`.
4. **The dashboard stays calm.** Motion is limited to page transitions, hover states, and status indicators.
5. **Folder separation:** copied code goes in `src/components/vendor/`, app-level wrappers in `src/components/ui/`, so parts are easy to find and swap.
6. Keep a `src/components/vendor/SOURCES.md` listing each copied component, where it came from, and the date.

---

## 3. Folder structure

```
src/
  api/            client.ts (fetch wrapper, 401 -> refresh -> retry), types.ts, hooks per resource
  components/
    ui/           app-level components built on tokens
    vendor/       copied from Componentry / React Bits / Uiverse (+ SOURCES.md)
  features/       auth, projects, deployments, env, logs
  pages/          Landing, Login, Register, Dashboard, NewProject,
                  ProjectDetail, DeploymentDetail, NotFound
  lib/            utils, siteUrl helpers, formatters
```

---

## 4. Phases

Each phase ends with something runnable. **The agent stops and summarizes after each phase, then waits for approval.**

### Phase 1: Foundation

- Scaffold Vite + React + TS, Tailwind, and shadcn setup.
- Define design tokens, fonts, and the dark theme.
- Build the API client:
  - Typed error shape
  - `credentials: "include"` on every request
  - Automatic `POST /auth/refresh` and one retry on 401, with a single in-flight refresh shared by concurrent requests
- Build the app shell: sidebar or topbar, user menu, and route guards using `GET /auth/me`.
- Add `.env.example` with `VITE_API_BASE` (empty in dev because of the proxy).

**Done when:** the app loads, and an unauthenticated user is redirected to login.

### Phase 2: Auth

- Login and Register pages (email/password).
- "Continue with GitHub" button that navigates to `/v1/auth/github/login`.
- Handle the post-OAuth redirect back to the dashboard.
- Logout.
- Inline form errors plus toasts.

**Done when:** both login methods work end to end against the local backend.

### Phase 3: Projects

- Dashboard with a project grid or list: name, repo, latest status badge, live URL (`slug` / `url`).
- New Project flow:
  - Repo URL (https only, no embedded credentials, validated the same way as the backend)
  - Branch
  - Install command, build command, output directory
  - Optional env var rows
  - "Deploy" action
- Empty, loading, and error states.

**Done when:** a project can be created and appears on the dashboard.

### Phase 4: Deployments and live logs

- Project detail page with tabs: **Deployments**, **Environment**, **Settings**.
- Deployments table: status badge, commit and time, duration, "Active" marker (`is_active`), permalink.
- Actions:
  - **Redeploy**
  - **Promote** (rollback or roll-forward): only on `success` deployments, with a confirmation dialog
- Deployment detail page with a **log viewer**:
  - Poll `GET /deployments/{id}/logs?after=<seq>`, starting at `0`
  - Send back `next_after` each time, stop when `done` is `true`
  - Poll every 1 to 2 seconds, back off after errors
  - Auto-scroll, with pause when the user scrolls up and a "Jump to latest" button
  - Monospace font, line numbers, copy button, ANSI color support
- Poll deployment status while it is `queued` or building. Stop once terminal.

**Done when:** a full flow works: push, watch logs stream, site goes live, roll back to an older deployment, and the live URL changes.

### Phase 5: Environment variables

- List, add or edit (`PUT`), delete (`DELETE .../env/{key}`).
- Mask values by default with a reveal toggle.
- Validate keys (match backend rules found in the code).
- Visible warning: *"Static-site build variables end up in the public bundle. Don't put secrets here."*

**Done when:** variables can be managed, and a redeploy picks them up.

### Phase 6: Landing page and polish

- Landing page:
  - Hero with a React Bits background and text effect
  - "How it works": Connect -> Build -> Deploy -> Rollback
  - Feature grid based on the backend features (sandboxed builds, encrypted env vars, instant rollback, stable URLs and permalinks, live logs)
  - Footer
- Polish:
  - Skeleton loaders
  - Keyboard accessibility and visible focus states
  - Responsive layouts down to 375px
  - Favicon and meta tags
  - 404 page
  - Contrast check and `prefers-reduced-motion` check

**Done when:** Lighthouse accessibility score is 90+ on landing, login, and dashboard, and the app is usable on mobile.

---

## 5. UI behavior guidelines

- **Status badges:** one shared mapping from status to color, icon, and label, reused everywhere.
  - Queued: neutral, pulsing
  - Building: animated
  - Success: green
  - Failed: red
  - Any status the backend returns that is not mapped falls back to a neutral badge showing the raw value.
- **URLs:** always show the stable `<slug>.<domain>` URL with copy and open buttons. Show the 12-character permalink as secondary.
- **Times:** relative ("2 min ago") with the exact timestamp in a tooltip.
- **Destructive actions** (delete project, delete variable, promote) always use a confirmation dialog.
- **Errors:** every failed request shows a human-readable message. Never swallow errors silently.
- **Loading:** use skeletons for lists and detail pages, not blank screens.

---

## 6. Rules for the agent

1. Read the backend code before writing any API type. Never invent endpoints or fields. Flag anything missing.
2. Read the docs of Componentry, React Bits, and Uiverse before using them. Do not guess install commands or props.
3. Do not hardcode URLs. Use the Vite proxy or `VITE_API_BASE`, and keep `.env.example` current.
4. Keep components small and typed. No `any`.
5. At the end of each phase run `npm run build`, the type-check, and the linter. Fix every error before moving on.
6. Do not add dependencies without saying why.
7. Never commit secrets or `.env` files.
8. After each phase, write a short summary covering:
   - What changed
   - What to test manually
   - Anything blocked by missing backend endpoints

---

## 7. Starter prompt for the agent

> You're building the frontend for Arthur, a static-site deployment platform. Read `PLAN.md` fully, then read the backend repo (`backend/internal/handlers`, `backend/internal/models`, `backend/sql/schema`) to learn the exact request and response shapes. Start with Phase 1 only. Before using Componentry, React Bits, or Uiverse, fetch their docs and follow their current install instructions. When you finish a phase, stop, summarize, list anything blocked by missing backend endpoints, and wait for my go-ahead.

---

## 8. Backend known limitations (to reflect in the UI copy)

- Static sites only. No long-running server apps.
- No custom domains and no pull request preview deployments yet.
- Single worker host, so builds may queue under load.
- No network isolation during builds yet.
