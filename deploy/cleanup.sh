#!/usr/bin/env bash
set -euo pipefail

# Tear down a preview environment created by deploy.sh. Invoked as root (via
# sudo) from the cleanup-preview CI job when a PR closes. Idempotent: every
# removal tolerates an already-absent target, so a double-close or a
# close-after-failed-deploy does not error.
#
# SAFETY — DRY RUN BY DEFAULT:
# DRY_RUN defaults to "true", so this script only LOGS what it would remove.
# That makes the first PR close after rollout a harmless no-op which still
# proves the CI -> server wiring (SSH, sudo, upload, script run) end to end.
# Once you have seen a clean dry-run log AND patched the live server's sudoers
# for cleanup.sh, flip the default below to "false" (one-line change, committed)
# to enable real teardown. You can also force a single real run manually with
# `DRY_RUN=false sudo cleanup.sh <safe_name>`.
DRY_RUN="${DRY_RUN:-true}"

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

# run executes a command, or just logs it when DRY_RUN is true.
run() {
  if [[ "$DRY_RUN" == "true" ]]; then
    echo "[cleanup] DRY RUN — would run: $*"
  else
    "$@"
  fi
}

if [[ "$DRY_RUN" == "true" ]]; then
  echo "[cleanup] DRY RUN enabled — no changes will be made. SAFE_NAME=$SAFE_NAME"
else
  echo "[cleanup] Tearing down preview SAFE_NAME=$SAFE_NAME"
fi

echo "[cleanup] Stopping and disabling $SERVICE_NAME"
run systemctl stop "$SERVICE_NAME" || true
run systemctl disable "$SERVICE_NAME" || true

echo "[cleanup] Removing systemd unit: $SERVICE_UNIT"
run rm -f "$SERVICE_UNIT"

echo "[cleanup] Removing Caddy site: $CADDY_SITE_DEST"
run rm -f "$CADDY_SITE_DEST"

echo "[cleanup] Removing branch data dir: $BRANCH_DATA_DIR"
run rm -rf "$BRANCH_DATA_DIR"

echo "[cleanup] Reloading systemd and Caddy"
run systemctl daemon-reload || true
run systemctl reload caddy || true

# Remove the app dir last; it holds this very script. On Linux the running
# process keeps its open inode, so removing the directory mid-run is safe.
echo "[cleanup] Removing app dir: $APP_DIR"
run rm -rf "$APP_DIR"

if [[ "$DRY_RUN" == "true" ]]; then
  echo "[cleanup] DRY RUN complete — nothing was removed for $SAFE_NAME."
else
  echo "[cleanup] Done. Preview $SAFE_NAME removed."
fi
