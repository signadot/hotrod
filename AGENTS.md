# AGENTS.md

## Cursor Cloud specific instructions

### Product overview
HotROD (Rides on Demand) is a microservices demo: Frontend (:8080), Location (:8081), Driver (:8082 healthz), Route (:8083 gRPC). All four run from a single Go binary via `go run ./cmd/hotrod all`.

### Infrastructure dependencies (Docker)
Three containers must be running before starting the app:
- **Redis** `redis:7.2.3` on port 6379
- **MariaDB** `mariadb:11.6` on port 3306 (database `location`, root password `abc`)
- **Kafka** `confluentinc/confluent-local:7.5.0` on port 9092

Start them:
```bash
sudo docker run -d --name redis -p 6379:6379 redis:7.2.3
sudo docker run -d --name mariadb -p 3306:3306 -e MYSQL_DATABASE=location -e MYSQL_ROOT_PASSWORD=abc -e MYSQL_ROOT_HOST='%' mariadb:11.6
sudo docker run -d --name kafka -p 9092:9092 -p 29093:29093 \
  -e KAFKA_NODE_ID=0 -e KAFKA_CONTROLLER_LISTENER_NAMES=CONTROLLER \
  -e KAFKA_LISTENERS=PLAINTEXT://0.0.0.0:9092,CONTROLLER://0.0.0.0:29093 \
  -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092 \
  -e KAFKA_CONTROLLER_QUORUM_VOTERS="0@localhost:29093" \
  -e CLUSTER_ID=2dd8oDpjQRGlRWObUIROyQ \
  confluentinc/confluent-local:7.5.0
```

### Running the app locally
Required environment variables for local development:
```bash
export REDIS_ADDR=localhost:6379
export MYSQL_ADDR=localhost:3306
export MYSQL_PASS=abc
export KAFKA_BROKER_ADDR=localhost:9092
export LOCATION_ADDR=localhost:8081
export ROUTE_ADDR=localhost:8083
export SIGNADOT_BASELINE_KIND=Deployment
export SIGNADOT_BASELINE_NAMESPACE=hotrod
export SIGNADOT_BASELINE_NAME=driver
```
Then: `go run ./cmd/hotrod all`

### Key gotchas
- The driver service **panics** without `SIGNADOT_BASELINE_KIND`, `SIGNADOT_BASELINE_NAMESPACE`, and `SIGNADOT_BASELINE_NAME` env vars.
- `routeserver.signadot.svc` errors in logs are harmless — the Signadot route server is not available locally; the driver still processes messages.
- OTEL trace export errors (`localhost:4318`) are harmless — Jaeger is optional.
- The MariaDB database and seed data are auto-created on first query (no manual migration needed).
- React frontend must be built before running the Go server: `make build-frontend-app` (uses `yarn` inside `services/frontend/react_app/`).

### Signadot CLI
Installed at `/usr/local/bin/signadot` (v1.5.0). Requires a `SIGNADOT_API_KEY` for authentication (`signadot auth login --with-api-key <key>`). Not needed for local-only dev, but required for sandbox workflows and smart tests.

### Lint / Test / Build
- **Go vet:** `go vet ./...`
- **Go tests:** `go test ./...`
- **React lint:** `cd services/frontend/react_app && yarn lint` (has pre-existing warnings)
- **Build binary:** `make build` (builds frontend + Go binary)
- **Build frontend only:** `make build-frontend-app`
