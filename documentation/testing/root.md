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

- [ ] **«Send inn arrangement» leder til opprettelse via Min Side**<br>
  **Gitt** at en innlogget bruker trykker «Send inn arrangement» på forsiden.<br>
  **Når** brukeren kommer til Min Side og trykker «Send inn arrangement» under Mine arrangementer.<br>
  **Så** skal et nytt arrangement opprettes som kladd, og brukeren sendes til skjemaet for det nye arrangementet.

- [ ] **Ikke-innlogget bruker får en vei videre fra «Send inn arrangement»**<br>
  **Gitt** at en ikke-innlogget bruker trykker «Send inn arrangement» på forsiden.<br>
  **Når** Min Side avviser brukeren.<br>
  **Så** skal brukeren få en tydelig beskjed om å logge inn og en vei videre til innlogging eller tilbake til arrangementslisten. Etter innlogging via «Logg inn» skal brukeren lande på Min Side.

### Program og arrangementskort

- [ ] **Upublisert program viser bare annonserte arrangementer**<br>
  **Gitt** at programmet ikke er publisert.<br>
  **Når** brukeren åpner forsiden.<br>
  **Så** skal forsiden vise en flat liste med annonserte arrangementer uten dagvelger og uten puljeinndeling.

- [ ] **Publisert program viser valgt dag**<br>
  **Gitt** at programmet er publisert.<br>
  **Når** brukeren åpner forsiden.<br>
  **Så** skal dagvelgeren vises med aktiv dag markert, og dagens tittel og tidsskjema vises før programoversikten og puljene.

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
  **Så** skal tittel, ingress og ikoner, og der kortet viser dem, system og arrangør, fremstå lesbare og høre til riktig arrangement.

- [ ] **Arrangementskort åpner riktig detaljside**<br>
  **Gitt** at et arrangementskort vises på forsiden.<br>
  **Når** brukeren trykker på kortet.<br>
  **Så** skal brukeren sendes til riktig arrangementside. Når programmet er publisert, skal valgt dag og pulje følge med som kontekst.

### Navigasjon og robusthet

- [ ] **Dagvelgeren bytter til riktig dag**<br>
  **Gitt** at programmet er publisert og har flere dager.<br>
  **Når** brukeren velger en annen dag i dagvelgeren.<br>
  **Så** skal forsiden vise program og puljer for valgt dag, dagen skal være markert som aktiv, og den sticky dagvelgeren skal ikke skjule overskrifter eller viktig informasjon.

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
