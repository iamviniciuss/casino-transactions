# 🎰 Casino Transactions System

This project is a backend service designed to manage and process casino transaction data. It receives bet and win events asynchronously through Kafka, persists them in a relational database, and provides a RESTful API to query transactions with filtering and pagination support.

---

## 📦 Stack & Technologies

- **Golang** – Main programming language
- **Kafka** – Event/message processing
- **PostgreSQL** – Data persistence
- **Fiber** – HTTP API
- **Testcontainers** – Integration tests
- **Docker & Docker Compose** – Infrastructure and services
- **uuid, testify, etc.** – Supporting tools

---

## Architecture & Design Decisions

### 1. **Clean Architecture Principles**
- Domain (`core`) is decoupled from infrastructure (`db`, `api`, `kafka`)
- Each adapter (Kafka, HTTP, DB) communicates via defined interfaces

### 2. **Kafka Adapter**
- Messages are processed via dynamic handlers
- New use cases can be registered using `RegisterHandler(key, handler)`
- Encourages separation of concerns and scalability

### 3. **Repository Pattern**
- Abstracts persistence logic
- Enables easy mocking in unit tests
- Allows flexible database migration

### 4. **Test Strategy**
- Unit tests with `testify/mock`
- Integration tests with `testcontainers-go` for PostgreSQL and Kafka
- Coverage support with `go test -coverprofile`

### 5. **Database**
- Chosen: PostgreSQL for reliability and SQL support
- `TIMESTAMPTZ` used for consistency in date/time
- Schema versioned via `init.sql` file

---

## 🧪 Running Locally

### 1. Build the environment

```bash
docker compose up --build -d
```

### 2. Run the Kafka producer (optional)

```bash
node scripts/kafka-producer.js
```

---

## 🔄 Note on Sample Data

If you're using the provided Kafka script (`scripts/transactions-kafka-producer.js`) to generate random transactions, the data is simulated using a fixed set of user IDs.

You can use the following `user_id` values when querying the API:

```
573a37e7-832a-4ecd-9691-41ff29afb955
912457bc-7eed-4170-aba5-8a13c35a9d8a
62d54d96-88f4-4111-8564-c043d710bdcd
```


---

## 🔎 API Example

### `GET /transactions`

Query Parameters:
- `user_id` (required)
- `transaction_type` = `bet` or `win`
- `limit` (default 20)
- `offset` (default 0)

**Response:**

```json
{
  "items": [ ... ],
  "total": 45,
  "limit": 10,
  "offset": 0
}
```

### cURL
```bash
curl -X GET "http://127.0.0.1:9095/transactions?user_id=62d54d96-88f4-4111-8564-c043d710bdcd&limit=10&offset=0"
```

---

## 🧪 Testing

### Run all tests:

```bash
make test
make test-coverage

Current coverage: 89.5%
```

---

## 📁 Project Structure

```
/cmd
  /api         # Entrypoint: starts the HTTP server
  /consumer    # Entrypoint: starts the Kafka consumer

/internal
  /api
    /http  # HTTP controllers (parse input, validate, call use cases)
    /router      # HTTP route definitions and bindings
  /consumer      # Kafka consumer loop and handler registration
  /core          # Domain logic and repository interfaces (Transaction, Filter, Errors)
  /use_case      # Application use cases (e.g., ProcessTransaction)
  /repository    # PostgreSQL-based repository implementations

/pkg
  /config        # Environment and configuration loader
  shared          # Shared utilities (e.g., logging, error handling)
    /http        # Fiber setup, error handling, query parameter abstraction
  /testutils     # Testcontainers setup for Kafka and PostgreSQL
  /mocks         # Unit test mocks (e.g., repository, use cases)
```




---

## Author
Vinicius ☕