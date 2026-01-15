package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const port = "9876"

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

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	switch r.URL.Path {
	case "/setup":
		// Trigger permission dialogs - run this from the Mac itself
		fmt.Println("Running setup to trigger permissions...")
		exec.Command("osascript", "-e", `tell application "Google Chrome" to get name`).Run()
		exec.Command("osascript", "-e", `tell application "Safari" to get name`).Run()
		exec.Command("osascript", "-e", `tell application "System Events" to get name`).Run()
		w.Write([]byte(`<!DOCTYPE html><html><body style="font-family:-apple-system;padding:40px;">
			<h1>Setup complete</h1>
			<p>Approve any permission dialogs that appeared, then use the remote from your phone.</p>
			<p><a href="/">Go to remote</a></p>
		</body></html>`))
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
	case "/":
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
    <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no">
    <meta name="apple-mobile-web-app-capable" content="yes">
    <meta name="apple-mobile-web-app-status-bar-style" content="black-translucent">
    <title>Remote</title>
    <style>
        * {
            box-sizing: border-box;
            touch-action: manipulation;
        }
        html, body {
            height: 100%;
            width: 100%;
            margin: 0;
            padding: 0;
            overflow: hidden;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, sans-serif;
            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;
            gap: 4vw;
            background: #000;
            padding: 20px;
        }
        button {
            font-size: 6vw;
            padding: 6vw 8vw;
            border-radius: 4vw;
            border: none;
            background: #333;
            color: white;
            cursor: pointer;
            -webkit-tap-highlight-color: transparent;
            user-select: none;
            touch-action: manipulation;
            transition: transform 0.1s ease, background 0.1s ease;
        }
        button:active,
        button.pressed {
            background: #444;
            transform: scale(0.97);
        }
        .play-pause {
            font-size: 8vw;
            padding: 12vw 16vw;
        }
        .volume-row {
            display: flex;
            gap: 4vw;
            width: 100%;
            justify-content: center;
        }
        .volume-row button {
            flex: 1;
            max-width: 40vw;
        }
    </style>
</head>
<body>
    <button class="play-pause" onclick="fetch('/toggle')">Play / Pause</button>
    <div class="volume-row">
        <button onclick="fetch('/vol-down')">Vol -</button>
        <button onclick="fetch('/vol-up')">Vol +</button>
    </div>
</body>
</html>`))
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
		if hostname != "" {
			fmt.Printf("URL: http://%s.local:%s\n", hostname, port)
		}
		fmt.Printf("URL: http://%s:%s\n", ip, port)
		fmt.Println("\nStop with: pkill media-remote")
		return
	}

	// Running as daemon - start server silently
	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+port, nil)
}
