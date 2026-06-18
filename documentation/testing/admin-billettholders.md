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

### Stabilitet og layout

- [ ] **Add- og delete-feil er tydelige**<br>
  **Gitt** at en add- eller delete-handling feiler.<br>
  **Når** admin utfører endringen.<br>
  **Så** skal feilmeldingen være tydelig og ikke etterlate inntrykk av at endringen likevel ble lagret.

- [ ] **Meldinger hører til riktig kort**<br>
  **Gitt** at admin jobber med flere billettholderkort på samme side.<br>
  **Når** flere endringer skjer etter hverandre.<br>
  **Så** skal meldinger og oppdateringer tilhøre riktig kort og ikke lekke til andre kort.

- [ ] **Billettholderkort fungerer på mobil**<br>
  **Gitt** at siden brukes på mobil eller smal skjerm.<br>
  **Når** mange billettholdere eller lange e-postadresser vises.<br>
  **Så** skal innholdet fortsatt være lesbart og trykkbart uten at kortene bryter sammen.

### Interessemodal

- [ ] **Interesser hører til riktig kort**<br>
  **Gitt** at det er interesser og tildelte spill på en billettholder.<br>
  **Når** interessemodalen blir åpnet.<br>
  **Så** skal interesser og tildelte spill vises uten at elementer mangler eller overlapper.
