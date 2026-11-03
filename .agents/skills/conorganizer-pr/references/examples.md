# PR examples

The maintainers confirmed these picks on 2026-11-03 from the 60 merged PRs before #702. The preview block CI appends is left out of every body below.

## Model: three to copy

### #653 Fix save blocked by manual seat on event removed from pulje

> Reported by an admin: when an event is taken out of a pulje after someone was manually placed on it, "Lagre fordeling" kept showing "1 ulagret endring" that could never be saved away.
>
> Cause: saving only replaced solver seats, so the manual seat on the removed event survived every save and kept differing from the emulated distribution.
>
> Fix: saving the distribution now also removes player seats on events no longer in the pulje. Manual seats on events still in the pulje and GM rows are untouched.
>
> Covered by `TestSaveConfirmedDistribution_ClearsManualSeatOnEventNoLongerInPulje`.

90 words. The title names the bug as the admin met it. The body is symptom, cause, fix, test, with no heading, because one change needs none.

### #695 Stop serving meetup-january-2026.lekeplassen.regncon.no

> ## Why
>
> `meetup-january-2026.lekeplassen.regncon.no` proxied to `127.0.0.1:8080`, where nothing listens, so Caddy answered every request with a 502 (about 46 a day from bots) and kept a certificate renewed for a site that is no longer in use. The app behind it (`meetupJan2026-main.service`, `/opt/meetupJan2026`) was disabled and dead, and has been removed from the server on 2026-10-22 with a backup in `/root`.
>
> ## What
>
> - Removes the site block from the checked-in Caddyfile and its row in the Caddy hostnames table.
> - `hetzner/scripts/addUser.sh` no longer links `/opt/meetupJan2026` into new admins' home directories.
>
> ## On the server, after merging
>
> The Caddyfile is a stow link into `/srv/configuration-as-code-repo/conorganizer`, so moving that checkout to the merge commit changes the file in place; no restow is needed, only a reload:
>
> ```bash
> sudo chown -R "$USER": /srv/configuration-as-code-repo/conorganizer
> git -C /srv/configuration-as-code-repo/conorganizer fetch origin
> git -C /srv/configuration-as-code-repo/conorganizer log --oneline HEAD..origin/main -- configuration-as-code
> git -C /srv/configuration-as-code-repo/conorganizer checkout --detach origin/main
> sudo caddy validate --config /etc/caddy/Caddyfile
> sudo systemctl reload caddy
> ```
>
> The `log` line shows every other configuration-as-code change the checkout picks up at the same time; apply their README steps if any are still pending.

150 words around a six-line command block. The rollout section is the only long part, and it is the part an admin needs. Under this skill the headings would be the summary without a heading, `## What changed` and `## After merging`; the shape is the same.

### #611 Add a publish switch for each pulje's room assignment

> Stacked on #605.
>
> - Migration: adds `puljer.rooms_published`, off by default.
> - Romfordeling tab: a "Publiser romfordeling" switch per pulje ("Synlig for alle" / "Ikke publisert").
> - Event page and print page: until a pulje is published, everyone (admins included) sees only the time, with no room, map or notes. Admins review the assignment in the Romfordeling tab.
> - Rooms are hidden on the server, so unpublished rooms never reach the HTML.
>
> Requires running the migration (goose up). Room assignments are hidden until each pulje is published.

80 words. One bullet per parallel change, user-facing strings quoted as they appear, and the migration line. Under this skill the migration line names the file: "Run `20260928120000_rooms_published.sql` on main and demo by hand (documentation/migrations.md)."

## Three to avoid

### #650 move kontoistillinger down and some small refactoring

Empty body. "Some small refactoring" names nothing, so main's history says nothing about what moved or why. A title that works: "Move Kontoadministrasjon to the bottom of Min side", with a two-sentence body.

### #599 583 bug kartene som vises i kart modalene er i light mode

Empty body. The branch slug as the title: issue number first, bokmål, no verb, and 24 files changed that nobody can account for from the log. A title that works: "Show the room maps in dark mode in the map dialogs", with "Closes #583" as the first body line.

### #671 Show backup age and SQLite sizes in Grafana from backup script metrics

The title names the outcome but runs to 70 characters, five over the limit; "Show backup age and SQLite sizes in Grafana" (43) says the same. The body ran to about 850 words: `## Summary`, `## What changed` with three `###` subsections, a six-row metrics table, `## Testing` with five bullets, `## After merge, on the server`, and `## Known limits`, plus an account of the review process ("An adversarial review (3 reviewers, each finding checked by another agent) led to the `du` and WAL fixes above."). A reader needs the metrics, the dashboard change, the tests and the four server steps. The rest was already in the README the PR updated.

The same PR in this skill's shape, 170 words around the steps:

> Several panels in the production dashboards showed "No data" because nothing exported the `conorganizer_*` metrics they query. The two backup scripts now publish them through node_exporter's textfile collector, and the backup age panels from #670 read the last successful run instead of the last timer trigger.
>
> ## What changed
>
> - `conorganizer-sqlite-backup` writes its last-success time, the newest backup's size, the `events.db` and WAL sizes and the backup directory size to `/var/lib/prometheus/node-exporter/` after each run; `conorganizer-images-backup` writes its last-success time. A failed run writes nothing, so the age panels grow.
> - Production Health: the two backup age panels read the new timestamps, and a "Backup directory size" panel joins them. Main Service Debugging drops "Optional" from its SQLite size panels.
> - The dashboards README describes the metrics; `configuration-as-code/README.md` and `documentation/deployment.md` get a paragraph each.
>
> ## Tests
>
> Both scripts ran end to end in a sandbox with a stub `zstd`, with and without a WAL, and the files they wrote parse with node_exporter's text parser. The panel queries run against production Prometheus and return no data until the scripts are deployed.
>
> ## After merging, on the server
>
> 1. Move the checkout to `origin/main` ("Repository location" in `configuration-as-code/README.md`), then `sudo stow --dir=configuration-as-code/stow --target=/ --restow scripts`.
> 2. `sudo systemctl start conorganizer-sqlite-backup.service conorganizer-images-backup.service`, or wait for the timers.
> 3. `curl -fsS http://127.0.0.1:9100/metrics | grep -E '^(conorganizer_|node_textfile_scrape_error)'` shows six `conorganizer_` lines and `node_textfile_scrape_error 0`.
> 4. Re-import both dashboard JSON files in Grafana.
