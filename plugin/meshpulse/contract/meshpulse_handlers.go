package contract

// meshpulse_handlers.go — CheckTx / DeliverTx implementations for MeshPulse
// transactions: SubmitMeasurement and ClaimReward.

import (
	"math/rand"
	"time"
)

// ─── SubmitMeasurement ─────────────────────────────────────────────────────────

// CheckMessageSubmitMeasurement performs stateless validation of a SubmitMeasurement tx.
func (c *Contract) CheckMessageSubmitMeasurement(msg *MessageSubmitMeasurement) *PluginCheckResponse {
	if len(msg.Address) != 20 {
		return &PluginCheckResponse{Error: ErrInvalidAddress()}
	}
	if msg.Ping == 0 && msg.Download == 0 && msg.Upload == 0 {
		return &PluginCheckResponse{Error: ErrInvalidMeasurement()}
	}
	return &PluginCheckResponse{
		Recipient:         msg.Address,
		AuthorizedSigners: [][]byte{msg.Address},
	}
}

// DeliverMessageSubmitMeasurement applies a SubmitMeasurement transaction:
//  1. Increments the global sequence counter to get a new measurement ID.
//  2. Stores the Measurement record.
//  3. Upserts the Contributor record (increments count, adds reward tokens).
//  4. Updates the global NetworkStats aggregates.
func (c *Contract) DeliverMessageSubmitMeasurement(msg *MessageSubmitMeasurement, fee uint64) *PluginDeliverResponse {
	var (
		seqQID, contribQID, statsQID = rand.Uint64(), rand.Uint64(), rand.Uint64()
		feeQID                       = rand.Uint64()
	)

	// ── Read current seq, contributor, network stats, fee pool ──────────────
	readResp, err := c.plugin.StateRead(c, &PluginStateReadRequest{
		Keys: []*PluginKeyRead{
			{QueryId: seqQID, Key: KeyForSeq()},
			{QueryId: contribQID, Key: KeyForContributor(msg.Address)},
			{QueryId: statsQID, Key: KeyForNetworkStats()},
			{QueryId: feeQID, Key: KeyForFeePool(c.Config.ChainId)},
		},
	})
	if err != nil {
		return &PluginDeliverResponse{Error: err}
	}
	if readResp.Error != nil {
		return &PluginDeliverResponse{Error: readResp.Error}
	}

	var seqBytes, contribBytes, statsBytes, feeBytes []byte
	for _, r := range readResp.Results {
		if len(r.Entries) == 0 {
			continue
		}
		switch r.QueryId {
		case seqQID:
			seqBytes = r.Entries[0].Value
		case contribQID:
			contribBytes = r.Entries[0].Value
		case statsQID:
			statsBytes = r.Entries[0].Value
		case feeQID:
			feeBytes = r.Entries[0].Value
		}
	}

	// ── Decode existing state ────────────────────────────────────────────────
	seq := decodeU64(seqBytes)
	newID := seq + 1

	contrib := &Contributor{Address: msg.Address}
	if pluginErr := unmarshalState(contribBytes, contrib); pluginErr != nil {
		return &PluginDeliverResponse{Error: pluginErr}
	}

	stats := &NetworkStats{}
	if pluginErr := unmarshalState(statsBytes, stats); pluginErr != nil {
		return &PluginDeliverResponse{Error: pluginErr}
	}

	feePool := &Pool{}
	if pluginErr := Unmarshal(feeBytes, feePool); pluginErr != nil {
		return &PluginDeliverResponse{Error: pluginErr}
	}

	// ── Build new Measurement ────────────────────────────────────────────────
	m := &Measurement{
		ID:        newID,
		Address:   msg.Address,
		Ping:      msg.Ping,
		Download:  msg.Download,
		Upload:    msg.Upload,
		ISP:       msg.Isp,
		Region:    msg.Region,
		Timestamp: uint64(time.Now().Unix()),
	}

	// ── Update Contributor ───────────────────────────────────────────────────
	contrib.TotalMeasurements++
	contrib.TokenBalance += MeshpRewardPerSubmit

	// ── Update NetworkStats (incremental average) ────────────────────────────
	stats.TotalMeasurements++
	stats.SumPing += msg.Ping
	stats.SumDownload += msg.Download
	stats.SumUpload += msg.Upload
	n := stats.TotalMeasurements
	stats.AvgPing = stats.SumPing / n
	stats.AvgDownload = stats.SumDownload / n
	stats.AvgUpload = stats.SumUpload / n

	// ── Fee pool: absorb the tx fee ──────────────────────────────────────────
	feePool.Amount += fee

	// ── Marshal updated state ────────────────────────────────────────────────
	mBytes, pluginErr := marshalState(m)
	if pluginErr != nil {
		return &PluginDeliverResponse{Error: pluginErr}
	}
	contribBytes, pluginErr = marshalState(contrib)
	if pluginErr != nil {
		return &PluginDeliverResponse{Error: pluginErr}
	}
	statsBytes, pluginErr = marshalState(stats)
	if pluginErr != nil {
		return &PluginDeliverResponse{Error: pluginErr}
	}
	feeBytes, pluginErr = Marshal(feePool)
	if pluginErr != nil {
		return &PluginDeliverResponse{Error: pluginErr}
	}

	// ── Write all state changes atomically ───────────────────────────────────
	writeResp, err := c.plugin.StateWrite(c, &PluginStateWriteRequest{
		Sets: []*PluginSetOp{
			{Key: KeyForSeq(), Value: encodeU64(newID)},
			{Key: KeyForMeasurement(newID), Value: mBytes},
			{Key: KeyForContributor(msg.Address), Value: contribBytes},
			{Key: KeyForNetworkStats(), Value: statsBytes},
			{Key: KeyForFeePool(c.Config.ChainId), Value: feeBytes},
		},
	})
	if err != nil {
		return &PluginDeliverResponse{Error: err}
	}
	// ── Update in-memory cache for UI server ─────────────────────────────────
	globalCache.AddMeasurement(m, *stats)
	globalCache.UpdateContributor(contrib)
	return &PluginDeliverResponse{Error: writeResp.Error}
}

