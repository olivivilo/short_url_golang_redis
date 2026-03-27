package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yourusername/short_url/internal/model"
)

// URLRepository handles Redis operations for URL entities.
type URLRepository struct {
	client *redis.Client
}

// NewURLRepository creates a new URL repository.
func NewURLRepository(client *redis.Client) *URLRepository {
	return &URLRepository{
		client: client,
	}
}

// Save saves a URL to Redis.
// DONE: Implement the save operation:
// 1. Create the Redis key: "url:{code}"
// 2. Use HSET to store all URL fields in a hash
// 3. Set TTL if expire_at is specified
// 4. Initialize visit counter: "url:{code}:visits" to 0
// 5. Set TTL on visit counter to match URL TTL
// 6. Use pipeline or transaction for atomicity
func (r *URLRepository) Save(ctx context.Context, url *model.URL) error {
	// DONE: Implement save logic
	// key := fmt.Sprintf("url:%s", url.Code)
	// visitsKey := fmt.Sprintf("url:%s:visits", url.Code)

	key := fmt.Sprintf("url:%s", url.Code)

	visitsKey := fmt.Sprintf("url:%s:visits", url.Code)

	values := map[string]any{
		"code":         url.Code,
		"short_url":    url.ShortURL,
		"original_url": url.OriginalURL,
		"created_at":   url.CreatedAt.Format(time.RFC3339Nano),
		"note":         url.Note,
	}

	if url.ExpireAt != nil {
		values["expire_at"] = url.ExpireAt.Format(time.RFC3339Nano)
	}

	if url.DeletedAt != nil {
		values["deleted_at"] = url.DeletedAt.Format(time.RFC3339Nano)
	}

	_, err := r.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, key, values)
		pipe.Set(ctx, visitsKey, 0, 0)

		if url.ExpireAt != nil {
			ttl := time.Until(*url.ExpireAt)
			if ttl > 0 {
				pipe.Expire(ctx, key, ttl)
				pipe.Expire(ctx, visitsKey, ttl)
			} else {
				pipe.Del(ctx, key, visitsKey)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("Save not implemented due to: %v", err)
	}

	// DONE: Use Redis pipeline for atomic operations
	// pipe := r.client.Pipeline()
	// ... pipeline operations ...
	// _, err := pipe.Exec(ctx)

	return nil
}

// Get retrieves a URL by its code.
// DONE: Implement the get operation:
// 1. Create the Redis key: "url:{code}"
// 2. Use HGETALL to retrieve all fields
// 3. Check if key exists (empty map means not found)
// 4. Parse timestamps from string to time.Time
// 5. Get visit count from "url:{code}:visits"
// 6. Construct and return the URL model
func (r *URLRepository) Get(ctx context.Context, code string) (*model.URL, error) {
	// DONE: Implement get logic
	// key := fmt.Sprintf("url:%s", code)
	// visitsKey := fmt.Sprintf("url:%s:visits", code)
	key := fmt.Sprintf("url:%s", code)
	visitsKey := fmt.Sprintf("url:%s:visits", code)

	// DONE: Use HGETALL to get all fields
	// fields, err := r.client.HGetAll(ctx, key).Result()
	// if err != nil {
	//     return nil, fmt.Errorf("failed to get URL: %w", err)
	// }
	fields, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get URL hash for %q: %w", code, err)
	}

	// DONE: Check if URL exists
	// if len(fields) == 0 {
	//     return nil, ErrURLNotFound
	// }
	if len(fields) == 0 {
		return nil, ErrURLNotFound
	}

	// DONE: Parse fields and construct URL model
	// DONE: Get visit count
	url := model.URL{
		Code:        fields["code"],
		ShortURL:    fields["short_url"],
		OriginalURL: fields["original_url"],
		Note:        fields["note"],
	}

	createdAt, err := time.Parse(time.RFC3339Nano, fields["created_at"])
	if err != nil {
		return nil, fmt.Errorf("invalid created_at %q: %w", fields["created_at"], err)
	}
	url.CreatedAt = createdAt

	if s, ok := fields["expire_at"]; ok && s != "" {
		t, err := time.Parse(time.RFC3339Nano, s)
		if err != nil {
			return nil, fmt.Errorf("invalid expire_at %q: %w", s, err)
		}
		url.ExpireAt = &t
	}

	if s, ok := fields["deleted_at"]; ok && s != "" {
		t, err := time.Parse(time.RFC3339Nano, s)
		if err != nil {
			return nil, fmt.Errorf("invalid deleted_at %q: %w", s, err)
		}
		url.DeletedAt = &t
	}

	visits, err := r.client.Get(ctx, visitsKey).Result()
	if err != nil {
		if err == redis.Nil {
			url.VisitCount = 0
		} else {
			return nil, fmt.Errorf("failed to get visits for %q: %w", code, err)
		}
	} else {
		n, err := strconv.ParseInt(visits, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid visits value %q: %w", visits, err)
		}
		url.VisitCount = n
	}

	return &url, nil
}

