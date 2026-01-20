#!/bin/bash

# Run HotRod services locally with Signadot sandbox environment
# Usage: ./run-local-services.sh [sandbox-name]

SANDBOX_NAME="${1:-distance-feature-test}"
KAFKA_BROKER="kafka-headless.hotrod-istio.svc:9092"
REDIS_ADDR="redis.hotrod-istio.svc:6379"
OTEL_ENDPOINT="http://jaeger.hotrod-istio.svc:4318"
LOCATION_ADDR="location.hotrod-istio.svc:8081"
ROUTE_ADDR="route.hotrod-istio.svc:8083"

# Check signadot connection
if ! signadot local status &>/dev/null; then
  echo "Error: signadot is not connected. Run 'signadot local connect' first."
  exit 1
fi

# Stop any existing services
./stop-local-services.sh

# Remove old log files
rm -f route.log driver.log frontend.log

# Start route service with Redis and OTEL
eval "$(signadot sandbox get-env $SANDBOX_NAME -l route-local 2>/dev/null | sed 's/#.*//')" && \
  export REDIS_ADDR="$REDIS_ADDR" && \
  export OTEL_EXPORTER_OTLP_ENDPOINT="$OTEL_ENDPOINT" && \
  go run cmd/hotrod/main.go route > route.log 2>&1 &
ROUTE_PID=$!

# Start driver service with Kafka, Redis, Route, and OTEL
eval "$(signadot sandbox get-env $SANDBOX_NAME -l driver-local 2>/dev/null | sed 's/#.*//')" && \
  export KAFKA_BROKER_ADDR="$KAFKA_BROKER" && \
  export REDIS_ADDR="$REDIS_ADDR" && \
  export ROUTE_ADDR="$ROUTE_ADDR" && \
  export OTEL_EXPORTER_OTLP_ENDPOINT="$OTEL_ENDPOINT" && \
  go run cmd/hotrod/main.go driver > driver.log 2>&1 &
DRIVER_PID=$!

# Start frontend service with Kafka, Redis, OTEL, and service addresses
eval "$(signadot sandbox get-env $SANDBOX_NAME -l frontend-local 2>/dev/null | sed 's/#.*//')" && \
  export KAFKA_BROKER_ADDR="$KAFKA_BROKER" && \
  export REDIS_ADDR="$REDIS_ADDR" && \
  export OTEL_EXPORTER_OTLP_ENDPOINT="$OTEL_ENDPOINT" && \
  export LOCATION_ADDR="$LOCATION_ADDR" && \
  export ROUTE_ADDR="$ROUTE_ADDR" && \
  go run cmd/hotrod/main.go frontend > frontend.log 2>&1 &
FRONTEND_PID=$!

echo "Route:    PID $ROUTE_PID    → route.log"
echo "Driver:   PID $DRIVER_PID   → driver.log"
echo "Frontend: PID $FRONTEND_PID → frontend.log"
