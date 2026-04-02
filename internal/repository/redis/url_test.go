package redis

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/yourusername/short_url/internal/id"
	"github.com/yourusername/short_url/internal/model"
)

func newTestRedis(t *testing.T) *redis.Client {
	t.Helper()

	addr := "localhost:6379"

	client := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   15,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		t.Fatalf("redis ping failed: %v", err)
	}

	if err := client.FlushDB(ctx).Err(); err != nil {
		t.Fatalf("redis flush before test failed: %v", err)
	}

	t.Cleanup(func() {
		clean_ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := client.FlushDB(clean_ctx).Err(); err != nil {
			t.Logf("redis flush after test failed: %v", err)
		}

		if err := client.Close(); err != nil {
			t.Logf("redis close failed: %v", err)
		}
	})

	return client
}

// TestSave tests saving a URL to Redis.
// DONE: Implement test cases:
// 1. Test saving a URL without expiry
// 2. Test saving a URL with expiry
// 3. Test saving a URL with note
// 4. Verify all fields are saved correctly
// 5. Verify visit counter is initialized to 0
// 6. Verify TTL is set correctly when expiry is specified
func TestSave(t *testing.T) {
	t.Run("save without expiry", func(t *testing.T) {
		client := newTestRedis(t)
		r := URLRepository{client: client}
		g := id.NewGenerator(r.client, 6)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		code, err := g.Generate(ctx)
		if err != nil {
			t.Fatalf("cannot generate code: %v", err)
		}

		now := time.Now().UTC().Truncate(time.Microsecond)
		url := &model.URL{
			Code:        code,
			ShortURL:    "http://short/" + code,
			OriginalURL: "https://example.com/no-expiry",
			CreatedAt:   now,
			VisitCount:  0,
			Note:        "",
		}

		if err := r.Save(ctx, url); err != nil {
			t.Fatalf("Save() error = %v", err)
		}

		key := fmt.Sprintf("url:%s", code)
		visitsKey := fmt.Sprintf("url:%s:visits", code)

		fields, err := client.HGetAll(ctx, key).Result()
		if err != nil {
			t.Fatalf("HGetAll() error = %v", err)
		}

		if got := fields["code"]; got != url.Code {
			t.Fatalf("code = %q, want %q", got, url.Code)
		}
		if got := fields["short_url"]; got != url.ShortURL {
			t.Fatalf("short_url = %q, want %q", got, url.ShortURL)
		}
		if got := fields["original_url"]; got != url.OriginalURL {
			t.Fatalf("original_url = %q, want %q", got, url.OriginalURL)
		}
		if got := fields["created_at"]; got != url.CreatedAt.Format(time.RFC3339Nano) {
			t.Fatalf("created_at = %q, want %q", got, url.CreatedAt.Format(time.RFC3339Nano))
		}
		if got := fields["note"]; got != url.Note {
			t.Fatalf("note = %q, want %q", got, url.Note)
		}
		if _, ok := fields["expire_at"]; ok {
			t.Fatalf("expire_at should not be set")
		}

		visits, err := client.Get(ctx, visitsKey).Int64()
		if err != nil {
			t.Fatalf("Get(visitsKey) error = %v", err)
		}
		if visits != 0 {
			t.Fatalf("visit counter = %d, want 0", visits)
		}

		ttl, err := client.TTL(ctx, key).Result()
		if err != nil {
			t.Fatalf("TTL(key) error = %v", err)
		}
		if ttl != -1 {
			t.Fatalf("TTL(key) = %v, want -1 (no expiry)", ttl)
		}
	})

	t.Run("save with expiry and note", func(t *testing.T) {
		client := newTestRedis(t)
		r := URLRepository{client: client}
		g := id.NewGenerator(r.client, 6)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		code, err := g.Generate(ctx)
		if err != nil {
			t.Fatalf("cannot generate code: %v", err)
		}

		now := time.Now().UTC().Truncate(time.Microsecond)
		expireAt := now.Add(10 * time.Minute)

		url := &model.URL{
			Code:        code,
			ShortURL:    "http://short/" + code,
			OriginalURL: "https://example.com/with-expiry",
			CreatedAt:   now,
			VisitCount:  0,
			Note:        "test note",
			ExpireAt:    &expireAt,
		}

		if err := r.Save(ctx, url); err != nil {
			t.Fatalf("Save() error = %v", err)
		}

		key := fmt.Sprintf("url:%s", code)
		visitsKey := fmt.Sprintf("url:%s:visits", code)

		fields, err := client.HGetAll(ctx, key).Result()
		if err != nil {
			t.Fatalf("HGetAll() error = %v", err)
		}

		if got := fields["code"]; got != url.Code {
			t.Fatalf("code = %q, want %q", got, url.Code)
		}
		if got := fields["short_url"]; got != url.ShortURL {
			t.Fatalf("short_url = %q, want %q", got, url.ShortURL)
		}
		if got := fields["original_url"]; got != url.OriginalURL {
			t.Fatalf("original_url = %q, want %q", got, url.OriginalURL)
		}
		if got := fields["created_at"]; got != url.CreatedAt.Format(time.RFC3339Nano) {
			t.Fatalf("created_at = %q, want %q", got, url.CreatedAt.Format(time.RFC3339Nano))
		}
		if got := fields["note"]; got != url.Note {
			t.Fatalf("note = %q, want %q", got, url.Note)
		}
		if got := fields["expire_at"]; got != url.ExpireAt.Format(time.RFC3339Nano) {
			t.Fatalf("expire_at = %q, want %q", got, url.ExpireAt.Format(time.RFC3339Nano))
		}

		visits, err := client.Get(ctx, visitsKey).Int64()
		if err != nil {
			t.Fatalf("Get(visitsKey) error = %v", err)
		}
		if visits != 0 {
			t.Fatalf("visit counter = %d, want 0", visits)
		}

		ttlKey, err := client.TTL(ctx, key).Result()
		if err != nil {
			t.Fatalf("TTL(key) error = %v", err)
		}
		ttlVisits, err := client.TTL(ctx, visitsKey).Result()
		if err != nil {
			t.Fatalf("TTL(visitsKey) error = %v", err)
		}

		if ttlKey <= 0 {
			t.Fatalf("TTL(key) = %v, want > 0", ttlKey)
		}
		if ttlVisits <= 0 {
			t.Fatalf("TTL(visitsKey) = %v, want > 0", ttlVisits)
		}

		expectedTTL := time.Until(expireAt)
		diff := ttlKey - expectedTTL
		if diff < 0 {
			diff = -diff
		}
		if diff > 2*time.Second {
			t.Fatalf("TTL(key) = %v, expected close to %v", ttlKey, expectedTTL)
		}
	})
}

