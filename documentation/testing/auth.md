# Autentisering og autorisering

Denne sjekklisten dekker `/auth`, `/auth/post-login`, `/auth/logout` og tilgangsfeil for beskyttede sider. Descope-flytene er med fordi konfigurasjonen eies av oss og må verifiseres før release.

## Roller

- Ikke-innlogget bruker
- Innlogget bruker

## Sjekkliste

### Innlogging og konto

- [ ] `Gitt at en ikke-innlogget bruker åpner innloggingssiden, når siden lastes, så skal Descope-komponenten vises uten brutte elementer, tomme flater eller tydelig mislykket lasting.`
- [ ] `Gitt at innloggingssiden vises på mobil og desktop, når brukeren skal bruke Descope-flyten, så skal komponenten være lesbar og brukbar uten at viktige felt eller handlinger forsvinner ut av layouten.`
- [ ] `Gitt at en ny bruker registrerer seg med gyldige opplysninger, når Descope-flyten fullføres, så skal brukeren ende i en innlogget tilstand i appen uten halvveis opprettet konto.`
- [ ] `Gitt at registrering, innlogging, e-postbekreftelse eller passordtilbakestilling avvises av Descope, når brukeren har oppgitt ugyldige eller ufullstendige data, så skal tilbakemeldingen være forståelig og ikke sende brukeren videre som om handlingen lyktes.`
- [ ] `Gitt at en bruker logger inn med gyldig konto, når Descope sender brukeren videre til post-login, så skal appen lande på forsiden i en innlogget tilstand.`
- [ ] `Gitt at en bruker logger inn for første gang, når post-login er fullført, så skal Min Side og andre beskyttede funksjoner være tilgjengelige uten delvis opprettet lokal bruker.`
- [ ] `Gitt at en innlogget bruker åpner innloggingssiden på nytt, når siden vises, så skal opplevelsen ikke skape tvil om brukerens faktiske innloggingsstatus.`
- [ ] `Gitt at tredjepartsinnholdet på innloggingssiden svikter, når Descope-komponenten ikke kan lastes normalt, så skal siden ikke fremstå som stille eller fullstendig ødelagt uten forståelig feiltilstand.`
- [ ] `Gitt at brukeren beveger seg inn og ut av autentiseringsflyten, når brukeren bruker tilbakeknapp og refresh, så skal appen oppføre seg konsistent og ikke havne i feil rolle eller halvveis innlogget tilstand.`

### Utlogging

- [ ] `Gitt at en innlogget bruker velger å logge ut, når utloggingen fullføres, så skal brukeren ikke lenger fremstå som innlogget i menyen eller i tilgangen til beskyttede sider.`
- [ ] `Gitt at en bruker nettopp har logget ut, når brukeren refresher siden eller åpner en tidligere beskyttet side, så skal appen fortsatt behandle brukeren som utlogget.`

### Autorisering

- [ ] `Gitt at en ikke-innlogget bruker prøver å åpne en beskyttet side, når brukeren blir avvist, så skal tilgangsfeilen være forståelig og tilby en tydelig vei videre til innlogging eller tilbake til appen.`
