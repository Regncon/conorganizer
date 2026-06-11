# Generelle tester

Denne sjekklisten dekker felles navigasjon, rolleopplevelse og tilgang på tvers av appen.

## Roller

- Ikke-innlogget bruker
- Innlogget bruker
- Admin

## Sjekkliste

### Alle roller

- [ ] `Gitt at brukeren åpner sider med hovednavigasjon, når sidene er ferdig lastet, så skal navigasjonen være stabil, lesbar og uten brutte eller uferdige elementer.`

- [ ] `Gitt at brukeren bruker appen på mobil og større skjerm, når hovednavigasjonen og brukermenyen vises, så skal lenker, knapper og menyer være lesbare og ikke overlappe eller havne utenfor skjermen.`

- [ ] `Gitt at brukeren navigerer med tastatur eller andre alternative navigasjonsformer, når fokus flyttes i hovednavigasjon og brukermeny, så skal det være tydelig hvor fokus er og hvilke handlinger som kan utføres.`

- [ ] `Gitt at brukeren navigerer via meny, tilbakeknapp og refresh, når rolle eller innloggingsstatus endres underveis, så skal appen ikke vise feil menytilstand eller ødelagte sider.`

- [ ] `Gitt at brukeren klikker raskt mellom tilgjengelige sider, når flere navigasjonshandlinger skjer tett etter hverandre, så skal appen ikke havne i duplikathandlinger, feilnavigasjon eller tydelig ustabil tilstand.`

- [ ] `Gitt at brukeren ser navigasjonen på tvers av sider, når appen brukes som helhet, så skal navigasjonen fremstå som ferdig og konsistent uten placeholder-preg, utilsiktet språkblanding eller visuelt forstyrrende detaljer.`

### Ikke innlogget bruker

- [ ] `Gitt at brukeren ikke er innlogget, når innloggingsinngangen brukes fra hovednavigasjonen, så skal brukeren komme til innloggingsflyten uten uventede feil eller feil side.`

- [ ] `Gitt at en ikke-innlogget bruker åpner en beskyttet side direkte, når tilgang avvises, så skal brukeren få en tydelig forklaring og en forståelig vei videre til innlogging eller tilbake til appen.`

### Innlogget bruker

- [ ] `Gitt at brukeren er innlogget, når brukeren logger ut fra brukermenyen, så skal appen tydelig gå over til utlogget tilstand.`

- [ ] `Gitt at brukeren nylig har logget ut, når brukeren refresher eller åpner en tidligere beskyttet side, så skal appen fortsatt behandle brukeren som utlogget.`

- [ ] `Gitt at eksterne lenker vises i brukermenyen, når brukeren åpner dem, så skal det være tydelig at innholdet ligger utenfor appen.`

- [ ] `Gitt at en bruker uten adminrettigheter åpner en adminside direkte, når tilgang avvises, så skal brukeren ikke se en halvveis eller misvisende adminvisning.`

### Admin

- [ ] `Gitt at brukeren er admin, når brukeren navigerer til Admin fra hovednavigasjonen, så skal brukeren bli sendt til adminområdet uten å møte feil rolle eller feil landingsside.`
