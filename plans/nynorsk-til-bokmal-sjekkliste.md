# Sjekkliste: nynorsk → bokmål

Fremdrift: 330 / 330 ✅ (ferdig 2026-09-24, ikke committet)

Kartlegging og bakgrunn: [nynorsk-til-bokmal.md](nynorsk-til-bokmal.md). Denne filen er fasit for hva som er gjort.

## Slik gjenopptar du

- Bare `[ ]` gjenstår. `[x]` er ferdig og skal ikke røres.
- Hver fiks krysses av med én gang filen den står i er rettet.
- Linjenumrene gjelder commit ee28b9ed. Er en linje flyttet, finn teksten i filen.
- Forslaget er et utgangspunkt: behold `%s`/`%v`, HTML og Datastar-uttrykk, og bruk domeneord fra `domeneordbok.md`.
- Ingenting committes; brukeren pusher alt i én commit til slutt.

## Batch A – Skjema 1  (101 / 101 ✅)

### [components/formsubmission/contact_info.templ](../components/formsubmission/contact_info.templ)

- [x] L24 · «Klarte ikkje å lese skjemadata» → «Klarte ikke å lese skjemadata»
- [x] L30 · «Arrangement-ID manglar» → «Arrangement-ID mangler»
- [x] L35 · «Klarte ikkje å oppdatere namnet for arrangementet i databasen» → «Klarte ikke å oppdatere navnet for arrangementet i databasen»
- [x] L39 · «Klarte ikkje å sende oppdatering» → «Klarte ikke å sende oppdatering»
- [x] L53 · «Klarte ikkje å lese skjemadata» → «Klarte ikke å lese skjemadata»
- [x] L59 · «Arrangement-ID manglar» → «Arrangement-ID mangler»
- [x] L64 · «Klarte ikkje å oppdatere e-posten for arrangementet i databasen» → «Klarte ikke å oppdatere e-posten for arrangementet i databasen»
- [x] L69 · «Klarte ikkje å sende oppdatering» → «Klarte ikke å sende oppdatering»
- [x] L82 · «Klarte ikkje å lese skjemadata» → «Klarte ikke å lese skjemadata»
- [x] L88 · «Arrangement-ID manglar» → «Arrangement-ID mangler»
- [x] L93 · «Klarte ikkje å oppdatere telefonnummeret for arrangementet i databasen» → «Klarte ikke å oppdatere telefonnummeret for arrangementet i databasen»
- [x] L98 · «Klarte ikkje å sende oppdatering» → «Klarte ikke å sende oppdatering»
- [x] L118 · «Namnet ditt» → «Navnet ditt»
- [x] L121 · «Bruk minst 2 teikn» → «Bruk minst 2 tegn»
- [x] L158 · «Bruk 8 siffer, eventuelt med +47 framfor» → «Bruk 8 siffer, eventuelt med +47 foran»

### [components/formsubmission/event_new.templ](../components/formsubmission/event_new.templ)

- [x] L86 · «Klarte ikkje å lese skjemadata» → «Klarte ikke å lese skjemadata»
- [x] L104 · «Klarte ikkje å lese skjemadata» → «Klarte ikke å lese skjemadata»
- [x] L266 · «Tittel må vere minst 3 teikn lang» → «Tittel må være minst 3 tegn lang»
- [x] L266 · «Tittel på spelmodul / arrangement» → «Tittel på spillmodul / arrangement»
- [x] L267 · «Introduksjon må vere minst 3 teikn lang» → «Introduksjon må være minst 3 tegn lang»
- [x] L268 · «Skildring» → «Beskrivelse»
- [x] L268 · «Skildring må vere minst 3 teikn lang» → «Beskrivelse må være minst 3 tegn lang»
- [x] L269 · «System må vere minst 3 teikn langt» → «System må være minst 3 tegn langt»
- [x] L270 · «Namn på arrangør» → «Navn på arrangør»
- [x] L270 · «Namn på arrangør må vere minst 3 teikn langt» → «Navn på arrangør må være minst 3 tegn langt»
- [x] L271 · «E-postadressa må vere gyldig» → «E-postadressen må være gyldig»
- [x] L272 · «Telefonnummeret må vere minst 8 siffer» → «Telefonnummeret må være minst 8 siffer»
- [x] L273 · «Vel ei gyldig aldersgruppe» → «Velg en gyldig aldersgruppe»
- [x] L274 · «Varigheit» → «Varighet»
- [x] L274 · «Vel ei gyldig varigheit» → «Velg en gyldig varighet»
- [x] L275 · «Maks tal på spelarar» → «Maks antall spillere»
- [x] L275 · «Maks tal på spelarar må vere mellom 1 og 18» → «Maks antall spillere må være mellom 1 og 18»
- [x] L276 · «Nybegynnarvennleg» → «Nybegynnervennlig»
- [x] L277 · «Kan haldast på engelsk» → «Kan holdes på engelsk»
- [x] L295 · «Takk for påmeldinga!» → «Takk for påmeldingen!»  _(borderline)_
- [x] L301 · «Klarte ikkje å lagre arrangementet. Prøv igjen. Feil:» → «Klarte ikke å lagre arrangementet. Prøv igjen. Feil:»
- [x] L581 · «Klarte ikkje å oppdatere arrangementet: %v» → «Klarte ikke å oppdatere arrangementet: %v»
- [x] L595 · «Fann ikkje arrangementet, eller ingen endringar vart gjorde» → «Fant ikke arrangementet, eller ingen endringer ble gjort»

### [components/formsubmission/other_details.templ](../components/formsubmission/other_details.templ)

