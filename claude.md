# Media Remote

A lightweight macOS remote control server that lets you control media playback from any device on your network via a web browser.

## Architecture

- **Server**: Go HTTP server (port 9876) running on macOS
- **Client**: Web UI accessible from any browser (optimized for mobile/iPhone)
- **System Integration**: Uses `osascript` (AppleScript) for macOS system control

## Files

| File | Purpose |
|------|---------|
| `media_remote.go` | Main Go server with embedded HTML UI |
| `media_remote.py` | Python alternative (simpler, port 8080) |
| `media-remote` | Compiled binary (ARM64/Apple Silicon) |
| `media-remote-intel` | Compiled binary (x86_64/Intel) |

## Build

```bash
# ARM64 (Apple Silicon)
go build -o media-remote media_remote.go

# Intel x86_64
GOARCH=amd64 go build -o media-remote-intel media_remote.go
```

## Run

```bash
./media-remote          # Starts server, forks to background
pkill media-remote      # Stop the server
```

Access at `http://<hostname>.local:9876` or `http://<ip>:9876`

## API Endpoints

| Endpoint | Action |
|----------|--------|
| `/` | Web UI |
| `/setup` | Trigger macOS permission dialogs |
| `/toggle` | Play/Pause (spacebar) |
| `/vol-up` | Volume +10 |
| `/vol-down` | Volume -10 |

## Requirements

- macOS (uses AppleScript for system control)
- Accessibility permissions (System Settings → Privacy & Security → Accessibility)
- Go 1.21+ for building (binaries included for convenience)
