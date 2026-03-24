package http

import (
	"net/http"

	"github.com/yourusername/short_url/internal/http/handler"
)

// NewRouter creates and configures the HTTP router with all routes and middleware.
// TODO: Implement router setup:
// 1. Create a new ServeMux
// 2. Apply global middleware (recover, logging, request ID)
// 3. Register health check routes
// 4. Register API routes with /api/v1 prefix
// 5. Register redirect route /r/{code}
// 6. Consider adding CORS middleware if needed
func NewRouter(urlHandler *handler.URLHandler, healthHandler *handler.HealthHandler) http.Handler {
	mux := http.NewServeMux()

	// TODO: Register health check routes
	// mux.HandleFunc("/healthz", healthHandler.Healthz)
	// mux.HandleFunc("/readyz", healthHandler.Readyz)

	// TODO: Register redirect route (no /api prefix)
	// mux.HandleFunc("/r/", urlHandler.RedirectURL)

	// TODO: Register API routes with /api/v1 prefix
	// mux.HandleFunc("POST /api/v1/urls", urlHandler.CreateURL)
	// mux.HandleFunc("GET /api/v1/urls/", urlHandler.GetURL)
	// mux.HandleFunc("DELETE /api/v1/urls/", urlHandler.DeleteURL)

	// TODO: Apply middleware chain
	// The order matters: Recover -> RequestID -> Logging
	// var handler http.Handler = mux
	// handler = middleware.Logging(handler)
	// handler = middleware.RequestID(handler)
	// handler = middleware.Recover(handler)

	// return handler

	// Temporary: return basic mux
	return mux
}

// Note: With standard library net/http, you need to manually parse path parameters.
// For production, consider using a router library like chi, gorilla/mux, or httprouter
// that provides better path parameter extraction.
//
// Alternative approach with method-based routing:
// TODO: Implement method-based routing helper:
// func methodRouter(handlers map[string]http.HandlerFunc) http.HandlerFunc {
//     return func(w http.ResponseWriter, r *http.Request) {
//         handler, ok := handlers[r.Method]
//         if !ok {
//             http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
//             return
//         }
//         handler(w, r)
//     }
// }
