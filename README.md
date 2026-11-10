# URL Shortener with Aggregation (WIP)
![Go Version](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go)
![License](https://img.shields.io/badge/License-MIT-green)

A backend service for shortening URLs and aggregating click statistics, written in Go.

**Current status:** authentication and role-based authorization are implemented; link shortening and aggregation are in active development.

## Current Features

- User registration with password hashing via **bcrypt**
- **JWT authentication** — the access token is issued on login and stored in an **HttpOnly, Secure, SameSite=Strict** cookie
- Token validation and role-based access control (**admin / user**) middleware
- Password change endpoint (protected)
- Admin endpoints for listing (paginated) and deleting users
- Health and readiness probes (`/health`, `/ready`)
- PostgreSQL schema and migrations for users
- Docker Compose environment with PostgreSQL and pgAdmin
- Dependency injection via `go.uber.org/dig`
- Configuration through environment variables (`.env`)

## Planned Features

- URL shortening and redirects
- Click aggregation and statistics
- Redis-based rate limiting / caching
- Refresh tokens and token rotation
- Unit and integration tests
- API documentation (OpenAPI/Swagger)

## Tech Stack

- **Language:** Go 1.25
- **HTTP framework:** Echo v4
- **Database:** PostgreSQL (pgx driver)
- **Auth:** JWT (`golang-jwt/jwt/v5`) stored in a cookie, bcrypt (`golang.org/x/crypto`)
- **DI:** `go.uber.org/dig`
- **Config:** `github.com/joho/godotenv`
- **Infrastructure:** Docker, Docker Compose

## Project Structure

```
├── cmd/api/                    # application entrypoint
├── internal/ 
│   ├── config/                 # environment configuration
│   ├── database/               # database connection setup
│   ├── domain/                 # domain models and interfaces
│   ├── dto/                    # request/response DTOs
│   ├── handler/                # HTTP handlers
│   ├── middleware/             # JWT and role-check middleware
│   ├── repositories/postgres/  # PostgreSQL repository implementations
│   ├── routes/                 # route registration
│   └── server/                 # server bootstrap
├── migrations/                 # SQL migrations
├── .env.example
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── go.sum
```

## Getting Started

### Prerequisites

- Go 1.25+
- Docker and Docker Compose

### 1. Clone the repository

```
git clone https://github.com/Anti6eptik/URL-shortener_with_agregation-WIP-.git
cd URL-shortener_with_agregation-WIP-
```
### 2. Configure environment
```
cp .env.example .env
```

  Fill in the values in `.env`. At minimum, set:

* `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`

* `JWT_SECRET`

* `BCRYPT_COST` (recommended 12–14)

* `ACCESS_TOKEN_EXPIRY`

### 3. Start infrastructure
```
docker compose up -d
```
  This starts PostgreSQL and pgAdmin. pgAdmin will be available on the port specified by `PGADMIN_PORT` (default `5050`).

### 4. Run migrations
Apply the SQL files from `migrations/` to your PostgreSQL instance (via pgAdmin, `psql`, or your preferred migration tool).

### 5. Run the application
```
go run ./cmd/api
```
The API will start on the port specified by `PORT` in `.env`.

## Authentication

The API uses **JWT access tokens delivered as cookies**.

- On successful `POST /api/login`, the server returns the user object in the response body **and** sets an auth cookie.
- Cookie parameters:
  - **Name:** `access_token`
  - **Path:** `/`
  - **HttpOnly:** `true` — not accessible from JavaScript, mitigates XSS token theft
  - **Secure:** `true` — sent only over HTTPS
  - **SameSite:** `Strict` — mitigates CSRF
  - **Expires:** set from `ACCESS_TOKEN_EXPIRY`
- Protected endpoints (`/api/authed/*`) require this cookie.
- The token is validated by middleware, which also checks the user's role for `/admin` routes.
- Because auth is cookie-based, browser clients get the token automatically; CLI clients must persist and resend cookies (see `curl` examples below).

> **Local development note:** `Secure: true` means browsers will not send the cookie over plain `http://`. For browser-based testing on `localhost`, either serve the API over HTTPS or temporarily set `Secure: false` in `SetAccessTokenCookie` for local runs only. `curl` and most API clients ignore the `Secure` flag and will send the cookie over HTTP, so the examples below work as-is.

> **Security notes:** `HttpOnly` + `SameSite=Strict` + `Secure` covers the most common cookie-based auth risks (XSS token theft and CSRF). Planned hardening: refresh tokens with rotation and an explicit CSRF token for state-changing endpoints.
## API Reference
Base URL: `http://localhost:8080`

### Service
`GET /health` — liveness probe (public)

`GET /ready` — readiness probe, checks dependencies (public)

### Auth
`POST /api/register` — register a new user (public)

`POST /api/login` — log in and receive a JWT access token (public)

`PUT /api/authed/change_password` — change password (authenticated)

### Admin
`GET /api/authed/admin/users` — list users, paginated (admin only)

`DELETE /api/authed/admin/delete/:id` — delete user by integer ID (admin only)

### Links (in development)
Not implemented yet. Planned endpoints: create short link, redirect by code, click statistics.

## Error Format
All errors are returned in a uniform shape:
```
{
  "error": "human-readable message",
  "kind": "optional_error_category"
}
```

## Usage Examples
For `curl`, use a cookie jar so the auth cookie is stored after login and sent with later requests.

Note: the auth cookie is marked `Secure`, but `curl` ignores that flag and will still send it over `http://localhost`.
### Health check
```
curl http://localhost:8080/health
```

### Register a user
```
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "strong-password-123"
  }'
```
Response:
```
{
  "id": 1,
  "email": "user@example.com",
  "created_at": "2026-01-01T12:00:00Z"
}
```
### Log in (stores the auth cookie)
```
curl -c cookies.txt -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "strong-password-123"
  }'
```
Response body (the JWT itself is set as a cookie, not returned in JSON):
```
{
  "id": 1,
  "email": "user@example.com",
  "created_at": "2026-01-01T12:00:00Z"
}
```
The file `cookies.txt` now contains the auth cookie. Reuse it with `-b cookies.txt` on all protected requests.

### Change password (authenticated)
```
curl -b cookies.txt -X PUT http://localhost:8080/api/authed/change_password \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "old_password": "strong-password-123",
    "new_password": "even-stronger-password-456"
  }'
```
### List users (admin only, paginated)
Log in first as an admin to obtain an admin cookie:

```
curl -c admin_cookies.txt -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin-password"
  }'

curl -b admin_cookies.txt "http://localhost:8080/api/authed/admin/users?page=1&limit=10"
```
Response:
```
{
  "users": [
    {
      "id": 1,
      "email": "admin@example.com",
      "role": "admin",
      "created_at": "2026-01-01T12:00:00Z",
      "updated_at": "2026-01-01T12:00:00Z"
    },
    {
      "id": 2,
      "email": "user@example.com",
      "role": "user",
      "created_at": "2026-01-01T12:00:00Z",
      "updated_at": "2026-01-01T12:00:00Z"
    }
  ],
  "total": 2,
  "page": 1,
  "total_pages": 1
}
```
### Delete a user (admin only)
```
curl -b admin_cookies.txt -X DELETE http://localhost:8080/api/authed/admin/delete/2
```
## Authorization errors
No cookie / invalid or expired token:
```
{
  "error": "unauthorized",
  "kind": "unauthorized"
}
```
Insufficient role (a `user` tries to access an admin route):
```
{
  "error": "forbidden",
  "kind": "forbidden"
}
```
Roadmap
- [x] User registration and login
- [x] JWT authentication (cookie-based)
- [x] Role-based middleware (admin/user)
- [x] Password change
- [x] Admin user management (list with pagination, delete)
- [x] Health and readiness probes
- [x] PostgreSQL schema and migrations
- [x] Docker Compose environment
- [ ] URL shortening
- [ ] Redirect handler
- [ ] Click aggregation and statistics
- [ ] Refresh tokens and token rotation
- [ ] Additional CSRF tokens for state-changing endpoints
- [ ] Redis integration
- [ ] Tests
- [ ] API documentation
## License
MIT
