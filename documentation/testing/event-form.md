# Arrangementsskjema

Denne sjekklisten dekker opprettelse og redigering av arrangementer i skjemaet under `/my-events/new/{id}` og adminredigering av arrangementer fra godkjenningsflyten. Filen dekker også inngangen til skjemaet fra `Send inn arrangement`.

## Roller

- Ikke-innlogget bruker
- Innlogget bruker
- Admin

## Sjekkliste

- [ ] `Gitt at en ikke-innlogget bruker trykker på Send inn arrangement fra forsiden, når navigasjonen skjer, så skal brukeren ikke få tilgang til skjemaet, men møtes med en tydelig og forståelig vei videre til innlogging.`
- [ ] `Gitt at en innlogget bruker trykker på Send inn arrangement fra forsiden eller fra egne oversikter, når navigasjonen skjer, så skal brukeren ende i riktig flyt for å opprette eller redigere arrangement og ikke på en uventet mellomside.`
- [ ] `Gitt at en innlogget bruker åpner skjemaet for et nytt eller eksisterende arrangement, når siden lastes, så skal skjemaet vises uten brutte felter, tomme kort eller uforståelig innhold.`
- [ ] `Gitt at skjemaet vises, når brukeren ser overskrifter og seksjoner, så skal status, arrangørinformasjon, arrangementsinformasjon og øvrige detaljer fremstå tydelig adskilt og forståelige.`
- [ ] `Gitt at brukeren endrer navn, e-post eller telefon i arrangørseksjonen, når feltet lagres, så skal endringen bli værende og ikke forsvinne ved oppdatering av siden.`
- [ ] `Gitt at brukeren endrer tittel, intro, type, system eller beskrivelse, når feltene lagres, så skal innholdet persistere og vises konsistent ved refresh og ved senere åpning av skjemaet.`
- [ ] `Gitt at brukeren endrer alder, varighet, nybegynnervennlighet, engelskstøtte, maks antall spillere eller merknader, når feltene lagres, så skal endringene være synlige og korrekte ved senere visning.`
- [ ] `Gitt at brukeren fyller inn ufullstendige, korte eller åpenbart svake data, når brukeren forsøker å sende inn arrangementet, så skal skjemaet ikke oppføre seg som om innsendingen var fullført uten at data faktisk er gyldige.`
- [ ] `Gitt at brukeren skriver ugyldige eller urealistiske verdier i felter som telefonnummer eller maks antall spillere, når brukeren lagrer eller sender inn, så skal brukeropplevelsen gjøre det tydelig om verdiene aksepteres eller avvises.`
- [ ] `Gitt at brukeren arbeider lenge i skjemaet, når flere felt endres etter hverandre, så skal tidligere lagrede felt ikke nullstilles, overskrives eller hoppe mellom verdier.`
- [ ] `Gitt at brukeren åpner skjemaet på nytt etter å ha gjort endringer, når siden lastes på nytt, så skal de siste lagrede verdiene vises og ikke eldre eller delvise versjoner av dataene.`
- [ ] `Gitt at brukeren åpner lenken for å laste opp bilde fra skjemaet, når bildeflyten åpnes og brukeren kommer tilbake, så skal arrangementet fortsatt være knyttet til riktig skjema og riktig arrangement.`
- [ ] `Gitt at brukeren forsøker å sende inn arrangementet når det er klart, når innsendingen lykkes, så skal status og videre flyt være tydelige og brukeren skal ikke bli stående igjen i tvil om at arrangementet faktisk ble sendt inn.`
- [ ] `Gitt at innsending av arrangement feiler, når brukeren forsøker å sende inn, så skal brukeren få en tydelig feiltilstand og ikke bli sendt videre som om innsendingen lyktes.`
- [ ] `Gitt at en vanlig bruker åpner et arrangement som allerede er godkjent, når skjemaet vises, så skal brukeren møtes med tydelig beskjed om at arrangementet ikke kan redigeres videre på vanlig måte.`
- [ ] `Gitt at en vanlig bruker prøver å åpne et arrangement som ikke tilhører dem, når skjemaet åpnes via direkte lenke, så skal brukeren ikke få tilgang til redigering og ikke se en misvisende delvis skjerm.`
- [ ] `Gitt at en admin åpner arrangementet via godkjenningsflyten, når skjemaet vises, så skal admin kunne redigere status og relevante felt uten å møte brukerbegrensningene som gjelder vanlige brukere.`
- [ ] `Gitt at en admin redigerer et arrangement i godkjenningsflyten, når endringene lagres, så skal både skjema og forhåndsvisning oppdatere seg konsistent.`
- [ ] `Gitt at skjemaet brukes på mobil, når mange felt og tekstområder fylles ut, så skal siden fortsatt være lesbar, skrollbar og brukbar uten at felter eller knapper havner utenfor skjermen.`
- [ ] `Gitt at skjemaet brukes med raske endringer i mange felt, når brukeren navigerer mellom felt og tilbake, så skal det ikke oppstå åpenbare race conditions, mistet innhold eller ustabil oppførsel.`

## Kan automatiseres

- Inngangen til skjemaet for ikke-innlogget og innlogget bruker egner seg godt for ende-til-ende-tester som verifiserer riktig adgang og redirect.
- Persistens av felter i skjemaet egner seg godt for ende-til-ende-tester og integrasjonstester som verifiserer at data blir liggende etter refresh.
- Innsending av arrangement, både vellykket og mislykket, egner seg godt for ende-til-ende-tester eller integrasjonstester.
- Begrensning for vanlig bruker på godkjente arrangementer egner seg godt for ende-til-ende-tester.
- Adminredigering med forhåndsvisning egner seg godt for ende-til-ende-tester som verifiserer at skjema og forhåndsvisning holder seg synkronisert.

