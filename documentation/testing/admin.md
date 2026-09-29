# Admin

Denne sjekklisten dekker hovedsiden for admin på `/admin`, altså inngangen til administrative funksjoner. Siden har overskriften «Adminverktøy» og er delt i tre seksjoner med verktøykort.

## Manuelt og automatisert

Innhold og oppførsel på adminforsiden er automatisert i Go-tester:

- `pages/admin/admin_page_test.go`: brødsmulestien viser Admin.
- `pages/admin/admin_landing_integration_test.go`: overskrift, seksjoner og kort i riktig rekkefølge med lenke, tittel, beskrivelse og grafikk; Publiser program-kortet er et bredt kort uten lenke med utskriftslenke og publiseringsbryter; alle kortlenker peker til registrerte ruter; grafikkfilene finnes og er gyldige SVG-er; publiseringsbryteren lagrer tilstanden og viser den ved neste visning; en ikke-admin kan ikke endre publiseringen.

Bekreftelsesdialogen for publisering, layouten på ulike skjermstørrelser og navigasjon med tilbakeknapp og refresh testes manuelt her.

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

- [ ] **Adminforsiden viser seksjoner og verktøykort**<br>
  **Gitt** at en admin åpner adminforsiden.<br>
  **Når** siden lastes ferdig etter liveoppdatering.<br>
  **Så** skal overskriften «Adminverktøy» vises med seksjonene i denne rekkefølgen: Puljefordeling med kortene «Fordeling av deltakere» og «Sett arrangementer i puljer», Programoppsett med «Publiser program», «Administrer rom» og «Godkjenn arrangementer», og Øvrige verktøy med «Administrer billettholdere» og «Tilbakemeldinger». Hvert kort skal ha grafikk, tittel og beskrivelse uten brutte bilder eller paneler.

- [ ] **Publiser program-kortet samler publisering og utskrift**<br>
  **Gitt** at en admin ser kortet Publiser program.<br>
  **Når** admin leser og bruker kortet.<br>
  **Så** skal kortet gå over hele raden, selve kortet skal ikke være en lenke, og det skal inneholde publiseringsbryteren og lenken «Hent utskriftsvennlig versjon av programmet», som åpner `/print`.

- [ ] **Adminvalg åpner riktig underside**<br>
  **Gitt** at en admin trykker på et verktøykort.<br>
  **Når** navigasjonen skjer.<br>
  **Så** skal riktig underside åpnes uten feil rolle eller uventet mellomtilstand: `/admin/puljefordeling/`, `/admin/puljeoppsett/`, `/admin/rooms/`, `/admin/approval/`, `/admin/billettholder/` eller `/admin/tilbakemeldinger/`.

- [ ] **Publisering av program krever bekreftelse**<br>
  **Gitt** at en admin bruker bryteren på kortet Publiser program.<br>
  **Når** admin først avbryter bekreftelsen og deretter bekrefter.<br>
  **Så** skal statusen være uendret etter avbrudd, og etter bekreftelse vise «Publisert» eller «Ikke publisert» i samsvar med lagret tilstand, også etter refresh.

- [ ] **Adminkort fungerer på alle skjermstørrelser**<br>
  **Gitt** at adminforsiden brukes på mobil og større skjerm.<br>
  **Når** kortene vises.<br>
  **Så** skal kortene ligge i én kolonne på smal skjerm og to kolonner på bredere skjerm, og være lesbare, klikkbare og visuelt stabile uten at tekst eller bilder kolliderer.

### Robusthet

- [ ] **Tilbakeknapp og refresh bevarer adminkontekst**<br>
  **Gitt** at admin går frem og tilbake mellom adminforsiden og underliggende adminsider.<br>
  **Når** brukeren bruker tilbakeknapp og refresh.<br>
  **Så** skal adminområdet fortsatt oppføre seg konsistent og tydelig som adminområde.
