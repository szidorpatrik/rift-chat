# Rift Chat API Documentation

## Base URL

All REST API endpoints (except the health check and WebTransport) are prefixed with `/api/v1`.

## Authentication Method

This API uses **HTTP-Only Cookies** for session management.

* **Access Token (`access_token`):** A short-lived JWT (1 hour) used to authorize requests to protected routes.
* **Refresh Token (`refresh_token`):** A long-lived token (30 days) used to securely rotate the access token when it expires.
* Both cookies are set with `HttpOnly`, `SameSite=Strict`, and conditionally `Secure` flags.

---

## System Endpoints

### Health Check

Check the operational status of the server.

* **URL:** `/healthz`
* **Method:** `GET`
* **Auth Required:** No

**Response: `200 OK`** (`text/plain; charset=utf-8`)

```text
OK
```

---

### WebTransport Stream

Establishes an upgraded WebTransport session to stream real-time events via a one-way server-to-client downlink. Subscribes the client session to a Valkey topic matching `user_events:{user_id}`.

* **URL:** `/wt`
* **Method:** `CONNECT` (WebTransport Protocol Upgrade)
* **Auth Required:** Yes

**Data Format:** Raw byte stream containing JSON payloads with newline characters (`\n`) appended at the end of each message event.

---

## User & Authentication Endpoints

### Register User

Creates a new user account. Sanitizes trailing/leading spaces on the username and email inputs before execution.

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

**Response: `201 Created`** (`application/json`)

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

Authenticates a user and initiates a session via cookies. Sets `Set-Cookie` headers for `access_token` and `refresh_token`. Includes an anti-timing attack mechanism if a non-existent user profile is requested.

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

**Response: `200 OK`** (`application/json`)

```json
{
  "id": "018f3a1b-2c3d-4e5f-8a9b-0c1d2e3f4a5b",
  "username": "test_user",
  "email": "user@example.com",
  "status": "online",
  "avatar_url": "[https://cdn.example.com/avatars/user.png](https://cdn.example.com/avatars/user.png)",
  "description": "An engineer working on Rift Chat"
}
```

---

### Session Check

Verifies the user's active session state using their context credentials and returns their fresh database profile.

* **URL:** `/api/v1/session`
* **Method:** `GET`
* **Auth Required:** Yes

**Response: `200 OK`** (`application/json`)

```json
{
  "id": "018f3a1b-2c3d-4e5f-8a9b-0c1d2e3f4a5b",
  "username": "test_user",
  "email": "user@example.com",
  "status": "online",
  "avatar_url": "[https://cdn.example.com/avatars/user.png](https://cdn.example.com/avatars/user.png)",
  "description": "An engineer working on Rift Chat"
}
```

---

### Refresh Tokens

Rotates old session tokens inside a single database transaction context. Revokes the utilized refresh token and appends two new `Set-Cookie` headers to the network response.

* **URL:** `/api/v1/refresh`
* **Method:** `POST`
* **Auth Required:** Requires valid `refresh_token` cookie

**Response: `204 No Content`**

---

### Revoke Tokens (Logout)

Ends the user's session by revoking the client's current refresh token within the database and instantly purging all authorization cookies from the user client.

* **URL:** `/api/v1/revoke`
* **Method:** `POST`
* **Auth Required:** No strict access auth (expects `refresh_token` cookie value)

**Response: `204 No Content`**

---

## Server Endpoints

### Get Public Servers

Retrieves an array of all servers set to public accessibility in the system.

* **URL:** `/api/v1/servers/public`
* **Method:** `GET`
* **Auth Required:** No

**Response: `200 OK`** (`application/json`)

```json
[
  {
    "id": "018f3a1b-2c3d-4e5f-8a9b-0c1d2e3f4a5b",
    "owner_id": "018f3a1b-2c3d-4e5f-8a9b-111111111111",
    "name": "Public Lounge",
    "icon_url": "[https://cdn.example.com/icons/server.png](https://cdn.example.com/icons/server.png)",
    "background_url": "[https://cdn.example.com/bg/server.jpg](https://cdn.example.com/bg/server.jpg)",
    "is_public": true,
    "description": "A place for global developer discussions.",
    "created_at": "2026-05-20T19:40:00Z",
    "updated_at": "2026-05-20T19:40:00Z"
  }
]
```

---

### Get User Servers

Retrieves an array of all servers the currently authenticated context user is an active member of.

* **URL:** `/api/v1/servers`
* **Method:** `GET`
* **Auth Required:** Yes

