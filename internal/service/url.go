package service

import (
	"context"
	"fmt"

	"github.com/yourusername/short_url/internal/id"
	"github.com/yourusername/short_url/internal/model"
	"github.com/yourusername/short_url/internal/repository/redis"
)

// URLService handles business logic for URL operations.
type URLService struct {
	repo      *redis.URLRepository
	generator *id.Generator
	baseURL   string
}

// NewURLService creates a new URL service.
func NewURLService(repo *redis.URLRepository, generator *id.Generator, baseURL string) *URLService {
	return &URLService{
		repo:      repo,
		generator: generator,
		baseURL:   baseURL,
	}
}

// CreateURL creates a new short URL.
// TODO: Implement the create operation:
// 1. Validate the original URL format
// 2. Parse expire_in duration if provided
// 3. Generate a unique short code using the ID generator
// 4. Calculate expire_at timestamp if expiry is specified
// 5. Construct the full short URL
// 6. Save to repository
// 7. Return the created URL model
func (s *URLService) CreateURL(ctx context.Context, req *model.CreateURLRequest) (*model.URL, error) {
	// TODO: Validate URL format
	// if err := s.validateURL(req.URL); err != nil {
	//     return nil, fmt.Errorf("invalid URL: %w", err)
	// }

	// TODO: Parse expiry duration
	// var expireAt *time.Time
	// if req.ExpireIn != "" {
	//     duration, err := time.ParseDuration(req.ExpireIn)
	//     if err != nil {
	//         return nil, fmt.Errorf("invalid expire_in format: %w", err)
	//     }
	//     t := time.Now().Add(duration)
	//     expireAt = &t
	// }

	// TODO: Generate short code
	// code, err := s.generator.Generate(ctx)
	// if err != nil {
	//     return nil, fmt.Errorf("failed to generate code: %w", err)
	// }

	// TODO: Construct URL model
	// urlModel := &model.URL{
	//     Code:        code,
	//     ShortURL:    fmt.Sprintf("%s/r/%s", s.baseURL, code),
	//     OriginalURL: req.URL,
	//     CreatedAt:   time.Now(),
	//     ExpireAt:    expireAt,
	//     VisitCount:  0,
	//     Note:        req.Note,
	// }

	// TODO: Save to repository
	// if err := s.repo.Save(ctx, urlModel); err != nil {
	//     return nil, fmt.Errorf("failed to save URL: %w", err)
	// }

	return nil, fmt.Errorf("not implemented")
}

// GetURL retrieves a URL by its code.
// TODO: Implement the get operation:
// 1. Retrieve URL from repository
// 2. Check if URL exists
// 3. Return the URL model
func (s *URLService) GetURL(ctx context.Context, code string) (*model.URL, error) {
	// TODO: Implement get logic
	// urlModel, err := s.repo.Get(ctx, code)
	// if err != nil {
	//     return nil, err
	// }

	return nil, fmt.Errorf("not implemented")
}

// RedirectURL handles the redirect logic for a short code.
// TODO: Implement the redirect operation:
// 1. Retrieve URL from repository
// 2. Check if URL exists (return 404 if not)
// 3. Check if URL is expired (return 410 if expired)
// 4. Check if URL is deleted (return 410 if deleted)
// 5. Increment visit counter
// 6. Return the original URL for redirection
func (s *URLService) RedirectURL(ctx context.Context, code string) (string, error) {
	// TODO: Implement redirect logic
	// urlModel, err := s.repo.Get(ctx, code)
	// if err != nil {
	//     if err == redis.ErrURLNotFound {
	//         return "", ErrURLNotFound
	//     }
	//     return "", fmt.Errorf("failed to get URL: %w", err)
	// }

	// TODO: Check expiry and deletion status
	// if urlModel.IsExpired() {
	//     return "", ErrURLExpired
	// }
	// if urlModel.IsDeleted() {
	//     return "", ErrURLDeleted
	// }

	// TODO: Increment visit counter
	// if _, err := s.repo.IncrementVisits(ctx, code); err != nil {
	//     // Log error but don't fail the redirect
	//     // In production, you might want to use a proper logger here
	// }

	// return urlModel.OriginalURL, nil

	return "", fmt.Errorf("not implemented")
}

// DeleteURL soft-deletes a URL.
// TODO: Implement the delete operation:
// 1. Check if URL exists
// 2. Call repository delete method
// 3. Return appropriate error if URL doesn't exist
func (s *URLService) DeleteURL(ctx context.Context, code string) error {
	// TODO: Implement delete logic
	// if err := s.repo.Delete(ctx, code); err != nil {
	//     return err
	// }

	return fmt.Errorf("not implemented")
}

// validateURL validates if a string is a valid URL.
// TODO: Implement URL validation:
// 1. Check if URL is not empty
// 2. Parse URL using net/url package
// 3. Check if scheme is http or https
// 4. Check URL length limits
// 5. Optionally check for blacklisted domains
func (s *URLService) validateURL(urlStr string) error {
	if urlStr == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	// TODO: Implement validation logic
	// parsedURL, err := url.Parse(urlStr)
	// if err != nil {
	//     return fmt.Errorf("invalid URL format: %w", err)
	// }

	// TODO: Check scheme
	// if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
	//     return fmt.Errorf("URL must use http or https scheme")
	// }

	// TODO: Check length
	// if len(urlStr) > maxURLLength {
	//     return fmt.Errorf("URL exceeds maximum length")
	// }

	return fmt.Errorf("not implemented")
}

// Common service errors
var (
	ErrURLNotFound    = fmt.Errorf("URL not found")
	ErrURLExpired     = fmt.Errorf("URL has expired")
	ErrURLDeleted     = fmt.Errorf("URL has been deleted")
	ErrInvalidURL     = fmt.Errorf("invalid URL")
	ErrInvalidExpiry  = fmt.Errorf("invalid expiry format")
	ErrCodeGeneration = fmt.Errorf("failed to generate code")
)
