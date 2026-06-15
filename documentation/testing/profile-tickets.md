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

- [ ] **Add- og delete-feil er tydelige**<br>
  **Gitt** at en add- eller delete-handling feiler.<br>
  **Når** admin utfører endringen.<br>
  **Så** skal feilmeldingen være tydelig og ikke etterlate inntrykk av at endringen likevel ble lagret.

### Meldinger og responsivitet

- [ ] **Billettsiden fungerer på mobil**<br>
  **Gitt** at brukeren bruker billettsiden på mobil.<br>
  **Når** mange billettholderkort eller lange e-postadresser vises.<br>
  **Så** skal innholdet fortsatt være lesbart og brukbart uten overlapp eller horisontal kollaps.
