# Billetter på Min Side

Denne sjekklisten dekker `/profile/tickets`, der innlogget bruker kan hente billetter, se billettholdere og legge til eller slette manuelle e-postadresser.

## Roller

- Innlogget bruker

## Sjekkliste

### Hente billetter

- [ ] **Hent billetter viser hentede billettholdere**<br>
  **Gitt** at brukeren trykker på Hent billetter.<br>
  **Når** billettene faktisk kan hentes.<br>
  **Så** skal billettholdere dukke opp uten at siden havner i en stille eller uavklart tilstand.

- [ ] **Henting har tydelig ventetilstand**<br>
  **Gitt** at brukeren trykker på Hent billetter.<br>
  **Når** henting pågår.<br>
  **Så** skal knappen og lasteindikatoren oppføre seg på en måte som gjør det tydelig at en handling er i gang.

- [ ] **Hentefeil gir tydelig feilmelding**<br>
  **Gitt** at henting av billetter feiler.<br>
  **Når** brukeren forsøker å hente billetter.<br>
  **Så** skal brukeren få en tydelig feilmelding og ikke en falsk bekreftelse på at alt gikk bra.

- [ ] **Raske henteklikk skaper ikke duplikater**<br>
  **Gitt** at brukeren trykker på Hent billetter flere ganger raskt.<br>
  **Når** siden håndterer forespørslene.<br>
  **Så** skal det ikke oppstå duplisering eller åpenbart ustabil oppførsel i resultatet.

### E-postadresser

- [ ] **Ny e-postadresse legges til riktig billettholder**<br>
  **Gitt** at brukeren legger til en ny manuell e-postadresse på en billettholder.<br>
  **Når** handlingen lykkes.<br>
  **Så** skal brukeren få en tydelig bekreftelse og se at e-postadressen faktisk er lagt til riktig billettholder.

- [ ] **Tom e-postadresse avvises**<br>
  **Gitt** at brukeren forsøker å legge til en tom e-postadresse.<br>
  **Når** handlingen sendes inn.<br>
  **Så** skal brukeren få en tydelig feilmelding og ingen ny e-postadresse skal legges til.

- [ ] **Duplikatadresse avvises**<br>
  **Gitt** at brukeren forsøker å legge til en e-postadresse som allerede finnes på samme billettholder.<br>
  **Når** handlingen sendes inn.<br>
  **Så** skal brukeren få en forståelig feilmelding og ingen duplikatadresse skal opprettes.

- [ ] **Oppdatering bevarer eksisterende kortdata**<br>
  **Gitt** at brukeren legger til en ny e-postadresse.<br>
  **Når** siden oppdateres.<br>
  **Så** skal tidligere data på siden fortsatt være intakte og ikke forsvinne eller byttes om mellom kortene.

- [ ] **Sletting fjerner riktig e-postadresse**<br>
  **Gitt** at en manuell e-postadresse finnes på en billettholder.<br>
  **Når** brukeren velger å slette den.<br>
  **Så** skal den slettes fra riktig billettholder og brukeren skal få en tydelig bekreftelse.

- [ ] **Slettefeil viser ufullført endring**<br>
  **Gitt** at sletting av e-postadresse feiler.<br>
  **Når** brukeren forsøker å slette.<br>
  **Så** skal brukeren få en feilmelding som gjør det tydelig at endringen ikke ble fullført.

### Meldinger og responsivitet

- [ ] **Meldinger hører til riktig billettholder**<br>
  **Gitt** at siden viser meldinger om vellykket eller mislykket endring.<br>
  **Når** flere handlinger utføres etter hverandre.<br>
  **Så** skal meldingene høre til riktig billettholder og ikke lekke over til andre kort.

- [ ] **Billettsiden fungerer på mobil**<br>
  **Gitt** at brukeren bruker billettsiden på mobil.<br>
  **Når** mange billettholderkort eller lange e-postadresser vises.<br>
  **Så** skal innholdet fortsatt være lesbart og brukbart uten overlapp eller horisontal kollaps.
