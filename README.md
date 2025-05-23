# CryptoFolio API

CryptoFolio is a backend API for a mobile app that tracks users' cryptocurrency assets and shows their total value in USD. This service allows users to securely manage assets and retrieve live prices.

---

## ⚡ Quickstart

```bash
git clone https://github.com/drashevskyi/cryptofolio.git
cd cryptofolio
export POSTGRES_DSN="postgres://postgres:postgres@localhost:5432/cryptofolio?sslmode=disable"
export JWT_SECRET="your-secret"
createdb cryptofolio
go mod tidy
go build ./cmd/server
./server
```

---

## 🚀 Features

- Full CRUD for crypto assets
- JWT-based stateless authentication
- Live USD value via CoinGecko API
- PostgreSQL backend
- Validates supported currencies and amount
- Includes test coverage for key endpoints

---

## 🧱 Requirements

- Go 1.21+
- PostgreSQL 13+
- CoinGecko API (no key required)
- Environment variables set (see below)

---

## ⚙️ Environment Setup

Set the following environment variables before running:

```bash
export POSTGRES_DSN="postgres://postgres:postgres@localhost:5432/cryptofolio?sslmode=disable"
export JWT_SECRET="your-secure-jwt-secret"
```

---

## 🧪 Running Locally

### 1. Clone and install dependencies

```bash
git clone https://github.com/drashevskyi/cryptofolio.git
cd cryptofolio
go mod tidy
```

### 2. Start PostgreSQL

Make sure your PostgreSQL instance is running and a database named `cryptofolio` exists:

```bash
createdb cryptofolio
```

### 3. Build and run

```bash
go build ./cmd/server
./server
```

Server will start at: [http://localhost:8080](http://localhost:8080)

---

## 🔐 Authentication

All routes (except `/login`) require a **JWT token** in the `Authorization` header.

### 🔑 Static Users

| Username | Password   |
|----------|------------|
| user1    | password1  |
| user2    | password2  |

### 🔁 Login Example

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user1", "password":"password1"}'
```

Returns:

```json
{
  "token": "eyJhbGciOi..."
}
```

Use the returned token in the `Authorization` header like this:

```bash
Authorization: Bearer eyJhbGciOi...
```

---

## 📡 API Endpoints

All endpoints require a valid JWT unless otherwise noted.

---

### 🔑 POST /login

Authenticate using static credentials and receive a JWT.

---

### ➕ POST /assets

Create a new crypto asset.

**Request Body:**

```json
{
  "label": "Trading Wallet",
  "currency": "BTC",
  "amount": 1.25
}
```

```bash
curl -X POST http://localhost:8080/assets \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "label": "My ETH Wallet",
    "currency": "ETH",
    "amount": 1.75
}'
```

---

### 📄 GET /assets

List all assets for the authenticated user.

```bash
curl http://localhost:8080/assets \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

---

### 🔍 GET /assets/{id}

Retrieve one asset by ID.

```bash
curl http://localhost:8080/assets/1 \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

**Example Response:**

```json
{
  "id": 3,
  "label": "Trading Wallet",
  "currency": "BTC",
  "amount": 1.25,
  "usd_value": 43125.00
}
```

---

### ✏️ PUT /assets/{id}

Update an existing asset.

**Request Body:**

```json
{
  "label": "Updated Label",
  "currency": "ETH",
  "amount": 2
}
```

```bash
curl -X PUT http://localhost:8080/assets/1 \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "label": "Updated Wallet Name",
    "currency": "BTC",
    "amount": 0.9
}'
```

---

### ❌ DELETE /assets/{id}

Delete an asset.

```bash
curl -X DELETE http://localhost:8080/assets/1 \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

---

### 💰 GET /assets/value/total

Returns the sum of all assets converted to USD using live exchange rates.

```bash
curl http://localhost:8080/assets/value/total \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

**Example Response:**

```json
{
  "total_usd": 98250.00
}
```

---

## 🧪 Running Tests

```bash
go test ./...
```

You’ll see output like:

```
ok  	cryptofolio/internal/handler	0.245s
```

---

## 📂 Project Structure

```
cryptofolio/
├── cmd/server/         # Entrypoint (main.go)
├── internal/
│   ├── auth/           # JWT auth logic
│   ├── config/         # Static users, API client config
│   ├── handler/        # API routes & handlers
│   ├── model/          # Asset structure + validation
│   ├── rate/           # CoinGecko rate fetcher
│   └── store/          # DB setup and test helpers
├── go.mod
├── README.md
├── ASSUMPTIONS.md
```

---

## 📝 Notes

- JWTs expire in 24 hours
- Only BTC, ETH, LTC are supported (enforced in code and DB)
- All inputs are validated at API level
- PostgreSQL is used for data integrity and concurrency
- Exchange rates are fetched dynamically from CoinGecko (no API key needed)

---
