# Arthur: Frontend

The dashboard and landing page for [Arthur](../README.md), a Vercel-style static-site deployment platform. It is a single-page app that talks to the Go API under `/v1`.

## Features

- Landing page (lazy-loaded animated background, respects `prefers-reduced-motion`)
- Email/password login and registration, and "Continue with GitHub" (full-page OAuth redirect)
- Project dashboard, plus a "new project" flow that triggers the first deployment
- Project detail with a deployments table: redeploy, promote (rollback) with confirmation, and a **Live** badge on the active deployment
- Deployment detail with a log viewer that polls for new lines and renders ANSI colors
- Environment variable management (list, set, delete) with masked values and a warning that static-site variables end up in the public bundle
- Delete project, copy buttons, toasts, skeleton loaders, 404 page

## Stack

React 19, TypeScript, Vite, Tailwind CSS 4, TanStack Query, React Router, React Hook Form + Zod, sonner (toasts), motion, lucide-react, anser (ANSI to HTML), oxlint.

## Getting started

Prerequisites: Node 24 and the Arthur API running on `http://localhost:8080` (see the [root README](../README.md); `make up` starts everything, including this app).

```bash
cd frontend
npm install
npm run dev        # http://localhost:5173
```

### Scripts

| Script | Does |
| --- | --- |
| `npm run dev` | Vite dev server with hot reload |
| `npm run build` | Type-check (`tsc -b`) and build to `dist/` |
| `npm run lint` | Lint with oxlint |
| `npm run preview` | Serve the production build locally |

## Configuration

| Variable | Default | Purpose |
| --- | --- | --- |
| `VITE_API_BASE` | empty | API origin. Leave empty in development so requests go through the Vite proxy. Set it at build time for production. |
| `VITE_SITE_SCHEME` | `http` | Scheme used to build deployed-site URLs |
| `VITE_SITE_BASE_DOMAIN` | `localhost:8080` | Base domain used to build deployed-site URLs (`<slug>.<domain>`) |
| `VITE_PROXY_TARGET` | `http://localhost:8080` | Dev proxy target for `/v1` (the Docker stack points it at the `api` service) |

`VITE_*` variables are compiled into the bundle, so never put secrets in them. These should match the backend's `SITE_SCHEME` and `SITE_BASE_DOMAIN`.

## How it talks to the backend

- **Same-origin in development.** Vite proxies `/v1` to the API, so the HttpOnly session cookies work without CORS tricks.
- **Cookies only.** Every request uses `credentials: "include"`; the app never touches tokens. In production the frontend and API must share a registrable domain (cookies are `SameSite=Lax`), and the backend's `FRONTEND_URL` must match the dashboard origin.
- **GitHub login** is a browser redirect to `/v1/auth/github/login`, not a fetch.
- **Registration** does not set cookies, so the UI logs in right after registering.
- **Log streaming** polls `GET /v1/deployments/{id}/logs?after=<seq>`, sends back `next_after` each time, and stops when `done` is `true`.

## Project structure

```
src/
  api/          client.ts (fetch wrapper), hooks.ts (TanStack Query hooks), types.ts (API types)
  components/
    ui/         shared UI: layout, sidebar, dialogs, buttons, inputs, status badge, route guards
    vendor/     adapted third-party components (see SOURCES.md)
  features/     auth/ projects/ (form schemas), deployments/ (LogViewer), env/ (EnvVarsPanel)
  lib/          ansi.ts, siteUrl.ts, utils.ts
  pages/        Landing, Login, Register, Dashboard, NewProject, ProjectDetail, DeploymentDetail, NotFound
  App.tsx       routes
```

The `@` alias points at `src/` (`import { cn } from "@/lib/utils"`).

### Routes

| Path | Access | Page |
| --- | --- | --- |
| `/` | public | Landing |
| `/login`, `/register` | public (redirects to dashboard when signed in) | Auth |
| `/dashboard` | signed in | Projects |
| `/projects/new` | signed in | Create a project |
| `/projects/:id` | signed in | Project detail: deployments, env vars, settings |
| `/projects/:projectId/deployments/:deploymentId` | signed in | Deployment logs |

## Production build

```bash
VITE_API_BASE=https://api.example.com \
VITE_SITE_SCHEME=https \
VITE_SITE_BASE_DOMAIN=sites.example.com \
npm run build
```

The output in `dist/` is static and can be served by any web server; configure it to fall back to `index.html` for client-side routes.

## Notes

- `PLAN.md` and `implementation.md` record the original build plan and status.
- Third-party component attributions are listed in `src/components/vendor/SOURCES.md`.
