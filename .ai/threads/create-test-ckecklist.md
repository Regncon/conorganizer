--- Codex Session ---
ID: 019db8b2-be49-71c1-8642-34cc1cc93dc8

#Create test checklist 

We are aproaching launch and we need a checklist to make sure we have tested everything before we launch.

There are some automated tests and we don't need to test those manually.

Have a look at the pages/ directory to get a sense of the different pages we have. 
Each page should have it's own checklist of things to test.

Other directories of interest are the components/ directory and the service/

The tests are written in Norwegian. 

There is a starting point in documentation/testing/manualTesting.md

Have a look at the codebase and gather relevant information to create a comprehensive checklist for manual testing.

Let's have a conversation about this before we start writing the checklist.

Time: 2026-04-23 06:58:23
--- Codex Run ---

```Markdown
Jeg går gjennom `pages/`, `components/`, `service/` og eksisterende `documentation/testing/manualTesting.md` for å kartlegge funksjonalitet, testflater og hva som allerede dekkes automatisk. Deretter oppsummerer jeg strukturen og foreslår hvordan vi bør avgrense den manuelle sjekklisten før vi skriver noe.
Jeg har nå oversikten over sider og den eksisterende manuelle lista. Neste steg er å koble sidene til faktisk funksjonalitet i koden, så vi kan skille mellom ting som bør testes manuelt og ting som allerede er godt dekket automatisk.
Jeg ser allerede noen tydelige testområder: mye skjer via inline oppdateringer/SSE, mange skjema-felter oppdateres enkeltvis, og admin/event-interesse har flere tilstandsoverganger som ikke bør stoles blindt på. Jeg vil nå hente rutekartet og se hvilke deler som ser ut til å være dekket av automatiske tester, så vi kan unngå unødvendig manuelt arbeid.
Bildet blir tydeligere: det finnes en del rene enhetstester i `service/checkIn`, `service/eventService` og noe i `components/formsubmission`, men de dekker ikke brukerflyt, UI-tilstand eller tilgangsopplevelse. Jeg henter nå de mest sentrale komponentene rundt profil, billettinnehavere og arrangementsinnsending, siden de sannsynligvis må gi struktur til den manuelle testen per side.
Jeg har nok kontekst til å begynne å snakke om struktur og avgrensning. Før jeg foreslår hvordan sjekklisten bør se ut, henter jeg én siste oversikt over synlige sider/overskrifter og særskilte funksjoner som lett blir glemt, som print-side og admin-redigering.
Jeg har gått gjennom kodebasen nok til å foreslå hvordan vi bør angripe den manuelle sjekklisten før vi skriver den.

Det viktigste funnet er at appen ikke bare består av “sider”, men av flere brukerflyter med mye live-oppdatering og inkrementelle endringer. Mange ting skjer via Datastar/SSE og små endepunkter for enkeltfelt, særlig i arrangementskjemaet og admin. Det betyr at en god manuell sjekkliste ikke bare bør si “siden laster”, men også teste at UI faktisk oppdateres riktig, at data persisteres, og at feiltilstander vises forståelig. Det ser også ut som eksisterende manuell liste i [documentation/testing/manualTesting.md](/home/cinmay/Documents/conorganizer/documentation/testing/manualTesting.md:1) er nyttig som start, men den er mer en løpende buglogg enn en launch-sjekkliste.

Jeg ville delt den manuelle testingen i disse hovedområdene:

1. Sidevise sjekklister
   `Forside`, `Innlogging`, `Min side`, `Mine billetter`, `Mine arrangementer`, `Nytt/rediger arrangement`, `Arrangementsdetalj`, `Printvennlig side`, `Admin`, `Godkjenning`, `Admin billettholdere`, `Admin legg til billettholder`.

2. Tverrgående sjekker
   Navigasjon, tilgangskontroll, responsivt oppsett, tomtilstander, feiltilstander, lasting/spinnere, persistering etter refresh, og språk/tekst.

3. Rollebaserte scenarier
   Uinnlogget bruker, vanlig innlogget bruker uten billetter, innlogget bruker med billetter/billettholdere, arrangør med egne arrangementer, admin.

Det som ser ut til å være dekket automatisk i dag er ganske snevert:
- `service/checkIn/*_test.go`: konvertering/tilknytning av billetter og billettholdere.
- `service/eventService/previous_next_test.go`: forrige/neste-navigasjon.
- `components/formsubmission/who_is_interested_test.go`: regler for interesse/“first choice”.

Det betyr at jeg ikke ville lagt mye manuell vekt på å re-verifisere de interne reglene i detalj, men heller fokusert manuelt på:
- at brukerflytene faktisk trigger riktig backendlogikk,
- at UI viser riktig resultat,
- at feil blir synlige for brukeren,
- at tilgang og rollebegrensninger fungerer i praksis.

Sider/flyter som ser mest launch-kritiske ut:
- `Mine arrangementer` og `Nytt/rediger arrangement` i [pages/myprofile/myevents/myevents.go](/home/cinmay/Documents/conorganizer/pages/myprofile/myevents/myevents.go:1)
- `Arrangementsdetalj` med interessevalg i [pages/event/event.go](/home/cinmay/Documents/conorganizer/pages/event/event.go:1)
- `Profil` og `Mine billetter` i [pages/profile/profile_index_templ.go](/home/cinmay/Documents/conorganizer/pages/profile/profile_index_templ.go:29) og [pages/profile/tickets/tickets_page.templ](/home/cinmay/Documents/conorganizer/pages/profile/tickets/tickets_page.templ:1)
- adminflytene i [pages/admin/admin.go](/home/cinmay/Documents/conorganizer/pages/admin/admin.go:1), [pages/admin/approval/editForm/edit_form_page.templ](/home/cinmay/Documents/conorganizer/pages/admin/approval/editForm/edit_form_page.templ:1) og [pages/admin/billettholder_admin/billettholder_admin.go](/home/cinmay/Documents/conorganizer/pages/admin/billettholder_admin/billettholder_admin.go:1)

Jeg ser også noen ting vi bør bestemme før jeg skriver selve sjekklisten:
- Skal sjekklisten være “minimum før launch”, eller også inkludere lavere prioritet kosmetiske ting?
- Vil du at jeg skal merke punkter med prioritet, f.eks. `Må testes`, `Bør testes`, `Nice to have`?
- Skal vi skrive den som rene avkrysningspunkter, eller også inkludere “testdata/oppsett” per side?
- Har dere faktisk tilgang til en admin-bruker for manuell test, eller må checklisten også beskrive hvordan admin-scenarier kan verifiseres?

Min anbefaling er at jeg lager en norsk launch-sjekkliste med:
- én seksjon per side/flyt,
- korte, konkrete avkrysningspunkter,
- egne seksjoner for `Tilgangskontroll`, `Feilhåndtering` og `Responsivt`,
- og en lett prioritering så teamet vet hva som må kjøres før lansering.

Hvis du vil, kan jeg neste runde foreslå en konkret disposisjon for `manualTesting.md` uten å fylle inn alle testpunktene ennå.
```

Time: 2026-04-23 06:58:23
--- Next Prompt ---

Sekklisten skal også inkludere lavere prioritet kosmetiske ting.
Ikke merk med prioritet, alt er viktig. Alt må fungere før launch.
Sekklisten skal være avkrysningspunkter. Vi går ut i fra at de som skal teste vet hvordan de skal teste, så det er ikke nødvendig å inkludere testdata/oppsett i sjekklisten.
Hvis det ikke er test oppsett så må den som tester tenke seg om. Vi vill at testerne skal tenke selv og ikke følge en oppskrift slavisk. 
Vi har full tilgang til admin-bruker for manuell test. 
Test miljøet er helt likt produksjon, så det er ingen spesielle hensyn å ta der. Vi kan også legge inn så mye testdata som vi ønsker, så det er ingen begrensning der heller.
Lag en makrdown fil pr side/flyt, og en hovedfil som linker til alle de andre.
Lage en egen fil for generelle terter som gjelder flere sider. 
Sørg for å ha god dekning av feilhåndtering. Test edge cases og ikke bare happy path. 
Hver fil bør inneholde sine egene edege cases og feilhåndtering.

La oss forsette dialogen før vi begynner å skrive sjekklistene.

Time: 2026-04-23 07:10:48
--- Codex Run ---

```Markdown
Jeg avgrenser dette som et sett med side-/flytspesifikke markdownfiler pluss én generell fil for tverrgående tester. Nå vil jeg kartlegge en fornuftig filstruktur i `documentation/testing/` og peke ut hvor grensene mellom filene bør gå, slik at vi unngår både hull og unødvendig duplisering.
Da ville jeg strukturert dette som en faktisk testpakke, ikke én stor fil. Basert på rutene som finnes nå, ser jeg for meg denne strukturen under `documentation/testing/`:

- `README.md` eller `index.md`
- `general.md`
- `root.md`
- `auth.md`
- `profile.md`
- `profile-tickets.md`
- `my-profile.md`
- `my-events.md`
- `event-form.md`
- `event-details.md`
- `print.md`
- `admin.md`
- `admin-approval.md`
- `admin-billettholders.md`
- `admin-add-billettholder.md`

Det viktigste jeg vil avklare før jeg skriver, er at kodebasen har to forskjellige “profil”-spor:
- `/my-profile` er en enklere eldre side med lenker til “Mine arrangement” og “Mine billetter”.
- `/profile` er en nyere “Min side” med egne seksjoner for arrangementer, program og billetter.
- `/my-profile/my-tickets` og `/profile/tickets` overlapper også tematisk, men er ikke samme flyt.

Jeg mener derfor at vi bør teste begge dersom begge skal eksistere ved launch. Hvis dere egentlig bare ser på én av dem som “den ekte” profilen, bør vi enten:
- fortsatt teste begge, fordi brukere kan nå begge rutene, eller
- eksplisitt behandle den andre som legacy og gi den en kortere sjekkliste.

Et annet viktig skille er hva som bør ligge i `general.md` versus i hver egen fil. Min anbefaling:
- `general.md`: hovednavigasjon, responsivitet, styling/kosmetiske avvik, lastetilstander, konsistens i språk/tekst, 404/403/500-lignende oppførsel, refresh/back-knapp, lenker, bilder/ikoner.
- Sidefilene: funksjon, edge cases og feilhåndtering som er spesifikk for akkurat den siden/flyten.

Jeg ville også bevisst unngått å lage egne “happy path”- og “edge case”-seksjoner i hver fil. Siden dere vil at testerne skal tenke selv, tror jeg det blir bedre med én flat liste per fil der funksjonelle, kosmetiske og negative scenarier står blandet. Det blir mindre oppskriftsaktig og mer som en reell kontrolliste.

To ting jeg vil ha bekreftet før vi starter:
- Skal `/root` regnes som en reell side som må inn i launch-sjekklisten, eller er den intern/eksperimentell og kan utelates?
- Vil dere at hovedfilen bare skal være en innholdsfortegnelse med lenker, eller også inneholde en kort beskrivelse av hva hver testfil dekker?
```

Time: 2026-04-23 07:10:48
--- Next Prompt ---
Du kan droppe å teste print siden. Den en ikke en del av den vanlige brukerreisen. 
 /my-profile/ og /my-profile/my-tickets er deprecated og skal ikke være en del av launch-sjekklisten.
/root er root routen for appen og må være med. Det en home, forsiden, og det er viktig at den fungerer.

Det er også viktig at vi identifiserer tester som kan automatiseres. Vi ønsker å ha så mye som mulig dekket av automatiske tester. 
En annen viktig ting er å fokusere på å teste oppførsel. Det er først og fremst oppførsel vi ønsker å teste. Trekk så mye som mulig inspirasjon fra Behavior-driven development.
Teksten på testene kan godt være litt lang og beskrivende. 

La oss forsette dialogen før vi begynner å skrive sjekklistene.


Time: 2026-04-23 08:50:32
--- Codex Run ---

```Markdown
Da ville jeg justert planen slik:

Filer som bør inngå:
- `index.md` som hovedfil med lenker
- `general.md`
- `root.md`
- `auth.md`
- `profile.md`
- `profile-tickets.md`
- `my-events.md`
- `event-form.md`
- `event-details.md`
- `admin.md`
- `admin-approval.md`
- `admin-billettholders.md`
- `admin-add-billettholder.md`

Det passer godt med det som faktisk er i aktiv bruk i rutene nå. `print` og deprecated `/my-profile*` holder vi helt ute.

Jeg er enig i at oppførsel bør være førende. I stedet for korte mekaniske punkter som “knapp fungerer”, ville jeg skrevet testene i en BDD-lignende stil, men fortsatt som avkrysningspunkter. Typisk:
- “Når en innlogget bruker åpner forsiden og det finnes publiserte arrangementer i flere puljer, skal arrangementene vises i riktig pulje og navigasjonselementene skal være konsistente.”
- “Når en bruker forsøker å sende inn et arrangement med ugyldige eller mangelfulle data, skal brukeren få tydelig tilbakemelding og ingen delvis lagring skal skje uten at det er forventet.”

Det gir dere sjekklister som tester oppførsel, ikke bare UI-elementer.

Når det gjelder automatisering, ser jeg to nivåer:

1. Godt egnet for automatisering nå
- Forside: visning av arrangementer per pulje, tomtilstander, navigasjon til arrangement.
- Auth: innlogging/utlogging, tilgangskontroll, redirecter.
- Profil: visning av egne arrangementer, billetter og festivalprogram.
- Mine arrangementer: opprette nytt arrangement, redigere eksisterende, validering, innsending.
- Arrangementsdetalj: interessevalg, feil ved ugyldige valg, forskjell på bruker med og uten billetter.
- Admin: tilgangskontroll til admin-ruter.
- Godkjenning: listevisning, statusoverganger, redigering av arrangement.
- Billettholder-admin: søk, koble e-post, slette e-post, konvertere billett til billettholder.

2. Fortsatt best manuelt, men kan delvis automatiseres senere
- Kosmetiske avvik og layoutkollisjoner.
- Tekstlig tydelighet og språk.
- Uventede kombinasjoner av data som ser “rare” ut, men ikke nødvendigvis feiler teknisk.
- Mer subjektive feilmeldinger og forståelighet.

Min anbefaling er derfor at hver markdown-fil får to deler:
- først selve sjekklisten med BDD-lignende avkrysningspunkter
- til slutt en kort seksjon som heter `Kan automatiseres`, med de punktene fra akkurat den filen som egner seg for E2E eller integrasjonstest

Det gir dere både launch-sjekklisten og en konkret backlog for mer automatisering.

Det jeg vil avklare før jeg skriver:
- Vil du at `Kan automatiseres` skal være i hver fil, eller vil du heller ha én samlet `automation-candidates.md`?
- Vil du at formuleringene skal være ren norsk BDD-stil som “Når … skal …”, eller vil du ha enda tydeligere `Gitt / Når / Så`-form?
```

Time: 2026-04-23 08:50:32
--- Next Prompt ---

Kan automatiseres kan godt være i hver enkelt fil. 
Bruk helst `Gitt / Når / Så`-form
Lag også et dokument som handler om hvordan vi lager tesene våre. Målgruppen er både utviklere og LLM AI kode assistenter.
Er det flere ting vi bør diskutere? 
Er det noen beste praksiser eller andre ting som omhandler testene vi bør snakke om?

Time: 2026-04-23 09:08:36
--- Codex Run ---

```Markdown
Ja. Det er noen ting som er verdt å avklare nå, fordi de vil påvirke kvaliteten på alle filene.

Det viktigste er språk og stil. Applikasjonen bruker en blanding av norsk språkføring i UI, og enkelte steder heller tydelig mot nynorsk. For sjekklistene bør vi velge én konsekvent form. Min anbefaling er:
- skriv testene i enkel, klar norsk
- bruk `Gitt / Når / Så`
- beskriv forventet oppførsel, ikke implementasjon
- unngå å referere til interne tekniske detaljer i selve sjekkpunktene

Eksempel på riktig nivå:
- `Gitt at brukeren er logget inn og har minst én billettholder, når brukeren åpner et arrangement med aktive puljer, så skal brukeren kunne registrere interesse uten å bli sendt ut av siden.`

Eksempel på feil nivå:
- `Gitt at Datastar sender SSE-oppdatering, når PATCH mottas, så skal DOM oppdateres.`

Beste praksis jeg ville lagt til grunn for disse filene:
- Ett sjekkpunkt skal verifisere én tydelig oppførsel.
- Hver fil skal være brukerorientert og følge én side eller flyt.
- Edge cases og feilhåndtering skal stå side om side med normal oppførsel, ikke i et eget “restledd”.
- Kosmetiske ting skal beskrives som observerbar oppførsel.
  Eksempel: `Så skal teksten være lesbar og ikke overlappe andre elementer`.
- `Kan automatiseres` bør peke på testnivå, ikke bare si “ja”.
  Eksempel: `Egner seg for E2E`, `Egner seg for integrasjonstest`, `Egner seg for enhetstest`.
- Checklistene bør være stabile nok til å tåle mindre UI-endringer, ellers blir de fort utdaterte.

Dokumentet om hvordan vi lager testene våre bør etter min mening dekke:
- formål med manuelle sjekklister
- hvordan vi skriver `Gitt / Når / Så`
- forskjellen på å teste oppførsel og implementasjon
- hvordan vi vurderer `Kan automatiseres`
- hvordan vi skriver tester som både mennesker og LLM-er kan utvide uten å forringe kvaliteten
- regler for språk, detaljnivå og filstruktur

Jeg ville laget det som noe i retning av:
- `documentation/testing/how-to-write-tests.md`

Ting det fortsatt kan være nyttig å avklare:
- Skal vi skrive alt på bokmål, eller vil dere at testene skal følge samme språkdrakt som produktet der det er naturlig?
- Vil dere at `Kan automatiseres` skal bruke faste etiketter, for eksempel `E2E`, `Integrasjon`, `Enhet`?
- Vil dere at hovedfilen også skal forklare hvilke roller som må brukes under test, eller bare lenke videre?

Min vurdering er at vi har avklart det meste som betyr noe strukturelt. Det eneste jeg ville låst før jeg begynner å skrive, er språkføring og formatet på `Kan automatiseres`.
```

Time: 2026-04-23 09:08:36
--- Next Prompt ---
Vi kan godt skrive alt på bokmål.
Ikke bruk etiketter. Lage heller en kort beskrivelse av hva slags test det egner seg for.
Jeg tror ikke roller bør ligge i hovedfilen, men heller i hver enkelt fil der det er relevant. Tydelig merke om det er admin, bruker, eller ikke-innlogget som skal teste.
La oss begynne med å lage `documentation/testing/how-to-write-tests.md` og root.md som en start. Da har vi både en guide for hvordan vi skriver testene, og et konkret eksempel på hvordan en sjekkliste kan se ut. Når vi har det på plass, kan vi bruke det som mal for resten av filene.

Time: 2026-04-23 09:20:12
--- Codex Run ---

```Markdown
Jeg lager først stilguiden for testene og en konkret sjekkliste for forsiden, slik at resten kan følge samme mønster. Før jeg skriver filene, henter jeg den faktiske oppførselen på forsiden og hovednavigasjonen, så `root.md` blir basert på kode og ikke antakelser.
Jeg har nok til å skrive begge filene. Jeg lager en stilguide som låser formatet for resten av arbeidet, og en `root.md` som bruker samme mønster med tydelige roller, `Gitt / Når / Så`, feilhåndtering, edge cases og en egen del for hva som egner seg for automatisering.
Jeg skriver nå to nye dokumenter i `documentation/testing/`: én stilguide og én konkret sjekkliste for forsiden. Jeg holder formatet stramt nok til å kunne brukes som mal for resten av filene.
Begge filene er skrevet. Jeg gjør en rask gjennomlesning nå for å sikre at tonen, `Gitt / Når / Så`-formatet og balansen mellom oppførsel, edge cases og automatiseringskandidater faktisk fungerer som mal før vi bygger videre på resten.
Jeg har lagt til [how-to-write-tests.md](/home/cinmay/Documents/conorganizer/documentation/testing/how-to-write-tests.md:1) og [root.md](/home/cinmay/Documents/conorganizer/documentation/testing/root.md:1).

`how-to-write-tests.md` låser formatet for resten: bokmål, avkrysningspunkter, `Gitt / Når / Så`, fokus på observerbar oppførsel, egne forventninger til edge cases, feilhåndtering, kosmetiske forhold og en avsluttende seksjon om hva som kan automatiseres.

`root.md` bruker samme mønster som en konkret mal for videre arbeid. Den dekker forsiden for ikke-innlogget bruker, innlogget bruker og admin, og inkluderer navigasjon, CTA, puljevisning, arrangementskort, snarveinavigasjon, feiltilstander, store datamengder og kosmetisk konsistens. Den avslutter også med en kort og konkret seksjon om automatiseringsmuligheter.

Ingen tester ble kjørt; dette er kun dokumentasjonsarbeid. Hvis du vil, kan vi nå gjøre en kort innholdsgjennomgang av disse to filene og justere stilen før jeg fortsetter med `general.md` eller neste side/flyt.
```

Time: 2026-04-23 09:20:12
--- Next Prompt ---
Lage generl.md nå.
Jeg har slettet noen irrelevante tester fra `root.md`.
Flytt alt som har med meny navigasjon å gjøre til general.md. Det er ikke spesifikt for forsiden, og det er bedre å samle det på ett sted.
Flytt navigasjon til `send inn arrangement` til `event-form.md`. Det er mer relevant der, og det gjør `root.md` mer fokusert på forsiden som helhet. Det skal være forskjellig oppførsel på `send inn arrangement`-knappen avhengig av om brukeren er logget inn eller ikke. 




Time: 2026-04-23 09:46:01
--- Codex Run ---

```Markdown
Jeg oppdaterer dokumentasjonen i tråd med den nye avgrensningen: `general.md` får alt som er tverrgående for meny og navigasjon, og `root.md` blir strammere rundt forsiden. Før jeg gjør endringen, leser jeg den nåværende `root.md` og ser på hvordan uinnlogget tilgang faktisk håndteres, slik at formuleringene om navigasjon blir riktige.
Jeg ser at uinnlogget tilgang til beskyttede sider gir en eksplisitt “Du har ikkje tilgang”-side med innloggingslenke. Jeg henter resten av menykomponenten før jeg skriver `general.md`, så den fanger både Hjem, Min Side, Admin og bruker-/innloggingsmenyen korrekt.
Jeg flytter nå all delt meny-/navigasjonsoppførsel ut av `root.md` og inn i en ny `general.md`. Samtidig lar jeg `root.md` stå igjen med det som faktisk er spesifikt for forsiden: puljer, arrangementskort, scrollnavigasjon på programmet og den visuelle helheten.
`general.md` er lagt til, og `root.md` er strammet inn slik at den ikke lenger eier felles menyoppførsel. Jeg tar en rask kontroll nå for å sikre at fordelingen mellom de to filene er ren, og at ingenting åpenbart relevant for forsiden ble flyttet ut ved et uhell.
Jeg har lagt til [general.md](/home/cinmay/Documents/conorganizer/documentation/testing/general.md:1) og flyttet felles meny-/navigasjonsoppførsel ut av [root.md](/home/cinmay/Documents/conorganizer/documentation/testing/root.md:1).

`general.md` samler nå hovednavigasjon, rollebasert synlighet i menyen, innlogging, utlogging, brukermeny, tilgang til beskyttede sider og generelle navigasjonsforventninger på tvers av appen. `root.md` er strammere og handler nå om forsiden som forside: CTA-seksjonen som innhold, puljer, arrangementskort, snarveinavigasjon, feiltilstander og helhetsopplevelse.

Jeg har bevisst ikke flyttet selve oppførselen for klikk på `Send inn arrangement` inn i noen fil ennå, siden du vil ha den i `event-form.md` og den filen ikke finnes ennå. Ingen tester ble kjørt; dette er kun dokumentasjonsendringer.

Neste naturlige steg er `event-form.md`, siden den nå eier den rolleavhengige oppførselen til `Send inn arrangement`-inngangen.
```

Time: 2026-04-23 09:46:01
--- Next Prompt ---

