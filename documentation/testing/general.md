# Generelle tester

Denne sjekklisten dekker felles navigasjon, rolleopplevelse og tilgang på tvers av appen. Den dekker også `/auth`, `/auth/post-login`, `/auth/logout` og tilgangsfeil for beskyttede sider. Descope-flytene er med fordi konfigurasjonen eies av oss og må verifiseres før release.

## Roller

- Ikke-innlogget bruker
- Innlogget bruker
- Admin

## Kommentarer
- Test nr. 1, "Hovednavigasjonen er stabil og lesbar", med "hovednavigasjonen" mener vi navigasjonsbaren med menyen. Dette inkluderer Hjem, Min Side, Admin, og Hamburgermenyen.

## Sjekkliste

### Alle roller

- [ ] **Hovednavigasjonen er stabil og lesbar**<br>
  **Gitt** at brukeren åpner sider med hovednavigasjon.<br>
  **Når** sidene er ferdig lastet.<br>
  **Så** skal navigasjonen være stabil, lesbar og uten brutte eller uferdige elementer.

- [ ] **Navigasjon og brukermeny passer på alle skjermstørrelser**<br>
  **Gitt** at brukeren bruker appen på mobil og større skjerm.<br>
  **Når** hovednavigasjonen og brukermenyen vises.<br>
  **Så** skal lenker, knapper og menyer være lesbare og ikke overlappe eller havne utenfor skjermen.

- [ ] **Fokus og alternative navigasjonsformer er tydelige**<br>
  **Gitt** at brukeren navigerer med tastatur eller andre alternative navigasjonsformer.<br>
  **Når** fokus flyttes i hovednavigasjon og brukermeny.<br>
  **Så** skal det være tydelig hvor fokus er og hvilke handlinger som kan utføres.

- [ ] **Raske sidebytter skaper ikke ustabilitet**<br>
  **Gitt** at brukeren klikker raskt mellom tilgjengelige sider.<br>
  **Når** flere navigasjonshandlinger skjer tett etter hverandre.<br>
  **Så** skal appen ikke havne i duplikathandlinger, feilnavigasjon eller tydelig ustabil tilstand.

- [ ] **Navigasjonen fremstår ferdig og konsistent**<br>
  **Gitt** at brukeren ser navigasjonen på tvers av sider.<br>
  **Når** appen brukes som helhet.<br>
  **Så** skal navigasjonen fremstå som ferdig og konsistent uten placeholder-preg, utilsiktet språkblanding eller visuelt forstyrrende detaljer.

### Ikke innlogget bruker

- [ ] **Innloggingsinngangen åpner riktig flyt**<br>
  **Gitt** at brukeren ikke er innlogget.<br>
  **Når** innloggingsinngangen brukes fra hovednavigasjonen.<br>
  **Så** skal brukeren komme til innloggingsflyten uten uventede feil eller feil side.

### Innlogget bruker

- [ ] **Utlogging endrer appen til utlogget tilstand**<br>
  **Gitt** at brukeren er innlogget.<br>
  **Når** brukeren logger ut fra brukermenyen.<br>
  **Så** skal appen tydelig gå over til utlogget tilstand.

- [ ] **Tidligere beskyttede sider forblir avvist etter utlogging**<br>
  **Gitt** at brukeren nylig har logget ut.<br>
  **Når** brukeren refresher eller åpner en tidligere beskyttet side.<br>
  **Så** skal appen fortsatt behandle brukeren som utlogget.

- [ ] **Eksterne menylenker markeres tydelig**<br>
  **Gitt** at eksterne lenker vises i brukermenyen.<br>
  **Når** brukeren åpner dem.<br>
  **Så** skal det være tydelig at innholdet ligger utenfor appen.

- [ ] **Ikke-admin får ingen halvveis adminvisning**<br>
  **Gitt** at en bruker uten adminrettigheter åpner en adminside direkte.<br>
  **Når** tilgang avvises.<br>
  **Så** skal brukeren ikke se en halvveis eller misvisende adminvisning.


### Authentisering

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

- [ ] **Passordtilbakestilling**<br>
  **Gitt** at en registrert bruker ønsker å tilbakestille passordet sitt.<br>
  **Når** brukreren oppgir et nytt passord.<br>
  **Så** skal tilbakemeldingen være forståelig og brukeren blir sendt til rett side.

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

#### Utlogging

- [ ] **Utlogging fjerner innlogget tilgang**<br>
  **Gitt** at en innlogget bruker velger å logge ut.<br>
  **Når** utloggingen fullføres.<br>
  **Så** skal brukeren ikke lenger fremstå som innlogget i menyen eller i tilgangen til beskyttede sider.

- [ ] **Refresh etter utlogging holder brukeren utlogget**<br>
  **Gitt** at en bruker nettopp har logget ut.<br>
  **Når** brukeren refresher siden eller åpner en tidligere beskyttet side.<br>
  **Så** skal appen fortsatt behandle brukeren som utlogget.

### Autorisering

- [ ] **Beskyttede sider forklarer avvist tilgang**<br>
  **Gitt** at en ikke-innlogget bruker åpner en beskyttet side direkte.<br>
  **Når** tilgang avvises.<br>
  **Så** skal brukeren få en tydelig forklaring og en forståelig vei videre til innlogging eller tilbake til appen.

- [ ] **Adminlenken åpner adminområdet riktig**<br>
  **Gitt** at brukeren er admin.<br>
  **Når** brukeren navigerer til Admin fra hovednavigasjonen.<br>
  **Så** skal brukeren bli sendt til adminområdet uten å møte feil rolle eller feil landingsside.
