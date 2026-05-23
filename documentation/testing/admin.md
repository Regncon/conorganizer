# Admin

Denne sjekklisten dekker hovedsiden for admin på `/admin`, altså inngangen til administrative funksjoner.

## Roller

- Bruker uten adminrettigheter
- Admin

## Sjekkliste

- [ ] `Gitt at en admin åpner adminforsiden, når siden lastes, så skal adminområdets hovedvalg vises tydelig og uten brutte paneler eller feil rolleopplevelse.`
- [ ] `Gitt at en admin åpner adminforsiden, når siden er ferdig lastet, så skal brødsmulestien tydelig vise at brukeren er i adminområdet.`
- [ ] `Gitt at en admin står på adminforsiden, når kortene for videre adminarbeid vises, så skal det være tydelig hva hvert valg leder til.`
- [ ] `Gitt at en admin velger å gå til godkjenning av arrangementer, når navigasjonen skjer, så skal brukeren havne på riktig underside uten å møte feil eller uventet mellomtilstand.`
- [ ] `Gitt at en admin velger å gå til billettholderoversikten, når navigasjonen skjer, så skal brukeren havne på riktig underside uten feil rolle eller feil side.`
- [ ] `Gitt at adminforsiden brukes på mobil og større skjerm, når kortene vises, så skal de være lesbare, klikkbare og visuelt stabile uten at tekst eller bilder kolliderer.`
- [ ] `Gitt at en bruker uten adminrettigheter prøver å åpne adminforsiden direkte, når siden lastes, så skal brukeren ikke få tilgang og heller ikke møte en misvisende halvveis adminvisning.`
- [ ] `Gitt at adminforsiden ikke kan laste nødvendig innhold som forventet, når siden vises, så skal brukeren ikke bli stående med en tilsynelatende tom adminside uten forklaring.`
- [ ] `Gitt at admin går frem og tilbake mellom adminforsiden og underliggende adminsider, når brukeren bruker tilbakeknapp og refresh, så skal adminområdet fortsatt oppføre seg konsistent og tydelig som adminområde.`

## Kan automatiseres

- Tilgangskontroll til adminforsiden egner seg godt for ende-til-ende-tester.
- Navigasjon fra adminforsiden til underliggende adminsider egner seg godt for ende-til-ende-tester.
- Responsiv presentasjon av adminkortene egner seg godt for nettleserbaserte ende-til-ende-tester på ulike skjermstørrelser.

