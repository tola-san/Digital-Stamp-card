# Digital Stamp Card API Design

This document is the API contract for the Digital Stamp Card backend. It separates endpoints that are live today from endpoints planned for the loyalty MVP.

## Base URLs

| Environment | Base URL |
| --- | --- |
| Production | `https://digital-stamp-card.onrender.com/api` |
| Local | `http://localhost:8080/api` |

All request and response bodies use JSON. Times use ISO 8601 UTC timestamps, and resource IDs are UUIDs.

## API conventions

### Success responses

Successful responses are wrapped in `data`:

```json
{
  "data": {}
}
```

### Error responses

Errors are wrapped in `error`:

```json
{
  "error": {
    "code": "DATABASE_UNAVAILABLE",
    "message": "database is unavailable"
  }
}
```

### Authentication

Customer and staff sessions will use secure, HTTP-only cookies. Browser requests must include credentials. Staff-only endpoints require an authenticated staff session; customer-only endpoints require an authenticated customer session.

## Implemented endpoint

### `GET /health`

Checks that the API and PostgreSQL connection are available.

**Response — `200 OK`**

```json
{
  "data": {
    "status": "api is running"
  }
}
```

**Response — `503 Service Unavailable`**

```json
{
  "error": {
    "code": "DATABASE_UNAVAILABLE",
    "message": "database is unavailable"
  }
}
```

## Planned MVP endpoints

The endpoints below are API design targets. They are not implemented yet.

### Customer session and profile

| Method | Path | Session | Purpose |
| --- | --- | --- | --- |
| `POST` | `/customers` | Public | Register a customer and create their card. |
| `POST` | `/customer-sessions` | Public | Start a customer session using a phone number. |
| `DELETE` | `/customer-sessions/current` | Customer | End the current customer session. |
| `GET` | `/customers/me` | Customer | Get the signed-in customer. |
| `GET` | `/customers/me/card` | Customer | Get stamp-card progress. |
| `GET` | `/customers/me/transactions` | Customer | Get stamp and reward history. |

#### `POST /customers`

```json
{
  "name": "Tola San",
  "phone": "+85512345678"
}
```

**Response — `201 Created`**

```json
{
  "data": {
    "customer": {
      "id": "4b0e624d-7d44-444a-aad1-010101010101",
      "name": "Tola San",
      "phone": "+85512345678"
    },
    "card": {
      "id": "a0e624d0-7d44-444a-aad1-010101010101",
      "stamp_count": 0,
      "required_stamps": 10
    }
  }
}
```

### Staff authentication and customers

| Method | Path | Session | Purpose |
| --- | --- | --- | --- |
| `POST` | `/staff-sessions` | Public | Sign in with email and password. |
| `DELETE` | `/staff-sessions/current` | Staff | Sign out. |
| `GET` | `/staff/me` | Staff | Get the signed-in staff member. |
| `GET` | `/staff/customers` | Staff | List and search customers. |
| `GET` | `/staff/customers/:customerId` | Staff | Get customer, card, and recent activity. |
| `POST` | `/staff/customers/:customerId/stamps` | Staff | Add or reverse a stamp with an audit record. |

#### `POST /staff-sessions`

```json
{
  "email": "staff@goodcoffee.com",
  "password": "not-returned-in-responses"
}
```

#### `POST /staff/customers/:customerId/stamps`

```json
{
  "stamp_delta": 1,
  "reason": "Manual purchase correction"
}
```

`stamp_delta` must be `1` to add a stamp or `-1` to reverse one. The API must reject a reversal that would make the card balance negative.

### Temporary stamp QR codes

| Method | Path | Session | Purpose |
| --- | --- | --- | --- |
| `POST` | `/staff/stamp-qrs` | Staff | Generate a one-time, short-lived stamp QR code. |
| `GET` | `/staff/stamp-qrs/current` | Staff | Get the staff member's active QR code. |
| `DELETE` | `/staff/stamp-qrs/:stampQrId` | Staff | Cancel an active code. |
| `GET` | `/staff/stamp-qrs` | Staff | View QR history. |
| `POST` | `/stamps/claims` | Customer | Claim one stamp from a QR token. |

#### `POST /staff/stamp-qrs`

**Response — `201 Created`**

```json
{
  "data": {
    "id": "1ee624d0-7d44-444a-aad1-010101010101",
    "token": "only-returned-when-created",
    "expires_at": "2026-09-06T09:10:00Z",
    "status": "ACTIVE"
  }
}
```

The database stores only a hash of `token`. A token is valid once, expires after 60 seconds, and becomes unavailable immediately after a successful claim.

#### `POST /stamps/claims`

```json
{
  "token": "scanned-qr-token"
}
```

The claim runs in one database transaction: validate the customer session and QR, increment the card, create a `STAMP_ADDED` transaction, and mark the QR as used. This prevents duplicate claims.

### Rewards

| Method | Path | Session | Purpose |
| --- | --- | --- | --- |
| `GET` | `/rewards` | Public | List active rewards. |
| `GET` | `/customers/me/rewards` | Customer | List rewards and customer eligibility. |
| `POST` | `/staff/reward-redemptions` | Staff | Redeem an eligible customer reward. |

#### `POST /staff/reward-redemptions`

```json
{
  "customer_id": "4b0e624d-7d44-444a-aad1-010101010101",
  "reward_id": "8fa624d0-7d44-444a-aad1-010101010101"
}
```

The API verifies eligibility, updates stamp balance according to the reward rule, and records a `REWARD_REDEEMED` transaction atomically.

## Resource shapes

### Customer

```json
{
  "id": "uuid",
  "name": "Tola San",
  "phone": "+85512345678",
  "created_at": "2026-09-06T09:00:00Z",
  "updated_at": "2026-09-06T09:00:00Z"
}
```

### Stamp card

```json
{
  "id": "uuid",
  "customer_id": "uuid",
  "stamp_count": 7,
  "required_stamps": 10,
  "created_at": "2026-09-06T09:00:00Z",
  "updated_at": "2026-09-06T09:00:00Z"
}
```

### Transaction

```json
{
  "id": "uuid",
  "customer_id": "uuid",
  "staff_id": "uuid",
  "stamp_qr_id": "uuid",
  "type": "STAMP_ADDED",
  "stamp_delta": 1,
  "created_at": "2026-09-06T09:00:00Z"
}
```

Allowed transaction types are `STAMP_ADDED`, `STAMP_REVERSED`, and `REWARD_REDEEMED`. QR states are `ACTIVE`, `USED`, `EXPIRED`, and `CANCELLED`.

## Status codes

| Status | Meaning |
| --- | --- |
| `200` | Request succeeded. |
| `201` | Resource created. |
| `204` | Request succeeded without a response body. |
| `400` | Invalid request body or malformed parameter. |
| `401` | No valid session. |
| `403` | Valid session without sufficient permission. |
| `404` | Resource or route not found. |
| `409` | Duplicate customer, stale QR, or reward conflict. |
| `422` | Valid JSON but business-rule validation failed. |
| `503` | Database unavailable. |
