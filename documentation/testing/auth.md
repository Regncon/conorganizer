# Autentisering

Denne sjekklisten dekker autentiseringsflyten på `/auth`, inkludert innlogging, utlogging, videreføring etter innlogging og oppførsel når brukeren mangler tilgang til beskyttede sider.

## Roller

- Ikke-innlogget bruker
- Innlogget bruker

## Sjekkliste

- [ ] `Gitt at en ikke-innlogget bruker åpner innloggingssiden, når siden lastes, så skal innloggingskomponenten vises uten brutte elementer, tomme flater eller tydelig mislykket lasting.`
- [ ] `Gitt at en ikke-innlogget bruker åpner innloggingssiden, når siden vises på mobil og desktop, så skal innloggingskomponenten være lesbar og brukbar uten at viktige felt eller handlinger forsvinner ut av layouten.`
- [ ] `Gitt at en ny bruker ønsker å opprette konto, når registreringsflyten brukes med gyldige opplysninger, så skal brukeren kunne fullføre registreringen uten å havne i en uklar eller halvveis opprettet tilstand.`
- [ ] `Gitt at en bruker forsøker å registrere seg med ugyldige eller ufullstendige opplysninger, når registreringen avvises, så skal brukeren få en tydelig tilbakemelding og ikke bli sendt videre som om registreringen lyktes.`
- [ ] `Gitt at en bruker forsøker å logge inn med gyldig konto, når innloggingen fullføres, så skal brukeren bli sendt videre til riktig post-login-flyt og ende på en innlogget tilstand i appen.`
- [ ] `Gitt at en bruker logger inn for første gang, når innloggingen fullføres, så skal brukeren kunne tas i bruk i appen uten at Min Side eller andre beskyttede funksjoner fremstår som delvis opprettet eller utilgjengelige.`
- [ ] `Gitt at en bruker forsøker å logge inn med ugyldige eller ufullstendige opplysninger, når innloggingen avvises, så skal brukeren få en forståelig tilbakemelding og ikke bli stående i en utydelig mellomtilstand.`
- [ ] `Gitt at en bruker må bekrefte e-post som del av registrering eller innlogging, når bekreftelsesflyten brukes, så skal det være tydelig hva som forventes av brukeren og om bekreftelsen faktisk lyktes eller feilet.`
- [ ] `Gitt at en bruker forsøker å bruke kode eller bekreftelsessteg som ikke lenger er gyldig, når flyten avvises, så skal brukeren få en tydelig feilmelding og ikke misvisende tekst som får det til å se ut som om et nytt steg allerede er sendt eller fullført.`
- [ ] `Gitt at en bruker skal tilbakestille passordet sitt, når glemt-passord-flyten brukes med gyldig konto, så skal brukeren kunne fullføre flyten uten å havne i en uklar tilstand om hva som ble sendt eller bekreftet.`
- [ ] `Gitt at en bruker oppgir feil eller tvetydig informasjon i glemt-passord-flyten, når systemet ikke kan fullføre normalt, så skal tilbakemeldingen være forståelig nok til at brukeren skjønner hva som kan prøves på nytt.`
- [ ] `Gitt at en bruker avbryter eller forlater innloggingsflyten underveis, når brukeren går tilbake til appen, så skal appen fortsatt tydelig vise at brukeren ikke er innlogget.`
- [ ] `Gitt at en ikke-innlogget bruker prøver å åpne en beskyttet side, når brukeren blir avvist, så skal tilgangsfeilen være forståelig og tilby en tydelig vei videre til innlogging eller tilbake til appen.`
- [ ] `Gitt at en ikke-innlogget bruker står på tilgangsfeilsiden, når brukeren velger å logge inn derfra, så skal brukeren sendes til innloggingssiden uten feil eller uventet omvei.`
- [ ] `Gitt at en innlogget bruker velger å logge ut, når utloggingen fullføres, så skal brukeren ikke lenger fremstå som innlogget i menyen eller i tilgangen til beskyttede sider.`
- [ ] `Gitt at en bruker nettopp har logget ut, når brukeren refresher siden eller åpner en tidligere beskyttet side, så skal appen fortsatt behandle brukeren som utlogget.`
- [ ] `Gitt at en innlogget bruker åpner innloggingssiden på nytt, når siden vises, så skal opplevelsen ikke skape tvil om brukerens faktiske innloggingsstatus.`
- [ ] `Gitt at nettverk eller tredjepartsinnhold på innloggingssiden svikter, når innloggingskomponenten ikke kan lastes normalt, så skal siden ikke fremstå som stille eller fullstendig ødelagt uten at det er mulig å forstå at noe gikk galt.`
- [ ] `Gitt at brukeren beveger seg inn og ut av autentiseringsflyten flere ganger, når brukeren bruker tilbakeknapp og refresh, så skal appen oppføre seg konsistent og ikke havne i feil rolle eller halvt innlogget tilstand.`
- [ ] `Gitt at autentiseringsrelatert tekst vises til brukeren, når innlogging, utlogging og avvisning skjer, så skal språket være forståelig og ikke gi motstridende signaler om hva som nettopp har skjedd.`

## Kan automatiseres

- Registrering, innlogging og glemt-passord-flyt egner seg godt for ende-til-ende-tester så lenge testmiljøet støtter stabile testkontoer og kontrollerbare bekreftelsessteg.
- Innlogging med gyldig og ugyldig bruker egner seg godt for ende-til-ende-tester som verifiserer overgang mellom ikke-innlogget og innlogget tilstand.
- Utlogging egner seg godt for en ende-til-ende-test som verifiserer at beskyttede sider ikke lenger er tilgjengelige etter at brukeren har logget ut.
- Tilgangsfeil for beskyttede sider egner seg godt for ende-til-ende-tester eller integrasjonstester som verifiserer at riktig side og riktig tekst vises til ikke-innlogget bruker.
- Post-login-flyten egner seg godt for en ende-til-ende-test som verifiserer at brukeren lander på riktig sted og faktisk kan bruke innlogget funksjonalitet etterpå.
