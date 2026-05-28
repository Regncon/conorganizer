#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/stow"

sudo stow --target=/ scripts
sudo stow --target=/ conorganizer
sudo stow --target=/ caddy
sudo stow --target=/ grafana

sudo systemctl daemon-reload

