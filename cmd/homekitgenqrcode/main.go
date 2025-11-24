package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/lordbasex/HomeKitGenQRCode/internal/generator"

	"github.com/spf13/cobra"
)

// Global variables for the generate command flags
var (
	password string // Setup password in format XXX-XX-XXX
	setupID  string // Setup ID (4 alphanumeric characters)
	mac      string // MAC address (12 hexadecimal characters)
	output   string // Output image file path
	category int    // HomeKit device category ID
)

// Variables for code command flags
var (
	codeCategory int    // HomeKit device category ID
	codeOutput   string // Output image file path
	codeSetupID  string // Setup ID (optional, auto-generated if not provided)
	codeMAC      string // MAC address (optional, auto-generated if not provided)
)

// rootCmd is the root command for the CLI application
var rootCmd = &cobra.Command{
	Use:   "homekitgenqrcode",
	Short: "Generate HomeKit QR code labels",
	Long: `Generate HomeKit QR code labels with device information.

This tool generates professional HomeKit setup labels with QR codes,
device codes, serial numbers, and other required information.

Examples:
  # Generate a QR code label with all parameters
  homekitgenqrcode generate --category 5 --password "482-91-573" --setup-id "HSPN" --mac "30AEA40506A0" --output example.png
  
  # Using short flags
  homekitgenqrcode generate -c 5 -p "482-91-573" -s "HSPN" -m "30AEA40506A0" -o example.png
  
  # Generate with auto-generated setup code (easiest way)
  homekitgenqrcode code -c 5 -o example.png
  
  # Generate with custom setup ID and MAC
  homekitgenqrcode code -c 5 -o example.png -s HSPN -m 30AEA40506A0
  
  # List all available categories
  homekitgenqrcode list-categories

For more documentation, visit: https://github.com/lordbasex/HomeKitGenQRCode`,
}

// generateCmd is the command for generating QR codes with all parameters specified
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a HomeKit QR code label",
	Long: `Generate a HomeKit QR code label with the specified parameters.

All parameters are required:
  - category: HomeKit device category ID (use 'list-categories' to see available options)
  - password: Setup password in format XXX-XX-XXX (e.g., 482-91-573)
  - setup-id: Setup ID with 4 alphanumeric characters (0-9, A-Z) (e.g., HSPN)
  - mac: MAC address with 12 hexadecimal characters (e.g., 30AEA40506A0)
  - output: Output image file path (PNG format, directory will be created if needed)`,
	RunE: runGenerate,
}

// listCategoriesCmd lists all available HomeKit device categories
var listCategoriesCmd = &cobra.Command{
	Use:   "list-categories",
	Short: "List all available HomeKit categories",
	Long:  "Display all available HomeKit device categories with their IDs.",
	Run:   runListCategories,
}

// codeCmd is the command for generating QR codes with auto-generated setup code
var codeCmd = &cobra.Command{
	Use:   "code",
	Short: "Generate a HomeKit QR code label with auto-generated setup code",
	Long: `Generate a HomeKit QR code label with automatically generated setup code.

This command automatically generates:
  - Setup code (password) in format XXX-XX-XXX
  - Setup ID (4 alphanumeric characters)
  - MAC address (12 hexadecimal characters)

You only need to provide:
  - category: HomeKit device category ID
  - output: Output image file path

Examples:
  # Generar con valores completamente automáticos
  homekitgenqrcode code -c 5 -o example.png
  
  # Generar con setup ID y MAC personalizados
  homekitgenqrcode code -c 5 -o example.png -s HSPN -m 30AEA40506A0
  
  # Solo MAC personalizado (setup ID se genera automáticamente)
  homekitgenqrcode code -c 5 -o example.png -m 30AEA40506A0
  
  # Solo setup ID personalizado (MAC se genera automáticamente)
  homekitgenqrcode code -c 5 -o example.png -s HSPN
  
  # Usando flags largos
  homekitgenqrcode code --category 5 --output example.png
  
  # Generar en directorio específico (se crea automáticamente)
  homekitgenqrcode code -c 5 -o output/example.png

For more documentation, visit: https://github.com/lordbasex/HomeKitGenQRCode`,
	RunE: runCode,
}