- [x] L25 · «Klarte ikkje å lese skjemadata» → «Klarte ikke å lese skjemadata»
- [x] L31 · «Arrangement-ID manglar» → «Arrangement-ID mangler»
- [x] L36 · «Klarte ikkje å oppdatere aldersgruppa for arrangementet i databasen» → «Klarte ikke å oppdatere aldersgruppen for arrangementet i databasen»
- [x] L41 · «Klarte ikkje å sende oppdatering» → «Klarte ikke å sende oppdatering»
- [x] L55 · «Klarte ikkje å lese skjemadata» → «Klarte ikke å lese skjemadata»
- [x] L61 · «Arrangement-ID manglar» → «Arrangement-ID mangler»
- [x] L66 · «Klarte ikkje å oppdatere varigheita for arrangementet i databasen» → «Klarte ikke å oppdatere varigheten for arrangementet i databasen»
- [x] L71 · «Klarte ikkje å sende oppdatering» → «Klarte ikke å sende oppdatering»
- [x] L85 · «Klarte ikkje å lese skjemadata» → «Klarte ikke å lese skjemadata»
- [x] L91 · «Arrangement-ID manglar» → «Arrangement-ID mangler»
- [x] L96 · «Klarte ikkje å oppdatere nybegynnarvennleg-statusen for arrangementet i databasen» → «Klarte ikke å oppdatere nybegynnervennlig-statusen for arrangementet i databasen»
- [x] L101 · «Klarte ikkje å sende oppdatering» → «Klarte ikke å sende oppdatering»
- [x] L115 · «Klarte ikkje å lese skjemadata» → «Klarte ikke å lese skjemadata»
- [x] L121 · «Arrangement-ID manglar» → «Arrangement-ID mangler»
- [x] L126 · «Klarte ikkje å oppdatere om arrangementet kan haldast på engelsk» → «Klarte ikke å oppdatere om arrangementet kan holdes på engelsk»
- [x] L131 · «Klarte ikkje å sende oppdatering» → «Klarte ikke å sende oppdatering»
- [x] L145 · «Klarte ikkje å lese skjemadata» → «Klarte ikke å lese skjemadata»
- [x] L151 · «Arrangement-ID manglar» → «Arrangement-ID mangler»
- [x] L156 · «Klarte ikkje å oppdatere maks tal på spelarar for arrangementet i databasen» → «Klarte ikke å oppdatere maks antall spillere for arrangementet i databasen»
- [x] L161 · «Klarte ikkje å sende oppdatering» → «Klarte ikke å sende oppdatering»
- [x] L175 · «Klarte ikkje å lese skjemadata» → «Klarte ikke å lese skjemadata»
- [x] L181 · «Arrangement-ID manglar» → «Arrangement-ID mangler»
- [x] L188 · «Klarte ikkje å oppdatere merknadane for arrangementet i databasen» → «Klarte ikke å oppdatere merknadene for arrangementet i databasen»
- [x] L193 · «Klarte ikkje å sende oppdatering» → «Klarte ikke å sende oppdatering»
- [x] L202 · «Viss arrangementet ditt er skremmande eller krev litt meir modne deltakarar, kan du velje dette for å setje ei attenårsgrense.» → «Hvis arrangementet ditt er skremmende eller krever litt mer modne deltakere, kan du velge dette for å sette en attenårsgrense.»
- [x] L204 · «På Regncon tel dei som er tolv år eller yngre som barn. Arrangement som passar for barn, blir sette opp i søndagspuljen.» → «På Regncon regnes de som er tolv år eller yngre som barn. Arrangementer som passer for barn, blir satt opp i søndagspuljen.»
- [x] L206 · «Standard aldersgruppe. Ho omfattar dei som kan delta i alle puljene, altså alle som er tretten år eller eldre.» → «Standard aldersgruppe. Den omfatter de som kan delta i alle puljene, altså alle som er tretten år eller eldre.»
- [x] L224 · «Vel denne viss arrangementet truleg varer mindre enn tre timar. Då veit deltakarane at dei må finne noko å gjere medan dei ventar på neste …» → «Velg denne hvis arrangementet trolig varer mindre enn tre timer. Da vet deltakerne at de må finne noe å gjøre mens de venter på neste pulje.»
- [x] L226 · «Vel denne viss arrangementet truleg varer seks timar eller meir. Då veit deltakarane at det kan påverke nattesøvnen eller den neste puljen.» → «Velg denne hvis arrangementet trolig varer seks timer eller mer. Da vet deltakerne at det kan påvirke nattesøvnen eller den neste puljen.»
- [x] L228 · «Ein vanleg pulje varer 4-5 timar på Regncon. Erfaringsmessig er dette varigheita til dei aller fleste arrangementa som blir haldne på conne…» → «En vanlig pulje varer 4-5 timer på Regncon. Erfaringsmessig er dette varigheten til de aller fleste arrangementene som blir holdt på connet.»
- [x] L245 · «Andre detaljar» → «Andre detaljer»
- [x] L266 · «Kva alder er arrangementet tilrådd for?» → «Hvilken alder er arrangementet anbefalt for?»
- [x] L270 · «Varigheit» → «Varighet»
- [x] L288 · «Kor lenge varer puljen?» → «Hvor lenge varer puljen?»
- [x] L294 · «Nybegynnarvennleg arrangement» → «Nybegynnervennlig arrangement»
- [x] L296 · «Valfritt» → «Valgfritt»
- [x] L311 · «Arrangementet er nybegynnarvennleg» → «Arrangementet er nybegynnervennlig»
- [x] L316 · «Arrangementet er nybegynnarvennleg» → «Arrangementet er nybegynnervennlig»
- [x] L318 · «Erfaringsmessig stemmer dette for dei fleste arrangement – men viss spelet er komplisert,» → «Erfaringsmessig stemmer dette for de fleste arrangementer – men hvis spillet er komplisert,»
- [x] L319 · «og du ikkje vil bruke tid på å forklare reglar, kan du la boksen stå tom.» → «og du ikke vil bruke tid på å forklare regler, kan du la boksen stå tom.»
- [x] L324 · «Kan haldast på engelsk» → «Kan holdes på engelsk»
- [x] L325 · «Valfritt» → «Valgfritt»
- [x] L340 · «Arrangementet kan haldast på engelsk» → «Arrangementet kan holdes på engelsk»
- [x] L345 · «Arrangementet kan haldast på engelsk» → «Arrangementet kan holdes på engelsk»
- [x] L347 · «Vel denne viss du kan halde arrangementet på engelsk» → «Velg denne hvis du kan holde arrangementet på engelsk»
- [x] L348 · «(så finn du eventuelt ut av det med deltakarane når arrangementet startar)» → «(så finner du eventuelt ut av det med deltakerne når arrangementet starter)»
- [x] L352 · «Maks tal på spelarar (utanom arrangør)» → «Maks antall spillere (utenom arrangør)»
- [x] L368 · «Kor mange kan vere med og spele?» → «Hvor mange kan være med og spille?»
- [x] L370 · «Har du mulighet til å hjelp dine medspillere?» → «Har du mulighet til å hjelpe dine medspillere?»  _(skrivefeil)_
- [x] L373 · «… enn du orginalt tenkte?» → «… enn du opprinnelig tenkte?»  _(skrivefeil)_
- [x] L377 · «Andre merknader – er det til dømes nokre tidspunkt det ikkje passar å halde arrangementet?» → «Andre merknader – er det for eksempel noen tidspunkt det ikke passer å holde arrangementet?»

