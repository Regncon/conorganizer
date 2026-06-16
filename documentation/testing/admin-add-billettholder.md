# Legg til billettholder i admin

Denne sjekklisten dekker `/admin/billettholder/add`, der admin kan se billetter fra check-in og konvertere relevante billetter til billettholdere.

## Roller

- Admin

## Sjekkliste

### Oversikt og konvertering

- [ ] **Billettoversikten laster uten brutte kort**<br>
  **Gitt** at en admin åpner siden for å legge til billettholder.<br>
  **Når** siden lastes.<br>
  **Så** skal oversikten over billetter vises uten brutte kort eller uforståelige feilmeldinger.

- [ ] **Konvertering oppretter riktig billettholder**<br>
  **Gitt** at en billett kan konverteres.<br>
  **Når** admin velger å konvertere den.<br>
  **Så** skal billetten bli til riktig billettholder uten at admin må gjette om handlingen faktisk lyktes.

- [ ] **Konverteringsfeil forklares tydelig**<br>
  **Gitt** at konvertering av billett feiler.<br>
  **Når** admin forsøker å konvertere.<br>
  **Så** skal admin få en tydelig feiltilstand og ikke stå igjen med en side som ser oppdatert ut uten at endringen faktisk ble gjort.

- [ ] **Flere konverteringer oppdaterer riktige kort**<br>
  **Gitt** at admin konverterer flere billetter etter hverandre.<br>
  **Når** siden oppdateres fortløpende.<br>
  **Så** skal riktig status vises på riktige kort og ikke blandes mellom billetter.

### Søk og datamengder

- [ ] **Søk viser forståelige resultater**<br>
  **Gitt** at admin bruker søk eller filtrering på siden.<br>
  **Når** relevante treff vises.<br>
  **Så** skal resultatene være forståelige og markeringen av søket ikke gjøre innholdet uleselig.

- [ ] **Tomt søk gir stabil opplevelse**<br>
  **Gitt** at admin bruker søk eller filtrering med tomt eller lite nyttig søk.<br>
  **Når** siden oppdateres.<br>
  **Så** skal brukeropplevelsen fortsatt være stabil og ikke gi inntrykk av at data har forsvunnet.

- [ ] **Mange billetter forblir lesbare**<br>
  **Gitt** at admin bruker siden med mange billetter og varierende data.<br>
  **Når** oversikten vises.<br>
  **Så** skal kortene fortsatt være lesbare og ikke bryte grid eller flyt.

### Navigasjon og refresh

- [ ] **Konvertert billettholder vises konsistent videre**<br>
  **Gitt** at admin navigerer tilbake til billettholderoversikten etter konvertering.<br>
  **Når** oversikten åpnes.<br>
  **Så** skal den nye billettholderen være håndtert konsistent med resten av systemet.

- [ ] **Siden fungerer på mobil**<br>
  **Gitt** at siden brukes på mobil eller smal skjerm.<br>
  **Når** lange navn, e-poster eller mange kort vises.<br>
  **Så** skal siden fortsatt være brukbar og ikke falle visuelt sammen.

- [ ] **Refresh viser lagret tilstand**<br>
  **Gitt** at admin refresher siden etter konvertering eller søk.<br>
  **Når** siden lastes inn igjen.<br>
  **Så** skal innholdet samsvare med faktisk lagret tilstand og ikke med en foreldet mellomtilstand.
