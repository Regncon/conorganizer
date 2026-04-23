# Forside

Denne sjekklisten dekker forsiden på `/`. Forsiden er en sentral inngang til appen og skal fungere for både ikke-innlogget bruker, innlogget bruker og admin.

## Roller

- Ikke-innlogget bruker
- Innlogget bruker
- Admin

## Sjekkliste

- [ ] `Gitt at brukeren åpner forsiden, når siden er ferdig lastet, så skal brødsmulestien vise at brukeren er på Hjem.`
- [ ] `Gitt at brukeren ikke er innlogget, når forsiden vises, så skal hovednavigasjonen ikke gi inntrykk av at brukeren har tilgang til Min Side eller adminfunksjoner.`
- [ ] `Gitt at brukeren er innlogget, når forsiden vises, så skal hovednavigasjonen gi tilgang til Min Side uten å vise feil rolle eller feil menyvalg.`
- [ ] `Gitt at brukeren er admin, når forsiden vises, så skal hovednavigasjonen gi tilgang til admin uten å skjule eller forveksle vanlige brukerfunksjoner.`
- [ ] `Gitt at brukeren åpner forsiden på mobil, når hovednavigasjonen brukes, så skal menyen være lesbar, trykkbar og stabil uten overlapp, avkuttede etiketter eller elementer som havner utenfor skjermen.`
- [ ] `Gitt at brukeren åpner forsiden på større skjerm, når hovednavigasjonen brukes, så skal logo, menyknapper og brukermeny oppføre seg konsistent og uten visuelle brudd.`
- [ ] `Gitt at brukeren åpner forsiden, når seksjonen for å sende inn arrangement vises, så skal innholdet være forståelig, lesbart og fremstå som en tydelig invitasjon til å registrere arrangement.`
- [ ] `Gitt at brukeren åpner forsiden, når knappen for å sende inn arrangement brukes, så skal brukeren bli sendt til riktig flyt for arrangementer uten å havne i en utydelig mellomtilstand.`
- [ ] `Gitt at brukeren åpner forsiden på en liten skjerm, når seksjonen for å sende inn arrangement vises, så skal tekst, knapp og illustrasjon være lesbare og ikke presse hverandre ut av layouten.`
- [ ] `Gitt at brukeren åpner forsiden på en større skjerm, når seksjonen for å sende inn arrangement vises, så skal tekst, knapp og illustrasjon være balansert og uten tomrom eller skjevheter som får innholdet til å se ødelagt ut.`
- [ ] `Gitt at det finnes godkjente arrangementer i én eller flere puljer, når brukeren åpner forsiden, så skal hver pulje vises med korrekt navn og tidspunkt.`
- [ ] `Gitt at det finnes godkjente arrangementer i flere puljer, når brukeren åpner forsiden, så skal arrangementene vises under riktig pulje og ikke lekke over i feil seksjon.`
- [ ] `Gitt at det ikke finnes arrangementer i en bestemt pulje, når brukeren åpner forsiden, så skal forsiden håndtere dette uten å vise ødelagt layout eller misvisende innhold i den puljen.`
- [ ] `Gitt at det ikke finnes noen godkjente arrangementer å vise, når brukeren åpner forsiden, så skal siden fortsatt fremstå som stabil og forståelig uten tomme kort eller ødelagte seksjoner.`
- [ ] `Gitt at forsiden viser arrangementskort, når brukeren leser kortene, så skal tittel, ingress, arrangør, system og ikoner fremstå lesbare og ikke være erstattet av åpenbart feil eller misvisende standardinnhold uten at det er forståelig hvorfor.`
- [ ] `Gitt at et arrangement mangler deler av innholdet sitt, når kortet vises på forsiden, så skal kortet fortsatt fremstå forståelig og ikke bryte layouten eller skape tvil om hva som er arrangementets faktiske data.`
- [ ] `Gitt at et arrangementskort vises på forsiden, når brukeren trykker på kortet, så skal brukeren sendes til riktig arrangementside og beholde riktig kontekst for valgt pulje.`
- [ ] `Gitt at brukeren navigerer mellom puljene via snarveinavigasjonen på forsiden, når brukeren trykker på en pulje, så skal siden scrolle til riktig seksjon uten å havne merkbart feil eller skjule seksjonsoverskriften bak sticky navigasjon.`
- [ ] `Gitt at brukeren scroller på forsiden, når snarveinavigasjonen er synlig, så skal den oppføre seg stabilt og ikke dekke viktig innhold på en måte som gjør siden vanskelig å bruke.`
- [ ] `Gitt at brukeren går direkte til en pulje via snarveinavigasjonen, når seksjonen blir synlig, så skal det være tydelig hvilken pulje brukeren har navigert til.`
- [ ] `Gitt at brukeren bruker tilbakeknappen etter å ha åpnet et arrangement fra forsiden, når brukeren kommer tilbake, så skal forsiden fortsatt være brukbar og ikke miste viktige deler av tilstanden sin.`
- [ ] `Gitt at brukeren refresher forsiden, når siden lastes på nytt, så skal innhold og navigasjon fortsatt vises korrekt uten at brukeren havner i en uforståelig tilstand.`
- [ ] `Gitt at brukeren åpner forsiden mens innhold eller data ikke kan lastes som forventet, når siden viser feiltilstand, så skal feilen være forståelig nok til at brukeren ikke sitter igjen med en tilsynelatende tom eller ødelagt side uten forklaring.`
- [ ] `Gitt at forsiden viser en feilmelding ved last av arrangementer, når brukeren ser feilen, så skal resten av siden fortsatt være brukbar så langt det lar seg gjøre.`
- [ ] `Gitt at brukeren bruker forsiden med tastatur eller andre alternative navigasjonsformer, når fokus flyttes mellom interaktive elementer, så skal det være mulig å forstå hvor brukeren befinner seg og hvilke handlinger som kan utføres.`
- [ ] `Gitt at forsiden brukes over tid med flere raske klikk på navigasjon og kort, når brukeren forflytter seg mellom sider, så skal det ikke oppstå åpenbare duplikathandlinger, feilnavigasjon eller ustabil oppførsel.`
- [ ] `Gitt at forsiden vises med ekte eller store datamengder, når mange arrangementer finnes i samme eller flere puljer, så skal siden fortsatt være lesbar, navigerbar og uten tydelige sammenbrudd i layout eller informasjonsstruktur.`


## Kan automatiseres

- Visning av forsiden for ikke-innlogget bruker egner seg godt for en ende-til-ende-test som verifiserer at riktig navigasjon og riktig hovedinnhold vises.
- Ulike navigasjonsvalg for ikke-innlogget bruker, innlogget bruker og admin egner seg godt for ende-til-ende-tester som verifiserer rollebasert oppførsel i toppnavigasjonen.
- Visning av puljer og arrangementer egner seg godt for ende-til-ende-tester eller integrasjonstester der databasen seedes med arrangementer i ulike puljer og med ulike datakombinasjoner.
- Klikk på arrangementskort og bevaring av riktig pulje i lenken egner seg godt for en ende-til-ende-test.
- Snarveinavigasjon mellom puljer egner seg godt for en nettleserbasert ende-til-ende-test som verifiserer scrolling og riktig ankeroppførsel.
- Feiltilstand ved manglende eller utilgjengelige arrangementsdata egner seg for en integrasjonstest eller ende-til-ende-test som verifiserer at brukeren ikke blir sittende igjen med en stille og uforståelig feil.
