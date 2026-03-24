package handler

import (
	"encoding/json"
	"net/http"

	"github.com/yourusername/short_url/internal/service"
)

// URLHandler handles HTTP requests for URL operations.
type URLHandler struct {
	service *service.URLService
}

// NewURLHandler creates a new URL handler.
func NewURLHandler(service *service.URLService) *URLHandler {
	return &URLHandler{
		service: service,
	}
}

// CreateURL handles POST /api/v1/urls
// TODO: Implement create URL handler:
// 1. Parse JSON request body
// 2. Validate request (check required fields)
// 3. Call service.CreateURL
// 4. Handle errors appropriately (400 for validation, 500 for server errors)
// 5. Return 201 Created with URL details
func (h *URLHandler) CreateURL(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement create logic
	// var req model.CreateURLRequest
	// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	//     respondError(w, http.StatusBadRequest, "Invalid request body")
	//     return
	// }

	// urlModel, err := h.service.CreateURL(r.Context(), &req)
	// if err != nil {
	//     // Map service errors to HTTP status codes
	//     respondError(w, http.StatusInternalServerError, "Failed to create URL")
	//     return
	// }

	// respondJSON(w, http.StatusCreated, urlModel)

	respondError(w, http.StatusNotImplemented, "Not implemented")
}

// GetURL handles GET /api/v1/urls/{code}
// TODO: Implement get URL handler:
// 1. Extract code from URL path
// 2. Call service.GetURL
// 3. Handle not found error (404)
// 4. Return 200 OK with URL details
func (h *URLHandler) GetURL(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement get logic
	// code := extractCodeFromPath(r.URL.Path)
	// if code == "" {
	//     respondError(w, http.StatusBadRequest, "Invalid code")
	//     return
	// }

	// urlModel, err := h.service.GetURL(r.Context(), code)
	// if err != nil {
	//     if errors.Is(err, redis.ErrURLNotFound) {
	//         respondError(w, http.StatusNotFound, "URL not found")
	//         return
	//     }
	//     respondError(w, http.StatusInternalServerError, "Failed to get URL")
	//     return
	// }

	// respondJSON(w, http.StatusOK, urlModel)

	respondError(w, http.StatusNotImplemented, "Not implemented")
}

// RedirectURL handles GET /r/{code}
// TODO: Implement redirect handler:
// 1. Extract code from URL path
// 2. Call service.RedirectURL
// 3. Handle errors:
//   - 404 for not found
//   - 410 for expired or deleted
//
// 4. Return 302 redirect to original URL
func (h *URLHandler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement redirect logic
	// code := extractCodeFromPath(r.URL.Path)
	// if code == "" {
	//     http.NotFound(w, r)
	//     return
	// }

	// originalURL, err := h.service.RedirectURL(r.Context(), code)
	// if err != nil {
	//     if errors.Is(err, service.ErrURLNotFound) {
	//         http.NotFound(w, r)
	//         return
	//     }
	//     if errors.Is(err, service.ErrURLExpired) || errors.Is(err, service.ErrURLDeleted) {
	//         http.Error(w, "Gone", http.StatusGone)
	//         return
	//     }
	//     http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	//     return
	// }

	// http.Redirect(w, r, originalURL, http.StatusFound)

	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// DeleteURL handles DELETE /api/v1/urls/{code}
// TODO: Implement delete URL handler:
// 1. Extract code from URL path
// 2. Call service.DeleteURL
// 3. Handle not found error (404)
// 4. Return 204 No Content on success
func (h *URLHandler) DeleteURL(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement delete logic
	// code := extractCodeFromPath(r.URL.Path)
	// if code == "" {
	//     respondError(w, http.StatusBadRequest, "Invalid code")
	//     return
	// }

	// if err := h.service.DeleteURL(r.Context(), code); err != nil {
	//     if errors.Is(err, redis.ErrURLNotFound) {
	//         respondError(w, http.StatusNotFound, "URL not found")
	//         return
	//     }
	//     respondError(w, http.StatusInternalServerError, "Failed to delete URL")
	//     return
	// }

	// w.WriteHeader(http.StatusNoContent)

	respondError(w, http.StatusNotImplemented, "Not implemented")
}

// extractCodeFromPath extracts the code parameter from the URL path.
// TODO: Implement path parameter extraction:
// 1. Parse the URL path
// 2. Extract the last segment as the code
// 3. Handle edge cases (empty path, trailing slashes)
func extractCodeFromPath(path string) string {
	// TODO: Implement extraction logic
	// This is a simple implementation; consider using a router library for production
	return ""
}

// respondJSON sends a JSON response.
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Log error in production
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// respondError sends a JSON error response.
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
