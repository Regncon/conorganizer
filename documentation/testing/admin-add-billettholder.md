# Legg til billettholder i admin

Denne sjekklisten dekker `/admin/billettholder/add`, der admin kan se billetter fra CheckIn og konvertere relevante billetter til billettholdere.

## Roller

- Admin

## Sjekkliste

### Oversikt og konvertering

- [ ] **Billettoversikten laster uten brutte kort**<br>
  **Gitt** at en admin åpner siden for å legge til billettholder.<br>
  **Når** siden lastes.<br>
  **Så** skal hvert billettkort vise bestilling, type, navn, e-post og alder uten brutte kort eller uforståelige feilmeldinger.

- [ ] **Konvertering oppretter billettholdere for hele bestillingen**<br>
  **Gitt** at en billett kan konverteres, og bestillingen har flere billetter.<br>
  **Når** admin velger «Konverter bestilling til deltagere» på én av billettene.<br>
  **Så** skal alle billettene på bestillingen, unntatt middag, bli billettholdere, og alle kortene på bestillingen skal deretter vise at billetten allerede er konvertert, uten at admin må gjette om handlingen faktisk lyktes.

- [ ] **Brukere på bestillingen får billettholderne**<br>
  **Gitt** at en bruker med en e-post på bestillingen allerede har logget inn.<br>
  **Når** admin konverterer bestillingen.<br>
  **Så** skal brukeren se billettholderne uten å trykke «Hent billetter».

- [ ] **Middagsbilletter og konverterte billetter kan ikke konverteres**<br>
  **Gitt** at en billett er en middagsbillett eller allerede er konvertert til billettholder.<br>
  **Når** kortet vises.<br>
  **Så** skal kortet vise en merknad om dette og ikke tilby knappen for å konvertere.

- [ ] **Konverteringsfeil forklares tydelig**<br>
  **Gitt** at konvertering av billett feiler.<br>
  **Når** admin forsøker å konvertere.<br>
  **Så** skal admin få en tydelig feiltilstand og ikke stå igjen med en side som ser oppdatert ut uten at endringen faktisk ble gjort.

- [ ] **Flere konverteringer oppdaterer riktige kort**<br>
  **Gitt** at admin konverterer flere billetter etter hverandre.<br>
  **Når** siden oppdateres fortløpende.<br>
  **Så** skal riktig status vises på riktige kort og ikke blandes mellom billetter.

### CheckIn og datamengder

- [ ] **Utilgjengelig CheckIn gir forståelig melding**<br>
  **Gitt** at CheckIn ikke svarer.<br>
  **Når** admin åpner siden.<br>
  **Så** skal siden si at den ikke får kontakt med CheckIn, eller at den viser sist lagrede informasjon når den har det, og ikke gi inntrykk av at billettene har forsvunnet.

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
  **Gitt** at admin refresher siden etter konvertering.<br>
  **Når** siden lastes inn igjen.<br>
  **Så** skal innholdet samsvare med faktisk lagret tilstand og ikke med en foreldet mellomtilstand.
