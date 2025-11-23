#!/usr/bin/env bash
#
# Setup a restricted 'deploy' user for GitHub Actions deployments.
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
DEPLOY_HOME="$APP_DIR"  # deploy user will use /opt/conorganizer as home

echo "[setup] Creating user '$DEPLOY_USER' (if missing)"

if ! id -u "$DEPLOY_USER" >/dev/null 2>&1; then
  adduser --system --home "$DEPLOY_HOME" --group "$DEPLOY_USER"
else
  echo "[setup] User already exists"
fi

mkdir -p "$APP_DIR"
chown -R "$DEPLOY_USER:$DEPLOY_GROUP" "$APP_DIR"

SSH_DIR="$DEPLOY_HOME/.ssh"
mkdir -p "$SSH_DIR"
chmod 700 "$SSH_DIR"
chown "$DEPLOY_USER:$DEPLOY_GROUP" "$SSH_DIR"

echo "[setup] Generating SSH keypair for deploy user"
sudo -u "$DEPLOY_USER" ssh-keygen -t ed25519 -N "" -f "$SSH_DIR/id_ed25519"

chmod 600 "$SSH_DIR/id_ed25519"
chmod 644 "$SSH_DIR/id_ed25519.pub"

echo "[setup] Configuring authorized_keys"
cp "$SSH_DIR/id_ed25519.pub" "$SSH_DIR/authorized_keys"
chmod 600 "$SSH_DIR/authorized_keys"
chown "$DEPLOY_USER:$DEPLOY_GROUP" "$SSH_DIR/authorized_keys"

# SSH hardening
if ! grep -q "Match User $DEPLOY_USER" "$SSHD_CONFIG"; then
  echo "[setup] Adding sshd Match block"
  cp "$SSHD_CONFIG" "${SSHD_CONFIG}.bak.$(date +%s)"

  cat >> "$SSHD_CONFIG" <<EOF

Match User $DEPLOY_USER
    PasswordAuthentication no
    PubkeyAuthentication yes
EOF

  systemctl reload sshd || systemctl restart sshd
fi

# Sudo rights (minimal)
echo "[setup] Installing sudoers file: $SUDOERS_FILE"

cat > "$SUDOERS_FILE" <<EOF
# Minimal sudo rights for deploy user
$DEPLOY_USER ALL=(root) NOPASSWD: \
  /bin/mv /opt/conorganizer/*, \
  /bin/chown $DEPLOY_USER:$DEPLOY_GROUP /opt/conorganizer/*, \
  /bin/chmod * /opt/conorganizer/*, \
  /bin/mv /opt/conorganizer/conorganizer.service /etc/systemd/system/conorganizer.service, \
  /bin/chmod 644 /etc/systemd/system/conorganizer.service, \
  /bin/systemctl daemon-reload, \
  /bin/systemctl restart conorganizer.service, \
  /bin/systemctl enable conorganizer.service, \
  /bin/systemctl is-active conorganizer.service
EOF

chmod 440 "$SUDOERS_FILE"
visudo -cf "$SUDOERS_FILE"

# OUTPUT KEYS FOR GITHUB ACTIONS
echo
echo "=============================================================="
echo " SSH KEYPAIR GENERATED FOR GITHUB ACTIONS DEPLOYMENT"
echo "=============================================================="
echo
echo "PRIVATE KEY (paste this into GitHub Secret: HETZNER_SSH_KEY):"
echo "--------------------------------------------------------------"
cat "$SSH_DIR/id_ed25519"
echo
echo "PUBLIC KEY (already installed on server):"
echo "--------------------------------------------------------------"
cat "$SSH_DIR/id_ed25519.pub"
echo
echo "=============================================================="
echo " Add these GitHub Secrets:"
echo "   HETZNER_HOST     = <your server IP>"
echo "   HETZNER_USER     = $DEPLOY_USER"
echo "   HETZNER_SSH_KEY  = (private key above)"
echo "   HETZNER_SSH_PORT = 22  (or your custom port)"
echo "=============================================================="
