#!/usr/bin/env bash
set -euo pipefail

APP_DIR="/opt/conorganizer"

# Runtime user/group for the service process
SERVICE_USER="conorganizer"
SERVICE_GROUP="www-data"

NEW_BIN="$APP_DIR/conorganizer.new"
CUR_BIN="$APP_DIR/conorganizer"
OLD_BIN="$APP_DIR/conorganizer.old"

SERVICE_NAME="conorganizer.service"
SERVICE_PATH="$APP_DIR/conorganizer.service"
SYSTEMD_UNIT="/etc/systemd/system/${SERVICE_NAME}"

echo "[deploy] Using APP_DIR=$APP_DIR"

# Sanity check: new binary must exist
if [[ ! -f "$NEW_BIN" ]]; then
  echo "[deploy] ERROR: new binary not found at $NEW_BIN" >&2
  exit 1
fi

# Install/overwrite systemd unit if present in APP_DIR
if [[ -f "$SERVICE_PATH" ]]; then
  echo "[deploy] Installing systemd unit: $SYSTEMD_UNIT"
  mv "$SERVICE_PATH" "$SYSTEMD_UNIT"
  chmod 644 "$SYSTEMD_UNIT"
fi

# Ensure ownership of the new binary (optional but nice)
chown "$SERVICE_USER:$SERVICE_GROUP" "$NEW_BIN" || true

# Backup current binary if it exists
if [[ -f "$CUR_BIN" ]]; then
  echo "[deploy] Backing up current binary to $OLD_BIN"
  mv "$CUR_BIN" "$OLD_BIN"
fi

echo "[deploy] Promoting new binary"
mv "$NEW_BIN" "$CUR_BIN"
chmod +x "$CUR_BIN"
chown "$SERVICE_USER:$SERVICE_GROUP" "$CUR_BIN" || true

echo "[deploy] Reloading systemd and restarting $SERVICE_NAME"
systemctl daemon-reload || true
systemctl restart "$SERVICE_NAME"

echo "[deploy] Checking service status..."
if systemctl.is-active --quiet "$SERVICE_NAME"; then
  echo "[deploy] Service is active."
else
  echo "[deploy] ERROR: service failed to start" >&2
  # Show recent logs to help debugging from CI logs
  journalctl -u "$SERVICE_NAME" -n 50 --no-pager || true
  exit 1
fi
