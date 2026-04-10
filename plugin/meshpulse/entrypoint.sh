#!/bin/bash
set -e

echo "=== MeshPulse DePIN Node ==="

# Start canopy node in background (data dir already initialized)
echo "[1/3] Starting Canopy node..."
./bin start &
CANOPY_PID=$!

# Wait for the plugin socket to appear
echo "[2/3] Waiting for plugin socket..."
TIMEOUT=90
COUNT=0
while [ ! -S "/tmp/plugin/plugin.sock" ] && [ $COUNT -lt $TIMEOUT ]; do
    sleep 1
    COUNT=$((COUNT + 1))
    if [ $((COUNT % 10)) -eq 0 ]; then
        echo "  ... waiting ${COUNT}s"
    fi
done

if [ ! -S "/tmp/plugin/plugin.sock" ]; then
    echo "ERROR: plugin.sock not found after ${TIMEOUT}s — is the plugin configured?"
    exit 1
fi

echo "[3/3] Starting MeshPulse plugin..."
./plugin/meshpulse/meshpulse &

echo ""
echo "✅ MeshPulse ready at http://localhost:8080"
echo ""

wait $CANOPY_PID
