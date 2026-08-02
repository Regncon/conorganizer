#!/usr/bin/env bash
set -euo pipefail

readonly repository_directory="$(cd "$(dirname "$0")" && pwd)"
readonly stow_directory="${repository_directory}/stow"

packages=(
    scripts
    caddy
    grafana
    loki
    prometheus
    promtail
    systemd
)

for package in "${packages[@]}"; do
    package_directory="${stow_directory}/${package}"

    if [[ -d "$package_directory" ]]; then
        sudo stow \
            --dir="$stow_directory" \
            --target=/ \
            --restow \
            "$package"
    else
        echo "Skipping missing Stow package: $package"
    fi
done

sudo systemctl daemon-reload
