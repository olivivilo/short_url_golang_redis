package handler

import (
	"net/http"

	"github.com/redis/go-redis/v9"
)

// HealthHandler handles health check requests.
type HealthHandler struct {
	redisClient *redis.Client
}

// NewHealthHandler creates a new health handler.
func NewHealthHandler(redisClient *redis.Client) *HealthHandler {
	return &HealthHandler{
		redisClient: redisClient,
	}
}

// Healthz handles GET /healthz - basic liveness check.
// TODO: Implement liveness check:
// 1. Return 200 OK with simple status message
// 2. This should always succeed if the process is running
// 3. Don't check external dependencies here
func (h *HealthHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement liveness check
	// respondJSON(w, http.StatusOK, map[string]string{
	//     "status": "ok",
	//     "timestamp": time.Now().Format(time.RFC3339),
	// })

	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Readyz handles GET /readyz - readiness check with Redis ping.
// TODO: Implement readiness check:
// 1. Ping Redis with a timeout
// 2. Return 200 OK if Redis is reachable
// 3. Return 503 Service Unavailable if Redis is down
// 4. Include Redis status in response
func (h *HealthHandler) Readyz(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement readiness check
	// ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	// defer cancel()

	// if err := h.redisClient.Ping(ctx).Err(); err != nil {
	//     respondJSON(w, http.StatusServiceUnavailable, map[string]interface{}{
	//         "status": "unavailable",
	//         "redis": "down",
	//         "error": err.Error(),
	//         "timestamp": time.Now().Format(time.RFC3339),
	//     })
	//     return
	// }

	// respondJSON(w, http.StatusOK, map[string]interface{}{
	//     "status": "ready",
	//     "redis": "up",
	//     "timestamp": time.Now().Format(time.RFC3339),
	// })

	respondJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}
