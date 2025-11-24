package generator

import (
	"fmt"
	"math/rand"
	"strings"
)

// GenerateHomeKitSetupCode generates a valid HomeKit setup code in format XXX-XX-XXX.
// Internally uses 8 digits and avoids trivial codes (sequences, repeated digits, simple patterns).
// Similar to HomeSpan's practical criteria for code generation.
func GenerateHomeKitSetupCode() string {
	for {
		// Generate 8-digit number (10000000-99999999 avoids 00000000)
		n := rand.Intn(90000000) + 10000000
		raw := fmt.Sprintf("%08d", n)

		if !isTooSimple(raw) {
			return fmt.Sprintf("%s-%s-%s", raw[0:3], raw[3:5], raw[5:8])
		}
	}
}

// PlainSetupCode converts a formatted setup code to plain format.
// Example: "613-80-755" -> "61380755"
func PlainSetupCode(code string) string {
	return strings.ReplaceAll(code, "-", "")
}

// IsValidSetupCode validates HomeKit setup code format.
// Accepts codes with or without dashes, but must be 8 digits and not "too simple".
func IsValidSetupCode(code string) bool {
	raw := PlainSetupCode(code)
	if len(raw) != 8 {
		return false
	}
	for _, ch := range raw {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return !isTooSimple(raw)
}

// isTooSimple checks if a setup code is too simple and should be rejected.
// Implements anti-simple rules (practical style similar to HomeSpan):
// - All digits the same (00000000, 11111111, ...)
// - Ascending or descending sequence (12345678 / 87654321)
// - Simple repetitive pattern (ABABABAB numeric: 12121212, 34343434, ...)
// - Additional blacklist of common simple codes
func isTooSimple(raw string) bool {
	// Check if all digits are the same (00000000, 11111111, ...)
	allSame := true
	for i := 1; i < 8; i++ {
		if raw[i] != raw[0] {
			allSame = false
			break
		}
	}
	if allSame {
		return true
	}

	// Check for ascending or descending sequence (12345678 / 87654321)
	asc, desc := true, true
	for i := 1; i < 8; i++ {
		if raw[i] != raw[i-1]+1 {
			asc = false
		}
		if raw[i] != raw[i-1]-1 {
			desc = false
		}
	}
	if asc || desc {
		return true
	}

	// Check for simple repetitive pattern (ABABABAB numeric: 12121212, 34343434, ...)
	if raw[0:2] == raw[2:4] &&
		raw[2:4] == raw[4:6] &&
		raw[4:6] == raw[6:8] {
		return true
	}

	// Additional minimal blacklist
	switch raw {
	case "12345678", "87654321", "00000000", "11111111":
		return true
	}

	return false
}
