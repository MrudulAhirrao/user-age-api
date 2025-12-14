# Implementation Reasoning & Design Decisions

## Overview
This project is a RESTful API designed to manage user data and perform dynamic business logic (age calculation). While the scope is small, I approached it as a foundation for a production-grade system, focusing on **maintainability**, **type safety**, and **clean separation of concerns**.

Below is the reasoning behind my architectural and technical choices.

## 1. Architectural Pattern: Clean/Layered Architecture
I structured the application into three distinct layers using **Dependency Injection**:
* **Handler Layer (`/internal/handler`):** Deals strictly with HTTP transport (parsing JSON, validating inputs, returning status codes).
* **Service Layer (`/internal/service`):** Contains the business logic. This is where the Age Calculation happens. It knows nothing about HTTP or SQL.
* **Repository Layer (`/db/sqlc`):** Handles raw data access.

**Why?**
This structure ensures that the business logic is isolated. If we decided to switch from HTTP to gRPC later, or swap Postgres for MySQL, the Service layer (and the core logic) wouldn't need to change. It also makes unit testing the logic (as done in `models/user_test.go`) trivial.

## 2. Key Design Decisions

### Dynamic Age Calculation (Service Layer)
The requirement was to store `dob` but return `age`.
* **Decision:** I chose to calculate age dynamically on every read (GET request) within the Service layer, rather than storing it in the database.
* **Reasoning:** "Age" is a volatile field—it changes automatically as time passes. Storing it would create a data synchronization nightmare (requiring daily cron jobs to update records). Calculating it on the fly is computationally cheap (O(1)) and ensures the data is always 100% accurate relative to the current server time.

### Database Access: SQLC vs. ORM
* **Decision:** I used **SQLC** instead of an ORM like GORM.
* **Reasoning:** While ORMs are convenient for prototyping, they often hide performance bottlenecks and runtime SQL errors. SQLC gives me the best of both worlds:
    1.  **Type Safety:** It generates Go code from raw SQL. If my SQL query is wrong, the code won't compile.
    2.  **Performance:** It uses the efficient `pgx` driver without reflection overhead.
    3.  **Clarity:** Looking at `query.sql` makes it immediately obvious what the database is doing.

### Date Handling
* **Decision:** Used `pgtype.Date` from `pgx/v5`.
* **Reasoning:** Go's standard `time.Time` includes a timestamp component, whereas the database column is strictly a `DATE`. Using `pgtype.Date` handles the mapping explicitly, preventing timezone shifts where a user born on "1990-05-10" might accidentally shift to "1990-05-09" due to UTC conversion.

## 3. Tech Stack Choices

* **GoFiber:** Chosen for its performance and low overhead. Its API is very similar to Express.js, making the handlers concise and readable.
* **Uber Zap:** I used Zap for structured, leveled logging. In a real production environment, structured logs are essential for ingestion by tools like Datadog or ELK.
* **Docker:** I included a multi-stage Dockerfile to ensure the application is portable and easy to deploy without worrying about the host OS environment.

## 4. Summary
I prioritized **code readability** and **robustness**. The application handles edge cases (like users born in the future or today), validates inputs before they hit the database, and uses strictly typed SQL interactions to prevent runtime surprises.