### [components/formsubmission/set_event_in_pulje_test.go](../components/formsubmission/set_event_in_pulje_test.go)

- [x] L21 · «INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players) VALUES ('e1', 'Spel', '', '', 'Ola', 'ola@x…» → «INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players) VALUES ('e1', 'Spill', '', '', 'Ola', 'ola@…»
- [x] L65 · «INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players) VALUES ('e1', 'Spel', '', '', 'Ola', 'ola@x…» → «INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players) VALUES ('e1', 'Spill', '', '', 'Ola', 'ola@…»
- [x] L112 · «INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players) VALUES ('e1', 'Spel', '', '', 'Ola', 'ola@x…» → «INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players) VALUES ('e1', 'Spill', '', '', 'Ola', 'ola@…»

### [components/formsubmission/statusCard.templ](../components/formsubmission/statusCard.templ)

- [x] L26 · «Klarte ikkje å lese skjemadata» → «Klarte ikke å lese skjemadata»
- [x] L32 · «Arrangement-ID manglar» → «Arrangement-ID mangler»
- [x] L37 · «Klarte ikkje å oppdatere statusen for arrangementet i databasen» → «Klarte ikke å oppdatere statusen for arrangementet i databasen»
- [x] L42 · «Klarte ikkje å sende oppdatering» → «Klarte ikke å sende oppdatering»
- [x] L72 · «(bruk knappen nedst i skjemaet)» → «(bruk knappen nederst i skjemaet)»

### [components/formsubmission/submit_section.templ](../components/formsubmission/submit_section.templ)

- [x] L18 · «Arrangement-ID manglar» → «Arrangement-ID mangler»
- [x] L24 · «Klarte ikkje å oppdatere statusen for arrangementet i databasen» → «Klarte ikke å oppdatere statusen for arrangementet i databasen»
- [x] L29 · «Klarte ikkje å sende oppdatering» → «Klarte ikke å sende oppdatering»
- [x] L46 · «Kladden blir lagra automatisk. Når du er klar til å sende arrangementet inn, treng du berre å klikke på knappen her!» → «Kladden blir lagret automatisk. Når du er klar til å sende arrangementet inn, trenger du bare å klikke på knappen her!»

## Batch B – Skjema 2  (84 / 84 ✅)

### [components/formsubmission/about_event.templ](../components/formsubmission/about_event.templ)

- [x] L28 · «Klarte ikkje å lese skjemadata» → «Klarte ikke å lese skjemadata»
- [x] L34 · «Arrangement-ID manglar» → «Arrangement-ID mangler»
- [x] L39 · «Klarte ikkje å oppdatere tittelen for arrangementet i databasen» → «Klarte ikke å oppdatere tittelen for arrangementet i databasen»
- [x] L51 · «Klarte ikkje å sende oppdatering» → «Klarte ikke å sende oppdatering»
- [x] L65 · «Klarte ikkje å lese skjemadata» → «Klarte ikke å lese skjemadata»
- [x] L71 · «Arrangement-ID manglar» → «Arrangement-ID mangler»
- [x] L78 · «Klarte ikkje å oppdatere introduksjonen for arrangementet i databasen» → «Klarte ikke å oppdatere introduksjonen for arrangementet i databasen»
- [x] L82 · «Klarte ikkje å sende oppdatering» → «Klarte ikke å sende oppdatering»
- [x] L96 · «Klarte ikkje å lese skjemadata» → «Klarte ikke å lese skjemadata»
- [x] L102 · «Arrangement-ID manglar» → «Arrangement-ID mangler»
- [x] L107 · «Klarte ikkje å oppdatere typen for arrangementet i databasen» → «Klarte ikke å oppdatere typen for arrangementet i databasen»
- [x] L111 · «Klarte ikkje å sende oppdatering» → «Klarte ikke å sende oppdatering»
- [x] L124 · «Klarte ikkje å lese skjemadata» → «Klarte ikke å lese skjemadata»
- [x] L130 · «Arrangement-ID manglar» → «Arrangement-ID mangler»
- [x] L135 · «Klarte ikkje å oppdatere systemet for arrangementet i databasen» → «Klarte ikke å oppdatere systemet for arrangementet i databasen»
- [x] L139 · «Klarte ikkje å sende oppdatering» → «Klarte ikke å sende oppdatering»
- [x] L153 · «Klarte ikkje å lese skjemadata» → «Klarte ikke å lese skjemadata»
- [x] L159 · «Arrangement-ID manglar» → «Arrangement-ID mangler»
- [x] L165 · «Klarte ikkje å oppdatere skildringa for arrangementet i databasen» → «Klarte ikke å oppdatere beskrivelsen for arrangementet i databasen»
- [x] L170 · «Klarte ikkje å sende oppdatering» → «Klarte ikke å sende oppdatering»
- [x] L202 · «Bruk mellom 3 og 50 teikn» → «Bruk mellom 3 og 50 tegn»
- [x] L248 · «Arrangementbilete» → «Arrangementsbilde»
- [x] L258 · «Finn eit bilete som presenterer arrangementet godt og» → «Finn et bilde som presenterer arrangementet godt og»
- [x] L259 · «gjer det lettare å kjenne att arrangementet.» → «gjør det lettere å kjenne igjen arrangementet.»
- [x] L269 · «Last opp bilete» → «Last opp bilde»
- [x] L272 · «Bytt bilete» → «Bytt bilde»
- [x] L295 · «Kva system er det?» → «Hvilket system er det?»
- [x] L298 · «Bruk minst 2 teikn» → «Bruk minst 2 tegn»
- [x] L303 · «Skildring av arrangementet» → «Beskrivelse av arrangementet»
- [x] L312 · «Skildring av arrangementet» → «Beskrivelse av arrangementet»

### [components/formsubmission/event_img_upload/event_img_upload.templ](../components/formsubmission/event_img_upload/event_img_upload.templ)

- [x] L30 · «Arrangement-ID manglar. Fekk:» → «Arrangement-ID mangler. Fikk:»
- [x] L37 · «Legg til bilete» → «Legg til bilde»
- [x] L53 · «Arrangement-ID manglar» → «Arrangement-ID mangler»
- [x] L61 · «Klarte ikkje å hente bilete frå skjemaet» → «Klarte ikke å hente bildet fra skjemaet»
- [x] L67 · «Klarte ikkje å erstatte biletet» → «Klarte ikke å erstatte bildet»
- [x] L75 · «Klarte ikkje å lagre biletet» → «Klarte ikke å lagre bildet»
- [x] L82 · «Klarte ikkje å lagre biletet» → «Klarte ikke å lagre bildet»
- [x] L87 · «Klarte ikkje å oppdatere revisjonsspor for arrangementsbiletet» → «Klarte ikke å oppdatere revisjonsspor for arrangementsbildet»
- [x] L119 · «Arrangement-ID manglar» → «Arrangement-ID mangler»
- [x] L125 · «Ugyldig bilettype» → «Ugyldig bildetype»
- [x] L132 · «Klarte ikkje å hente bilete frå skjemaet» → «Klarte ikke å hente bildet fra skjemaet»
- [x] L141 · «Klarte ikkje å lagre biletet» → «Klarte ikke å lagre bildet»
- [x] L148 · «Klarte ikkje å lagre biletet» → «Klarte ikke å lagre bildet»
- [x] L153 · «Klarte ikkje å oppdatere revisjonsspor for arrangementsbiletet» → «Klarte ikke å oppdatere revisjonsspor for arrangementsbildet»
- [x] L237 · «Klarte ikkje å hente arrangement:» → «Klarte ikke å hente arrangement:»
- [x] L242 · «Fann ikkje arrangementet» → «Fant ikke arrangementet»
- [x] L246 · «Lagra» → «Lagret»  _(borderline)_
- [x] L428 · «Heim» → «Hjem»
- [x] L431 · «Legg til bilete» → «Legg til bilde»
- [x] L459 · «Legg til bilete» → «Legg til bilde»
- [x] L467 · «Ver venleg og ikkje bruk AI/KI-genererte bilete» → «Vær vennlig og ikke bruk AI/KI-genererte bilder»
- [x] L470 · «Sidan AI/KI-genererte bilete er eit betent tema, ønskjer vi i RegnCon ikkje» → «Siden AI/KI-genererte bilder er et betent tema, ønsker vi i RegnCon ikke»
- [x] L471 · «at arrangørar skal bruke slike bilete som hovudbilete for arrangementet sitt.» → «at arrangører skal bruke slike bilder som hovedbilde for arrangementet sitt.»
- [x] L474 · «Om du er i tvil om bruk av bilete, anbefaler vi å finne bilete som er ope» → «Om du er i tvil om bruk av bilder, anbefaler vi å finne bilder som er åpent»
- [x] L475 · «tilgjengelege for bruk, eller bilete som elles blir brukte til å demonstrere» → «tilgjengelige for bruk, eller bilder som ellers blir brukt til å demonstrere»
- [x] L476 · «og promotere spelet/systemet du skal bruke på arrangementet.» → «og promotere spillet/systemet du skal bruke på arrangementet.»
- [x] L483 · «Last opp bilete» → «Last opp bilde»
- [x] L517 · «Lastar bilete. Ver venleg og vent...» → «Laster bilde. Vær vennlig og vent...»
- [x] L518 · «Klikk eller dra og slipp for å bytte bilete» → «Klikk eller dra og slipp for å bytte bilde»
- [x] L529 · «$_uploadProgress === 100 ? 'Lastar bilete. Ver venleg og vent...' : ''» → «$_uploadProgress === 100 ? 'Laster bilde. Vær vennlig og vent...' : ''»
- [x] L549 · «Lastar bilete» → «Laster bilde»
- [x] L562 · «Lastar bilete» → «Laster bilde»

### [components/formsubmission/feedback_message.go](../components/formsubmission/feedback_message.go)

- [x] L4 · «Klarte ikkje å lagre endringa. Prøv igjen. Kontakt styret dersom problemet held fram.» → «Klarte ikke å lagre endringen. Prøv igjen. Kontakt styret dersom problemet vedvarer.»
- [x] L5 · «Klarte ikkje å lagre endringa. Prøv igjen. Sjekk logger dersom problemet held fram.» → «Klarte ikke å lagre endringen. Prøv igjen. Sjekk logger dersom problemet vedvarer.»

### [components/formsubmission/puljefordeling.templ](../components/formsubmission/puljefordeling.templ)

- [x] L24 · «Arrangement-ID manglar» → «Arrangement-ID mangler»
- [x] L33 · «Klarte ikkje å lese skjemadata» → «Klarte ikke å lese skjemadata»
- [x] L39 · «Klarte ikkje å oppdatere puljefordeling» → «Klarte ikke å oppdatere puljefordeling»
- [x] L45 · «Klarte ikkje å sende oppdatering» → «Klarte ikke å sende oppdatering»
- [x] L60 · «Arrangement-ID manglar» → «Arrangement-ID mangler»
- [x] L67 · «Pulje-ID manglar» → «Pulje-ID mangler»
- [x] L81 · «Klarte ikkje å lese skjemadata» → «Klarte ikke å lese skjemadata»
- [x] L101 · «Klarte ikkje å oppdatere arrangementet i puljen» → «Klarte ikke å oppdatere arrangementet i puljen»
- [x] L107 · «Klarte ikkje å sende oppdatering» → «Klarte ikke å sende oppdatering»
- [x] L122 · «Arrangement-ID manglar» → «Arrangement-ID mangler»
- [x] L129 · «Pulje-ID manglar» → «Pulje-ID mangler»
- [x] L140 · «Klarte ikkje å lese skjemadata» → «Klarte ikke å lese skjemadata»
- [x] L147 · «Klarte ikkje å oppdatere rommet for puljen» → «Klarte ikke å oppdatere rommet for puljen»
- [x] L153 · «Klarte ikkje å sende oppdatering» → «Klarte ikke å sende oppdatering»
- [x] L235 · «Klarte ikkje å hente arrangementspulje:» → «Klarte ikke å hente arrangementspulje:»
- [x] L240 · «Klarte ikkje å hente rom» → «Klarte ikke å hente rom»
- [x] L287 · «Ikkje tildelt» → «Ikke tildelt»
- [x] L326 · «Klarte ikkje å hente puljer:» → «Klarte ikke å hente puljer:»

### [static/js/error_feedback.js](../static/js/error_feedback.js)

- [x] L6 · «"Klarte ikkje å lagre endringa. Prøv igjen."» → «Klarte ikke å lagre endringen. Prøv igjen.»

### [static/web_components/banner_cropper.js](../static/web_components/banner_cropper.js)

- [x] L35 · «'Klarte ikkje å lagre endringa. Prøv igjen. Kontakt styret dersom problemet held fram.'» → «Klarte ikke å lagre endringen. Prøv igjen. Kontakt styret dersom problemet vedvarer.»

## Batch C – Admin, modeller, tjenester  (83 / 83 ✅)

### [models/billettholder.go](../models/billettholder.go)

- [x] L87 · «spelar» → «spiller»

### [models/event-model.go](../models/event-model.go)

- [x] L20 · «Godkjend» → «Godkjent»
- [x] L22 · «Forkasta» → «Forkastet»  _(borderline)_
- [x] L40 · «Rollespel» → «Rollespill»
- [x] L42 · «Brettspel» → «Brettspill»
- [x] L44 · «Kortspel» → «Kortspill»
- [x] L46 · «Anna» → «Annet»
- [x] L67 · «Vaksne (18+)» → «Voksne (18+)»
- [x] L76 · «Barnevennleg» → «Barnevennlig»
- [x] L78 · «Eigna for vaksne» → «Egnet for voksne»
- [x] L104 · «Vanleg pulje» → «Vanlig pulje»
- [x] L106 · «Kortare (2-3 timar)» → «Kortere (2-3 timer)»
- [x] L108 · «Lengre (6+ timar)» → «Lengre (6+ timer)»
- [x] L117 · «Varer under 3 timar» → «Varer under 3 timer»
- [x] L119 · «Varer over 6 timar» → «Varer over 6 timer»

### [models/interests-model.go](../models/interests-model.go)

- [x] L25 · «Ikkje interessert» → «Ikke interessert»

### [pages/admin/admin.go](../pages/admin/admin.go)

- [x] L129 · «Finner ikkje rom» → «Finner ikke rom»

### [pages/admin/admin_page.templ](../pages/admin/admin_page.templ)

- [x] L74 · «Dra spel inn i puljane og sjå kor mange spel, 18+ og nybegynnarvennlege kvar pulje har. Åtvarar om ein spelleiar har to spel i same pulje.» → «Dra spill inn i puljene og se hvor mange spill, 18+ og nybegynnervennlige hver pulje har. Advarer om en spilleder har to spill i samme pulj…»
- [x] L82 · «Her kan du ligge til, modifisere …» → «Her kan du legge til, modifisere …»  _(dialekt)_

### [pages/admin/approval/editForm/edit_form_index_test.go](../pages/admin/approval/editForm/edit_form_index_test.go)

- [x] L15 · «Gitt redigeringssida for eit arrangement.» → «Gitt redigeringssiden for et arrangement.»
- [x] L16 · «Når sida blir rendra.» → «Når siden blir rendret.»  _(borderline)_
- [x] L17 · «Så inneheld ho ikkje lenger spelartildeling.» → «Så inneholder den ikke lenger spillertildeling.»

### [pages/admin/billettholder_admin/billettholder_admin_page.templ](../pages/admin/billettholder_admin/billettholder_admin_page.templ)

- [x] L85 · «Bilettholdere» → «Billettholdere»  _(skrivefeil)_

### [pages/admin/puljefordeling_tab_assignment_test.go](../pages/admin/puljefordeling_tab_assignment_test.go)

- [x] L92 · «Gitt ei SQLite-database med berre éin tilgjengeleg tilkopling.» → «Gitt en SQLite-database med bare én tilgjengelig tilkobling.»
- [x] L93 · «Når dialogen lastar interesser og spelarleiarar.» → «Når dialogen laster interesser og spilledere.»
- [x] L94 · «Så blir kvar spørjing ferdig før den neste bruker tilkoplinga.» → «Så blir hver spørring ferdig før den neste bruker tilkoblingen.»
- [x] L175 · «Når en administrator endrer interessa i puljefordeling.» → «Når en administrator endrer interessen i puljefordeling.»
- [x] L176 · «Så blir interessa oppdatert.» → «Så blir interessen oppdatert.»
- [x] L201 · «Gitt en åpen pulje med interesse for eitt arrangement.» → «Gitt en åpen pulje med interesse for ett arrangement.»
- [x] L202 · «Når administratoren vel Ikkje interessert.» → «Når administratoren velger Ikke interessert.»
- [x] L203 · «Så blir berre interessa for det arrangementet fjerna.» → «Så blir bare interessen for det arrangementet fjernet.»
- [x] L224 · «Gitt en publisert pulje med ei interesse.» → «Gitt en publisert pulje med en interesse.»
- [x] L225 · «Når administratoren prøver å fjerne interessa.» → «Når administratoren prøver å fjerne interessen.»
- [x] L226 · «Så blir endringa avvist og interessa står urørt.» → «Så blir endringen avvist og interessen står urørt.»
- [x] L274 · «Ikkje interessert» → «Ikke interessert»
- [x] L503 · «VALUES ('evA','Voksenspel',…)» → «VALUES ('evA','Voksenspill',…)»

### [pages/admin/puljefordeling_tab_test.go](../pages/admin/puljefordeling_tab_test.go)

- [x] L318 · «VALUES ('evA','Voksenspel',…)» → «VALUES ('evA','Voksenspill',…)»
- [x] L332 · «Gitt ein deltakar under 18 som er manuelt plassert i eit 18+-arrangement.» → «Gitt en deltaker under 18 som er manuelt plassert i et 18+-arrangement.»
- [x] L334 · «Så skal flisa merkast med «Under 18».» → «Så skal flisen merkes med «Under 18».»
- [x] L379 · «Gitt ein spelleiar under 18 på eit 18+-arrangement.» → «Gitt en spilleder under 18 på et 18+-arrangement.»
- [x] L381 · «Så skal spelleiar-lina merkast med «Under 18».» → «Så skal spilleder-linjen merkes med «Under 18».»
- [x] L387 · «VALUES ('evA','Voksenspel',…)» → «VALUES ('evA','Voksenspill',…)»
- [x] L408 · «VALUES ('evA','Voksenspel',…)» → «VALUES ('evA','Voksenspill',…)»

### [pages/admin/puljeoppsett.templ](../pages/admin/puljeoppsett.templ)

- [x] L43 · «Klarte ikkje å byggje brettet:» → «Klarte ikke å bygge brettet:»
- [x] L47 · «%d spelleiar-kollisjon(ar)» → «%d spilleder-kollisjon(er)»
- [x] L64 · «%d spel» → «%d spill»
- [x] L84 · «%d spel · %d 18+ · %d nyb.» → «%d spill · %d 18+ · %d nyb.»
- [x] L90 · «spel i denne puljen» → «spill i denne puljen»
- [x] L90 · «⚠ Spelleiar-kollisjon —» → «⚠ Spilleder-kollisjon —»
- [x] L219 · «Arrangement-ID manglar» → «Arrangement-ID mangler»
- [x] L229 · «Klarte ikkje å oppdatere puljen» → «Klarte ikke å oppdatere puljen»

### [pages/admin/puljeoppsett_render_test.go](../pages/admin/puljeoppsett_render_test.go)

- [x] L18 · «Ledig spel» → «Ledig spill»
- [x] L23 · «1 spel» → «1 spill»
- [x] L23 · «Ledig spel» → «Ledig spill»
- [x] L89 · «Spelleiar-kollisjon» → «Spilleder-kollisjon»

### [pages/admin/puljeoppsett_route_test.go](../pages/admin/puljeoppsett_route_test.go)

- [x] L25 · «Spel» → «Spill»
- [x] L63 · «Spel» → «Spill»

### [pages/admin/rooms/rooms_assignment_page.templ](../pages/admin/rooms/rooms_assignment_page.templ)

- [x] L235 · «Klarte ikkje å tildele rommet. Prøv igjen.» → «Klarte ikke å tildele rommet. Prøv igjen.»
- [x] L525 · «Romfordelig for …» → «Romfordeling for …»  _(skrivefeil)_
- [x] L549 · «vær den førte til å … ligge inn …» → «vær den første til å … legge inn …»  _(dialekt/skrivefeil)_

### [pages/admin/rooms/rooms_page.templ](../pages/admin/rooms/rooms_page.templ)

- [x] L92 · «Er du sikker på at du vill slette dette rommet?» → «Er du sikker på at du vil slette dette rommet?»  _(skrivefeil)_
- [x] L291 · «Ingen rom funnet, vær den førte til å et rom nå!» → «Ingen rom funnet, vær den første til å legge inn et rom nå!»  _(skrivefeil)_
- [x] L298 · «Ligg til nytt rom» → «Legg til nytt rom»  _(dialekt)_

### [pages/admin/rooms/rooms_page_test.go](../pages/admin/rooms/rooms_page_test.go)

- [x] L53 · «"Ligg til nytt rom"» → «"Legg til nytt rom"»  _(dialekt, følger rooms_page.templ 298)_
- [x] L92 · «"Romfordelig for Fredag kveld"» → «"Romfordeling for Fredag kveld"»  _(skrivefeil, følger rooms_assignment_page.templ 525)_

### [service/puljefordeling/emulate_test.go](../service/puljefordeling/emulate_test.go)

- [x] L251 · «Vaksen» → «Voksen»  _(borderline)_
- [x] L270 · «Vaksen Voksdal» → «Voksen Voksdal»  _(borderline)_

### [service/puljefordeling/tildelinger.go](../service/puljefordeling/tildelinger.go)

- [x] L166 · «fjern %s-tildeling for billettholder %d frå %s i %s: %w» → «fjern %s-tildeling for billettholder %d fra %s i %s: %w»
- [x] L170 · «les fjerna %s-tildeling: %w» → «les fjernet %s-tildeling: %w»  _(borderline)_
- [x] L318 · «%w: dra-og-slipp kan berre flytte ein spelar» → «%w: dra-og-slipp kan bare flytte én spiller»
- [x] L367 · «Bruk Legg til for å plassere ein GM som spelar.» → «Bruk Legg til for å plassere en GM som spiller.»
- [x] L396 · «Legg til som spelar på «%s»» → «Legg til som spiller på «%s»»
- [x] L403 · «Legg til som spelar på «%s»» → «Legg til som spiller på «%s»»

### [service/puljefordeling/tildelinger_test.go](../service/puljefordeling/tildelinger_test.go)

- [x] L59 · «Legg til som spelar på «Alpha»» → «Legg til som spiller på «Alpha»»
- [x] L141 · «kan-ikkje-overstyre» → «kan-ikke-overstyre»
- [x] L485 · «Tilskodar» → «Tilskuer»

### [service/rooms/rooms_status_test.go](../service/rooms/rooms_status_test.go)

- [x] L24 · «Laurdag» → «Lørdag»

### [service/rooms/rooms_test.go](../service/rooms/rooms_test.go)

- [x] L24 · «Dette er eit gyldig rom» → «Dette er et gyldig rom»

### [service/rooms/rooms_update_test.go](../service/rooms/rooms_update_test.go)

- [x] L25 · «Dette er ei oppdatert note» → «Dette er en oppdatert note»

### [service/rooms/rooms_validation.go](../service/rooms/rooms_validation.go)

- [x] L16 · «Rom namn kan ikkje berre innehalde mellomrom» → «Romnavn kan ikke bare inneholde mellomrom»
- [x] L22 · «Namn kan ikkje vere lengre enn 50 teikn» → «Navn kan ikke være lengre enn 50 tegn»
- [x] L27 · «Romnummer er påkravd» → «Romnummer er påkrevd»
- [x] L33 · «Rom nummer kan ikkje vere lengre enn 10 teikn» → «Romnummer kan ikke være lengre enn 10 tegn»

## Batch D – Sider, profil, tester, dokumentasjon  (62 / 62 ✅)

### [components/header/menu.templ](../components/header/menu.templ)

- [x] L154 · «Regncon framside» → «Regncon forside»  _(borderline)_
- [x] L626 · «Vanlege Spørsmål» → «Vanlige Spørsmål»
- [x] L881 · «Vanlege Spørsmål» → «Vanlige Spørsmål»

### [components/header/menu_test.go](../components/header/menu_test.go)

- [x] L36 · «Så skal brukeren bare få navigasjonslenker til forsiden, egen profil, utlogging og vanlege spørsmål.» → «Så skal brukeren bare få navigasjonslenker til forsiden, egen profil, utlogging og vanlige spørsmål.»
- [x] L68 · «Så skal brukeren få navigasjonslenker til forsiden, egen profil, utlogging, adminområdene og vanlege spørsmål.» → «Så skal brukeren få navigasjonslenker til forsiden, egen profil, utlogging, adminområdene og vanlige spørsmål.»

### [components/profile/my_program.templ](../components/profile/my_program.templ)

- [x] L147 · «Programmet for Regncon er ikkje publisert enno» → «Programmet for Regncon er ikke publisert ennå»
- [x] L149 · «Vi arbeider framleis med planlegginga og vil leggje ut programmet så snart det er klart.» → «Vi arbeider fortsatt med planleggingen og vil legge ut programmet så snart det er klart.»
- [x] L150 · «Følg med på nettsida for oppdateringar og nyheiter.» → «Følg med på nettsiden for oppdateringer og nyheter.»

### [components/profile/my_program_render_test.go](../components/profile/my_program_render_test.go)

- [x] L125 · «Programmet for Regncon er ikkje publisert enno» → «Programmet for Regncon er ikke publisert ennå»

### [components/profile/my_program_test.go](../components/profile/my_program_test.go)

- [x] L224 · «Gitt en manuell spillerplassering i ei open pulje.» → «Gitt en manuell spillerplassering i en åpen pulje.»
- [x] L225 · «Når festivalprogrammet blir lasta.» → «Når festivalprogrammet blir lastet.»  _(borderline)_
- [x] L226 · «Så viser programmet arrangementet med ein gong, slik varselet i interessedialogen gjer.» → «Så viser programmet arrangementet med en gang, slik varselet i interessedialogen gjør.»
- [x] L258 · «Gitt en manuell spillerplassering i ei låst pulje.» → «Gitt en manuell spillerplassering i en låst pulje.»
- [x] L259 · «Når festivalprogrammet blir lasta.» → «Når festivalprogrammet blir lastet.»  _(borderline)_
- [x] L260 · «Så viser programmet arrangementet sjølv om puljefordelinga ikkje er ferdig.» → «Så viser programmet arrangementet selv om puljefordelingen ikke er ferdig.»

### [components/profile/my_tickets.templ](../components/profile/my_tickets.templ)

- [x] L29 · «Billettar» → «Billetter»

### [components/timeschedule.templ](../components/timeschedule.templ)

- [x] L193 · «For de som har kjøpt middagsbilett» → «For de som har kjøpt middagsbillett»  _(skrivefeil)_

### [documentation/testing/admin-approval.md](../documentation/testing/admin-approval.md)

- [x] L66 · «Godkjenningsflyten er arbeidsbar på ulike skjermer» → «Godkjenningsflyten er brukbar på ulike skjermer»  _(avgjort: brukbar)_
- [x] L69 · «… fortsatt være lesbar og arbeidsbar.» → «… fortsatt være lesbar og brukbar.»  _(avgjort: brukbar)_

### [documentation/testing/general.md](../documentation/testing/general.md)

- [x] L97 · «**Når** brukreren oppgir et nytt passord.» → «**Når** brukeren oppgir et nytt passord.»  _(skrivefeil)_

### [documentation/testing/root.md](../documentation/testing/root.md)

- [x] L62 · «… eller gjømme viktig informasjon.» → «… eller gjemme viktig informasjon.»  _(dialekt)_

### [domeneordbok.md](../domeneordbok.md)

- [x] L23 · «… som skal spiller innen for tidspunktet som styret har valg …» → «… som skal spilles innenfor tidspunktet styret har valgt …»  _(grammatikk)_
- [x] L45 · «… som alle så vil kan melde seg på.Som for eksempel …» → «… som alle som vil, kan melde seg på. For eksempel …»  _(grammatikk)_
- [x] L49 · «… til og vise html» → «… til å vise html»  _(grammatikk)_

### [pages/event/event.go](../pages/event/event.go)

- [x] L42 · «Du har ikkje tilgang til å endre interessa til denne billettheldaren. Kontakt styret.» → «Du har ikke tilgang til å endre interessen til denne billettholderen. Kontakt styret.»
- [x] L45 · «Denne pulja er ikkje tilgjengeleg for dette arrangementet.» → «Denne puljen er ikke tilgjengelig for dette arrangementet.»
- [x] L48 · «Pulja er låst. Du kan ikkje melde eller endre interesse lenger medan vi fordeler spelarar.» → «Puljen er låst. Du kan ikke melde eller endre interesse lenger mens vi fordeler spillere.»
- [x] L51 · «Puljefordelinga er klar. Gå til profilen din for å sjå kva du fekk.» → «Puljefordelingen er klar. Gå til profilen din for å se hva du fikk.»
- [x] L56 · «Det oppstod ein feil då interessa skulle lagrast. Prøv igjen, eller kontakt styret dersom feilen held fram.» → «Det oppstod en feil da interessen skulle lagres. Prøv igjen, eller kontakt styret dersom feilen fortsetter.»
- [x] L249 · «Vel billetthelder før du melder interesse.» → «Velg billettholder før du melder interesse.»
- [x] L256 · «Vel pulje før du melder interesse.» → «Velg pulje før du melder interesse.»

### [pages/event/event_interest_test.go](../pages/event/event_interest_test.go)

- [x] L522 · «Gitt at ein billettholder under 18 år er tildelt eit arrangement i puljen.» → «Gitt at en billettholder under 18 år er tildelt et arrangement i puljen.»
- [x] L523 · «Når eit 18-års arrangement i same pulje vert vist.» → «Når et 18-års arrangement i samme pulje blir vist.»
- [x] L524 · «Så skal berre aldersvarselet vises, ikkje tildelingsvarselet.» → «Så skal bare aldersvarselet vises, ikke tildelingsvarselet.»
- [x] L592 · «Vel billetthelder» → «Velg billettholder»
- [x] L592 · «Vel pulje» → «Velg pulje»

### [pages/event/event_interest_update_test.go](../pages/event/event_interest_update_test.go)

- [x] L83 · «… men den gamle publiseringsflagget står av.» → «… men det gamle publiseringsflagget står av.»  _(skrivefeil)_

### [pages/event/event_page_admin_test.go](../pages/event/event_page_admin_test.go)

- [x] L58 · «Denne pulja er ikkje tilgjengeleg for dette arrangementet.» → «Denne puljen er ikke tilgjengelig for dette arrangementet.»
- [x] L59 · «Pulja er låst. Du kan ikkje melde eller endre interesse lenger medan vi fordeler spelarar.» → «Puljen er låst. Du kan ikke melde eller endre interesse lenger mens vi fordeler spillere.»
- [x] L60 · «Puljefordelinga er klar. Gå til profilen din for å sjå kva du fekk.» → «Puljefordelingen er klar. Gå til profilen din for å se hva du fikk.»

### [pages/event/event_visibility_test.go](../pages/event/event_visibility_test.go)

- [x] L419 · «… men den gamle puljeflagget står av.» → «… men det gamle puljeflagget står av.»  _(skrivefeil)_

### [pages/login/login.go](../pages/login/login.go)

- [x] L109 · «Velkomen tilbake til Regncon 2026!» → «Velkommen tilbake til Regncon 2026!»

### [pages/login/login.templ](../pages/login/login.templ)

- [x] L22 · «Velkomen tilbake!» → «Velkommen tilbake!»
- [x] L23 · «Det ser ut til at du allereie er innlogga og du treng ikkje logge inn på nytt.» → «Det ser ut til at du allerede er innlogget, og du trenger ikke logge inn på nytt.»
- [x] L24 · «Du vert sendt vidare til framsida om» → «Du blir sendt videre til forsiden om»
- [x] L24 · «sekund.» → «sekunder.»  _(borderline)_
- [x] L25 · «Om vidaresendinga ikkje fungerer, kan du bruke denna lenkja til» → «Om videresendingen ikke fungerer, kan du bruke denne lenken til»
- [x] L25 · «framsida» → «forsiden»  _(borderline)_

### [pages/print-friendly/print-friendly-page.templ](../pages/print-friendly/print-friendly-page.templ)

- [x] L27 · «Regncon program uttriftsvennlig» → «Regncon program utskriftsvennlig»  _(skrivefeil)_

### [pages/profile/newevent/new_page.templ](../pages/profile/newevent/new_page.templ)

- [x] L50 · «Ta kontakt med RegnCon styret på Dicord …» → «Ta kontakt med RegnCon-styret på Discord …»  _(skrivefeil)_

### [service/checkIn/assign_users_test.go](../service/checkIn/assign_users_test.go)

- [x] L80 · «Gitt at ein billettholder har fått lagt til ei manuell e-postadresse, og ein eksisterande brukar har same e-postadresse med annan casing.» → «Gitt at en billettholder har fått lagt til en manuell e-postadresse, og en eksisterende bruker har samme e-postadresse med annen casing.»
- [x] L81 · «Når e-postadressa blir forsona mot brukarar.» → «Når e-postadressen blir avstemt mot brukere.»
- [x] L82 · «Så skal billettholderen få ei varig brukar-tilknyting.» → «Så skal billettholderen få en varig brukertilknytning.»
- [x] L110 · «Gitt at ein billettholder allereie er knytt til ein brukar via ei manuell e-postadresse.» → «Gitt at en billettholder allerede er knyttet til en bruker via en manuell e-postadresse.»
- [x] L111 · «Når same e-postforsoning køyrer på nytt.» → «Når samme e-postavstemming kjører på nytt.»
- [x] L112 · «Så skal det framleis berre finnast ei brukar-tilknyting.» → «Så skal det fortsatt bare finnes én brukertilknytning.»
- [x] L141 · «Gitt at ei manuell e-postadresse er fjerna frå ein billettholder, og ingen attverande e-postadresser på billettholderen samsvarer med bruka…» → «Gitt at en manuell e-postadresse er fjernet fra en billettholder, og ingen gjenværende e-postadresser på billettholderen samsvarer med bruk…»
- [x] L142 · «Når e-postadressa blir forsona mot brukar-tilknytingar.» → «Når e-postadressen blir avstemt mot brukertilknytninger.»
- [x] L143 · «Så skal den varige brukar-tilknytinga fjernast.» → «Så skal den varige brukertilknytningen fjernes.»
- [x] L174 · «Gitt at ei manuell e-postadresse er fjerna frå ein billettholder, men ei anna attverande e-postadresse på same billettholder framleis samsv…» → «Gitt at en manuell e-postadresse er fjernet fra en billettholder, men en annen gjenværende e-postadresse på samme billettholder fortsatt sa…»
- [x] L175 · «Når e-postadressa blir forsona mot brukar-tilknytingar.» → «Når e-postadressen blir avstemt mot brukertilknytninger.»
- [x] L176 · «Så skal den varige brukar-tilknytinga behaldast.» → «Så skal den varige brukertilknytningen beholdes.»

## Ikke endre (bevisst beholdt)

- `https://www.regncon.no/vanlege-sporsmal/`: ekstern lenke (`components/header/menu.templ`, `pages/event/event_interest_panel.templ` og tester)
- «verken», «sidene», «Attende» (ordenstall), «Inga» (navn): gyldig bokmål
- Domeneord fra `domeneordbok.md`: Billettholder, Pulje, Interesse, Førstevalg, Påmelding, Lordag, Sondag

## Sluttkontroll

- [x] Rester på tvers av batcher er rettet
- [x] Overflødige cSpell-ord fjernet fra `.vscode/settings.json` (bilete, Arrangementbilete, gjer, interessa, lagrast, lettare, pulja, tilgjengeleg)
- [x] `go tool templ generate`
- [x] `go build ./...`
- [x] `go run ./cmd/testreport` er grønn
- [x] Ny nynorsk-skanning: bare godkjente rester
- [x] Diffen er gjennomgått

### Resultat

- 56 filer endret, bare tekst i strenger, maler, tester og dokumentasjon (ingen kode, identifikatorer eller lagrede verdier).
- `go run ./cmd/testreport`: 0 feilede tester, 0 feilede pakker.
- Sluttskanning av 40 476 tekster: eneste gjenværende nynorsk er den eksterne lenken `regncon.no/vanlege-sporsmal/` (beholdt med vilje).
- Neste steg senere: CI-sjekk (`cmd/nynorskcheck` med tillatelsesliste) så ny nynorsk ikke sniker seg inn.
