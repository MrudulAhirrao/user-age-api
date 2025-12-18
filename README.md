# User Age API

A high-performance, production-ready RESTful API built with **Go** and **Fiber**. This service manages user profiles and demonstrates dynamic business logic by calculating a user's age in real-time based on their Date of Birth.

Designed with **Clean Architecture** principles and **Type-Safe SQL**.

---

## Tech Stack

| Component | Technology | Reasoning |
| :--- | :--- | :--- |
| **Language** | [Go (Golang)](https://go.dev/) | 1.23+ |
| **Framework** | [GoFiber v2](https://gofiber.io/) | Express-style, high-performance HTTP framework. |
| **Database** | [PostgreSQL](https://www.postgresql.org/) | Robust, relational database storage. |
| **Data Access** | [SQLC](https://sqlc.dev/) | Generates type-safe Go code from raw SQL (using `pgx/v5`). |
| **Logging** | [Uber Zap](https://github.com/uber-go/zap) | Structured, leveled logging for observability. |
| **Validation** | [Validator v10](https://github.com/go-playground/validator) | Struct-based input validation. |
| **Container** | [Docker](https://www.docker.com/) | Containerization for consistent environments. |

---

## Project Structure

The project follows the Standard Go Project Layout to ensure separation of concerns:

```bash
├── cmd/server/       # Application entry point
├── db/
│   ├── migrations/   # SQL Schema definitions
│   ├── query.sql     # Raw SQL queries for SQLC
│   └── sqlc/         # Auto-generated Go database code (DO NOT EDIT)
├── internal/
│   ├── handler/      # HTTP Layer: Request parsing & Validation
│   ├── service/      # Business Logic Layer: Age Calculation lives here
│   ├── models/       # API Data structures & Unit Tests
│   └── logger/       # Logger configuration
├── Dockerfile        # Multi-stage Docker build
└── reasoning.md      # Documentation of design choices
```
---

## Getting Started
Follow these steps to run the application locally.

### Prerequisites
  - Go 1.23+ installed
  - Docker installed and running
### 1. Start the Database
We use Docker to spin up a temporary PostgreSQL instance
```bash
docker run --name db \
  -e POSTGRES_USER=user \
  -e POSTGRES_PASSWORD=secret \
  -e POSTGRES_DB=userdb \
  -p 5432:5432 \
  -d postgres:alpine
```
### 2. Initialize the Schema
Create the users table in the running container:
```bash
docker exec -i db psql -U user -d userdb -c "CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT NOT NULL, dob DATE NOT NULL);"
```
### 3. Run the Application

Start the Go server:

```Bash

go run cmd/server/main.go
```
The server will start on port 3000.

## API Documentation
### 1. Create User
POST /users

Request Body:
```bash
JSON

{
  "name": "Alice",
  "dob": "1995-05-20"
}
```
Response:
```bash
JSON

{
  "id": 1,
  "name": "Alice",
  "dob": "1995-05-20",
  "age": 29
}
```
### 2. Get User (Calculates Age)
GET /users/:id

Returns the user details with the dynamically calculated age.
Response:
```bash
JSON

{
  "id": 1,
  "name": "Alice",
  "dob": "1995-05-20",
  "age": 29
}
```
### 3. List All Users
GET /users

Returns a list of all users with their current ages.

### 4. Update User
PUT /users/:id

Request Body:
```bash
JSON

{
  "name": "Alice Updated",
  "dob": "1996-01-01"
}
```
### 5. Delete User
DELETE /users/:id

Returns HTTP 204 No Content.

## Testing
### Unit Tests
The core logic (Age Calculation) is covered by unit tests.

```Bash
go test -v ./internal/models
Docker Build
```
To verify the application builds correctly in a containerized environment:
```bash
docker build -t user-age-api .
```
