# Admin

Denne sjekklisten dekker hovedsiden for admin på `/admin`, altså inngangen til administrative funksjoner.

## Roller

- Bruker uten adminrettigheter
- Admin

## Sjekkliste

### Hovedvalg og navigasjon

- [ ] **Adminforsiden viser hovedvalg**<br>
  **Gitt** at en admin åpner adminforsiden.<br>
  **Når** siden lastes ferdig etter liveoppdatering.<br>
  **Så** skal adminområdets hovedvalg vises uten brutte paneler eller feil rolleopplevelse.

- [ ] **Adminvalg åpner riktig underside**<br>
  **Gitt** at en admin velger å gå til et underliggende adminområde.<br>
  **Når** navigasjonen skjer.<br>
  **Så** skal riktig underside åpnes uten feil rolle eller uventet mellomtilstand.

- [ ] **Adminkort fungerer på alle skjermstørrelser**<br>
  **Gitt** at adminforsiden brukes på mobil og større skjerm.<br>
  **Når** kortene vises.<br>
  **Så** skal de være lesbare, klikkbare og visuelt stabile uten at tekst eller bilder kolliderer.

### Tilgang og robusthet

- [ ] **Ikke-admin avvises tydelig**<br>
  **Gitt** at en bruker uten adminrettigheter prøver å åpne adminforsiden direkte.<br>
  **Når** siden lastes.<br>
  **Så** skal brukeren ikke få tilgang og heller ikke møte en misvisende halvveis adminvisning.

- [ ] **Lastingsfeil forklares**<br>
  **Gitt** at adminforsiden ikke kan laste nødvendig innhold som forventet.<br>
  **Når** siden vises.<br>
  **Så** skal brukeren ikke bli stående med en tilsynelatende tom adminside uten forklaring.

- [ ] **Tilbakeknapp og refresh bevarer adminkontekst**<br>
  **Gitt** at admin går frem og tilbake mellom adminforsiden og underliggende adminsider.<br>
  **Når** brukeren bruker tilbakeknapp og refresh.<br>
  **Så** skal adminområdet fortsatt oppføre seg konsistent og tydelig som adminområde.
