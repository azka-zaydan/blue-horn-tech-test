# Blue Horn Tech Test

This repository contains a full-stack application with a Go backend and a React (TypeScript) frontend, orchestrated via Docker Compose. It uses PostgreSQL as the database.

## Table of Contents

- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Environment Variables](#environment-variables)
  - [Running with Docker Compose](#running-with-docker-compose)
- [Backend](#backend)
  - [Development Mode](#development-mode)
  - [Unit Testing](#unit-testing)
- [Frontend](#frontend)
- [Database](#database)
- [Development](#development)
- [Useful Links](#useful-links)

---

## Project Structure

```
blue-horn-tech-test/
├── backend/         # Go backend service
│   ├── Dockerfile
│   ├── .env.example
│   └── ...
├── frontend/        # React frontend (TypeScript)
│   ├── Dockerfile
│   ├── .env.example
│   └── ...
├── docker-compose.yaml
├── .env.example     # Root env for docker-compose (frontend)
└── ...
```

---

## Getting Started

### Prerequisites

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/)
- (For local dev outside Docker) [Go](https://go.dev/) and [Node.js](https://nodejs.org/)

### Environment Variables

Copy the example environment files and adjust as needed:

```bash
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env
cp .env.example .env
```

#### Root `.env.example`

```dotenv
VITE_API_URL=http://localhost:8080/api
```

#### Backend `.env.example`

```dotenv
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=evvlogger
```

#### Frontend `.env.example`

```dotenv
VITE_API_URL=http://localhost:8080/api
VITE_API_HOST=localhost
VITE_API_PORT=8080
VITE_API_PROTOCOL=http
VITE_API_BASE_PATH=/api
VITE_APP_TITLE=My React App
VITE_ENV=development
```

---

## Running with Docker Compose

To start the full stack (database, backend, frontend):

```bash
docker-compose up --build
```

- Frontend: http://localhost:3000
- Backend API: http://localhost:8080/api
- PostgreSQL: localhost:5432 (default user/pass: postgres/postgres)

To stop and remove containers:

```bash
docker-compose down
```

---

## Backend

### Development Mode

Use the Makefile for fast local development (auto-reloads on file changes):

```bash
cd backend
make dev
```

> This uses [air](https://github.com/cosmtrek/air) for hot-reloading.

### Unit Testing

Unit tests are implemented for both repository and service layers, using mocks for the database and dependencies.

You’ll find unit tests in files like:

- `backend/src/domains/schedule/repository/schedule_repo_test.go`
- `backend/src/domains/schedule/service/schedule_svc_test.go`
- `backend/src/domains/task/repository/task_repo_test.go`
- `backend/src/domains/task/service/task_svc_test.go`

Typical test structure uses [sqlmock](https://github.com/DATA-DOG/go-sqlmock), [gomock](https://github.com/golang/mock), and [testify](https://github.com/stretchr/testify):

```go
func TestGetSchedules(t *testing.T) {
	initMocks(t)
	// ...mock SQL expectations...
	// ...call repo.GetSchedules() and assert results...
}
```

To run all backend unit tests:

```bash
cd backend
go test ./...
```

_You can add this to your Makefile for convenience:_
```makefile
test:
	go test ./...
```

---

## Frontend

- Written in **React** with **TypeScript**
- Located in `/frontend`
- Configured via `/frontend/.env`
- Talks to backend using the `VITE_API_URL` variable

To run locally (without Docker):

```bash
cd frontend
cp .env.example .env
yarn install
yarn dev
```

#### Frontend Testing

The frontend uses [Vitest](https://vitest.dev/) for testing:

```bash
cd frontend
yarn run test
```

See [frontend/README.md](frontend/README.md) for more details.

---

## Database

- Uses **PostgreSQL** (see `docker-compose.yaml`)
- Init SQL can be mounted via `backend/migration/init.sql`
- Default DB config:
  - Host: `localhost`
  - Port: `5432`
  - User: `postgres`
  - Password: `postgres`
  - Database: `evvlogger`

---

## Useful Links

- [Backend Directory](https://github.com/azka-zaydan/blue-horn-tech-test/tree/main/backend)
- [Frontend Directory](https://github.com/azka-zaydan/blue-horn-tech-test/tree/main/frontend)

---

## License

MIT
