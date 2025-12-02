#!/usr/bin/env bash
set -euo pipefail

# Base directory for all branch deployments
ROOT_DIR="/opt/conorganizer"

# Data root on the mounted volume
DATA_ROOT="/mnt/HC_Volume_103911252"
MAIN_DATA_DIR="$DATA_ROOT/main"

# Branch-safe name is passed as first argument from CI, e.g. "main", "274-..."
SAFE_NAME="${1:-}"

if [[ -z "$SAFE_NAME" ]]; then
  echo "[deploy] ERROR: SAFE_NAME (first argument) is not set." >&2
  exit 1
fi

APP_DIR="$ROOT_DIR/$SAFE_NAME"

# Per-branch binary name: conorganizer-main, conorganizer-feature-x, etc.
BIN_NAME="conorganizer-${SAFE_NAME}"
NEW_BIN_SRC="$ROOT_DIR/conorganizer.new"
CUR_BIN="$APP_DIR/$BIN_NAME"
OLD_BIN="$APP_DIR/${BIN_NAME}.old"

# Per-branch systemd service
SERVICE_NAME="conorganizer-${SAFE_NAME}.service"
SERVICE_SRC="$ROOT_DIR/$SERVICE_NAME"
SERVICE_UNIT="/etc/systemd/system/$SERVICE_NAME"

# Per-branch Caddy site
CADDY_SITE_SRC="$ROOT_DIR/caddy-${SAFE_NAME}.caddy"
CADDY_SITE_DEST="/etc/caddy/sites-enabled/conorganizer-${SAFE_NAME}.caddy"

# Per-branch data paths
BRANCH_DATA_DIR="$DATA_ROOT/$SAFE_NAME"
BRANCH_DB_DIR="$BRANCH_DATA_DIR/database"
BRANCH_IMG_DIR="$BRANCH_DATA_DIR/event-images"

# Runtime user/group
SERVICE_USER="deploy"
SERVICE_GROUP="www-data"

echo "[deploy] Deploying branch SAFE_NAME=$SAFE_NAME"
echo "[deploy] APP_DIR=$APP_DIR"
echo "[deploy] BIN_NAME=$BIN_NAME"

# --- Sanity checks on input files ---

if [[ ! -f "$NEW_BIN_SRC" ]]; then
  echo "[deploy] ERROR: new binary not found at $NEW_BIN_SRC" >&2
  exit 1
fi

if [[ ! -f "$SERVICE_SRC" ]]; then
  echo "[deploy] ERROR: service file not found at $SERVICE_SRC" >&2
  exit 1
fi

if [[ ! -f "$CADDY_SITE_SRC" ]]; then
  echo "[deploy] ERROR: Caddy site file not found at $CADDY_SITE_SRC" >&2
  exit 1
fi

# --- Ensure per-branch data dirs (copy from main if needed) ---

if [[ "$SAFE_NAME" != "main" ]]; then
  if [[ ! -d "$MAIN_DATA_DIR" ]]; then
    echo "[deploy] ERROR: main data dir $MAIN_DATA_DIR does not exist; cannot clone data." >&2
    exit 1
  fi

  echo "[deploy] Ensuring data directory for branch: $BRANCH_DATA_DIR"
  mkdir -p "$BRANCH_DATA_DIR"

  if [[ ! -d "$BRANCH_DB_DIR" ]]; then
    echo "[deploy] Copying database from main to $BRANCH_DB_DIR"
    mkdir -p "$BRANCH_DATA_DIR"
    cp -a "$MAIN_DATA_DIR/database" "$BRANCH_DATA_DIR/"
  else
    echo "[deploy] Database dir already exists for branch: $BRANCH_DB_DIR (skipping copy)"
  fi

  if [[ ! -d "$BRANCH_IMG_DIR" ]]; then
    echo "[deploy] Copying event-images from main to $BRANCH_IMG_DIR"
    mkdir -p "$BRANCH_DATA_DIR"
    cp -a "$MAIN_DATA_DIR/event-images" "$BRANCH_DATA_DIR/"
  else
    echo "[deploy] Event-images dir already exists for branch: $BRANCH_IMG_DIR (skipping copy)"
  fi
else
  echo "[deploy] SAFE_NAME=main, not cloning data directories."
fi

# --- Prepare app directory ---

echo "[deploy] Ensuring app directory exists: $APP_DIR"
mkdir -p "$APP_DIR"
chown -R "$SERVICE_USER:$SERVICE_GROUP" "$APP_DIR"

# --- Install / update systemd unit ---

echo "[deploy] Installing systemd unit: $SERVICE_UNIT"
mv "$SERVICE_SRC" "$SERVICE_UNIT"
chmod 644 "$SERVICE_UNIT"

# --- Install / update Caddy site ---

echo "[deploy] Installing Caddy site: $CADDY_SITE_DEST"
mkdir -p /etc/caddy/sites-enabled
mv "$CADDY_SITE_SRC" "$CADDY_SITE_DEST"
chmod 644 "$CADDY_SITE_DEST"

echo "[deploy] Reloading Caddy"
systemctl reload caddy

# --- Promote new binary ---

# Backup current binary if it exists
if [[ -f "$CUR_BIN" ]]; then
  echo "[deploy] Backing up current binary to $OLD_BIN"
  mv "$CUR_BIN" "$OLD_BIN"
fi

echo "[deploy] Promoting new binary to $CUR_BIN"
mv "$NEW_BIN_SRC" "$CUR_BIN"
chmod +x "$CUR_BIN"
chown "$SERVICE_USER:$SERVICE_GROUP" "$CUR_BIN" || true

# --- Restart systemd service ---

echo "[deploy] Reloading systemd and restarting $SERVICE_NAME"
systemctl daemon-reload || true
systemctl enable "$SERVICE_NAME" || true
systemctl restart "$SERVICE_NAME"

echo "[deploy] Checking service status..."
if systemctl is-active --quiet "$SERVICE_NAME"; then
  echo "[deploy] Service is active."
else
  echo "[deploy] ERROR: service failed to start" >&2
  journalctl -u "$SERVICE_NAME" -n 50 --no-pager || true
  exit 1
fi

echo "[deploy] Done. Service: $SERVICE_NAME, Binary: $CUR_BIN"
