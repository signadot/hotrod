# Signadot + Claude Code Demo

## The Bug

The `demo-buggy` branch has a distance feature with an intentional bug: driver service uses meters as kilometers without conversion, causing frontend to display incorrect distances.

## Prompt to setup the Sandbox (one time)
```
Create a Signadot sandbox named "demo-distance-feature" with local mappings for the route, driver, and frontend services (use the same ports locally) on the demo cluster. Use local mapping names: route-local, driver-local, and frontend-local. Use a dummy, non-conflicting port for the driver service as it does not expose a port. Save this sandbox spec in a file.
```

## Demo Prompt

```
The distance display feature has been implemented across the route, driver, and frontend services, and an E2E Playwright test has been added.

Start the local services using the sandbox environment, run the E2E Playwright test to verify the feature works, and if it fails, investigate the logs from all three services (route.log, driver.log, frontend.log), identify the root cause, fix the bug, build and rerun the services changed and verify all tests pass. Continuously execute this loop until the tests pass.

Important:
- After making any code changes to the frontend React app (services/frontend/react_app/), rebuild it with `make build-frontend-app` before restarting services. Go code changes are picked up automatically by `go run`.
- Run Playwright tests with a single worker: `npm run e2e:playwright -- --workers=1`
```

## Pre-Demo Setup

```bash
./demo-setup.sh
```

This script:
- Resets to clean `demo-buggy` branch state
- Builds the frontend React app
- Cleans up any existing sandbox
- Stops running services and clears logs

## Build Requirements

- **Go code**: Automatically compiled by `go run` - no manual build needed
- **Frontend React app**: Must be rebuilt with `make build-frontend-app` after any changes to `services/frontend/react_app/`

## What Claude Does

1. **Creates sandbox** via MCP server with local mappings for route/driver/frontend
2. **Starts local services** using `./run-local-services.sh demo-distance-feature`
3. **Runs E2E test** - fails (distance not displaying correctly)
4. **Investigates logs**:
   - `route.log`: returns distance in meters (e.g., 5234 meters)
   - `driver.log`: uses raw meters value as km (bug - shows "5234 km away")
   - `frontend.log`: receives malformed distance data
5. **Fixes** `services/driver/consumer.go` - adds meters-to-km conversion (`distance/1000`)
6. **Restarts services and re-runs test** - passes

## Reset for Next Demo

Just run `./demo-setup.sh` again - it reverts code changes, rebuilds the frontend, cleans up the sandbox, and prepares for the next run.

## Architecture

```
Developer Machine                    Kubernetes Cluster
┌────────────────────┐              ┌─────────────────┐
│ Frontend :8080     │              │ Location svc    │
│ Driver   :8082     │◄──routing───►│ Redis           │
│ Route    :8083     │   key        │ Kafka           │
└────────────────────┘              └─────────────────┘
     (local, editable)                (remote, shared)
```
