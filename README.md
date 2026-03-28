# Staxa Demo: Echo + PostgreSQL — Contact Manager

A fullstack contact manager application built with [Echo](https://echo.labstack.com/) (Go) and PostgreSQL. This is a demo/starter template designed for deployment on the [Staxa](https://staxa.dev) PaaS platform.

## Tech Stack

- **Framework**: Echo v4 (Go 1.22)
- **Database**: PostgreSQL 16
- **Template Engine**: Go html/template
- **Database Driver**: pgx v5

## Features

- Full CRUD for contacts (name, email, phone, notes)
- Server-rendered HTML pages + JSON API
- Health check endpoint (`GET /healthz`) with DB connectivity verification
- Auto-migration on startup
- 3 sample contacts seeded on first run

## Local Development

### Prerequisites

- Go 1.22+
- PostgreSQL 16

### Setup

```bash
# Create the database
createdb echo_postgres

# Set environment variable
export DATABASE_URL="postgresql://postgres:postgres@localhost:5432/echo_postgres"

# Install dependencies
go mod tidy

# Run the application
go run .
```

The app will be available at http://localhost:8080.

## Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DATABASE_URL` | Yes | — | PostgreSQL connection string |
| `PORT` | No | `8080` | Port the server listens on |

## Docker

```bash
# Build the image
docker build -t staxa-demo-echo-postgres .

# Run the container
docker run -e DATABASE_URL="postgresql://postgres:postgres@host.docker.internal:5432/echo_postgres" -p 8080:8080 staxa-demo-echo-postgres
```

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/contacts` | List all contacts |
| `POST` | `/api/contacts` | Create a contact |
| `GET` | `/api/contacts/:id` | Get one contact |
| `PUT` | `/api/contacts/:id` | Update a contact |
| `DELETE` | `/api/contacts/:id` | Delete a contact |
| `GET` | `/healthz` | Health check |

## Deployment

This app is designed for deployment on [Staxa](https://staxa.dev). Push to GitHub and deploy using the Echo + PostgreSQL template.
