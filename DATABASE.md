# Database — Erisco Blog

## Connection
| Item | Value |
|---|---|
| Host | `localhost` |
| Port | `5432` |
| Database | `eriscoodb` |
| User | `eriscoo` |
| Password | `********` |

## Tables

### users

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role_id INT REFERENCES user_role(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

| Column | Type | Constraints |
|---|---|---|
| id | SERIAL | PK, auto increment |
| name | VARCHAR(100) | NOT NULL |
| email | VARCHAR(255) | UNIQUE, NOT NULL |
| password_hash | VARCHAR(255) | NOT NULL |
| role_id | INT | FK → user_role(id) |
| created_at | TIMESTAMP | DEFAULT NOW() |
| updated_at | TIMESTAMP | DEFAULT NOW() |

### user_role

```sql
CREATE TABLE user_role (
    id SERIAL PRIMARY KEY,
    role_name VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

| Column | Type | Constraints |
|---|---|---|
| id | SERIAL | PK, auto increment |
| role_name | VARCHAR(50) | UNIQUE, NOT NULL |
| created_at | TIMESTAMP | DEFAULT NOW() |
| updated_at | TIMESTAMP | DEFAULT NOW() |

**Seed data:**
```sql
INSERT INTO user_role (role_name) VALUES ('admin'), ('user');
```

### posts

```sql
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    body TEXT,
    image_url VARCHAR(500),
    categories TEXT,
    tags TEXT,
    created_by INT REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) DEFAULT 'draft',
    published_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

| Column | Type | Description |
|---|---|---|
| `id` | SERIAL PK | Primary key, auto increment |
| `title` | VARCHAR(255) | Article title (required) |
| `slug` | VARCHAR(255) | URL-friendly version of title. Unique, so no two articles share the same slug. Example: `"how-to-cook-fried-rice"` |
| `body` | TEXT | Article content (can be long) |
| `image_url` | VARCHAR(500) | Path to article thumbnail/header image. Stored as string path, not binary |
| `categories` | TEXT | Comma-separated IDs from categories table. Example: `"1,3,5"` |
| `tags` | TEXT | Comma-separated IDs from tags table. Example: `"2,4,7"` |
| `created_by` | INT | Foreign key to users table. Determines the author. CASCADE → if user is deleted, their posts are also deleted |
| `status` | VARCHAR(20) | Publication status. Draft → still editing, Published → already published, Archived → no longer displayed |
| `published_at` | TIMESTAMP | Date & time article was published. NULL while draft, auto-filled when status changes to 'published' |
| `created_at` | TIMESTAMP | Date article was created (auto) |
| `updated_at` | TIMESTAMP | Date article was last edited (auto) |

### categories

```sql
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);
```

| Column | Type | Constraints |
|---|---|---|
| id | SERIAL | PK |
| name | VARCHAR(100) | UNIQUE, NOT NULL |

**Seed data:**
```sql
INSERT INTO categories (name) VALUES ('technology'), ('programming'), ('design'), ('business'), ('lifestyle');
```

### tags

```sql
CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);
```

| Column | Type | Constraints |
|---|---|---|
| id | SERIAL | PK |
| name | VARCHAR(100) | UNIQUE, NOT NULL |

**Seed data:**
```sql
INSERT INTO tags (name) VALUES ('golang'), ('react'), ('javascript'), ('tutorial'), ('technology');
```

### user_profile

```sql
CREATE TABLE user_profile (
    user_id INT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    bio VARCHAR(200),
    avatar_url VARCHAR(500),
    website VARCHAR(255),
    location VARCHAR(100),
    phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

| Column | Type | Description |
|---|---|---|
| user_id | INT PK | FK → users(id) CASCADE |
| bio | VARCHAR(200) | Short bio |
| avatar_url | VARCHAR(500) | Profile photo path |
| website | VARCHAR(255) | Website |
| location | VARCHAR(100) | Location |
| phone | VARCHAR(20) | Phone number |
| created_at | TIMESTAMP | Auto |
| updated_at | TIMESTAMP | Auto |
