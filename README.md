# Book_Store_Service


##  Project Status (TODO)


###  Core Foundations

* ~~Initialize Go module (`go.mod`)~~
* ~~Project structure organized by responsibility (`models`, `store`, `handlers`, `router`)~~
* ~~Configuration loading via environment variables~~
* ~~Central HTTP router using `net/http`~~
* ~~Consistent JSON response helpers~~

---

###  Data Models

* ~~Book model with nested Author~~
* ~~Author model~~
* ~~Customer model~~
* ~~Order + OrderItem models~~
* ~~Address model~~
* ~~SalesReport + BookSales models~~
* ~~Proper JSON tags on all structs~~

---

###  Interfaces & Stores

* ~~`BookStore` interface defined~~
* ~~In-memory `MemStore` implementing `BookStore`~~
* ~~Thread-safe access using `sync.RWMutex`~~
* ~~Auto-incrementing IDs~~
* ~~JSON file persistence (`database.json`)~~
* ~~Load data on startup, save on mutation~~

---

###  Book API (Fully Functional)

* ~~POST `/books` – Create book~~
* ~~GET `/books/{id}` – Retrieve book by ID~~
* ~~PUT `/books/{id}` – Update book~~
* ~~DELETE `/books/{id}` – Delete book~~
* ~~GET `/books?title=...` – Search books~~
* ~~Nested author creation & normalization~~
* ~~Correct HTTP status codes~~
* ~~Error responses in JSON~~
* ~~Handle Author Book relationship~~
---

###  Testing 

* ~~Handler unit tests using `httptest`~~
* ~~Mock store implementing interfaces~~
* ~~Tests isolated from persistence~~
* ~~Covers create, read, delete, search paths~~

### Using Postman with 1000 POST requests

### Finished in 21 seconds with average response time of 9ms
* ~~Concurrent request tests~~ 
---


##  WORK TO DO

### Authors API

* ~~POST `/authors` – create author~~
* ~~GET `/authors/{id}` – retrieve author by ID~~
* ~~PUT `/authors/{id}` – update author~~
* ~~DELETE `/authors/{id}` – delete author~~
* ~~GET `/authors` – list all authors~~
* ~~In-memory author store with mutex~~
* ~~JSON persistence for authors~~
* ⬜ Author handler unit tests

---

### Customers API

* ~~POST `/customers` – create customer~~
* ~~GET `/customers/{id}` – retrieve customer~~
* ~~PUT `/customers/{id}` – update customer~~
* ~~DELETE `/customers/{id}` – delete customer~~
* ~~GET `/customers` – list customers~~
* ~~In-memory customer store with mutex~~
* ~~JSON persistence for customers~~
* ~~Customer handler~~ 
* ⬜ Customers tests
---

### Orders API

* ~~⬜ POST `/orders` – place an order~~
* ~~⬜ GET `/orders/{id}` – retrieve order~~
* ~~⬜ GET `/orders?customer_id=` – order history per customer~~
* ~~⬜ Stock validation on order creation~~
* ~~⬜ Order status lifecycle (`pending`, `paid`, `shipped`, `cancelled`)~~
* ~~⬜ Automatic stock decrement on purchase~~
* ~~⬜ In-memory order store with mutex~~
* ~~⬜ JSON persistence for orders~~
* ~~⬜ Order handler~~
* ~~⬜ Concurrency tests (simultaneous orders)~~

---

* ~~⬜ Concurrent request tests~~

---
---

##  Background Job 

### Periodic Sales Report Generation

* ~~⬜ Background goroutine with `time.Ticker`~~
* ~~⬜ Context-based lifecycle management~~
* ~~⬜ Last 24h of orders~~
* ~~⬜ Calculate:~~

  * ~~Total revenue~~
  * ~~Total orders~~
  * ~~Top-selling books~~
* ~~⬜ Persist reports to `output-reports/`~~
* ~~⬜ Filename format: `report_YYYYMMDDHHMM.json`~~

### Reports API

* ~~⬜ GET `/reports/sales`~~
* ~~⬜ Filter by date range~~

---

##  Concurrency & Context 

* ~~⬜ Pass `context.Context` through handlers~~
* ⬜ Cancel background jobs on shutdown
* ~~⬜ Respect request cancellation in long operations~~

---

##  Logging 

* ~~⬜ Request logging middleware~~
* ~~⬜ Error logging~~
* ~~⬜ Background job logging~~

---
## Metrics

### Customers Metrics

* ~~⬜ Total Customers~~

### Books Metrics

* ~~⬜ total Books~~
* ~~⬜ Books per genre ~~
* ~~⬜ Out of stock Books~

### Authors Metrics

* ~~⬜ Total Authors~~
* ~~⬜ Books per author~~

### API Metrics

* ~~⬜ Requests per Endpoint~~

---

##  Documentation

* ~~README created and maintained~~
* ⬜ Swagger / OpenAPI spec
* ⬜ Endpoint examples for all resources

---

##  Running the Project

```bash
go run main.go
```

Environment variables (optional):

```bash
export PORT=8080
export DB_PATH=internal/db/database.json
```

---

##  Testing

Run all tests:

```bash
go test ./... -v
```

Run handler tests only:

```bash
go test ./internal/http/handlers/HandlersTests -v
```

---

##  Next Milestones 

~~1. Finish Authors API~~
~~2. Implement Customers~~
~~3. Implement Orders~~
~~4. Add background sales report generator~~
~~5. Add graceful shutdown with contexts~~
6. Write Swagger spec


