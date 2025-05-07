# CryptoFolio API

CryptoFolio is a backend API for a mobile app that tracks users' cryptocurrency assets and shows their total value in USD. This service allows users to securely manage assets and retrieve live prices.

---

## 🚀 Features

- Full CRUD for crypto assets
- JWT-based stateless authentication
- Live USD value via CoinGecko API
- PostgreSQL backend (production-ready)
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
git clone https://github.com/youruser/cryptofolio.git
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

Use this token in all subsequent requests:

```
Authorization: Bearer <your_token_here>
```

---

## 📡 API Endpoints

All endpoints require a valid JWT.

### 🔐 POST /login

Authenticate a static user and receive a JWT.

---

### ➕ POST /assets

Create a new crypto asset.

**Body:**

```json
{
  "label": "Trading Wallet",
  "currency": "BTC",
  "amount": 1.25
}
```

### 📄 GET /assets

List all assets for the authenticated user.

---

### 🔍 GET /assets/{id}

Retrieve one asset by ID.

**Response:**

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

**Body:**

```json
{
  "label": "Updated Label",
  "currency": "ETH",
  "amount": 2
}
```

---

### ❌ DELETE /assets/{id}

Delete an asset.

---

### 💰 GET /assets/value/total

Returns the sum of all assets converted to USD using live exchange rates.

**Response:**

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
│   ├── config/         # Static users, Valid Currencies, Api client url
│   ├── handler/        # API routes & handlers
│   ├── model/          # Asset structure
│   ├── rate/           # CoinGecko rate fetcher
│   └── store/          # DB setup and test helpers
├── go.mod
├── README.md
├── ASSUMPTIONS.md
```

---

## 📝 Notes

- JWTs expire in 24 hours
- Only BTC, ETH, LTC are supported
- All inputs are validated at API level
- PostgreSQL used for concurrency and precision

---
