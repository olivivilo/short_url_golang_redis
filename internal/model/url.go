package model

import "time"

// URL represents a short URL entity with all its metadata.
type URL struct {
	Code        string     `json:"code"`
	ShortURL    string     `json:"short_url"`
	OriginalURL string     `json:"original_url"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpireAt    *time.Time `json:"expire_at,omitempty"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	VisitCount  int64      `json:"visit_count"`
	Note        string     `json:"note,omitempty"`
}

// CreateURLRequest represents the request body for creating a short URL.
type CreateURLRequest struct {
	URL      string `json:"url"`
	ExpireIn string `json:"expire_in,omitempty"` // Go duration format, e.g., "24h", "7d"
	Note     string `json:"note,omitempty"`
}

// IsExpired checks if the URL has expired.
func (u *URL) IsExpired() bool {
	if u.ExpireAt == nil {
		return false
	}
	return time.Now().After(*u.ExpireAt)
}

// IsDeleted checks if the URL has been soft-deleted.
func (u *URL) IsDeleted() bool {
	return u.DeletedAt != nil
}

// IsAccessible checks if the URL can be accessed (not expired and not deleted).
func (u *URL) IsAccessible() bool {
	return !u.IsExpired() && !u.IsDeleted()
}
