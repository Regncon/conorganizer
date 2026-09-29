# Billettholdere i admin

Denne sjekklisten dekker `/admin/billettholder`, der admin får oversikt over alle billettholdere, kan filtrere på førstevalg og spilleder, se interesser og tildelinger per pulje og vedlikeholde manuelle e-postadresser.

## Roller

- Admin

## Sjekkliste

### Oversikt

- [ ] **Billettholdergrid er responsivt og lesbart**<br>
  **Gitt** at billettholderoversikten inneholder mange deltakere.<br>
  **Når** siden vises.<br>
  **Så** skal grid være responsive og kort forbli lesbare og brukbare uten sammenfallende innhold.

- [ ] **Tellere og merker viser førstevalg og spilleder**<br>
  **Gitt** at noen billettholdere er tildelt som spiller på et arrangement de har gitt `Veldig interessert`, og noen er tildelt som spilleder.<br>
  **Når** oversikten vises.<br>
  **Så** skal tellerne Totalt, Uten førstevalg og GM/DM stemme, og hvert kort skal vise «Har fått førstevalg» eller «Har ikke fått førstevalg» og «Spilleder (GM/DM)» eller «Ikke spilleder».

- [ ] **Filtrene begrenser listen**<br>
  **Gitt** at admin slår på filteret «Uten førstevalg», «GM/DM» eller begge.<br>
  **Når** listen oppdateres.<br>
  **Så** skal bare billettholdere som oppfyller alle valgte filtre vises, knappene skal vises som valgt, og hele listen skal komme tilbake når filtrene slås av.

- [ ] **Legg til billettholder åpner billettsiden**<br>
  **Gitt** at admin er på billettholderoversikten.<br>
  **Når** admin velger «Legg til billettholder».<br>
  **Så** skal siden for å konvertere billetter fra CheckIn åpnes.

### E-postadresser

- [ ] **Manuell e-post legges til riktig billettholder**<br>
  **Gitt** at admin skriver en ny e-postadresse på ett av flere billettholderkort.<br>
  **Når** admin trykker Enter eller «Legg til epost».<br>
  **Så** skal adressen vises under «Andre epostadresser på samme bestilling» på riktig kort, feltet tømmes og en bekreftelse vises på samme kort.

- [ ] **Tom eller duplisert e-post avvises**<br>
  **Gitt** at admin forsøker å legge til en tom adresse eller en adresse som allerede finnes på billettholderen.<br>
  **Når** handlingen utføres.<br>
  **Så** skal kortet vise en tydelig feilmelding og ingen ny adresse lagres.

- [ ] **Bare manuelle e-postadresser kan slettes**<br>
  **Gitt** at en billettholder har billett-e-post, e-post fra samme bestilling og manuelt lagt til e-post.<br>
  **Når** kortet vises og admin sletter den manuelle adressen.<br>
  **Så** skal bare manuelle adresser ha «Slett», og etter sletting skal adressen være borte fra riktig kort med en bekreftelse.

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

### Interessedialog

- [ ] **Interesser hører til riktig kort**<br>
  **Gitt** at det er interesser og tildelte spill på en billettholder.<br>
  **Når** admin åpner dialogen med knappen «Interesser (antall) Tildelt (antall)».<br>
  **Så** skal dialogen gjelde riktig billettholder, antallene stemme, og interesser og tildelte spill vises uten at elementer mangler eller overlapper.

- [ ] **Puljefaner viser tildelt først**<br>
  **Gitt** at billettholderen har interesser eller tildelinger i flere puljer.<br>
  **Når** admin bytter mellom puljefanene.<br>
  **Så** skal hver fane vise antall interesser og tildelinger og et spilledermerke når billettholderen er spilleder i puljen, og panelet skal vise tildelte arrangementer først og deretter interesser gruppert etter interessenivå.

- [ ] **Billettholder uten interesser får tom tilstand**<br>
  **Gitt** at billettholderen verken har interesser eller tildelinger.<br>
  **Når** dialogen åpnes.<br>
  **Så** skal dialogen si at det ikke finnes interesser eller tildelinger, uten tomme faner.

- [ ] **Arrangement i dialogen åpner redigering**<br>
  **Gitt** at dialogen er åpen.<br>
  **Når** admin velger et arrangement i listen.<br>
  **Så** skal dialogen lukkes og redigeringssiden for arrangementet åpnes.

- [ ] **Åpen dialog tåler liveoppdatering**<br>
  **Gitt** at dialogen er åpen mens en annen admin endrer billettholdere eller interesser.<br>
  **Når** siden oppdateres.<br>
  **Så** skal dialogen forbli åpen og ikke endre kortets høyde eller griddets layout.
