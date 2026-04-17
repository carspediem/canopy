package contract

// meshpulse_server.go — lightweight HTTP server for the MeshPulse UI.
//
// State is served from the in-memory globalCache (see meshpulse_cache.go),
// which is updated by DeliverTx handlers. This avoids calling StateRead
// outside of FSM hook callbacks (which the Canopy protocol does not allow).

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// StartUIServer starts the MeshPulse HTTP server on :8080.
func StartUIServer(p *Plugin, cfg Config) {
	s := &meshpulseServer{plugin: p, cfg: cfg}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.serveUI)
	mux.HandleFunc("/api/stats", s.handleStats)
	mux.HandleFunc("/api/feed", s.handleFeed)
	mux.HandleFunc("/api/measurements", s.handleFeed) // alias for /api/feed
	mux.HandleFunc("/api/leaderboard", s.handleLeaderboard)
	mux.HandleFunc("/api/contributor", s.handleContributor)
	mux.HandleFunc("/api/submit-measurement", s.handleSubmitMeasurement)

	log.Printf("🌐 MeshPulse UI → http://localhost:8080")
	go func() {
		if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Printf("UI server error: %v", err)
		}
	}()
}

type meshpulseServer struct {
	plugin *Plugin
	cfg    Config
}

// ─── HTTP helpers ─────────────────────────────────────────────────────────────

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func hexAddr(addr []byte) string { return hex.EncodeToString(addr) }

// ─── Handlers ─────────────────────────────────────────────────────────────────

func (s *meshpulseServer) serveUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, meshpulseHTML)
}

func (s *meshpulseServer) handleStats(w http.ResponseWriter, _ *http.Request) {
	stats := globalCache.Stats()
	writeJSON(w, stats)
}

func (s *meshpulseServer) handleFeed(w http.ResponseWriter, _ *http.Request) {
	type measurementJSON struct {
		ID        uint64 `json:"id"`
		Address   string `json:"address"`
		Ping      uint64 `json:"ping"`
		Download  uint64 `json:"download"`
		Upload    uint64 `json:"upload"`
		ISP       string `json:"isp"`
		Region    string `json:"region"`
		Timestamp uint64 `json:"timestamp"`
	}
	items := globalCache.Feed(50)
	feed := make([]measurementJSON, len(items))
	for i, m := range items {
		feed[i] = measurementJSON{
			ID:        m.ID,
			Address:   hexAddr(m.Address),
			Ping:      m.Ping,
			Download:  m.Download,
			Upload:    m.Upload,
			ISP:       m.ISP,
			Region:    m.Region,
			Timestamp: m.Timestamp,
		}
	}
	writeJSON(w, feed)
}

func (s *meshpulseServer) handleLeaderboard(w http.ResponseWriter, _ *http.Request) {
	type entry struct {
		Address           string `json:"address"`
		TotalMeasurements uint64 `json:"totalMeasurements"`
		TokenBalance      uint64 `json:"tokenBalance"`
	}
	top := globalCache.Leaderboard(20)
	result := make([]entry, len(top))
	for i, c := range top {
		result[i] = entry{
			Address:           hexAddr(c.Address),
			TotalMeasurements: c.TotalMeasurements,
			TokenBalance:      c.TokenBalance,
		}
	}
	writeJSON(w, result)
}

func (s *meshpulseServer) handleContributor(w http.ResponseWriter, r *http.Request) {
	addrHex := r.URL.Query().Get("address")
	if addrHex == "" {
		writeErr(w, "missing address param", 400)
		return
	}
	if _, err := hex.DecodeString(addrHex); err != nil || len(addrHex) != 40 {
		writeErr(w, "invalid address (need 40 hex chars)", 400)
		return
	}
	c := globalCache.GetContributor(addrHex)
	if c == nil {
		writeJSON(w, map[string]any{
			"address": addrHex, "totalMeasurements": 0, "tokenBalance": 0,
		})
		return
	}
	writeJSON(w, map[string]any{
		"address":           hexAddr(c.Address),
		"totalMeasurements": c.TotalMeasurements,
		"tokenBalance":      c.TokenBalance,
	})
}

// handleSubmitMeasurement accepts a POST with speed-test results, signs a
// submit_measurement transaction with the keystore key, and submits it on-chain.
func (s *meshpulseServer) handleSubmitMeasurement(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		writeErr(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Ping     uint64 `json:"ping"`
		Download uint64 `json:"download"`
		Upload   uint64 `json:"upload"`
		ISP      string `json:"isp"`
		Region   string `json:"region"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	txHash, err := SubmitMeasurementTx(req.Ping, req.Download, req.Upload, req.ISP, req.Region)
	if err != nil {
		log.Printf("SubmitMeasurementTx error: %v", err)
		writeErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]string{"txHash": txHash})
}
