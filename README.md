# go-user-api

REST API на Go (Gin + GORM + PostgreSQL), построенный по структуре NestJS.

## Стек

- **Gin** — HTTP-роутер (аналог Express/Fastify)
- **GORM** — ORM (аналог TypeORM/Sequelize)
- **PostgreSQL** — база данных
- **godotenv** — переменные окружения из `.env`

## Запуск

### 1. Клонируй репозиторий

```bash
git clone <repo-url>
cd go-user-api
```

### 2. Настрой окружение

```bash
cp .env.example .env
```

Заполни `.env` своими значениями:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=userapi
```

### 3. Запусти

```bash
go run main.go
```

Сервер стартует на `http://localhost:3000`. Таблица `users` создастся автоматически (AutoMigrate).

## Эндпоинты

| Метод  | URL             | Описание              |
|--------|-----------------|-----------------------|
| GET    | /users          | Список всех юзеров    |
| GET    | /users/:id      | Юзер по ID            |
| POST   | /users          | Создать юзера         |
| PATCH  | /users/:id      | Обновить юзера        |
| DELETE | /users/:id      | Удалить юзера         |

### Примеры запросов

**Создать юзера**
```bash
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Ivan", "email": "ivan@example.com", "age": 25}'
```

**Обновить юзера**
```bash
curl -X PATCH http://localhost:3000/users/<id> \
  -H "Content-Type: application/json" \
  -d '{"name": "Ivan Petrov"}'
```

**Удалить юзера**
```bash
curl -X DELETE http://localhost:3000/users/<id>
```

## Структура проекта

```
.
├── main.go                  # точка входа, DI, подключение к БД
└── internal/
    └── user/
        ├── model.go         # User struct + DTO (аналог entity + dto в NestJS)
        ├── repository.go    # работа с БД через GORM (аналог TypeORM Repository)
        ├── service.go       # бизнес-логика (аналог @Injectable Service)
        └── handler.go       # HTTP-хендлеры (аналог @Controller)
```