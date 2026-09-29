# Manuelle tester

Denne mappen inneholder launch-sjekklistene for manuell testing av Conorganizer.

## Testfiler

- [ ] [Generelle tester](./general.md)
- [ ] [Forside](./root.md)
- [ ] [Min Side](./profile.md)
- [ ] [Billetter på Min Side](./profile-tickets.md)
- [ ] [Arrangementsskjema](./event-form.md)
- [ ] [Arrangementsdetaljer](./event-details.md)
- [ ] [Admin](./admin.md)
- [ ] [Godkjenning av arrangementer](./admin-approval.md)
- [ ] [Billettholdere i admin](./admin-billettholders.md)
- [ ] [Legg til billettholder i admin](./admin-add-billettholder.md)
- [ ] [Romadministrasjon i admin](./admin-rooms.md)

## Guide

- [Hvordan vi skriver manuelle tester](./how-to-write-tests.md)

## Automatisert testoversikt

- Se [Automatiserte tester](../automated-tests.md) for hvordan Go-testene skrives, og hvordan BDD-rapporten lages med `go tool task test:report` lokalt og i CI.

## Relatert dokumentasjon

Disse dokumentene beskriver forventet oppførsel som sjekklistene tester:

- [Tilgangskontroll og feilsider](../access-control-and-error-pages.md)
- [Puljestatus og publisering](../pulje-status-and-publishing.md)
- [Billettholdere](../billettholdere.md)
- [Romfordeling](../room-assignment.md)
- [UI-konvensjoner](../ui-conventions.md)

## Bygge PDF-er

- Kjør `./build-pdfs.sh` for å lage én PDF per `*.md` i denne mappen. Skriptet bytter selv til mappen sin, så det kan kjøres fra hvor som helst. Det krever `pandoc` og `xelatex`.
- Pandoc dropper rå `<br>` når det lager LaTeX, og da ville Gitt/Når/Så slått seg sammen til ett avsnitt. Skriptet bruker derfor et midlertidig Lua-filter som gjør `<br>` om til ekte linjeskift.
- PDF-er uten tilhørende `.md`-fil, som `auth.pdf`, `report.pdf` og `PR-558-manual-test-checklist.pdf`, er rester og lages ikke av skriptet.

## Dekningsinventar

Launch-sjekklistene dekker disse aktive sidene og flytene:

- `/` og oppdatert forsidestruktur fra `/root/api/` dekkes av [Forside](./root.md).
- `/auth`, `/auth/post-login` og `/auth/logout` dekkes av seksjonene Authentisering og Autorisering i [Generelle tester](./general.md).
- `/profile` dekkes av [Min Side](./profile.md).
- `/profile/descope-profile` (Descope-profilwidgeten bak «Reset passord» på Min Side) dekkes av seksjonen Authentisering i [Generelle tester](./general.md).
- `/profile/tickets` dekkes av [Billetter på Min Side](./profile-tickets.md).
- `/profile/new/{id}` og tilhørende skjema- og bildeopplastingsflyt dekkes av [Arrangementsskjema](./event-form.md).
- `/event/{id}` og interesseflyten under `/event/api/{id}` dekkes av [Arrangementsdetaljer](./event-details.md).
- `/admin` (Adminverktøy) dekkes av [Admin](./admin.md).
- `/admin/approval` og `/admin/approval/edit/{id}` dekkes av [Godkjenning av arrangementer](./admin-approval.md).
- `/admin/billettholder` dekkes av [Billettholdere i admin](./admin-billettholders.md).
- `/admin/billettholder/add` dekkes av [Legg til billettholder i admin](./admin-add-billettholder.md).
- `/admin/rooms` og `/admin/rooms/assignment/{pulje}` dekkes av [Romadministrasjon i admin](./admin-rooms.md).

Disse aktive adminsidene lenkes fra adminforsiden, men har ennå ingen egen launch-sjekkliste:

- `/admin/puljefordeling/` og `/admin/puljefordeling/{pulje}`.
- `/admin/puljeoppsett/`.
- `/admin/tilbakemeldinger/` og QR-kodesiden `/admin/tilbakemeldinger/qr`.

Tilbakemeldinger har ennå ingen manuell sjekkliste. Det gjelder både skjemaet `/tilbakemelding`, som innloggede brukere når fra brukermenyen og Min Side, og adminlisten `/admin/tilbakemeldinger/` med filtrering på kategori, sletting etter bekreftelse og utskrift av QR-kode. Oppførselen dekkes foreløpig av Go-testene i `pages/feedback/`, `service/feedback/` og `pages/admin/feedback_admin*_test.go`. Lenkene fra brukermenyen og Min Side sjekkes i [Generelle tester](./general.md) og [Min Side](./profile.md).

Disse rutene er bevisst ikke egne launch-sjekklister:

- `/print`, fordi printvennlig side ikke er del av den vanlige launch-reisen.
- `/auth/test`, fordi dette er en diagnostisk rute og ikke en brukerflyt.
- Rene API- og liveoppdateringsruter, fordi de testes gjennom siden eller flyten som eier oppførselen.

## Bruk

- Start med [Generelle tester](./general.md) og [Forside](./root.md) for å verifisere grunnleggende navigasjon og synlig innhold.
- Kjør deretter rollebaserte og funksjonelle tester i de relevante filene.
- Bruk `go tool task test:report` for å sammenligne manuelle sjekkpunkter med automatiserte tester, se [Automatiserte tester](../automated-tests.md).
