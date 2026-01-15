# Media Remote

A lightweight macOS remote control server that lets you control media playback from any device on your network via a web browser.

## Architecture

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│  Your Phone     │────▶│  Go Server       │────▶│  macOS System   │
│  (browser)      │     │  (API only)      │     │  (osascript)    │
└────────┬────────┘     └──────────────────┘     └─────────────────┘
         │
         ▼
┌─────────────────┐
│  GitHub Pages   │  ← UI hosted here, update by pushing to repo
│  ui/index.html  │
└─────────────────┘
```

- **Server**: Go HTTP API server (port 9876) running on macOS target
- **UI**: Hosted on GitHub Pages - update by pushing, no rebuild needed
- **System Integration**: Uses `osascript` (AppleScript) for macOS control

## Files

| File | Purpose |
|------|---------|
| `media_remote.go` | Go API server (no embedded UI) |
| `media_remote.py` | Python alternative (simpler, port 8080) |
| `ui/index.html` | Remote UI (hosted on GitHub Pages) |

## Usage

### On the target Mac (one-time setup):
```bash
go build -o media-remote media_remote.go
./media-remote
# Note the server address shown (e.g., 192.168.1.100:9876)
```

### On your phone:
1. Open https://tomreinert.github.io/media-remote/ui/
2. Enter the server address from above
3. (Optional) Add to Home Screen for fullscreen experience

## Updating the UI

Just push changes to `ui/index.html` - GitHub Pages auto-deploys. No need to touch the target Mac.

## API Endpoints

| Endpoint | Action |
|----------|--------|
| `/ping` | Health check (returns "pong") |
| `/setup` | Trigger macOS permission dialogs |
| `/toggle` | Play/Pause (spacebar) |
| `/vol-up` | Volume +10 |
| `/vol-down` | Volume -10 |

## Build (both architectures)

```bash
# ARM64 (Apple Silicon)
go build -o media-remote media_remote.go

# Intel x86_64
GOARCH=amd64 go build -o media-remote-intel media_remote.go
```

## Requirements

- macOS (uses AppleScript for system control)
- Accessibility permissions (System Settings → Privacy & Security → Accessibility)
- Go 1.21+ for building
