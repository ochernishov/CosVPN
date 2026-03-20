package admin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type ServerStatus struct {
	Uptime        string `json:"uptime"`
	CPU           int    `json:"cpu"`
	RAM           int    `json:"ram"`
	Disk          int    `json:"disk"`
	VPNUp         bool   `json:"vpnUp"`
	Port          int    `json:"port"`
	PublicIP      string `json:"publicIP"`
	TotalClients  int    `json:"totalClients"`
	OnlineClients int    `json:"onlineClients"`
	ObfsMode      string `json:"obfsMode"`
}

type Client struct {
	Name          string `json:"name"`
	IP            string `json:"ip"`
	Online        bool   `json:"online"`
	LastHandshake string `json:"lastHandshake"`
	TransferUp    int64  `json:"transferUp"`
	TransferDown  int64  `json:"transferDown"`
	PublicKey     string `json:"publicKey"`
}

type NewClient struct {
	Name   string `json:"name"`
	IP     string `json:"ip"`
	Config string `json:"config"`
	QRData []byte `json:"-"`
}

type Settings struct {
	ObfsMode   string   `json:"obfuscationMode"`
	ObfsKeySet bool     `json:"obfuscationKeySet"`
	Port       int      `json:"listenPort"`
	DNS        []string `json:"dns"`
	MTU        int      `json:"mtu"`
	Subnet     string   `json:"subnet"`
}

type WgCtl struct {
	configDir string
}

func NewWgCtl(configDir string) *WgCtl {
	return &WgCtl{configDir: configDir}
}

func runCmd(cmd string) (string, error) {
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func (w *WgCtl) Status() (ServerStatus, error) {
	var s ServerStatus

	// Check VPN status
	wgOut, err := runCmd("wg show wg0 listen-port 2>/dev/null")
	if err == nil && wgOut != "" {
		s.VPNUp = true
		port, _ := strconv.Atoi(wgOut)
		s.Port = port
	}

	// Uptime
	uptimeOut, _ := runCmd("uptime -p 2>/dev/null || uptime")
	s.Uptime = uptimeOut

	// CPU usage
	cpuOut, _ := runCmd("top -bn1 2>/dev/null | grep 'Cpu(s)' | awk '{print int($2+$4)}'")
	cpuVal, _ := strconv.Atoi(cpuOut)
	s.CPU = cpuVal

	// RAM usage
	ramOut, _ := runCmd("free -m 2>/dev/null | awk '/Mem:/{print int($3/$2*100)}'")
	ramVal, _ := strconv.Atoi(ramOut)
	s.RAM = ramVal

	// Disk usage
	diskOut, _ := runCmd("df -h / 2>/dev/null | awk 'NR==2{print int($5)}'")
	diskVal, _ := strconv.Atoi(diskOut)
	s.Disk = diskVal

	// Public IP
	ipOut, _ := runCmd("curl -s --connect-timeout 3 ifconfig.me 2>/dev/null")
	s.PublicIP = ipOut

	// Client counts
	clients, err := w.ListClients()
	if err == nil {
		s.TotalClients = len(clients)
		for _, c := range clients {
			if c.Online {
				s.OnlineClients++
			}
		}
	}

	// Obfuscation mode
	settings, err := w.GetSettings()
	if err == nil {
		s.ObfsMode = settings.ObfsMode
	}

	return s, nil
}

func (w *WgCtl) ListClients() ([]Client, error) {
	// Get WireGuard dump
	dumpOut, err := runCmd("wg show wg0 dump 2>/dev/null")
	if err != nil {
		return []Client{}, nil
	}

	lines := strings.Split(dumpOut, "\n")
	if len(lines) < 2 {
		return []Client{}, nil
	}

	// Build name map from client directories
	nameMap := make(map[string]string)
	clientsDir := filepath.Join(w.configDir, "clients")
	entries, _ := os.ReadDir(clientsDir)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		confPath := filepath.Join(clientsDir, entry.Name(), entry.Name()+".conf")
		confData, err := os.ReadFile(confPath)
		if err != nil {
			continue
		}
		// Extract client's public key by computing it from private key
		for _, line := range strings.Split(string(confData), "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "PrivateKey") {
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					privKey := strings.TrimSpace(strings.Join(parts[1:], "="))
					pubKey, err := runCmd(fmt.Sprintf("echo '%s' | wg pubkey 2>/dev/null", privKey))
					if err == nil && pubKey != "" {
						nameMap[pubKey] = entry.Name()
					}
				}
			}
		}
	}

	var clients []Client
	// Skip first line (interface line), parse peer lines
	for _, line := range lines[1:] {
		if strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) < 8 {
			continue
		}

		pubKey := fields[0]
		// fields[1] = preshared key
		// fields[2] = endpoint
		allowedIPs := fields[3]
		latestHandshake, _ := strconv.ParseInt(fields[4], 10, 64)
		transferRx, _ := strconv.ParseInt(fields[5], 10, 64)
		transferTx, _ := strconv.ParseInt(fields[6], 10, 64)

		name := nameMap[pubKey]
		if name == "" {
			name = pubKey[:8] + "..."
		}

		ip := strings.Split(allowedIPs, "/")[0]

		online := false
		lastHandshake := "never"
		if latestHandshake > 0 {
			hsTime := time.Unix(latestHandshake, 0)
			if time.Since(hsTime) < 3*time.Minute {
				online = true
			}
			lastHandshake = hsTime.Format(time.RFC3339)
		}

		clients = append(clients, Client{
			Name:          name,
			IP:            ip,
			Online:        online,
			LastHandshake: lastHandshake,
			TransferUp:    transferTx,
			TransferDown:  transferRx,
			PublicKey:     pubKey,
		})
	}

	return clients, nil
}

