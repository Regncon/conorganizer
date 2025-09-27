# Important for Linux users
This folder is used when binding the volume in docker compose and for inheriting permission

If you're encountering issues on Linux you can try to expose your user id and group in `~/.bashrc`
```bash
export GID="$(id -g)"
export GNAME="$(id -gn)"
```

You can then apply this in `docker-compose.yaml`
```yaml
services:
    backup-service:
        user: "${UID}:${GID}"
        build:
            context: .
            dockerfile: Dockerfile.local
        container_name: backup-dev
```
