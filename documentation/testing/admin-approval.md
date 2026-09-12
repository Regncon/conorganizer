# Godkjenning av arrangementer

Denne sjekklisten dekker `/admin/approval` og `/admin/approval/edit/{id}`, altså oversikten over arrangementer til godkjenning og adminredigering av enkeltarrangementer.

## Roller

- Admin

## Sjekkliste

### Oversikt og redigering

- [ ] **Tom seksjon bryter ikke adminvisningen**<br>
  **Gitt** at det ikke finnes arrangementer i en av seksjonene.<br>
  **Når** siden vises.<br>
  **Så** skal resten av siden fortsatt fremstå korrekt og ikke som om hele adminvisningen feiler.

- [ ] **Riktig arrangement åpnes i redigeringsflyten**<br>
  **Gitt** at en admin åpner et arrangement fra godkjenningslisten.<br>
  **Når** redigeringssiden lastes.<br>
  **Så** skal riktig arrangement vises i skjema og forhåndsvisning.

- [ ] **Skjema og forhåndsvisning oppdateres sammen**<br>
  **Gitt** at admin redigerer felt i arrangementskjemaet fra godkjenningsflyten.<br>
  **Når** endringene lagres.<br>
  **Så** skal skjema og forhåndsvisning oppdatere seg konsistent.

- [ ] **Statusendring gir tydelig ny tilstand**<br>
  **Gitt** at admin endrer status på arrangementet.<br>
  **Når** endringen lagres.<br>
  **Så** skal statusendringen oppføre seg tydelig og ikke etterlate tvil om arrangementets nye tilstand.

- [ ] **Manglende arrangement gir forståelig feil**<br>
  **Gitt** at et arrangement ikke finnes eller ikke kan lastes.<br>
  **Når** admin forsøker å åpne det i redigeringsflyten.<br>
  **Så** skal admin møte en forståelig feiltilstand og ikke en halvferdig redigeringsvisning.

### Stabilitet og layout

- [ ] **Flere adminhandlinger holder data stabilt**<br>
  **Gitt** at flere adminhandlinger utføres etter hverandre på samme side.<br>
  **Når** siden oppdateres fortløpende.<br>
  **Så** skal innholdet forbli stabilt og ikke vise gamle eller blandede data mellom seksjonene.

- [ ] **Godkjenningsflyten er arbeidsbar på ulike skjermer**<br>
  **Gitt** at admin bruker godkjenningsflyten på større og mindre skjermer.<br>
  **Når** skjema, forhåndsvisning og interesseoversikt vises samtidig.<br>
  **Så** skal siden fortsatt være lesbar og arbeidsbar.

- [ ] **Refresh viser korrekt arrangementsdata**<br>
  **Gitt** at admin refresher siden midt i redigeringsarbeidet.<br>
  **Når** siden lastes inn igjen.<br>
  **Så** skal korrekt arrangementsdata vises.