// TestGet tests retrieving a URL from Redis.
// DONE: Implement test cases:
// 1. Test getting an existing URL
// 2. Test getting a non-existent URL (should return ErrURLNotFound)
// 3. Test getting a URL with all optional fields
// 4. Verify visit count is retrieved correctly
func TestGet(t *testing.T) {
	t.Run("get existing URL", func(t *testing.T) {
		client := newTestRedis(t)
		r := URLRepository{client: client}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		now := time.Now().UTC().Truncate(time.Microsecond)
		code := "abc123"

		url := &model.URL{
			Code:        code,
			ShortURL:    "http://short/" + code,
			OriginalURL: "https://example.com/basic",
			CreatedAt:   now,
			Note:        "",
		}

		if err := r.Save(ctx, url); err != nil {
			t.Fatalf("Save() error = %v", err)
		}

		got, err := r.Get(ctx, code)
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}

		if got.Code != url.Code {
			t.Fatalf("Code = %q, want %q", got.Code, url.Code)
		}
		if got.ShortURL != url.ShortURL {
			t.Fatalf("ShortURL = %q, want %q", got.ShortURL, url.ShortURL)
		}
		if got.OriginalURL != url.OriginalURL {
			t.Fatalf("OriginalURL = %q, want %q", got.OriginalURL, url.OriginalURL)
		}
		if !got.CreatedAt.Equal(url.CreatedAt) {
			t.Fatalf("CreatedAt = %v, want %v", got.CreatedAt, url.CreatedAt)
		}
		if got.VisitCount != 0 {
			t.Fatalf("VisitCount = %d, want 0", got.VisitCount)
		}
		if got.ExpireAt != nil {
			t.Fatalf("ExpireAt = %v, want nil", got.ExpireAt)
		}
		if got.DeletedAt != nil {
			t.Fatalf("DeletedAt = %v, want nil", got.DeletedAt)
		}
	})

	t.Run("get non-existent URL", func(t *testing.T) {
		client := newTestRedis(t)
		r := URLRepository{client: client}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		_, err := r.Get(ctx, "doesnotexist")
		if !errors.Is(err, ErrURLNotFound) {
			t.Fatalf("Get() error = %v, want ErrURLNotFound", err)
		}
	})

	t.Run("get URL with all optional fields", func(t *testing.T) {
		client := newTestRedis(t)
		r := URLRepository{client: client}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		now := time.Now().UTC().Truncate(time.Microsecond)
		expireAt := now.Add(30 * time.Minute)
		deletedAt := now.Add(5 * time.Minute)
		code := "full123"

		key := fmt.Sprintf("url:%s", code)
		visitsKey := fmt.Sprintf("url:%s:visits", code)

		fields := map[string]any{
			"code":         code,
			"short_url":    "http://short/" + code,
			"original_url": "https://example.com/full",
			"created_at":   now.Format(time.RFC3339Nano),
			"expire_at":    expireAt.Format(time.RFC3339Nano),
			"deleted_at":   deletedAt.Format(time.RFC3339Nano),
			"note":         "has optional fields",
		}

		if err := client.HSet(ctx, key, fields).Err(); err != nil {
			t.Fatalf("HSet() error = %v", err)
		}
		if err := client.Set(ctx, visitsKey, 7, 0).Err(); err != nil {
			t.Fatalf("Set() visits error = %v", err)
		}

		got, err := r.Get(ctx, code)
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}

		if got.Code != code {
			t.Fatalf("Code = %q, want %q", got.Code, code)
		}
		if got.Note != "has optional fields" {
			t.Fatalf("Note = %q, want %q", got.Note, "has optional fields")
		}
		if got.VisitCount != 7 {
			t.Fatalf("VisitCount = %d, want 7", got.VisitCount)
		}
		if got.ExpireAt == nil || !got.ExpireAt.Equal(expireAt) {
			t.Fatalf("ExpireAt = %v, want %v", got.ExpireAt, expireAt)
		}
		if got.DeletedAt == nil || !got.DeletedAt.Equal(deletedAt) {
			t.Fatalf("DeletedAt = %v, want %v", got.DeletedAt, deletedAt)
		}
	})

	t.Run("get retrieves visit count correctly", func(t *testing.T) {
		client := newTestRedis(t)
		r := URLRepository{client: client}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		now := time.Now().UTC().Truncate(time.Microsecond)
		code := "visits1"

		url := &model.URL{
			Code:        code,
			ShortURL:    "http://short/" + code,
			OriginalURL: "https://example.com/visits",
			CreatedAt:   now,
			Note:        "count test",
		}

		if err := r.Save(ctx, url); err != nil {
			t.Fatalf("Save() error = %v", err)
		}

		visitsKey := fmt.Sprintf("url:%s:visits", code)
		if err := client.Set(ctx, visitsKey, 42, 0).Err(); err != nil {
			t.Fatalf("Set() visits error = %v", err)
		}

		got, err := r.Get(ctx, code)
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}
		if got.VisitCount != 42 {
			t.Fatalf("VisitCount = %d, want 42", got.VisitCount)
		}
	})
}

