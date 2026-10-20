#!/usr/bin/env bash
set -euo pipefail

# Tear down a PR preview environment created by deploy.sh. Runs as root (via
# sudo) from the cleanup-preview CI job when a PR closes. Idempotent: every
# removal tolerates an already-absent target, so a double-close or a
# close-after-failed deploy does not error.
#
# DRY_RUN=true only logs what would be removed and needs no root:
#   DRY_RUN=true bash cleanup.sh <PR_NUMBER>-merge
# sudo drops DRY_RUN, so `DRY_RUN=true sudo cleanup.sh` removes the preview.
DRY_RUN="${DRY_RUN:-false}"

DATA_ROOT="/mnt/HC_Volume_103911252/environments"

SAFE_NAME="${1:-}"

if [[ -z "$SAFE_NAME" ]]; then
  echo "[cleanup] ERROR: SAFE_NAME (first argument) is not set." >&2
  exit 1
fi

if [[ ! "$SAFE_NAME" =~ ^[0-9]+-merge$ ]]; then
  echo "[cleanup] ERROR: refusing to clean up SAFE_NAME=$SAFE_NAME; only PR previews named <PR_NUMBER>-merge can be removed." >&2
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
# A failed reload leaves the old route live in Caddy, so let it fail the job.
run systemctl reload caddy

# Remove the app dir last; it holds this very script. On Linux the running
# process keeps its open inode, so removing the directory mid-run is safe.
echo "[cleanup] Removing app dir: $APP_DIR"
run rm -rf "$APP_DIR"

if [[ "$DRY_RUN" == "true" ]]; then
  echo "[cleanup] DRY RUN complete — nothing was removed for $SAFE_NAME."
else
  echo "[cleanup] Done. Preview $SAFE_NAME removed."
fi
