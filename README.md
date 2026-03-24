# Short URL Service

A production-ready short URL service built with Go and Redis, following enterprise best practices.

## Features

- ✅ Create short URLs with optional expiration
- ✅ Redirect to original URLs with visit tracking
- ✅ Query URL details and statistics
- ✅ Soft delete URLs
- ✅ Health checks (liveness and readiness)
- ✅ Graceful shutdown
- ✅ Request ID tracking
- ✅ Structured logging
- ✅ Panic recovery
- ✅ Comprehensive testing (unit + integration)

## Architecture

See [`ARCHITECTURE.md`](ARCHITECTURE.md) for detailed design decisions, data models, and API specifications.

## Quick Start

### Prerequisites

- Go 1.22 or higher
- Redis 6.0 or higher
- Docker and Docker Compose (optional, for containerized setup)

### Local Development

1. **Clone the repository**
   ```bash
   cd /root/Prjs/Tests/short_url
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Start Redis**
   ```bash
   # Using Docker
   docker run -d -p 6379:6379 --name redis redis:7-alpine
   
   # Or using Docker Compose
   docker-compose up -d redis
   ```

4. **Set environment variables** (optional)
   ```bash
   export SERVER_PORT=8080
   export REDIS_ADDR=localhost:6379
   export BASE_URL=http://localhost:8080
   export LOG_LEVEL=info
   ```

5. **Run the service**
   ```bash
   go run cmd/shorturl/main.go
   ```

6. **Test the service**
   ```bash
   # Health check
   curl http://localhost:8080/healthz
   
   # Create a short URL
   curl -X POST http://localhost:8080/api/v1/urls \
     -H "Content-Type: application/json" \
     -d '{"url":"https://example.com","expire_in":"24h"}'
   
   # Redirect (replace {code} with actual code from above)
   curl -L http://localhost:8080/r/{code}
   ```

## API Reference

### Create Short URL

```bash
POST /api/v1/urls
Content-Type: application/json

{
  "url": "https://example.com/very/long/path",
  "expire_in": "24h",  # Optional: Go duration format (e.g., "1h", "7d", "30m")
  "note": "Campaign link"  # Optional
}

Response 201 Created:
{
  "code": "000001",
  "short_url": "http://localhost:8080/r/000001",
  "original_url": "https://example.com/very/long/path",
  "created_at": "2026-03-23T08:00:00Z",
  "expire_at": "2026-03-24T08:00:00Z",
  "deleted_at": null,
  "visit_count": 0,
  "note": "Campaign link"
}
```

### Redirect to Original URL

```bash
GET /r/{code}

Response 302 Found:
Location: https://example.com/very/long/path

Response 404 Not Found:
URL not found

Response 410 Gone:
URL expired or deleted
```

### Get URL Details

```bash
GET /api/v1/urls/{code}

Response 200 OK:
{
  "code": "000001",
  "short_url": "http://localhost:8080/r/000001",
  "original_url": "https://example.com/very/long/path",
  "created_at": "2026-03-23T08:00:00Z",
  "expire_at": "2026-03-24T08:00:00Z",
  "deleted_at": null,
  "visit_count": 42,
  "note": "Campaign link"
}

Response 404 Not Found:
{
  "error": "URL not found"
}
```

### Delete Short URL

```bash
DELETE /api/v1/urls/{code}

Response 204 No Content

Response 404 Not Found:
{
  "error": "URL not found"
}
```

### Health Checks

```bash
# Liveness check (always returns 200 if process is running)
GET /healthz

Response 200 OK:
{
  "status": "ok"
}

# Readiness check (checks Redis connectivity)
GET /readyz

Response 200 OK:
{
  "status": "ready"
}

Response 503 Service Unavailable:
{
  "status": "unavailable"
}
```

## Configuration

Configuration is loaded from environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | HTTP server port | `8080` |
| `REDIS_ADDR` | Redis address | `localhost:6379` |
| `REDIS_PASSWORD` | Redis password | `` |
| `REDIS_DB` | Redis database number | `0` |
| `REDIS_POOL_SIZE` | Redis connection pool size | `50` |
| `BASE_URL` | Base URL for short links | `http://localhost:8080` |
| `LOG_LEVEL` | Logging level | `info` |
| `MIN_CODE_LENGTH` | Minimum short code length | `6` |
| `DEFAULT_TTL` | Default TTL in seconds (0 = no expiry) | `0` |
| `MAX_URL_LENGTH` | Maximum URL length | `2048` |
| `SHUTDOWN_TIMEOUT` | Graceful shutdown timeout (seconds) | `30` |

## Testing

### Run Unit Tests

```bash
# Run all unit tests
go test ./internal/... -v

# Run specific package tests
go test ./internal/id -v
go test ./internal/service -v
go test ./internal/repository/redis -v

# Run with coverage
go test ./internal/... -cover
go test ./internal/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Run Integration Tests

```bash
# Start Redis for testing
docker-compose up -d redis

