# Generelle tester

Denne sjekklisten dekker felles navigasjon, rolleopplevelse og tilgang på tvers av appen.

## Roller

- Ikke-innlogget bruker
- Innlogget bruker
- Admin

## Sjekkliste

### Alle roller

- [ ] **Hovednavigasjonen er stabil og lesbar**<br>
  **Gitt** at brukeren åpner sider med hovednavigasjon.<br>
  **Når** sidene er ferdig lastet.<br>
  **Så** skal navigasjonen være stabil, lesbar og uten brutte eller uferdige elementer.

- [ ] **Navigasjon og brukermeny passer på alle skjermstørrelser**<br>
  **Gitt** at brukeren bruker appen på mobil og større skjerm.<br>
  **Når** hovednavigasjonen og brukermenyen vises.<br>
  **Så** skal lenker, knapper og menyer være lesbare og ikke overlappe eller havne utenfor skjermen.

- [ ] **Fokus og alternative navigasjonsformer er tydelige**<br>
  **Gitt** at brukeren navigerer med tastatur eller andre alternative navigasjonsformer.<br>
  **Når** fokus flyttes i hovednavigasjon og brukermeny.<br>
  **Så** skal det være tydelig hvor fokus er og hvilke handlinger som kan utføres.

- [ ] **Rolle- og innloggingsstatus holder seg konsistent**<br>
  **Gitt** at brukeren navigerer via meny, tilbakeknapp og refresh.<br>
  **Når** rolle eller innloggingsstatus endres underveis.<br>
  **Så** skal appen ikke vise feil menytilstand eller ødelagte sider.

- [ ] **Raske sidebytter skaper ikke ustabilitet**<br>
  **Gitt** at brukeren klikker raskt mellom tilgjengelige sider.<br>
  **Når** flere navigasjonshandlinger skjer tett etter hverandre.<br>
  **Så** skal appen ikke havne i duplikathandlinger, feilnavigasjon eller tydelig ustabil tilstand.

- [ ] **Navigasjonen fremstår ferdig og konsistent**<br>
  **Gitt** at brukeren ser navigasjonen på tvers av sider.<br>
  **Når** appen brukes som helhet.<br>
  **Så** skal navigasjonen fremstå som ferdig og konsistent uten placeholder-preg, utilsiktet språkblanding eller visuelt forstyrrende detaljer.

### Ikke innlogget bruker

- [ ] **Innloggingsinngangen åpner riktig flyt**<br>
  **Gitt** at brukeren ikke er innlogget.<br>
  **Når** innloggingsinngangen brukes fra hovednavigasjonen.<br>
  **Så** skal brukeren komme til innloggingsflyten uten uventede feil eller feil side.

- [ ] **Beskyttede sider forklarer avvist tilgang**<br>
  **Gitt** at en ikke-innlogget bruker åpner en beskyttet side direkte.<br>
  **Når** tilgang avvises.<br>
  **Så** skal brukeren få en tydelig forklaring og en forståelig vei videre til innlogging eller tilbake til appen.

### Innlogget bruker

- [ ] **Utlogging endrer appen til utlogget tilstand**<br>
  **Gitt** at brukeren er innlogget.<br>
  **Når** brukeren logger ut fra brukermenyen.<br>
  **Så** skal appen tydelig gå over til utlogget tilstand.

- [ ] **Tidligere beskyttede sider forblir avvist etter utlogging**<br>
  **Gitt** at brukeren nylig har logget ut.<br>
  **Når** brukeren refresher eller åpner en tidligere beskyttet side.<br>
  **Så** skal appen fortsatt behandle brukeren som utlogget.

- [ ] **Eksterne menylenker markeres tydelig**<br>
  **Gitt** at eksterne lenker vises i brukermenyen.<br>
  **Når** brukeren åpner dem.<br>
  **Så** skal det være tydelig at innholdet ligger utenfor appen.

- [ ] **Ikke-admin får ingen halvveis adminvisning**<br>
  **Gitt** at en bruker uten adminrettigheter åpner en adminside direkte.<br>
  **Når** tilgang avvises.<br>
  **Så** skal brukeren ikke se en halvveis eller misvisende adminvisning.

### Admin

- [ ] **Adminlenken åpner adminområdet riktig**<br>
  **Gitt** at brukeren er admin.<br>
  **Når** brukeren navigerer til Admin fra hovednavigasjonen.<br>
  **Så** skal brukeren bli sendt til adminområdet uten å møte feil rolle eller feil landingsside.
