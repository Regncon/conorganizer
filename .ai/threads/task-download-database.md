--- Codex Session ---
ID: 019dfcc7-92c0-71c2-bf6e-e8c2cd8a05bb

# Download database

I want to add downloading of the database to the Taskfile.yml file. 
The database is stored on the server at /mnt/HC_Volume_103911252/main/database/events.db
The script needs to work on both linux, mac and Microslop Windows.
The developers has ssh accsess to the server. I don't want to make the database public so it needs to be behind a password.
Let's discuss how to implement this in the Taskfile.yml file.

Time: 2026-05-06 12:17:28
--- Codex Run ---

```Markdown
I’ll first look at the existing `Taskfile.yml` so the options fit the project’s current task style rather than inventing a parallel convention.
The current `download` task only mirrors public event images with `wget`; database download should probably be a separate task because it needs SSH credentials and has a different risk profile. I’m checking the repo paths and ignore rules so the local target is correct and we don’t accidentally encourage committing the DB.
The DB is already ignored via `database/*`, and tests expect `database/events.db`, so the local destination can stay simple. The main decision is how much shell logic to embed in the Taskfile versus moving cross-platform behavior into a small script.
One important implementation detail: if production SQLite uses WAL mode, copying only `events.db` can miss data sitting in `events.db-wal`. I’m checking whether the app configures SQLite journaling so the download task does not quietly produce stale or invalid local copies.
I’d implement this as a private SSH-based task, not HTTP/basic auth. SSH already gives the “behind a password/key” boundary, and the Taskfile should never contain a database password.

Recommended shape:

```yaml
tasks:
    download:
        deps:
            - download:db
            - download:images

    download:db:
        desc: Download production SQLite database over SSH
        requires:
            vars: [DB_SSH_TARGET]
        cmds:
            - |
                set -eu
                mkdir -p database
                tmp="database/events.db.download"
                rm -f "$tmp"

                ssh -p "${DB_SSH_PORT:-22}" "$DB_SSH_TARGET" 'sh -s' > "$tmp" <<'REMOTE_SCRIPT'
                set -eu
                remote_db="/mnt/HC_Volume_103911252/main/database/events.db"
                tmp="$(mktemp /tmp/conorganizer-events.XXXXXX.db)"
                trap 'rm -f "$tmp"' EXIT
                sqlite3 "$remote_db" ".timeout 5000" ".backup $tmp"
                cat "$tmp"
                REMOTE_SCRIPT

                if command -v sqlite3 >/dev/null 2>&1; then
                    sqlite3 "$tmp" "PRAGMA integrity_check;"
                fi

                mv "$tmp" database/events.db

    download:images:
        cmds:
            - mkdir -p ./local-event-images
            - wget -r -np -nH -N -R html,tmp --cut-dirs=1 -P ./local-event-images https://main.lekeplassen.regncon.no/event-images/
```

Developers would add this to their ignored `.env`:

```dotenv
DB_SSH_TARGET="username@example.com"
# DB_SSH_PORT="22"
```

Key points:

- Use `ssh`, not a public URL.
- Do not store SSH passwords in `.env`; let SSH prompt or use an SSH key/agent.
- Use SQLite’s `.backup` on the server before downloading, instead of copying the live `events.db` directly. That avoids missing WAL data.
- `database/events.db` is already ignored by `.gitignore`, so the local target is safe.
- This requires `sqlite3` on the server and OpenSSH client locally. Task’s own shell syntax is cross-platform, but external commands like `ssh` still need to exist in `PATH` per Task’s docs: https://taskfile.dev/docs/getting-started. The `requires` guard is also supported by Task: https://taskfile.dev/docs/guide.

One concern: the current image download still uses `wget`, which is not reliably available on macOS or Windows. If `task download` must be fully cross-platform, we should either keep database download as `task download:db`, or replace the image mirroring with something cross-platform too.
```

Time: 2026-05-06 12:17:28
--- Next Prompt ---
I like the idea of splitting the download task into two subtasks, one for the database and one for the images.
Could we make the script simpler by using sftp? We don't need to worry about the WAL file. Just download the `events.db` file directly.

Time: 2026-05-06 12:28:27
--- Codex Run ---

```Markdown
Yes, if we intentionally ignore WAL and just copy the DB file, the task can be much simpler. I’m checking the exact OpenSSH `sftp` invocation style so the Taskfile version stays portable and doesn’t rely on interactive-only behavior.
Yes. If we do not care about WAL, `sftp` makes this much simpler.
```

I’d split it like this:

```yaml
tasks:
    download:
        desc: Download database and event images
        deps:
            - download:db
            - download:images

    download:db:
        desc: Download production database over SFTP
        requires:
            vars: [DB_SSH_TARGET]
        cmds:
            - |
                set -eu
                mkdir -p database
                tmp="database/events.db.download"
                rm -f "$tmp"

                sftp -P "${DB_SSH_PORT:-22}" -b - "$DB_SSH_TARGET" <<SFTP
                get /mnt/HC_Volume_103911252/main/database/events.db $tmp
                SFTP

                mv "$tmp" database/events.db

    download:images:
        desc: Download event images
        cmds:
            - mkdir -p ./local-event-images
            - wget -r -np -nH -N -R html,tmp --cut-dirs=1 -P ./local-event-images https://main.lekeplassen.regncon.no/event-images/
```

Then each developer has this in their ignored `.env`:

```dotenv
DB_SSH_TARGET="username@example.com"
# DB_SSH_PORT="22"
```

I’d keep the temp file + `mv` even in the simple version. It prevents a failed transfer from replacing a working local `database/events.db` with a partial file.

One caveat: `sftp -b -` is best with SSH keys or ssh-agent. Password prompting may work in a normal terminal, but batch mode is less friendly to password auth than `scp`. If developers use SSH keys, this is the cleanest version. If they mostly type SSH passwords manually, `scp` may be simpler operationally.
```

Time: 2026-05-06 12:28:27
--- Next Prompt ---
Ok, let'g go with scp. Where does $tmp come from?  The host ip is fixt and it's 46.62.159.49
If we use scp do we still need the tmp file?

Time: 2026-05-06 12:44:07
--- Codex Run ---

```Markdown
`$tmp` was just a local shell variable I introduced inside the script. We can avoid that entirely and use the temp filename directly.

With fixed host IP and `scp`, I’d use:

```yaml
download:db:
    desc: Download production database over SCP
    cmds:
        - |
            set -eu
            : "${DB_SSH_USER:?Set DB_SSH_USER in .env}"
            mkdir -p database
            scp -P "${DB_SSH_PORT:-22}" "${DB_SSH_USER}@46.62.159.49:/mnt/HC_Volume_103911252/main/database/events.db" "database/events.db.download"
            mv "database/events.db.download" "database/events.db"
```

`.env`:

```dotenv
DB_SSH_USER="your-server-username"
# DB_SSH_PORT="22"
```

We do not strictly need the temp file with `scp`; this would also work:

```bash
scp ... database/events.db
```

But I’d keep `database/events.db.download`. If the transfer fails halfway, the existing local `database/events.db` stays intact. Only after a successful `scp` do we replace it with `mv`.
```

Time: 2026-05-06 12:44:07
--- Next Prompt ---
Ok, we can hard code the port to 22 as well. Let'r keep things simple.


Time: 2026-05-06 12:49:37
--- Codex Run ---

```Markdown
I’ll make the Taskfile change now: split `download` into `download:db` and `download:images`, keep the database copy as a plain `scp` over port 22, and keep the temp file so interrupted transfers do not clobber the current local DB.
I spotted that the first patch left the old task name in place, which would duplicate `download`. I’m correcting that to `download:images` and then I’ll re-read the relevant section.
The Taskfile section now has the intended shape: `download` depends on `download:db` and `download:images`, and the DB copy uses the fixed IP and fixed port. I’m doing a lightweight YAML parse check since the local `task` binary is not installed in this environment.
One small docs tweak is needed because `task download` now requires `DB_SSH_USER`. I’ll add that near the existing database download note so new developers know what to put in `.env`.
Implemented.

`download` now fans out to `download:db` and `download:images` in [Taskfile.yml](/home/cinmay/Documents/conorganizer/Taskfile.yml:80). The DB task uses:

```
```bash
scp -P 22 "${DB_SSH_USER}@46.62.159.49:/mnt/HC_Volume_103911252/main/database/events.db" "database/events.db.download"
mv "database/events.db.download" "database/events.db"
```

So developers only need this in `.env`:

```dotenv
DB_SSH_USER="their-server-username"
```

I also added a short README note at [README.md](/home/cinmay/Documents/conorganizer/README.md:54).

Verification: `Taskfile.yml` parses as YAML. I could not run `task` itself because the `task` binary is not installed in this environment. Note that `download:db` is the cross-platform DB task assuming OpenSSH is installed; `download:images` still depends on `wget`.
```

Time: 2026-05-06 12:49:37
--- Next Prompt ---
Task is not installed as a separate binary, but it is available as a Go tool. go tool task should work.  

