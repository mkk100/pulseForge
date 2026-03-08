# pulseForge

Minimal Go API for users and posts, backed by Postgres.

The app uses a small layered structure:

- `internal/http`: HTTP handlers, request parsing, auth middleware
- `internal/service`: business logic / orchestration
- `internal/repo`: SQL and database access
- `cmd/api`: app bootstrap and integration test
- `migrations`: schema setup and teardown

## Stack

- Go
- `net/http`
- `pgxpool`
- Postgres via Docker Compose
- Minimal JWT-style auth for authenticated post creation

## Requirements

- Go installed locally
- Docker and Docker Compose available

## Environment

Current local DB config in [`.env`](/Users/minkaungkhant/Desktop/Full%20Time%20Job%20Hunt/Interview%20Prep/Projects/pulseForge/.env):

```env
DATABASE_URL=postgres://pulseforge:pulseforge@localhost:5433/pulseforge?sslmode=disable
```

Optional:

```env
JWT_SECRET=your-local-secret
```

If `JWT_SECRET` is not set, the app falls back to a dev secret in [`internal/http/auth.go`](/Users/minkaungkhant/Desktop/Full%20Time%20Job%20Hunt/Interview%20Prep/Projects/pulseForge/internal/http/auth.go).

## Run Locally

Start Postgres:

```bash
make db-create
```

Apply migrations:

```bash
make db-insert
```

Run the API:

```bash
make backend
```

The server listens on `http://localhost:8080`.

## API

### Health

```bash
curl -i http://localhost:8080/health
```

### Create user

```bash
curl -i -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"userName":"thomas"}'
```

Example response:

```json
{"userId":1}
```

### Get user ID by name

```bash
curl -i "http://localhost:8080/users/id?name=thomas"
```

Example response:

```json
{"userId":1,"userName":"thomas"}
```

### Login

Current login is intentionally simple: if the username exists, the server issues a signed token.

```bash
curl -i -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"userName":"thomas"}'
```

Example response:

```json
{"token":"<token>","userId":1}
```

### Create post

`POST /posts` is protected. Send the token in the `Authorization` header.

```bash
curl -i -X POST http://localhost:8080/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{"title":"hello","description":"world"}'
```

Example response:

```json
{"postId":1}
```

### List posts

```bash
curl -i "http://localhost:8080/posts?limit=10"
```

Example response:

```json
{
  "posts": [
    {
      "id": 1,
      "title": "hello",
      "description": "world",
      "userId": 1,
      "createdAt": "2026-03-08T00:00:00Z"
    }
  ]
}
```

## Auth Flow

1. Create a user with `POST /users`.
2. Login with `POST /login`.
3. Receive a signed token.
4. Send that token as `Authorization: Bearer <token>` on protected routes.
5. Middleware verifies the token and puts `userID` into request context.
6. `POST /posts` uses that authenticated `userID` instead of trusting a client-supplied one.

## Database

Migrations live in [`migrations/init.up.sql`](/Users/minkaungkhant/Desktop/Full%20Time%20Job%20Hunt/Interview%20Prep/Projects/pulseForge/migrations/init.up.sql) and [`migrations/init.down.sql`](/Users/minkaungkhant/Desktop/Full%20Time%20Job%20Hunt/Interview%20Prep/Projects/pulseForge/migrations/init.down.sql).

Reset the schema:

```bash
make db-reset
make db-insert
```

## Tests

Integration test:

- [`cmd/api/main_integration_test.go`](/Users/minkaungkhant/Desktop/Full%20Time%20Job%20Hunt/Interview%20Prep/Projects/pulseForge/cmd/api/main_integration_test.go)

What it covers:

- connects to real Postgres
- truncates tables for repeatability
- creates a user
- logs in
- creates a post with auth
- fetches posts and checks the inserted data

Run tests:

```bash
make test
```

## Project Structure

```text
cmd/api/                app entrypoint + integration test
internal/http/          handlers and auth middleware
internal/service/       business logic layer
internal/repo/          SQL and pgx access
migrations/             schema files
docker-compose.yml      local Postgres
Makefile                local commands
```

## Notes

- The current auth flow is identity wiring, not full production auth.
- Login currently checks that a username exists and then issues a token.
- A realistic next step is password hashing plus real credential verification.
