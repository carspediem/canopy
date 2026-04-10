package contract

// meshpulse_state.go — on-chain state structs and key helpers for MeshPulse.
//
// State is stored as JSON-encoded bytes under well-known prefixed keys in the
// Canopy key-value store.  We avoid adding more .proto files for internal state
// to keep the plugin self-contained; the FSM never needs to decode these bytes
// directly — only the plugin contract does.

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
)

// ─── Token constants ───────────────────────────────────────────────────────────

const (
	// MeshpRewardPerSubmit is the $MESHP (in micro-denomination) awarded per measurement.
	MeshpRewardPerSubmit uint64 = 10_000_000 // 10 $MESHP
)

// ─── State key prefixes (must not collide with base plugin: 1, 2, 7) ──────────

var (
	measurementPrefix = []byte{10} // keyed by uint64 ID
	contributorPrefix = []byte{11} // keyed by address bytes
	networkStatsKey   = []byte{12} // single global entry
	seqKey            = []byte{13} // monotonic measurement counter
)

// ─── On-chain state types ──────────────────────────────────────────────────────

// Measurement records a single speed-test result submitted on-chain.
type Measurement struct {
	ID        uint64 `json:"id"`
	Address   []byte `json:"address"`
	Ping      uint64 `json:"ping"`     // milliseconds
	Download  uint64 `json:"download"` // Kbps
	Upload    uint64 `json:"upload"`   // Kbps
	ISP       string `json:"isp"`
	Region    string `json:"region"`
	Timestamp uint64 `json:"timestamp"` // Unix seconds
}

// Contributor aggregates a user's activity and $MESHP balance.
type Contributor struct {
	Address           []byte `json:"address"`
	TotalMeasurements uint64 `json:"totalMeasurements"`
	TokenBalance      uint64 `json:"tokenBalance"`
}

// NetworkStats holds aggregate network-wide speed-test statistics.
type NetworkStats struct {
	TotalMeasurements uint64 `json:"totalMeasurements"`
	AvgPing           uint64 `json:"avgPing"`
	AvgDownload       uint64 `json:"avgDownload"`
	AvgUpload         uint64 `json:"avgUpload"`
	// running sums used to recompute averages incrementally
	SumPing     uint64 `json:"sumPing"`
	SumDownload uint64 `json:"sumDownload"`
	SumUpload   uint64 `json:"sumUpload"`
}

// ─── State key constructors ────────────────────────────────────────────────────

// KeyForMeasurement returns the store key for a Measurement by its uint64 ID.
func KeyForMeasurement(id uint64) []byte {
	b := encodeU64(id)
	return JoinLenPrefix(measurementPrefix, b)
}

// KeyForContributor returns the store key for a Contributor by address.
func KeyForContributor(addr []byte) []byte {
	return JoinLenPrefix(contributorPrefix, addr)
}

// KeyForNetworkStats returns the store key for the singleton NetworkStats entry.
func KeyForNetworkStats() []byte {
	return networkStatsKey
}

// KeyForSeq returns the store key for the measurement sequence counter.
func KeyForSeq() []byte {
	return seqKey
}

// ─── JSON marshal helpers ──────────────────────────────────────────────────────

// UnmarshalState is exported for use by the HTTP server package.
func UnmarshalState(b []byte, v any) *PluginError { return unmarshalState(b, v) }

// DecodeU64 is exported for the HTTP server package.
func DecodeU64(b []byte) uint64 { return decodeU64(b) }

func marshalState(v any) ([]byte, *PluginError) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, NewError(50, "meshpulse", fmt.Sprintf("state marshal failed: %s", err))
	}
	return b, nil
}

func unmarshalState(b []byte, v any) *PluginError {
	if len(b) == 0 {
		return nil
	}
	if err := json.Unmarshal(b, v); err != nil {
		return NewError(51, "meshpulse", fmt.Sprintf("state unmarshal failed: %s", err))
	}
	return nil
}

// ─── Uint64 codec helpers ──────────────────────────────────────────────────────

func encodeU64(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func decodeU64(b []byte) uint64 {
	if len(b) < 8 {
		return 0
	}
	return binary.BigEndian.Uint64(b)
}
