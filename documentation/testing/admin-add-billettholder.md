# Legg til billettholder i admin

Denne sjekklisten dekker `/admin/billettholder/add`, der admin kan se billetter fra check-in og konvertere relevante billetter til billettholdere.

## Roller

- Admin

## Sjekkliste

- [ ] `Gitt at en admin åpner siden for å legge til billettholder, når siden lastes, så skal oversikten over billetter vises uten brutte kort eller uforståelige feilmeldinger.`
- [ ] `Gitt at billetter vises på siden, når admin leser kortene, så skal bestilling, type, navn, e-post og alder være tydelige og høre til riktig billett.`
- [ ] `Gitt at en billett allerede er konvertert til billettholder, når kortet vises, så skal siden gjøre dette tydelig og ikke gi inntrykk av at samme billett kan konverteres på nytt uten videre.`
- [ ] `Gitt at en billett er en middagsbillett, når kortet vises, så skal siden tydelig markere dette og ikke presentere billetten som en vanlig deltakerbillett som enkelt kan konverteres.`
- [ ] `Gitt at en billett kan konverteres, når admin velger å konvertere den, så skal billetten bli til riktig billettholder uten at admin må gjette om handlingen faktisk lyktes.`
- [ ] `Gitt at konvertering av billett feiler, når admin forsøker å konvertere, så skal admin få en tydelig feiltilstand og ikke stå igjen med en side som ser oppdatert ut uten at endringen faktisk ble gjort.`
- [ ] `Gitt at admin konverterer flere billetter etter hverandre, når siden oppdateres fortløpende, så skal riktig status vises på riktige kort og ikke blandes mellom billetter.`
- [ ] `Gitt at admin bruker søk eller filtrering på siden, når relevante treff vises, så skal resultatene være forståelige og markeringen av søket ikke gjøre innholdet uleselig.`
- [ ] `Gitt at admin bruker søk eller filtrering med tomt eller lite nyttig søk, når siden oppdateres, så skal brukeropplevelsen fortsatt være stabil og ikke gi inntrykk av at data har forsvunnet.`
- [ ] `Gitt at admin bruker siden med mange billetter og varierende data, når oversikten vises, så skal kortene fortsatt være lesbare og ikke bryte grid eller flyt.`
- [ ] `Gitt at admin navigerer tilbake til billettholderoversikten etter konvertering, når oversikten åpnes, så skal den nye billettholderen være håndtert konsistent med resten av systemet.`
- [ ] `Gitt at siden brukes på mobil eller smal skjerm, når lange navn, e-poster eller mange kort vises, så skal siden fortsatt være brukbar og ikke falle visuelt sammen.`
- [ ] `Gitt at admin refresher siden etter konvertering eller søk, når siden lastes inn igjen, så skal innholdet samsvare med faktisk lagret tilstand og ikke med en foreldet mellomtilstand.`

## Kan automatiseres

- Visning av billetter med ulike typer og tilstander egner seg godt for ende-til-ende-tester med seedet testdata.
- Konvertering av billett til billettholder egner seg godt for ende-til-ende-tester og integrasjonstester.
- Håndtering av allerede konverterte billetter og middagsbilletter egner seg godt for integrasjonstester og ende-til-ende-tester.
- Søk og filtrering egner seg godt for ende-til-ende-tester som verifiserer oppdatert visning og stabil markering av treff.
