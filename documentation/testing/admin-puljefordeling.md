# Puljefordeling i admin

Denne sjekklisten dekker `/admin/puljefordeling/{pulje}`, der admin kan se interesser og administrere manuelle GM-, spiller- og førstevalgtildelinger.

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
  **Så** skal spillertildelingen vises på riktig arrangement, og interessen for samme arrangement og pulje skal ikke lenger vises som et vanlig interessevalg.

- [ ] **Spiller kan tildeles som førstevalg**<br>
  **Gitt** at admin har valgt en pulje, et arrangement og en billettholder.<br>
  **Når** billettholderen tildeles som førstevalg.<br>
  **Så** skal tildelingen markeres som førstevalg på riktig arrangement.

- [ ] **Valgdialogen viser alle tre handlinger**<br>
  **Gitt** at admin åpner tildelingsdialogen for et arrangement.<br>
  **Når** en billettholder velges.<br>
  **Så** skal handlingene for GM, spiller og førstevalg være tydelige og knyttet til valgt pulje.

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
  **Så** skal deltakeren plasseres som manuell spiller, interessen skal bli stående i listen, og tildelingsikonet skal oppdateres uten at dialogen lukkes.

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
  **Gitt** at et arrangement har en manuelt tildelt spiller.<br>
  **Når** admin fjerner spillertildelingen.<br>
  **Så** skal bare den valgte spillertildelingen forsvinne.

- [ ] **Førstevalg kan fjernes**<br>
  **Gitt** at et arrangement har en manuelt tildelt førstevalgspiller.<br>
  **Når** admin fjerner førstevalget.<br>
  **Så** skal både den manuelle plassen og førstevalgmarkeringen for dette arrangementet og denne puljen fjernes.

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

- [ ] **Refresh viser lagret tildelingstilstand**<br>
  **Gitt** at admin har utført flere tildelingshandlinger.<br>
  **Når** siden lastes på nytt.<br>
  **Så** skal GM-, spiller- og førstevalgstatus samsvare med de lagrede tildelingene.
