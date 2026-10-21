#!/usr/bin/env bash
#
echo "Setup a restricted 'deploy' user for GitHub Actions deployments on Ubuntu."
echo "Generates an SSH keypair automatically."

set -euo pipefail

DEPLOY_USER="deploy"
DEPLOY_GROUP="deploy"
APP_DIR="/opt/conorganizer"
SSHD_CONFIG="/etc/ssh/sshd_config"
SUDOERS_FILE="/etc/sudoers.d/deploy-conorganizer"
SSH_SERVICE="ssh"   # Ubuntu default

# Fresh servers only: on an existing account every chown, chmod and cp below
# would follow whatever the deploy user has put in its home.
if id -u "$DEPLOY_USER" >/dev/null 2>&1; then
  echo "[setup] User '$DEPLOY_USER' already exists. This script is for a fresh server; edit $SUDOERS_FILE with 'visudo -f' instead (hetzner/hetzner-admin.md)." >&2
  exit 1
fi

echo "[setup] Creating user '$DEPLOY_USER'"
adduser --system --home "$APP_DIR" --group "$DEPLOY_USER"

echo "[setup] Ensuring '$DEPLOY_USER' has a usable shell (/bin/bash)"
usermod -s /bin/bash "$DEPLOY_USER"

echo "[setup] Ensuring $APP_DIR exists and is owned by $DEPLOY_USER"
mkdir -p "$APP_DIR"
chown -R "$DEPLOY_USER:$DEPLOY_GROUP" "$APP_DIR"

DEPLOY_HOME="$APP_DIR"
SSH_DIR="$DEPLOY_HOME/.ssh"
AUTH_KEYS="$SSH_DIR/authorized_keys"

echo "[setup] Setting up SSH directory at $SSH_DIR"
mkdir -p "$SSH_DIR"
chmod 700 "$SSH_DIR"
chown "$DEPLOY_USER:$DEPLOY_GROUP" "$SSH_DIR"

# Generate keypair only if it doesn't already exist
KEY_FILE="$SSH_DIR/id_ed25519"
PUB_KEY_FILE="$SSH_DIR/id_ed25519.pub"

if [[ -f "$KEY_FILE" ]]; then
  echo "[setup] SSH keypair already exists at $KEY_FILE, skipping generation"
else
  echo "[setup] Generating SSH keypair for $DEPLOY_USER"
  sudo -u "$DEPLOY_USER" ssh-keygen -t ed25519 -N "" -f "$KEY_FILE"
  chmod 600 "$KEY_FILE"
  chmod 644 "$PUB_KEY_FILE"
fi

echo "[setup] Configuring authorized_keys"
cp "$PUB_KEY_FILE" "$AUTH_KEYS"
chmod 600 "$AUTH_KEYS"
chown "$DEPLOY_USER:$DEPLOY_GROUP" "$AUTH_KEYS"

# SSH hardening: key-only for deploy user
if ! grep -q "Match User $DEPLOY_USER" "$SSHD_CONFIG"; then
  echo "[setup] Adding Match block for $DEPLOY_USER in $SSHD_CONFIG"
  cp "$SSHD_CONFIG" "${SSHD_CONFIG}.bak.$(date +%s)"

  cat >> "$SSHD_CONFIG" <<EOF

Match User $DEPLOY_USER
    PasswordAuthentication no
    PubkeyAuthentication yes
EOF

  echo "[setup] Reloading $SSH_SERVICE"
  systemctl reload "$SSH_SERVICE" || systemctl restart "$SSH_SERVICE" || true
else
  echo "[setup] Match block for $DEPLOY_USER already present in $SSHD_CONFIG"
fi

# Minimal sudo rights. This writes the final state only; an existing sudoers
# file is edited with `visudo -f` instead.
echo "[setup] Writing sudoers file: $SUDOERS_FILE"
SUDOERS_TMP="$(mktemp)"
cat > "$SUDOERS_TMP" <<EOF
# Minimal sudo rights for GitHub Actions deploy user 'deploy'.
# Only these two scripts may run as root. Stow installs them from the admin
# checkout /srv/configuration-as-code-repo/conorganizer, which root already
# trusts for the backup units. They validate the environment name themselves,
# so nothing the deploy user can write ever runs as root.
deploy ALL=(root) NOPASSWD: \
    /usr/local/bin/conorganizer-deploy, \
    /usr/local/bin/conorganizer-cleanup
EOF

visudo -cf "$SUDOERS_TMP"
install -m 0440 -o root -g root "$SUDOERS_TMP" "$SUDOERS_FILE"
rm -f "$SUDOERS_TMP"
sudo -l -U "$DEPLOY_USER"

# conorganizer-deploy refuses preview deploys unless the data root is a
# root-owned directory that only root can write.
DATA_ROOT="/mnt/HC_Volume_103911252/environments"
if mountpoint -q "$(dirname "$DATA_ROOT")"; then
  echo "[setup] Ensuring $DATA_ROOT is root-owned and writable only by root"
  install -d -o root -g root -m 0755 "$DATA_ROOT"
else
  echo "[setup] WARNING: $(dirname "$DATA_ROOT") is not mounted; create $DATA_ROOT as root:root 0755 after mounting it" >&2
fi

echo
echo "=============================================================="
echo " SSH KEYPAIR FOR GITHUB ACTIONS DEPLOYMENT"
echo "=============================================================="
echo
echo "PRIVATE KEY (paste this into GitHub Secret: HETZNER_SSH_KEY):"
echo "--------------------------------------------------------------"
cat "$KEY_FILE"
echo
echo "PUBLIC KEY (already installed in authorized_keys):"
echo "--------------------------------------------------------------"
cat "$PUB_KEY_FILE"
echo
echo "=============================================================="
echo " Add these GitHub Secrets in your repository:"
echo "   HETZNER_HOST     = <your server IP or hostname>"
echo "   HETZNER_USER     = $DEPLOY_USER"
echo "   HETZNER_SSH_KEY  = (private key above)"
echo "   HETZNER_SSH_PORT = 22  (or your custom port)"
echo "=============================================================="
