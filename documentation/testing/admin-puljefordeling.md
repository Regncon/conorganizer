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
