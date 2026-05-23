# Product Catalog Service

This project is a backend HTTP service built in Go.

The project includes two main components:

1. Rate-limited API
2. Product catalog with media support

All data is stored in memory using Go maps and structs.

---

# Tech Stack

- Go
- net/http
- encoding/json
- sync.Mutex
- sync.RWMutex

No external frameworks or databases were used.

---

# Project Structure

```text
cmd/server/main.go

internal/
├── handlers/
├── limiter/
├── models/
├── store/
└── utils/
```

---

# How To Run

## 1. Clone repository

```bash
git clone <your-repo-link>
```

## 2. Move into project

```bash
cd product-catalog-service
```

## 3. Run server

```bash
go run ./cmd/server
```

Server starts on:

```text
http://localhost:8080
```

---

# Part 1 - Rate Limited API

## POST /request

Accepts requests with per-user rate limiting.

### Request Body

```json
{
  "user_id": "anuj",
  "payload": {
    "message": "hello"
  }
}
```

### Success Response

```json
{
  "message": "Request accepted",
  "accepted_requests": 1
}
```

### Rate Limit Exceeded

Status:
```text
429 Too Many Requests
```

Response:
```json
{
  "error": "Rate limit exceeded"
}
```

### Validation Errors

Status:
```text
400 Bad Request
```

---

# Rate Limiting Approach

A fixed 1-minute window approach is used.

Each user is allowed:
- maximum 5 accepted requests per minute

The implementation is concurrency-safe using sync.Mutex.

The check and increment operations are performed inside the same mutex lock to prevent race conditions.

---

# GET /stats

Returns per-user request statistics.

Example response:

```json
{
  "anuj": {
    "accepted_requests": 5,
    "rejected_requests": 1,
    "window_start": "2026-05-23T14:00:00Z"
  }
}
```

Rejected requests are tracked within the current window.

---

# Part 2 - Product Catalog

## Data Model Design

Products and media are stored separately.

### Lightweight Product Store

Stores:
- id
- name
- sku
- image_count
- video_count
- created_at

### Separate Media Store

Stores:
- image_urls
- video_urls

This separation keeps the list endpoint lightweight and avoids loading unnecessary media data.

---

# POST /products

Creates a new product.

### Example Request

```json
{
  "name": "Gaming Mouse",
  "sku": "SKU-001",
  "image_urls": [
    "https://cdn.example.com/mouse1.jpg"
  ],
  "video_urls": [
    "https://cdn.example.com/demo.mp4"
  ]
}
```

### Success Response

Status:
```text
201 Created
```

---

# GET /products

Returns paginated lightweight product data.

Supports:
- limit
- offset

Example:

```text
GET /products?limit=10&offset=0
```

The list endpoint intentionally does NOT return:
- image_urls
- video_urls

to keep responses lightweight and scalable.

---

# GET /products/{id}

Returns full product detail including:
- image_urls
- video_urls

---

# POST /products/{id}/media

Appends new media URLs to an existing product.

Example request:

```json
{
  "image_urls": [
    "https://cdn.example.com/mouse2.jpg"
  ]
}
```

---

# Validation Rules

The API validates:

- non-empty name
- non-empty sku
- unique sku
- valid URLs
- maximum URL length of 2048 characters
- maximum 20 URLs per request

Supported URL schemes:
- http://
- https://

---

# Concurrency

- sync.Mutex is used for rate limiter safety
- sync.RWMutex is used for product storage

RWMutex allows multiple readers while maintaining safe writes.

---

# Production Limitations

This project uses in-memory storage only.

Current limitations:
- data is lost after server restart
- single instance only
- no distributed rate limiting
- no persistence layer
- not horizontally scalable

---

# Production Improvements

In a production environment I would use:

- PostgreSQL for persistent storage
- Redis for distributed rate limiting
- Object storage/CDN for media
- Authentication and authorization
- Structured logging
- Automated tests
- Docker deployment

---

# Example curl Commands

## Create Product

```bash
curl -X POST http://localhost:8080/products \
-H "Content-Type: application/json" \
-d '{
  "name":"Gaming Mouse",
  "sku":"SKU-001"
}'
```

## Get Products

```bash
curl http://localhost:8080/products
```

## Get Product Detail

```bash
curl http://localhost:8080/products/1
```

---

# AI Usage Disclosure

AI tools were used for:
- brainstorming
- understanding concurrency approaches
- improving README structure
- debugging small issues

All implementation, testing, and final integration were completed manually.