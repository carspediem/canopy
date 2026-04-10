# MeshPulse ⚡ — DePIN Speed Network

> A Canopy plugin that turns every internet connection into a DePIN node.
> Users run browser-based speed tests, submit results on-chain, and earn **$MESHP** tokens.

Built for the **Canopy Vibe Coding Contest**.

---

## Architecture

```
main.go
  └─ contract.StartPlugin(cfg)
       ├─ Connects to Canopy FSM via plugin.sock
       ├─ Registers MeshPulse tx types with the FSM
       ├─ Handles CheckTx / DeliverTx for all transactions
       └─ Starts HTTP UI server on :8080
```

### On-Chain State

| Type | Key | Description |
|------|-----|-------------|
| `Measurement` | prefix `0x0a` + uint64 ID | Single speed-test result |
| `Contributor` | prefix `0x0b` + address | Per-user stats + $MESHP balance |
| `NetworkStats` | `0x0c` | Global aggregate averages |
| Sequence counter | `0x0d` | Monotonic measurement ID |

### Transactions

| Name | Type URL | What it does |
|------|----------|--------------|
| `submit_measurement` | `types.MessageSubmitMeasurement` | Records ping/download/upload + awards **10 $MESHP** |
| `claim_reward` | `types.MessageClaimReward` | Transfers earned $MESHP to a wallet |
| `send` | `types.MessageSend` | Standard CNPY transfer (base template) |

### Frontend (single HTML file, no dependencies)

- **Page 1 – Speed Test**: animated gauge + run test + submit to chain
- **Page 2 – My Stats**: wallet lookup, balance, measurement count, rank tier
- **Page 3 – Live Feed**: last 50 on-chain measurements, real-time
- **Page 4 – Leaderboard**: top 20 contributors by measurement count

---

## How to Run

### Prerequisites

- Go 1.21+
- A running Canopy node (provides `plugin.sock`)

### 1. Build and start Canopy

```bash
# From the canopy repo root:
make build/canopy
canopy start
```

### 2. Run MeshPulse

```bash
cd plugin/meshpulse   # this directory
go run .
# or
make build && ./meshpulse
```

The plugin will:
1. Wait for `plugin.sock` to appear (Canopy FSM must be running)
2. Perform the handshake, registering `submit_measurement` and `claim_reward`
3. Start the UI at **http://localhost:8080**

### 3. Submit a measurement via CLI

```bash
canopy tx submit_measurement \
  --address <your-address> \
  --ping 12 \
  --download 50000 \
  --upload 20000 \
  --isp "Comcast" \
  --region "US-East"
```

Or use the browser UI — it builds the transaction payload and copies it to your clipboard.

---

## Token Economics

- Every `SubmitMeasurement` tx awards **10 $MESHP** (= 10,000,000 in micro-denomination)
- `ClaimReward` transfers earned $MESHP to any Canopy wallet address
- Fees are collected into the chain's fee pool for validator distribution

## File Layout

```
meshpulse/
├── main.go                          # entry point
├── go.mod
├── contract/
│   ├── contract.go                  # ContractConfig + routing
│   ├── plugin.go                    # FSM socket protocol (base template)
│   ├── meshpulse.pb.go              # protobuf types for custom messages
│   ├── meshpulse_state.go           # on-chain state structs + key helpers
│   ├── meshpulse_handlers.go        # CheckTx / DeliverTx implementations
│   ├── meshpulse_server.go          # HTTP API server
│   ├── meshpulse_ui.go              # embedded single-file HTML UI
│   ├── error.go                     # error definitions
│   ├── account.pb.go                # Account / Pool protobuf types
│   ├── tx.pb.go                     # Transaction / MessageSend types
│   └── plugin.pb.go                 # FSM ↔ Plugin protocol types
└── proto/
    ├── meshpulse.proto              # source proto for custom messages
    └── (other .proto files)
```
