# Godkjenning av arrangementer

Denne sjekklisten dekker `/admin/approval` og `/admin/approval/edit/{id}`, altså oversikten over arrangementer til godkjenning og adminredigering av enkeltarrangementer med interesse- og tildelingsarbeid.

## Roller

- Admin

## Sjekkliste

- [ ] `Gitt at en admin åpner godkjenningssiden, når siden lastes, så skal innsendte og godkjente arrangementer vises i tydelige seksjoner uten å blandes sammen.`
- [ ] `Gitt at det finnes arrangementer til godkjenning, når siden vises, så skal hvert arrangement fremstå som et tydelig valg videre til redigering.`
- [ ] `Gitt at det ikke finnes arrangementer i en av seksjonene, når siden vises, så skal resten av siden fortsatt fremstå korrekt og ikke som om hele adminvisningen feiler.`
- [ ] `Gitt at en admin åpner et arrangement fra godkjenningslisten, når redigeringssiden lastes, så skal riktig arrangement vises i skjema, forhåndsvisning og tilhørende interesse-/tildelingsvisning.`
- [ ] `Gitt at admin redigerer felt i arrangementskjemaet fra godkjenningsflyten, når endringene lagres, så skal skjema og forhåndsvisning oppdatere seg konsistent.`
- [ ] `Gitt at admin endrer status på arrangementet, når endringen lagres, så skal statusendringen oppføre seg tydelig og ikke etterlate tvil om arrangementets nye tilstand.`
- [ ] `Gitt at admin bruker forrige- og neste-navigasjon i redigeringsvisningen, når brukeren går mellom arrangementer, så skal riktig arrangement åpnes i riktig rekkefølge uten forvirring.`
- [ ] `Gitt at et arrangement ikke finnes eller ikke kan lastes, når admin forsøker å åpne det i redigeringsflyten, så skal admin møte en forståelig feiltilstand og ikke en halvferdig redigeringsvisning.`
- [ ] `Gitt at admin ser oversikten over interesserte og tildelte personer på arrangementsredigeringen, når listene vises, så skal de fremstå forståelige og høre til riktig arrangement og riktig pulje.`
- [ ] `Gitt at admin legger til en deltaker som spiller via godkjenningsflyten, når handlingen lykkes, så skal tildelingen vises riktig og ikke havne på feil arrangement eller feil pulje.`
- [ ] `Gitt at admin legger til en deltaker som GM via godkjenningsflyten, når handlingen lykkes, så skal GM-rollen vises riktig og ikke forveksles med vanlig spiller.`
- [ ] `Gitt at admin endrer status for en allerede tildelt person mellom spiller, GM og fjernet, når handlingen utføres, så skal resultatet oppdateres tydelig og konsistent.`
- [ ] `Gitt at en tildelingshandling feiler, når admin forsøker å oppdatere spiller- eller GM-status, så skal feilen være tydelig nok til at admin forstår at endringen ikke ble fullført.`
- [ ] `Gitt at flere adminhandlinger utføres etter hverandre på samme side, når siden oppdateres fortløpende, så skal innholdet forbli stabilt og ikke vise gamle eller blandede data mellom seksjonene.`
- [ ] `Gitt at admin bruker godkjenningsflyten på større og mindre skjermer, når skjema, forhåndsvisning og interesseoversikt vises samtidig, så skal siden fortsatt være lesbar og arbeidsbar.`
- [ ] `Gitt at admin refresher siden midt i redigeringsarbeidet, når siden lastes inn igjen, så skal korrekt arrangementsdata og korrekt tildelingsstatus vises.`

## Kan automatiseres

- Visning av innsendte og godkjente arrangementer i egne seksjoner egner seg godt for ende-til-ende-tester.
- Redigering av arrangement fra adminflyten egner seg godt for ende-til-ende-tester som verifiserer både skjema og forhåndsvisning.
- Statusendringer og navigasjon mellom arrangementer egner seg godt for ende-til-ende-tester.
- Tildeling av spiller og GM egner seg godt for integrasjonstester og ende-til-ende-tester.
- Feiltilstander ved ugyldige tildelinger eller manglende arrangement egner seg godt for integrasjonstester.

