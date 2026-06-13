# Billettholdere i admin

Denne sjekklisten dekker `/admin/billettholder`, der admin får oversikt over alle billettholdere og kan vedlikeholde manuelle e-postadresser.

## Roller

- Admin

## Sjekkliste

### Oversikt

- [ ] **Billettholdergrid er responsivt og lesbart**<br>
  **Gitt** at billettholderoversikten inneholder mange deltakere.<br>
  **Når** siden vises.<br>
  **Så** skal grid være responsive og kort forbli lesbare og brukbare uten sammenfallende innhold.

### E-postvedlikehold

- [ ] **Ny e-postadresse vises på riktig kort**<br>
  **Gitt** at admin legger til en manuell e-postadresse på en billettholder.<br>
  **Når** handlingen lykkes.<br>
  **Så** skal bekreftelsen vises på riktig kort og den nye adressen vises på riktig billettholder.

- [ ] **Tom e-postadresse avvises**<br>
  **Gitt** at admin forsøker å legge til en tom e-postadresse.<br>
  **Når** handlingen utføres.<br>
  **Så** skal admin få en tydelig feilmelding og ingen adresse skal legges til.

- [ ] **Duplikatadresse avvises tydelig**<br>
  **Gitt** at admin forsøker å legge til en e-postadresse som allerede finnes på samme billettholder.<br>
  **Når** handlingen utføres.<br>
  **Så** skal siden avvise duplikatet tydelig og uten å skape uklar tilstand.

- [ ] **Sletting fjerner riktig adresse**<br>
  **Gitt** at admin sletter en manuell e-postadresse.<br>
  **Når** handlingen lykkes.<br>
  **Så** skal adressen fjernes fra riktig kort og ikke bli stående igjen på siden som om den fortsatt eksisterer.

- [ ] **Brukertilknytning ryddes opp ved sletting**<br>
  **Gitt** at sletting av e-postadresse medfører at bruker-tilknytning må ryddes opp.<br>
  **Når** handlingen lykkes.<br>
  **Så** skal resultatet fremstå konsistent og ikke etterlate spor av delvis sletting i brukeropplevelsen.

- [ ] **Add- og delete-feil er tydelige**<br>
  **Gitt** at en add- eller delete-handling feiler.<br>
  **Når** admin utfører endringen.<br>
  **Så** skal feilmeldingen være tydelig og ikke etterlate inntrykk av at endringen likevel ble lagret.

### Stabilitet og layout

- [ ] **Meldinger hører til riktig kort**<br>
  **Gitt** at admin jobber med flere billettholderkort på samme side.<br>
  **Når** flere endringer skjer etter hverandre.<br>
  **Så** skal meldinger og oppdateringer tilhøre riktig kort og ikke lekke til andre kort.

- [ ] **Billettholderkort fungerer på mobil**<br>
  **Gitt** at siden brukes på mobil eller smal skjerm.<br>
  **Når** mange billettholdere eller lange e-postadresser vises.<br>
  **Så** skal innholdet fortsatt være lesbart og trykkbart uten at kortene bryter sammen.
