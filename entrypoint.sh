#!/bin/bash
set -euo pipefail

export DISPLAY=:99
Xvfb "$DISPLAY" -screen 0 1920x1080x24 -nolisten tcp &
XVFB_PID=$!
trap 'kill "$XVFB_PID" 2>/dev/null' EXIT

for _ in $(seq 1 50); do
  if [ -e /tmp/.X11-unix/X99 ]; then
    break
  fi
  sleep 0.1
done

exec /app/lidl-coupons
