# Go Events REST API

A small REST API built with `gin` and `sqlite` for user authentication, event management, and event registration.

## What This Project Does

- User signup and login
- JWT-based authentication
- CRUD operations for events
- Register/cancel registration for events
- SQLite persistence with automatic table creation

## Tech Stack

- Go
- [gin-gonic/gin](https://github.com/gin-gonic/gin)
- SQLite (`github.com/mattn/go-sqlite3`)
- JWT (`github.com/golang-jwt/jwt/v5`)

## Project Structure

- `main.go`: application entry point
- `db/`: DB initialization and schema creation
- `models/`: domain models and persistence methods
- `routes/`: HTTP handlers and route registration
- `middleware/`: auth middleware
- `utils/`: password hashing and JWT utilities
- `api-test/`: HTTP request samples for quick manual testing

## Prerequisites

- Go (use the version declared in `go.mod`)
- CGO enabled (required by `go-sqlite3`)

## Configuration

Environment variables:

- `JWT_SECRET` (optional): secret used to sign JWT tokens

If `JWT_SECRET` is not set, the app uses a development fallback (`dev-secret-change-me`). Set your own secret for any non-local usage.

## Run the API

```bash
go run .
```

The server starts on `http://localhost:8080`.

On startup, SQLite database file `api.db` is created in the project root (if missing), and tables are initialized automatically.

## Test

```bash
go test ./...
```

## API Endpoints

Public:

- `POST /signup`
- `POST /login`
- `GET /events`
- `GET /events/:id`

Authenticated (send `Authorization: Bearer <token>`):

- `POST /events`
- `PUT /events/:id`
- `DELETE /events/:id`
- `POST /events/:id/register`
- `DELETE /events/:id/register`

## Example Flow (curl)

Create user:

```bash
curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","email":"alice@example.com","password":"password123"}'
```

Login:

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"password123"}'
```

Create event (replace `<TOKEN>`):

```bash
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{
    "name":"Go Meetup",
    "description":"Monthly backend meetup",
    "location":"Cluj",
    "date":"2026-08-10T18:00:00Z"
  }'
```

## Notes

- Users can update/delete only events they created.
- Duplicate event registration by the same user returns conflict.
- Passwords are hashed with bcrypt and never returned in API responses.
