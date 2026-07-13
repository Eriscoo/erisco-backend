# Infrastructure — Erisco Blog

## Stack

| Layer | Tech |
|---|---|
| Database | PostgreSQL 17 (Docker) |
| Backend | Go 1.26 + Gin + Clean Architecture |
| Auth | JWT Bearer Token (HS256, 72h expiry) |
| Env | `.env` + godotenv |
| Upload | Multipart → `./uploads/profile/` |

## Docker

```bash
docker compose up -d    # start
docker compose down     # stop
docker compose down -v  # stop + delete volume (data)
```

| Service | Image | Container | Port | Volume |
|---|---|---|---|---|
| postgres | `postgres:17-alpine` | `eriscoo-postgres` | `5432` | `eriscoo-postgres-volume` |

## Environment

| Variable | Default |
|---|---|
| `DATABASE_URL` | `postgres://username:password@localhost:5432/dbname?sslmode=disable` |
| `JWT_SECRET` | `change-me-in-production` |
| `PORT` | `8080` |
| `UPLOAD_DIR` | `./uploads` |

## Backend Architecture

```
cmd/server/main.go              → DI wiring
│
├── transport/                   → HTTP layer
│   ├── handler/
│   │   ├── auth/                 → Register, Login, GetMe
│   │   ├── tags/                 → CRUD tags
│   │   ├── categories/           → CRUD categories
│   │   ├── profile/              → GET/PUT profile
│   │   └── upload/               → POST /upload
│   ├── middleware/auth.go       → Bearer token guard
│   └── router/router.go        → Route definitions
│
├── application/                 → Business logic / use cases
│   ├── port.go                  → Interfaces (repositories, token service)
│   ├── auth/service.go          → Auth use case
│   ├── tags/service.go          → Tags use case
│   ├── categories/service.go    → Categories use case
│   └── profile/service.go       → Profile use case
│
├── domain/                      → Core entities & errors
│   ├── user.go / tag.go / category.go / user_profile.go
│   └── errors.go
│
└── infrastructure/              → Adapters
    ├── config/config.go         → Load .env
    ├── auth/jwt_service.go      → JWT generate & validate
    ├── persistence/             → PostgreSQL (user, tag, category, profile repos)
    └── docs/                    → Generated swagger docs
```

### Dependency Rule

```
domain ← application ← infrastructure, transport
```

All layers depend inward. Domain knows nothing.

## API

All endpoints under `/api/v1/`. See `API.md` for request & response details.
