# Hetzner admin documentation

## Add a user
To add a user to the Hetzner cloud VM, run the script as root from the repository
checkout on the server:

```bash
sudo bash hetzner/scripts/addUser.sh <username>
```

The script:

- creates the user with a home directory and `/bin/bash` as shell
- adds the user to the groups `adm`, `sudo`, `www-data`, `users` and `deploy`
- generates a random password (needs `pwgen`), prints it, and expires it so
  the user must change it at first login
- links `/opt/conorganizer` into the user's home directory

To do the same by hand:

```bash
adduser <username>
usermod -aG adm,sudo,www-data,users,deploy <username>
passwd <username>
```

Set up SSH keys
**ToDo: add more details here**

A user who runs `go tool task download:*` locally (as `DB_SSH_USER`) must be
able to run `conorganizer-export-db` and read the `main` and `demo`
database files and `event-images` directories listed below.

## Deploy user
`hetzner/scripts/setup_deploy_user.sh` creates the restricted `deploy` user
that GitHub Actions deploys as. Its home is `/opt/conorganizer`, SSH login is
key-only, and sudo is limited to the two scripts
`/usr/local/bin/conorganizer-deploy` and `/usr/local/bin/conorganizer-cleanup`,
which Stow installs from the admin-owned server checkout
`/srv/configuration-as-code-repo/conorganizer`. CI uploads only the binary,
`goose` and the migrations into `/opt/conorganizer/<env>/`; everything that
runs as root comes from that checkout, so the deploy user can trigger a deploy
or a cleanup but cannot change what either does (see
[Service user and ownership](../documentation/deployment.md#service-user-and-ownership)).
The script is for a fresh server and refuses to run when the `deploy` user
already exists, because its chown, chmod and cp steps would follow whatever
that user has put in its home. It validates the sudoers file with `visudo -cf`
before installing it. On a server that is set up, edit the live file with
`sudo visudo -f /etc/sudoers.d/deploy-conorganizer` and check it with
`sudo -l -U deploy`. The script
generates an SSH keypair and prints the values for the `HETZNER_HOST`,
`HETZNER_USER`, `HETZNER_SSH_KEY` and `HETZNER_SSH_PORT` GitHub secrets.
The workflow only reads the first three; it hard-codes `port: 22`
(`.github/workflows/buildAndTest.yml`), so a custom SSH port also needs a
workflow change.

## Environments on the server

| Environment | Service | Port | Data directory |
| --- | --- | --- | --- |
| main | `conorganizer-main.service` | `19080` | `/mnt/HC_Volume_103911252/environments/main/` |
| demo | `conorganizer-demo.service` | `19081` | `/mnt/HC_Volume_103911252/environments/demo/` |
| restored | `conorganizer-restored.service` | `19082` | `/mnt/HC_Volume_103911252/environments/restored/` |

Each data directory holds `database/events.db` and `event-images/`. The
directory above them, `/mnt/HC_Volume_103911252/environments`, must be owned
by `root:root` with mode 755 (`stat -c '%U:%G %a'`), or `conorganizer-deploy`
refuses to deploy a preview; `setup_deploy_user.sh` creates it that way when
the volume is mounted. The
binaries live in `/opt/conorganizer/<env>/`, next to the `goose` binary and
`migrations/` that CI deploys for `conorganizer-sqlite-migrate`, and the
services run as `deploy:www-data`. The systemd units and Caddy config are managed in
`configuration-as-code/`.

See [documentation/deployment.md](../documentation/deployment.md) for
pull request environments and CI/CD.
