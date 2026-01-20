#!/bin/bash

# Stop all HotRod local services by killing processes on ports 8080, 8082, 8083

KILLED_PIDS=""

for PORT in 8080 8082 8083; do
  PIDS=$(lsof -ti:$PORT 2>/dev/null)
  if [ -n "$PIDS" ]; then
    kill -9 $PIDS 2>/dev/null
    KILLED_PIDS="$KILLED_PIDS $PIDS"
  fi
done

if [ -n "$KILLED_PIDS" ]; then
  echo "Stopped processes:$KILLED_PIDS"
else
  echo "No services running"
fi
