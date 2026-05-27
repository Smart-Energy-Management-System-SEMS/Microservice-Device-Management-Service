#!/usr/bin/env sh
set -eu

APP_BINARY="${APP_BINARY:-./app}"
KEEP_ALIVE_ENABLED="${KEEP_ALIVE_ENABLED:-false}"
KEEP_ALIVE_INTERVAL_SECONDS="${KEEP_ALIVE_INTERVAL_SECONDS:-600}"
KEEP_ALIVE_PATH="${KEEP_ALIVE_PATH:-/api/v1/device-management/health}"

start_keep_alive() {
  if [ "$KEEP_ALIVE_ENABLED" != "true" ]; then
    echo "keep-alive disabled"
    return
  fi

  BASE_URL="${KEEP_ALIVE_URL:-${RENDER_EXTERNAL_URL:-}}"
  if [ -z "$BASE_URL" ]; then
    echo "keep-alive disabled: KEEP_ALIVE_URL or RENDER_EXTERNAL_URL is required"
    return
  fi

  TARGET_URL="${BASE_URL%/}${KEEP_ALIVE_PATH}"
  echo "keep-alive enabled: pinging ${TARGET_URL} every ${KEEP_ALIVE_INTERVAL_SECONDS}s"

  (
    while true; do
      sleep "$KEEP_ALIVE_INTERVAL_SECONDS"
      if command -v curl >/dev/null 2>&1; then
        curl -fsS "$TARGET_URL" >/dev/null 2>&1 || true
      elif command -v wget >/dev/null 2>&1; then
        wget -qO- "$TARGET_URL" >/dev/null 2>&1 || true
      else
        echo "keep-alive stopped: curl or wget is required"
        exit 0
      fi
    done
  ) &
}

start_keep_alive
exec "$APP_BINARY"
