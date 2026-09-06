# Puljefordeling i admin

Denne sjekklisten dekker `/admin/puljefordeling/{pulje}`, der admin kan se interesser og administrere manuelle GM- og spillertildelinger uten å endre interessene.

## Roller

- Admin

## Sjekkliste

### Tildeling

- [ ] **GM tildeles riktig arrangement og pulje**<br>
  **Gitt** at admin har valgt en pulje og et arrangement.<br>
  **Når** en billettholder tildeles som GM.<br>
  **Så** skal GM-tildelingen vises på riktig arrangement uten å bli fremstilt som en spillertildeling.

- [ ] **Spiller tildeles riktig arrangement og pulje**<br>
  **Gitt** at admin har valgt en pulje og et arrangement.<br>
  **Når** en billettholder tildeles som spiller.<br>
  **Så** skal spillertildelingen vises på riktig arrangement uten at eksisterende interesser endres.

- [ ] **Valgdialogen viser de to tildelingshandlingene**<br>
  **Gitt** at admin åpner tildelingsdialogen for et arrangement.<br>
  **Når** en billettholder velges.<br>
  **Så** skal handlingene for GM og spiller være tydelige og knyttet til valgt pulje, uten en handling for å sette førstevalg.

### Interesselisten

- [ ] **Interesser vises for valgt arrangement**<br>
  **Gitt** at flere har meldt interesse på ulike nivåer for arrangementet.<br>
  **Når** admin åpner tildelingsdialogen.<br>
  **Så** skal interessene vises under det vanlige søket, gruppert som veldig interessert, interessert og litt interessert.

- [ ] **Navn og billettype vises på to linjer**<br>
  **Gitt** at en interessert deltaker har en billettype fra Chicken.no.<br>
  **Når** interessen vises i dialogen.<br>
  **Så** skal navnet og statusikonene stå på første linje, og billettypen på den andre.

- [ ] **Statusikoner gir riktig informasjon**<br>
  **Gitt** interesserte som er under 18, allerede tildelt i denne puljen eller har fått førstevalget sitt i en tidligere pulje.<br>
  **Når** interesselisten vises.<br>
  **Så** skal hver relevant status ha et ikon med forklarende hjelpetekst, og allerede tildelte skal vises mørkere.

- [ ] **Klikk på en interesse plasserer spilleren direkte**<br>
  **Gitt** at en deltaker har en interesse for valgt arrangement.<br>
  **Når** admin klikker på interessen.<br>
  **Så** skal deltakeren plasseres som manuell spiller og dialogen lukkes. Når dialogen åpnes igjen, skal interessen fortsatt stå i listen med oppdatert tildelingsikon.

- [ ] **Beholdt veldig interesse teller som førstevalg**<br>
  **Gitt** at deltakeren er veldig interessert og ikke har fått førstevalget sitt i en tidligere pulje.<br>
  **Når** admin plasserer deltakeren fra interesselisten.<br>
  **Så** skal plasseringen telle som deltakerens førstevalg.

- [ ] **Aldersbekreftelse gjelder også interesselisten**<br>
  **Gitt** at en deltaker under 18 er interessert i et 18+-arrangement.<br>
  **Når** admin klikker på interessen.<br>
  **Så** skal den eksisterende aldersadvarselen kreve bekreftelse før plasseringen lagres.

### Fjerning

- [ ] **GM kan fjernes**<br>
  **Gitt** at et arrangement har en manuelt tildelt GM.<br>
  **Når** admin fjerner GM-tildelingen.<br>
  **Så** skal bare den valgte GM-tildelingen forsvinne.

- [ ] **Spiller kan fjernes**<br>
  **Gitt** at et arrangement har en manuelt tildelt spiller med en eksisterende interesse.<br>
  **Når** admin fjerner spillertildelingen.<br>
  **Så** skal bare den valgte spillertildelingen forsvinne, mens interessen beholdes.

### Avgrensning og robusthet

- [ ] **Påmeldte vises uten ordinær solverplass**<br>
  **Gitt** at en billettholder er direkte påmeldt og fortsatt har en ordinær interesse i samme pulje.<br>
  **Når** admin åpner og lagrer puljefordelingen.<br>
  **Så** skal billettholderen vises som påmeldt på riktig arrangement, ikke få en ordinær solverplass i puljen, og beholde påmeldingskilden etter lagring.

- [ ] **Påmeldinger beholdes uavhengig av manuelle tildelinger**<br>
  **Gitt** at billettholderen har en direkte påmelding og admin endrer en manuell tildeling i samme pulje.<br>
  **Når** admin lagrer endringen.<br>
  **Så** skal den uavhengige påmeldingen beholdes.

- [ ] **Handlinger påvirker bare valgt pulje**<br>
  **Gitt** at billettholderen har data i flere puljer.<br>
  **Når** admin tildeler eller fjerner en rolle i én pulje.<br>
  **Så** skal tildelinger og interesser i andre puljer forbli uendret.

- [ ] **Tildelingshandlinger endrer aldri interesser**<br>
  **Gitt** at billettholderen har interesser på arrangementet som tildeles og på andre arrangementer.<br>
  **Når** admin tildeler, flytter eller fjerner en spiller eller GM.<br>
  **Så** skal alle interessene og nivåene deres være uendret.

- [ ] **Refresh viser lagret tildelingstilstand**<br>
  **Gitt** at admin har utført flere tildelingshandlinger.<br>
  **Når** siden lastes på nytt.<br>
  **Så** skal GM- og spillerstatus samsvare med de lagrede tildelingene, mens førstevalg fortsatt utledes fra lagrede interesser.
