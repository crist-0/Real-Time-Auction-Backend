# Real-Time Auction Backend

A real-time auction backend built with Go, PostgreSQL, and eventually Redis. The project focuses on financial correctness, transactional integrity, concurrency control, and real-time bidding.

The system is being developed incrementally, starting with wallet fundamentals and progressing toward a concurrent, real-time auction platform.

## Goals

* Build a reliable wallet and ledger system.
* Maintain correct available and locked balances.
* Implement atomic deposits, transfers, and reservations.
* Understand PostgreSQL transactions, row-level locks, isolation, and deadlocks.
* Support concurrent bidding without compromising financial correctness.
* Introduce real-time auction updates as the system evolves.

## Tech Stack

| Technology    | Purpose                                                             |
| ------------- | ------------------------------------------------------------------- |
| Go            | Backend application                                                 |
| chi           | HTTP routing                                                        |
| PostgreSQL    | Primary database and source of truth                                |
| pgx / pgxpool | PostgreSQL driver and connection pooling                            |
| Redis         | Planned: real-time event distribution and supporting infrastructure |
| Raw SQL       | Database queries and transaction management                         |

Redis is not currently required. PostgreSQL remains the authoritative source for wallet balances, reservations, and auction state.

## Current Status

The project is under active development.

* [x] Initialize the Go module.
* [x] Set up a basic HTTP server using chi.
* [x] Implement `GET /health`.
* [x] Initialize the Git repository.
* [ ] Configure PostgreSQL connection pooling.
* [ ] Add a database-backed wallet endpoint.
* [ ] Implement deposits and transfers.
* [ ] Implement wallet reservations.
* [ ] Implement release and capture operations.
* [ ] Add concurrency and transaction tests.
* [ ] Implement auctions and bidding.
* [ ] Add real-time event delivery.

Features marked incomplete are planned, not yet implemented.

## Project Structure

The structure will evolve as features are added.

```text
auction-backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── http/
│   └── wallet/
├── migrations/
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

* `cmd/server`: application entry point.
* `internal/http`: HTTP routing and handlers.
* `internal/wallet`: wallet business logic and database access.
* `migrations`: versioned database schema changes.

Some directories may not exist yet; they will be added as development progresses.

## Prerequisites

Install the following:

* Go — a version compatible with the project's `go.mod`.
* PostgreSQL — for persistent storage.
* Git — for version control.

Redis is not required for the current stage.

## Getting Started

### 1. Clone the repository

```bash
git clone <repository-url>
cd auction-backend
```

If you're working from your existing local repository, skip this step.

### 2. Install Go dependencies

```bash
go mod download
```

### 3. Configure the database

Create or select a local PostgreSQL database for development.

Set the `DATABASE_URL` environment variable to your PostgreSQL connection string.

**Windows Command Prompt:**

```bat
set DATABASE_URL=postgres://USERNAME:PASSWORD@localhost:5432/layer1
```

**PowerShell:**

```powershell
$env:DATABASE_URL = "postgres://USERNAME:PASSWORD@localhost:5432/layer1"
```

Replace the placeholders with your local credentials and database name.

Do not commit real credentials to Git. The environment variable will be used by the application when PostgreSQL connectivity is implemented.

### 4. Run the application

From the project root:

```bash
go run ./cmd/server
```

The HTTP server currently listens on port `3000`.

### 5. Check the health endpoint

```bash
curl -i http://localhost:3000/health
```

Expected response:

```http
HTTP/1.1 200 OK
Content-Type: application/json
```

```json
{
  "status": "Ok"
}
```

The health endpoint currently checks only that the HTTP application is responding. A separate database readiness check will be introduced later.

## Planned API

The following endpoints are planned and are not all implemented yet.

| Method | Endpoint                            | Purpose                        |
| ------ | ----------------------------------- | ------------------------------ |
| GET    | `/health`                           | Basic HTTP health check        |
| GET    | `/users/{userID}/wallet`            | Retrieve wallet balances       |
| POST   | `/users/{userID}/wallet/deposit`    | Deposit funds                  |
| POST   | `/transfers`                        | Transfer funds between wallets |
| POST   | `/wallet/reservations`              | Reserve available funds        |
| POST   | `/wallet/reservations/{id}/release` | Release reserved funds         |
| POST   | `/wallet/reservations/{id}/capture` | Capture reserved funds         |

The final API design may change as requirements become clearer.

## Financial Correctness

PostgreSQL is the source of truth for financial state.

The wallet model will distinguish:

* **Available balance:** funds that can be spent or reserved.
* **Locked balance:** funds reserved for an operation, such as an auction bid.

An important invariant is:

```text
available_balance >= 0
locked_balance >= 0
```

For a wallet without other balance adjustments:

```text
available_balance + locked_balance = total wallet funds
```

Transfers and reservation operations must use database transactions. Related balance changes and ledger entries must commit atomically or roll back together.

Concurrent operations will be protected using appropriate database constraints, atomic updates, and row-level locking. Lock acquisition order will be consistent when multiple wallets are involved.

## Development Principles

* Prefer raw SQL while learning database behavior.
* Keep business logic separate from HTTP handlers.
* Keep transaction boundaries explicit.
* Validate inputs at API boundaries and enforce critical invariants in PostgreSQL.
* Use parameterized SQL queries.
* Test concurrent operations, insufficient funds, rollback behavior, and failure scenarios.
* Keep secrets and real financial data out of version control.
* Introduce additional infrastructure only when a concrete requirement justifies it.

## Roadmap

### Stage 1 — Wallet Fundamentals

Create the wallet API, connect PostgreSQL, and implement deposits and transfers with ledger entries.

### Stage 2 — Reservations

Implement reserve, release, and capture operations with transactional balance updates.

### Stage 3 — Concurrency

Test simultaneous transfers and reservations, investigate deadlocks, and verify database invariants under load.

### Stage 4 — Auction Engine

Add auctions, bids, bid validation, and reservation management.

### Stage 5 — Real-Time Updates

Introduce Redis and WebSocket or Server-Sent Events (SSE) for distributing auction updates to connected clients.

## License

Not yet specified.
