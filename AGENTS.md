# AGENTS.md — arci-api

## Project

Go 1.25 module `arci.it`. Gin HTTP API backed by MariaDB. JWT auth.

## Developer commands

- **Run locally**: `go run ./cmd/arci/main.go` (requires MariaDB running and env vars set)
- **Start full stack**: `docker compose up` — builds and runs API on `:8080` with MariaDB
- **Build**: `go build -o main ./cmd/arci/main.go`
- **No test suite exists** — do not invent `go test` commands

## Required environment variables

| Var | Required | Default |
|-----|----------|---------|
| `DB_HOST` | yes | — |
| `DB_USER` | yes | — |
| `DB_NAME` | yes | — |
| `JWT_SECRET` | yes | — |
| `DB_PORT` | no | `3306` |
| `DB_PASSWORD` | no | `""` |

When running via `docker compose up`, these are loaded from `.env`.

## Architecture

```
cmd/arci/main.go          — single entrypoint, wires DB + routes, listens :8080
pkg/arci/db/              — raw database/sql + MySQL driver (connect.go, auth.go, events.go, partecipation.go)
pkg/arci/routes/          — Gin handlers + JWT middleware (auth.go, events.go, partecipation.go, middlewares.go)
schema.sql                — DB init, mounted into MariaDB container on first start
```

- `pkg/arci/db` uses a **package-level global `db *sql.DB`** set by `ConnectDatabase()`. All DB functions reference it directly.
- Auth middleware reads `Authorization` header, strips `Bearer ` prefix, validates JWT, sets `member_id`, `showname`, `is_admin` in Gin context.

## API routes

- **Public**: `POST /register`, `POST /login`
- **Protected** (JWT required): `GET/POST /events`, `POST/DELETE /events/:event_id/partecipate`, `GET /roles`
- **Admin** (JWT + `is_admin=true`): `POST /roles`, `DELETE /roles/:id`

## Gotchas

- `Partecipation` is the table name (Italian spelling) — do not "correct" to `Participation`
- DB columns use mixed Italian/English names (`titolo`, `descrizione`, `data`, `nome_ruolo`, etc.) — match schema.sql exactly
- No CI, no Makefile, no linter/formatter config — bare Go toolchain only
- `go.mod` marks all deps as `// indirect` — likely needs `go mod tidy` if new imports are added
