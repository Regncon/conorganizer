# Arrangementsdetaljer

Denne sjekklisten dekker `/event/{id}`, altså den publiserte detaljvisningen for et arrangement, inkludert visning av detaljer, bilder, puljer, forrige/neste-navigasjon og interesseflyten.

## Roller

- Ikke-innlogget bruker
- Innlogget bruker uten billetter
- Innlogget bruker med billetter
- Admin

## Sjekkliste

### Innhold og robusthet

- [ ] **Arrangementet vises som en helhetlig detaljside**<br>
  **Gitt** at brukeren åpner et gyldig arrangement.<br>
  **Når** siden lastes.<br>
  **Så** skal tittel, introduksjon, bilde, detaljer og beskrivelse vises som en helhetlig arrangementsvisning uten brutte hovedseksjoner.

- [ ] **Lang beskrivelse bryter ikke layouten**<br>
  **Gitt** at brukeren åpner et arrangement med lang eller innholdsrik beskrivelse.<br>
  **Når** siden vises.<br>
  **Så** skal teksten være lesbar og ikke bryte layouten eller forsvinne på en uventet måte.

- [ ] **Arrangementsegenskaper vises korrekt**<br>
  **Gitt** at arrangementet har ulike egenskaper som aldersgruppe, varighet, nybegynnervennlighet eller engelskstøtte.<br>
  **Når** siden vises.<br>
  **Så** skal disse egenskapene fremstå korrekt og uten motstridende signaler.

- [ ] **Manglende valgfri informasjon håndteres robust**<br>
  **Gitt** at arrangementet mangler deler av valgfri informasjon.<br>
  **Når** siden vises.<br>
  **Så** skal detaljsiden fortsatt fremstå robust og ikke se ødelagt ut.

- [ ] **Bilder skaleres og vises riktig**<br>
  **Gitt** at brukeren åpner et arrangement med bilder.<br>
  **Når** siden vises på mobil og større skjerm.<br>
  **Så** skal riktige bilder brukes og fremstå som en bevisst del av siden og ikke som ødelagte eller feilskalerte flater.

- [ ] **Manglende arrangement gir forståelig feiltilstand**<br>
  **Gitt** at brukeren åpner et arrangement som ikke finnes eller ikke kan lastes riktig.<br>
  **Når** siden vises.<br>
  **Så** skal brukeren møte en forståelig feiltilstand og ikke en halvferdig arrangementsvisning.

### Interesseflyt

- [ ] **Interessepanelet forklarer handlingen**<br>
  **Gitt** at en innlogget bruker med billetter åpner arrangementsiden.<br>
  **Når** interessepanelet vises.<br>
  **Så** skal det være tydelig at brukeren kan melde interesse og hva denne handlingen innebærer.

- [ ] **Interessemodalen viser tydelige valg**<br>
  **Gitt** at en innlogget bruker med billetter åpner interessemodalen.<br>
  **Når** modalen vises.<br>
  **Så** skal valg av billettholder, pulje og interesse fremstå tydelig og uten at brukeren må gjette hva som skal gjøres.

- [ ] **Lagret interesse oppdaterer tilstanden**<br>
  **Gitt** at brukeren registrerer interesse med gyldige valg.<br>
  **Når** valget lagres.<br>
  **Så** skal tilstanden oppdatere seg uten at modalen eller siden havner i en ødelagt eller forvirrende tilstand.

- [ ] **Lukket modal etterlater siden stabil**<br>
  **Gitt** at brukeren lukker interessemodalen etter å ha gjort valg.<br>
  **Når** modalen forsvinner.<br>
  **Så** skal resten av arrangementsiden fortsatt være stabil og brukbar.

### Mobil

- [ ] **Mobilvisning er lesbar og navigerbar**<br>
  **Gitt** at arrangementsiden brukes på mobil.<br>
  **Når** brukeren skroller gjennom bilde, header, detaljer, interessepanel og beskrivelse.<br>
  **Så** skal siden forbli lesbar og navigerbar uten overlapp eller merkbar layoutkollaps.
