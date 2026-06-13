# Manuelle tester

Denne mappen inneholder launch-sjekklistene for manuell testing av Conorganizer.

## Testfiler

- [ ] [Generelle tester](./general.md)
- [ ] [Forside](./root.md)
- [ ] [Autentisering](./auth.md)
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

- Kjør `task test:report` lokalt for å se hvilke Go-tester som kjøres og hvilken BDD-kommentar hver test dekker.
- GitHub Actions skriver samme rapport til CI-loggen. Rapporten lagres ikke som artefakt og skal ikke committes.

## Dekningsinventar

Launch-sjekklistene dekker disse aktive sidene og flytene:

- `/` og oppdatert forsidestruktur fra `/root/api/` dekkes av [Forside](./root.md).
- `/auth`, `/auth/post-login` og `/auth/logout` dekkes av [Autentisering](./auth.md).
- `/profile` dekkes av [Min Side](./profile.md).
- `/profile/tickets` dekkes av [Billetter på Min Side](./profile-tickets.md).
- `/profile/new/{id}` og tilhørende skjema- og bildeopplastingsflyt dekkes av [Arrangementsskjema](./event-form.md).
- `/event/{id}` og interesseflyten under `/event/api/{id}` dekkes av [Arrangementsdetaljer](./event-details.md).
- `/admin` dekkes av [Admin](./admin.md).
- `/admin/approval` og `/admin/approval/edit/{id}` dekkes av [Godkjenning av arrangementer](./admin-approval.md).
- `/admin/billettholder` dekkes av [Billettholdere i admin](./admin-billettholders.md).
- `/admin/billettholder/add` dekkes av [Legg til billettholder i admin](./admin-add-billettholder.md).
- `/admin/rooms` og `/admin/rooms/assignment/{pulje}` dekkes av [Romadministrasjon i admin](./admin-rooms.md).

Disse rutene er bevisst ikke egne launch-sjekklister:

- `/print`, fordi printvennlig side ikke er del av den vanlige launch-reisen.
- `/auth/test`, fordi dette er en diagnostisk rute og ikke en brukerflyt.
- Rene API- og liveoppdateringsruter, fordi de testes gjennom siden eller flyten som eier oppførselen.

## Bruk

- Start med [Generelle tester](./general.md) og [Forside](./root.md) for å verifisere grunnleggende navigasjon og synlig innhold.
- Kjør deretter rollebaserte og funksjonelle tester i de relevante filene.
- Bruk `task test:report` for å sammenligne manuelle sjekkpunkter med automatiserte tester.
