package contract

// meshpulse_rpc.go — on-chain transaction submission for MeshPulse.
//
// Builds, signs, and submits a submit_measurement transaction to the Canopy node
// RPC at http://localhost:50002/v1/tx using the keystore key approach.
//
// BLS12-381 signing is inlined to avoid a circular import with the crypto package.

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/drand/kyber"
	bls12381 "github.com/drand/kyber-bls12381"
	"github.com/drand/kyber/pairing"
	"github.com/drand/kyber/sign/bdn"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	// keystoreAddress is the MeshPulse node address used for signing.
	keystoreAddress = "aaec98c808751dff160dba5dbaedeaa1a1b23b9c"
	// keystorePrivKey is the BLS12-381 private key (32 bytes hex).
	keystorePrivKey = "02d0f5a7af1d85af18697a3e653346114337c353d345f18cf14d470358966361"

	// Canopy node RPC endpoints
	queryRPCURL = "http://localhost:50002"

	// networkID and chainID for Canopy
	meshNetworkID = uint64(1)
	meshChainID   = uint64(1)

	// defaultTxFee is the fee for submit_measurement transactions (uCNPY).
	defaultTxFee = uint64(10000)

	// submitMeasurementTypeURL is the protobuf type URL for MessageSubmitMeasurement.
	submitMeasurementTypeURL = "type.googleapis.com/types.MessageSubmitMeasurement"
	// submitMeasurementMsgType is the message type name registered in the plugin.
	submitMeasurementMsgType = "submit_measurement"
)

// SubmitMeasurementTx builds, signs, and submits a submit_measurement transaction.
// Returns the transaction hash on success.
func SubmitMeasurementTx(ping, download, upload uint64, isp, region string) (string, error) {
	// 1. Get the current block height for CreatedHeight
	height, err := getChainHeight()
	if err != nil {
		height = 0 // non-fatal: proceed with height 0
	}

	// 2. Decode the node address bytes
	addrBytes, err := hex.DecodeString(keystoreAddress)
	if err != nil {
		return "", fmt.Errorf("invalid keystore address: %w", err)
	}

	// 3. Build the MessageSubmitMeasurement proto and marshal it
	msg := &MessageSubmitMeasurement{
		Address:  addrBytes,
		Ping:     ping,
		Download: download,
		Upload:   upload,
		Isp:      isp,
		Region:   region,
	}
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return "", fmt.Errorf("marshal measurement msg: %w", err)
	}

	// 4. Wrap in Any
	msgAny := &anypb.Any{
		TypeUrl: submitMeasurementTypeURL,
		Value:   msgBytes,
	}

	// 5. Compute sign bytes using the Transaction proto (no Signature field)
	txTime := uint64(time.Now().UnixMicro())
	signBytes, err := getSignBytes(submitMeasurementMsgType, msgAny, txTime, height, defaultTxFee, "", meshNetworkID, meshChainID)
	if err != nil {
		return "", fmt.Errorf("get sign bytes: %w", err)
	}

	// 6. Load the BLS12-381 key and sign
	scalar, pubKeyBytes, err := parseBLSKey(keystorePrivKey)
	if err != nil {
		return "", fmt.Errorf("load BLS key: %w", err)
	}
	scheme := bdn.NewSchemeOnG2(newBLSSuite())
	signature, err := scheme.Sign(scalar, signBytes)
	if err != nil {
		return "", fmt.Errorf("BLS sign: %w", err)
	}

	// 7. Build the transaction JSON payload
	// For plugin-only message types use msgTypeUrl + msgBytes for exact byte control
	// (mirrors the tutorial's buildSignAndSendTx approach)
	txPayload := map[string]interface{}{
		"type":       submitMeasurementMsgType,
		"msgTypeUrl": submitMeasurementTypeURL,
		"msgBytes":   hex.EncodeToString(msgBytes),
		"signature": map[string]string{
			"publicKey": hex.EncodeToString(pubKeyBytes),
			"signature": hex.EncodeToString(signature),
		},
		"time":          txTime,
		"createdHeight": height,
		"fee":           defaultTxFee,
		"memo":          "",
		"networkID":     meshNetworkID,
		"chainID":       meshChainID,
		// Also include the human-readable msg field for nodes that prefer it
		"msg": map[string]interface{}{
			"address":  base64.StdEncoding.EncodeToString(addrBytes),
			"ping":     ping,
			"download": download,
			"upload":   upload,
			"isp":      isp,
			"region":   region,
		},
	}

	txJSON, err := json.Marshal(txPayload)
	if err != nil {
		return "", fmt.Errorf("marshal tx: %w", err)
	}

	// 8. POST to /v1/tx
	return submitTxJSON(txJSON)
}

// getChainHeight fetches the current block height from the Canopy node.
func getChainHeight() (uint64, error) {
	resp, err := http.Post(queryRPCURL+"/v1/query/height", "application/json", bytes.NewBufferString("{}"))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Height uint64 `json:"height"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}
	return result.Height, nil
}

// submitTxJSON POSTs the serialized transaction JSON to /v1/tx and returns the tx hash.
func submitTxJSON(txJSON []byte) (string, error) {
	resp, err := http.Post(queryRPCURL+"/v1/tx", "application/json", bytes.NewBuffer(txJSON))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}
	var txHash string
	if err := json.Unmarshal(body, &txHash); err != nil {
		return "", fmt.Errorf("parse tx hash: %v, body: %s", err, string(body))
	}
	return txHash, nil
}

// ──────────────────────────────────────────────────────────────────────────────
// Inlined BLS12-381 helpers (avoids circular import with the crypto package)
// ──────────────────────────────────────────────────────────────────────────────

func newBLSSuite() pairing.Suite { return bls12381.NewBLS12381Suite() }

// parseBLSKey decodes a 32-byte hex BLS12-381 private key and returns
// the kyber.Scalar and the corresponding public key bytes.
func parseBLSKey(hexKey string) (kyber.Scalar, []byte, error) {
	bz, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, nil, err
	}
	suite := newBLSSuite()
	scalar := suite.G2().Scalar()
	if err := scalar.UnmarshalBinary(bz); err != nil {
		return nil, nil, err
	}
	// Derive public key on G1
	pubPoint := suite.G1().Point().Mul(scalar, suite.G1().Point().Base())
	pubBytes, err := pubPoint.MarshalBinary()
	if err != nil {
		return nil, nil, err
	}
	return scalar, pubBytes, nil
}

// getSignBytes returns the canonical protobuf bytes for signing a transaction.
// Uses deterministic marshaling so it matches lib.Transaction.GetSignBytes() exactly.
func getSignBytes(msgType string, msg *anypb.Any, txTime, createdHeight, fee uint64, memo string, networkID, chainID uint64) ([]byte, error) {
	tx := &Transaction{
		MessageType:   msgType,
		Msg:           msg,
		Signature:     nil, // omitted for signing
		CreatedHeight: createdHeight,
		Time:          txTime,
		Fee:           fee,
		Memo:          memo,
		NetworkId:     networkID,
		ChainId:       chainID,
	}
	return proto.MarshalOptions{Deterministic: true}.Marshal(tx)
}
