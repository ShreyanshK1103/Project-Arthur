# Arthur

Arthur is a Vercel-style deployment platform for static sites. Connect a GitHub repository and Arthur clones it, builds it inside a locked-down Docker container, stores the output in S3, and serves it on a stable subdomain. Every push to the tracked branch can trigger a new deployment, and any earlier successful deployment can be promoted back to live (rollback) without rebuilding.

This is a monorepo:

| Directory | What it is |
| --- | --- |
| [`backend/`](backend/README.md) | Go API, build worker, Postgres schema and migrations |
| [`frontend/`](frontend/README.md) | React dashboard and landing page |
| `scripts/` | End-to-end smoke test |
| `docker-compose.yml`, `Makefile` | One-command local stack |

## Features

- **GitHub login** (OAuth) plus email/password, with HttpOnly cookie sessions and refresh-token rotation
- **Deployments from GitHub**: manual deploys, redeploys, and push webhooks (HMAC-verified)
- **Sandboxed builds**: install and build run in a `node:24` container with memory, CPU and PID limits, all capabilities dropped, a non-root user, and timeouts
- **Encrypted environment variables** per project (AES-256-GCM at rest), redacted from build logs
- **Stable project URLs and rollback**: `<slug>.<domain>` always serves the active deployment, and every deployment keeps an immutable permalink
- **Live build logs** in the dashboard via cursor-based polling, with ANSI colors rendered
- **Stuck-job recovery**: builds orphaned by a worker crash are marked failed

## How it works

```
 Browser (dashboard)                GitHub
        |                              |  push webhook
        v                              v
   +---------+   HTTP /v1        +-----------+
   | Frontend| ----------------> |  Go API   | --- serves <slug>.<domain> from S3
   +---------+                   +-----+-----+
                                       |
                                  queue in Postgres
                                       |
                                 +-----v-----+   docker run node:24   +-----------+
                                 |  Worker   | ---------------------> |  Sandbox  |
                                 +-----+-----+   (limited build)      +-----------+
                                       |
                                  upload output
                                       v
                                 S3 (private bucket)
```

1. The API stores a deployment as `queued` in Postgres.
2. The worker claims jobs with `FOR UPDATE SKIP LOCKED`, shallow-clones the branch, and builds in a container.
3. The output directory is uploaded to S3 under `<deploymentID>/`, and the deployment becomes the project's active one.
4. A request to `<slug>.<site domain>` resolves slug, then active deployment, then streams the file from S3.

Details, API reference and security notes are in the [backend README](backend/README.md).

## Tech stack

- **Backend:** Go, chi, PostgreSQL, sqlc, goose, Docker, AWS S3
- **Frontend:** React 19, TypeScript, Vite, Tailwind CSS 4, TanStack Query, React Router

## Quick start (Docker)

Prerequisites: Docker with Compose, `make`, a private S3 bucket with credentials, and a GitHub OAuth App (callback `http://localhost:8080/v1/auth/github/callback`).

```bash
cp backend/.env.example backend/.env     # fill in the values (see comments inside)
make up
```

| Service | URL |
| --- | --- |
| Dashboard | http://localhost:5173 |
| API health | http://localhost:8080/v1/healthz |
| Postgres | localhost:5433 (user, password and db: `arthur`) |

Deployed sites are served at `http://<slug>.localhost:8080` (browsers resolve `*.localhost` automatically).

> **Docker Desktop (macOS/Windows):** the containerised worker needs the host's Docker socket and a shared temp path, which only works with native Docker Engine. Use `make up-lite` (everything except the worker) and `make worker-local` (runs the worker on your machine) instead.

### Make targets

Run `make` to see them all.

| Target | Does |
| --- | --- |
| `make up` | Build and start the full stack |
| `make up-lite` | Start everything except the worker |
| `make worker-local` | Run the worker on the host (pair with `up-lite`) |
| `make down` | Stop the stack |
| `make restart` / `make rebuild` | Restart / rebuild images without cache |
| `make logs s=<service>` | Follow logs (`db`, `api`, `worker`, `frontend`) |
| `make ps` | Show running services |
| `make migrate` | Run database migrations |
| `make psql` | Open a psql shell |
| `make test` | Run the backend Go tests |
| `make sqlc` | Regenerate sqlc code |
| `make smoke REPO=<url>` | End-to-end test against a public npm repo |
| `make clean` | Remove containers and volumes (asks first) |

## Running without Docker

See [backend/README.md](backend/README.md#getting-started) and [frontend/README.md](frontend/README.md#getting-started). In short: Postgres + migrations, `go run ./cmd/api`, `go run ./cmd/worker`, then `npm run dev` in `frontend/`.

## Testing

```bash
make test                                              # Go unit tests
make smoke REPO=https://github.com/<you>/<vite-repo>   # full flow: auth, deploy, logs, rollback, authorization, delete
```

The smoke test needs the API and worker running, plus `curl` and `jq`.

## Known limitations / roadmap

- Static sites only; no long-running server apps
- No network isolation during builds yet, so restrict signup to people you trust
- Single worker host
- No custom domains or pull-request preview deployments
- S3 objects are not yet cleaned up when a project is deleted

## License

Add a license of your choice (for example MIT) before publishing.
