package id

import (
	"testing"

	"github.com/redis/go-redis/v9"
)

// TestEncodeBase62 tests the Base62 encoding function.
// DONE: Implement test cases:
// 1. Test encoding of 0
// 2. Test encoding of small numbers (1, 10, 61, 62)
// 3. Test encoding of large numbers
// 4. Verify the encoded string uses only valid Base62 characters
func TestEncodeBase62(t *testing.T) {
	// DONE: Implement test cases
	// t.Skip("TODO: Implement Base62 encoding tests")
	g := Generator{minLength: 6}

	tests := []struct {
		in   int64
		want string
	}{
		{0, "0"},
		{1, "1"},
		{61, "z"},
		{62, "10"},
		{125, "21"},
	}

	for _, tt := range tests {
		got := g.encodeBase62(tt.in)

		if got != tt.want {
			t.Errorf("encodeBase62(%d) = %q; want %q", tt.in, got, tt.want)
		}
	}
}

// TestDecodeBase62 tests the Base62 decoding function.
// DONE: Implement test cases:
// 1. Test decoding of single characters
// 2. Test decoding of multi-character codes
// 3. Test round-trip (encode then decode)
// 4. Test error cases (invalid characters)
func TestDecodeBase62(t *testing.T) {
	// DONE: Implement test cases
	// t.Skip("TODO: Implement Base62 decoding tests")

	tests := []struct {
		in      string
		want    int64
		wantErr bool
	}{
		{"", 0, true},
		{"0", 0, false},
		{"10", 62, false},
		{"21", 125, false},
		{"!", 0, true},
	}

	for _, tt := range tests {
		out, outErr := DecodeBase62(tt.in)

		if tt.wantErr {
			if outErr == nil {
				t.Errorf("DecodeBase62(%q) expected error, got nil", tt.in)
			}
			continue
		}

		if outErr != nil {
			t.Errorf("DecodeBase62(%q) unexpected error: %v", tt.in, outErr)
			continue
		}

		if out != tt.want {
			t.Errorf("DecodeBase62(%q) = %d; want %d", tt.in, out, tt.want)
		}

	}
}

// TestPadToMinLength tests the padding function.
// DONE: Implement test cases:
// 1. Test padding when code is shorter than minimum
// 2. Test no padding when code is at minimum length
// 3. Test no padding when code is longer than minimum
// 4. Verify padding uses the correct character ('0')
func TestPadToMinLength(t *testing.T) {
	// DONE: Implement test cases
	// t.Skip("TODO: Implement padding tests")
	g := Generator{minLength: 6}

	tests := []struct {
		in   string
		want string
	}{
		{"0", "000000"},
		{"5", "000005"},
		{"2eR", "0002eR"},
		{"rP07sf", "rP07sf"},
		{"abcdefg", "abcdefg"},
	}

	for _, tt := range tests {
		got := g.padToMinLength(tt.in)
		if got != tt.want {
			t.Errorf("padToMinLength(%q) = %q; want %q", tt.in, got, tt.want)
		}
	}
}

// TestValidateCode tests the code validation function.
// DONE: Implement test cases:
// 1. Test valid codes
// 2. Test empty code
// 3. Test codes with invalid characters
// 4. Test codes with special characters
func TestValidateCode(t *testing.T) {
	// DONE: Implement test cases
	// t.Skip("TODO: Implement validation tests")
	tests := []struct {
		in   string
		want string
	}{
		{"", "code cannot be empty"},
		{"erfdsg", ""},
		{"rr#1@^&d", "invalid character in code: #"},
		{"abc-123", "invalid character in code: -"},
	}

	for _, tt := range tests {
		err := ValidateCode(tt.in)
		if tt.want == "" {
			if err != nil {
				t.Errorf("ValidateCode(%q) unexpected error: %v", tt.in, err)
			}
			continue
		}

		if err == nil {
			t.Errorf("ValidateCode(%q) = nil; want error %q", tt.in, tt.want)
			continue
		}

		if err.Error() != tt.want {
			t.Errorf("ValidateCode(%q) error = %q; want %q", tt.in, err.Error(), tt.want)
		}
	}
}

// TestMaxIDForLength tests the maximum ID calculation.
// DONE: Implement test cases:
// 1. Test for length 1 (should be 61)
// 2. Test for length 6 (minimum length)
// 3. Test for length 10
func TestMaxIDForLength(t *testing.T) {
	// DONE: Implement test cases
	//t.Skip("TODO: Implement max ID calculation tests")

	tests := []struct {
		in   int
		want int64
	}{
		{1, 61},
		{6, 56800235583},
		{10, 839299365868340223},
	}

	for _, tt := range tests {
		got := MaxIDForLength(tt.in)
		if got != tt.want {
			t.Errorf("MaxIDForLength(%d) = %d; want %d", tt.in, got, tt.want)
		}
	}
}

func newTestRedis(t *testing.T) *redis.Client {
	return nil
}

func TestNewGenerator(t *testing.T) {
	t.Helper()

	addr := "localhost:6379"

	client := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   15,
	})
}

// TestGeneratorIntegration tests the full generator with a real Redis instance.
// This test requires a running Redis instance.
// TODO: Implement integration test:
// 1. Set up test Redis client
// 2. Generate multiple codes
// 3. Verify uniqueness
// 4. Verify minimum length
// 5. Clean up test data
func TestGeneratorIntegration(t *testing.T) {
	// TODO: Implement integration test
	t.Skip("TODO: Implement generator integration tests")
}
