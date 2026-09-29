# Godkjenning av arrangementer

Denne sjekklisten dekker `/admin/approval` og `/admin/approval/edit/{id}`, altså oversikten over alle arrangementer gruppert etter status og adminredigering av enkeltarrangementer med status, puljer og rom. Spillertildeling gjøres i puljefordelingen og ikke på redigeringssiden.

## Roller

- Admin

## Sjekkliste

### Oversikt og redigering

- [ ] **Arrangementer grupperes etter status i fast rekkefølge**<br>
  **Gitt** at det finnes arrangementer med ulike statuser.<br>
  **Når** admin åpner godkjenningssiden.<br>
  **Så** skal alle arrangementer vises, hvert under sin status i rekkefølgen Innsendt, Godkjent, Annonsert, Kladd, Forkastet, og hvert arrangement skal vise tittel og system.

- [ ] **Tom seksjon bryter ikke adminvisningen**<br>
  **Gitt** at det ikke finnes arrangementer med en av statusene.<br>
  **Når** siden vises.<br>
  **Så** skal seksjonen for den statusen ikke vises, og resten av siden skal fortsatt fremstå korrekt.

- [ ] **Riktig arrangement åpnes i redigeringsflyten**<br>
  **Gitt** at en admin åpner et arrangement fra godkjenningslisten.<br>
  **Når** redigeringssiden lastes.<br>
  **Så** skal riktig arrangement vises i skjema, forhåndsvisning og puljefordelingskortet.

- [ ] **Skjema og forhåndsvisning oppdateres sammen**<br>
  **Gitt** at admin redigerer felt i arrangementskjemaet fra godkjenningsflyten.<br>
  **Når** endringene lagres.<br>
  **Så** skal skjema og forhåndsvisning oppdatere seg konsistent.

- [ ] **Statusendring gir tydelig ny tilstand**<br>
  **Gitt** at admin endrer status på arrangementet til Kladd, Innsendt, Godkjent, Annonsert eller Forkastet.<br>
  **Når** endringen lagres.<br>
  **Så** skal statusendringen oppføre seg tydelig, og arrangementet skal vises under den nye statusen på godkjenningssiden.

- [ ] **Manglende arrangement gir forståelig feil**<br>
  **Gitt** at et arrangement ikke finnes eller ikke kan lastes.<br>
  **Når** admin forsøker å åpne det i redigeringsflyten.<br>
  **Så** skal admin møte en forståelig feiltilstand og ikke en halvferdig redigeringsvisning.

### Puljer og rom

- [ ] **Med i puljefordeling lagres**<br>
  **Gitt** at admin endrer «Med i puljefordeling» på et arrangement.<br>
  **Når** endringen lagres.<br>
  **Så** skal valget være det samme etter refresh og ikke påvirke andre arrangementer.

- [ ] **Puljevalg lagres på riktig pulje**<br>
  **Gitt** at admin krysser av eller fjerner «Legg til i pulje» for én pulje.<br>
  **Når** endringen lagres.<br>
  **Så** skal bare den puljen endres, og romvalget for puljen skal bare være mulig når arrangementet er med i puljen.

- [ ] **Romvalg lagres per pulje**<br>
  **Gitt** at arrangementet er med i flere puljer.<br>
  **Når** admin velger rom under «Legg til rom» i én pulje.<br>
  **Så** skal rommet bare lagres for den puljen og vises i romfordelingen for samme pulje.

- [ ] **Feil ved puljer og rom er tydelig**<br>
  **Gitt** at lagring av puljevalg eller romvalg feiler.<br>
  **Når** admin gjør endringen.<br>
  **Så** skal feilen være tydelig nok til at admin forstår at endringen ikke ble fullført.

### Stabilitet og layout

- [ ] **Flere adminhandlinger holder data stabilt**<br>
  **Gitt** at flere adminhandlinger utføres etter hverandre på samme side.<br>
  **Når** siden oppdateres fortløpende.<br>
  **Så** skal innholdet forbli stabilt og ikke vise gamle eller blandede data mellom seksjonene.

- [ ] **Godkjenningsflyten er brukbar på ulike skjermer**<br>
  **Gitt** at admin bruker godkjenningsflyten på større og mindre skjermer.<br>
  **Når** skjema, puljefordelingskort og forhåndsvisning vises samtidig.<br>
  **Så** skal siden fortsatt være lesbar og brukbar.

- [ ] **Refresh viser korrekt data, puljer og rom**<br>
  **Gitt** at admin refresher siden midt i redigeringsarbeidet.<br>
  **Når** siden lastes inn igjen.<br>
  **Så** skal korrekt arrangementsdata, status, puljevalg og romvalg vises.
