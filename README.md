# Book Store API Service

A comprehensive REST API for managing an online bookstore built with Go. This service handles book catalogs, author management, customer registration/authentication, order processing, and automated sales reporting every 24 hours.

---

## Setup

```bash
cd "/home/moad/Desktop/Projects/Book Store"
go mod download

JWT_KEY=$(openssl rand -base64 64)
cat <<EOF > .env
JWT_SECRET=$JWT_KEY
SERVER_PORT=8080
DB_PATH=database.json
REPORT_OUTPUT_DIR=output-reports/
REPORT_INTERVAL=24h
EOF
```

Run server:

```bash
go run main.go
```

---

## Test Flow

Follow this order exactly. Every step depends on the previous state.

---

### 1) Create Author

```bash
curl -X POST http://localhost:8080/authors/ \
-H "Content-Type: application/json" \
-d '{
  "first_name": "George",
  "last_name": "Orwell",
  "bio": "English novelist"
}'
```

---

### 2) Create Book

```bash
curl -X POST http://localhost:8080/books/ \
-H "Content-Type: application/json" \
-d '{
  "title": "1984",
  "author": { "id": 0 },
  "genres": ["dystopia", "political"],
  "price": 15.99,
  "published_at": "1949-06-08T00:00:00Z",
  "stock": 10
}'
```

---

### 3) Search Books (Get All Books)

```bash
curl http://localhost:8080/books/
```

---

### 4) Get Book By ID

```bash
curl http://localhost:8080/books/0
```

---

### 5) Update Book

```bash
curl -X PUT http://localhost:8080/books/0 \
-H "Content-Type: application/json" \
-d '{
  "title": "1984 - Updated Edition",
  "author": { "id": 0 },
  "genres": ["dystopia", "political", "fiction"],
  "price": 17.99,
  "published_at": "1949-06-08T00:00:00Z",
  "stock": 15
}'
```

---

### 6) Search Books with Filters

Search by author:
```bash
curl "http://localhost:8080/books/?author=george"
```

Search by genre:
```bash
curl "http://localhost:8080/books/?genre=dystopia"
```

Search by price range:
```bash
curl "http://localhost:8080/books/?min_price=10&max_price=40"
```

Search by title:
```bash
curl "http://localhost:8080/books/?title=1984"
```

Sort books (by price, title, etc.):
```bash
curl "http://localhost:8080/books/?sort_by=price&sort_order=asc"
```

---

### 7) Delete Book

```bash
curl -X DELETE http://localhost:8080/books/0
```

---

### 8) Register Customer

```bash
curl -X POST http://localhost:8080/customers/register \
-H "Content-Type: application/json" \
-d '{
  "name": "Alice",
  "email": "alice@example.com",
  "password": "password123",
  "address": {
    "street": "123 Maple St",
    "city": "Springfield",
    "state": "IL",
    "postal_code": "62704",
    "country": "USA"
  }
}'
```

Response includes JWT token for authentication.

---

### 9) Login + Export Token

```bash
export TOKEN=$(curl -s -X POST http://localhost:8080/customers/login \
-H "Content-Type: application/json" \
-d '{"email": "alice@example.com", "password": "password123"}' | jq -r .token)
```

---

### 10) Customer Endpoints

#### Get All Customers

```bash
curl http://localhost:8080/customers
```

#### Get Customer by ID

**Authentication required**

```bash
curl -H "Authorization: Customer $TOKEN" http://localhost:8080/customers/0
```

#### Update Customer

**Authentication required**

```bash
curl -X PUT http://localhost:8080/customers/0 \
-H "Authorization: Customer $TOKEN" \
-H "Content-Type: application/json" \
-d '{
  "name": "Alice Walker",
  "email": "alice.updated@example.com",
  "address": {
    "street": "456 Oak Ave",
    "city": "Springfield",
    "state": "IL",
    "postal_code": "62704",
    "country": "USA"
  }
}'
```

#### Delete Customer

**Authentication required**

```bash
curl -X DELETE http://localhost:8080/customers/0 \
-H "Authorization: Customer $TOKEN"
```

---

### 11) Order Endpoints

All order endpoints require authentication.

#### Create Order

```bash
curl -X POST http://localhost:8080/orders/ \
-H "Authorization: Customer $TOKEN" \
-H "Content-Type: application/json" \
-d '{
  "items": [
    { "book_id": 0, "quantity": 2 }
  ]
}'
```

#### Get Order by ID

```bash
curl -H "Authorization: Customer $TOKEN" http://localhost:8080/orders/0
```

#### List All Orders

```bash
curl -H "Authorization: Customer $TOKEN" http://localhost:8080/orders
```

#### Filter Orders by Status

```bash
curl -H "Authorization: Customer $TOKEN" "http://localhost:8080/orders?status=created"
```

Other status values: `completed`, `cancelled`

#### Get Orders in Time Range

```bash
curl -H "Authorization: Customer $TOKEN" \
"http://localhost:8080/orders?start_date=2026-01-01T00:00:00Z&end_date=2026-12-31T23:59:59Z"
```

#### Update Order Status

Complete an order:
```bash
curl -X PUT "http://localhost:8080/orders/0?status=completed" \
-H "Authorization: Customer $TOKEN"
```

Cancel an order:
```bash
curl -X PUT "http://localhost:8080/orders/0?status=cancelled" \
-H "Authorization: Customer $TOKEN"
```

---

### 12) Sales Reports

#### List All Reports

```bash
curl http://localhost:8080/reports/sales
```

#### Get Latest Report

```bash
curl http://localhost:8080/reports/sales/latest
```

#### Get Reports Since N Days Ago

```bash
curl "http://localhost:8080/reports/sales?since=7"
```

#### Generate New Report Manually

```bash
curl -X POST http://localhost:8080/reports/sales/generate
```

---

### 13) Metrics System

#### Total Customers

```bash
curl "http://localhost:8080/metrics?total_customers=1"
```

#### Total Books

```bash
curl "http://localhost:8080/metrics?total_books=1"
```

#### Total Authors

```bash
curl "http://localhost:8080/metrics?total_authors=1"
```

#### Out of Stock Books

```bash
curl "http://localhost:8080/metrics?out_of_stock_books=1"
```

#### Books Per Author

```bash
curl "http://localhost:8080/metrics?books_per_author=1"
```

#### Books Per Genre

```bash
curl "http://localhost:8080/metrics?books_per_genre=1&genre=dystopia"
```

#### API Hits Counter

```bash
curl http://localhost:8080/metrics/hits
```

---

## Authentication Model

Custom JWT-based authentication scheme:

```
Authorization: Customer <JWT_TOKEN>
```

---

## Shutdown

Press `CTRL+C` for graceful shutdown.
