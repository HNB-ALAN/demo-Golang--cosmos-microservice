# Service-04 – Basic Production Setup (Minimal)

Goal: run Service-04 in a minimal “production-like” mode using a standalone Cosmos daemon (uscd) and point the service to it, with health checks and basic acceptance.

This is the shortest path to a stable baseline before full hardening.

---

## 0) Prerequisites
- Go 1.24.4+
- Linux host with build tools (CGO enabled)
- RocksDB libs (for local RocksDB if needed)
- PostgreSQL, Redis, Kafka reachable (or disabled if not required for your smoke test)

---

## 1) Build and run the Cosmos daemon (single-node)
```bash
cd SERVICES/service-04-usc-blockchain-core/block-chain-cosmos/daemon/uscd
# Build
go build -o uscd ./cmd/main.go

# Initialize node home and genesis
./uscd init validator --chain-id usc-1

# Start node (single-node, default ports)
./uscd start
```
Expected:
- Node starts and logs height increases (> 0) within a few seconds.
- RPC typically at 26657 (if enabled in config).

Tips:
- Home directory: ~/.uscd
- Configs: ~/.uscd/config
- To reset: stop node, remove ~/.uscd and re-init.

---

## 2) Configure Service-04 to use daemon mode
Edit `configs/config.yaml` (or provide envs) to point to the running node.

Required keys (example):
- cosmos:
  - mode: daemon
  - rpcAddr: http://127.0.0.1:26657
  - chainId: usc-1
  - timeouts: { dialMs: 2000, requestMs: 5000 }
  - retries: { max: 3, backoffMs: 300 }

If using env vars, set:
```bash
export COSMOS_MODE=daemon
export COSMOS_RPC_ADDR=http://127.0.0.1:26657
export COSMOS_CHAIN_ID=usc-1
```

---

## 3) Run Service-04 against the node
```bash
cd SERVICES/service-04-usc-blockchain-core
# Option A: use existing binary
./service-04-service
# Option B: run from source
go run cmd/main.go
```
Expected:
- Service binds on gRPC :8004
- Startup logs indicate Cosmos client initialized (daemon mode)

Health check (port):
```bash
ss -ltnp | grep 8004 || true
```

---

## 4) Basic smoke tests
- gRPC health: use grpcurl or grpc_health_probe to verify service is healthy.
- Transaction (optional if CLI wired): submit a sample tx via `uscd tx` and query balance.
- Logs: ensure no fatal errors on service during node restarts.

Example (if CLI wired for usc_coin):
```bash
# Example only; adjust to your actual module/CLI once wired
uscd tx usc-coin transfer <from> <to> 10usc --chain-id usc-1 --yes
uscd query bank balances <to>
```

---

## 5) Minimal acceptance criteria
- Node reaches height > 0 and stays stable for > 5 minutes.
- Service-04 stays healthy (gRPC health OK) and reports node connectivity OK.
- If node restarts, service reconnects (retries/backoff) without crash.
- p95 service latency < 100ms for simple read endpoints (in staging conditions).

---

## 6) Next steps (beyond minimal)
- Real Msg/keeper state changes (usc_coin) + integration tests.
- RBAC interceptors, secrets via env/manager, tracing/metrics dashboards.
- K8s manifests with readiness/liveness, resource limits, and alerts.

For a full plan, see `docs/PRODUCTION_IMPLEMENTATION_PLAN.md`.

