# Billetter på Min Side

Denne sjekklisten dekker `/profile/tickets`, der innlogget bruker kan hente billetter, se billettholdere og legge til eller slette manuelle e-postadresser.

## Roller

- Innlogget bruker

## Sjekkliste

### Hente billetter

- [ ] **Tom side tilbyr å hente eller kjøpe billetter**<br>
  **Gitt** at brukeren ikke har billettholdere knyttet til seg.<br>
  **Når** `/profile/tickets` vises.<br>
  **Så** vises «Ingen billetter funnet på denne e-postadressen» med knappene «Hent billetter» og «Kjøp billetter». «Kjøp billetter» åpner Checkin i ny fane.

- [ ] **Hent billetter viser hentede billettholdere**<br>
  **Gitt** at brukeren trykker på Hent billetter.<br>
  **Når** billettene faktisk kan hentes.<br>
  **Så** skal billettholdere dukke opp uten at siden havner i en stille eller uavklart tilstand, og meldingen skal være «1 billett hentet!» eller «N billetter hentet!».

- [ ] **Ingen nye billetter gir nøytral melding**<br>
  **Gitt** at brukeren allerede har hentet billettene sine eller ikke har noen billetter på e-postadressen.<br>
  **Når** brukeren trykker på Hent billetter.<br>
  **Så** vises «Ingen nye billetter funnet.» og ingen feilmelding.

- [ ] **Henting har tydelig ventetilstand**<br>
  **Gitt** at brukeren trykker på Hent billetter.<br>
  **Når** henting pågår.<br>
  **Så** skal knappen vise en lasteindikator og være deaktivert, slik at det er tydelig at en handling er i gang.

- [ ] **Hentefeil gir tydelig feilmelding**<br>
  **Gitt** at henting av billetter feiler.<br>
  **Når** brukeren forsøker å hente billetter.<br>
  **Så** skal brukeren få en tydelig feilmelding og ikke en falsk bekreftelse på at alt gikk bra. Når Checkin ikke svarer og ingen tidligere data finnes, sier meldingen at vi ikke får kontakt med Checkin. Når bare oppdateringen feiler, vises en advarsel om at sist lagrede informasjon brukes.

- [ ] **Raske henteklikk skaper ikke duplikater**<br>
  **Gitt** at brukeren trykker på Hent billetter flere ganger raskt.<br>
  **Når** siden håndterer forespørslene.<br>
  **Så** skal det ikke oppstå duplisering eller åpenbart ustabil oppførsel i resultatet.

### E-postadresser

- [ ] **Manuell e-post kan legges til og slettes**<br>
  **Gitt** at brukeren skriver en ny e-postadresse på et billettholderkort og trykker «Legg til epost» eller Enter.<br>
  **Når** adressen lagres.<br>
  **Så** vises «Epostadressen … er lagt til» på samme kort, og adressen får en «Slett»-knapp. Bare manuelle adresser kan slettes; billett-e-posten og andre adresser fra samme bestilling har ingen slettehandling.

- [ ] **Add- og delete-feil er tydelige**<br>
  **Gitt** at en add- eller delete-handling feiler, for eksempel med tomt felt eller en adresse som allerede finnes på billettholderen.<br>
  **Når** brukeren utfører endringen.<br>
  **Så** skal feilmeldingen vises på riktig kort, være tydelig («Tomt felt for epostadresse», «Epostadressen … finnes allerede for denne bilettholderen») og ikke etterlate inntrykk av at endringen likevel ble lagret.

### Meldinger og responsivitet

- [ ] **Billettsiden fungerer på mobil**<br>
  **Gitt** at brukeren bruker billettsiden på mobil.<br>
  **Når** mange billettholderkort eller lange e-postadresser vises.<br>
  **Så** skal innholdet fortsatt være lesbart og brukbart uten overlapp eller horisontal kollaps.
