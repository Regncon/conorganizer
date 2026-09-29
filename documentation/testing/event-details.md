# Arrangementsdetaljer

Denne sjekklisten dekker `/event/{id}`, altså detaljvisningen for et arrangement, inkludert synlighet, visning av detaljer, bilder, puljer, forrige/neste-navigasjon og interesseflyten.

## Roller

- Ikke-innlogget bruker
- Innlogget bruker uten billetter
- Innlogget bruker med billetter
- Eier av arrangementet
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
  **Så** skal brukeren møte en forståelig feiltilstand og ikke en halvferdig arrangementsvisning. En arrangement-ID som ikke finnes gir 404-siden «Vi fant ikke arrangementet du leter etter».

### Synlighet

- [ ] **Annonsert arrangement er offentlig**<br>
  **Gitt** at arrangementet har status Annonsert.<br>
  **Når** en ikke-innlogget bruker, en innlogget bruker, eieren og admin åpner siden.<br>
  **Så** ser alle hele arrangementet uten varselbanner.

- [ ] **Eier og admin ser ikke-annonserte arrangementer med banner**<br>
  **Gitt** at arrangementet har status Kladd, Innsendt eller Godkjent.<br>
  **Når** eieren eller admin åpner siden.<br>
  **Så** vises hele arrangementet med banneret «Arrangementet er ikke annonsert». For status Forkastet vises i stedet banneret «Arrangementet er forkastet».

- [ ] **Andre får beskjed om at arrangementet ikke er annonsert**<br>
  **Gitt** at arrangementet har status Kladd, Innsendt eller Godkjent.<br>
  **Når** en ikke-innlogget bruker eller en innlogget bruker som ikke er eier åpner siden.<br>
  **Så** vises bare meldingen «Arrangementet er ikke annonsert ennå» (også som sidetittel), uten detaljer om arrangementet. Å være tildelt som spiller eller GM gir ikke tilgang.

- [ ] **Forkastet arrangement er utilgjengelig for andre**<br>
  **Gitt** at arrangementet har status Forkastet.<br>
  **Når** en ikke-innlogget bruker eller en innlogget bruker som ikke er eier åpner siden.<br>
  **Så** vises bare meldingen «Arrangementet er ikke tilgjengelig» (også som sidetittel), og siden svarer med HTTP 410 Gone.

### Romkart

- [ ] **Romkart vises først etter publisering**<br>
  **Gitt** at arrangementet har et tildelt rom.<br>
  **Når** programmet er publisert, romfordelingen for puljen er publisert og arrangementet har en aktiv romtildeling i puljen, uavhengig av puljens status og det gamle publiseringsflagget per arrangement.<br>
  **Så** vises romnavn og «Se rommet» for rom med kjent kart. Før programmet er publisert, skal verken puljer, klokkeslett, romnavn eller kart rendres.

- [ ] **Upublisert romfordeling skjuler rommet**<br>
  **Gitt** at programmet er publisert, men romfordelingen for puljen ikke er publisert.<br>
  **Når** en bruker eller admin åpner arrangementet.<br>
  **Så** vises puljen med klokkeslett som en tidsrad uten rom, men ikke romnavn, romnotat eller kartknapp. Admin ser det samme som alle andre.

- [ ] **Offentlige romnotater vises, admin-notater aldri**<br>
  **Gitt** at rommet har både offentlige notater og admin-notater, og romfordelingen for puljen er publisert.<br>
  **Når** detaljsiden vises og brukeren åpner romkartet.<br>
  **Så** vises de offentlige notatene under rommet i romlisten og i kartmodalen, som ren tekst også når de inneholder HTML. Admin-notatene vises ikke noe sted på siden. Rom uten offentlige notater har ingen notatlinje.

- [ ] **Hver pulje viser sitt eget rom**<br>
  **Gitt** at arrangementet har ulike rom i to puljer og programmet er publisert.<br>
  **Når** brukeren åpner hvert romkart.<br>
  **Så** vises riktig pulje, romnavn, etasje og SVG-kart, også uten innlogging.

- [ ] **Samme rom gir én knapp med alle tidspunktene**<br>
  **Gitt** at arrangementet har samme rom i flere puljer og programmet er publisert.<br>
  **Når** detaljsiden vises.<br>
  **Så** vises én kartknapp per rom med puljenavn og klokkeslett i kronologisk rekkefølge og større tekst. Det finnes ingen separat «Pulje(r)»-liste. Puljer uten rom vises som vanlige tidsrader uten kartknapp.

- [ ] **Kartet forblir åpent ved liveoppdateringer**<br>
  **Gitt** at kartmodalen er åpen.<br>
  **Når** en endring i interesse, arrangement eller rom utløser en NATS/Datastar-oppdatering.<br>
  **Så** forblir samme modal åpen. Endring av romnavn eller romtildeling oppdaterer innholdet uten å lukke modalen, også når romknapper slås sammen eller deles opp.

- [ ] **Tilbaketrukket romkart lukkes**<br>
  **Gitt** at kartmodalen er åpen.<br>
  **Når** romtildelingen fjernes, romfordelingen for puljen avpubliseres eller programmet avpubliseres.<br>
  **Så** forsvinner kartet. Ny publisering skal ikke åpne modalen automatisk. Rom uten kjent SVG viser rominformasjon uten kartknapp.

