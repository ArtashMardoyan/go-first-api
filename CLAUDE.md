# go-first-api

Gin + GORM + PostgreSQL REST API. Go 1.26.

## Critical workflow rules

**NEVER commit or push without explicit user instruction.** Wait for the user to say "commit" or "commit and push" before running any `git commit` or `git push`. No exceptions.

**Large tasks require a plan first.** For any new feature, refactor, or cross-cutting change:
1. Describe the approach and list files to change
2. Wait for user approval before writing any code

**Keep the Postman collection in sync.** Whenever you add, remove, or change a route (method, path, request body, or query params), update `postman/go-first-api.postman_collection.json` in the same change — add/edit the matching request under the right folder. See "API documentation (Postman)" below.

## Commands

```bash
go run ./cmd/server/         # start server (port 3000)
go build ./...               # compile check
golangci-lint run            # lint
golangci-lint fmt            # auto-fix formatting
./scripts/deploy.sh          # full deploy: build → ECR push → App Runner → Postman sync
./scripts/sync-postman.sh    # push Postman collection to the cloud workspace only
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
2. Add to `db.AutoMigrate(...)` in `cmd/server/main.go`
3. Wire repo → service → handler in `main.go`
4. Call `handler.RegisterRoutes(r, jwtMiddleware)`
5. Add a folder + requests for the module's routes to `postman/go-first-api.postman_collection.json`

## API documentation (Postman)

`postman/` holds the importable API collection — keep it current whenever routes change.

- `go-first-api.postman_collection.json` — requests grouped into folders per module (Auth, Users, Posts).
- `scripts/sync-postman.sh` — pushes the collection to the Postman cloud workspace via the Postman API. Run `POSTMAN_API_KEY=PMAK-... ./scripts/sync-postman.sh` after editing the collection (or as a deploy step). The API key is read from the env, never hardcoded.
- `go-first-api.local.postman_environment.json` / `go-first-api.prod.postman_environment.json` — environments (display names "06.1 GO First API Local" / "06.2 GO First API Prod") with `API_URL`, `accessToken`, `userId`, `postId`, `userEmail`, `userName` (Local defaults `API_URL` to `http://localhost:3000`; Prod is left blank).

Conventions for new requests:
- **Order requests within each folder by method: GET → POST → PATCH → DELETE** (mirrors the Go handlers).
- Auth is **collection-level** Bearer `{{accessToken}}`; requests inherit it. Public routes (login, create user) set `"auth": { "type": "noauth" }`.
- Use `{{API_URL}}` as the base and `{{userId}}` / `{{postId}}` for path params; capture new ids in a `test` script via `pm.environment.set(...)`.
- Request bodies must match the DTO json tags exactly (lowercase: `name`, `email`, `title`, ...), satisfy validation (e.g. `password` min 8, `age` min 1), and use the real route (e.g. post status is `PATCH /posts/:id/status`).

## Code style

Separate logical steps inside functions with a blank line — same as TypeScript style:

```go
func (s *Service) Create(ctx context.Context, dto CreateDTO) (User, error) {
    // step 1: validate
    hashed, err := bcrypt.GenerateFromPassword(...)
    if err != nil {
        return User{}, err
    }

    // step 2: build entity
    u := User{Name: dto.Name, ...}

    // step 3: persist
    if err := s.repo.Create(ctx, &u); err != nil {
        return User{}, err
    }

    return u, nil
}
```

Rule: one blank line between each distinct phase (validate → build → persist → return).
`gofumpt` enforces no blank line at start/end of a block, but allows them in the middle.

## Key conventions

- **Errors**: sentinel vars (`var ErrNotFound = errors.New(...)`) in the package that owns the type; compare with `errors.Is`, never `err.Error() == "..."`
- **Context**: always thread `ctx context.Context` through repo and service methods
- **Response format**: always use `shared.OK` / `shared.Error` / `shared.Created`
- **Pagination**: bind `shared.PaginationQuery` from query params, call `.Normalize()`, return `shared.PaginatedResult[T]`
- **Auth**: JWT stored in `shared.ContextUserKey`; retrieve with two-value type assertion `u, ok := val.(user.User)` — never one-value (panics)
- **Repository**: interface in `repository.go`, GORM impl in `repository_gorm.go` (unexported struct)
- **Imports**: 3 groups separated by blank lines — stdlib / external / internal (`go-first-api/...`) — enforced by goimports

## Common mistakes to avoid

- **Never** use `err.Error() == "..."` — use sentinel errors and `errors.Is`
- **Never** read `os.Getenv` outside of `config/config.go`
- **Never** store JWT secret in a global variable — pass via constructor
- **Never** use one-value type assertion `val.(T)` — always `v, ok := val.(T)`
- **Always** call `q.Normalize()` before using `PaginationQuery`
- **Always** return `[]T{}` not `nil` for empty lists in paginated results

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
- Middleware depends only on `user.Repository` interface, not on services
- No global variables
- No circular dependencies