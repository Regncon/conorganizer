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

### Romkart

- [ ] **Romkart vises først etter publisering**<br>
  **Gitt** at arrangementet har et tildelt rom.<br>
  **Når** både programmet og puljefordelingen er publisert, og arrangementet er aktivt og publisert i puljen.<br>
  **Så** vises romnavn og «Se rommet» for rom med kjent kart. Før publisering skal verken romnavn eller kart rendres.

- [ ] **Hver pulje viser sitt eget rom**<br>
  **Gitt** at arrangementet har ulike rom i to publiserte puljer.<br>
  **Når** brukeren åpner hvert romkart.<br>
  **Så** vises riktig pulje, romnavn, etasje og SVG-kart, også uten innlogging.

- [ ] **Kartet forblir åpent ved liveoppdateringer**<br>
  **Gitt** at kartmodalen er åpen.<br>
  **Når** en endring i interesse, arrangement eller rom utløser en NATS/Datastar-oppdatering.<br>
  **Så** forblir samme modal åpen. Endring av romnavn eller romtildeling oppdaterer innholdet uten å lukke modalen.

- [ ] **Tilbaketrukket romkart lukkes**<br>
  **Gitt** at kartmodalen er åpen.<br>
  **Når** romtildelingen fjernes eller programmet/puljefordelingen avpubliseres.<br>
  **Så** forsvinner kartet. Ny publisering skal ikke åpne modalen automatisk. Rom uten kjent SVG viser rominformasjon uten kartknapp.

- [ ] **Kartmodalen fungerer med tastatur og på mobil**<br>
  **Gitt** at brukeren åpner «Se rommet» med tastatur eller på en smal skjerm.<br>
  **Når** modalen vises og lukkes med Escape, lukkeknappen eller bakgrunnen.<br>
  **Så** holdes tastaturfokus i modalen mens den er åpen og returnerer til knappen ved lukking. Hele kartet beholder proporsjonene, og innholdet kan rulles ved liten skjermhøyde.

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
