//go:build wasm
// +build wasm

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"syscall/js"

	"github.com/lordbasex/HomeKitGenQRCode/internal/generator"
)

// GenerateLabelRequest represents the request from JavaScript
type GenerateLabelRequest struct {
	Category int    `json:"category"`
	Password string `json:"password"`
	SetupID  string `json:"setupId"`
	MAC      string `json:"mac"`
}

// GenerateLabelResponse represents the response to JavaScript
type GenerateLabelResponse struct {
	ImageBase64 string `json:"imageBase64"`
	Error       string `json:"error,omitempty"`
}

// generateHomeKitLabel is the JavaScript function wrapper for generating QR code labels
func generateHomeKitLabel(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return js.ValueOf(map[string]interface{}{
			"error": "expected 1 argument (JSON string)",
		})
	}

	// Parse JSON request
	var req GenerateLabelRequest
	if err := json.Unmarshal([]byte(args[0].String()), &req); err != nil {
		return js.ValueOf(map[string]interface{}{
			"error": "invalid JSON: " + err.Error(),
		})
	}

	// Generate image bytes
	imageBytes, err := generator.GenerateHomeKitLabelBytes(
		req.Category,
		req.Password,
		req.SetupID,
		req.MAC,
	)
	if err != nil {
		return js.ValueOf(map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Convert to base64 data URL
	imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)

	response := GenerateLabelResponse{
		ImageBase64: "data:image/png;base64," + imageBase64,
	}

	jsonResponse, _ := json.Marshal(response)
	return js.ValueOf(string(jsonResponse))
}

// generateRandomCode generates random setup code, setup ID, and MAC address
func generateRandomCode(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return js.ValueOf(map[string]interface{}{
			"error": "expected 1 argument (category)",
		})
	}

	category := args[0].Int()
	setupCode := generator.GenerateHomeKitSetupCode()
	setupID := generateRandomSetupID()
	mac := generateRandomMAC()

	result := map[string]interface{}{
		"setupCode": setupCode,
		"setupID":   setupID,
		"mac":       mac,
		"category":  category,
	}

	jsonResponse, _ := json.Marshal(result)
	return js.ValueOf(string(jsonResponse))
}

// listCategories returns all available HomeKit categories
func listCategories(this js.Value, args []js.Value) interface{} {
	categories := make(map[int]string)
	for id, name := range generator.CategoryReference {
		categories[id] = name
	}

	jsonResponse, _ := json.Marshal(categories)
	return js.ValueOf(string(jsonResponse))
}

// generateRandomSetupID generates a random 4-character setup ID (0-9, A-Z)
func generateRandomSetupID() string {
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, 4)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

// generateRandomMAC generates a random 12-character hexadecimal MAC address
func generateRandomMAC() string {
	const hexChars = "0123456789ABCDEF"
	result := make([]byte, 12)
	for i := range result {
		result[i] = hexChars[rand.Intn(len(hexChars))]
	}
	return string(result)
}

// ValidationRequest represents a validation request
type ValidationRequest struct {
	Category int    `json:"category"`
	Password string `json:"password"`
	SetupID  string `json:"setupId"`
	MAC      string `json:"mac"`
}

// ValidationResponse represents a validation response
type ValidationResponse struct {
	Valid bool   `json:"valid"`
	Error string `json:"error,omitempty"`
}

// validateInputs validates all input parameters (same as CLI)
func validateInputs(category int, password, setupID, mac string) error {
	// Validate category
	if category < 1 {
		return fmt.Errorf("category must be a positive number")
	}
	if _, exists := generator.CategoryReference[category]; !exists {
		return fmt.Errorf("invalid category ID: %d. Use 'list-categories' to see available categories", category)
	}

	// Validate password format: XXX-XX-XXX
	if err := validatePassword(password); err != nil {
		return err
	}

	// Validate setup ID format: 4 alphanumeric characters (0-9, A-Z)
	if err := validateSetupID(setupID); err != nil {
		return err
	}

	// Validate MAC address format: 12 hexadecimal characters
	if err := validateMAC(mac); err != nil {
		return err
	}

	return nil
}

