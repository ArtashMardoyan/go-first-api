# go-first-api

A production-ready REST API built with Go, Gin, GORM, and PostgreSQL. Feature-based module architecture inspired by NestJS.

## Stack

- **Go 1.26**
- **Gin** — HTTP router (like Express/Fastify)
- **GORM** — ORM (like TypeORM/Sequelize)
- **PostgreSQL** — database
- **JWT** — authentication
- **golangci-lint** — linting

## Getting started

### 1. Clone the repository

```bash
git clone https://github.com/ArtashMardoyan/go-first-api.git
cd go-first-api
```

### 2. Set up environment

```bash
cp .env.example .env
```

Fill in your values:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=go-first-api
JWT_SECRET=your_secret_key
```

### 3. Run

```bash
go run ./cmd/server/
```

Server starts at `http://localhost:3000`. Database tables are created automatically via AutoMigrate.

## Project structure

```
go-first-api/
├── cmd/
│   └── server/
│       └── main.go               # entry point — wires all dependencies
├── internal/
│   ├── config/
│   │   └── config.go             # loads and validates env variables
│   ├── infrastructure/
│   │   ├── database/             # database connection
│   │   └── middleware/           # JWT middleware
│   ├── modules/
│   │   ├── auth/                 # login, /me endpoint
│   │   ├── user/                 # user CRUD
│   │   └── post/                 # post CRUD
│   └── shared/                   # response helpers, pagination, context keys
├── migrations/
├── .env
├── .golangci.yml
└── go.mod
```

Each module contains: `entity.go`, `dto.go`, `repository.go` (interface), `repository_gorm.go`, `service.go`, `handler.go`

## API

### Auth

| Method | URL        | Auth | Description        |
|--------|------------|------|--------------------|
| POST   | /auth/login | —   | Login, get JWT     |
| GET    | /auth/me   | ✓    | Current user       |

### Users

| Method | URL          | Auth | Description              |
|--------|--------------|------|--------------------------|
| GET    | /users       | ✓    | List users (paginated)   |
| GET    | /users/:id   | ✓    | Get user by ID           |
| POST   | /users       | —    | Register (create user)   |
| PATCH  | /users       | ✓    | Update own profile       |
| DELETE | /users       | ✓    | Delete own account       |

### Posts

| Method | URL                   | Auth | Description                  |
|--------|-----------------------|------|------------------------------|
| GET    | /posts                | ✓    | List all posts (paginated)   |
| GET    | /posts/:id            | ✓    | Get post by ID (with user)   |
| GET    | /posts/user/:userId   | ✓    | Get posts by user            |
| POST   | /posts                | ✓    | Create post                  |
| PATCH  | /posts/:id/status     | ✓    | Update post status (owner)   |
| DELETE | /posts/:id            | ✓    | Delete post (owner)          |

## Response format

**Success:**
```json
{
  "success": true,
  "message": "user retrieved",
  "data": {},
  "error": null
}
```

**Error:**
```json
{
  "success": false,
  "message": "user not found",
  "data": null,
  "error": "user not found"
}
```

**Paginated:**
```json
{
  "success": true,
  "message": "posts retrieved",
  "data": {
    "items": [],
    "meta": {
      "page": 1,
      "limit": 10,
      "total": 100
    }
  }
}
```

## Example requests

**Login:**
```bash
curl -X POST http://localhost:3000/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}'
```

**Create user:**
```bash
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Artash", "email": "artash@example.com", "password": "password123", "age": 25}'
```

**Create post (authenticated):**
```bash
curl -X POST http://localhost:3000/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"title": "Hello", "body": "My first post"}'
```

## Development

```bash
go build ./...          # compile check
golangci-lint run       # lint
golangci-lint fmt       # auto-fix formatting
```