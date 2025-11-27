# AnimeVerse API

A production-ready RESTful Web API built with Go, serving as the backend for the AnimeVerse mobile application. This project demonstrates a structured Go application with authentication, database interactions, and third-party API integration.

## Features

*   **Synchronous REST API**: Full CRUD operations for managing user anime lists and viewing anime details.
*   **Authentication & Authorization**: Secure user registration, activation, and authentication using tokens. Role-based permissions (e.g., `animes:read`, `animes:write`).
*   **Third-Party Integration**: Integrates with the **Jikan API** (Unofficial MyAnimeList API) to fetch and enrich anime data.
*   **Advanced Querying**: Supports pagination, filtering, and sorting on list endpoints.
*   **Graceful Shutdown**: Ensures the server closes cleanly and finishes in-flight requests.
*   **Security**: Rate limiting, CORS configuration, and environment-based secret management.

## Architecture

The project follows a standard Go project layout:

*   `cmd/api`: Application entry point (`main.go`), routes, middleware, and server configuration.
*   `internal/data`: Data models, validation logic, and database interactions (DAO pattern).
*   `internal/mailer`: Email sending functionality (e.g., for user activation).
*   `internal/validator`: Input validation helpers.
*   `migrations`: SQL database migration files.

## Prerequisites

*   **Go**: Version 1.25 or later.
*   **PostgreSQL**: A running PostgreSQL database instance.
*   **Migrate CLI**: For running database migrations (`golang-migrate`).
*   **Make**: For running build and utility commands.

## Quickstart

### 1. Configuration

Create a `.envrc` file (or set environment variables) with your configuration. You can use the provided `.envrc` as a template.

```bash
export ANIMEVERSE_DB_DSN='postgres://user:password@localhost/animeverse?sslmode=disable'
```

### 2. Database Setup

Run the database migrations to set up the schema:

```bash
make db/migrations/up
```

### 3. Running the Server

Start the API server:

```bash
make run/api
```

The server will start on port `4000` (default) in `development` mode.

## API Documentation

### Endpoints

**Healthcheck**
*   `GET /v1/healthcheck`: Check system status.

**Authentication**
*   `POST /v1/users`: Register a new user.
*   `PUT /v1/users/activated`: Activate a user account.
*   `POST /v1/tokens/authentication`: Generate an authentication token.

**Animes**
*   `GET /v1/animes`: List animes (supports pagination, filtering).
*   `POST /v1/animes`: Create a new anime (Admin only).
*   `GET /v1/animes/:id`: Get anime details.
*   `PATCH /v1/animes/:id`: Update anime details (Admin only).
*   `DELETE /v1/animes/:id`: Delete an anime (Admin only).

**User Anime List**
*   `GET /v1/user_anime_list`: Get the authenticated user's anime list.
*   `POST /v1/user_anime_list`: Add an anime to the user's list.
*   `GET /v1/user_anime_list/:id`: Get a specific list entry.
*   `PATCH /v1/user_anime_list/:id`: Update a list entry (e.g., status, score).
*   `DELETE /v1/user_anime_list/:id`: Remove an anime from the list.

### Third-Party Integration: Jikan API

This API uses the [Jikan API](https://jikan.moe/) to fetch anime metadata when adding new entries. This ensures rich data availability without manual entry.

## Development

*   **Create Migration**: `make db/migrations/new name=migration_name`
*   **Connect to DB**: `make db/psql`

