# AGENTS.md

## Cursor Cloud specific instructions

### Overview

HotROD (Rides on Demand) is a microservices ride-hailing demo written in Go with a React (Vite/TypeScript) frontend. A single Go binary (`cmd/hotrod`) runs 4 services: **frontend** (:8080), **location** (:8081), **driver** (:8082 healthz), **route** (:8083 gRPC). Infrastructure deps: Redis, MySQL (MariaDB), Kafka.

### Infrastructure (Docker)

Three Docker containers provide infrastructure. Docker is pre-installed; the nested-VM workarounds (fuse-overlayfs, iptables-legacy) are already configured. Start/restart containers:

```bash
sudo dockerd &>/tmp/dockerd.log &
sleep 3
sudo docker start redis mariadb kafka
```

If containers don't exist yet (first-time setup):

```bash
sudo docker run -d --name redis -p 6379:6379 redis:7.2.3
sudo docker run -d --name mariadb -p 3306:3306 -e MYSQL_ROOT_PASSWORD=abc -e MYSQL_DATABASE=location -e MYSQL_ROOT_HOST=% mariadb:11.6
sudo docker run -d --name kafka -p 9092:9092 -p 29093:29093 confluentinc/confluent-local:7.5.0
```

### Running all services locally

Default env vars point to k8s DNS names. Override them for localhost. See `k8s/base/*.yaml` for the full env var list per service. Key overrides:

| Variable | Local value |
|---|---|
| `REDIS_ADDR` | `localhost:6379` |
| `MYSQL_ADDR` | `localhost:3306` |
| `MYSQL_PASS` | `abc` |
| `KAFKA_BROKER_ADDR` | `localhost:9092` |
| `LOCATION_ADDR` | `localhost:8081` |
| `ROUTE_ADDR` | `localhost:8083` |

Additionally, the driver service **requires** three baseline env vars or it panics on startup. Use these dummy values for local dev:

| Variable | Local value |
|---|---|
| `*_BASELINE_KIND` | `Deployment` |
| `*_BASELINE_NAMESPACE` | `hotrod` |
| `*_BASELINE_NAME` | `driver` |

(The `*` prefix matches the vendor-specific env prefix used in `pkg/config/` and `k8s/base/driver.yaml`.)

Run all services: `go run ./cmd/hotrod all` with the above env vars exported.

**Non-obvious gotchas:**
- The three baseline env vars above are **required** or the driver service panics. See `services/driver/consumer.go` line 40.
- The vendor routeserver is unavailable locally; the driver's routing watcher retries in the background — this is harmless.
- OTEL trace export errors to `:4318` are harmless (no local Jaeger). Set `OTEL_EXPORTER_TYPE=stdout` to log traces to stdout instead, or ignore the errors.
- The React frontend must be built before running the Go binary: `make build-frontend-app` (uses `yarn && yarn build`). The built assets are Go-embedded via `//go:embed`.

### Lint / Test / Build

- **Go tests:** `go test ./...`
- **Go build:** `make build` (builds frontend first, then Go binary to `dist/`)
- **Frontend lint:** `cd services/frontend/react_app && yarn lint` (has pre-existing warnings/errors)
- **Frontend dev server:** `cd services/frontend/react_app && yarn dev` (Vite dev server on :5173, for React-only development)
- **Hot-reload:** `air` (requires [air](https://github.com/air-verse/air) installed; uses `.air.toml` config)
