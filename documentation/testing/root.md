# Forside

Denne sjekklisten dekker forsiden på `/`. Forsiden er en sentral inngang til appen og skal fungere for både ikke-innlogget bruker, innlogget bruker og admin.

## Roller

- Ikke-innlogget bruker
- Innlogget bruker
- Admin

## Sjekkliste

### Førsteinntrykk og layout

- [ ] **Brødsmulestien viser Hjem**<br>
  **Gitt** at brukeren åpner forsiden.<br>
  **Når** siden er ferdig lastet.<br>
  **Så** skal brødsmulestien vise at brukeren er på Hjem.

- [ ] **Innsendingsseksjonen inviterer tydelig til registrering**<br>
  **Gitt** at brukeren åpner forsiden.<br>
  **Når** seksjonen for å sende inn arrangement vises.<br>
  **Så** skal innholdet være forståelig, lesbart og fremstå som en tydelig invitasjon til å registrere arrangement.

- [ ] **Innsendingsseksjonen fungerer på liten skjerm**<br>
  **Gitt** at brukeren åpner forsiden på en liten skjerm.<br>
  **Når** seksjonen for å sende inn arrangement vises.<br>
  **Så** skal tekst, knapp og illustrasjon være lesbare og ikke presse hverandre ut av layouten.

- [ ] **Innsendingsseksjonen er balansert på større skjerm**<br>
  **Gitt** at brukeren åpner forsiden på en større skjerm.<br>
  **Når** seksjonen for å sende inn arrangement vises.<br>
  **Så** skal tekst, knapp og illustrasjon være balansert og uten tomrom eller skjevheter som får innholdet til å se ødelagt ut.

### Program og arrangementskort

- [ ] **Puljer vises med riktig navn og tidspunkt**<br>
  **Gitt** at det finnes publiserte arrangementer i én eller flere puljer.<br>
  **Når** brukeren åpner forsiden etter at programmet er publisert.<br>
  **Så** skal hver pulje vises med korrekt navn og tidspunkt.

- [ ] **Arrangementer ligger under riktig pulje**<br>
  **Gitt** at det finnes publiserte arrangementer i flere puljer.<br>
  **Når** brukeren åpner forsiden etter at programmet er publisert.<br>
  **Så** skal arrangementene vises under riktig pulje og ikke lekke over i feil seksjon.

- [ ] **Arrangementskort viser riktig lesbar informasjon**<br>
  **Gitt** at forsiden viser arrangementskort.<br>
  **Når** brukeren leser kortene.<br>
  **Så** skal tittel, ingress, arrangør, system og ikoner fremstå lesbare og høre til riktig arrangement.

- [ ] **Arrangementskort åpner riktig detaljside**<br>
  **Gitt** at et arrangementskort vises på forsiden.<br>
  **Når** brukeren trykker på kortet.<br>
  **Så** skal brukeren sendes til riktig arrangementside og beholde riktig kontekst for valgt pulje.

### Navigasjon og robusthet

- [ ] **Snarveier scroller til riktig pulje**<br>
  **Gitt** at brukeren navigerer mellom puljene via snarveinavigasjonen på forsiden.<br>
  **Når** brukeren trykker på en pulje.<br>
  **Så** skal siden scrolle til riktig seksjon uten å havne merkbart feil eller skjule seksjonsoverskriften bak sticky navigasjon.

- [ ] **Snarveinavigasjonen dekker ikke viktig innhold**<br>
  **Gitt** at brukeren scroller på forsiden.<br>
  **Når** snarveinavigasjonen er synlig.<br>
  **Så** skal den oppføre seg stabilt og ikke dekke viktig innhold på en måte som gjør siden vanskelig å bruke.

- [ ] **Valgt pulje er tydelig etter navigasjon**<br>
  **Gitt** at brukeren går direkte til en pulje via snarveinavigasjonen.<br>
  **Når** seksjonen blir synlig.<br>
  **Så** skal det være tydelig hvilken pulje brukeren har navigert til.

- [ ] **Tilbakeknapp bevarer brukbar forside**<br>
  **Gitt** at brukeren bruker tilbakeknappen etter å ha åpnet et arrangement fra forsiden.<br>
  **Når** brukeren kommer tilbake.<br>
  **Så** skal forsiden fortsatt være brukbar og ikke miste viktige deler av tilstanden sin.

- [ ] **Refresh viser forsiden korrekt**<br>
  **Gitt** at brukeren refresher forsiden.<br>
  **Når** siden lastes på nytt.<br>
  **Så** skal innhold og forsideseksjonene fortsatt vises korrekt uten at brukeren havner i en uforståelig tilstand.

- [ ] **Feiltilstand er brukervennlig**<br>
  **Gitt** at forsiden ikke kan laste innhold eller arrangementsdata som forventet.<br>
  **Når** siden viser feiltilstand.<br>
  **Så** skal feilen være brukervennlig og ikke vise tekniske detaljer.

- [ ] **Raske klikk skaper ikke feilnavigasjon**<br>
  **Gitt** at forsiden brukes over tid med flere raske klikk på navigasjon og kort.<br>
  **Når** brukeren forflytter seg mellom sider.<br>
  **Så** skal det ikke oppstå åpenbare duplikathandlinger, feilnavigasjon eller ustabil oppførsel.

- [ ] **Store datamengder beholder lesbar struktur**<br>
  **Gitt** at forsiden vises med ekte eller store datamengder.<br>
  **Når** mange arrangementer finnes i samme eller flere puljer.<br>
  **Så** skal siden fortsatt være lesbar, navigerbar og uten tydelige sammenbrudd i layout eller informasjonsstruktur.
