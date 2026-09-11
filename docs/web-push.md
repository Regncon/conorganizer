# Web Push-varsler

Web Push er avslått når variablene under ikke er satt. Alle tre må settes for å slå på tjenesten:

```text
WEB_PUSH_PUBLIC_KEY=<offentlig VAPID-nøkkel>
WEB_PUSH_PRIVATE_KEY=<privat VAPID-nøkkel>
WEB_PUSH_SUBJECT=mailto:drift@example.com
```

Lag et nytt nøkkelpar lokalt med:

```sh
go run ./cmd/varslerkeys
```

Lagre den private nøkkelen som en hemmelighet i driftsmiljøet. Ikke legg nøklene i Git, logger eller dokumentasjon. `WEB_PUSH_SUBJECT` må være en `mailto:`-adresse eller en HTTPS-URL som identifiserer avsenderen.

Databasen må migreres med Goose-migrasjonen `20260910180000_add_web_push_varsler.sql` før denne versjonen av applikasjonen tas i bruk, også når Web Push er avslått. Publisering lagrer resultathistorikk i de nye tabellene uavhengig av push-konfigurasjonen. Nettleseren kan bare registrere HTTPS-endepunkter hos de støttede push-leverandørene. Et endepunkt som allerede tilhører en annen konto blir avvist; det flyttes aldri automatisk mellom kontoer.

Arbeideren bruker en varig kø med tidsbegrenset lease og avgrenset retry. HTTP 404 og 410 fjerner abonnementet. Levering kan skje mer enn én gang dersom push-tjenesten mottar meldingen, men forbindelsen feiler før svaret når applikasjonen; den stabile `tag`-verdien lar klienten erstatte samme varsel.
