#!/usr/bin/env python3
"""Simple media remote control server - control browser playback from your phone."""

import subprocess
from http.server import HTTPServer, BaseHTTPRequestHandler
import socket

def get_local_ip():
    """Get the local IP address for this machine."""
    s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    try:
        s.connect(("8.8.8.8", 80))
        return s.getsockname()[0]
    finally:
        s.close()

def send_space():
    """Send spacebar keystroke to the frontmost app using AppleScript."""
    script = 'tell application "System Events" to keystroke space'
    subprocess.run(["osascript", "-e", script])

class RemoteHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path == "/toggle" or self.path == "/":
            send_space()
            self.send_response(200)
            self.send_header("Content-type", "text/html")
            self.send_header("Access-Control-Allow-Origin", "*")
            self.end_headers()
            self.wfile.write(b"""
                <html><body style="font-family: -apple-system; display: flex;
                justify-content: center; align-items: center; height: 100vh; margin: 0;">
                <button onclick="fetch('/toggle')" style="font-size: 48px; padding: 40px 80px;
                border-radius: 20px; border: none; background: #007AFF; color: white;">
                Play/Pause</button>
                </body></html>
            """)
        else:
            self.send_response(404)
            self.end_headers()

    def log_message(self, format, *args):
        print(f"Request: {args[0]}")

if __name__ == "__main__":
    port = 8080
    ip = get_local_ip()
    print(f"\n🎬 Media Remote Server Running!")
    print(f"\nOpen this URL on your iPhone:")
    print(f"  http://{ip}:{port}")
    print(f"\nOr use this in an iOS Shortcut:")
    print(f"  http://{ip}:{port}/toggle")
    print(f"\nPress Ctrl+C to stop\n")

    server = HTTPServer(("0.0.0.0", port), RemoteHandler)
    server.serve_forever()
