# Generelle tester

Denne sjekklisten dekker oppførsel som går igjen på tvers av flere sider og flyter. Dette er tester som ikke hører naturlig hjemme bare på én side, og som bør vurderes samlet for å sikre en konsistent brukeropplevelse i hele appen.

## Roller

- Ikke-innlogget bruker
- Innlogget bruker
- Admin

## Sjekkliste

- [ ] `Gitt at brukeren åpner en side med hovednavigasjon, når siden er ferdig lastet, så skal navigasjonen fremstå som en stabil og konsistent del av appen uten brutte elementer, feil plassering eller visuelt uferdige tilstander.`
- [ ] `Gitt at brukeren åpner appen på mobil, når hovednavigasjonen brukes på tvers av relevante sider, så skal menyen være lesbar, trykkbar og stabil uten overlapp, avkuttede etiketter eller elementer som havner utenfor skjermen.`
- [ ] `Gitt at brukeren åpner appen på større skjerm, når hovednavigasjonen brukes på tvers av relevante sider, så skal logo, menyknapper og brukermeny oppføre seg konsistent og uten visuelle brudd.`
- [ ] `Gitt at brukeren ser Hjem i hovednavigasjonen, når brukeren trykker på Hjem, så skal brukeren bli sendt til forsiden uten å havne i en mellomtilstand eller på feil side.`
- [ ] `Gitt at brukeren ikke er innlogget, når hovednavigasjonen vises, så skal navigasjonen ikke gi inntrykk av at brukeren har tilgang til Min Side eller adminfunksjoner.`
- [ ] `Gitt at brukeren ikke er innlogget, når innloggingsknappen vises i hovednavigasjonen, så skal det være tydelig at dette er riktig inngang til innlogging.`
- [ ] `Gitt at brukeren ikke er innlogget, når brukeren trykker på innloggingsknappen fra hovednavigasjonen, så skal brukeren bli sendt til innloggingsflyten uten å møte uventede feilmeldinger eller feil side.`
- [ ] `Gitt at brukeren er innlogget, når hovednavigasjonen vises, så skal Min Side være tilgjengelig og presentert på en måte som ikke kan forveksles med adminfunksjoner.`
- [ ] `Gitt at brukeren er innlogget, når brukeren trykker på Min Side fra hovednavigasjonen, så skal brukeren bli sendt til riktig side for egen profil og eget innhold.`
- [ ] `Gitt at brukeren er innlogget, når brukermenyen åpnes, så skal menyen åpne og lukke seg på en stabil måte uten å dekke til eller bryte annen viktig navigasjon.`
- [ ] `Gitt at brukeren er innlogget, når brukermenyen vises, så skal den inneholde valg som er relevante for vanlig bruker og ikke vise adminvalg til brukere som ikke er admin.`
- [ ] `Gitt at brukeren er admin, når hovednavigasjonen vises, så skal admin være tilgjengelig uten at vanlige brukerfunksjoner skjules eller endrer betydning.`
- [ ] `Gitt at brukeren er admin, når brukeren trykker på Admin fra hovednavigasjonen, så skal brukeren bli sendt til adminområdet uten å møte feil rolle eller feil landingsside.`
- [ ] `Gitt at brukeren er admin, når brukermenyen åpnes, så skal adminspesifikke valg vises tydelig og ikke blandes sammen med valg som gjelder vanlig bruker.`
- [ ] `Gitt at brukeren ikke er admin, når brukeren navigerer rundt i appen, så skal adminfunksjoner ikke vises i navigasjonen eller presenteres som tilgjengelige.`
- [ ] `Gitt at brukeren er innlogget, når brukeren velger å logge ut fra brukermenyen, så skal brukeren bli logget ut og etterpå møte en tilstand som tydelig viser at brukeren ikke lenger er innlogget.`
- [ ] `Gitt at brukeren nylig har logget ut, når brukeren navigerer videre i appen, så skal navigasjonen oppføre seg som for en ikke-innlogget bruker og ikke etterlate inntrykk av at brukeren fortsatt er innlogget.`
- [ ] `Gitt at en ikke-innlogget bruker prøver å gå direkte til en beskyttet side, når siden åpnes, så skal brukeren møte en tydelig tilgangsfeil med en forståelig vei videre til innlogging eller tilbake til appen.`
- [ ] `Gitt at en ikke-innlogget bruker havner på en tilgangsfeil, når siden vises, så skal teksten og handlingene på siden være tydelige nok til at brukeren forstår hvorfor siden ikke er tilgjengelig.`
- [ ] `Gitt at en bruker uten adminrettigheter prøver å gå direkte til en adminside, når siden åpnes, så skal brukeren ikke få tilgang og heller ikke møte en halvveis eller misvisende adminvisning.`
- [ ] `Gitt at brukeren navigerer mellom sider via hovednavigasjon og brukermeny, når brukeren bruker tilbakeknapp og refresh, så skal appen oppføre seg stabilt og ikke vise feil rolle, feil menyvalg eller ødelagte tilstander.`
- [ ] `Gitt at brukeren navigerer raskt mellom tilgjengelige sider via menyen, når flere klikk skjer tett etter hverandre, så skal appen ikke havne i duplikathandlinger, feilnavigasjon eller tydelig ustabil tilstand.`
- [ ] `Gitt at brukeren bruker appen med tastatur eller andre alternative navigasjonsformer, når fokus flyttes mellom interaktive elementer i meny og brukermeny, så skal det være mulig å forstå hvor brukeren befinner seg og hvilke handlinger som kan utføres.`
- [ ] `Gitt at eksterne lenker vises i brukermenyen, når brukeren åpner dem, så skal det være tydelig at brukeren forlater eller åpner innhold utenfor appens egne sider.`
- [ ] `Gitt at språk og navngivning brukes i hovednavigasjon og brukermeny, når brukeren beveger seg mellom sider, så skal tekstene være konsistente og ikke skape tvil om hvor brukeren er eller hva et menyvalg gjør.`
- [ ] `Gitt at brukeren ser navigasjonen på tvers av sider, når appen brukes som helhet, så skal navigasjonen fremstå som ferdig og konsistent uten placeholder-preg, utilsiktet språkblanding eller visuelt forstyrrende detaljer.`

## Kan automatiseres

- Synlighet av menyvalg for ikke-innlogget bruker, innlogget bruker og admin egner seg godt for ende-til-ende-tester som verifiserer rollebasert navigasjon.
- Navigasjon til Hjem, Min Side, innlogging og admin egner seg godt for ende-til-ende-tester som verifiserer at riktige sider åpnes fra hovednavigasjonen.
- Logg ut-flyten egner seg godt for en ende-til-ende-test som verifiserer at brukerstatus og menyinnhold endres korrekt etter utlogging.
- Beskyttelse av innloggede sider og adminsider egner seg godt for ende-til-ende-tester eller integrasjonstester som verifiserer at brukeren møtes med riktig tilgangsoppførsel.
- Tastaturnavigasjon i hovedmenyen og brukermenyen egner seg godt for nettleserbaserte ende-til-ende-tester som verifiserer fokusrekkefølge og brukbarhet.
