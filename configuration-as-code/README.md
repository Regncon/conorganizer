# Configuration as Code

Server configuration for the Conorganizer production server.

This directory uses GNU Stow to symlink configuration files into the root filesystem.

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