func (w *WgCtl) AddClient(name string) (*NewClient, error) {
	addScript := filepath.Join(w.configDir, "add-client.sh")
	out, err := runCmd(fmt.Sprintf("bash '%s' '%s'", addScript, name))
	if err != nil {
		return nil, fmt.Errorf("add-client.sh failed: %s: %w", out, err)
	}

	confPath := filepath.Join(w.configDir, "clients", name, name+".conf")
	confData, err := os.ReadFile(confPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read client config: %w", err)
	}

	config := string(confData)

	// Extract IP from config
	ip := ""
	for _, line := range strings.Split(config, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Address") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				ip = strings.TrimSpace(parts[1])
			}
		}
	}

	// Generate QR
	qrData, _ := runCmd(fmt.Sprintf("qrencode -o - -t PNG < '%s' 2>/dev/null", confPath))

	return &NewClient{
		Name:   name,
		IP:     ip,
		Config: config,
		QRData: []byte(qrData),
	}, nil
}

func (w *WgCtl) RemoveClient(name string) error {
	clientDir := filepath.Join(w.configDir, "clients", name)
	confPath := filepath.Join(clientDir, name+".conf")

	confData, err := os.ReadFile(confPath)
	if err != nil {
		return fmt.Errorf("client config not found: %w", err)
	}

	// Extract private key and compute public key
	var pubKey string
	for _, line := range strings.Split(string(confData), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "PrivateKey") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				privKey := strings.TrimSpace(strings.Join(parts[1:], "="))
				pubKey, _ = runCmd(fmt.Sprintf("echo '%s' | wg pubkey 2>/dev/null", privKey))
			}
		}
	}

	if pubKey != "" {
		runCmd(fmt.Sprintf("wg set wg0 peer '%s' remove 2>/dev/null", pubKey))
	}

	if err := os.RemoveAll(clientDir); err != nil {
		return fmt.Errorf("failed to remove client directory: %w", err)
	}

	return nil
}

func (w *WgCtl) GetClientConfig(name string) (string, error) {
	confPath := filepath.Join(w.configDir, "clients", name, name+".conf")
	data, err := os.ReadFile(confPath)
	if err != nil {
		return "", fmt.Errorf("config not found: %w", err)
	}
	return string(data), nil
}

func (w *WgCtl) GenerateQR(name string) ([]byte, error) {
	confPath := filepath.Join(w.configDir, "clients", name, name+".conf")
	out, err := exec.Command("bash", "-c", fmt.Sprintf("qrencode -o - -t PNG < '%s'", confPath)).Output()
	if err != nil {
		return nil, fmt.Errorf("qrencode failed: %w", err)
	}
	return out, nil
}

