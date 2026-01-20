# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

HotROD (Rides on Demand) is a demo microservices application based on the Jaeger HotROD demo, modified to showcase Signadot sandboxes. It simulates a ride-hailing service with distributed tracing via OpenTelemetry.

## Architecture

The application consists of four Go microservices:

- **frontend** (port 8080) - HTTP server serving the React web UI and API endpoints (`/dispatch`, `/splash`, `/notifications`). Entry point for user requests.
- **location** (port 8081) - HTTP service managing location data via MySQL database. Provides CRUD operations for pickup/dropoff locations.
- **route** (port 8083) - gRPC service calculating routes and ETAs between locations. Uses Redis for caching.
- **driver** (port 8082) - Kafka consumer processing driver dispatch requests. Finds best available drivers based on ETA.

**Communication flow:**
1. Frontend receives dispatch request → resolves locations via location service
2. Frontend publishes to Kafka topic `driver`
3. Driver service consumes message → calls route service for ETAs → finds best driver
4. Notifications stored in Redis and polled by frontend

**Infrastructure dependencies:** Kafka, Redis, MySQL, Jaeger (for tracing)

## Build Commands

```bash
# Build the Go binary (includes React frontend build)
make build

# Build React frontend only
make build-frontend-app

# Build Docker image for local development
make dev-build-docker

# Generate protobuf files for route service
make generate-proto
```

## Running Services

**Run single service:**
```bash
go run ./cmd/hotrod <service>   # service: frontend, location, route, driver, all
```

**Run with hot-reload (frontend development):**
```bash
air
```

**Run all services locally with Signadot sandbox:**
```bash
./run-local-services.sh [sandbox-name]
```

## Testing

```bash
# Run Go tests
go test ./...

# Run Playwright e2e tests
npm run e2e:playwright

# Run Signadot Smart Tests
signadot st run --sandbox=<sandbox-name> --publish
```

Smart tests are Starlark files in `smart-tests/` directory.

## Environment Variables

Key service configuration via environment:
- `KAFKA_BROKER_ADDR` - Kafka broker address
- `REDIS_ADDR` - Redis address
- `MYSQL_ADDR` - MySQL address
- `LOCATION_ADDR` - Location service address
- `ROUTE_ADDR` - Route service address
- `OTEL_EXPORTER_OTLP_ENDPOINT` - OpenTelemetry collector endpoint

## Kubernetes Deployment

```bash
kubectl -n <namespace> apply -k k8s/overlays/prod/<devmesh|istio>
```

Kustomize overlays support: devmesh, istio, linkerd (with optional gateway API variants).

## Signadot Integration

Local development with Signadot sandboxes:
```bash
signadot local connect
signadot sandbox apply -f <sandbox-spec.yaml>
eval "$(signadot sandbox get-env <sandbox-name>)"
```

Sandbox templates are in `.signadot/` directory.
