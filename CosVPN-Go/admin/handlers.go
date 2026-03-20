package admin

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
)

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func HandleLogin(password string, rl *RateLimiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)
		if !rl.Allow(ip) {
			jsonError(w, "too many attempts", http.StatusTooManyRequests)
			return
		}
		var req struct {
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			jsonError(w, "invalid request", http.StatusBadRequest)
			return
		}
		if req.Password != password {
			jsonError(w, "wrong password", http.StatusUnauthorized)
			return
		}
		token, _ := GenerateJWT(password)
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   86400,
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":true}`))
	}
}

func HandleStatus(wg *WgCtl) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, err := wg.Status()
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse(w, status)
	}
}

func HandleListClients(wg *WgCtl) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clients, err := wg.ListClients()
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse(w, clients)
	}
}

func HandleAddClient(wg *WgCtl, log *EventLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
			jsonError(w, "name is required", http.StatusBadRequest)
			return
		}
		client, err := wg.AddClient(req.Name)
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Add("client_add", req.Name, "Client created")

		resp := map[string]interface{}{
			"name":   client.Name,
			"ip":     client.IP,
			"config": client.Config,
		}
		if len(client.QRData) > 0 {
			resp["qr"] = base64.StdEncoding.EncodeToString(client.QRData)
		}
		jsonResponse(w, resp)
	}
}

func HandleDeleteClient(wg *WgCtl, log *EventLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		if name == "" {
			jsonError(w, "client name is required", http.StatusBadRequest)
			return
		}
		if err := wg.RemoveClient(name); err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Add("client_remove", name, "Client removed")
		jsonResponse(w, map[string]bool{"ok": true})
	}
}

func HandleClientQR(wg *WgCtl) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		if name == "" {
			jsonError(w, "client name is required", http.StatusBadRequest)
			return
		}
		qrData, err := wg.GenerateQR(name)
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "image/png")
		w.Write(qrData)
	}
}

func HandleClientConf(wg *WgCtl) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		if name == "" {
			jsonError(w, "client name is required", http.StatusBadRequest)
			return
		}
		config, err := wg.GetClientConfig(name)
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Disposition", "attachment; filename="+name+".conf")
		w.Write([]byte(config))
	}
}

func HandleGetSettings(wg *WgCtl) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		settings, err := wg.GetSettings()
		if err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse(w, settings)
	}
}

func HandleUpdateSettings(wg *WgCtl, log *EventLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var s Settings
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			jsonError(w, "invalid settings", http.StatusBadRequest)
			return
		}
		if err := wg.UpdateSettings(s); err != nil {
			jsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Add("settings", "", "Settings updated")
		jsonResponse(w, map[string]bool{"ok": true})
	}
}

func HandleLogs(log *EventLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limitStr := r.URL.Query().Get("limit")
		limit := 50
		if limitStr != "" {
			if v, err := strconv.Atoi(limitStr); err == nil && v > 0 {
				limit = v
			}
		}
		entries := log.Get(limit)
		jsonResponse(w, entries)
	}
}
