package generator

import (
	"fmt"
	"math/rand"
	"strings"
)

// CategoryReference maps HomeKit category IDs to their human-readable names.
// This map contains all supported HomeKit device categories.
var CategoryReference = map[int]string{
	1: "Other", 2: "Bridge", 3: "Fan", 4: "Garage Door Opener", 5: "Light",
	6: "Lock", 7: "Outlet", 8: "Switch", 9: "Thermostat", 10: "Sensor",
	11: "Security system", 12: "Door", 13: "Window", 14: "Window covering",
	15: "Programmable switch", 16: "Range extender", 17: "IP camera",
	18: "Video doorbell", 19: "Air purifier", 20: "Heater", 21: "Air conditioner",
	22: "Humidifier", 23: "Dehumidifier", 24: "Apple TV", 26: "Speaker",
	27: "Airport", 28: "Sprinkler", 29: "Faucet", 30: "Shower head",
	31: "Television", 32: "Target remote",
}

// Note: As of Go 1.20+, rand.Seed is deprecated.
// The global random generator is automatically seeded, so we don't need init()

// randLetter generates a random uppercase letter (A-Z)
func randLetter() string {
	return string('A' + rune(rand.Intn(26)))
}

// GenerateDeviceCode generates a device code matching the Python implementation format.
// Format: XX{category}X{X}/X (e.g., AB3C2DE/F)
// The code consists of:
//   - 2 uppercase letters
//   - Category ID (1 digit)
//   - 1 uppercase letter
//   - 1 random digit (0-9)
//   - 2 uppercase letters
//   - "/"
//   - 1 uppercase letter
func GenerateDeviceCode(category int) string {
	return fmt.Sprintf("%s%d%s%d%s/%s",
		randLetter()+randLetter(), // 2 uppercase letters
		category,                  // category ID
		randLetter(),              // 1 uppercase letter
		rand.Intn(10),             // 1 digit
		randLetter()+randLetter(), // 2 uppercase letters
		randLetter(),              // 1 uppercase letter
	)
}

// GenerateSerial generates a serial number matching the Python implementation format.
// Format: X{X}X{X}X{X}X{X}{X}{X}X{X} (12 characters total)
// Pattern: Letter-Digit-Letter-Letter-Letter-Digit-Letter-Digit-Digit-Digit-Letter-Letter
func GenerateSerial() string {
	pattern := []string{
		randLetter(),              // 0: letter
		fmt.Sprint(rand.Intn(10)), // 1: digit
		randLetter(),              // 2: letter
		randLetter(),              // 3: letter
		randLetter(),              // 4: letter
		fmt.Sprint(rand.Intn(10)), // 5: digit
		randLetter(),              // 6: letter
		fmt.Sprint(rand.Intn(10)), // 7: digit
		fmt.Sprint(rand.Intn(10)), // 8: digit
		fmt.Sprint(rand.Intn(10)), // 9: digit
		randLetter(),              // 10: letter
		randLetter(),              // 11: letter
	}
	return strings.Join(pattern, "")
}

// GenerateCSN generates a CSN (Customer Serial Number) matching the Python implementation format.
// Format: {20 digits}{3 letters}{4 digits}{letter}{digit}{letter}{3 digits}
// Total length: 20 + 3 + 4 + 1 + 1 + 1 + 3 = 33 characters
func GenerateCSN() string {
	var b strings.Builder

	// Generate 20 digits
	for i := 0; i < 20; i++ {
		b.WriteString(fmt.Sprint(rand.Intn(10)))
	}

	// Generate 3 letters
	for i := 0; i < 3; i++ {
		b.WriteString(randLetter())
	}

	// Generate 4 digits
	for i := 0; i < 4; i++ {
		b.WriteString(fmt.Sprint(rand.Intn(10)))
	}

	// Generate letter + digit + letter
	b.WriteString(randLetter())
	b.WriteString(fmt.Sprint(rand.Intn(10)))
	b.WriteString(randLetter())

	// Generate 3 digits
	for i := 0; i < 3; i++ {
		b.WriteString(fmt.Sprint(rand.Intn(10)))
	}

	return b.String()
}