# Run integration tests
go test ./test/integration/... -v

# Set custom Redis address for tests
TEST_REDIS_ADDR=localhost:6379 go test ./test/integration/... -v
```

### Run All Tests

```bash
# Run all tests (unit + integration)
go test ./... -v

# Run with race detector
go test ./... -race

# Run with coverage
go test ./... -cover -coverprofile=coverage.out
```

## Docker Deployment

### Build Docker Image

```bash
docker build -t short-url-service:latest -f deploy/Dockerfile .
```

### Run with Docker Compose

```bash
# Start all services (Redis + App)
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Run Docker Container Manually

```bash
# Start Redis
docker run -d --name redis -p 6379:6379 redis:7-alpine

# Start the service
docker run -d \
  --name short-url \
  -p 8080:8080 \
  -e REDIS_ADDR=redis:6379 \
  -e BASE_URL=http://localhost:8080 \
  --link redis \
  short-url-service:latest
```

## Development Workflow

### Project Structure

```
short_url/
├── cmd/shorturl/          # Application entry point
├── internal/              # Private application code
│   ├── config/           # Configuration management
│   ├── http/             # HTTP layer (handlers, middleware, router)
│   ├── service/          # Business logic
│   ├── repository/       # Data access layer
│   ├── id/               # ID generation (Base62)
│   └── model/            # Domain models
├── test/                  # Integration tests
│   ├── integration/      # Integration test suites
│   └── testutil/         # Test utilities
├── deploy/                # Deployment files
│   ├── Dockerfile        # Container image
│   └── docker-compose.yml # Local stack
├── configs/               # Configuration examples
├── .vscode/               # VSCode settings
├── go.mod                 # Go module definition
├── README.md              # This file
└── ARCHITECTURE.md        # Architecture documentation
```

### Adding New Features

1. **Define the model** in `internal/model/`
2. **Add repository methods** in `internal/repository/redis/`
3. **Implement business logic** in `internal/service/`
4. **Add HTTP handlers** in `internal/http/handler/`
5. **Register routes** in `internal/http/router.go`
6. **Write tests** (unit + integration)
7. **Update documentation**

### Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Use `golint` for linting
- Write clear comments and TODOs
- Keep functions small and focused
- Use meaningful variable names

### Debugging

Use VSCode launch configurations (see `.vscode/launch.json`):

- **Launch Server**: Start the service with debugger attached
- **Debug Tests**: Run tests with debugger
- **Attach to Process**: Attach to running process

## TODO: Implementation Tasks

This project provides a complete skeleton with TODOs for you to implement. Here are the main areas:

### Core Implementation
- [ ] Complete Base62 encoding/decoding in `internal/id/generator.go`
- [ ] Implement Redis operations in `internal/repository/redis/url.go`
- [ ] Implement business logic in `internal/service/url.go`
- [ ] Complete HTTP handlers in `internal/http/handler/`
- [ ] Implement middleware in `internal/http/middleware/`
- [ ] Wire up router in `internal/http/router.go`
- [ ] Complete main function in `cmd/shorturl/main.go`

### Testing
- [ ] Write unit tests for ID generator
- [ ] Write unit tests for repository layer
- [ ] Write unit tests for service layer
- [ ] Write integration tests for full workflows
- [ ] Write HTTP endpoint tests

### Configuration & Deployment
- [ ] Implement configuration validation
- [ ] Complete Docker Compose setup
- [ ] Add environment-specific configs

## Troubleshooting

### Redis Connection Issues

```bash
# Check if Redis is running
docker ps | grep redis

# Test Redis connection
redis-cli ping

# View Redis logs
docker logs redis
```

### Port Already in Use

```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>

# Or use a different port
export SERVER_PORT=8081
```

### Module Import Errors

```bash
# Clean module cache
go clean -modcache

# Re-download dependencies
go mod download

# Verify dependencies
go mod verify
```

## Performance Considerations

- **Redis Connection Pooling**: Configured via `REDIS_POOL_SIZE`
- **HTTP Timeouts**: Read/Write/Idle timeouts configured in server
- **Graceful Shutdown**: Ensures in-flight requests complete
- **Visit Counter**: Separate key for efficient increments
- **TTL Management**: Automatic cleanup via Redis expiration

## Security Considerations

- **Input Validation**: URL format and length validation
- **No Authentication**: As per requirements (add if needed)
- **Rate Limiting**: Not implemented (add if needed)
- **HTTPS**: Use reverse proxy (nginx, Caddy) in production

## License

This is a learning project. Use it as you wish.

## Contributing

This is a personal learning project. Feel free to fork and modify for your own learning.

## Support

For questions or issues, refer to:
- [`ARCHITECTURE.md`](ARCHITECTURE.md) for design details
- Code comments and TODOs for implementation guidance
- Go documentation: https://golang.org/doc/
- Redis documentation: https://redis.io/documentation