// TestDelete tests soft-deleting a URL.
// DONE: Implement test cases:
// 1. Test deleting an existing URL
// 2. Test deleting a non-existent URL (should return error)
// 3. Verify deleted_at is set correctly
// 4. Verify URL can still be retrieved after deletion
func TestDelete(t *testing.T) {
	t.Run("delete existing URL", func(t *testing.T) {
		client := newTestRedis(t)
		r := URLRepository{client: client}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		now := time.Now().UTC().Truncate(time.Microsecond)
		code := "del123"

		url := &model.URL{
			Code:        code,
			ShortURL:    "http://short/" + code,
			OriginalURL: "https://example.com/delete",
			CreatedAt:   now,
			Note:        "",
		}

		if err := r.Save(ctx, url); err != nil {
			t.Fatalf("Save() error = %v", err)
		}

		if err := r.Delete(ctx, code); err != nil {
			t.Fatalf("Delete() error = %v", err)
		}

		key := fmt.Sprintf("url:%s", code)
		deletedAt, err := client.HGet(ctx, key, "deleted_at").Result()
		if err != nil {
			t.Fatalf("HGet(deleted_at) error = %v", err)
		}
		if deletedAt == "" {
			t.Fatalf("deleted_at is empty")
		}

		if _, err := time.Parse(time.RFC3339Nano, deletedAt); err != nil {
			t.Fatalf("deleted_at parse error = %v", err)
		}
	})

	t.Run("delete non-existent URL", func(t *testing.T) {
		client := newTestRedis(t)
		r := URLRepository{client: client}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		err := r.Delete(ctx, "missing123")
		if err == nil {
			t.Fatalf("Delete() error = nil, want non-nil")
		}
	})

	t.Run("URL can still be retrieved after deletion", func(t *testing.T) {
		client := newTestRedis(t)
		r := URLRepository{client: client}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		now := time.Now().UTC().Truncate(time.Microsecond)
		code := "soft123"

		url := &model.URL{
			Code:        code,
			ShortURL:    "http://short/" + code,
			OriginalURL: "https://example.com/soft-delete",
			CreatedAt:   now,
			Note:        "soft deleted",
		}

		if err := r.Save(ctx, url); err != nil {
			t.Fatalf("Save() error = %v", err)
		}

		if err := r.Delete(ctx, code); err != nil {
			t.Fatalf("Delete() error = %v", err)
		}

		got, err := r.Get(ctx, code)
		if err != nil {
			t.Fatalf("Get() after delete error = %v", err)
		}
		if got == nil {
			t.Fatalf("Get() after delete returned nil URL")
		}
		if got.DeletedAt == nil {
			t.Fatalf("DeletedAt = nil, want non-nil")
		}
	})
}

