package id

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/redis/go-redis/v9"
)

const (
	// Base62 alphabet: 0-9, A-Z, a-z
	base62Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	base           = 62
)

// Generator generates unique short codes using auto-increment ID + Base62 encoding.
type Generator struct {
	redis     *redis.Client
	minLength int
}

// NewGenerator creates a new ID generator.
func NewGenerator(redisClient *redis.Client, minLength int) *Generator {
	return &Generator{
		redis:     redisClient,
		minLength: minLength,
	}
}

// Generate generates a new unique short code.
// It uses Redis INCR to get a unique ID, then encodes it to Base62.
// TODO: Implement the ID generation logic:
// 1. Use Redis INCR on "global:url_id" to get a unique ID
// 2. Encode the ID to Base62
// 3. Pad to minimum length if necessary
// 4. Handle Redis errors appropriately
func (g *Generator) Generate(ctx context.Context) (string, error) {
	// TODO: Implement Redis INCR operation
	// id, err := g.redis.Incr(ctx, "global:url_id").Result()
	// if err != nil {
	//     return "", fmt.Errorf("failed to generate ID: %w", err)
	// }

	// TODO: Encode ID to Base62
	// code := g.encodeBase62(id)

	// TODO: Pad to minimum length
	// code = g.padToMinLength(code)

	return "", fmt.Errorf("not implemented")
}

// encodeBase62 encodes a number to Base62 string.
// TODO: Implement Base62 encoding algorithm:
// 1. Handle zero case
// 2. Convert number to base 62 using the alphabet
// 3. Reverse the result string
func (g *Generator) encodeBase62(num int64) string {
	if num == 0 {
		return string(base62Alphabet[0])
	}

	// TODO: Implement encoding logic
	var result strings.Builder
	// ... encoding logic here ...

	return result.String()
}

// DecodeBase62 decodes a Base62 string back to a number.
// This is useful for testing and debugging.
// TODO: Implement Base62 decoding algorithm:
// 1. Iterate through each character
// 2. Find its position in the alphabet
// 3. Calculate the numeric value
func DecodeBase62(code string) (int64, error) {
	// TODO: Implement decoding logic
	var result int64
	// ... decoding logic here ...

	return result, fmt.Errorf("not implemented")
}

// padToMinLength pads the code to the minimum length by prepending '0' characters.
// TODO: Implement padding logic:
// 1. Check if code is already at or above minimum length
// 2. Calculate padding needed
// 3. Prepend '0' characters (first character in alphabet)
func (g *Generator) padToMinLength(code string) string {
	if len(code) >= g.minLength {
		return code
	}

	// TODO: Implement padding
	padding := g.minLength - len(code)
	return strings.Repeat(string(base62Alphabet[0]), padding) + code
}

// ValidateCode validates if a code is a valid Base62 string.
// TODO: Implement validation logic:
// 1. Check if code is empty
// 2. Check if all characters are in the Base62 alphabet
// 3. Check length constraints
func ValidateCode(code string) error {
	if code == "" {
		return fmt.Errorf("code cannot be empty")
	}

	// TODO: Implement character validation
	// for _, ch := range code {
	//     if !strings.ContainsRune(base62Alphabet, ch) {
	//         return fmt.Errorf("invalid character in code: %c", ch)
	//     }
	// }

	return fmt.Errorf("not implemented")
}

// MaxIDForLength calculates the maximum ID that can be represented with the given code length.
// This is useful for capacity planning.
func MaxIDForLength(length int) int64 {
	return int64(math.Pow(base, float64(length))) - 1
}
