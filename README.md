# Gin Event App

Event management REST API built with Go and Gin. Supports user registration/login (JWT), CRUD for events, and attendee management. Uses SQLite with SQL migrations and ships with Swagger docs and a Bruno API collection.

## Features

- JWT auth: Register and login to receive a bearer token
- Events: Create, read, update, delete
- Attendees: Add/remove users to/from events, list attendees of an event, list events for a user
- SQLite storage with SQL migrations
- Auto-loaded env vars via .env
- Swagger UI at /swagger
- Bruno collection included in `burno/gin-event-app`

## Tech stack

- Go (Gin, bcrypt, golang-jwt)
- SQLite (mattn/go-sqlite3)
- Migrations (golang-migrate)
- Swagger (swaggo)

## Project structure

```
cmd/
  api/          # HTTP server, routes, handlers, middleware
  migrate/      # Migration runner and SQL files
internal/
  database/     # Models for users, events, attendees (raw SQL)
  env/          # Env helpers
  helpers/      # Context and response helpers
burno/gin-event-app  # Bruno API collection
```

## Requirements

- Go (as specified in go.mod)

## Configuration

Create a `.env` in the project root (auto-loaded):

```
PORT=8000
JWT_SECRET=your-super-secret
```

Defaults: `PORT=8000`, `JWT_SECRET=secret-123123`.

## Database & migrations

The API uses a local SQLite file `data.db` in the project root.

Run migrations:

```
# Up
go run ./cmd/migrate up

# Down
go run ./cmd/migrate down
```

## Run the API

```
go run ./cmd/api
```

Server starts on `http://localhost:8000`.

## API docs (Swagger)

- Swagger UI: `http://localhost:8000/swagger/index.html`

## Auth

- Register: `POST /api/v1/auth/register` (email, password, name)
- Login: `POST /api/v1/auth/login` → returns `{ token }`
- For protected routes, set header: `Authorization: Bearer <token>`

## Endpoints overview

Public

- GET `/api/v1/events` — list events
- GET `/api/v1/events/:id` — get event by id
- GET `/api/v1/events/:id/attendees` — list attendees for an event
- GET `/api/v1/attendees/:id/events` — list events by user
- POST `/api/v1/auth/register` — register
- POST `/api/v1/auth/login` — login

Protected (Bearer token)

- POST `/api/v1/events` — create event (owner = current user)
- PUT `/api/v1/events/:id` — update owned event
- DELETE `/api/v1/events/:id` — delete owned event
- POST `/api/v1/events/:id/attendees/:userId` — add attendee (owner only)
- DELETE `/api/v1/events/:id/attendees/:userId` — remove attendee (owner only)

Request/response schemas are documented in Swagger and in the Bruno collection.

## Bruno collection

A ready-to-use Bruno collection is included under `burno/gin-event-app`. Import it into Bruno, set the base URL to `http://localhost:8000`, and run requests.
