# Åpen påmelding

Denne sjekklisten dekker åpen påmelding fra `/event/{id}`, tilhørende visning på `/profile` og administrativ konfigurasjon. Databaseregler og kildehåndtering dekkes av automatiserte tester; punktene her prioriterer den sammenhengende brukeropplevelsen og liveoppdateringene.

## Roller

- Innlogget bruker med billettholder
- Admin

## Sjekkliste

### Konfigurasjon

- [ ] **Bare admin ser innstillingen for åpen påmelding**<br>
  **Gitt** at et arrangement redigeres av henholdsvis en vanlig bruker og en admin.<br>
  **Når** skjemaet vises.<br>
  **Så** skal bare admin kunne se og endre innstillingen for åpen påmelding.

- [ ] **Lagret innstilling endrer interessevelgeren**<br>
  **Gitt** at admin har markert et arrangement for åpen påmelding.<br>
  **Når** en billettholder åpner arrangementet etter at programmet er publisert.<br>
  **Så** skal `Meld deg på` erstatte `Veldig interessert`, mens de lavere interessenivåene fortsatt vises.

### På- og avmelding

- [ ] **Konsekvensen forklares før påmelding**<br>
  **Gitt** at billettholderen åpner interessevelgeren for et arrangement med åpen påmelding.<br>
  **Når** `Meld deg på` vises.<br>
  **Så** skal det forklares at plassen bekreftes direkte, at ordinær fordeling i puljen stopper, og at andre åpne påmeldinger fortsatt er mulige.

- [ ] **Direkte påmelding bekrefter plassen**<br>
  **Gitt** at programmet er publisert, puljen er åpen og arrangementet har åpen påmelding.<br>
  **Når** billettholderen velger `Meld deg på`.<br>
  **Så** skal valget endres til `Meld deg av`, og arrangementet skal vises på profilen uten manuell oppdatering av siden.

- [ ] **Avmelding fjerner bare den valgte plassen**<br>
  **Gitt** at billettholderen er påmeldt et arrangement og har andre interesser eller påmeldinger.<br>
  **Når** billettholderen velger `Meld deg av`.<br>
  **Så** skal bare plassen og eventuell interesse for dette arrangementet i denne puljen fjernes.

- [ ] **Manuell spillertildeling kan meldes av**<br>
  **Gitt** at admin har tildelt billettholderen som spiller på et vanlig arrangement.<br>
  **Når** billettholderen åpner arrangementet mens programmet er publisert og puljen er åpen.<br>
  **Så** skal `Meld deg av` vises og kunne fjerne den manuelle tildelingen.

- [ ] **Flere åpne arrangementer kan velges i samme pulje**<br>
  **Gitt** at minst to arrangementer i samme pulje har åpen påmelding.<br>
  **Når** billettholderen melder seg på begge.<br>
  **Så** skal begge påmeldingene beholdes og vises i advarselen og på profilen.

- [ ] **Påmelding gjelder bare valgt pulje**<br>
  **Gitt** at samme arrangement tilbys i flere puljer.<br>
  **Når** billettholderen melder seg på eller av i én pulje.<br>
  **Så** skal påmelding, avmelding og blokkering av interesser bare påvirke den valgte puljen.

### Interesser og blokkering

- [ ] **Lavere interesser fungerer som reservevalg**<br>
  **Gitt** at billettholderen ikke er direkte påmeldt i puljen.<br>
  **Når** billettholderen velger `Interessert` eller `Litt interessert` på et arrangement med åpen påmelding.<br>
  **Så** skal interessen lagres som et vanlig reservevalg.

- [ ] **Bekreftet plass blokkerer vanlige interesser**<br>
  **Gitt** at billettholderen er manuelt tildelt eller påmeldt et arrangement i puljen.<br>
  **Når** interessevelgeren åpnes for et annet arrangement i samme pulje.<br>
  **Så** skal vanlige interessevalg være deaktivert, og advarselen skal lenke til alle bekreftede arrangementer i puljen.

- [ ] **Bekreftet plass blokkerer ikke flere direkte påmeldinger**<br>
  **Gitt** at billettholderen allerede har en bekreftet plass i puljen.<br>
  **Når** et annet arrangement med åpen påmelding åpnes.<br>
  **Så** skal billettholderen fortsatt kunne velge `Meld deg på` for dette arrangementet.

- [ ] **GM kan ikke gjøre endringer i puljen**<br>
  **Gitt** at billettholderen er tildelt som GM i puljen.<br>
  **Når** interessevelgeren åpnes for et arrangement i samme pulje.<br>
  **Så** skal påmelding, avmelding og vanlige interesser være deaktivert, og advarselen skal lenke til GM-arrangementet.

- [ ] **Aldersgrense blokkerer påmelding og interesse**<br>
  **Gitt** at billettholderen er under 18 år og arrangementet har 18-årsgrense.<br>
  **Når** interessevelgeren åpnes.<br>
  **Så** skal påmelding og vanlige interesser være deaktivert med en tydelig forklaring.

### Publisering og puljestatus

- [ ] **Upublisert program åpner ikke påmelding**<br>
  **Gitt** at arrangementet har åpen påmelding, men programmet ikke er publisert.<br>
  **Når** billettholderen åpner arrangementet.<br>
  **Så** skal interesse- og påmeldingsvalgene ikke være tilgjengelige.

- [ ] **Låst pulje fryser alle valg**<br>
  **Gitt** at billettholderen har interesser eller påmeldinger i en pulje som er `Locked`.<br>
  **Når** interessevelgeren åpnes.<br>
  **Så** skal påmelding, avmelding og interesseendringer være deaktivert.

### Profil

- [ ] **Bekreftede tildelinger vises før fullføring**<br>
  **Gitt** at programmet er publisert og billettholderen har en manuell spillertildeling, påmelding eller GM-tildeling i en åpen eller låst pulje.<br>
  **Når** profilen vises.<br>
  **Så** skal den bekreftede tildelingen vises, og interessenivåene for puljen skal være skjult.

- [ ] **Solverresultatet skjules frem til fullføring**<br>
  **Gitt** at solveren har opprettet en spillertildeling i en åpen eller låst pulje.<br>
  **Når** profilen vises før og etter at puljen settes til `Completed`.<br>
  **Så** skal resultatet være skjult før fullføring og vises etterpå.

### Robusthet

- [ ] **Raske eller gjentatte valg lager ikke duplikater**<br>
  **Gitt** at billettholderen kan melde seg på et arrangement.<br>
  **Når** samme handling aktiveres flere ganger tett etter hverandre.<br>
  **Så** skal bare én påmelding finnes, og interessevelgeren skal ende i en forståelig tilstand.

- [ ] **Liveoppdatering bevarer valgt billettholder og pulje**<br>
  **Gitt** at brukeren har valgt billettholder og pulje i interessevelgeren.<br>
  **Når** en påmelding eller avmelding lagres og siden oppdateres automatisk.<br>
  **Så** skal riktig billettholder, pulje, advarsel og knappetilstand fortsatt vises.
