#!/bin/bash
# Pre-demo setup script - idempotent, safe to run multiple times

set -e

echo "=== Demo Setup ==="

# Ensure on demo-buggy branch
CURRENT_BRANCH=$(git branch --show-current)
if [ "$CURRENT_BRANCH" != "demo-buggy" ]; then
  echo "Switching to demo-buggy branch..."
  git checkout demo-buggy
else
  echo "Already on demo-buggy branch"
fi

# Discard any local changes to ensure clean state
echo "Resetting any local changes..."
git checkout -- .

# Build frontend React app (required after git checkout)
# Note: Go code is compiled automatically by 'go run'
echo "Building frontend React app..."
make build-frontend-app

# Clean up any existing sandbox
echo "Cleaning up existing sandbox..."
signadot sandbox delete demo-distance-feature 2>/dev/null || true

# Stop any running services
echo "Stopping any running services..."
./stop-local-services.sh 2>/dev/null || true

# Clear logs
echo "Clearing log files..."
rm -f route.log driver.log frontend.log

echo ""
echo "== Read the SIGNADOT_DEMO_SCENARIO.md for instructions on how to run the demo =="