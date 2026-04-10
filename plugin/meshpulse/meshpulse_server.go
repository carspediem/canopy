package contract

// meshpulse_server.go — lightweight HTTP server for the MeshPulse UI.
//
// Started from StartPlugin() to expose:
//   GET /             → single-file HTML UI
//   GET /api/stats    → NetworkStats JSON
//   GET /api/feed     → last 50 measurements JSON
//   GET /api/leaderboard → top 20 contributors JSON
//   GET /api/contributor?address=<hex> → single contributor JSON
//
// State reads are done via the plugin's existing StateRead() path so no
// second socket connection is needed.

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sort"
)

// StartUIServer starts the MeshPulse HTTP server on :8080.
// It receives a reference to the running Plugin so it can read chain state.
func StartUIServer(p *Plugin, cfg Config) {
	s := &meshpulseServer{plugin: p, cfg: cfg}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.serveUI)
	mux.HandleFunc("/api/stats", s.handleStats)
	mux.HandleFunc("/api/feed", s.handleFeed)
	mux.HandleFunc("/api/leaderboard", s.handleLeaderboard)
	mux.HandleFunc("/api/contributor", s.handleContributor)

	log.Printf("🌐 MeshPulse UI → http://localhost:8080")
	go func() {
		if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Printf("UI server error: %v", err)
		}
	}()
}

// meshpulseServer holds shared state for HTTP handlers.
type meshpulseServer struct {
	plugin *Plugin
	cfg    Config
}

// newCtx creates a throw-away contract context for read-only state access.
func (s *meshpulseServer) newCtx() *Contract {
	return &Contract{
		Config:    s.cfg,
		FSMConfig: s.plugin.fsmConfig,
		plugin:    s.plugin,
		fsmId:     rand.Uint64(),
	}
}

// readKey performs a single key lookup in the chain state.
func (s *meshpulseServer) readKey(key []byte) ([]byte, error) {
	qid := rand.Uint64()
	resp, pErr := s.plugin.StateRead(s.newCtx(), &PluginStateReadRequest{
		Keys: []*PluginKeyRead{{QueryId: qid, Key: key}},
	})
	if pErr != nil {
		return nil, fmt.Errorf("%s", pErr.Msg)
	}
	if resp.Error != nil {
		return nil, fmt.Errorf("%s", resp.Error.Msg)
	}
	for _, r := range resp.Results {
		if r.QueryId == qid && len(r.Entries) > 0 {
			return r.Entries[0].Value, nil
		}
	}
	return nil, nil
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

// ─── Handlers ─────────────────────────────────────────────────────────────────

func (s *meshpulseServer) serveUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, meshpulseHTML)
}

func (s *meshpulseServer) handleStats(w http.ResponseWriter, r *http.Request) {
	b, err := s.readKey(KeyForNetworkStats())
	if err != nil {
		writeErr(w, err.Error(), 500)
		return
	}
	stats := &NetworkStats{}
	if pErr := unmarshalState(b, stats); pErr != nil {
		writeErr(w, pErr.Msg, 500)
		return
	}
	writeJSON(w, stats)
}

func (s *meshpulseServer) handleFeed(w http.ResponseWriter, r *http.Request) {
	seqB, err := s.readKey(KeyForSeq())
	if err != nil {
		writeErr(w, err.Error(), 500)
		return
	}
	seq := decodeU64(seqB)
	if seq == 0 {
		writeJSON(w, []any{})
		return
	}
	const maxFeed = 50
	start := uint64(1)
	if seq > maxFeed {
		start = seq - maxFeed + 1
	}
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
	var feed []measurementJSON
	for id := seq; id >= start; id-- {
		b2, e := s.readKey(KeyForMeasurement(id))
		if e != nil || len(b2) == 0 {
			continue
		}
		m := &Measurement{}
		if pErr := unmarshalState(b2, m); pErr != nil {
			continue
		}
		feed = append(feed, measurementJSON{
			ID:        m.ID,
			Address:   hex.EncodeToString(m.Address),
			Ping:      m.Ping,
			Download:  m.Download,
			Upload:    m.Upload,
			ISP:       m.ISP,
			Region:    m.Region,
			Timestamp: m.Timestamp,
		})
	}
	writeJSON(w, feed)
}

func (s *meshpulseServer) handleLeaderboard(w http.ResponseWriter, r *http.Request) {
	seqB, err := s.readKey(KeyForSeq())
	if err != nil {
		writeErr(w, err.Error(), 500)
		return
	}
	seq := decodeU64(seqB)
	seen := map[string]*Contributor{}
	// Walk measurements to discover unique contributors, then read their records.
	for id := seq; id >= 1 && len(seen) < 100; id-- {
		b2, e := s.readKey(KeyForMeasurement(id))
		if e != nil || len(b2) == 0 {
			continue
		}
		m := &Measurement{}
		if pErr := unmarshalState(b2, m); pErr != nil {
			continue
		}
		key := hex.EncodeToString(m.Address)
		if _, ok := seen[key]; !ok {
			cb, ce := s.readKey(KeyForContributor(m.Address))
			if ce == nil && len(cb) > 0 {
				c := &Contributor{}
				if pErr := unmarshalState(cb, c); pErr == nil {
					seen[key] = c
				}
			}
		}
	}
	type entry struct {
		Address           string `json:"address"`
		TotalMeasurements uint64 `json:"totalMeasurements"`
		TokenBalance      uint64 `json:"tokenBalance"`
	}
	result := make([]entry, 0, len(seen))
	for addrHex, c := range seen {
		result = append(result, entry{
			Address:           addrHex,
			TotalMeasurements: c.TotalMeasurements,
			TokenBalance:      c.TokenBalance,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].TotalMeasurements > result[j].TotalMeasurements
	})
	if len(result) > 20 {
		result = result[:20]
	}
	writeJSON(w, result)
}

func (s *meshpulseServer) handleContributor(w http.ResponseWriter, r *http.Request) {
	addrHex := r.URL.Query().Get("address")
	if addrHex == "" {
		writeErr(w, "missing address param", 400)
		return
	}
	addr, err := hex.DecodeString(addrHex)
	if err != nil || len(addr) != 20 {
		writeErr(w, "invalid address (need 40 hex chars)", 400)
		return
	}
	b, e := s.readKey(KeyForContributor(addr))
	if e != nil {
		writeErr(w, e.Error(), 500)
		return
	}
	if len(b) == 0 {
		writeJSON(w, map[string]any{
			"address": addrHex, "totalMeasurements": 0, "tokenBalance": 0,
		})
		return
	}
	c := &Contributor{}
	if pErr := unmarshalState(b, c); pErr != nil {
		writeErr(w, pErr.Msg, 500)
		return
	}
	writeJSON(w, map[string]any{
		"address":           hex.EncodeToString(c.Address),
		"totalMeasurements": c.TotalMeasurements,
		"tokenBalance":      c.TokenBalance,
	})
}