// init initializes the CLI commands and their flags
func init() {
	// Generate command flags
	generateCmd.Flags().IntVarP(&category, "category", "c", 0, "HomeKit category ID (required)")
	generateCmd.Flags().StringVarP(&password, "password", "p", "", "Setup password in format XXX-XX-XXX (required)")
	generateCmd.Flags().StringVarP(&setupID, "setup-id", "s", "", "Setup ID: 4 alphanumeric characters (0-9, A-Z) (required)")
	generateCmd.Flags().StringVarP(&mac, "mac", "m", "", "MAC address: 12 hexadecimal characters (required)")
	generateCmd.Flags().StringVarP(&output, "output", "o", "", "Output image file path (PNG format) (required)")

	// Mark flags as required
	generateCmd.MarkFlagRequired("category")
	generateCmd.MarkFlagRequired("password")
	generateCmd.MarkFlagRequired("setup-id")
	generateCmd.MarkFlagRequired("mac")
	generateCmd.MarkFlagRequired("output")

	// Code command flags
	codeCmd.Flags().IntVarP(&codeCategory, "category", "c", 0, "HomeKit category ID (required)")
	codeCmd.Flags().StringVarP(&codeOutput, "output", "o", "", "Output image file path (PNG format) (required)")
	codeCmd.Flags().StringVarP(&codeSetupID, "setup-id", "s", "", "Setup ID: 4 alphanumeric characters (0-9, A-Z) (optional, auto-generated if not provided)")
	codeCmd.Flags().StringVarP(&codeMAC, "mac", "m", "", "MAC address: 12 hexadecimal characters (optional, auto-generated if not provided)")

	codeCmd.MarkFlagRequired("category")
	codeCmd.MarkFlagRequired("output")

	// Add commands to root
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(listCategoriesCmd)
	rootCmd.AddCommand(codeCmd)
}

// runGenerate executes the generate command
// It validates all inputs and generates the HomeKit QR code label
func runGenerate(cmd *cobra.Command, args []string) error {
	// Trim whitespace from inputs
	password = strings.TrimSpace(password)
	setupID = strings.TrimSpace(strings.ToUpper(setupID))
	mac = strings.TrimSpace(strings.ToUpper(mac))
	output = strings.TrimSpace(output)

	// Validate all inputs
	if err := validateInputs(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	// Ensure output directory exists
	if err := ensureOutputDirectory(output); err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}

	// Generate the HomeKit label
	err := generator.GenerateHomeKitLabel(category, password, setupID, mac, output)
	if err != nil {
		return fmt.Errorf("error generating label: %w", err)
	}

	fmt.Printf("\n✅ QR-code opgeslagen als: %s\n", output)
	return nil
}

// runListCategories executes the list-categories command
// It displays all available HomeKit device categories sorted by ID
func runListCategories(cmd *cobra.Command, args []string) {
	fmt.Println("Available HomeKit Categories:")
	fmt.Println(strings.Repeat("=", 50))

	// Get all categories and sort by ID
	ids := make([]int, 0, len(generator.CategoryReference))
	for id := range generator.CategoryReference {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	for _, id := range ids {
		fmt.Printf("  %2d: %s\n", id, generator.CategoryReference[id])
	}
	fmt.Println()
}

// runCode executes the code command
// It generates a setup code automatically and optionally generates setup ID and MAC address
func runCode(cmd *cobra.Command, args []string) error {
	// Validate category
	if codeCategory < 1 {
		return fmt.Errorf("category must be a positive number")
	}
	if _, exists := generator.CategoryReference[codeCategory]; !exists {
		return fmt.Errorf("invalid category ID: %d. Use 'list-categories' to see available categories", codeCategory)
	}

	// Validate output path
	codeOutput = strings.TrimSpace(codeOutput)
	if codeOutput == "" {
		return fmt.Errorf("output path cannot be empty")
	}
	if !strings.HasSuffix(strings.ToLower(codeOutput), ".png") {
		return fmt.Errorf("output file must have .png extension")
	}

	// Generate setup code automatically
	setupCode := generator.GenerateHomeKitSetupCode()

	// Generate setup ID if not provided
	if codeSetupID == "" {
		codeSetupID = generateRandomSetupID()
	} else {
		codeSetupID = strings.TrimSpace(strings.ToUpper(codeSetupID))
		if err := validateSetupID(codeSetupID); err != nil {
			return fmt.Errorf("invalid setup ID: %w", err)
		}
	}

	// Generate MAC address if not provided
	if codeMAC == "" {
		codeMAC = generateRandomMAC()
	} else {
		codeMAC = strings.TrimSpace(strings.ToUpper(codeMAC))
		if err := validateMAC(codeMAC); err != nil {
			return fmt.Errorf("invalid MAC address: %w", err)
		}
	}

	// Display generated values
	fmt.Println("Generated HomeKit Setup Information:")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("  Setup Code:    %s\n", setupCode)
	fmt.Printf("  Setup ID:      %s\n", codeSetupID)
	fmt.Printf("  MAC Address:   %s\n", formatMACDisplay(codeMAC))
	fmt.Printf("  Category:      %d (%s)\n", codeCategory, generator.CategoryReference[codeCategory])
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println()

	// Ensure output directory exists
	if err := ensureOutputDirectory(codeOutput); err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}

	// Generate the HomeKit label
	err := generator.GenerateHomeKitLabel(codeCategory, setupCode, codeSetupID, codeMAC, codeOutput)
	if err != nil {
		return fmt.Errorf("error generating label: %w", err)
	}

	fmt.Printf("✅ QR-code opgeslagen als: %s\n", codeOutput)
	return nil
}

