# Billetter på Min Side

Denne sjekklisten dekker `/profile/tickets`, der innlogget bruker kan hente billetter, se billettholdere og legge til eller slette manuelle e-postadresser.

## Roller

- Innlogget bruker

## Sjekkliste

- [ ] `Gitt at en innlogget bruker åpner billettsiden, når siden lastes, så skal brødsmulesti, overskrift og hovedinnhold vises uten brutte paneler eller uforståelige feilmeldinger.`
- [ ] `Gitt at brukeren allerede har billettholdere knyttet til seg, når billettsiden vises, så skal billettholderne presenteres med navn, billettype, aldersinformasjon og relevante e-postadresser.`
- [ ] `Gitt at brukeren ikke har billettholdere knyttet til seg, når billettsiden vises, så skal tomtilstanden være forståelig og gi mening uten å se ut som en ren teknisk feil.`
- [ ] `Gitt at brukeren ikke har billetter på sin e-postadresse, når tomtilstanden vises, så skal videre handlingsvalg være forståelige og ikke gi inntrykk av at appen er låst.`
- [ ] `Gitt at brukeren trykker på Hent billetter, når billettene faktisk kan hentes, så skal billettholdere dukke opp uten at siden havner i en stille eller uavklart tilstand.`
- [ ] `Gitt at brukeren trykker på Hent billetter, når henting pågår, så skal knappen og lasteindikatoren oppføre seg på en måte som gjør det tydelig at en handling er i gang.`
- [ ] `Gitt at henting av billetter feiler, når brukeren forsøker å hente billetter, så skal brukeren få en tydelig feilmelding og ikke en falsk bekreftelse på at alt gikk bra.`
- [ ] `Gitt at brukeren trykker på Hent billetter flere ganger raskt, når siden håndterer forespørslene, så skal det ikke oppstå duplisering eller åpenbart ustabil oppførsel i resultatet.`
- [ ] `Gitt at en billettholder har flere e-postadresser, når kortet vises, så skal det være tydelig hvilke e-postadresser som er billettadresse og hvilke som er andre tilknyttede adresser.`
- [ ] `Gitt at brukeren legger til en ny manuell e-postadresse på en billettholder, når handlingen lykkes, så skal brukeren få en tydelig bekreftelse og se at e-postadressen faktisk er lagt til riktig billettholder.`
- [ ] `Gitt at brukeren forsøker å legge til en tom e-postadresse, når handlingen sendes inn, så skal brukeren få en tydelig feilmelding og ingen ny e-postadresse skal legges til.`
- [ ] `Gitt at brukeren forsøker å legge til en e-postadresse som allerede finnes på samme billettholder, når handlingen sendes inn, så skal brukeren få en forståelig feilmelding og ingen duplikatadresse skal opprettes.`
- [ ] `Gitt at brukeren legger til en ny e-postadresse, når siden oppdateres, så skal tidligere data på siden fortsatt være intakte og ikke forsvinne eller byttes om mellom kortene.`
- [ ] `Gitt at en manuell e-postadresse finnes på en billettholder, når brukeren velger å slette den, så skal den slettes fra riktig billettholder og brukeren skal få en tydelig bekreftelse.`
- [ ] `Gitt at e-postadressen ikke er manuell, når brukeren forsøker eller forventer å kunne slette den, så skal siden tydelig håndheve at slike adresser ikke kan slettes på samme måte.`
- [ ] `Gitt at sletting av e-postadresse feiler, når brukeren forsøker å slette, så skal brukeren få en feilmelding som gjør det tydelig at endringen ikke ble fullført.`
- [ ] `Gitt at siden viser meldinger om vellykket eller mislykket endring, når flere handlinger utføres etter hverandre, så skal meldingene høre til riktig billettholder og ikke lekke over til andre kort.`
- [ ] `Gitt at brukeren bruker billettsiden på mobil, når mange billettholderkort eller lange e-postadresser vises, så skal innholdet fortsatt være lesbart og brukbart uten overlapp eller horisontal kollaps.`
- [ ] `Gitt at brukeren refresher siden etter å ha hentet billetter eller endret e-postadresser, når siden lastes på nytt, så skal resultatet samsvare med faktisk lagret tilstand.`

## Kan automatiseres

- Henting av billetter egner seg godt for ende-til-ende-tester som verifiserer både tomtilstand, vellykket henting og tydelig feiltilstand.
- Legg til og slett manuell e-postadresse egner seg godt for ende-til-ende-tester eller integrasjonstester som verifiserer at endringene havner på riktig billettholder.
- Duplikat- og tom-felt-validering egner seg godt for integrasjonstester og ende-til-ende-tester.
- Oppdatering av suksess- og feilmeldinger på riktig kort egner seg godt for nettleserbaserte ende-til-ende-tester.