// validatePassword validates password format XXX-XX-XXX (same as CLI)
func validatePassword(pwd string) error {
	// Trim whitespace
	pwd = strings.TrimSpace(pwd)
	// Format XXX-XX-XXX has 10 characters: 3 + 1 + 2 + 1 + 3 = 10
	if len(pwd) != 10 {
		return fmt.Errorf("invalid password length (%d). Expected format: XXX-XX-XXX (e.g., 613-80-755)", len(pwd))
	}

	parts := strings.Split(pwd, "-")
	if len(parts) != 3 {
		return fmt.Errorf("invalid password format. Expected format: XXX-XX-XXX (e.g., 613-80-755)")
	}

	if len(parts[0]) != 3 || len(parts[1]) != 2 || len(parts[2]) != 3 {
		return fmt.Errorf("invalid password format. Expected format: XXX-XX-XXX (e.g., 613-80-755)")
	}

	for i, part := range parts {
		if _, err := strconv.Atoi(part); err != nil {
			return fmt.Errorf("password part %d contains non-numeric characters. Expected format: XXX-XX-XXX", i+1)
		}
	}

	return nil
}

// validateSetupID validates setup ID format (same as CLI)
func validateSetupID(id string) error {
	// Trim whitespace and convert to uppercase
	id = strings.TrimSpace(strings.ToUpper(id))
	if len(id) != 4 {
		return fmt.Errorf("invalid setup ID length. Expected 4 alphanumeric characters (0-9, A-Z)")
	}

	for i, r := range id {
		if !((r >= '0' && r <= '9') || (r >= 'A' && r <= 'Z')) {
			return fmt.Errorf("invalid character '%c' at position %d. Setup ID must contain only 0-9 and A-Z", r, i+1)
		}
	}

	return nil
}

// validateMAC validates MAC address format (same as CLI)
func validateMAC(macAddr string) error {
	// Trim whitespace and convert to uppercase
	macAddr = strings.TrimSpace(strings.ToUpper(macAddr))
	if len(macAddr) != 12 {
		return fmt.Errorf("invalid MAC address length. Expected 12 hexadecimal characters (e.g., AABBCCDDEEFF)")
	}

	for i, r := range macAddr {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')) {
			return fmt.Errorf("invalid character '%c' at position %d. MAC address must contain only hexadecimal characters (0-9, A-F)", r, i+1)
		}
	}

	return nil
}

// validateInputsWASM is the JavaScript function wrapper for validation
func validateInputsWASM(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return js.ValueOf(map[string]interface{}{
			"valid": false,
			"error": "expected 1 argument (JSON string)",
		})
	}

	// Parse JSON request
	var req ValidationRequest
	if err := json.Unmarshal([]byte(args[0].String()), &req); err != nil {
		return js.ValueOf(map[string]interface{}{
			"valid": false,
			"error": "invalid JSON: " + err.Error(),
		})
	}

	// Validate inputs
	if err := validateInputs(req.Category, req.Password, req.SetupID, req.MAC); err != nil {
		response := ValidationResponse{
			Valid: false,
			Error: err.Error(),
		}
		jsonResponse, _ := json.Marshal(response)
		return js.ValueOf(string(jsonResponse))
	}

	response := ValidationResponse{
		Valid: true,
	}
	jsonResponse, _ := json.Marshal(response)
	return js.ValueOf(string(jsonResponse))
}

func main() {
	c := make(chan struct{}, 0)

	// Export JavaScript functions
	js.Global().Set("generateHomeKitLabel", js.FuncOf(generateHomeKitLabel))
	js.Global().Set("generateRandomCode", js.FuncOf(generateRandomCode))
	js.Global().Set("listCategories", js.FuncOf(listCategories))
	js.Global().Set("validateInputs", js.FuncOf(validateInputsWASM))

	// Keep the program running
	<-c
}

