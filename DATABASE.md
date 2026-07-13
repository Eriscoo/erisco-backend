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

| Column | Type | Fungsi |
|---|---|---|
| `id` | SERIAL PK | Primary key, auto increment |
| `title` | VARCHAR(255) | Judul artikel (wajib) |
| `slug` | VARCHAR(255) | Versi URL-friendly dari title. Unik, jadi tidak boleh ada judul yang sama slug-nya. Contoh: `"cara-memasak-nasi-goreng"` |
| `body` | TEXT | Isi konten artikel (bisa panjang) |
| `image_url` | VARCHAR(500) | Path ke file gambar thumbnail/header artikel. Disimpan sebagai string path, bukan file binary |
| `categories` | TEXT | Comma-separated ID dari tabel categories. Contoh: `"1,3,5"` |
| `tags` | TEXT | Comma-separated ID dari tabel tags. Contoh: `"2,4,7"` |
| `created_by` | INT | Foreign key ke tabel users. Menentukan siapa penulis artikel. CASCADE → jika user dihapus, post-nya ikut terhapus |
| `status` | VARCHAR(20) | Status publikasi. Draft → masih diedit, Published → sudah terbit, Archived → sudah tidak ditampilkan |
| `published_at` | TIMESTAMP | Tanggal & waktu artikel diterbitkan. NULL selama masih draft, diisi otomatis saat status berubah jadi 'published' |
| `created_at` | TIMESTAMP | Tanggal artikel dibuat (otomatis) |
| `updated_at` | TIMESTAMP | Tanggal artikel terakhir diedit (otomatis) |

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

| Column | Type | Keterangan |
|---|---|---|
| user_id | INT PK | FK → users(id) CASCADE |
| bio | VARCHAR(200) | Bio singkat |
| avatar_url | VARCHAR(500) | Path foto profil |
| website | VARCHAR(255) | Website |
| location | VARCHAR(100) | Lokasi |
| phone | VARCHAR(20) | No telepon |
| created_at | TIMESTAMP | Auto |
| updated_at | TIMESTAMP | Auto |
