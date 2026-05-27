#!/usr/bin/env sh
set -eu

APP_BINARY="${APP_BINARY:-./app}"
KEEP_ALIVE_ENABLED="${KEEP_ALIVE_ENABLED:-false}"

start_keep_alive() {
  if [ "$KEEP_ALIVE_ENABLED" != "true" ]; then
    echo "keep-alive disabled"
    return
  fi

  sh scripts/keepalive.sh &
}

start_keep_alive
exec "$APP_BINARY"
