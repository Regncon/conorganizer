# Configuration as Code

Server configuration for the Conorganizer production server.

This directory uses GNU Stow to symlink configuration files into the root filesystem.

## Repository location

The production repository is stored at:

```text
/srv/configuration-as-code-repo/conorganizer/configuration-as-code
```

## Install

From the repository root:

```bash
./configuration-as-code/install.sh
```

## Fix permissions
```bash
sudo find configuration-as-code/stow -type d -exec chmod 755 {} \;
sudo find configuration-as-code/stow -type f -exec chmod 644 {} \;
```


## Find all stowed files

List symlinks below `/etc` and `/usr/local` that resolve into this repository:

```bash
sudo find /etc /usr/local -type l -print0 |
while IFS= read -r -d '' symlink_path; do
    resolved_target="$(readlink -f -- "$symlink_path" 2>/dev/null || true)"

    case "$resolved_target" in
        /srv/configuration-as-code-repo/conorganizer/configuration-as-code/stow/*)
            printf '%s -> %s\n' "$symlink_path" "$resolved_target"
            ;;
    esac
done
```

