package middleware

import (
	"net/http"
)

// responseWriter wraps http.ResponseWriter to capture status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.written += n
	return n, err
}

// Logging logs HTTP requests with method, path, status, duration, and size.
// TODO: Implement logging middleware:
// 1. Wrap the response writer to capture status code and bytes written
// 2. Record start time
// 3. Call next handler
// 4. Calculate duration
// 5. Log request details (method, path, status, duration, size)
// 6. Consider adding request ID from context
// 7. In production, use a structured logger (e.g., slog, zap, zerolog)
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement logging logic
		// start := time.Now()
		// wrapped := newResponseWriter(w)

		// next.ServeHTTP(wrapped, r)

		// duration := time.Since(start)
		// log.Printf(
		//     "%s %s %d %v %d bytes",
		//     r.Method,
		//     r.URL.Path,
		//     wrapped.statusCode,
		//     duration,
		//     wrapped.written,
		// )

		// Temporary pass-through until implemented
		next.ServeHTTP(w, r)
	})
}

// LoggingWithRequestID logs HTTP requests including the request ID from context.
// TODO: Implement enhanced logging with request ID:
// 1. Extract request ID from context
// 2. Include request ID in log output
// 3. Log additional fields like user agent, remote addr
func LoggingWithRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement enhanced logging
		// requestID := GetRequestID(r.Context())
		// start := time.Now()
		// wrapped := newResponseWriter(w)

		// next.ServeHTTP(wrapped, r)

		// duration := time.Since(start)
		// log.Printf(
		//     "[%s] %s %s %d %v %d bytes - %s",
		//     requestID,
		//     r.Method,
		//     r.URL.Path,
		//     wrapped.statusCode,
		//     duration,
		//     wrapped.written,
		//     r.UserAgent(),
		// )

		// Temporary pass-through until implemented
		next.ServeHTTP(w, r)
	})
}