// ─── ClaimReward ──────────────────────────────────────────────────────────────

// CheckMessageClaimReward performs stateless validation of a ClaimReward tx.
func (c *Contract) CheckMessageClaimReward(msg *MessageClaimReward) *PluginCheckResponse {
	if len(msg.FromAddress) != 20 {
		return &PluginCheckResponse{Error: ErrInvalidAddress()}
	}
	if len(msg.ToAddress) != 20 {
		return &PluginCheckResponse{Error: ErrInvalidAddress()}
	}
	if msg.Amount == 0 {
		return &PluginCheckResponse{Error: ErrInvalidAmount()}
	}
	return &PluginCheckResponse{
		Recipient:         msg.ToAddress,
		AuthorizedSigners: [][]byte{msg.FromAddress},
	}
}

// DeliverMessageClaimReward transfers earned $MESHP tokens from the contributor's
// in-plugin balance to their main Canopy account (or another address).
func (c *Contract) DeliverMessageClaimReward(msg *MessageClaimReward, fee uint64) *PluginDeliverResponse {
	var (
		contribQID, toQID, feeQID = rand.Uint64(), rand.Uint64(), rand.Uint64()
	)

	// ── Read contributor, destination account, fee pool ──────────────────────
	readResp, err := c.plugin.StateRead(c, &PluginStateReadRequest{
		Keys: []*PluginKeyRead{
			{QueryId: contribQID, Key: KeyForContributor(msg.FromAddress)},
			{QueryId: toQID, Key: KeyForAccount(msg.ToAddress)},
			{QueryId: feeQID, Key: KeyForFeePool(c.Config.ChainId)},
		},
	})
	if err != nil {
		return &PluginDeliverResponse{Error: err}
	}
	if readResp.Error != nil {
		return &PluginDeliverResponse{Error: readResp.Error}
	}

	var contribBytes, toBytes, feeBytes []byte
	for _, r := range readResp.Results {
		if len(r.Entries) == 0 {
			continue
		}
		switch r.QueryId {
		case contribQID:
			contribBytes = r.Entries[0].Value
		case toQID:
			toBytes = r.Entries[0].Value
		case feeQID:
			feeBytes = r.Entries[0].Value
		}
	}

	// ── Decode ───────────────────────────────────────────────────────────────
	contrib := &Contributor{Address: msg.FromAddress}
	if pluginErr := unmarshalState(contribBytes, contrib); pluginErr != nil {
		return &PluginDeliverResponse{Error: pluginErr}
	}

	toAcct := &Account{}
	if pluginErr := Unmarshal(toBytes, toAcct); pluginErr != nil {
		return &PluginDeliverResponse{Error: pluginErr}
	}

	feePool := &Pool{}
	if pluginErr := Unmarshal(feeBytes, feePool); pluginErr != nil {
		return &PluginDeliverResponse{Error: pluginErr}
	}

	// ── Validate contributor has enough earned tokens ─────────────────────────
	totalNeeded := msg.Amount + fee
	if contrib.TokenBalance < totalNeeded {
		return &PluginDeliverResponse{Error: ErrInsufficientFunds()}
	}

	// ── Apply state changes ──────────────────────────────────────────────────
	contrib.TokenBalance -= totalNeeded
	toAcct.Amount += msg.Amount
	feePool.Amount += fee

	// ── Marshal ──────────────────────────────────────────────────────────────
	contribBytes, pluginErr := marshalState(contrib)
	if pluginErr != nil {
		return &PluginDeliverResponse{Error: pluginErr}
	}
	toBytes, pluginErr = Marshal(toAcct)
	if pluginErr != nil {
		return &PluginDeliverResponse{Error: pluginErr}
	}
	feeBytes, pluginErr = Marshal(feePool)
	if pluginErr != nil {
		return &PluginDeliverResponse{Error: pluginErr}
	}

	// ── Write ─────────────────────────────────────────────────────────────────
	sets := []*PluginSetOp{
		{Key: KeyForContributor(msg.FromAddress), Value: contribBytes},
		{Key: KeyForAccount(msg.ToAddress), Value: toBytes},
		{Key: KeyForFeePool(c.Config.ChainId), Value: feeBytes},
	}
	// Delete contributor record if fully drained
	var writeResp *PluginStateWriteResponse
	if contrib.TokenBalance == 0 {
		writeResp, err = c.plugin.StateWrite(c, &PluginStateWriteRequest{
			Sets:    sets[1:], // keep toAcct + feePool
			Deletes: []*PluginDeleteOp{{Key: KeyForContributor(msg.FromAddress)}},
		})
	} else {
		writeResp, err = c.plugin.StateWrite(c, &PluginStateWriteRequest{Sets: sets})
	}
	if err != nil {
		return &PluginDeliverResponse{Error: err}
	}
	return &PluginDeliverResponse{Error: writeResp.Error}
}