- [ ] **Kartmodalen fungerer med tastatur og på mobil**<br>
  **Gitt** at brukeren åpner «Se rommet» med tastatur eller på en smal skjerm.<br>
  **Når** modalen vises og lukkes med Escape, lukkeknappen eller bakgrunnen.<br>
  **Så** holdes tastaturfokus i modalen mens den er åpen og returnerer til knappen ved lukking. Hele kartet beholder proporsjonene, og innholdet kan rulles ved liten skjermhøyde.

### Forrige/neste

- [ ] **Før publisering følger navigasjonen den alfabetiske listen**<br>
  **Gitt** at programmet ikke er publisert.<br>
  **Når** brukeren åpner et annonsert arrangement.<br>
  **Så** peker forrige/neste til naboene i forsidens alfabetiske liste over annonserte arrangementer, med vanlige lenker til `/event/{id}` uten pulje eller dato.

- [ ] **Etter publisering følger navigasjonen dagens rekkefølge på forsiden**<br>
  **Gitt** at programmet er publisert og brukeren kommer fra forsiden med `?date=` og `?pulje=` i lenken.<br>
  **Når** brukeren blar med forrige/neste.<br>
  **Så** følger navigasjonen samme dag og samme rekkefølge som forsiden: først dagens programarrangementer, deretter arrangementene i puljefordeling pulje for pulje. Lenkene beholder `date` og `pulje`, for eksempel `/event/{id}?date=2026-10-10&pulje=LordagMorgen`.

- [ ] **Samme arrangement i flere puljer får riktige naboer**<br>
  **Gitt** at et arrangement ligger i to puljer samme dag og programmet er publisert.<br>
  **Når** siden åpnes med hver av puljene i `?pulje=`.<br>
  **Så** viser forrige/neste naboene til akkurat den forekomsten på forsiden.

- [ ] **Ingen navigasjon rundt kantene**<br>
  **Gitt** at brukeren står på første eller siste arrangement i listen.<br>
  **Når** siden vises.<br>
  **Så** mangler henholdsvis forrige eller neste helt, uten deaktivert knapp, og navigasjonen hopper ikke rundt til andre enden.

- [ ] **Ugyldig pulje gir ingen forrige/neste**<br>
  **Gitt** at programmet er publisert.<br>
  **Når** siden åpnes uten `?pulje=`, med en ugyldig verdi eller med en pulje arrangementet ikke ligger i den dagen.<br>
  **Så** vises ingen forrige/neste-navigasjon, og resten av siden fungerer som vanlig.

### Interesseflyt

- [ ] **Interessepanelet venter på publisert program**<br>
  **Gitt** at programmet ikke er publisert og arrangementet er med i puljefordelingen.<br>
  **Når** en bruker med eller uten billett åpner arrangementsiden.<br>
  **Så** viser panelet bare «Interessevalget åpner når programmet er publisert.», uten «Meld interesse»-knapp og uten «Hent billett»-lenke.

- [ ] **Bruker uten billett sendes til billettsiden**<br>
  **Gitt** at programmet er publisert og brukeren ikke er innlogget eller ikke har billett.<br>
  **Når** interessepanelet vises.<br>
  **Så** står det «For å kunne melde interesse må du først hente billetten din.» med knappen «Hent billett» til `/profile/tickets`.

- [ ] **Interessepanelet forklarer handlingen**<br>
  **Gitt** at programmet er publisert og en innlogget bruker med billetter åpner et arrangement i puljefordelingen.<br>
  **Når** interessepanelet vises.<br>
  **Så** skal det være tydelig at brukeren kan melde interesse og hva denne handlingen innebærer, med lenken «Hvordan fungerer interessevalget?» og knappen «Meld interesse».

- [ ] **Programarrangement trenger ikke interesse**<br>
  **Gitt** at arrangementet ligger i en pulje, men ikke er med i puljefordelingen.<br>
  **Når** siden vises.<br>
  **Så** forklarer panelet at dette er et programarrangement som er åpent for alle, og det finnes ingen «Meld interesse»-knapp. Har brukeren meldt interesse for andre arrangementer i samme pulje, vises et varsel med lenke til Min Side.

- [ ] **Interessemodalen viser tydelige valg**<br>
  **Gitt** at en innlogget bruker med billetter åpner interessemodalen.<br>
  **Når** modalen vises.<br>
  **Så** skal valg av billettholder, pulje og interesse fremstå tydelig og uten at brukeren må gjette hva som skal gjøres. Puljevelgeren vises bare når arrangementet ligger i flere puljer.

- [ ] **18+-arrangement viser anbefaling for billettholdere under 18**<br>
  **Gitt** at en billettholder under 18 år åpner interessemodalen for et 18+-arrangement, og ikke allerede er tildelt et arrangement i puljen.<br>
  **Når** modalen vises.<br>
  **Så** står det «Anbefalt for voksne (18+)», og billettholderen kan likevel velge og lagre interesse. Bytter brukeren til en billettholder over 18, forsvinner notisen.

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
