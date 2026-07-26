# Subscription Aggregation Service

REST service for storing and aggregating user online subscriptions.

The service provides CRUDL operations for subscription records and calculates total subscription cost for a selected period with filtering by user ID and subscription name.

## Task Scope And Expectations

The application implements the requirements from the test assignment:

- CRUDL operations for subscription records:
  - create subscription;
  - get subscription by ID;
  - update subscription;
  - delete subscription;
  - list subscriptions by user ID.
- Total cost calculation for subscriptions over a selected period.
- Filtering total cost by:
  - `user_id`;
  - `service_name`.
- PostgreSQL as persistent storage.
- Database migrations for schema initialization.
- Configuration via environment variables.
- Structured logging.
- Swagger/OpenAPI documentation.
- Docker Compose setup for local launch.

A subscription record contains:

- service name;
- monthly price in rubles;
- user ID in UUID format;
- subscription start date in `MM-YYYY` format;
- optional subscription end date in `MM-YYYY` format.

User existence is not checked. User management is outside of this service scope.

Subscription price is stored as an integer number of rubles.

## Tech Stack

- Go
- Echo
- PostgreSQL
- pgx
- sqlc
- goose migrations
- oapi-codegen
- Docker Compose
- slog structured logging

## Project Structure

```text
cmd/app                         Application entrypoint
api/openapi.yaml                 OpenAPI specification
internal/app                     Application wiring and lifecycle
internal/config                  Environment configuration
internal/db                      PostgreSQL connection and migrations
internal/db/sqlc                 sqlc queries and generated storage code
internal/models                  Internal domain models
internal/repository             Database access layer
internal/service                Business logic layer
internal/transport/http          HTTP handlers and middleware
pkg/swagger-ui                   Embedded Swagger UI assets
```

## How To Run

### 1. Prepare Environment File

Create `.env` from the example file:

```bash
cp .env.example .env
```

Default values from `.env.example`:

```env
SERVER_PORT=8080

POSTGRES_USER=service
POSTGRES_PASSWORD=superstrongpassword
POSTGRES_DB=service
POSTGRES_MAX_CONNS=100
```

The Docker Compose file builds `POSTGRES_URL` from these values and passes it to the application container.

### 2. Start The Service

```bash
docker-compose up --build
```

### 3. Stop The Service

```bash
docker-compose down
```

## API Documentation

Swagger UI is available at:

```text
http://localhost:8080/api/v1/swagger/
```

OpenAPI JSON is available at:

```text
http://localhost:8080/api/v1/openapi.json
```

## API Examples

The examples below use `curl` and assume that the service is running at `localhost:8080`.

Set helper variables:

```bash
baseUrl="http://localhost:8080/api/v1"
userId="60601fee-2bf1-4721-ae6f-7636e79a0cba"
otherUserId="11111111-1111-1111-1111-111111111111"
```

### Create Subscription

```bash
curl -i -X POST "$baseUrl/subscriptions" \
  -H "Content-Type: application/json" \
  -d '{"service_name":"Yandex Plus","price":400,"user_id":"60601fee-2bf1-4721-ae6f-7636e79a0cba","start_date":"07-2025"}'
```

Example response:

```json
{
  "id": "f8db037d-e682-4d62-bc89-1ad9f9048b10",
  "price": 400,
  "service_name": "Yandex Plus",
  "start_date": "07-2025",
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba"
}
```

Save the returned `id` and use it in requests that work with a specific subscription:

```bash
subscriptionId="f8db037d-e682-4d62-bc89-1ad9f9048b10"
```

### Create Subscription With End Date

```bash
curl -i -X POST "$baseUrl/subscriptions" \
  -H "Content-Type: application/json" \
  -d '{"service_name":"Yandex Plus","price":400,"user_id":"60601fee-2bf1-4721-ae6f-7636e79a0cba","start_date":"07-2025","end_date":"12-2025"}'
```

Example response:

```json
{
  "end_date": "12-2025",
  "id": "337bd74c-b5c4-41be-9017-1a46776f8f0f",
  "price": 400,
  "service_name": "Yandex Plus",
  "start_date": "07-2025",
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba"
}
```

### List User Subscriptions

```bash
curl -i "$baseUrl/subscriptions?user_id=$userId"
```

Example response:

```json
[
  {
    "end_date": "12-2025",
    "id": "337bd74c-b5c4-41be-9017-1a46776f8f0f",
    "price": 400,
    "service_name": "Yandex Plus",
    "start_date": "07-2025",
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba"
  },
  {
    "id": "f8db037d-e682-4d62-bc89-1ad9f9048b10",
    "price": 400,
    "service_name": "Yandex Plus",
    "start_date": "07-2025",
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba"
  }
]
```

### Get Subscription By ID

```bash
curl -i "$baseUrl/subscriptions/$subscriptionId"
```

Example response:

```json
{
  "id": "f8db037d-e682-4d62-bc89-1ad9f9048b10",
  "price": 400,
  "service_name": "Yandex Plus",
  "start_date": "07-2025",
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba"
}
```

### Update Subscription

```bash
curl -i -X PUT "$baseUrl/subscriptions/$subscriptionId" \
  -H "Content-Type: application/json" \
  -d '{"service_name":"Yandex Plus","price":500,"user_id":"60601fee-2bf1-4721-ae6f-7636e79a0cba","start_date":"07-2025","end_date":"12-2025"}'
```

Example response:

