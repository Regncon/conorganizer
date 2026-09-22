#!/usr/bin/env sh
set -eu

script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
project_root=$(dirname "$script_dir")
certificate_directory="$project_root/tmp/caddy"
certificate_path="$certificate_directory/conorganizer-caddy-root.crt"
certificate_nickname="Conorganizer Docker Caddy CA"

mkdir -p "$certificate_directory"
cd "$project_root"

docker compose cp caddy:/data/caddy/pki/authorities/local/root.crt "$certificate_path"

# Chrome and Firefox on Linux read their own NSS databases instead of the
# system store, so the certificate is added there as well when possible.
trust_in_nss_databases() {
    if ! command -v certutil >/dev/null 2>&1; then
        printf '%s\n' \
            "certutil was not found, so browsers may still reject the certificate." \
            "Install libnss3-tools (Debian/Ubuntu), nss-tools (Fedora) or nss (Arch) and run this script again." >&2
        return
    fi

    for nss_directory in "$HOME/.pki/nssdb" "$HOME"/.mozilla/firefox/*/ "$HOME"/snap/firefox/common/.mozilla/firefox/*/; do
        if [ -f "$nss_directory/cert9.db" ]; then
            certutil -d "sql:$nss_directory" -D -n "$certificate_nickname" >/dev/null 2>&1 || true
            certutil -d "sql:$nss_directory" -A -t "C,," -n "$certificate_nickname" -i "$certificate_path"
            echo "Trusted the certificate in $nss_directory"
        fi
    done
}

case "$(uname -s)" in
    Linux)
        if command -v update-ca-certificates >/dev/null 2>&1; then
            sudo install -m 0644 "$certificate_path" /usr/local/share/ca-certificates/conorganizer-caddy-root.crt
            sudo update-ca-certificates
        elif command -v update-ca-trust >/dev/null 2>&1; then
            sudo install -m 0644 "$certificate_path" /etc/pki/ca-trust/source/anchors/conorganizer-caddy-root.crt
            sudo update-ca-trust
        elif command -v trust >/dev/null 2>&1; then
            sudo trust anchor --store "$certificate_path"
        else
            echo "No supported system certificate tool was found. Trust this certificate manually: $certificate_path" >&2
            exit 1
        fi
        trust_in_nss_databases
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
    "Open https://localhost on the port configured by HTTPS_PORT (default: 7331)." \
    "To trust it on a phone, install this file on the phone: $certificate_path"
