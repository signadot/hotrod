#!/bin/bash

# Run HotRod services locally with Signadot sandbox environment
# Usage: ./run-local-services.sh <sandbox-name> [service...]
# Examples:
#   ./run-local-services.sh my-sandbox                        # Start all services
#   ./run-local-services.sh my-sandbox frontend               # Start only frontend
#   ./run-local-services.sh my-sandbox frontend driver        # Start frontend and driver
#   ./run-local-services.sh my-sandbox location frontend      # Start location and frontend
#   ./run-local-services.sh my-sandbox route driver frontend location # Start all four

set -e

if [ -z "$1" ]; then
  echo "Usage: $0 <sandbox-name> [service...]"
  echo ""
  echo "Services: frontend, driver, route, location (default: all)"
  echo ""
  echo "Examples:"
  echo "  $0 my-sandbox                  # Start all services"
  echo "  $0 my-sandbox frontend         # Start only frontend"
  echo "  $0 my-sandbox frontend driver  # Start frontend and driver"
  echo "  $0 my-sandbox location         # Start only location"
  exit 1
fi

SANDBOX_NAME="$1"
shift

# If no services specified, start all
if [ $# -eq 0 ]; then
  SERVICES="route driver location frontend"
else
  SERVICES="$*"
fi

KAFKA_BROKER="kafka-headless.hotrod-istio.svc:9092"
REDIS_ADDR="redis.hotrod-istio.svc:6379"
OTEL_ENDPOINT="http://jaeger.hotrod-istio.svc:4318"
LOCATION_ADDR="location.hotrod-istio.svc:8081"
ROUTE_ADDR="route.hotrod-istio.svc:8083"
MYSQL_ADDR="mysql.hotrod-istio.svc:3306"

# Check signadot connection
if ! signadot local status &>/dev/null; then
  echo "Error: signadot is not connected. Run 'signadot local connect' first."
  exit 1
fi

# Get port for a service
get_port() {
  case $1 in
    route) echo 8083 ;;
    driver) echo 8082 ;;
    location) echo 8081 ;;
    frontend) echo 8080 ;;
  esac
}

# Try to resolve the local mapping name for the sandbox
resolve_local_name() {
  local sandbox=$1
  local svc=$2
  # Try common naming conventions
  for name in "${svc}-local" "local-${svc}" "${svc}"; do
    if signadot sandbox get-env "$sandbox" -l "$name" &>/dev/null; then
      echo "$name"
      return
    fi
  done
  echo ""
}

# Stop services that are being started (by port)
for SVC in $SERVICES; do
  PORT=$(get_port "$SVC")
  if [ -n "$PORT" ]; then
    PIDS=$(lsof -ti:$PORT 2>/dev/null || true)
    if [ -n "$PIDS" ]; then
      kill -9 $PIDS 2>/dev/null || true
      echo "Stopped existing process on port $PORT"
    fi
  fi
done

rm -f route.log driver.log frontend.log location.log

for SVC in $SERVICES; do
  LOCAL_NAME=$(resolve_local_name "$SANDBOX_NAME" "$SVC")
  if [ -z "$LOCAL_NAME" ]; then
    echo "Warning: Could not find local mapping for $SVC in sandbox $SANDBOX_NAME, starting without sandbox env"
    ENV_CMD="true"
  else
    ENV_CMD="eval \"\$(signadot sandbox get-env $SANDBOX_NAME -l $LOCAL_NAME 2>/dev/null | sed 's/#.*//')\""
  fi

  case $SVC in
    route)
      eval "$ENV_CMD" && \
        export REDIS_ADDR="$REDIS_ADDR" && \
        export OTEL_EXPORTER_OTLP_ENDPOINT="$OTEL_ENDPOINT" && \
        go run cmd/hotrod/main.go route > route.log 2>&1 &
      echo "Route:    PID $!  → route.log     (port 8083)"
      ;;
    driver)
      eval "$ENV_CMD" && \
        export KAFKA_BROKER_ADDR="$KAFKA_BROKER" && \
        export REDIS_ADDR="$REDIS_ADDR" && \
        export ROUTE_ADDR="$ROUTE_ADDR" && \
        export OTEL_EXPORTER_OTLP_ENDPOINT="$OTEL_ENDPOINT" && \
        go run cmd/hotrod/main.go driver > driver.log 2>&1 &
      echo "Driver:   PID $!  → driver.log    (port 8082)"
      ;;
    location)
      eval "$ENV_CMD" && \
        export MYSQL_ADDR="$MYSQL_ADDR" && \
        export OTEL_EXPORTER_OTLP_ENDPOINT="$OTEL_ENDPOINT" && \
        go run cmd/hotrod/main.go location > location.log 2>&1 &
      echo "Location: PID $!  → location.log  (port 8081)"
      ;;
    frontend)
      eval "$ENV_CMD" && \
        export KAFKA_BROKER_ADDR="$KAFKA_BROKER" && \
        export REDIS_ADDR="$REDIS_ADDR" && \
        export OTEL_EXPORTER_OTLP_ENDPOINT="$OTEL_ENDPOINT" && \
        export LOCATION_ADDR="$LOCATION_ADDR" && \
        export ROUTE_ADDR="$ROUTE_ADDR" && \
        go run cmd/hotrod/main.go frontend > frontend.log 2>&1 &
      echo "Frontend: PID $!  → frontend.log  (port 8080)"
      ;;
    *)
      echo "Unknown service: $SVC (valid: frontend, driver, route, location)"
      ;;
  esac
done

echo ""
echo "Services started. Logs: tail -f {route,driver,frontend,location}.log"