```json
{
  "end_date": "12-2025",
  "id": "f8db037d-e682-4d62-bc89-1ad9f9048b10",
  "price": 500,
  "service_name": "Yandex Plus",
  "start_date": "07-2025",
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba"
}
```

### Delete Subscription

```bash
curl -i -X DELETE "$baseUrl/subscriptions/$subscriptionId"
```

Expected response status:

```text
HTTP/1.1 204 No Content
```

### Prepare Data For Total Cost Calculation

The following requests create a small dataset for checking period overlap and filters.

Subscription fully inside the checked service and user:

```bash
curl -i -X POST "$baseUrl/subscriptions" \
  -H "Content-Type: application/json" \
  -d '{"service_name":"Yandex Plus","price":500,"user_id":"60601fee-2bf1-4721-ae6f-7636e79a0cba","start_date":"07-2025","end_date":"12-2025"}'
```

Subscription started before the requested period:

```bash
curl -i -X POST "$baseUrl/subscriptions" \
  -H "Content-Type: application/json" \
  -d '{"service_name":"Yandex Plus","price":400,"user_id":"60601fee-2bf1-4721-ae6f-7636e79a0cba","start_date":"01-2025","end_date":"12-2025"}'
```

Subscription ended inside the requested period:

```bash
curl -i -X POST "$baseUrl/subscriptions" \
  -H "Content-Type: application/json" \
  -d '{"service_name":"Yandex Plus","price":300,"user_id":"60601fee-2bf1-4721-ae6f-7636e79a0cba","start_date":"07-2025","end_date":"08-2025"}'
```

Subscription without `end_date`:

```bash
curl -i -X POST "$baseUrl/subscriptions" \
  -H "Content-Type: application/json" \
  -d '{"service_name":"Yandex Plus","price":200,"user_id":"60601fee-2bf1-4721-ae6f-7636e79a0cba","start_date":"07-2025"}'
```

Another service for the same user:

```bash
curl -i -X POST "$baseUrl/subscriptions" \
  -H "Content-Type: application/json" \
  -d '{"service_name":"Netflix","price":700,"user_id":"60601fee-2bf1-4721-ae6f-7636e79a0cba","start_date":"07-2025","end_date":"12-2025"}'
```

Same service for another user:

```bash
curl -i -X POST "$baseUrl/subscriptions" \
  -H "Content-Type: application/json" \
  -d '{"service_name":"Yandex Plus","price":999,"user_id":"11111111-1111-1111-1111-111111111111","start_date":"07-2025","end_date":"12-2025"}'
```

### Calculate Total Cost

Calculate `Yandex Plus` cost for July-September 2025:

```bash
curl -i "$baseUrl/subscriptions/cost?user_id=$userId&service_name=Yandex%20Plus&start_period=07-2025&end_period=09-2025"
```

Expected response:

```json
{
  "total_cost": 3900
}
```

Calculation:

```text
500 * 3 months = 1500
400 * 3 months = 1200
300 * 2 months = 600
200 * 3 months = 600
Total = 3900
```

Calculate `Yandex Plus` cost for one month:

```bash
curl -i "$baseUrl/subscriptions/cost?user_id=$userId&service_name=Yandex%20Plus&start_period=07-2025&end_period=07-2025"
```

Expected response:

```json
{
  "total_cost": 1400
}
```

Calculate another service:

```bash
curl -i "$baseUrl/subscriptions/cost?user_id=$userId&service_name=Netflix&start_period=07-2025&end_period=09-2025"
```

Expected response:

```json
{
  "total_cost": 2100
}
```

### Validation Error Examples

Invalid date format:

```bash
curl -i "$baseUrl/subscriptions/cost?user_id=$userId&service_name=Yandex%20Plus&start_period=2025-07&end_period=09-2025"
```

Expected response status:

```text
HTTP/1.1 400 Bad Request
```

End period before start period:

```bash
curl -i "$baseUrl/subscriptions/cost?user_id=$userId&service_name=Yandex%20Plus&start_period=09-2025&end_period=07-2025"
```

Expected response status:

```text
HTTP/1.1 400 Bad Request
```

## Cost Calculation Rules

The service treats subscription price as a monthly price.

The total cost endpoint calculates only the months where the subscription period overlaps with the requested period.

Both period boundaries are inclusive.

Example:

```text
Subscription:
service_name = Yandex Plus
price        = 500
start_date   = 07-2025
end_date     = 12-2025

Request period:
start_period = 07-2025
end_period   = 09-2025

Result:
500 * 3 months = 1500
```

If `end_date` is omitted, the subscription is treated as active after `start_date`.

## Solution Notes

The project is split into several layers:

- HTTP transport layer receives requests and maps service errors to HTTP responses.
- Service layer validates input and converts API date strings into internal `time.Time` values.
- Repository layer works with PostgreSQL through sqlc-generated queries.
- Database schema is initialized with goose migrations.
- OpenAPI is the source of HTTP contract and generated request/response types.

The date format used by the API is:

```text
MM-YYYY
```

Example:

```text
07-2025
```

The service does not check whether a user exists. It only stores the provided `user_id`.

## Logging

The service uses structured JSON logs through Go `slog`.

Logged events include:

- application startup;
- graceful shutdown;
- HTTP requests;
- HTTP request errors;
- unexpected handler errors.

## Docker Compose Services

`compose.yml` starts:

- `service` - Go application container;
- `postgres` - PostgreSQL database container.

PostgreSQL data is stored in a named Docker volume:

```text
pg_data
```

To fully reset the database, run:

```bash
docker-compose down -v
```
