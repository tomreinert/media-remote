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
	default:
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html><html><body style="font-family:-apple-system;padding:40px;text-align:center;">
			<h1>Media Remote API</h1>
			<p>Server is running. Use the remote UI at:</p>
			<p><a href="https://tomreinert.github.io/media-remote/ui/">tomreinert.github.io/media-remote/ui/</a></p>
		</body></html>`))
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

		fmt.Println("Media Remote API started in background")
		fmt.Println("\nServer address (enter this in the remote UI):")
		if hostname != "" {
			fmt.Printf("  %s.local:%s\n", hostname, port)
		}
		fmt.Printf("  %s:%s\n", ip, port)
		fmt.Println("\nRemote UI: https://tomreinert.github.io/media-remote/ui/")
		fmt.Println("\nStop with: pkill media-remote")
		return
	}

	// Running as daemon - start server silently
	http.HandleFunc("/", corsMiddleware(handler))
	http.ListenAndServe(":"+port, nil)
}
