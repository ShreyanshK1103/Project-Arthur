# Arthur Frontend — Implementation Status

All planned phases are implemented. Phase 1 was already present; Phases 2–6 plus the backend project list/get/delete gap are done.

## Completed

| Step | Status |
|------|--------|
| Phase 1 Foundation | Done (pre-existing) |
| Phase 2 Auth | Done — login, register (+ auto-login), GitHub OAuth redirect, logout toasts |
| Backend GET/DELETE projects | Done — `GET /v1/projects`, `GET /v1/projects/{id}`, `DELETE /v1/projects/{id}` |
| Phase 3 Projects | Done — dashboard grid, new project + first deployment, project detail shell |
| Phase 4 Deployments/logs | Done — table, redeploy, promote confirm, ANSI log viewer with polling |
| Phase 5 Env vars | Done — list/set/delete, mask, backend key rules, public-bundle warning |
| Phase 6 Landing/polish | Done — landing with lazy PixelGrid, reduced-motion, focus styles, 404 |

## Manual test checklist

1. `npm run dev` (frontend) + backend on `:8080`
2. Register → lands on dashboard (register then login for cookies)
3. Login with email/password; GitHub OAuth via `/v1/auth/github/login`
4. Create project → redirected to deployment logs
5. Promote a successful non-active deployment; delete project from Settings
6. Env vars: add/delete; confirm warning copy
7. Landing `/` works logged out and logged in

## Notes

- Register API does not set cookies; UI logs in immediately after register.
- ANSI logs use `anser` (see `src/lib/ansi.ts`).
- Componentry Pixel Canvas pattern adapted locally as `PixelGrid` (see `SOURCES.md`).
