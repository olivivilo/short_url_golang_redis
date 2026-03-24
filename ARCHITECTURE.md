# Short URL Service - Architecture

## Overview
Industrial-grade short URL service built with Go and Redis, following enterprise best practices for code organization, testing, and deployment.

## Design Decisions

### 1. Short Code Generation
- **Strategy**: Auto-increment ID + Base62 encoding
- **Minimum Length**: 6 characters (left-padded)
- **Alphabet**: `0-9A-Za-z` (62 characters)
- **Collision**: Theoretically impossible with sequential IDs
- **Scalability**: Supports ~56.8 billion URLs (62^6)

### 2. API Design
- **Style**: RESTful JSON API
- **Version**: `/api/v1` prefix
- **Redirect**: `302 Found` (temporary redirect)
- **Error Codes**:
  - `404 Not Found`: Code doesn't exist
  - `410 Gone`: Expired or deleted
  - `201 Created`: Successful creation
  - `204 No Content`: Successful deletion

### 3. Redis Data Model

#### Key Patterns
```
url:{code}          -> Hash (metadata)
url:{code}:visits   -> String (counter)
global:url_id       -> String (auto-increment)
```

#### Hash Fields (`url:{code}`)
- `original_url`: Target URL
- `created_at`: RFC3339 timestamp
- `expire_at`: RFC3339 timestamp (optional)
- `deleted_at`: RFC3339 timestamp (optional, soft delete)
- `note`: User-provided note/tag (optional)

#### TTL Strategy
- Redis TTL set on both `url:{code}` and `url:{code}:visits`
- Automatic cleanup when expired
- Expired URLs return `410 Gone`

### 4. Statistics
- **Visit Counter**: Incremented only on successful redirects
- **No Counting**: Expired, deleted, or non-existent URLs
- **Persistence**: Counter stored separately for flexible TTL management

### 5. Expiration Semantics
- **Expired URL**: Returns `410 Gone`
- **Deleted URL**: Returns `410 Gone`
- **Not Found**: Returns `404 Not Found`
- **No Stats After Expiry**: Visit counter stops incrementing

## Project Structure

```
short_url/
├── cmd/
│   └── shorturl/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── http/
│   │   ├── handler/
│   │   │   ├── health.go        # Health check handlers
│   │   │   └── url.go           # URL CRUD handlers
│   │   ├── middleware/
│   │   │   ├── logging.go       # Request logging
│   │   │   ├── recover.go       # Panic recovery
│   │   │   └── requestid.go     # Request ID injection
│   │   └── router.go            # Route definitions
│   ├── service/
│   │   └── url.go               # Business logic layer
│   ├── repository/
│   │   └── redis/
│   │       └── url.go           # Redis data access
│   ├── id/
│   │   └── generator.go         # ID generation & Base62
│   └── model/
│       └── url.go               # Domain models
├── test/
│   ├── integration/
│   │   └── url_test.go          # Integration tests
│   └── testutil/
│       └── redis.go             # Test utilities
├── configs/
│   └── config.example.yaml      # Example configuration
├── deploy/
│   ├── docker-compose.yml       # Local development stack
│   └── Dockerfile               # Application container
├── .vscode/
│   ├── launch.json              # Debug configurations
│   └── tasks.json               # Build tasks
├── go.mod
├── go.sum
├── README.md
└── ARCHITECTURE.md              # This file
```

## Component Responsibilities

### HTTP Layer (`internal/http/`)
- Request parsing and validation
- Response formatting
- Middleware chain (logging, recovery, request ID)
- Route registration
- HTTP status code mapping

### Service Layer (`internal/service/`)
- Business logic orchestration
- Input validation
- Error handling and mapping
- Transaction coordination
- Expiry calculation

### Repository Layer (`internal/repository/redis/`)
- Redis operations (CRUD)
- Key management
- TTL handling
- Atomic operations (INCR, HSETNX)
- Connection pooling

### ID Generator (`internal/id/`)
- Auto-increment ID fetching
- Base62 encoding/decoding
- Minimum length padding

## API Endpoints

### Create Short URL
```
POST /api/v1/urls
Content-Type: application/json

{
  "url": "https://example.com/very/long/path",
  "expire_in": "24h",           // Optional: Go duration format
  "note": "Campaign link"        // Optional
}

Response 201:
{
  "code": "abc123",
  "short_url": "http://localhost:8080/r/abc123",
  "original_url": "https://example.com/very/long/path",
  "created_at": "2026-03-23T08:00:00Z",
  "expire_at": "2026-03-24T08:00:00Z",
  "deleted_at": null,
  "visit_count": 0
}
```

### Redirect
```
GET /r/{code}

Response 302: Location: <original_url>
Response 404: Code not found
Response 410: Expired or deleted
```

### Get URL Details
```
GET /api/v1/urls/{code}

Response 200:
{
  "code": "abc123",
  "short_url": "http://localhost:8080/r/abc123",
  "original_url": "https://example.com/very/long/path",
  "created_at": "2026-03-23T08:00:00Z",
  "expire_at": "2026-03-24T08:00:00Z",
  "deleted_at": null,
  "visit_count": 42
}
```

### Delete Short URL
```
DELETE /api/v1/urls/{code}

Response 204: No Content
Response 404: Code not found
```

### Health Checks
```
GET /healthz    # Process health (always returns 200)
GET /readyz     # Readiness check (pings Redis)
```

## Configuration

Environment variables:
- `SERVER_PORT`: HTTP server port (default: 8080)
- `REDIS_ADDR`: Redis address (default: localhost:6379)
- `REDIS_PASSWORD`: Redis password (optional)
- `REDIS_DB`: Redis database number (default: 0)
- `BASE_URL`: Base URL for short links (default: http://localhost:8080)
- `LOG_LEVEL`: Logging level (default: info)

## Testing Strategy

### Unit Tests
- ID generator (Base62 encoding/decoding, padding)
- Service layer (business logic, validation)
- Middleware (logging, recovery, request ID)

### Integration Tests
- Full API workflow (create → redirect → stats → delete)
- Redis operations (TTL, expiry, soft delete)
- Error scenarios (not found, expired, deleted)

### Test Utilities
- Redis container management (testcontainers or docker-compose)
- Test data fixtures
- Assertion helpers

## Deployment

### Local Development
```bash
docker-compose up -d redis
go run cmd/shorturl/main.go
```

### Docker
```bash
docker-compose up
```

### Production Considerations (Future)
- Horizontal scaling (stateless service)
- Redis cluster/sentinel for HA
- Rate limiting middleware
- Metrics and monitoring (Prometheus)
- Distributed tracing (OpenTelemetry)
- API gateway integration

## Graceful Shutdown
1. Stop accepting new connections
2. Wait for in-flight requests (timeout: 30s)
3. Close Redis connections
4. Exit with appropriate status code

## Security Considerations
- Input validation (URL format, length limits)
- Rate limiting (future)
- No authentication (as per requirements)
- CORS headers (if needed for frontend)

## Performance Targets
- P99 latency < 50ms (redirect)
- P99 latency < 100ms (create/query)
- Throughput > 1000 req/s (single instance)
- Redis connection pooling (min: 10, max: 100)

## Future Enhancements (Out of Scope)
- Custom short codes
- Analytics dashboard
- Batch operations
- QR code generation
- Link preview
- Custom domains
- User authentication
- Multi-tenancy
