package admin

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"
)

//go:embed static/*
var staticFiles embed.FS

func StartServer(addr, password, wgConfigDir string) {
	wg := NewWgCtl(wgConfigDir)
	eventLog := NewEventLogger(100)
	rl := NewRateLimiter(5, time.Minute)

	mux := http.NewServeMux()

	// Login (public)
	mux.HandleFunc("POST /api/login", HandleLogin(password, rl))

	// Protected API
	api := http.NewServeMux()
	api.HandleFunc("GET /api/status", HandleStatus(wg))
	api.HandleFunc("GET /api/clients", HandleListClients(wg))
	api.HandleFunc("POST /api/clients", HandleAddClient(wg, eventLog))
	api.HandleFunc("DELETE /api/clients/{name}", HandleDeleteClient(wg, eventLog))
	api.HandleFunc("GET /api/clients/{name}/qr", HandleClientQR(wg))
	api.HandleFunc("GET /api/clients/{name}/conf", HandleClientConf(wg))
	api.HandleFunc("GET /api/settings", HandleGetSettings(wg))
	api.HandleFunc("PUT /api/settings", HandleUpdateSettings(wg, eventLog))
	api.HandleFunc("GET /api/logs", HandleLogs(eventLog))

	mux.Handle("/api/", AuthMiddleware(api, password))

	// Static files (fallback)
	staticFS, _ := fs.Sub(staticFiles, "static")
	mux.Handle("/", http.FileServer(http.FS(staticFS)))

	eventLog.Add("settings", "", "Dashboard started")

	// Check for TLS certificates
	certFile := os.Getenv("COSVPN_TLS_CERT")
	keyFile := os.Getenv("COSVPN_TLS_KEY")

	if certFile != "" && keyFile != "" {
		log.Printf("CosVPN Dashboard starting on %s (HTTPS)", addr)
		if err := http.ListenAndServeTLS(addr, certFile, keyFile, mux); err != nil {
			log.Printf("Dashboard HTTPS error: %v", err)
		}
	} else {
		log.Printf("CosVPN Dashboard starting on %s (HTTP)", addr)
		if err := http.ListenAndServe(addr, mux); err != nil {
			log.Printf("Dashboard error: %v", err)
		}
	}
}
