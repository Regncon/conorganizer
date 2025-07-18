# Manuelle tester

## Sjekkliste for manuell testing

### index

- [x] **Last inn siden** : Test at siden laster uten feil.
- [x] **Nytt event** : Test at knappen for å opprette et nytt arrangement fungerer og leder til riktig side.
   * funker, men får noen ms med "hvis du ser denne meldingen..."


### Menu

- [x] **Navigasjon** : Test at alle menylenker fungerer og leder til riktige sider.
    * finner ikke noen lenke til profilen, men det funker å gå til /my-profile
- [x] **Responsiv design** : Test at menyen fungerer på forskjellige skjermstørrelser (mobil, nettbrett, desktop).
    * funker, men behøver vi hamburgerknapp på desktop?
- [x] **Hjem knapp** : Test at "Hjem" knappen leder tilbake til forsiden.
- [x] **Logg inn og ut** : Test at menyen byttes mellom "Logg inn" knappen når brukeren er logget ut.

### Authentisering og autorisasjon

- [x] **Registrering** : Test at en ny bruker kan registrere seg uten problemer.
    * Jeg får helt fint lov til å sette passordet mitt til "password", lurer på om vi kanskje bør legge listen høyere
- [x] **Innlogging** : Test at en eksisterende bruker kan logge inn med riktige legitimasjonsbeskrivelser.
- [x] **Log ut** : Test at brukeren kan logge ut uten problemer.
    * "you are not logged in" bør gjerne erstattes av "you are logged out"?
- [x] **Tilgangskontroll** : Test at brukere uten riktig tilgang ikke kan se beskyttede sider.
    * ikke med et uhell i alle fall, har ikke forsøkt å hacke (det er for varmt med hettegenser)
- [x] **Feilmeldinger** : Test at brukeren får riktige feilmeldinger ved feil innlogging eller registrering.
    * Når man er for sent ute med koden får man siden på nytt med feilmelding, men det står fortsatt "We've sent a message containing a 6-digit code...". Det er muligens litt tvetydig, for den er jo ikke sendt på nytt.
    * når jeg skriver "kattepus@motherfucker.ded" så blir jeg sendt videre til "We've sent a message containing a 6-digit
 code to kat*****@mot********.ded" - jeg *antar* at det betyr at den kommer uansett hva du skriver, det kan muligens bli en liten hodepine for folk som feilstaver eposten sin.
- [x] **Glemt passord** : Test at funksjonen for å tilbakestille passord fungerer som forventet.
- [x] **E-postbekreftelse** : Test at e-postbekreftelse fungerer for nye brukere.
    * Det er ikke åpenbart at en epost fra "loving-ochre-moore" er fra Regncon :p
- [ ] **Brukerprofil** : Test at brukeren kan se og redigere sin profil.
    * ser ingen respons når jeg trykker "koble til epost" eller bossknappen på "Kobledt epost 1", og "koblet" staves uten "d" (og burde kanskje være "tilkoblet"?)
    * symbolet på bossknappen funker ikke i så liten størrelse
    * de gråtende bildene på "Min profil"-siden er gjerne placeholders og bør erstattes
- [ ] **Administrator** : Test at administratorer har tilgang til administrative funksjoner som brukerhåndtering.
    * christer vet ikke hvordan han logger inn som admin

### my-events

- [x] **Liste over arrangementer** : Test at listen over brukerens arrangementer vises korrekt.
- [x] **Opprett arrangement** : Test at en bruker kan opprette et nytt arrangement.
    * burde "Om arrangøren" heller være en dropdown av billettholdere tilknyttet kontoen? Har vi kanskje diskutert
    dette allerede?
    * grådig slitsomt med den popupen når jeg forsøker å skrive all-caps med shift
    * navn og telefonnummer ser ut til å forsvinne av og til? Men så fyller jeg de inn på nytt, så funker det :l
    * placeholder-teksten forsvinner ikke når jeg klikker i boksen, har man brukt "value" i stedet for "placeholder"?
- [x] **Rediger arrangement** : Test at en bruker kan redigere et eksisterende arrangement.
    * navn og telefonnummer forsvinner når jeg skal redigere
    * når jeg testet punktet under driver også navnet, tlf og eposten å forsvinner - på ett tidspunkt ble telefonnummeret erstattet med "9", og eposten ble nettopp erstattet med "c" når jeg fylte inn navnet
    * ser det står "nytt arrangement" som overskrift når jeg redigerer ett arrangement som alt er sendt inn?
- [ ] **Påkrvde felt** : Test at alle påkrevde felt er riktig validert ved opprettelse og redigering av arrangement.
    * kan sette telefonnummeret som "9" og navnet som "c"
- [x] **Rettigheter** : Test at brukeren kun kan redigere eller slette sine egne arrangementer.
    * får forskjellig feilmelding når jeg går inn på urlen til et arrangement når jeg er logget inn som en annen bruker (ustylet side med "event not found") vs når jeg ikke er logget inn: fin stailet side med "du har ikke tillatelse" elns. Samme med admin-siden. Det er gjerne meningen, men nevner det nå likevel.
### None happy path

- [ ] **Opprettelse feiler** : Sjekk at brykeren får fornuftig feilmelding ved feil under opprettelse av arrangement.
    * testet å rename databasefilen når jeg opprettet arrangement, fikk "Failed to update the status for event in the database" - så det funker i alle fall som det skal.
    * OI! Jeg får ingen feilmelding når jeg redigerer et arrangement og databasen er renamet! Jeg trykker "send inn" og den redirekter meg stille og pent til arrangementsiden, og ingenting er lagret!
