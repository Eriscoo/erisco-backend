# Migration — Erisco Blog

Migrations use [golang-migrate](https://github.com/golang-migrate/migrate) (v4).  
SQL files live in `migrations/`. Migrations are run manually, on-demand — they do not run automatically on server start.

## Quick reference

| Action | Command |
|---|---|
| Apply all pending migrations | `go run ./cmd/migrate/` |
| Create new migration pair | `migrate create -ext sql -dir migrations -seq <name>` |
| Roll back last migration | `migrate -path migrations -database "$DATABASE_URL" down 1` |
| Roll back all migrations | `migrate -path migrations -database "$DATABASE_URL" down` |
| View current version | `migrate -path migrations -database "$DATABASE_URL" version` |
| Force version (dirty fix) | `migrate -path migrations -database "$DATABASE_URL" force <N>` |

> `$DATABASE_URL` is read from `.env`. See the Environment section below.

---

## How to run migrations

### 1. Go command (recommended)

Built-in command, reads `.env` automatically. For applying pending migrations:

```powershell
go run ./cmd/migrate/
```

### 2. migrate CLI (advanced)

For operations beyond `up` (down, force, version). Set `$DATABASE_URL` first.

**PowerShell:**
```powershell
$env:DATABASE_URL = 'postgres://user:password@localhost:5432/dbname?sslmode=disable'
migrate -path migrations -database $env:DATABASE_URL up
migrate -path migrations -database $env:DATABASE_URL down 1
migrate -path migrations -database $env:DATABASE_URL version
```

**Linux/macOS:**
```bash
export DATABASE_URL='postgres://user:password@localhost:5432/dbname?sslmode=disable'
migrate -path migrations -database "$DATABASE_URL" up
migrate -path migrations -database "$DATABASE_URL" down 1
migrate -path migrations -database "$DATABASE_URL" version
```

---

## Creating a new migration

```powershell
migrate create -ext sql -dir migrations -seq add_comments_table
# generates:
#   migrations/<N>_add_comments_table.up.sql
#   migrations/<N>_add_comments_table.down.sql
```

Edit `.up.sql` for the DDL changes and `.down.sql` for the reverse.  
Run `go run ./cmd/migrate/` to apply.

---

## Environment

`DATABASE_URL` is stored in `.env` (do not commit this file). Format:

```
DATABASE_URL=postgres://<user>:<password>@<host>:<port>/<db>?sslmode=disable
```

---

## Troubleshooting

| Problem | Fix |
|---|---|
| `error: no change` | All migrations are already applied. Check with `migrate version`. |
| `dirty=true` | A migration failed part-way. Fix the issue, then run `migrate force <version>`. |
| `password authentication failed` | The password in `DATABASE_URL` is wrong. Check `docker-compose.yml`. |
| `SSL is not enabled` | Add `?sslmode=disable` to the connection string. |

---

## Current migrations

| Version | Description |
|---|---|
| 1 | init_schema — all CREATE TABLE statements |
| 2 | seed_data — user_role, users, categories, tags |

## File structure

```
cmd/migrate/main.go                          → command to run migrations
migrations/embed.go                          → embeds all *.sql files
migrations/000001_init_schema.{up,down}.sql  → DDL
migrations/000002_seed_data.{up,down}.sql    → seed data
internal/infrastructure/persistence/migration.go → RunMigrations logic
```
