# Romadministrasjon i admin

Denne sjekklisten dekker `/admin/rooms` og `/admin/rooms/assignment/{pulje}`, der admin kan administrere rom og se romfordeling per pulje.

## Roller

- Admin

## Sjekkliste

### Romoversikt

- [ ] **Romoversikten grupperer rom forståelig**<br>
  **Gitt** at en admin åpner romoversikten.<br>
  **Når** siden lastes.<br>
  **Så** skal rom vises gruppert på en forståelig måte uten brutte paneler eller tomme kort.

- [ ] **Tom romoversikt forklarer neste steg**<br>
  **Gitt** at det ikke finnes rom ennå.<br>
  **Når** romoversikten vises.<br>
  **Så** skal tomtilstanden forklare hva admin kan gjøre videre.

### Romendringer

- [ ] **Nytt rom vises i riktig etasje**<br>
  **Gitt** at admin legger til et nytt rom med gyldige verdier.<br>
  **Når** handlingen lagres.<br>
  **Så** skal rommet dukke opp i riktig etasje med riktige detaljer.

- [ ] **Ugyldige romverdier avvises**<br>
  **Gitt** at admin prøver å lagre rom med ugyldige verdier.<br>
  **Når** valideringen kjøres.<br>
  **Så** skal feilen være tydelig og rommet skal ikke lagres som gyldig data.

- [ ] **Romendringer påvirker bare riktig rom**<br>
  **Gitt** at admin endrer et eksisterende rom.<br>
  **Når** handlingen lagres.<br>
  **Så** skal oppdaterte verdier vises uten at rom-ID eller andre rom blir påvirket.

- [ ] **Sletting fjerner bare riktig rom**<br>
  **Gitt** at admin sletter et rom.<br>
  **Når** handlingen bekreftes.<br>
  **Så** skal bare riktig rom fjernes fra oversikten.

### Romfordeling

- [ ] **Romfordeling viser riktig pulje**<br>
  **Gitt** at admin åpner romfordeling for en pulje.<br>
  **Når** siden lastes.<br>
  **Så** skal rom og tildelte arrangementer høre til riktig pulje.

- [ ] **Manglende romtildelinger er tydelige**<br>
  **Gitt** at arrangementer mangler rom i en pulje.<br>
  **Når** romfordelingen vises.<br>
  **Så** skal de være tydelige som manglende tildelinger.

- [ ] **Romtildeling flytter arrangementet riktig**<br>
  **Gitt** at admin tildeler et arrangement til et rom.<br>
  **Når** handlingen lykkes.<br>
  **Så** skal arrangementet vises under riktig rom og ikke fortsatt som manglende rom.

- [ ] **Deaktivert rom med tildelinger varsles tydelig**<br>
  **Gitt** at et deaktivert rom har tildelte arrangementer.<br>
  **Når** romfordelingen vises.<br>
  **Så** skal advarselen være tydelig nok til at admin kan rette opp fordelingen.

### Mobil

- [ ] **Romadministrasjon fungerer på smal skjerm**<br>
  **Gitt** at romadministrasjonen brukes på mobil eller smal skjerm.<br>
  **Når** mange rom og arrangementer vises.<br>
  **Så** skal innholdet fortsatt være lesbart og mulig å arbeide med.
