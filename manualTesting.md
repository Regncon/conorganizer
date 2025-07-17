# Manuelle tester

## Sjekkliste for manuell testing

### index

- [ ] **Last inn siden** : Test at siden laster uten feil.
- [ ] **Nytt event** : Test at knappen for å opprette et nytt arrangement fungerer og leder til riktig side.

### Menu

- [ ] **Navigasjon** : Test at alle menylenker fungerer og leder til riktige sider.
- [ ] **Responsiv design** : Test at menyen fungerer på forskjellige skjermstørrelser (mobil, nettbrett, desktop).
- [ ] **Hjem knapp** : Test at "Hjem" knappen leder tilbake til forsiden.
- [ ] **Logg inn og ut** : Test at menyen byttes mellom "Logg inn" knappen når brukeren er logget ut.

### Authentisering og autorisasjon

- [ ] **Registrering** : Test at en ny bruker kan registrere seg uten problemer.
- [ ] **Innlogging** : Test at en eksisterende bruker kan logge inn med riktige legitimasjonsbeskrivelser.
- [ ] **Log ut** : Test at brukeren kan logge ut uten problemer.
- [ ] **Tilgangskontroll** : Test at brukere uten riktig tilgang ikke kan se beskyttede sider.
- [ ] **Feilmeldinger** : Test at brukeren får riktige feilmeldinger ved feil innlogging eller registrering.
- [ ] **Glemt passord** : Test at funksjonen for å tilbakestille passord fungerer som forventet.
- [ ] **E-postbekreftelse** : Test at e-postbekreftelse fungerer for nye brukere.
- [ ] **Brukerprofil** : Test at brukeren kan se og redigere sin profil.
- [ ] **Administrator** : Test at administratorer har tilgang til administrative funksjoner som brukerhåndtering.

### my-events

- [ ] **Liste over arrangementer** : Test at listen over brukerens arrangementer vises korrekt.
- [ ] **Opprett arrangement** : Test at en bruker kan opprette et nytt arrangement.
- [ ] **Rediger arrangement** : Test at en bruker kan redigere et eksisterende arrangement.
- [ ] **Påkrvde felt** : Test at alle påkrevde felt er riktig validert ved opprettelse og redigering av arrangement.
- [ ] **Rettigheter** : Test at brukeren kun kan redigere eller slette sine egne arrangementer.

### None happy path

- [ ] **Opprettelse feiler** : Sjekk at brykeren får fornuftig feilmelding ved feil under opprettelse av arrangement.
