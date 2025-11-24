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

.PHONY: all clean windows darwin-amd64 darwin-arm64 linux-amd64 linux-arm64 linux-arm help

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

# Create release packages (binaries only, assets are embedded)
release: all
	@echo "$(YELLOW)Creating release packages...$(NC)"
	@mkdir -p $(BUILD_DIR)/packages
	@# Windows package
	@mkdir -p $(BUILD_DIR)/packages/windows-amd64
	@cp $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe $(BUILD_DIR)/packages/windows-amd64/$(APP_NAME).exe
	@cd $(BUILD_DIR)/packages && zip -r windows-amd64.zip windows-amd64
	@# macOS Intel package
	@mkdir -p $(BUILD_DIR)/packages/darwin-amd64
	@cp $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 $(BUILD_DIR)/packages/darwin-amd64/$(APP_NAME)
	@cd $(BUILD_DIR)/packages && tar -czf darwin-amd64.tar.gz darwin-amd64
	@# macOS Apple Silicon package
	@mkdir -p $(BUILD_DIR)/packages/darwin-arm64
	@cp $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 $(BUILD_DIR)/packages/darwin-arm64/$(APP_NAME)
	@cd $(BUILD_DIR)/packages && tar -czf darwin-arm64.tar.gz darwin-arm64
	@# Linux amd64 package
	@mkdir -p $(BUILD_DIR)/packages/linux-amd64
	@cp $(BUILD_DIR)/$(APP_NAME)-linux-amd64 $(BUILD_DIR)/packages/linux-amd64/$(APP_NAME)
	@cd $(BUILD_DIR)/packages && tar -czf linux-amd64.tar.gz linux-amd64
	@# Linux ARM64 package
	@mkdir -p $(BUILD_DIR)/packages/linux-arm64
	@cp $(BUILD_DIR)/$(APP_NAME)-linux-arm64 $(BUILD_DIR)/packages/linux-arm64/$(APP_NAME)
	@cd $(BUILD_DIR)/packages && tar -czf linux-arm64.tar.gz linux-arm64
	@# Linux ARM package
	@mkdir -p $(BUILD_DIR)/packages/linux-arm
	@cp $(BUILD_DIR)/$(APP_NAME)-linux-arm $(BUILD_DIR)/packages/linux-arm/$(APP_NAME)
	@cd $(BUILD_DIR)/packages && tar -czf linux-arm.tar.gz linux-arm
	@echo "$(GREEN)✓ Release packages created in $(BUILD_DIR)/packages/$(NC)"

