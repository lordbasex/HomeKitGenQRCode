# Makefile for HomeKitGenQRCode
# Builds binaries for multiple platforms

# Application name
APP_NAME := homekitgenqrcode
VERSION := 1.0.0
BUILD_DIR := dist
CMD_DIR := cmd/homekitgenqrcode

# Build flags
LDFLAGS := -s -w -X main.version=$(VERSION)

# Colors for output
GREEN := \033[0;32m
YELLOW := \033[0;33m
NC := \033[0m # No Color

.PHONY: all clean windows darwin-amd64 darwin-arm64 linux-amd64 linux-arm64 linux-arm help release

# Default target
all: clean windows darwin-amd64 darwin-arm64 linux-amd64 linux-arm64 linux-arm
	@echo "$(GREEN)✓ All binaries built successfully!$(NC)"

# Help target
help:
	@echo "Available targets:"
	@echo "  all              - Build all platforms"
	@echo "  windows          - Build for Windows (amd64)"
	@echo "  darwin-amd64     - Build for macOS Intel (amd64)"
	@echo "  darwin-arm64     - Build for macOS Apple Silicon (arm64)"
	@echo "  linux-amd64     - Build for Linux 64-bit (amd64)"
	@echo "  linux-arm64     - Build for Linux ARM64 (arm64) - Raspberry Pi 4+"
	@echo "  linux-arm       - Build for Linux ARM (32-bit) - Raspberry Pi 3 and older"
	@echo "  release          - Upload uncompressed binaries to GitHub release"
	@echo "  clean            - Remove build directory"
	@echo ""
	@echo "Example: make all"

# Create build directory
$(BUILD_DIR):
	@mkdir -p $(BUILD_DIR)

# Windows (amd64)
windows: $(BUILD_DIR)
	@echo "$(YELLOW)Building for Windows (amd64)...$(NC)"
	@GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe ./$(CMD_DIR)
	@echo "$(GREEN)✓ Windows binary created: $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe$(NC)"

# macOS Intel (amd64)
darwin-amd64: $(BUILD_DIR)
	@echo "$(YELLOW)Building for macOS Intel (amd64)...$(NC)"
	@GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 ./$(CMD_DIR)
	@echo "$(GREEN)✓ macOS Intel binary created: $(BUILD_DIR)/$(APP_NAME)-darwin-amd64$(NC)"

# macOS Apple Silicon (arm64)
darwin-arm64: $(BUILD_DIR)
	@echo "$(YELLOW)Building for macOS Apple Silicon (arm64)...$(NC)"
	@GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 ./$(CMD_DIR)
	@echo "$(GREEN)✓ macOS Apple Silicon binary created: $(BUILD_DIR)/$(APP_NAME)-darwin-arm64$(NC)"

# Linux 64-bit (amd64)
linux-amd64: $(BUILD_DIR)
	@echo "$(YELLOW)Building for Linux 64-bit (amd64)...$(NC)"
	@GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 ./$(CMD_DIR)
	@echo "$(GREEN)✓ Linux 64-bit binary created: $(BUILD_DIR)/$(APP_NAME)-linux-amd64$(NC)"

# Linux ARM64 (arm64) - Raspberry Pi 4 and newer
linux-arm64: $(BUILD_DIR)
	@echo "$(YELLOW)Building for Linux ARM64 (arm64) - Raspberry Pi 4+...$(NC)"
	@GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-linux-arm64 ./$(CMD_DIR)
	@echo "$(GREEN)✓ Linux ARM64 binary created: $(BUILD_DIR)/$(APP_NAME)-linux-arm64$(NC)"

# Linux ARM (32-bit) - Raspberry Pi 3 and older
linux-arm: $(BUILD_DIR)
	@echo "$(YELLOW)Building for Linux ARM (32-bit) - Raspberry Pi 3 and older...$(NC)"
	@GOOS=linux GOARCH=arm go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-linux-arm ./$(CMD_DIR)
	@echo "$(GREEN)✓ Linux ARM binary created: $(BUILD_DIR)/$(APP_NAME)-linux-arm$(NC)"

# Clean build directory
clean:
	@echo "$(YELLOW)Cleaning build directory...$(NC)"
	@rm -rf $(BUILD_DIR)
	@echo "$(GREEN)✓ Clean complete$(NC)"

# Install dependencies
deps:
	@echo "$(YELLOW)Downloading dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)✓ Dependencies updated$(NC)"

# Build for current platform
local:
	@echo "$(YELLOW)Building for current platform...$(NC)"
	@go build -ldflags "$(LDFLAGS)" -o $(APP_NAME) ./$(CMD_DIR)
	@echo "$(GREEN)✓ Local binary created: $(APP_NAME)$(NC)"

# Upload binaries to GitHub release (uncompressed, assets are embedded)
release: all
	@echo "$(YELLOW)Uploading binaries to GitHub release v$(VERSION)...$(NC)"
	@gh release delete v$(VERSION) --yes 2>/dev/null || true
	@echo "## Release v$(VERSION)" > /tmp/release-notes.txt
	@echo "" >> /tmp/release-notes.txt
	@echo "This release includes standalone binaries for all major platforms. **No external files needed** - all assets are embedded in the binary." >> /tmp/release-notes.txt
	@echo "" >> /tmp/release-notes.txt
	@echo "### Supported Platforms" >> /tmp/release-notes.txt
	@echo "- **Windows** (amd64) - \`homekitgenqrcode-windows-amd64.exe\`" >> /tmp/release-notes.txt
	@echo "- **macOS Intel** (amd64) - \`homekitgenqrcode-darwin-amd64\`" >> /tmp/release-notes.txt
	@echo "- **macOS Apple Silicon** (arm64) - \`homekitgenqrcode-darwin-arm64\`" >> /tmp/release-notes.txt
	@echo "- **Linux 64-bit** (amd64) - \`homekitgenqrcode-linux-amd64\`" >> /tmp/release-notes.txt
	@echo "- **Linux ARM64** (arm64) - Raspberry Pi 4+ - \`homekitgenqrcode-linux-arm64\`" >> /tmp/release-notes.txt
	@echo "- **Linux ARM** (32-bit) - Raspberry Pi 3 and older - \`homekitgenqrcode-linux-arm\`" >> /tmp/release-notes.txt
	@echo "" >> /tmp/release-notes.txt
	@echo "### Installation" >> /tmp/release-notes.txt
	@echo "1. Download the binary for your platform" >> /tmp/release-notes.txt
	@echo "2. Make it executable (Linux/macOS): \`chmod +x homekitgenqrcode-*\`" >> /tmp/release-notes.txt
	@echo "3. Run it directly - no dependencies needed!" >> /tmp/release-notes.txt
	@echo "" >> /tmp/release-notes.txt
	@echo "### Usage" >> /tmp/release-notes.txt
	@echo "\`\`\`bash" >> /tmp/release-notes.txt
	@echo "# Quick start (auto-generate all values)" >> /tmp/release-notes.txt
	@echo "./homekitgenqrcode code -c 5 -o example.png" >> /tmp/release-notes.txt
	@echo "" >> /tmp/release-notes.txt
	@echo "# With all parameters" >> /tmp/release-notes.txt
	@echo "./homekitgenqrcode generate -c 5 -p \"613-80-755\" -s \"ABCD\" -m \"AABBCCDDEEFF\" -o example.png" >> /tmp/release-notes.txt
	@echo "\`\`\`" >> /tmp/release-notes.txt
	@echo "" >> /tmp/release-notes.txt
	@echo "For more information, see the [README](https://github.com/lordbasex/HomeKitGenQRCode/blob/main/README.md)." >> /tmp/release-notes.txt
	@gh release create v$(VERSION) --title "v$(VERSION) - Release" --notes-file /tmp/release-notes.txt \
		$(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe \
		$(BUILD_DIR)/$(APP_NAME)-darwin-amd64 \
		$(BUILD_DIR)/$(APP_NAME)-darwin-arm64 \
		$(BUILD_DIR)/$(APP_NAME)-linux-amd64 \
		$(BUILD_DIR)/$(APP_NAME)-linux-arm64 \
		$(BUILD_DIR)/$(APP_NAME)-linux-arm
	@rm -f /tmp/release-notes.txt
	@echo "$(GREEN)✓ Release v$(VERSION) created with uncompressed binaries$(NC)"
