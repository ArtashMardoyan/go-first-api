# go-first-api

Gin + GORM + PostgreSQL REST API. Go 1.26.

## Commands

```bash
go run ./cmd/server/         # start server (port 3000)
go build ./...               # compile check
golangci-lint run            # lint
golangci-lint fmt            # auto-fix formatting
```

## Architecture

```
cmd/server/main.go           — wiring only: load config, init deps, start server
internal/
  config/config.go           — all env vars in one place (config.Load())
  infrastructure/
    database/db.go           — Connect(*config.DBConfig) *gorm.DB
    middleware/auth.go       — JWT(userRepo, jwtSecret) gin.HandlerFunc
  modules/                   — feature modules; new features go here
    user/                    — entity, dto, status, repository (interface + gorm), service, handler
    post/
    auth/
  shared/                    — response helpers, pagination, context key
migrations/                  — SQL migration files (not yet wired)
```

## Adding a new module

1. Create `internal/modules/<name>/` with: `entity.go`, `dto.go`, `repository.go` (interface), `repository_gorm.go`, `service.go`, `handler.go`
2. Add to `db.AutoMigrate(...)` in `main.go`
3. Wire repo → service → handler in `main.go`
4. Call `handler.RegisterRoutes(r, jwtMiddleware)`

## Key conventions

- **Errors**: sentinel vars (`var ErrNotFound = errors.New(...)`) in the package that owns the type; compare with `errors.Is`
- **Context**: always thread `ctx context.Context` through repo and service methods
- **Response format**: always use `shared.OK` / `shared.Error` / `shared.Created`
- **Pagination**: bind `shared.PaginationQuery` from query params, call `.Normalize()`, return `shared.PaginatedResult[T]`
- **Auth**: JWT stored in `shared.ContextUserKey`; retrieve with two-value type assertion `val.(user.User)` — never one-value
- **Repository**: interface in `repository.go`, GORM impl in `repository_gorm.go` (unexported struct)
- **Imports**: 3 groups — stdlib / external / internal (`go-first-api/...`) — enforced by goimports

## Response format

```json
{ "success": true,  "message": "...", "data": {},   "error": null }
{ "success": false, "message": "...", "data": null, "error": "..." }
```

Paginated data inside `data`:
```json
{ "items": [], "meta": { "page": 1, "limit": 10, "total": 100 } }
```

## Environment variables

```
DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME
JWT_SECRET   — required, server refuses to start without it
```

## Dependency direction

```
handler → service → repository interface → gorm repository → database
```

- Services depend only on interfaces, never on concrete types
- Middleware must not depend on services — only on `user.Repository`
- No global variables