# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a casino transactions system built in Go that processes bet and win events asynchronously through Kafka, persists them in PostgreSQL, and provides a REST API for querying transactions with filtering and pagination.

## Core Architecture

The system follows Clean Architecture principles with clear separation of concerns:

- **Domain Layer**: `internal/module/transaction/core/` - Contains business logic, Transaction entities, and repository interfaces
- **Application Layer**: `internal/module/transaction/use_case/` - Contains ProcessTransaction use case
- **Infrastructure Layer**: 
  - HTTP API: `internal/module/transaction/api/http/` - Fiber-based REST controllers
  - Kafka Consumer: `internal/module/transaction/consumer/` - Message processing handlers
  - Repository: `internal/module/transaction/repository/` - PostgreSQL implementations
- **Shared Components**: `pkg/shared/` - Reusable HTTP utilities, Kafka abstractions

## Development Commands

### Testing
```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

### Code Quality
```bash
# Lint code (uses golangci-lint with extensive rule set)
golangci-lint run

# The project uses revive linter with config in revive.toml
```

### Running Locally
```bash
# Start all services (API, consumer, PostgreSQL, Kafka)
docker compose up --build -d

# Generate test transactions
node scripts/transactions-kafka-producer.js
```

## Key Technical Details

### Kafka Integration
- Uses dynamic handler registration via `RegisterHandler(key, handler)`
- Handlers are registered in the consumer for different message types
- Kafka client abstraction in `pkg/shared/message_broker/kafka.go`

### Database Schema
- PostgreSQL with TIMESTAMPTZ for consistent datetime handling
- Schema initialization via `init.sql`
- Repository pattern for data access abstraction

### HTTP Framework
- Uses Fiber (Go's Express.js equivalent) with custom abstractions
- Security headers automatically applied
- CORS enabled
- Custom query parameter handling in `pkg/shared/http/`

### Testing Strategy
- Unit tests with `testify/mock` for mocking
- Integration tests use `testcontainers-go` for PostgreSQL and Kafka
- Test utilities in `pkg/test_utils/`
- Current coverage: ~89.5%

## Transaction Types
- `bet` - User placing a bet
- `win` - User winning a bet

## Sample Test Data
When using the Kafka producer script, use these user IDs for testing:
- `573a37e7-832a-4ecd-9691-41ff29afb955`
- `912457bc-7eed-4170-aba5-8a13c35a9d8a`
- `62d54d96-88f4-4111-8564-c043d710bdcd`

## API Endpoints
- `GET /transactions` - Query transactions with filters (user_id, transaction_type, limit, offset)
- `GET /health` - Health check endpoint

## Configuration
- Environment-based configuration using `envconfig`
- Configuration struct in `pkg/config/config.go`
- Docker services use environment variables for database and Kafka connections