package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

const port = "9876"
const uiURL = "https://raw.githubusercontent.com/tomreinert/media-remote/main/ui/index.html"

var cachedUI string

func sendKey(key string) {
	exec.Command("osascript", "-e", fmt.Sprintf(`tell application "System Events" to keystroke "%s"`, key)).Run()
}

func sendSpace() {
	exec.Command("osascript", "-e", `tell application "System Events" to keystroke space`).Run()
}

func volumeUp() {
	exec.Command("osascript", "-e", `set volume output volume ((output volume of (get volume settings)) + 10)`).Run()
}

func volumeDown() {
	exec.Command("osascript", "-e", `set volume output volume ((output volume of (get volume settings)) - 10)`).Run()
}

func getLocalHostname() string {
	out, err := exec.Command("scutil", "--get", "LocalHostName").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "localhost"
	}
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).IP.String()
}

func fetchUI() string {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(uiURL)
	if err != nil {
		fmt.Println("Failed to fetch UI:", err)
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println("Failed to fetch UI: status", resp.StatusCode)
		return ""
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read UI response:", err)
		return ""
	}
	fmt.Println("UI fetched successfully from GitHub")
	return string(body)
}

const fallbackUI = `<!DOCTYPE html>
<html>
<head>
    <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no">
    <meta name="apple-mobile-web-app-capable" content="yes">
    <title>Remote</title>
    <style>
        * { box-sizing: border-box; touch-action: manipulation; }
        html, body { height: 100%; width: 100%; margin: 0; padding: 0; overflow: hidden; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, sans-serif;
            display: flex; flex-direction: column; justify-content: center; align-items: center;
            gap: 4vw; background: #000; padding: 20px;
        }
        button {
            font-size: 6vw; padding: 6vw 8vw; border-radius: 4vw; border: none;
            background: #333; color: white; cursor: pointer;
            -webkit-tap-highlight-color: transparent; user-select: none;
            touch-action: manipulation; transition: transform 0.1s ease, background 0.1s ease;
        }
        button:active { background: #444; transform: scale(0.97); }
        .play-pause { font-size: 8vw; padding: 12vw 16vw; }
        .row { display: flex; gap: 3vw; width: 100%; justify-content: center; }
        .row button { flex: 1; max-width: 45vw; }
        .row-3 button { max-width: 30vw; padding: 5vw 2vw; font-size: 4.5vw; }
        .section-label { font-size: 3vw; color: #666; margin-top: 2vw; margin-bottom: -1vw; }
        .title { font-size: 5vw; color: #fff; font-weight: 600; margin-bottom: auto; }
    </style>
</head>
<body>
    <div class="title">Remote</div>
    <button class="play-pause" onclick="fetch('/toggle')">Play / Pause</button>
    <div class="row">
        <button onclick="fetch('/vol-down')">Vol -</button>
        <button onclick="fetch('/mute')">Mute</button>
        <button onclick="fetch('/vol-up')">Vol +</button>
    </div>
    <div class="row">
        <button onclick="fetch('/fullscreen')">Fullscreen</button>
    </div>
    <div class="section-label">Zattoo</div>
    <div class="row row-3">
        <button onclick="fetch('/ch-down')">CH -</button>
        <button onclick="fetch('/ch-list')">List</button>
        <button onclick="fetch('/ch-up')">CH +</button>
    </div>
</body>
</html>`

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/ping":
		w.Write([]byte("pong"))
	case "/setup":
		fmt.Println("Running setup to trigger permissions...")
		exec.Command("osascript", "-e", `tell application "Google Chrome" to get name`).Run()
		exec.Command("osascript", "-e", `tell application "Safari" to get name`).Run()
		exec.Command("osascript", "-e", `tell application "System Events" to get name`).Run()
		w.Write([]byte("OK"))
	case "/toggle":
		sendSpace()
		fmt.Println("Toggle received")
		w.Write([]byte("OK"))
	case "/vol-up":
		volumeUp()
		fmt.Println("Volume up")
		w.Write([]byte("OK"))
	case "/vol-down":
		volumeDown()
		fmt.Println("Volume down")
		w.Write([]byte("OK"))
	case "/mute":
		sendKey("m")
		fmt.Println("Mute toggle")
		w.Write([]byte("OK"))
	case "/fullscreen":
		sendKey("f")
		fmt.Println("Fullscreen toggle")
		w.Write([]byte("OK"))
	case "/ch-up":
		sendKey("l")
		fmt.Println("Channel up")
		w.Write([]byte("OK"))
	case "/ch-down":
		sendKey("j")
		fmt.Println("Channel down")
		w.Write([]byte("OK"))
	case "/ch-list":
		sendKey("h")
		fmt.Println("Channel list toggle")
		w.Write([]byte("OK"))
	case "/escape":
		exec.Command("osascript", "-e", `tell application "System Events" to key code 53`).Run()
		fmt.Println("Escape")
		w.Write([]byte("OK"))
	case "/enter":
		exec.Command("osascript", "-e", `tell application "System Events" to key code 36`).Run()
		fmt.Println("Enter")
		w.Write([]byte("OK"))
	case "/":
		w.Header().Set("Content-Type", "text/html")
		if cachedUI != "" {
			w.Write([]byte(cachedUI))
		} else {
			w.Write([]byte(fallbackUI))
		}
	default:
		http.NotFound(w, r)
	}
}

func main() {
	// Check if running as daemon
	isDaemon := len(os.Args) > 1 && os.Args[1] == "--daemon"

	if !isDaemon {
		// Fork to background and exit
		hostname := getLocalHostname()
		ip := getLocalIP()

		cmd := exec.Command(os.Args[0], "--daemon")
		cmd.Start()

		fmt.Println("Media Remote started in background")
		fmt.Println("\nOpen on your phone:")
		if hostname != "" {
			fmt.Printf("  http://%s.local:%s\n", hostname, port)
		}
		fmt.Printf("  http://%s:%s\n", ip, port)
		fmt.Println("\nStop with: pkill media-remote")
		return
	}

	// Running as daemon - fetch UI and start server
	cachedUI = fetchUI()
	http.HandleFunc("/", corsMiddleware(handler))
	http.ListenAndServe(":"+port, nil)
}