// TestIncrementVisits tests incrementing the visit counter.
// DONE: Implement test cases:
// 1. Test incrementing from 0 to 1
// 2. Test multiple increments
// 3. Test incrementing for non-existent URL
// 4. Verify returned count is correct
func TestIncrementVisits(t *testing.T) {
	t.Run("increment from 0 to 1", func(t *testing.T) {
		client := newTestRedis(t)
		r := URLRepository{client: client}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		now := time.Now().UTC().Truncate(time.Microsecond)
		code := "inc001"

		url := &model.URL{
			Code:        code,
			ShortURL:    "http://short/" + code,
			OriginalURL: "https://example.com/inc1",
			CreatedAt:   now,
		}

		if err := r.Save(ctx, url); err != nil {
			t.Fatalf("Save() error = %v", err)
		}

		got, err := r.IncrementVisits(ctx, code)
		if err != nil {
			t.Fatalf("IncrementVisits() error = %v", err)
		}
		if got != 1 {
			t.Fatalf("IncrementVisits() = %d, want 1", got)
		}
	})

	t.Run("multiple increments", func(t *testing.T) {
		client := newTestRedis(t)
		r := URLRepository{client: client}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		now := time.Now().UTC().Truncate(time.Microsecond)
		code := "inc002"

		url := &model.URL{
			Code:        code,
			ShortURL:    "http://short/" + code,
			OriginalURL: "https://example.com/inc2",
			CreatedAt:   now,
		}

		if err := r.Save(ctx, url); err != nil {
			t.Fatalf("Save() error = %v", err)
		}

		var got int64
		var err error
		for i := 0; i < 5; i++ {
			got, err = r.IncrementVisits(ctx, code)
			if err != nil {
				t.Fatalf("IncrementVisits() error on iteration %d = %v", i, err)
			}
		}

		if got != 5 {
			t.Fatalf("final IncrementVisits() = %d, want 5", got)
		}
	})

	t.Run("increment non-existent URL", func(t *testing.T) {
		client := newTestRedis(t)
		r := URLRepository{client: client}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		_, err := r.IncrementVisits(ctx, "missing999")
		if err == nil {
			t.Fatalf("IncrementVisits() error = nil, want non-nil")
		}
	})

	t.Run("returned count is correct", func(t *testing.T) {
		client := newTestRedis(t)
		r := URLRepository{client: client}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		now := time.Now().UTC().Truncate(time.Microsecond)
		code := "inc003"

		url := &model.URL{
			Code:        code,
			ShortURL:    "http://short/" + code,
			OriginalURL: "https://example.com/inc3",
			CreatedAt:   now,
		}

		if err := r.Save(ctx, url); err != nil {
			t.Fatalf("Save() error = %v", err)
		}

		want := int64(3)
		for i := int64(1); i <= want; i++ {
			got, err := r.IncrementVisits(ctx, code)
			if err != nil {
				t.Fatalf("IncrementVisits() error = %v", err)
			}
			if got != i {
				t.Fatalf("IncrementVisits() = %d, want %d", got, i)
			}
		}
	})
}

