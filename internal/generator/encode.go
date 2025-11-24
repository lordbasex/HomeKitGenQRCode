package generator

import (
	"fmt"
	"strconv"
	"strings"
)

// base36 contains the characters used for base36 encoding (0-9, A-Z)
const base36 = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// GenHomeKitSetupURI generates a HomeKit setup URI from the provided parameters.
// The URI format is: X-HM://{encoded_payload}{setupID}
//
// Parameters:
//   - category: HomeKit device category ID
//   - password: Setup password in format XXX-XX-XXX (will be converted to plain format)
//   - setupID: Setup ID (4 alphanumeric characters)
//
// The payload is encoded as follows:
//   - version (3 bits): Currently 0
//   - reserved (4 bits): Currently 0
//   - category (8 bits): Device category ID
//   - flags (4 bits): Currently 2
//   - password (27 bits): 8-digit password without dashes
//
// The payload is then base36 encoded to create the URI.
func GenHomeKitSetupURI(category int, password, setupID string) string {
	version := 0
	reserved := 0
	flags := 2

	// Build payload by bit-shifting and ORing values
	payload := 0
	payload |= (version & 0x7)
	payload <<= 4
	payload |= (reserved & 0xF)
	payload <<= 8
	payload |= (category & 0xFF)
	payload <<= 4
	payload |= (flags & 0xF)
	payload <<= 27

	// Convert password to plain format and add to payload
	num, _ := strconv.Atoi(strings.ReplaceAll(password, "-", ""))
	payload |= (num & 0x7FFFFFF)

	// Encode payload to base36 (9 characters)
	out := make([]byte, 9)
	for i := 8; i >= 0; i-- {
		out[i] = base36[payload%36]
		payload /= 36
	}

	return fmt.Sprintf("X-HM://%s%s", string(out), setupID)
}
