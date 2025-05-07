---

## Login & Authentication Limitations

- The `/login` endpoint accepts credentials over HTTPS and returns a signed JWT if valid. However:
  - No rate limiting is implemented — repeated brute-force login attempts are not blocked.
  - No account lockout policy, audit logging, or login history is tracked.
  - Passwords are stored in plain-text in memory; in production, passwords would be securely hashed (e.g., using bcrypt).
  - No refresh tokens or session renewal flow is implemented — users must log in again after 24h.
  - There is no support for multi-factor authentication (MFA).
  - All error responses are generic and do not distinguish between incorrect usernames or passwords (intentional to prevent user enumeration).

## JWT Authentication

- The secret key (`JWT_SECRET`) is read from an environment variable, but in a real application secrets should be stored securely using a secret manager (e.g., AWS Secrets Manager or Vault).
- JWTs include an expiration (`exp`) of 24 hours, but:
  - There is no mechanism for token revocation (e.g., logout, blacklisting).
  - No refresh token flow is implemented, which would be necessary for longer sessions or mobile resilience.
  - Tokens are stored entirely client-side; there is no session tracking or invalidation server-side.
- Only basic claims (`username`, `exp`, `iss`) are used — no roles or permissions.

## Exchange Rate Fetching (CoinGecko)

- The URL is stored in a config constant, but:
  - There is no retry mechanism or backoff strategy if the third-party API fails.
  - No caching is implemented, so frequent requests (e.g., `/assets` or `/assets/{id}`) may trigger rate limits or slowdowns.
  - Only three currencies (BTC, ETH, LTC) are supported; adding more requires changing code.
- The request uses a fixed 5-second timeout to prevent hangs, but no additional observability (logging, metrics) is in place for rate failures.
- In a production application, better to have:
  - Caching (e.g., in-memory or Redis)
  - API key and request authentication
  - More robust error handling and retry logic

---



