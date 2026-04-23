# Arrangementsdetaljer

Denne sjekklisten dekker `/event/{id}`, altså den publiserte detaljvisningen for et arrangement, inkludert visning av detaljer, bilder, puljer, forrige/neste-navigasjon og interesseflyten.

## Roller

- Ikke-innlogget bruker
- Innlogget bruker uten billetter
- Innlogget bruker med billetter
- Admin

## Sjekkliste

- [ ] `Gitt at brukeren åpner et gyldig arrangement, når siden lastes, så skal tittel, introduksjon, bilde, detaljer og beskrivelse vises som en helhetlig arrangementsvisning uten brutte hovedseksjoner.`
- [ ] `Gitt at brukeren åpner et arrangement med lang eller innholdsrik beskrivelse, når siden vises, så skal teksten være lesbar og ikke bryte layouten eller forsvinne på en uventet måte.`
- [ ] `Gitt at et arrangement har ett eller flere puljer, når siden vises, så skal puljene presenteres tydelig og ikke blandes med andre arrangementers data.`
- [ ] `Gitt at arrangementet har ulike egenskaper som aldersgruppe, varighet, nybegynnervennlighet eller engelskstøtte, når siden vises, så skal disse egenskapene fremstå korrekt og uten motstridende signaler.`
- [ ] `Gitt at arrangementet mangler deler av valgfri informasjon, når siden vises, så skal detaljsiden fortsatt fremstå robust og ikke se ødelagt ut.`
- [ ] `Gitt at brukeren åpner et arrangement med bilder, når siden vises på mobil og større skjerm, så skal riktige bilder brukes og fremstå som en bevisst del av siden og ikke som ødelagte eller feilskalerte flater.`
- [ ] `Gitt at brukeren åpner et arrangement som ikke finnes eller ikke kan lastes riktig, når siden vises, så skal brukeren møte en forståelig feiltilstand og ikke en halvferdig arrangementsvisning.`
- [ ] `Gitt at brukeren åpner et arrangement fra forsiden med valgt puljekontekst, når detaljsiden vises, så skal videre oppførsel rundt puljer og navigasjon være i tråd med det arrangementet brukeren valgte.`
- [ ] `Gitt at brukeren bruker forrige- og neste-navigasjon på arrangementsiden, når brukeren går videre mellom arrangementer, så skal navigasjonen holde brukeren innenfor riktig logisk rekkefølge uten å hoppe uforståelig.`
- [ ] `Gitt at brukeren bruker tilbakeknappen etter å ha åpnet et arrangement, når brukeren går tilbake, så skal forrige side fortsatt være brukbar og ikke miste tydelig kontekst.`
- [ ] `Gitt at en ikke-innlogget bruker åpner arrangementsiden, når interessepanelet vises, så skal siden ikke late som om brukeren kan registrere interesse uten å ha riktig tilgang eller billettgrunnlag.`
- [ ] `Gitt at en innlogget bruker uten billetter åpner arrangementsiden, når interessepanelet vises, så skal siden tydelig forklare at brukeren først må hente billett og gi en klar vei videre til billettsiden.`
- [ ] `Gitt at en innlogget bruker med billetter åpner arrangementsiden, når interessepanelet vises, så skal det være tydelig at brukeren kan melde interesse og hva denne handlingen innebærer.`
- [ ] `Gitt at en innlogget bruker med billetter åpner interessemodalen, når modalen vises, så skal valg av billettholder, pulje og interesse fremstå tydelig og uten at brukeren må gjette hva som skal gjøres.`
- [ ] `Gitt at brukeren prøver å registrere interesse uten å ha valgt nødvendig billettholder eller pulje, når handlingen utføres, så skal brukeren få en tydelig og forståelig feilmelding i stedet for en stille feil.`
- [ ] `Gitt at brukeren registrerer interesse med gyldige valg, når valget lagres, så skal tilstanden oppdatere seg uten at modalen eller siden havner i en ødelagt eller forvirrende tilstand.`
- [ ] `Gitt at brukerens valg ikke kan lagres, når brukeren forsøker å melde interesse, så skal siden tydelig vise at handlingen feilet og ikke late som om interessen ble registrert.`
- [ ] `Gitt at brukeren lukker interessemodalen etter å ha gjort valg, når modalen forsvinner, så skal resten av arrangementsiden fortsatt være stabil og brukbar.`
- [ ] `Gitt at admin åpner en arrangementsdetaljside, når siden vises, så skal admin kunne se tydelig inngang til administrativ redigering uten at den vanlige arrangementsvisningen blir ødelagt.`
- [ ] `Gitt at arrangementsiden brukes på mobil, når brukeren skroller gjennom bilde, header, detaljer, interessepanel og beskrivelse, så skal siden forbli lesbar og navigerbar uten overlapp eller merkbar layoutkollaps.`

## Kan automatiseres

- Visning av arrangementsdetaljer med ulike datakombinasjoner egner seg godt for ende-til-ende-tester og integrasjonstester.
- Interessepanelets ulike oppførsel for ikke-innlogget bruker, bruker uten billetter og bruker med billetter egner seg godt for ende-til-ende-tester.
- Feilmeldinger ved manglende billettholder eller pulje i interesseflyten egner seg godt for ende-til-ende-tester.
- Forrige- og neste-navigasjon mellom arrangementer egner seg godt for ende-til-ende-tester.
- Adminlenken på arrangementsdetaljsiden egner seg godt for ende-til-ende-tester med adminrolle.