func (w *WgCtl) GetSettings() (Settings, error) {
	var s Settings

	// Parse wg0.conf
	wgConfPath := filepath.Join(w.configDir, "wg0.conf")
	wgConf, err := os.ReadFile(wgConfPath)
	if err == nil {
		for _, line := range strings.Split(string(wgConf), "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "ListenPort") {
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					s.Port, _ = strconv.Atoi(strings.TrimSpace(parts[1]))
				}
			}
		}
	}

	// Parse cosvpn.conf for obfuscation settings
	cosConf, err := os.ReadFile(filepath.Join(w.configDir, "cosvpn.conf"))
	if err == nil {
		for _, line := range strings.Split(string(cosConf), "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "ObfuscationMode") || strings.HasPrefix(line, "obfuscation_mode") {
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					s.ObfsMode = strings.TrimSpace(parts[1])
				}
			}
			if strings.HasPrefix(line, "ObfuscationKey") || strings.HasPrefix(line, "obfuscation_key") {
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 && strings.TrimSpace(parts[1]) != "" {
					s.ObfsKeySet = true
				}
			}
			if strings.HasPrefix(line, "DNS") || strings.HasPrefix(line, "dns") {
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					dnsStr := strings.TrimSpace(parts[1])
					for _, d := range strings.Split(dnsStr, ",") {
						d = strings.TrimSpace(d)
						if d != "" {
							s.DNS = append(s.DNS, d)
						}
					}
				}
			}
			if strings.HasPrefix(line, "MTU") || strings.HasPrefix(line, "mtu") {
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					s.MTU, _ = strconv.Atoi(strings.TrimSpace(parts[1]))
				}
			}
			if strings.HasPrefix(line, "Subnet") || strings.HasPrefix(line, "subnet") {
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					s.Subnet = strings.TrimSpace(parts[1])
				}
			}
		}
	}

	if s.DNS == nil {
		s.DNS = []string{}
	}

	return s, nil
}

func (w *WgCtl) UpdateSettings(s Settings) error {
	// Update cosvpn.conf
	cosConfPath := filepath.Join(w.configDir, "cosvpn.conf")
	var lines []string

	existing, err := os.ReadFile(cosConfPath)
	if err == nil {
		lines = strings.Split(string(existing), "\n")
	}

	updated := make(map[string]bool)
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "ObfuscationMode") || strings.HasPrefix(trimmed, "obfuscation_mode") {
			lines[i] = fmt.Sprintf("ObfuscationMode = %s", s.ObfsMode)
			updated["obfs"] = true
		}
		if strings.HasPrefix(trimmed, "DNS") || strings.HasPrefix(trimmed, "dns") {
			lines[i] = fmt.Sprintf("DNS = %s", strings.Join(s.DNS, ", "))
			updated["dns"] = true
		}
		if strings.HasPrefix(trimmed, "MTU") || strings.HasPrefix(trimmed, "mtu") {
			lines[i] = fmt.Sprintf("MTU = %d", s.MTU)
			updated["mtu"] = true
		}
		if strings.HasPrefix(trimmed, "Subnet") || strings.HasPrefix(trimmed, "subnet") {
			lines[i] = fmt.Sprintf("Subnet = %s", s.Subnet)
			updated["subnet"] = true
		}
	}

	if !updated["obfs"] {
		lines = append(lines, fmt.Sprintf("ObfuscationMode = %s", s.ObfsMode))
	}
	if !updated["dns"] {
		lines = append(lines, fmt.Sprintf("DNS = %s", strings.Join(s.DNS, ", ")))
	}
	if !updated["mtu"] {
		lines = append(lines, fmt.Sprintf("MTU = %d", s.MTU))
	}
	if !updated["subnet"] {
		lines = append(lines, fmt.Sprintf("Subnet = %s", s.Subnet))
	}

	if err := os.WriteFile(cosConfPath, []byte(strings.Join(lines, "\n")), 0600); err != nil {
		return fmt.Errorf("failed to write cosvpn.conf: %w", err)
	}

	// Update ListenPort in wg0.conf if changed
	if s.Port > 0 {
		wgConfPath := filepath.Join(w.configDir, "wg0.conf")
		wgData, err := os.ReadFile(wgConfPath)
		if err == nil {
			wgLines := strings.Split(string(wgData), "\n")
			for i, line := range wgLines {
				if strings.HasPrefix(strings.TrimSpace(line), "ListenPort") {
					wgLines[i] = fmt.Sprintf("ListenPort = %d", s.Port)
				}
			}
			os.WriteFile(wgConfPath, []byte(strings.Join(wgLines, "\n")), 0600)
		}
	}

	// Apply WireGuard config
	runCmd(fmt.Sprintf("wg syncconf wg0 <(wg-quick strip '%s/wg0.conf') 2>/dev/null", w.configDir))

	return nil
}
