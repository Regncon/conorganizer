# Admin

Denne sjekklisten dekker hovedsiden for admin på `/admin`, altså inngangen til administrative funksjoner.

## Roller

- Bruker uten adminrettigheter
- Admin

## Sjekkliste

### Tilgang

- [ ] **Ikke-admin får tydelig avvisning**<br>
  **Gitt** at en innlogget bruker uten adminrettigheter åpner `/admin` direkte.<br>
  **Når** siden lastes.<br>
  **Så** skal brukeren få siden «Du har ikke tilgang» med en knapp tilbake til arrangementslisten, og ingen adminkort skal vises.

### Hovedvalg og navigasjon

- [ ] **Adminforsiden viser hovedvalg**<br>
  **Gitt** at en admin åpner adminforsiden.<br>
  **Når** siden lastes ferdig etter liveoppdatering.<br>
  **Så** skal kortene Publiser program, Arrangementer til godkjenning, Billettholderoversikt, Puljefordeling, Puljeoppsett, Administrer rom og Utskriftsvennlig program vises uten brutte paneler eller feil rolleopplevelse.

- [ ] **Adminvalg åpner riktig underside**<br>
  **Gitt** at en admin velger å gå til et underliggende adminområde.<br>
  **Når** navigasjonen skjer.<br>
  **Så** skal riktig underside åpnes uten feil rolle eller uventet mellomtilstand.

- [ ] **Publisering av program krever bekreftelse**<br>
  **Gitt** at en admin bruker bryteren på kortet Publiser program.<br>
  **Når** admin først avbryter bekreftelsen og deretter bekrefter.<br>
  **Så** skal statusen være uendret etter avbrudd, og etter bekreftelse vise «Publisert» eller «Ikke publisert» i samsvar med lagret tilstand, også etter refresh.

- [ ] **Adminkort fungerer på alle skjermstørrelser**<br>
  **Gitt** at adminforsiden brukes på mobil og større skjerm.<br>
  **Når** kortene vises.<br>
  **Så** skal de være lesbare, klikkbare og visuelt stabile uten at tekst eller bilder kolliderer.

### Robusthet

- [ ] **Tilbakeknapp og refresh bevarer adminkontekst**<br>
  **Gitt** at admin går frem og tilbake mellom adminforsiden og underliggende adminsider.<br>
  **Når** brukeren bruker tilbakeknapp og refresh.<br>
  **Så** skal adminområdet fortsatt oppføre seg konsistent og tydelig som adminområde.