// generateRandomSetupID generates a random 4-character setup ID.
// Characters are selected from 0-9 and A-Z (36 possible characters).
func generateRandomSetupID() string {
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, 4)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

// generateRandomMAC generates a random 12-character hexadecimal MAC address.
// Characters are selected from 0-9 and A-F (16 possible characters).
func generateRandomMAC() string {
	const hexChars = "0123456789ABCDEF"
	result := make([]byte, 12)
	for i := range result {
		result[i] = hexChars[rand.Intn(len(hexChars))]
	}
	return string(result)
}

// formatMACDisplay formats a MAC address for display by adding colons every 2 characters.
// Example: "30AEA40506A0" -> "30:AE:A4:05:06:A0"
// If the MAC address is not 12 characters, returns it unchanged.
func formatMACDisplay(mac string) string {
	if len(mac) != 12 {
		return mac
	}
	var parts []string
	for i := 0; i < 12; i += 2 {
		parts = append(parts, mac[i:i+2])
	}
	return strings.Join(parts, ":")
}

// validateInputs validates all input parameters for the generate command.
// Returns an error if any validation fails.
func validateInputs() error {
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

	// Validate output path
	if output == "" {
		return fmt.Errorf("output path cannot be empty")
	}
	if !strings.HasSuffix(strings.ToLower(output), ".png") {
		return fmt.Errorf("output file must have .png extension")
	}

	return nil
}

// validatePassword validates password format XXX-XX-XXX.
// The password must be exactly 10 characters: 3 digits, dash, 2 digits, dash, 3 digits.
// All parts must be numeric.
func validatePassword(pwd string) error {
	// Trim whitespace
	pwd = strings.TrimSpace(pwd)
	// Format XXX-XX-XXX has 10 characters: 3 + 1 + 2 + 1 + 3 = 10
	if len(pwd) != 10 {
		return fmt.Errorf("invalid password length (%d). Expected format: XXX-XX-XXX (e.g., 482-91-573)", len(pwd))
	}

	parts := strings.Split(pwd, "-")
	if len(parts) != 3 {
		return fmt.Errorf("invalid password format. Expected format: XXX-XX-XXX (e.g., 482-91-573)")
	}

	if len(parts[0]) != 3 || len(parts[1]) != 2 || len(parts[2]) != 3 {
		return fmt.Errorf("invalid password format. Expected format: XXX-XX-XXX (e.g., 482-91-573)")
	}

	for i, part := range parts {
		if _, err := strconv.Atoi(part); err != nil {
			return fmt.Errorf("password part %d contains non-numeric characters. Expected format: XXX-XX-XXX", i+1)
		}
	}

	return nil
}

// validateSetupID validates setup ID format (4 alphanumeric: 0-9, A-Z).
// The setup ID must be exactly 4 characters and contain only alphanumeric characters.
// The input is automatically converted to uppercase.
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

// validateMAC validates MAC address format (12 hexadecimal characters).
// The MAC address must be exactly 12 characters and contain only hexadecimal characters (0-9, A-F).
// The input is automatically converted to uppercase.
func validateMAC(macAddr string) error {
	// Trim whitespace and convert to uppercase
	macAddr = strings.TrimSpace(strings.ToUpper(macAddr))
	if len(macAddr) != 12 {
		return fmt.Errorf("invalid MAC address length. Expected 12 hexadecimal characters (e.g., 30AEA40506A0)")
	}

	for i, r := range macAddr {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')) {
			return fmt.Errorf("invalid character '%c' at position %d. MAC address must contain only hexadecimal characters (0-9, A-F)", r, i+1)
		}
	}

	return nil
}

// ensureOutputDirectory creates the output directory if it doesn't exist.
// If the output path is in the current directory (empty or "."), no directory is created.
// The directory is created with permissions 0755.
func ensureOutputDirectory(outputPath string) error {
	outputDir := filepath.Dir(outputPath)

	// If outputDir is empty or ".", it means the file is in the current directory
	if outputDir == "" || outputDir == "." {
		return nil
	}

	// Check if directory exists
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		// Create directory with permissions 0755
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory '%s': %w", outputDir, err)
		}
		fmt.Printf("📁 Created output directory: %s\n", outputDir)
	}

	return nil
}

// main is the entry point of the application.
// It executes the root command and handles any errors.
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}
}
