# Media Remote

A lightweight macOS remote control server that lets you control media playback from any device on your network via a web browser.

## Architecture

```
┌─────────────────┐         ┌──────────────────┐         ┌─────────────────┐
│  Your Phone     │────────▶│  Go Server       │────────▶│  macOS System   │
│  (browser)      │         │  (port 9876)     │         │  (osascript)    │
└─────────────────┘         └────────┬─────────┘         └─────────────────┘
                                     │
                            On startup, fetches UI from:
                                     │
                                     ▼
                     ┌───────────────────────────────┐
                     │  GitHub Raw                   │
                     │  raw.githubusercontent.com/   │
                     │  .../main/ui/index.html       │
                     └───────────────────────────────┘
```

**How it works:**
1. Go server starts and fetches `ui/index.html` from GitHub
2. Serves that UI at `/` over HTTP
3. UI makes API calls to the same server (relative paths)
4. If GitHub fetch fails (offline), falls back to embedded UI

## Files

| File | Purpose |
|------|---------|
| `media_remote.go` | Go server - fetches UI from GitHub, serves API |
| `ui/index.html` | The UI - edit this on GitHub to update remotely |

## Updating the UI

1. Edit `ui/index.html` on GitHub (or push changes)
2. Wait ~1-2 min for GitHub's raw CDN cache to clear
3. On target Mac: `pkill media-remote && ./media-remote`

No rebuild needed - the server fetches fresh UI on each restart.

## Initial Setup (target Mac)

```bash
# Build (one-time, or after changing media_remote.go)
GOARCH=amd64 go build -o media-remote-intel media_remote.go

# Run
./media-remote

# Stop
pkill media-remote
```

## On Your Phone

Just open the URL shown when the server starts:
```
http://<ip>:9876
```
Add to Home Screen for fullscreen experience.

## API Endpoints

| Endpoint | Action |
|----------|--------|
| `/` | Serves the UI |
| `/ping` | Health check (returns "pong") |
| `/setup` | Trigger macOS permission dialogs |
| `/toggle` | Play/Pause (spacebar keystroke) |
| `/vol-up` | Volume +10 |
| `/vol-down` | Volume -10 |

## Requirements

- macOS (uses AppleScript for system control)
- Accessibility permissions (System Settings → Privacy & Security → Accessibility)
- Go 1.21+ for building (not needed if using pre-built binary)
- Internet connection on target Mac (for fetching UI from GitHub)
