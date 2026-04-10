package contract

// meshpulse_cache.go — in-memory read cache for the MeshPulse UI server.
//
// The Canopy plugin protocol only allows StateRead/StateWrite calls during
// FSM hook callbacks (CheckTx, DeliverTx, etc.).  The HTTP UI server runs
// outside those callbacks, so it reads from this cache instead.
//
// The cache is updated by DeliverMessageSubmitMeasurement and
// DeliverMessageClaimReward after every successful transaction.

import "sync"

// globalCache is the singleton cache instance shared by handlers and the UI server.
var globalCache = &MeshCache{
	contributors: map[string]*Contributor{},
}

// MeshCache holds an in-memory snapshot of the on-chain MeshPulse state.
type MeshCache struct {
	mu           sync.RWMutex
	stats        NetworkStats
	measurements []*Measurement // newest-first, capped at 200
	contributors map[string]*Contributor
}

// AddMeasurement records a new measurement and refreshes stats.
func (c *MeshCache) AddMeasurement(m *Measurement, stats NetworkStats) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// prepend
	c.measurements = append([]*Measurement{m}, c.measurements...)
	if len(c.measurements) > 200 {
		c.measurements = c.measurements[:200]
	}
	c.stats = stats
}

// UpdateContributor upserts a contributor record.
func (c *MeshCache) UpdateContributor(contrib *Contributor) {
	c.mu.Lock()
	defer c.mu.Unlock()
	key := string(contrib.Address)
	c.contributors[key] = contrib
}

// Stats returns a copy of the current network stats.
func (c *MeshCache) Stats() NetworkStats {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.stats
}

// Feed returns up to n recent measurements.
func (c *MeshCache) Feed(n int) []*Measurement {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if n > len(c.measurements) {
		n = len(c.measurements)
	}
	out := make([]*Measurement, n)
	copy(out, c.measurements[:n])
	return out
}

// Leaderboard returns contributors sorted by measurement count (top n).
func (c *MeshCache) Leaderboard(n int) []*Contributor {
	c.mu.RLock()
	defer c.mu.RUnlock()
	all := make([]*Contributor, 0, len(c.contributors))
	for _, v := range c.contributors {
		all = append(all, v)
	}
	// simple insertion sort (small n)
	for i := 1; i < len(all); i++ {
		for j := i; j > 0 && all[j].TotalMeasurements > all[j-1].TotalMeasurements; j-- {
			all[j], all[j-1] = all[j-1], all[j]
		}
	}
	if n > len(all) {
		n = len(all)
	}
	return all[:n]
}

// GetContributor returns a contributor by hex address string, or nil.
func (c *MeshCache) GetContributor(addrHex string) *Contributor {
	c.mu.RLock()
	defer c.mu.RUnlock()
	// contributors map is keyed by raw bytes string; search by hex is done below
	for _, v := range c.contributors {
		if hexAddr(v.Address) == addrHex {
			return v
		}
	}
	return nil
}
