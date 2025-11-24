# HomeKitGenQRCode

A Go application for generating HomeKit QR code labels with device information.

## Description

This tool generates professional HomeKit setup labels with QR codes, device codes, serial numbers, and other required information. It was created for the [HomeSpan](https://github.com/HomeSpan/HomeSpan/) project, taking inspiration from the original Python implementation by [AchimPieters/esp32-homekit-qrcode](https://github.com/AchimPieters/esp32-homekit-qrcode), but rewritten in Go for better performance, easier distribution, and improved cross-platform compatibility.

## Features

- Generate HomeKit QR code labels with all required information
- Support for all HomeKit device categories
- Automatic generation of setup codes, setup IDs, and MAC addresses
- Professional label formatting matching Apple's HomeKit standards
- Command-line interface with multiple subcommands

## Installation

### From Source

```bash
git clone https://github.com/lordbasex/HomeKitGenQRCode.git
cd HomeKitGenQRCode
go build ./cmd/homekitgenqrcode
```

### Using Go Install

```bash
go install github.com/lordbasex/HomeKitGenQRCode/cmd/homekitgenqrcode@latest
```

## Usage

### Quick Start (Recommended)

Generate a QR code label with auto-generated setup code:

```bash
homekitgenqrcode code -c 5 -o example.png
```

### Generate with All Parameters

```bash
homekitgenqrcode generate --category 5 --password "482-91-573" --setup-id "HSPN" --mac "30AEA40506A0" --output example.png
```

### List Available Categories

```bash
homekitgenqrcode list-categories
```

## Commands

### `code` - Auto-generate setup code (Easiest)

Automatically generates setup code, setup ID, and MAC address:

```bash
homekitgenqrcode code -c <category> -o <output.png>
```

Options:
- `-c, --category`: HomeKit category ID (required)
- `-o, --output`: Output image file path (required)
- `-s, --setup-id`: Custom setup ID (optional, auto-generated if not provided)
- `-m, --mac`: Custom MAC address (optional, auto-generated if not provided)

### `generate` - Manual generation

Generate with all parameters manually specified:

```bash
homekitgenqrcode generate -c <category> -p <password> -s <setup-id> -m <mac> -o <output.png>
```

Options:
- `-c, --category`: HomeKit category ID (required)
- `-p, --password`: Setup password in format XXX-XX-XXX (required)
- `-s, --setup-id`: Setup ID: 4 alphanumeric characters (0-9, A-Z) (required)
- `-m, --mac`: MAC address: 12 hexadecimal characters (required)
- `-o, --output`: Output image file path (required)

### `list-categories` - List available categories

Display all available HomeKit device categories:

```bash
homekitgenqrcode list-categories
```

## Examples

```bash
# Generate with completely automatic values
homekitgenqrcode code -c 5 -o example.png

# Generate with custom setup ID and MAC
homekitgenqrcode code -c 5 -o example.png -s HSPN -m 30AEA40506A0

# Generate in a specific directory (will be created automatically)
homekitgenqrcode code -c 5 -o output/example.png

# Using long flags
homekitgenqrcode generate --category 5 --password "482-91-573" --setup-id "HSPN" --mac "30AEA40506A0" --output example.png
```

## Requirements

- Go 1.24.0 or later
- Assets folder with required fonts and template image

## License

This project is open source and available for use.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Credits

This project was created for [HomeSpan](https://github.com/HomeSpan/HomeSpan/) - a robust Arduino library for creating ESP32-based HomeKit devices.

The original idea and concept came from [AchimPieters/esp32-homekit-qrcode](https://github.com/AchimPieters/esp32-homekit-qrcode), a Python-based QR code generator. This Go implementation provides:

- **Better Performance**: Compiled Go binaries are faster and more efficient than interpreted Python scripts
- **Easier Distribution**: Single binary executable, no Python runtime or dependencies required
- **Cross-Platform**: Works seamlessly on Windows, macOS, and Linux without additional setup
- **Enhanced CLI**: Modern command-line interface using Cobra with auto-completion support

## Related Projects

- [HomeSpan](https://github.com/HomeSpan/HomeSpan/) - HomeKit Library for the Arduino-ESP32
- [esp32-homekit-qrcode](https://github.com/AchimPieters/esp32-homekit-qrcode) - Original Python implementation

