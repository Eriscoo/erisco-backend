# API — Erisco Blog

**Base URL:** `http://localhost:8080/api/v1`

## Authentication

All endpoints requiring auth use **Bearer Token** (JWT).
Header:
```
Authorization: Bearer <token>
```

---

## POST /register

Register a new user. Role defaults to `user` (id: 2).

### Request
```json
{
  "name": "string (required)",
  "email": "string (required, valid email)",
  "password": "string (required, min 6 characters)"
}
```

### Response 201 (Created)
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### Response 400 (Bad Request — validation failed)
```json
{
  "error": "Key: 'Name' Error:Field validation for 'Name' failed on the 'required' tag"
}
```

### Response 409 (Conflict — email already registered)
```json
{
  "error": "email already exists"
}
```

---

## POST /login

Login with email & password.

### Request
```json
{
  "email": "string (required, valid email)",
  "password": "string (required)"
}
```

### Response 200 (Success)
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### Response 401 (Unauthorized — invalid email or password)
```json
{
  "error": "invalid email or password"
}
```

---

## GET /me

Get currently logged-in user data.

### Headers
```
Authorization: Bearer <token>
```

### Response 200 (Success)
```json
{
  "user_id": 1
}
```

### Response 401 (Unauthorized — token invalid/expired)
```json
{
  "error": "invalid or expired token"
}
```

---

## GET /tags

Get all tags.

### Response 200 (Success)
```json
[
  {
    "id": 1,
    "name": "golang"
  },
  {
    "id": 2,
    "name": "javascript"
  }
]
```

### Response 500 (Internal Server Error)
```json
{
  "error": "failed to fetch tags"
}
```

---

## POST /tags

Create a new tag (requires auth).

### Headers
```
Authorization: Bearer <token>
Content-Type: application/json
```

### Request
```json
{
  "name": "string (required)"
}
```

### Response 201 (Created)
```json
{
  "id": 6,
  "name": "golang"
}
```

### Response 400 (Bad Request — validation failed)
```json
{
  "error": "Key: 'Name' Error:Field validation for 'Name' failed on the 'required' tag"
}
```

### Response 409 (Conflict — tag already exists)
```json
{
  "error": "tag already exists"
}
```

---

## PUT /tags/{id}

Update tag name by ID (requires auth).

### Headers
```
Authorization: Bearer <token>
Content-Type: application/json
```

### Request
```json
{
  "name": "string (required)"
}
```

### Response 200 (Success)
```json
{
  "id": 1,
  "name": "golang-updated"
}
```

### Response 400 (Bad Request)
```json
{
  "error": "invalid id"
}
```

### Response 404 (Not Found)
```json
{
  "error": "tag not found"
}
```

### Response 409 (Conflict)
```json
{
  "error": "tag already exists"
}
```

---

## DELETE /tags/{id}

Delete tag by ID (requires auth).

### Headers
```
Authorization: Bearer <token>
```

### Response 200 (Success)
```json
{
  "message": "tag deleted"
}
```

### Response 400 (Bad Request)
```json
{
  "error": "invalid id"
}
```

### Response 404 (Not Found)
```json
{
  "error": "tag not found"
}
```

---

## GET /categories

Get all categories.

### Response 200 (Success)
```json
[
  {
    "id": 1,
    "name": "technology"
  },
  {
    "id": 2,
    "name": "programming"
  }
]
```

### Response 500 (Internal Server Error)
```json
{
  "error": "failed to fetch categories"
}
```

---

## POST /categories

Create a new category (requires auth).

### Headers
```
Authorization: Bearer <token>
Content-Type: application/json
```

### Request
```json
{
  "name": "string (required)"
}
```

### Response 201 (Created)
```json
{
  "id": 6,
  "name": "technology"
}
```

### Response 400 (Bad Request)
```json
{
  "error": "Key: 'Name' Error:Field validation for 'Name' failed on the 'required' tag"
}
```

### Response 409 (Conflict)
```json
{
  "error": "category already exists"
}
```

---

## PUT /categories/{id}

Update category name by ID (requires auth).

### Headers
```
Authorization: Bearer <token>
Content-Type: application/json
```

### Request
```json
{
  "name": "string (required)"
}
```

### Response 200 (Success)
```json
{
  "id": 1,
  "name": "devops"
}
```

### Response 400 (Bad Request)
```json
{
  "error": "invalid id"
}
```

### Response 404 (Not Found)
```json
{
  "error": "category not found"
}
```

### Response 409 (Conflict)
```json
{
  "error": "category already exists"
}
```

---

## DELETE /categories/{id}

Delete category by ID (requires auth).

### Headers
```
Authorization: Bearer <token>
```

### Response 200 (Success)
```json
{
  "message": "category deleted"
}
```

### Response 400 (Bad Request)
```json
{
  "error": "invalid id"
}
```

### Response 404 (Not Found)
```json
{
  "error": "category not found"
}
```

---

## GET /profile/{user_id}

Get user profile data by ID.

### Headers
```
Authorization: Bearer <token>
```

### Response 200 (Success)
```json
{
  "user_id": 1,
  "bio": "Full-stack developer",
  "avatar_url": "/uploads/avatar.jpg",
  "website": "https://erisco.dev",
  "location": "Jakarta, Indonesia",
  "phone": "+6281234567890",
  "created_at": "2026-07-11T00:00:00Z",
  "updated_at": "2026-07-11T00:00:00Z"
}
```

### Response 404 (Not Found)
```json
{
  "error": "profile not found"
}
```

---

## PUT /profile/{user_id}

Create or update user profile (can only edit own profile).

### Headers
```
Authorization: Bearer <token>
Content-Type: application/json
```

### Request
```json
{
  "bio": "string",
  "avatar_url": "string",
  "website": "string",
  "location": "string",
  "phone": "string"
}
```

### Response 200 (Success)
```json
{
  "user_id": 1,
  "bio": "...",
  "avatar_url": "...",
  "website": "...",
  "location": "...",
  "phone": "...",
  "created_at": "...",
  "updated_at": "..."
}
```

### Response 403 (Forbidden)
```json
{
  "error": "you can only update your own profile"
}
```

---

## Convention

| Item | Rule |
|---|---|---|
| Base path | `/api/v1/` |
| Auth | Bearer JWT (exp: 72 hours) |
| Request body | JSON, `Content-Type: application/json` |
| Success | 200, 201 |
| Validation failed | 400 |
| Auth failed | 401 |
| Duplicate data | 409 |
| Server error | 500 |
