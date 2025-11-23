#!/usr/bin/env bash
set -euo pipefail

APP_DIR="/opt/conorganizer"
NEW_BIN="$APP_DIR/conorganizer.new"
CUR_BIN="$APP_DIR/conorganizer"
OLD_BIN="$APP_DIR/conorganizer.old"
SERVICE_NAME="conorganizer.service"
SERVICE_PATH="$APP_DIR/conorganizer.service"
SYSTEMD_UNIT="/etc/systemd/system/conorganizer.service"

echo "[deploy] Using APP_DIR=$APP_DIR"

if [[ ! -f "$NEW_BIN" ]]; then
  echo "[deploy] ERROR: new binary not found at $NEW_BIN" >&2
  exit 1
fi

# Install/overwrite systemd unit if present
if [[ -f "$SERVICE_PATH" ]]; then
  echo "[deploy] Installing systemd unit: $SYSTEMD_UNIT"
  sudo mv "$SERVICE_PATH" "$SYSTEMD_UNIT"
  sudo chmod 644 "$SYSTEMD_UNIT"
fi

# Ensure ownership (optional)
sudo chown deploy:deploy "$NEW_BIN" || true

# Backup current binary if it exists
if [[ -f "$CUR_BIN" ]]; then
  echo "[deploy] Backing up current binary to $OLD_BIN"
  sudo mv "$CUR_BIN" "$OLD_BIN"
fi

echo "[deploy] Promoting new binary"
sudo mv "$NEW_BIN" "$CUR_BIN"
sudo chmod +x "$CUR_BIN"
sudo chown deploy:deploy "$CUR_BIN" || true

echo "[deploy] Restarting service: $SERVICE_NAME"
sudo systemctl daemon-reload || true
sudo systemctl restart "$SERVICE_NAME"

echo "[deploy] Checking service status..."
if sudo systemctl is-active --quiet "$SERVICE_NAME"; then
  echo "[deploy] Service is active."
else
  echo "[deploy] ERROR: service failed to start" >&2
  exit 1
fi
