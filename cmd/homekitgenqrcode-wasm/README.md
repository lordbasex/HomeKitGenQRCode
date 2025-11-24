# HomeKit QR Code Generator - WebAssembly (WASM) Version

A web-based version of the HomeKit QR Code Generator that runs entirely in the browser using WebAssembly. Generate professional HomeKit setup labels with QR codes directly from your web browser - no server-side processing required!

## Features

- 🚀 **100% Client-Side**: All processing happens in your browser using WebAssembly
- 🎨 **Modern Web Interface**: Clean, responsive UI with real-time validation
- 🔍 **Debug Mode**: CLI-like logging in browser console for troubleshooting
- 📥 **Download Support**: One-click download of generated QR code labels
- ✅ **Full Validation**: Same validation logic as the CLI version
- 🎲 **Random Generation**: Auto-generate setup codes, setup IDs, and MAC addresses
- 📱 **Responsive Design**: Works on desktop, tablet, and mobile devices

## Requirements

- A modern web browser with WebAssembly support:
  - Chrome/Edge 57+
  - Firefox 52+
  - Safari 11+
- Node.js and npm (for installing http-server)

## Installation

### Step 1: Install http-server

Install `http-server` globally using npm:

```bash
sudo npm install http-server -g
```

**Note**: On Windows, you may not need `sudo`. On macOS/Linux, `sudo` is required for global installation.

### Step 2: Build WASM (if not already built)

If you haven't built the WASM version yet, run:

```bash
cd /path/to/HomeKitGenQRCode
make wasm
```

This will create:
- `homekitgenqrcode.wasm` - The WebAssembly binary
- `wasm_exec.js` - Go's WebAssembly runtime
- `index.html` - The web interface

## Usage

### Starting the Server

1. Navigate to the WASM directory:

```bash
cd cmd/homekitgenqrcode-wasm
```

2. Start the HTTP server:

```bash
http-server -p 8080
```

**Alternative options:**
- Use a different port: `http-server -p 3000`
- Enable CORS (if needed): `http-server -p 8080 --cors`
- Open browser automatically: `http-server -p 8080 -o`

### Accessing the Application

Open your web browser and navigate to:

```
http://localhost:8080
```

The application will load automatically. You should see the HomeKit QR Code Generator interface.

## How to Use

### 1. Select a Category

Choose a HomeKit device category from the dropdown menu (e.g., "5: Light", "8: Switch", etc.).

### 2. Enter Setup Information

You have two options:

#### Option A: Manual Entry
- **Setup Code**: Enter in format `XXX-XX-XXX` (e.g., `613-80-755`)
- **Setup ID**: Enter 4 alphanumeric characters (0-9, A-Z) (e.g., `ABCD`)
- **MAC Address**: Enter 12 hexadecimal characters (0-9, A-F) (e.g., `AABBCCDDEEFF`)

#### Option B: Generate Random Values
Click the "🎲 Generate Random Values" button to automatically fill all fields with valid random values.

### 3. Generate QR Code

Click the "✨ Generate QR Code" button. The application will:
- Validate all inputs (same validation as CLI)
- Generate the QR code label
- Display the result below

### 4. Download the Image

Once generated, click the "⬇️ Download Image" button to save the PNG file to your computer.

The downloaded file will be named: `homekit-qrcode-{category}-{timestamp}.png`

Example: `homekit-qrcode-light-2024-11-24T14-30-25.png`

## Debug Mode

Enable debug mode to see detailed CLI-like logs in your browser's console:

1. Check the "Enable Debug Mode" checkbox
2. Open your browser's Developer Console:
   - **Chrome/Edge**: Press `F12` or `Ctrl+Shift+I` (Windows/Linux) / `Cmd+Option+I` (Mac)
   - **Firefox**: Press `F12` or `Ctrl+Shift+K` (Windows/Linux) / `Cmd+Option+K` (Mac)
   - **Safari**: Enable Developer menu in Preferences, then press `Cmd+Option+C`

3. All operations will be logged with timestamps and color-coded messages:
   - ✅ Success (green)
   - ❌ Error (red)
   - ℹ️ Info (blue)

### Example Debug Output

```
[14:30:15] ✅ WASM module loaded successfully
[14:30:15] ℹ️ Loading HomeKit categories...
[14:30:15] ✅ Loaded 32 categories
[14:30:20] ℹ️ Generating random values for category: 5
[14:30:20] ✅ Generated Setup Code: 613-80-755
[14:30:20] ✅ Generated Setup ID: ABCD
[14:30:20] ✅ Generated MAC Address: AABBCCDDEEFF
[14:30:25] ℹ️ Starting QR code generation...
[14:30:25] ✅ All inputs validated successfully
[14:30:25] ℹ️ Generating HomeKit QR code label...
[14:30:26] ✅ QR code generated successfully
[14:30:26] ✅ QR-code generated and displayed
```

