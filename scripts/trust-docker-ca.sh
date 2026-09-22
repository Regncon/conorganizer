#!/usr/bin/env sh
set -eu

script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
project_root=$(dirname "$script_dir")
certificate_directory="$project_root/tmp/caddy"
certificate_path="$certificate_directory/conorganizer-caddy-root.crt"

mkdir -p "$certificate_directory"
cd "$project_root"

docker compose cp caddy:/data/caddy/pki/authorities/local/root.crt "$certificate_path"

case "$(uname -s)" in
    Linux)
        sudo install -m 0644 "$certificate_path" /usr/local/share/ca-certificates/conorganizer-caddy-root.crt
        sudo update-ca-certificates
        ;;
    Darwin)
        sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain "$certificate_path"
        ;;
    *)
        echo "Unsupported operating system. Trust this certificate manually: $certificate_path" >&2
        exit 1
        ;;
esac

printf '%s\n' \
    "Trusted the Conorganizer Docker certificate." \
    "Open https://localhost on the port configured by HTTPS_PORT (default: 7331)."
