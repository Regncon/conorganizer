# Mine arrangementer

Denne sjekklisten dekker `/my-events`, altså oversikten over brukerens egne arrangementer og inngangen til å opprette eller åpne arrangementer for videre arbeid.

## Roller

- Innlogget bruker

## Sjekkliste

- [ ] `Gitt at en innlogget bruker åpner Mine arrangementer, når siden lastes, så skal siden vise en tydelig oversikt over egne arrangementer uten brutte kort eller tomme seksjoner som ser utilsiktede ut.`
- [ ] `Gitt at brukeren åpner Mine arrangementer, når siden er ferdig lastet, så skal brødsmulestien tydelig vise at brukeren er på Mine arrangementer.`
- [ ] `Gitt at brukeren ikke har noen arrangementer ennå, når siden vises, så skal opprettelsesinngangen være tydelig og siden skal ikke se tom eller ødelagt ut.`
- [ ] `Gitt at brukeren har ett eller flere arrangementer, når siden vises, så skal hvert kort representere riktig arrangement og ikke blande data mellom ulike arrangementer.`
- [ ] `Gitt at et arrangement mangler tittel, system eller ingress, når kortet vises, så skal kortet fortsatt være forståelig og ikke bryte layouten.`
- [ ] `Gitt at et arrangement er i kladd, når brukeren åpner det fra oversikten, så skal brukeren sendes til riktig redigeringsflyt.`
- [ ] `Gitt at et arrangement er innsendt eller publisert, når brukeren åpner det fra oversikten, så skal brukeren sendes til riktig side for den statusen og ikke feilaktig inn i kladdredigering.`
- [ ] `Gitt at brukeren velger å legge til nytt arrangement, når handlingen lykkes, så skal brukeren få opprettet et nytt arrangement og sendes videre til riktig skjema for det nye arrangementet.`
- [ ] `Gitt at opprettelse av nytt arrangement feiler, når brukeren forsøker å opprette fra oversikten, så skal brukeren ikke havne i en stille feil eller på en side som later som om arrangementet ble opprettet.`
- [ ] `Gitt at brukeren har mange arrangementer, når oversikten vises, så skal kortene fortsatt være lesbare og navigerbare uten å skape kaotisk layout.`
- [ ] `Gitt at brukeren åpner Mine arrangementer på mobil, når kortene og opprettelsesinngangen vises, så skal alle interaktive elementer være trykkbare og lesbare uten å bryte grid-oppsettet.`
- [ ] `Gitt at brukeren refresher siden etter å ha opprettet eller endret arrangementer, når siden vises igjen, så skal oversikten samsvare med faktisk lagret tilstand.`
- [ ] `Gitt at brukeren prøver å åpne et arrangement som ikke tilhører dem via direkte lenke, når siden åpnes, så skal brukeren ikke få tilgang til å redigere arrangementet og heller ikke møte en misvisende delvis visning.`
- [ ] `Gitt at brukeren bruker tilbakeknapp mellom oversikten og underliggende arrangementsider, når brukeren kommer tilbake til oversikten, så skal siden fortsatt fremstå stabil og oppdatert.`

## Kan automatiseres

- Visning av tom og ikke-tom oversikt egner seg godt for ende-til-ende-tester med seedet database.
- Opprettelse av nytt arrangement fra oversikten egner seg godt for ende-til-ende-tester som verifiserer redirect til riktig skjema.
- Lenkeoppførsel basert på arrangementsstatus egner seg godt for ende-til-ende-tester.
- Beskyttelse mot å åpne andres arrangementer egner seg godt for ende-til-ende-tester eller integrasjonstester.