// TestExists tests checking URL existence.
// DONE: Implement test cases:
// 1. Test existing URL returns true
// 2. Test non-existent URL returns false
// 3. Test after deletion (should still return true)
func TestExists(t *testing.T) {
	t.Run("existing URL returns true", func(t *testing.T) {
		client := newTestRedis(t)
		r := URLRepository{client: client}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		now := time.Now().UTC().Truncate(time.Microsecond)
		code := "exists1"

		url := &model.URL{
			Code:        code,
			ShortURL:    "http://short/" + code,
			OriginalURL: "https://example.com/exists1",
			CreatedAt:   now,
		}

		if err := r.Save(ctx, url); err != nil {
			t.Fatalf("Save() error = %v", err)
		}

		got, err := r.Exists(ctx, code)
		if err != nil {
			t.Fatalf("Exists() error = %v", err)
		}
		if !got {
			t.Fatalf("Exists() = false, want true")
		}
	})

	t.Run("non-existent URL returns false", func(t *testing.T) {
		client := newTestRedis(t)
		r := URLRepository{client: client}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		got, err := r.Exists(ctx, "missing-exists")
		if err != nil {
			t.Fatalf("Exists() error = %v", err)
		}
		if got {
			t.Fatalf("Exists() = true, want false")
		}
	})

	t.Run("after deletion still returns true", func(t *testing.T) {
		client := newTestRedis(t)
		r := URLRepository{client: client}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		now := time.Now().UTC().Truncate(time.Microsecond)
		code := "exists2"

		url := &model.URL{
			Code:        code,
			ShortURL:    "http://short/" + code,
			OriginalURL: "https://example.com/exists2",
			CreatedAt:   now,
		}

		if err := r.Save(ctx, url); err != nil {
			t.Fatalf("Save() error = %v", err)
		}

		if err := r.Delete(ctx, code); err != nil {
			t.Fatalf("Delete() error = %v", err)
		}

		got, err := r.Exists(ctx, code)
		if err != nil {
			t.Fatalf("Exists() error = %v", err)
		}
		if !got {
			t.Fatalf("Exists() after delete = false, want true")
		}
	})
}
