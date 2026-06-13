# Romadministrasjon i admin

Denne sjekklisten dekker `/admin/rooms` og `/admin/rooms/assignment/{pulje}`, der admin kan administrere rom og se romfordeling per pulje.

## Roller

- Admin

## Sjekkliste

- [ ] `Gitt at en admin åpner romoversikten, når siden lastes, så skal rom vises gruppert på en forståelig måte uten brutte paneler eller tomme kort.`
- [ ] `Gitt at det ikke finnes rom ennå, når romoversikten vises, så skal tomtilstanden forklare hva admin kan gjøre videre.`
- [ ] `Gitt at admin legger til et nytt rom med gyldige verdier, når handlingen lagres, så skal rommet dukke opp i riktig etasje med riktige detaljer.`
- [ ] `Gitt at admin prøver å lagre rom med ugyldige verdier, når valideringen kjøres, så skal feilen være tydelig og rommet skal ikke lagres som gyldig data.`
- [ ] `Gitt at admin endrer et eksisterende rom, når handlingen lagres, så skal oppdaterte verdier vises uten at rom-ID eller andre rom blir påvirket.`
- [ ] `Gitt at admin sletter et rom, når handlingen bekreftes, så skal bare riktig rom fjernes fra oversikten.`
- [ ] `Gitt at admin åpner romfordeling for en pulje, når siden lastes, så skal rom og tildelte arrangementer høre til riktig pulje.`
- [ ] `Gitt at arrangementer mangler rom i en pulje, når romfordelingen vises, så skal de være tydelige som manglende tildelinger.`
- [ ] `Gitt at admin tildeler et arrangement til et rom, når handlingen lykkes, så skal arrangementet vises under riktig rom og ikke fortsatt som manglende rom.`
- [ ] `Gitt at et deaktivert rom har tildelte arrangementer, når romfordelingen vises, så skal advarselen være tydelig nok til at admin kan rette opp fordelingen.`
- [ ] `Gitt at romadministrasjonen brukes på mobil eller smal skjerm, når mange rom og arrangementer vises, så skal innholdet fortsatt være lesbart og mulig å arbeide med.`
