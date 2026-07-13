# Erisco Blog — Backend

Golang v1.26.5 + Gin + PostgreSQL 17

## Setup

### 1. Start Database
```bash
docker compose up -d
```

### 2. Run Backend
```bash
go run ./cmd/server
# → http://localhost:8080
```

## Environment Variables

Buat file `.env`:

```env
DATABASE_URL=postgres://username:password@localhost:5432/dbname?sslmode=disable&timezone=Asia/Jakarta
JWT_SECRET=super-secret-key-change-in-production
PORT={PORTHERE}
UPLOAD_DIR={DIRECTORYHERE}
APP_ENV={APPENVHERE}
```

## API Endpoints

Semua endpoint di bawah `/api/v1/`:

| Method | Endpoint | Auth | Deskripsi |
|---|---|---|---|
| POST | `/register` | No | Register user |
| POST | `/login` | No | Login → JWT token |
| GET | `/me` | Bearer | User info (id, name, email) |
| GET | `/tags` | No | List tags |
| POST | `/tags` | Bearer | Create tag |
| PUT | `/tags/{id}` | Bearer | Update tag |
| DELETE | `/tags/{id}` | Bearer | Delete tag |
| GET | `/categories` | No | List categories |
| POST | `/categories` | Bearer | Create category |
| PUT | `/categories/{id}` | Bearer | Update category |
| DELETE | `/categories/{id}` | Bearer | Delete category |
| GET | `/profile/{id}` | Bearer | Get user profile |
| PUT | `/profile/{id}` | Bearer | Create/update profile |
| POST | `/upload` | Bearer | Upload image (multipart) |

## Swagger

`GET /swagger/index.html` — dokumentasi API interaktif.
