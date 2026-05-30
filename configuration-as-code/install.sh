#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/stow"

for package in scripts conorganizer caddy grafana prometheus; do
    if [[ -d "$package" ]]; then
        sudo stow --target=/ "$package"
    else
        echo "Skipping missing stow package: $package"
    fi
done

sudo systemctl daemon-reload
