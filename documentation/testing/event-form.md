# Arrangementsskjema

Denne sjekklisten dekker opprettelse og redigering av arrangementer under `/profile/new/{id}` og adminredigering fra godkjenningsflyten under `/admin/approval/edit/{id}`.

## Roller

- Ikke-innlogget bruker
- Innlogget bruker
- Admin

## Sjekkliste

### Skjema og lagring

- [ ] **Bare eget arrangement kan åpnes**<br>
  **Gitt** at en innlogget bruker åpner `/profile/new/{id}` for et arrangement som tilhører en annen bruker.<br>
  **Når** siden lastes.<br>
  **Så** vises 404-siden og ikke skjemaet.

- [ ] **Skjemaet laster uten brutte felt**<br>
  **Gitt** at en innlogget bruker åpner skjemaet for et nytt eller eksisterende arrangement.<br>
  **Når** siden lastes.<br>
  **Så** skal skjemaet vises uten brutte felter, tomme kort eller uforståelig innhold.

- [ ] **Seksjoner er tydelig adskilt**<br>
  **Gitt** at skjemaet vises.<br>
  **Når** brukeren ser overskrifter og seksjoner.<br>
  **Så** skal Status, «Om arrangøren», «Om arrangementet» og «Andre detaljer» fremstå tydelig adskilt og forståelige. Admin ser i tillegg «Puljefordeling».

- [ ] **Data relatert til arrangement og arrangør lagres og beholdes**<br>
  **Gitt** at brukeren endrer data.<br>
  **Når** feltet mister fokus og lagres automatisk.<br>
  **Så** skal «Lagret» vises, og endringen skal bli værende og ikke forsvinne ved oppdatering av siden.

### Validering og innsending

- [ ] **Svake data sendes ikke som gyldig arrangement**<br>
  **Gitt** at brukeren fyller inn ufullstendige, korte eller åpenbart svake data.<br>
  **Når** brukeren forsøker å sende inn arrangementet.<br>
  **Så** skal skjemaet ikke oppføre seg som om innsendingen var fullført uten at data faktisk er gyldige.

- [ ] **Ugyldige verdier får tydelig respons**<br>
  **Gitt** at brukeren skriver ugyldige eller urealistiske verdier i felter som telefonnummer eller maks antall spillere.<br>
  **Når** brukeren lagrer eller sender inn.<br>
  **Så** skal brukeropplevelsen gjøre det tydelig om verdiene aksepteres eller avvises.

- [ ] **Lange redigeringsøkter bevarer lagrede felt**<br>
  **Gitt** at brukeren arbeider lenge i skjemaet.<br>
  **Når** flere felt endres etter hverandre.<br>
  **Så** skal tidligere lagrede felt ikke nullstilles, overskrives eller hoppe mellom verdier.

- [ ] **Gjenåpning viser siste lagrede verdier**<br>
  **Gitt** at brukeren åpner skjemaet på nytt etter å ha gjort endringer.<br>
  **Når** siden lastes på nytt.<br>
  **Så** skal de siste lagrede verdiene vises og ikke eldre eller delvise versjoner av dataene.

### Bilde og adminflyt

- [ ] **Bildeflyt beholder riktig arrangement**<br>
  **Gitt** at brukeren åpner lenken for å laste opp bilde fra skjemaet.<br>
  **Når** bildeflyten åpnes og brukeren kommer tilbake.<br>
  **Så** skal arrangementet fortsatt være knyttet til riktig skjema og riktig arrangement.

- [ ] **Vanlig bruker kan bare velge Kladd eller Forkastet**<br>
  **Gitt** at en vanlig bruker åpner statusvalget for en kladd, et innsendt eller et forkastet arrangement.<br>
  **Når** valgene vises.<br>
  **Så** kan brukeren velge Kladd eller Forkastet, mens Innsendt er deaktivert med teksten «bruk knappen nederst i skjemaet».

- [ ] **Vellykket innsending gir tydelig status**<br>
  **Gitt** at brukeren trykker «Send inn» når arrangementet er klart.<br>
  **Når** innsendingen lykkes.<br>
  **Så** får arrangementet status Innsendt, og brukeren sendes tilbake til Min Side der arrangementet vises som innsendt.

- [ ] **Innsendingsfeil stopper videre flyt**<br>
  **Gitt** at innsending av arrangement feiler.<br>
  **Når** brukeren forsøker å sende inn.<br>
  **Så** skal brukeren få en tydelig feilmelding ved «Send inn»-knappen og ikke bli sendt videre som om innsendingen lyktes.

- [ ] **Admin kan redigere fra godkjenningsflyten**<br>
  **Gitt** at en admin åpner arrangementet via godkjenningsflyten.<br>
  **Når** skjemaet vises.<br>
  **Så** skal admin kunne redigere alle statuser (også Godkjent og Annonsert), puljefordeling og relevante felt uten å møte brukerbegrensningene som gjelder vanlige brukere.

- [ ] **Adminendringer oppdaterer skjema og forhåndsvisning**<br>
  **Gitt** at en admin redigerer et arrangement i godkjenningsflyten.<br>
  **Når** endringene lagres.<br>
  **Så** skal både skjema og forhåndsvisning oppdatere seg konsistent.

- [ ] **Forhåndsvisningen er bare en visning**<br>
  **Gitt** at forhåndsvisningen vises, i godkjenningsflyten eller for admin på `/profile/new/{id}`.<br>
  **Når** admin ser på forhåndsvisningen.<br>
  **Så** vises arrangementet uten interessepanel og uten forrige/neste-navigasjon. Vanlige brukere ser ingen forhåndsvisning.

### Mobil og stabilitet

- [ ] **Skjemaet er brukbart på mobil**<br>
  **Gitt** at skjemaet brukes på mobil.<br>
  **Når** mange felt og tekstområder fylles ut.<br>
  **Så** skal siden fortsatt være lesbar, skrollbar og brukbar uten at felter eller knapper havner utenfor skjermen.

- [ ] **Raske feltendringer skaper ikke mistet innhold**<br>
  **Gitt** at skjemaet brukes med raske endringer i mange felt.<br>
  **Når** brukeren navigerer mellom felt og tilbake.<br>
  **Så** skal det ikke oppstå åpenbare race conditions, mistet innhold eller ustabil oppførsel.
