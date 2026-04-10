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

- **Landing page**: hero section, feature pills, "Connect Wallet" button
- **Wallet modal**: enter your Canopy address to connect (persisted in localStorage)
- **Node Dashboard**: ACTIVE status with pulsing indicator, live uptime counter, $MESHP earned counter (auto-increments every 30s), real-time measurement feed, Claim Rewards button
- **My Stats**: wallet lookup, balance, measurement count, rank tier
- **Live Feed**: last 50 on-chain measurements
- **Leaderboard**: top 20 contributors by measurement count

### In-memory cache

HTTP handlers read from an in-memory `globalCache` instead of calling `StateRead` directly. This is required by the Canopy FSM protocol — `StateRead`/`StateWrite` are only permitted inside `CheckTx`/`DeliverTx` callbacks. The cache is populated by `DeliverTx` handlers on every confirmed measurement.

---

## How to Run

### Docker (recommended)

```bash
# From the canopy repo root:
docker build -f plugin/meshpulse/Dockerfile -t meshpulse .
docker run -p 8080:8080 meshpulse
```

Open **http://localhost:8080**

### Live demo

**http://204.168.151.179:8080**

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
├── Dockerfile                       # multi-stage build (canopy + plugin + alpine)
├── entrypoint.sh                    # starts canopy then plugin
├── init.expect                      # non-interactive key gen during docker build
├── contract/
│   ├── contract.go                  # ContractConfig + routing
│   ├── plugin.go                    # FSM socket protocol + StartUIServer call
│   ├── meshpulse.pb.go              # protobuf types for custom messages
│   ├── meshpulse_state.go           # on-chain state structs + key helpers
│   ├── meshpulse_handlers.go        # CheckTx / DeliverTx + cache population
│   ├── meshpulse_cache.go           # in-memory cache (globalCache)
│   ├── meshpulse_server.go          # HTTP API server (reads from cache)
│   ├── meshpulse_ui.go              # embedded single-file HTML/JS UI
│   ├── error.go                     # error definitions
│   ├── account.pb.go                # Account / Pool protobuf types
│   ├── tx.pb.go                     # Transaction / MessageSend types
│   └── plugin.pb.go                 # FSM ↔ Plugin protocol types
└── proto/
    ├── meshpulse.proto              # source proto for custom messages
    └── (other .proto files)
```
