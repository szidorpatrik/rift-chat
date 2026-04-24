# Rift Chat API Documentation

## Base URL

All API endpoints (except the health check) are prefixed with `/api/v1`.

## Authentication Method

This API uses **HTTP-Only Cookies** for session management.

* **Access Token (`access_token`):** A short-lived JWT (1 hour) used to authorize requests to protected routes.
* **Refresh Token (`refresh_token`):** A long-lived token (30 days) used to securely rotate the access token when it expires.
* Both cookies are set with `HttpOnly`, `SameSite=Strict`, and conditionally `Secure` (based on the environment) flags.

---

## System Endpoints

### Health Check

Check the operational status of the server.

* **URL:** `/healthz`
* **Method:** `GET`
* **Auth Required:** No

**Response: `200 OK`**

```text
OK
```

---

## User Endpoints

### Register User

Creates a new user account.

* **URL:** `/api/v1/users`
* **Method:** `POST`
* **Auth Required:** No

**Request Body** (`application/json`)

```json
{
  "username": "test_user",
  "email": "user@example.com",
  "password": "SecurePassword123!"
}
```

**Response: `201 Created`**

```json
{
  "id": "018f3a1b-2c3d-4e5f-8a9b-0c1d2e3f4a5b",
  "username": "test_user",
  "email": "user@example.com",
  "created_at": "2026-04-24T12:00:00Z",
  "updated_at": "2026-04-24T12:00:00Z"
}
```

**Error Responses**

* `400 Bad Request`: Invalid JSON, missing fields, or failed validation (e.g., regex mismatch).
* `409 Conflict`: Username or Email already exists.
* `500 Internal Server Error`: Server failure.

---

## Authentication Endpoints

### Login

Authenticates a user and initiates a session via cookies.

* **URL:** `/api/v1/login`
* **Method:** `POST`
* **Auth Required:** No

**Request Body** (`application/json`)

```json
{
  "username": "test_user",
  "password": "SecurePassword123!"
}
```

**Response: `200 OK`**
Returns basic user info. **Crucially, this endpoint sets `Set-Cookie` headers for `access_token` and `refresh_token`.**

```json
{
  "id": "018f3a1b-2c3d-4e5f-8a9b-0c1d2e3f4a5b",
  "username": "test_user",
  "email": "user@example.com"
}
```

**Error Responses**

* `400 Bad Request`: Invalid JSON body.
* `401 Unauthorized`: Incorrect username or password.
* `500 Internal Server Error`: Server failure.

---

### Refresh Tokens

Rotates the session tokens. The frontend should call this when the access token expires.

* **URL:** `/api/v1/refresh`
* **Method:** `POST`
* **Auth Required:** Requires a valid `refresh_token` cookie.

**Request Body**
None.

**Response: `200 OK`**
Returns an empty string. **This endpoint sets new `Set-Cookie` headers, rotating both the `access_token` and the `refresh_token`.**

```json
""
```

**Error Responses**

* `401 Unauthorized`: Missing `refresh_token` cookie, or the token is invalid/revoked.
* `500 Internal Server Error`: Server failure during rotation.

---

### Revoke Tokens (Logout)

Ends the user's session, revokes the refresh token in the database, and clears the client's cookies.

* **URL:** `/api/v1/revoke`
* **Method:** `POST`
* **Auth Required:** No authentication required, but expects a valid `refresh_token` cookie to identify the session to revoke.

**Request Body**
None.

**Response: `200 OK`**
Returns an empty string. **This endpoint sends `Set-Cookie` headers with `MaxAge: -1` to instantly delete the `access_token` and `refresh_token` cookies on the client.**

```json
""
```
