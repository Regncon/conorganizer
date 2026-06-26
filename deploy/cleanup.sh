#!/usr/bin/env bash
set -euo pipefail

# Tear down a preview environment created by deploy.sh. Invoked as root (via
# sudo) from the cleanup-preview CI job when a PR closes. Idempotent: every
# removal tolerates an already-absent target, so a double-close or a
# close-after-failed-deploy does not error.

DATA_ROOT="/mnt/HC_Volume_103911252/environments"

SAFE_NAME="${1:-}"

if [[ -z "$SAFE_NAME" ]]; then
  echo "[cleanup] ERROR: SAFE_NAME (first argument) is not set." >&2
  exit 1
fi

if [[ "$SAFE_NAME" == "main" ]]; then
  echo "[cleanup] ERROR: refusing to clean up SAFE_NAME=main (production)." >&2
  exit 1
fi

SERVICE_NAME="conorganizer-${SAFE_NAME}.service"
SERVICE_UNIT="/etc/systemd/system/$SERVICE_NAME"
CADDY_SITE_DEST="/etc/caddy/sites-enabled/conorganizer-${SAFE_NAME}.caddy"
APP_DIR="/opt/conorganizer/${SAFE_NAME}"
BRANCH_DATA_DIR="$DATA_ROOT/$SAFE_NAME"

echo "[cleanup] Tearing down preview SAFE_NAME=$SAFE_NAME"

echo "[cleanup] Stopping and disabling $SERVICE_NAME"
systemctl stop "$SERVICE_NAME" 2>/dev/null || true
systemctl disable "$SERVICE_NAME" 2>/dev/null || true

echo "[cleanup] Removing systemd unit: $SERVICE_UNIT"
rm -f "$SERVICE_UNIT"

echo "[cleanup] Removing Caddy site: $CADDY_SITE_DEST"
rm -f "$CADDY_SITE_DEST"

echo "[cleanup] Removing branch data dir: $BRANCH_DATA_DIR"
rm -rf "$BRANCH_DATA_DIR"

echo "[cleanup] Reloading systemd and Caddy"
systemctl daemon-reload || true
systemctl reload caddy || true

# Remove the app dir last; it holds this very script. On Linux the running
# process keeps its open inode, so removing the directory mid-run is safe.
echo "[cleanup] Removing app dir: $APP_DIR"
rm -rf "$APP_DIR"

echo "[cleanup] Done. Preview $SAFE_NAME removed."
