# Billettholdere i admin

Denne sjekklisten dekker `/admin/billettholder`, der admin får oversikt over alle billettholdere og kan vedlikeholde manuelle e-postadresser.

## Roller

- Admin

## Sjekkliste

- [ ] `Gitt at en admin åpner billettholderoversikten, når siden lastes, så skal oversikten vises med tydelig brødsmulesti, overskrift og innhold uten brutte kort eller tydelig manglende data.`
- [ ] `Gitt at billettholderoversikten inneholder mange deltakere, når siden vises, så skal grid og kort forbli lesbare og brukbare uten sammenfallende innhold.`
- [ ] `Gitt at en billettholder vises i adminoversikten, når kortet leses, så skal bestilling, type, navn, alder og relevante e-postadresser fremstå tydelig.`
- [ ] `Gitt at en billettholder har flere e-postadresser, når kortet vises, så skal det være forståelig hvilke som er billettadresse og hvilke som er andre tilknyttede adresser.`
- [ ] `Gitt at admin legger til en manuell e-postadresse på en billettholder, når handlingen lykkes, så skal bekreftelsen vises på riktig kort og den nye adressen vises på riktig billettholder.`
- [ ] `Gitt at admin forsøker å legge til en tom e-postadresse, når handlingen utføres, så skal admin få en tydelig feilmelding og ingen adresse skal legges til.`
- [ ] `Gitt at admin forsøker å legge til en e-postadresse som allerede finnes på samme billettholder, når handlingen utføres, så skal siden avvise duplikatet tydelig og uten å skape uklar tilstand.`
- [ ] `Gitt at admin sletter en manuell e-postadresse, når handlingen lykkes, så skal adressen fjernes fra riktig kort og ikke bli stående igjen på siden som om den fortsatt eksisterer.`
- [ ] `Gitt at admin forsøker å slette en ikke-manual e-postadresse, når handlingen utføres eller forventes, så skal systemet tydelig håndheve at slike adresser ikke kan slettes som en manuell adresse.`
- [ ] `Gitt at sletting av e-postadresse medfører at bruker-tilknytning må ryddes opp, når handlingen lykkes, så skal resultatet fremstå konsistent og ikke etterlate spor av delvis sletting i brukeropplevelsen.`
- [ ] `Gitt at en add- eller delete-handling feiler, når admin utfører endringen, så skal feilmeldingen være tydelig og ikke etterlate inntrykk av at endringen likevel ble lagret.`
- [ ] `Gitt at admin jobber med flere billettholderkort på samme side, når flere endringer skjer etter hverandre, så skal meldinger og oppdateringer tilhøre riktig kort og ikke lekke til andre kort.`
- [ ] `Gitt at admin klikker seg videre til å legge til billettholder fra oversikten, når navigasjonen skjer, så skal riktig underside åpnes uten feil rolle eller feil kontekst.`
- [ ] `Gitt at siden brukes på mobil eller smal skjerm, når mange billettholdere eller lange e-postadresser vises, så skal innholdet fortsatt være lesbart og trykkbart uten at kortene bryter sammen.`
- [ ] `Gitt at admin refresher oversikten etter å ha endret e-postadresser, når siden lastes inn igjen, så skal den reflektere den faktiske lagrede tilstanden.`

## Kan automatiseres

- Visning av billettholderkort med ulike datakombinasjoner egner seg godt for ende-til-ende-tester og integrasjonstester.
- Legg til og slett manuelle e-postadresser egner seg godt for ende-til-ende-tester og integrasjonstester.
- Feilhåndtering for tomme og dupliserte e-postadresser egner seg godt for integrasjonstester.
- Riktig plassering av suksess- og feilmeldinger på riktig kort egner seg godt for ende-til-ende-tester.

