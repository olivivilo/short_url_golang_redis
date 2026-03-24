package middleware

import (
	"context"
	"net/http"
)

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

const requestIDKey contextKey = "request_id"

// RequestID adds a unique request ID to each request context.
// TODO: Implement request ID middleware:
// 1. Check if X-Request-ID header is present
// 2. If present, use it; otherwise generate a new UUID
// 3. Add request ID to context
// 4. Set X-Request-ID response header
// 5. Call next handler with updated context
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement request ID logic
		// requestID := r.Header.Get("X-Request-ID")
		// if requestID == "" {
		//     requestID = uuid.New().String()
		// }

		// ctx := context.WithValue(r.Context(), requestIDKey, requestID)
		// w.Header().Set("X-Request-ID", requestID)

		// next.ServeHTTP(w, r.WithContext(ctx))

		// Temporary pass-through until implemented
		next.ServeHTTP(w, r)
	})
}

// GetRequestID retrieves the request ID from the context.
// TODO: Implement request ID retrieval:
// 1. Get value from context
// 2. Type assert to string
// 3. Return empty string if not found
func GetRequestID(ctx context.Context) string {
	// TODO: Implement retrieval logic
	// if requestID, ok := ctx.Value(requestIDKey).(string); ok {
	//     return requestID
	// }
	return ""
}
