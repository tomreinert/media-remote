# Media Remote

A lightweight macOS remote control server for Netflix, Zattoo, and other media apps. Control playback from your phone via a web browser.

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
| `ui/index.html` | The UI - edit on GitHub or via Claude Code web |

## Updating

### UI changes (most common)
No build needed! Just:
1. Edit `ui/index.html` (via GitHub web, Claude Code, or push)
2. Wait ~1-2 min for GitHub raw CDN cache
3. On target Mac: `pkill media-remote && ./media-remote`

### Server/API changes (rare)
Only rebuild if changing `media_remote.go`:
```bash
GOARCH=amd64 go build -o media-remote-intel media_remote.go
```

## Initial Setup (target Mac)

```bash
# Build once
GOARCH=amd64 go build -o media-remote-intel media_remote.go

# Run
./media-remote

# Stop
pkill media-remote
```

## On Your Phone

Open the URL shown when server starts:
```
http://<ip>:9876
```
Add to Home Screen for fullscreen experience.

## API Endpoints

| Endpoint | Action |
|----------|--------|
| `/` | Serves the UI |
| `/ping` | Health check |
| `/toggle` | Play/Pause (spacebar) |
| `/vol-up` | Volume +10 |
| `/vol-down` | Volume -10 |
| `/mute` | Mute toggle (m key) |
| `/fullscreen` | Fullscreen toggle (f key) |
| `/ch-up` | Next channel (l key) - Zattoo |
| `/ch-down` | Previous channel (j key) - Zattoo |
| `/ch-list` | Channel list (h key) - Zattoo |
| `/escape` | Escape key |
| `/enter` | Enter key |

## Requirements

- macOS with Accessibility permissions
- Internet on target Mac (for fetching UI from GitHub)
- Go 1.21+ only if rebuilding server
