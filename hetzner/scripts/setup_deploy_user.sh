#!/usr/bin/env bash
#
# Setup a restricted 'deploy' user for GitHub Actions deployments on Ubuntu.
# Generates an SSH keypair automatically.
#
# Usage:
#   sudo ./setup_deploy_user.sh
#

set -euo pipefail

DEPLOY_USER="deploy"
DEPLOY_GROUP="deploy"
APP_DIR="/opt/conorganizer"
SSHD_CONFIG="/etc/ssh/sshd_config"
SUDOERS_FILE="/etc/sudoers.d/deploy-conorganizer"
SSH_SERVICE="ssh"   # Ubuntu default

echo "[setup] Creating user '$DEPLOY_USER' (if missing)"

if ! id -u "$DEPLOY_USER" >/dev/null 2>&1; then
  # system user with specified home
  adduser --system --home "$APP_DIR" --group "$DEPLOY_USER"
else
  echo "[setup] User '$DEPLOY_USER' already exists, skipping creation"
fi

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

# Minimal sudo rights
echo "[setup] Writing sudoers file: $SUDOERS_FILE"
cat > "$SUDOERS_FILE" <<EOF
# Minimal sudo rights for GitHub Actions deploy user 'deploy'
deploy ALL=(root) NOPASSWD: \
    /usr/bin/mv, \
    /usr/bin/chown, \
    /usr/bin/chmod, \
    /usr/bin/systemctl, \
    /opt/conorganizer/*/deploy.sh
EOF

chmod 440 "$SUDOERS_FILE"
visudo -cf "$SUDOERS_FILE"

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
