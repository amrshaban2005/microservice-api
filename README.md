# Banking Microservice API

A Go REST API for banking operations: customer lookup, account creation, and deposits/withdrawals. The service uses [Gin](https://github.com/gin-gonic/gin) for HTTP, PostgreSQL for persistence, and delegates authorization to a separate auth service.

## Features

- List and retrieve customers (optional filter by `active` / `inactive` status)
- Open new savings or checking accounts (minimum opening balance: 5000)
- Record deposits and withdrawals with balance checks on withdraw
- JWT-based authorization via an external auth service

## Tech stack

| Layer        | Technology                                      |
| ------------ | ----------------------------------------------- |
| Language     | Go 1.26+                                        |
| HTTP         | Gin                                             |
| Database     | PostgreSQL 15 (`pgx/v5`)                        |
| Shared lib   | [`banking-lib`](https://github.com/amrshaban2005/banking-lib) (errors, logging) |

## Architecture

```
HTTP (Gin)  →  app handlers  →  service  →  domain repositories  →  PostgreSQL
                    ↓
            Auth middleware  →  remote auth service
```

| Package   | Role                                              |
| --------- | ------------------------------------------------- |
| `app`     | Routes, handlers, auth middleware                 |
| `service` | Business rules and validation                     |
| `domain`  | Entities and database access                      |
| `dto`     | Request/response shapes                           |

## Prerequisites

- Go 1.26 or later
- Docker (for PostgreSQL via Compose)
- A running **auth service** that exposes `GET /auth/verify` (see [Authorization](#authorization))

## Quick start

### 1. Start the database

From the project root:

```bash
docker compose -f resources/docker/docker-compose.yml up -d
```

This starts PostgreSQL and applies the schema from `resources/docker/initdb/banking.sql`. Configure credentials and ports in `resources/docker/docker-compose.yml` to match your environment.

### 2. Configure environment

Set the following variables before starting the API (for example via your shell, a `.env` file, or `start.sh`):

| Variable         | Required | Description                    |
| ---------------- | -------- | ------------------------------ |
| `SERVER_ADDRESS` | yes      | HTTP bind address              |
| `SERVER_PORT`    | yes      | HTTP port                      |
| `DB_USER`        | yes      | PostgreSQL user                |
| `DB_PASSWD`      | no*      | PostgreSQL password            |
| `DB_ADDR`        | yes      | PostgreSQL host                |
| `DB_PORT`        | yes      | PostgreSQL port                |
| `DB_NAME`        | yes      | Database name                  |

\*`DB_PASSWD` is not checked at startup but is needed for a successful DB connection.

### 3. Run the API

Using the provided script (after setting the required variables in `start.sh` or your environment):

```bash
chmod +x start.sh
./start.sh
```

Or run directly once the variables above are set:

```bash
go run main.go
```

## Authorization

Every protected route expects an `Authorization` header:

```http
Authorization: Bearer <JWT>
```

The middleware calls the auth service verify endpoint with the token, permission name, and relevant route parameters (for example `customer_id`, `account_id`).

Permissions used by this API:

| Permission            | Route                                      |
| --------------------- | ------------------------------------------ |
| `customers:read_all`  | `GET /customers`                           |
| `customers:read_one`  | `GET /customers/:id`                       |
| `accounts:create`     | `POST /customers/:id/account`              |
| `transactions:create` | `POST /customers/:id/account/:account_id`  |

Obtain a JWT from your auth service before calling protected routes.

## API reference

### Customers

**List customers**

```http
GET /customers
GET /customers?status=active
GET /customers?status=inactive
```

**Get one customer**

```http
GET /customers/:id
```

### Accounts

**Open account**

```http
POST /customers/:id/account
Content-Type: application/json

{
  "account_type": "saving",
  "amount": 5000
}
```

- `account_type`: `saving` or `checking`
- `amount`: must be ≥ 5000

### Transactions

**Deposit or withdraw**

```http
POST /customers/:id/account/:account_id
Content-Type: application/json

{
  "transaction_type": "deposit",
  "amount": 100.50
}
```

- `transaction_type`: `deposit` or `withdraw`
- `amount`: must be > 0
- Withdrawals fail if the account balance is insufficient

## Project layout

```
microservice-api/
├── app/                 # HTTP layer (routes, handlers, middleware)
├── domain/              # Models and repository implementations
├── dto/                 # API request/response types
├── service/             # Application services
├── resources/docker/    # Docker Compose and DB init scripts
├── main.go              # Entry point
├── start.sh             # Local run helper with env vars
└── go.mod
```

## Development

Build:

```bash
go build -o microservice-api .
```

Run tests (when present):

```bash
go test ./...
```

Stop the database:

```bash
docker compose -f resources/docker/docker-compose.yml down
```
