# Min Side

Denne sjekklisten dekker `/profile`: innlogget oversikt, egne arrangementer, kort oppsummering av billetter, eget festivalprogram, lenke videre til billettadministrasjon og invitasjonen til å gi tilbakemelding.

## Roller

- Innlogget bruker

## Sjekkliste

### Oversikt

- [ ] **Min Side viser en helhetlig oversikt**<br>
  **Gitt** at en innlogget bruker åpner Min Side.<br>
  **Når** siden lastes.<br>
  **Så** skal siden vises som en helhetlig oversikt uten brutte seksjoner eller tydelig manglende innhold.

- [ ] **Min Side inviterer til å gi tilbakemelding**<br>
  **Gitt** at en innlogget bruker åpner Min Side.<br>
  **Når** boksen «Har du en tilbakemelding?» vises over Kontoadministrasjon.<br>
  **Så** skal boksen forklare at brukeren kan si hva hen synes om nettsiden eller festivalen, og knappen «Gi tilbakemelding» skal åpne `/tilbakemelding`.

### Mine arrangementer

- [ ] **Nytt arrangement åpner riktig skjema**<br>
  **Gitt** at brukeren trykker «Send inn arrangement» i Mine arrangementer.<br>
  **Når** handlingen lykkes.<br>
  **Så** skal brukeren få opprettet et nytt arrangement med status Kladd og sendes videre til skjemaet `/profile/new/{id}` for det nye arrangementet.

- [ ] **Arrangementskort lenker etter status**<br>
  **Gitt** at brukeren har arrangementer med ulike statuser.<br>
  **Når** brukeren trykker på kortene i Mine arrangementer.<br>
  **Så** åpner Innsendt, Godkjent og Annonsert arrangementsiden `/event/{id}`, mens Kladd og Forkastet åpner skjemaet `/profile/new/{id}`.

- [ ] **Arrangement uten tittel vises som «Mangler navn»**<br>
  **Gitt** at et av brukerens arrangementer mangler tittel.<br>
  **Når** Mine arrangementer vises.<br>
  **Så** viser kortet «Mangler navn» og lenker fortsatt til arrangementet.

- [ ] **Godkjent arrangement kan ikke redigeres av vanlig bruker**<br>
  **Gitt** at et arrangement har status Godkjent og brukeren ikke er admin.<br>
  **Når** brukeren åpner `/profile/new/{id}` direkte.<br>
  **Så** vises «Arrangementet er allerede godkjent eller publisert» og beskjed om å ta kontakt med RegnCon-styret, i stedet for skjemaet.

- [ ] **Mange arrangementer forblir lesbare**<br>
  **Gitt** at brukeren har mange arrangementer.<br>
  **Når** seksjonen for Mine arrangementer vises.<br>
  **Så** skal kortene fortsatt være lesbare og navigerbare uten å skape kaotisk layout.

- [ ] **Refresh viser lagret arrangementstilstand**<br>
  **Gitt** at brukeren refresher siden etter å ha opprettet eller endret arrangementer.<br>
  **Når** siden vises igjen.<br>
  **Så** skal oversikten samsvare med faktisk lagret tilstand.

- [ ] **Tilbakeknapp bevarer stabil Min Side**<br>
  **Gitt** at brukeren bruker tilbakeknapp mellom Min Side og underliggende arrangementsider.<br>
  **Når** brukeren kommer tilbake til Min Side.<br>
  **Så** skal siden fortsatt fremstå stabil og oppdatert.

### Billetter og festivalprogram

- [ ] **Billettseksjonen håndterer manglende billettholdere**<br>
  **Gitt** at brukeren ikke har billettholdere knyttet til seg.<br>
  **Når** Min Side vises.<br>
  **Så** viser billettseksjonen «Ingen billetter er knyttet til profilen din ennå.» og knappen «Hent billetter», uten at hele siden fremstår som feil eller mangelfull.

- [ ] **Festivalprogrammet venter på publisert program**<br>
  **Gitt** at programmet ikke er publisert.<br>
  **Når** Min Side vises.<br>
  **Så** viser Mitt festivalprogram bare «Programmet for Regncon er ikke publisert ennå», også om brukeren har interesser eller tildelinger.

- [ ] **Festivalprogrammet viser tildelinger etter reglene**<br>
  **Gitt** at programmet er publisert og billettholderen har interesser og tildelinger i ulike puljer.<br>
  **Når** Mitt festivalprogram vises.<br>
  **Så** vises GM-tildelinger og manuelle spillerplasseringer med en gang, mens spillerplasseringer fra puljefordelingen først vises når puljen er Fullført. Før det vises billettholderens interesser for puljen. En pulje med synlig tildeling viser ikke interesselisten.

- [ ] **Programarrangement viser til beskrivelsen for tidspunkt**<br>
  **Gitt** at programmet er publisert og billettholderen er GM eller manuelt satt opp på et programarrangement.<br>
  **Når** Mitt festivalprogram vises.<br>
  **Så** viser kortet for programarrangementet «Se beskrivelse for tidspunkt». Puljeoverskriften beholder klokkeslettet, og kort for arrangementer i puljefordelingen har ikke denne teksten.

- [ ] **Tomme programpunkter forklares tydelig**<br>
  **Gitt** at programmet er publisert og brukeren ikke har noe i en eller flere puljer i sitt festivalprogram.<br>
  **Når** Min Side vises.<br>
  **Så** viser puljen «Vi har ingenting å vise her ennå.» med knappen «Se arrangementer», og tomtilstanden ser ikke ut som om data har forsvunnet.

- [ ] **Billettlenken åpner riktig underside**<br>
  **Gitt** at brukeren trykker seg videre til billettsiden fra Min Side.<br>
  **Når** navigasjonen skjer.<br>
  **Så** skal brukeren ende på `/profile/tickets` uten feil kontekst eller feil rolle. Knappen heter «Hent billetter» uten billetter og «Mine billetter» når brukeren har billetter.

### Visning og navigasjon

- [ ] **Min Side fungerer på mobil**<br>
  **Gitt** at brukeren bruker Min Side på mobil.<br>
  **Når** de ulike seksjonene vises under hverandre.<br>
  **Så** skal kort, paneler og overskrifter være lesbare og ikke overlappe.

- [ ] **Kolonner er balansert på større skjerm**<br>
  **Gitt** at brukeren bruker Min Side på større skjerm.<br>
  **Når** innholdet fordeles i kolonner.<br>
  **Så** skal seksjonene fremstå balansert og uten uventede tomrom eller kollisjon mellom paneler.
