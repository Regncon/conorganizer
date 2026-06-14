# Arrangementsskjema

Denne sjekklisten dekker opprettelse og redigering av arrangementer under `/profile/new/{id}` og adminredigering fra godkjenningsflyten.

## Roller

- Ikke-innlogget bruker
- Innlogget bruker
- Admin

## Sjekkliste

### Skjema og lagring

- [ ] **Skjemaet laster uten brutte felt**<br>
  **Gitt** at en innlogget bruker åpner skjemaet for et nytt eller eksisterende arrangement.<br>
  **Når** siden lastes.<br>
  **Så** skal skjemaet vises uten brutte felter, tomme kort eller uforståelig innhold.

- [ ] **Seksjoner er tydelig adskilt**<br>
  **Gitt** at skjemaet vises.<br>
  **Når** brukeren ser overskrifter og seksjoner.<br>
  **Så** skal status, arrangørinformasjon, arrangementsinformasjon og øvrige detaljer fremstå tydelig adskilt og forståelige.

- [ ] **Arrangørdata lagres og beholdes**<br>
  **Gitt** at brukeren endrer navn, e-post eller telefon i arrangørseksjonen.<br>
  **Når** feltet lagres.<br>
  **Så** skal endringen bli værende og ikke forsvinne ved oppdatering av siden.

- [ ] **Arrangementstekst og metadata persisterer**<br>
  **Gitt** at brukeren endrer tittel, intro, type, system eller beskrivelse.<br>
  **Når** feltene lagres.<br>
  **Så** skal innholdet persistere og vises konsistent ved refresh og ved senere åpning av skjemaet.

- [ ] **Detaljfelt vises korrekt etter lagring**<br>
  **Gitt** at brukeren endrer alder, varighet, nybegynnervennlighet, engelskstøtte, maks antall spillere eller merknader.<br>
  **Når** feltene lagres.<br>
  **Så** skal endringene være synlige og korrekte ved senere visning.

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

- [ ] **Vellykket innsending gir tydelig status**<br>
  **Gitt** at brukeren forsøker å sende inn arrangementet når det er klart.<br>
  **Når** innsendingen lykkes.<br>
  **Så** skal status og videre flyt være tydelige og brukeren skal ikke bli stående igjen i tvil om at arrangementet faktisk ble sendt inn.

- [ ] **Innsendingsfeil stopper videre flyt**<br>
  **Gitt** at innsending av arrangement feiler.<br>
  **Når** brukeren forsøker å sende inn.<br>
  **Så** skal brukeren få en tydelig feiltilstand og ikke bli sendt videre som om innsendingen lyktes.

- [ ] **Admin kan redigere fra godkjenningsflyten**<br>
  **Gitt** at en admin åpner arrangementet via godkjenningsflyten.<br>
  **Når** skjemaet vises.<br>
  **Så** skal admin kunne redigere status og relevante felt uten å møte brukerbegrensningene som gjelder vanlige brukere.

- [ ] **Adminendringer oppdaterer skjema og forhåndsvisning**<br>
  **Gitt** at en admin redigerer et arrangement i godkjenningsflyten.<br>
  **Når** endringene lagres.<br>
  **Så** skal både skjema og forhåndsvisning oppdatere seg konsistent.

### Mobil og stabilitet

- [ ] **Skjemaet er brukbart på mobil**<br>
  **Gitt** at skjemaet brukes på mobil.<br>
  **Når** mange felt og tekstområder fylles ut.<br>
  **Så** skal siden fortsatt være lesbar, skrollbar og brukbar uten at felter eller knapper havner utenfor skjermen.

- [ ] **Raske feltendringer skaper ikke mistet innhold**<br>
  **Gitt** at skjemaet brukes med raske endringer i mange felt.<br>
  **Når** brukeren navigerer mellom felt og tilbake.<br>
  **Så** skal det ikke oppstå åpenbare race conditions, mistet innhold eller ustabil oppførsel.