**Response: `200 OK`** (`application/json`)

```json
[
  {
    "id": "018f3a1b-2c3d-4e5f-8a9b-0c1d2e3f4a5b",
    "owner_id": "018f3a1b-2c3d-4e5f-8a9b-111111111111",
    "name": "My Workspace",
    "icon_url": "",
    "background_url": "",
    "is_public": false,
    "description": "A completely private server",
    "created_at": "2026-05-20T19:40:00Z",
    "updated_at": "2026-05-20T19:40:00Z"
  }
]
```

---

### Create Server

Creates a server configuration and executes an internal atomized transaction to immediately join the creator as its owner. Inputs are stripped of trailing whitespaces during sanitization validation phases.

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

**Response: `201 Created`** (`application/json`)

```json
{
  "id": "018f3a1b-2c3d-4e5f-8a9b-333333333333",
  "owner_id": "018f3a1b-2c3d-4e5f-8a9b-111111111111",
  "name": "New Public Server",
  "icon_url": "",
  "background_url": "",
  "is_public": true,
  "description": "A place to chat.",
  "created_at": "2026-05-20T21:40:00Z",
  "updated_at": "2026-05-20T21:40:00Z"
}
```

---

### Join Server

Joins an existing public server using its structural ID identifier. Returns a `409 Conflict` if membership relations already exist for this user node inside the target server.

* **URL:** `/api/v1/servers/join`
* **Method:** `POST`
* **Auth Required:** Yes

**Request Body** (`application/json`)

```json
{
  "server_id": "018f3a1b-2c3d-4e5f-8a9b-0c1d2e3f4a5b"
}
```

**Response: `200 OK`** (`application/json`)

```json
{
  "id": "018f3a1b-2c3d-4e5f-8a9b-0c1d2e3f4a5b",
  "owner_id": "018f3a1b-2c3d-4e5f-8a9b-111111111111",
  "name": "Public Lounge",
  "icon_url": "[https://cdn.example.com/icons/server.png](https://cdn.example.com/icons/server.png)",
  "background_url": "[https://cdn.example.com/bg/server.jpg](https://cdn.example.com/bg/server.jpg)",
  "is_public": true,
  "description": "A place for global developer discussions.",
  "created_at": "2026-05-20T19:40:00Z",
  "updated_at": "2026-05-20T19:40:00Z"
}
```

---

### Get Server Members

Retrieves the structural list of profile objects representing members joined to the targeted server reference.

* **URL:** `/api/v1/servers/{server_id}/members`
* **Method:** `GET`
* **Auth Required:** Yes (Requires server membership context validation)

**Response: `200 OK`** (`application/json`)

```json
[
  {
    "id": "018f3a1b-2c3d-4e5f-8a9b-111111111111",
    "username": "test_user",
    "status": "online",
    "avatar_url": "[https://cdn.example.com/avatars/user.png](https://cdn.example.com/avatars/user.png)",
    "description": "An engineer working on Rift Chat"
  }
]
```

---

## Channel Endpoints

### Get Server Channels

Retrieves all text, voice, or category channels linked directly with the specific server resource.

* **URL:** `/api/v1/servers/{server_id}/channels`
* **Method:** `GET`
* **Auth Required:** Yes (Requires server membership context validation)

**Response: `200 OK`** (`application/json`)

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

Creates a new channel layout inside a targeted server framework.

**Validation Mapping Rules:**

* Valid channel types are: `0` (text), `1` (voice), `2` (category).
* **Type `0` (Text) Transformation Rule:** When creating a text channel, all inner whitespace segments are converted to standard hyphens (`-`), and the channel name string is transformed into full lowercasing automatically.
* **URL:** `/api/v1/servers/{server_id}/channels`
* **Method:** `POST`
* **Auth Required:** Yes (Requires server membership context validation)

**Request Body** (`application/json`)

```json
{
  "name": "General Chat",
  "type": 0,
  "position": 1
}
```

**Response: `201 Created`** (`application/json`)

```json
{
  "id": "018f3a1b-2c3d-4e5f-8a9b-222222222222",
  "server_id": "018f3a1b-2c3d-4e5f-8a9b-0c1d2e3f4a5b",
  "name": "general-chat",
  "type": 0,
  "position": 1
}
```

---

### Delete Channel

Deletes a specific channel by its unique identifier tracking parameter.

* **URL:** `/api/v1/servers/{server_id}/channels/{channel_id}`
* **Method:** `DELETE`
* **Auth Required:** Yes (Requires server membership context validation)

**Response: `204 No Content`**