## Validation

The WASM version uses the same validation logic as the CLI:

- **Category**: Must be a valid HomeKit category ID (1-32, excluding 25)
- **Password**: Must match format `XXX-XX-XXX` where X is a digit (0-9)
- **Setup ID**: Must be exactly 4 alphanumeric characters (0-9, A-Z)
- **MAC Address**: Must be exactly 12 hexadecimal characters (0-9, A-F)

Validation errors show detailed messages, including:
- Character position for invalid characters
- Expected format examples
- Specific error descriptions

## Troubleshooting

### WASM Module Not Loading

**Problem**: "Failed to load WASM" error

**Solutions**:
1. Ensure you're serving via HTTP (not `file://` protocol)
2. Check that `homekitgenqrcode.wasm` and `wasm_exec.js` are in the same directory
3. Verify browser console for specific errors
4. Try a different browser

### Port Already in Use

**Problem**: `http-server` fails to start because port 8080 is in use

**Solution**: Use a different port:
```bash
http-server -p 3000
```

### CORS Errors

**Problem**: CORS errors when loading WASM files

**Solution**: Start server with CORS enabled:
```bash
http-server -p 8080 --cors
```

### Validation Errors

**Problem**: Validation fails even with correct format

**Solutions**:
1. Enable debug mode to see detailed error messages
2. Check for extra spaces (they're automatically trimmed)
3. Ensure Setup ID and MAC are uppercase (auto-converted)
4. Verify password format matches exactly: `XXX-XX-XXX`

### Image Not Downloading

**Problem**: Download button doesn't work

**Solutions**:
1. Check browser console for errors
2. Ensure pop-up blocker isn't blocking downloads
3. Try right-clicking the image and "Save Image As..."

## File Structure

```
cmd/homekitgenqrcode-wasm/
├── README.md                 # This file
├── index.html                # Web interface
├── main.go                   # WASM source code
├── homekitgenqrcode.wasm     # Compiled WebAssembly binary
└── wasm_exec.js              # Go WebAssembly runtime
```

## Technical Details

- **WebAssembly**: Compiled from Go using `GOOS=js GOARCH=wasm`
- **Size**: ~9.7 MB WASM binary (includes all fonts and assets embedded)
- **Runtime**: Go's WebAssembly runtime (`wasm_exec.js`)
- **Image Format**: PNG, 300 DPI ready for printing
- **Validation**: Server-side validation logic running client-side

## Comparison with CLI Version

| Feature | CLI | WASM |
|---------|-----|------|
| Generate QR Code | ✅ | ✅ |
| Random Value Generation | ✅ | ✅ |
| List Categories | ✅ | ✅ |
| Full Validation | ✅ | ✅ |
| Debug Logging | ✅ | ✅ |
| Download Image | ✅ | ✅ |
| Cross-Platform | ✅ | ✅ |
| No Installation | ❌ | ✅ |
| Works Offline | ✅ | ✅* |

*After initial load, works offline if all files are cached

## Browser Compatibility

| Browser | Minimum Version | Status |
|---------|----------------|--------|
| Chrome | 57+ | ✅ Fully Supported |
| Edge | 57+ | ✅ Fully Supported |
| Firefox | 52+ | ✅ Fully Supported |
| Safari | 11+ | ✅ Fully Supported |
| Opera | 44+ | ✅ Fully Supported |

## Security Notes

- All processing happens client-side - no data is sent to any server
- Generated images are created entirely in your browser
- No tracking or analytics included
- You can host this on your own server for complete privacy

## Development

### Rebuilding WASM

If you modify `main.go`, rebuild with:

```bash
cd /path/to/HomeKitGenQRCode
make wasm
```

### Testing Locally

1. Start server: `http-server -p 8080`
2. Open browser: `http://localhost:8080`
3. Enable debug mode and check console
4. Test all features

## License

Same license as the main HomeKitGenQRCode project.

## Related

- [Main Project README](../../README.md)
- [CLI Documentation](../../README.md#usage)
- [GitHub Repository](https://github.com/lordbasex/HomeKitGenQRCode)

## Support

For issues, questions, or contributions, please visit the [main project repository](https://github.com/lordbasex/HomeKitGenQRCode).

