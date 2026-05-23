# Min Side

Denne sjekklisten dekker `/profile`, altså den innloggede oversiktssiden med egne arrangementer, eget festivalprogram og lenke videre til billettadministrasjon.

## Roller

- Innlogget bruker

## Sjekkliste

- [ ] `Gitt at en innlogget bruker åpner Min Side, når siden lastes, så skal siden vises som en helhetlig oversikt uten brutte seksjoner eller tydelig manglende innhold.`
- [ ] `Gitt at en innlogget bruker åpner Min Side, når siden er ferdig lastet, så skal brødsmulestien tydelig vise at brukeren befinner seg på Min Side.`
- [ ] `Gitt at brukeren har egne arrangementer, når seksjonen for egne arrangementer vises, så skal arrangementene fremstå som brukerens egne og ha lenker som passer arrangementets status.`
- [ ] `Gitt at brukeren ikke har egne arrangementer, når Min Side vises, så skal seksjonen for egne arrangementer håndteres på en forståelig måte uten å se ødelagt eller tom ut på en utilsiktet måte.`
- [ ] `Gitt at brukeren har et arrangement i kladd eller ikke-ferdig status, når brukeren åpner det fra Min Side, så skal brukeren sendes til redigering og ikke til den publiserte arrangementsvisningen.`
- [ ] `Gitt at brukeren har et innsendt eller publisert arrangement, når brukeren åpner det fra Min Side, så skal brukeren bli sendt til riktig arrangementsvisning og ikke til feil redigeringsflyt.`
- [ ] `Gitt at brukeren har billettinnehavere knyttet til seg, når seksjonen for billetter vises på Min Side, så skal navn, billettype og e-post fremstå lesbart og uten å kuttes eller blandes sammen.`
- [ ] `Gitt at brukeren ikke har billettinnehavere knyttet til seg, når Min Side vises, så skal billettseksjonen håndteres uten at hele siden fremstår som feil eller mangelfull.`
- [ ] `Gitt at brukeren åpner Min Side, når seksjonen for festivalprogram vises, så skal tildelte arrangementer eller registrerte interesser vises i riktig pulje der det finnes data.`
- [ ] `Gitt at brukeren ikke har noe i en eller flere puljer i sitt festivalprogram, når Min Side vises, så skal tomtilstanden være forståelig og ikke se ut som om data har forsvunnet.`
- [ ] `Gitt at brukeren både har tildelte arrangementer og registrerte interesser, når festivalprogrammet vises, så skal tildelte arrangementer og interesser presenteres på en måte som ikke skaper tvil om hva som er faktisk tildeling og hva som bare er interesse.`
- [ ] `Gitt at datagrunnlaget for Min Side er ufullstendig eller delvis mangler, når siden vises, så skal siden fortsatt være brukbar og ikke bryte layouten eller stoppe all visning.`
- [ ] `Gitt at brukeren trykker seg videre til billettsiden fra Min Side, når navigasjonen skjer, så skal brukeren ende på riktig underside for billetter uten feil kontekst eller feil rolle.`
- [ ] `Gitt at brukeren bruker Min Side på mobil, når de ulike seksjonene vises under hverandre, så skal kort, paneler og overskrifter være lesbare og ikke overlappe.`
- [ ] `Gitt at brukeren bruker Min Side på større skjerm, når innholdet fordeles i kolonner, så skal seksjonene fremstå balansert og uten uventede tomrom eller kollisjon mellom paneler.`
- [ ] `Gitt at brukeren refresher Min Side eller kommer tilbake fra en underside, når siden vises igjen, så skal innholdet fortsatt stemme med brukerens faktiske data og rolle.`
- [ ] `Gitt at brukeren ser Min Side som helhet, når siden er ferdig lastet, så skal den fremstå som en ferdig del av produktet uten placeholder-preg eller visuelt forstyrrende detaljer.`

## Kan automatiseres

- Visning av egne arrangementer, billetter og festivalprogram egner seg godt for ende-til-ende-tester der brukeren seedes med forskjellige datakombinasjoner.
- Riktig lenkeoppførsel fra egne arrangementer egner seg godt for ende-til-ende-tester som verifiserer at status bestemmer hvor brukeren sendes.
- Tomtilstander og delvise data egner seg godt for integrasjonstester eller ende-til-ende-tester med seedet database.
- Layout og rekkefølge mellom seksjonene egner seg godt for nettleserbaserte ende-til-ende-tester på mobil og desktop.

