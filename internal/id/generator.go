package id

import (
	"context"
	"fmt"
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
// DONE: Implement the ID generation logic:
// 1. Use Redis INCR on "global:url_id" to get a unique ID
// 2. Encode the ID to Base62
// 3. Pad to minimum length if necessary
// 4. Handle Redis errors appropriately
func (g *Generator) Generate(ctx context.Context) (string, error) {
	// DONE: Implement Redis INCR operation
	// id, err := g.redis.Incr(ctx, "global:url_id").Result()
	// if err != nil {
	//     return "", fmt.Errorf("failed to generate ID: %w", err)
	// }

	id, err := g.redis.Incr(ctx, "global:url_id").Result()

	if err != nil {
		return "", err
	}

	// DONE: Encode ID to Base62
	// code := g.encodeBase62(id)

	code := g.encodeBase62(id)

	// TODO: Pad to minimum length
	// code = g.padToMinLength(code)
	code = g.padToMinLength(code)

	return code, nil
}

// encodeBase62 encodes a number to Base62 string.
// DONE: Implement Base62 encoding algorithm:
// 1. Handle zero case
// 2. Convert number to base 62 using the alphabet
// 3. Reverse the result string
func (g *Generator) encodeBase62(num int64) string {
	if num == 0 {
		return string(base62Alphabet[0])
	}

	// DONE: Implement encoding logic
	var b []byte
	for num > 0 {
		rem := num % 62
		b = append(b, base62Alphabet[rem])
		num /= 62
	}

	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

	return string(b)
}

// DecodeBase62 decodes a Base62 string back to a number.
// This is useful for testing and debugging.
// DONE: Implement Base62 decoding algorithm:
// 1. Iterate through each character
// 2. Find its position in the alphabet
// 3. Calculate the numeric value
func DecodeBase62(code string) (int64, error) {
	// DONE: Implement decoding logic
	// ... decoding logic here ...

	if len(code) == 0 {
		return 0, fmt.Errorf("cannot decode empty string")
	}

	var result int64
	for i := 0; i < len(code); i++ {
		j := strings.Index(base62Alphabet, string(code[i]))
		if j < 0 {
			return 0, fmt.Errorf("cannot decode string due to illegal char: %c", code[i])
		}
		result = result*62 + int64(j)
	}

	return result, nil
}

// padToMinLength pads the code to the minimum length by prepending '0' characters.
// DONE: Implement padding logic:
// 1. Check if code is already at or above minimum length
// 2. Calculate padding needed
// 3. Prepend '0' characters (first character in alphabet)
func (g *Generator) padToMinLength(code string) string {
	if len(code) >= g.minLength {
		return code
	}

	// DONE: Implement padding
	padding := g.minLength - len(code)
	return strings.Repeat(string(base62Alphabet[0]), padding) + code
}

// ValidateCode validates if a code is a valid Base62 string.
// DONE: Implement validation logic:
// 1. Check if code is empty
// 2. Check if all characters are in the Base62 alphabet
// 3. Check length constraints
func ValidateCode(code string) error {
	if code == "" {
		return fmt.Errorf("code cannot be empty")
	}

	// DONE: Implement character validation
	// for _, ch := range code {
	//     if !strings.ContainsRune(base62Alphabet, ch) {
	//         return fmt.Errorf("invalid character in code: %c", ch)
	//     }
	// }

	for _, c := range code {
		if !strings.ContainsRune(base62Alphabet, c) {
			return fmt.Errorf("invalid character in code: %c", c)
		}
	}

	return nil
}

// MaxIDForLength calculates the maximum ID that can be represented with the given code length.
// This is useful for capacity planning.
func MaxIDForLength(length int) int64 {
	if length <= 0 {
		return 0
	}

	var result int64 = 1
	for i := 0; i < length; i++ {
		result *= 62
	}
	return result - 1

}
