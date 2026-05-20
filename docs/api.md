# Rift Chat API Documentation

## Base URL

All REST API endpoints (except the health check and WebTransport) are prefixed with `/api/v1`.

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

### WebTransport Stream

Establishes a WebTransport session to stream real-time events via a one-way downlink.

* **URL:** `/wt`
* **Method:** `CONNECT` (WebTransport Protocol)
* **Auth Required:** Optional (Unauthenticated connections receive a nil UUID for testing)

**Data Format:** Raw byte stream containing JSON payloads appended with newline characters (`\n`).

---

## User & Authentication Endpoints

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

---

### Login

Authenticates a user and initiates a session via cookies. Sets `Set-Cookie` headers for `access_token` and `refresh_token`.

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

```json
{
  "id": "018f3a1b-2c3d-4e5f-8a9b-0c1d2e3f4a5b",
  "username": "test_user",
  "email": "user@example.com",
  "status": "online"
}
```

---

### Session Check

Verifies the user's active session and returns their profile.

* **URL:** `/api/v1/session`
* **Method:** `GET`
* **Auth Required:** Yes

**Response: `200 OK`**

```json
{
  "id": "018f3a1b-2c3d-4e5f-8a9b-0c1d2e3f4a5b",
  "username": "test_user",
  "email": "user@example.com",
  "status": "online"
}
```

---

### Refresh Tokens

Rotates the session tokens. Sets new `Set-Cookie` headers.

* **URL:** `/api/v1/refresh`
* **Method:** `POST`
* **Auth Required:** Requires valid `refresh_token` cookie

**Response: `204 No Content`**

---

### Revoke Tokens (Logout)

Ends the user's session and clears client cookies instantly.

* **URL:** `/api/v1/revoke`
* **Method:** `POST`
* **Auth Required:** No strict auth, expects `refresh_token` cookie

**Response: `204 No Content`**

---

## Server Endpoints

### Get Public Servers

Retrieves a list of all public servers.

* **URL:** `/api/v1/servers/public`
* **Method:** `GET`
* **Auth Required:** No

**Response: `200 OK`**

```json
[
  {
    "id": "018f3a1b-2c3d-4e5f-8a9b-0c1d2e3f4a5b",
    "name": "Public Server",
    "is_public": true
  }
]
```

---

### Get User Servers

Retrieves servers the authenticated user is a member of.

* **URL:** `/api/v1/servers`
* **Method:** `GET`
* **Auth Required:** Yes

**Response: `200 OK`**

```json
[
  {
    "id": "018f3a1b-2c3d-4e5f-8a9b-0c1d2e3f4a5b",
    "owner_id": "018f3a1b-2c3d-4e5f-8a9b-111111111111",
    "name": "My Server",
    "is_public": false,
    "description": "A private server"
  }
]
```

---

### Create Server

Creates a server and automatically joins the creator as owner.

* **URL:** `/api/v1/servers`
* **Method:** `POST`
* **Auth Required:** Yes

**Request Body** (`application/json`)

```json
{
  "name": "New Public Server",
  "description": "A place to chat.",
  "is_public": true
}
```

**Response: `201 Created`**

---

### Join Server

Joins an existing public server.

* **URL:** `/api/v1/servers/join`
* **Method:** `POST`
* **Auth Required:** Yes

**Request Body** (`application/json`)

```json
{
  "server_id": "018f3a1b-2c3d-4e5f-8a9b-0c1d2e3f4a5b"
}
```

**Response: `200 OK`**

---

### Get Server Members

Retrieves the list of members for a specific server.

* **URL:** `/api/v1/servers/{server_id}/members`
* **Method:** `GET`
* **Auth Required:** Yes (Requires server membership)

**Response: `200 OK`**

```json
[
  {
    "id": "018f3a1b-2c3d-4e5f-8a9b-111111111111",
    "username": "test_user",
    "status": "online"
  }
]
```

---

## Channel Endpoints

### Get Server Channels

Retrieves all channels associated with a specific server.

* **URL:** `/api/v1/servers/{server_id}/channels`
* **Method:** `GET`
* **Auth Required:** Yes (Requires server membership)

**Response: `200 OK`**

```json
[
  {
    "id": "018f3a1b-2c3d-4e5f-8a9b-222222222222",
    "server_id": "018f3a1b-2c3d-4e5f-8a9b-0c1d2e3f4a5b",
    "name": "general-chat",
    "type": 0,
    "position": 1
  }
]
```

---

### Create Channel

Creates a new channel within a server. Channel types: `0` (text), `1` (voice), `2` (category). Text channels format names to lowercase and replace spaces with hyphens.

* **URL:** `/api/v1/servers/{server_id}/channels`
* **Method:** `POST`
* **Auth Required:** Yes (Requires server membership)

**Request Body** (`application/json`)

```json
{
  "name": "General Chat",
  "type": 0,
  "position": 1
}
```

**Response: `201 Created`**

---

### Delete Channel

Deletes a specific channel by ID.

* **URL:** `/api/v1/servers/{server_id}/channels/{channel_id}`
* **Method:** `DELETE`
* **Auth Required:** Yes (Requires server membership)

**Response: `204 No Content`**
