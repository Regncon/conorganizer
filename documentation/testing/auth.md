# Autentisering og autorisering

Denne sjekklisten dekker `/auth`, `/auth/post-login`, `/auth/logout` og tilgangsfeil for beskyttede sider. Descope-flytene er med fordi konfigurasjonen eies av oss og må verifiseres før release.

## Roller

- Ikke-innlogget bruker
- Innlogget bruker

## Sjekkliste

### Innlogging og konto

- [ ] **Descope-komponenten laster lesbart**<br>
  **Gitt** at en ikke-innlogget bruker åpner innloggingssiden.<br>
  **Når** siden lastes.<br>
  **Så** skal Descope-komponenten vises uten brutte elementer, tomme flater eller tydelig mislykket lasting.

- [ ] **Innloggingsflyten fungerer på mobil og desktop**<br>
  **Gitt** at innloggingssiden vises på mobil og desktop.<br>
  **Når** brukeren skal bruke Descope-flyten.<br>
  **Så** skal komponenten være lesbar og brukbar uten at viktige felt eller handlinger forsvinner ut av layouten.

- [ ] **Registrering gir innlogget tilstand**<br>
  **Gitt** at en ny bruker registrerer seg med gyldige opplysninger.<br>
  **Når** Descope-flyten fullføres.<br>
  **Så** skal brukeren ende i en innlogget tilstand i appen uten halvveis opprettet konto.

- [ ] **Descope-avvisninger gir forståelig tilbakemelding**<br>
  **Gitt** at registrering, innlogging, e-postbekreftelse eller passordtilbakestilling avvises av Descope.<br>
  **Når** brukeren har oppgitt ugyldige eller ufullstendige data.<br>
  **Så** skal tilbakemeldingen være forståelig og ikke sende brukeren videre som om handlingen lyktes.

- [ ] **Gyldig innlogging lander på forsiden**<br>
  **Gitt** at en bruker logger inn med gyldig konto.<br>
  **Når** Descope sender brukeren videre til post-login.<br>
  **Så** skal appen lande på forsiden i en innlogget tilstand.

- [ ] **Førstegangsinnlogging oppretter lokal bruker riktig**<br>
  **Gitt** at en bruker logger inn for første gang.<br>
  **Når** post-login er fullført.<br>
  **Så** skal Min Side og andre beskyttede funksjoner være tilgjengelige uten delvis opprettet lokal bruker.

- [ ] **Innlogget bruker møter ingen uklar innloggingsstatus**<br>
  **Gitt** at en innlogget bruker åpner innloggingssiden på nytt.<br>
  **Når** siden vises.<br>
  **Så** skal opplevelsen ikke skape tvil om brukerens faktiske innloggingsstatus.

- [ ] **Descope-lastingsfeil gir forståelig feiltilstand**<br>
  **Gitt** at tredjepartsinnholdet på innloggingssiden svikter.<br>
  **Når** Descope-komponenten ikke kan lastes normalt.<br>
  **Så** skal siden ikke fremstå som stille eller fullstendig ødelagt uten forståelig feiltilstand.

- [ ] **Tilbakeknapp og refresh bevarer riktig auth-tilstand**<br>
  **Gitt** at brukeren beveger seg inn og ut av autentiseringsflyten.<br>
  **Når** brukeren bruker tilbakeknapp og refresh.<br>
  **Så** skal appen oppføre seg konsistent og ikke havne i feil rolle eller halvveis innlogget tilstand.

### Utlogging

- [ ] **Utlogging fjerner innlogget tilgang**<br>
  **Gitt** at en innlogget bruker velger å logge ut.<br>
  **Når** utloggingen fullføres.<br>
  **Så** skal brukeren ikke lenger fremstå som innlogget i menyen eller i tilgangen til beskyttede sider.

- [ ] **Refresh etter utlogging holder brukeren utlogget**<br>
  **Gitt** at en bruker nettopp har logget ut.<br>
  **Når** brukeren refresher siden eller åpner en tidligere beskyttet side.<br>
  **Så** skal appen fortsatt behandle brukeren som utlogget.

### Autorisering

- [ ] **Beskyttet side gir tydelig vei videre**<br>
  **Gitt** at en ikke-innlogget bruker prøver å åpne en beskyttet side.<br>
  **Når** brukeren blir avvist.<br>
  **Så** skal tilgangsfeilen være forståelig og tilby en tydelig vei videre til innlogging eller tilbake til appen.