// Delete soft-deletes a URL by setting the deleted_at timestamp.
// TODO: Implement the delete operation:
// 1. Check if URL exists
// 2. Set "deleted_at" field in the hash
// 3. Optionally adjust TTL to expire soon (e.g., 7 days)
// 4. Return error if URL doesn't exist
func (r *URLRepository) Delete(ctx context.Context, code string) error {
	// TODO: Implement delete logic
	// key := fmt.Sprintf("url:%s", code)

	// TODO: Check existence
	// exists, err := r.client.Exists(ctx, key).Result()
	// if err != nil {
	//     return fmt.Errorf("failed to check URL existence: %w", err)
	// }
	// if exists == 0 {
	//     return ErrURLNotFound
	// }

	// TODO: Set deleted_at timestamp
	// now := time.Now().Format(time.RFC3339)
	// err = r.client.HSet(ctx, key, "deleted_at", now).Err()

	return fmt.Errorf("not implemented")
}

// IncrementVisits increments the visit counter for a URL.
// TODO: Implement the increment operation:
// 1. Use INCR on "url:{code}:visits"
// 2. Return the new count
// 3. Handle the case where the key doesn't exist
func (r *URLRepository) IncrementVisits(ctx context.Context, code string) (int64, error) {
	// TODO: Implement increment logic
	// visitsKey := fmt.Sprintf("url:%s:visits", code)
	// count, err := r.client.Incr(ctx, visitsKey).Result()
	// if err != nil {
	//     return 0, fmt.Errorf("failed to increment visits: %w", err)
	// }

	return 0, fmt.Errorf("not implemented")
}

// Exists checks if a URL code exists in Redis.
// TODO: Implement the exists check:
// 1. Use EXISTS command on "url:{code}"
// 2. Return true if exists, false otherwise
func (r *URLRepository) Exists(ctx context.Context, code string) (bool, error) {
	// TODO: Implement exists check
	// key := fmt.Sprintf("url:%s", code)
	// count, err := r.client.Exists(ctx, key).Result()
	// if err != nil {
	//     return false, fmt.Errorf("failed to check existence: %w", err)
	// }
	// return count > 0, nil

	return false, fmt.Errorf("not implemented")
}

// parseTime parses a time string in RFC3339 format.
// Returns nil if the string is empty.
func parseTime(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil, fmt.Errorf("failed to parse time: %w", err)
	}
	return &t, nil
}

// formatTime formats a time pointer to RFC3339 string.
// Returns empty string if the pointer is nil.
func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}

// calculateTTL calculates the TTL duration from an expiry time.
// Returns 0 if expireAt is nil (no expiry).
func calculateTTL(expireAt *time.Time) time.Duration {
	if expireAt == nil {
		return 0
	}
	ttl := time.Until(*expireAt)
	if ttl < 0 {
		return 0
	}
	return ttl
}

// Common errors
var (
	ErrURLNotFound = fmt.Errorf("URL not found")
	ErrURLExpired  = fmt.Errorf("URL has expired")
	ErrURLDeleted  = fmt.Errorf("URL has been deleted")
)
