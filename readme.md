# BOOKCABIN – Flight Search Service

A lightweight **flight search aggregation service** built in Go.
The service queries multiple airline providers concurrently, normalizes results, applies filters/sorting, and returns unified search results via a **GET-only REST API**.

This project follows **Clean Architecture** and **SOLID principles**, making it easy to extend (new providers, caching layer, transport) without touching core business logic.

---

## 🚀 How to Run

### 1️⃣ Prerequisites

* Go **1.21+** recommended
* Unix-like OS (macOS / Linux / WSL)

Verify Go:

```bash
go version
```

---

### 2️⃣ Clone & Install

```bash
git clone <your-repo-url>
cd bookcabin
go mod tidy
```

---

### 3️⃣ Run the API Server

```bash
go run ./cmd/api
```

The server will start on:

```
http://localhost:8080
```

---

## 🔍 Search API (GET only)

### Endpoint

```
GET /search
```

### Example Request

```bash
curl "http://localhost:8080/search?origin=CGK&destination=DPS&departure_date=2025-12-15&passengers=1&cabin_class=economy&airlines=GA,ID&min_price=500000&max_price=2000000&sort_by=price_asc"
```

### Required Query Params

| Param          | Description              |
| -------------- | ------------------------ |
| origin         | Origin airport code      |
| destination    | Destination airport code |
| departure_date | YYYY-MM-DD               |
| passengers     | Number of passengers     |
| cabin_class    | economy / business       |

### Optional Filters

| Param              | Description                         |
| ------------------ | ----------------------------------- |
| min_price          | Minimum price (IDR)                 |
| max_price          | Maximum price (IDR)                 |
| max_stops          | Maximum allowed stops               |
| airlines           | Airline codes (CSV or repeated)     |
| max_duration       | Max duration (minutes)              |
| earliest_departure | HH:MM                               |
| latest_departure   | HH:MM                               |
| earliest_arrival   | HH:MM                               |
| latest_arrival     | HH:MM                               |
| sort_by            | price_asc, price_desc, duration_asc, duration_desc, departure_asc, arrival_asc best_value |

---

## 🧱 Project Structure (Clean Architecture)

```
cmd/api
  └── main.go            # Application entry point

internal/
  common/                # Shared utilities (sorting, helpers)
  domain/                # Core business models & rules
    ├── flight.go        # Flight entity
    └── search.go        # SearchRequest, FlightFilter

  handler/               # HTTP layer (transport)
    └── flight_handler.go

  service/               # Use cases / business logic
    ├── flight_interface.go
    └── search_flights.go

  provider/              # External airline integrations
    ├── airasia.go
    ├── batik.go
    ├── garuda.go
    ├── lion.go

  infra/                 # Infrastructure concerns
    └── cache.go         # In-memory TTL cache

  mock/                  # Mock providers & fixtures
    ├── *.json
    └── mock_*.go
```

---

## 🧠 Clean Architecture Mapping

| Layer              | Responsibility                             |
| ------------------ | ------------------------------------------ |
| Handler            | HTTP parsing, validation, response         |
| Service (Use Case) | Orchestrates providers, caching, filtering |
| Domain             | Pure business rules & entities             |
| Provider           | External airline data sources              |
| Infra              | Cache, IO, implementations                 |

### Dependency Rule

> **Dependencies always point inward**
> Outer layers depend on inner layers, never the reverse.

---

## 🧩 SOLID Principles Applied

### ✅ Single Responsibility

* Handlers only parse HTTP
* Use cases only orchestrate logic
* Providers only fetch airline data

### ✅ Open / Closed

* Add new airline provider without touching existing logic
* Implement `FlightProvider` interface

### ✅ Liskov Substitution

* All providers implement the same interface

### ✅ Interface Segregation

* Small focused interfaces (`FlightProvider`, `Repository`)

### ✅ Dependency Inversion

* Use cases depend on abstractions, not implementations

---

## ⚡ Concurrency & Performance

* Providers called **concurrently** using goroutines
* Context timeout (5s)
* In-memory cache for raw provider results
* Filters & sorting applied after cache

---

## 🧠 Caching Strategy

* Cache **raw provider results**
* Keyed by origin, destination, date, pax, cabin
* TTL: **~3 minutes** (configurable)
* Filters do NOT affect cache key

---

## 📜 Swagger (OpenAPI)

### Install swag

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### Add Swagger Annotations (example)

```go
// Search flights
// @Summary Search flights
// @Description Search flights with filters
// @Tags Flights
// @Accept json
// @Produce json
// @Param origin query string true "Origin airport"
// @Param destination query string true "Destination airport"
// @Param departure_date query string true "Departure date"
// @Param passengers query int true "Passengers"
// @Param cabin_class query string true "Cabin class"
// @Param airlines query string false "Airlines CSV"
// @Success 200 {object} map[string][]domain.Flight
// @Router /search [get]
```

### Generate Docs

```bash
swag init -g cmd/api/main.go
```

### Swagger UI

```
http://localhost:8080/swagger/index.html
```

---

## 🧪 Testing

* Mock providers in `internal/mock`
* Replace real providers in `main.go`
* Deterministic test data via JSON fixtures

---

## 🚧 Future Improvements

* Redis cache
* Rate limiting
* Pagination
* Circuit breaker per provider
* Stale-while-revalidate

---

## 👨‍💻 Author Notes

This codebase is intentionally designed for:

* System design interviews
* Real-world extensibility
* Production-ready patterns

---

Happy flying ✈️
