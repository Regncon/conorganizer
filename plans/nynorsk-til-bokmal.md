# Kartlegging: nynorsk → bokmål

Fullstendig oversikt over nynorsk i koden (kartlagt på commit 5887e4ac, kontrollert på nytt mot 8e0d0369). Målet er å gjøre om alt til bokmål.

**310 forekomster i 43 filer.** 127 vises garantert for brukeren, 21 kanskje, 89 gjør det ikke (feilmeldinger Datastar ignorerer), 72 er testtekst.

Ingen av de kartlagte filene er endret mellom 5887e4ac og 8e0d0369, så linjenumrene stemmer fortsatt.

## Hvordan kartleggingen ble gjort

Ikke søkeord, men full gjennomgang:

1. En parser (Go `go/ast` + templ sin egen parser) hentet ut **hver** tekstbit fra alle `.go`- og `.templ`-filer: string-literaler, kommentarer, tekstnoder og attributtverdier. JS, SQL, CSS og konfigurasjon ble hentet ut med regex (string-literaler, kommentarer, SQL-verdier).
2. Det ga 14 066 tekstbiter, 7 405 unike. Hver unike tekst ble lest og vurdert – ikke bare de som traff nynorskord.
3. Treffene ble koblet tilbake til alle fil/linje-steder de forekommer.
4. Egen sporing av alle veier backendtekst når frontend (`http.Error`, SSE-signaler, `PatchElementTempl`, toast/errorfeedback, validering, `Label()`, `{ err.Error() }` i maler).
5. Kryssjekket mot en tidligere søkeordbasert runde: ingen nynorsk som søkeordrunden fant, mangler her.

**Kontroll i etterkant (mot 8e0d0369):**

- Ny uttrekking fra HEAD: all ny tekst fra de 16 nye commitene er bokmål.
- Nynorsk ord- og formsjekk over **hver linje i hver sporet fil**, uten forhåndsfilter (dekker også JS-maler over flere linjer og konfigfiler). Fant én bom: testtittelen «Voksenspel» (4 steder, lagt til under).
- Uavhengig ny gjennomlesing av alle 1 350 norskliknende tekster som første runde vurderte som bokmål: 0 nynorsk.

- «Flagg, ikke avgjør»-runde: alle 15 249 unike tekster (ingen forhåndsfilter) lest på nytt med tre utfall: nynorsk / usikker / klart ikke. Alt usikkert er løftet til avgjørelse (se «Flagget som usikker – avgjort»). Lange tekster (> 500 tegn) er sjekket helt ut. Ingen ny nynorsk utover «Voksenspel».

**Ikke gjennomgått:** `.ai/threads/*` (AI-samtalelogger), tekst i bilder, generert `*_templ.go` (gjenskapes fra `.templ`), `static/datastar.js`. Dokumentasjon (`documentation/`, `plans/`, `README.md`, `domeneordbok.md`) er gjennomgått separat: ingen nynorsk.

**Utenfor koden:** innhold i databasen (arrangementer, rom og puljenavn som admin har skrevet inn), billettnavn fra CheckIn og sider på regncon.no.

## Når teksten brukeren? (backend → frontend)

| Mekanisme | Nynorsk | Vises? |
|---|---|---|
| `http.Error` på Datastar `@put/@post/@delete` (formsubmission, puljeoppsett) | 84 | **Nei** – Datastar viser ikke body; brukeren ser `data-feedback-message` |
| `http.Error` i bildeopplasting (`event_img_upload.templ`) | 7 | **Ja** – `upload-with-progress.js` legger responsteksten i `$_uploadError` |
| `http.Error` i `event_new.templ` (eksempel-rute uten UI) | 3 | Kanskje |
| Modell-labels `Label()`/`BadgeLabel()` (`models/`) | 15 | **Ja** (unntatt `EventPlayerRole.Label()` «spelar», ubrukt) |
| `interestErrorMessage`-signal (`pages/event/event.go`) | 7 | **Ja** |
| Romvalidering `RoomFormErrors.AddError` (`service/rooms/rooms_validation.go`) | 4 | **Ja** |
| `FormTitle` «Finner ikkje rom» (`pages/admin/admin.go`) | 1 | **Ja** |
| `errorfeedback` standardmeldinger (`feedback_message.go`, `error_feedback.js`, `errorfeedback.Attrs`) | 4 | **Ja** |
| `Tildelingsvarsel.Handling` / `fmt.Sprintf` (`service/puljefordeling/tildelinger.go`) | 6 | **Ja** |
| Hjelpetekster, breadcrumbs, sidetitler, toast «Lagra», adminkort | 13 | **Ja** |
| Markup i `.templ` (tekst og attributter som placeholder/title/aria-label) | resten av `ui` | **Ja** |
| `fmt.Errorf` i `tildelinger.go` (166/170/318) | 3 | **Nei** – brukes bare i tester / blir til generisk «Ugyldig tildeling» |
| `banner_cropper.js` `UPLOAD_ERROR_MESSAGE` | 1 | Nei – ubrukt konstant |
| slog-loggmeldinger | 0 | – |

Alle `{ err.Error() }` i maler viser engelske feil – ingen nynorsk kommer den veien.

**Dynamisk tekst som et rent literal-søk kan gå glipp av** (ta med når dere endrer):
- `pages/event/event.go`: «Vel billetthelder f\u00f8r…» er skrevet med `\u00f8`-escape.
- `interestErrorMessageFromError` velger melding ved `strings.Contains` på engelsk feiltekst.
- `components/interest_indicator.templ` 23 bygger `"%s er %s"` + `InterestLevel.Label()` med små bokstaver, så tooltipen viser «ikkje interessert» til `models/interests-model.go` 25 er endret.
- `"Arrangement-ID manglar. Fekk: "+eventId` og `Sprintf("Klarte ikkje å oppdatere arrangementet: %v")` settes sammen ved kjøring.

## Ting å passe på ved omskriving

- **Tester som må endres sammen med produksjonsteksten** (eksakt match):
  - `pages/event/event.go` 45/48/51 ↔ `pages/event/event_page_admin_test.go` 58–60
  - `pages/event/event.go` 249/256 ↔ `pages/event/event_interest_test.go` 592 (negativ «inneholder ikke»-sjekk – blir stille grønn hvis bare produksjon endres)
  - `components/profile/my_program.templ` 147 ↔ `components/profile/my_program_render_test.go` 125
  - `models/interests-model.go` 25 («Ikkje interessert») ↔ `pages/admin/puljefordeling_tab_assignment_test.go` 274
  - `pages/admin/puljeoppsett.templ` 64/90 («spel», «Spelleiar-kollisjon») ↔ `pages/admin/puljeoppsett_render_test.go` 18/23/89
  - `service/puljefordeling/tildelinger.go` 396/403 («Legg til som spelar») ↔ `service/puljefordeling/tildelinger_test.go` 59
  - `service/puljefordeling/emulate_test.go` 251 ↔ 270 (testperson «Vaksen»)
- **Kjør `templ generate`** etter endringer i `.templ`-filer.
- **Ekstern lenke** `regncon.no/vanlege-sporsmal/` (`components/header/menu.templ`, `pages/event/event_interest_panel.templ` og tester) er nynorsk, men peker til regncon.no – ikke endre med mindre nettsiden endres.
- **cSpell-ord** i `.vscode/settings.json` som blir overflødige: bilete, Arrangementbilete, gjer, interessa, lagrast, lettare, pulja, tilgjengeleg (behold Lordag/Sondag).
- `models/event-model.go`: etter omskriving blir label for Approved/Archived («Godkjend»/«Forkasta») lik lagret verdi («Godkjent»/«Forkastet»), så de case-ene kan fjernes.
- **Ingen nynorsk er lagret i databasen** eller brukt som nøkkel/enum – ingen migrering trengs.

## Gjentatte fraser (enkle søk-og-erstatt)

| Nynorsk | Bokmål | Antall |
|---|---|---|
| Arrangement-ID manglar | Arrangement-ID mangler | 22 |
| Klarte ikkje å lese skjemadata | Klarte ikke å lese skjemadata | 20 |
| Klarte ikkje å sende oppdatering | Klarte ikke å sende oppdatering | 19 |
| Klarte ikkje å lagre biletet | Klarte ikke å lagre bildet | 4 |
| Legg til bilete | Legg til bilde | 3 |
| Pulje-ID manglar | Pulje-ID mangler | 2 |

## Flagget som usikker – avgjort

Ikke nynorsk, men heller ikke standard bokmål. **Avgjort 2026-09-23: rettes i samme PR-er som nynorsken**, med forslaget i tabellen («arbeidsbar» → «brukbar»).

| Sted | Ord | Vurdering | Forslag |
|---|---|---|---|
| `pages/admin/admin_page.templ` 82 | ligge til | dialekt/skrivefeil | legge til |
| `pages/admin/rooms/rooms_page.templ` 298 (↔ `rooms_page_test.go` 53) | Ligg til | dialekt/skrivefeil | Legg til |
| `pages/admin/rooms/rooms_assignment_page.templ` 549 | ligge inn, førte | dialekt/skrivefeil | legge inn, første |
| `pages/admin/rooms/rooms_page.templ` 291 | førte | skrivefeil, mangler også verb | første |
| `pages/admin/rooms/rooms_page.templ` 92 | vill | skrivefeil | vil |
| `documentation/testing/root.md` 62 | gjømme | dialekt/sideform | gjemme |
| `documentation/testing/admin-approval.md` 66, 69 | arbeidsbar | ikke et vanlig ord | brukbar |
| `pages/print-friendly/print-friendly-page.templ` 27 | uttriftsvennlig | skrivefeil | utskriftsvennlig |

Mekaniske treff vurdert som gyldig bokmål (ingen endring): «verken», «sidene», «Attende» (ordenstall), «Inga» (navn).

## Ikke nynorsk, men funnet underveis (skrivefeil)

- `pages/admin/rooms/rooms_assignment_page.templ` «Romfordelig»
- `documentation/testing/general.md` 97 «brukreren»; `domeneordbok.md` 23/45/49 grammatikk
- `pages/profile/newevent/new_page.templ` 50 «RegnCon styret på Dicord» → «RegnCon-styret på Discord»
- `pages/admin/billettholder_admin/billettholder_admin_page.templ` 85 «Bilettholdere» → «Billettholdere»
- `components/timeschedule.templ` 193 «middagsbilett» → «middagsbillett»
- `components/formsubmission/other_details.templ` 370 «å hjelp dine» → «å hjelpe dine», 373 «orginalt» → «opprinnelig»
- `pages/event/event_interest_update_test.go` 83 og `pages/event/event_visibility_test.go` 419 «den gamle …flagget» → «det gamle …flagget»

## Oversikt per mappe

| Mappe | Forekomster |
|---|---|
| `components/formsubmission` | 181 |
| `pages/admin` | 40 |
| `pages/event` | 16 |
| `models/event-model.go` | 14 |
| `service/checkIn` | 12 |
| `components/profile` | 11 |
| `service/puljefordeling` | 11 |
| `components/header` | 7 |
| `pages/login` | 7 |
| `service/rooms` | 7 |
| `models/billettholder.go` | 1 |
| `models/interests-model.go` | 1 |
| `static/js` | 1 |
| `static/web_components` | 1 |

## Per fil

Kolonner: **Type** = ui (synlig markup/label), http.Error, error, js, test, kommentar. **Vises?** = når teksten faktisk til brukeren (ja/nei/kanskje; – for test/kommentar). **Grad** = nynorsk eller borderline (radikal bokmålsform).

### [components/formsubmission/about_event.templ](../components/formsubmission/about_event.templ) (30)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 28 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å lese skjemadata | Klarte ikke å lese skjemadata |
| 34 | http.Error | nei | nynorsk | manglar | Arrangement-ID manglar | Arrangement-ID mangler |
| 39 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å oppdatere tittelen for arrangementet i databasen | Klarte ikke å oppdatere tittelen for arrangementet i databasen |
| 51 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å sende oppdatering | Klarte ikke å sende oppdatering |
| 65 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å lese skjemadata | Klarte ikke å lese skjemadata |
| 71 | http.Error | nei | nynorsk | manglar | Arrangement-ID manglar | Arrangement-ID mangler |
| 78 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å oppdatere introduksjonen for arrangementet i databasen | Klarte ikke å oppdatere introduksjonen for arrangementet i databasen |
| 82 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å sende oppdatering | Klarte ikke å sende oppdatering |
| 96 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å lese skjemadata | Klarte ikke å lese skjemadata |
| 102 | http.Error | nei | nynorsk | manglar | Arrangement-ID manglar | Arrangement-ID mangler |
| 107 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å oppdatere typen for arrangementet i databasen | Klarte ikke å oppdatere typen for arrangementet i databasen |
| 111 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å sende oppdatering | Klarte ikke å sende oppdatering |
| 124 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å lese skjemadata | Klarte ikke å lese skjemadata |
| 130 | http.Error | nei | nynorsk | manglar | Arrangement-ID manglar | Arrangement-ID mangler |
| 135 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å oppdatere systemet for arrangementet i databasen | Klarte ikke å oppdatere systemet for arrangementet i databasen |
| 139 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å sende oppdatering | Klarte ikke å sende oppdatering |
| 153 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å lese skjemadata | Klarte ikke å lese skjemadata |
| 159 | http.Error | nei | nynorsk | manglar | Arrangement-ID manglar | Arrangement-ID mangler |
| 165 | http.Error | nei | nynorsk | ikkje, skildringa | Klarte ikkje å oppdatere skildringa for arrangementet i databasen | Klarte ikke å oppdatere beskrivelsen for arrangementet i databasen |
| 170 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å sende oppdatering | Klarte ikke å sende oppdatering |
| 202 | ui | ja | nynorsk | teikn | Bruk mellom 3 og 45 teikn | Bruk mellom 3 og 45 tegn |
| 248 | ui | ja | nynorsk | Arrangementbilete | Arrangementbilete | Arrangementsbilde |
| 258 | ui | ja | nynorsk | eit, bilete | Finn eit bilete som presenterer arrangementet godt og | Finn et bilde som presenterer arrangementet godt og |
| 259 | ui | ja | nynorsk | gjer, lettare, kjenne att | gjer det lettare å kjenne att arrangementet. | gjør det lettere å kjenne igjen arrangementet. |
| 269 | ui | ja | nynorsk | bilete | Last opp bilete | Last opp bilde |
| 272 | ui | ja | nynorsk | bilete | Bytt bilete | Bytt bilde |
| 295 | ui | ja | nynorsk | Kva | Kva system er det? | Hvilket system er det? |
| 298 | ui | ja | nynorsk | teikn | Bruk minst 2 teikn | Bruk minst 2 tegn |
| 303 | ui | ja | nynorsk | Skildring | Skildring av arrangementet | Beskrivelse av arrangementet |
| 312 | ui | ja | nynorsk | Skildring | Skildring av arrangementet | Beskrivelse av arrangementet |

### [components/formsubmission/contact_info.templ](../components/formsubmission/contact_info.templ) (15)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 24 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å lese skjemadata | Klarte ikke å lese skjemadata |
| 30 | http.Error | nei | nynorsk | manglar | Arrangement-ID manglar | Arrangement-ID mangler |
| 35 | http.Error | nei | nynorsk | ikkje, namnet | Klarte ikkje å oppdatere namnet for arrangementet i databasen | Klarte ikke å oppdatere navnet for arrangementet i databasen |
| 39 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å sende oppdatering | Klarte ikke å sende oppdatering |
| 53 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å lese skjemadata | Klarte ikke å lese skjemadata |
| 59 | http.Error | nei | nynorsk | manglar | Arrangement-ID manglar | Arrangement-ID mangler |
| 64 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å oppdatere e-posten for arrangementet i databasen | Klarte ikke å oppdatere e-posten for arrangementet i databasen |
| 69 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å sende oppdatering | Klarte ikke å sende oppdatering |
| 82 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å lese skjemadata | Klarte ikke å lese skjemadata |
| 88 | http.Error | nei | nynorsk | manglar | Arrangement-ID manglar | Arrangement-ID mangler |
| 93 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å oppdatere telefonnummeret for arrangementet i databasen | Klarte ikke å oppdatere telefonnummeret for arrangementet i databasen |
| 98 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å sende oppdatering | Klarte ikke å sende oppdatering |
| 118 | ui | ja | nynorsk | Namnet | Namnet ditt | Navnet ditt |
| 121 | ui | ja | nynorsk | teikn | Bruk minst 2 teikn | Bruk minst 2 tegn |
| 158 | ui | ja | nynorsk | framfor | Bruk 8 siffer, eventuelt med +47 framfor | Bruk 8 siffer, eventuelt med +47 foran |

### [components/formsubmission/event_img_upload/event_img_upload.templ](../components/formsubmission/event_img_upload/event_img_upload.templ) (32)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 30 | http.Error | ja | nynorsk | manglar, Fekk | Arrangement-ID manglar. Fekk: | Arrangement-ID mangler. Fikk: |
| 37 | ui | ja | nynorsk | bilete | Legg til bilete | Legg til bilde |
| 53 | http.Error | ja | nynorsk | manglar | Arrangement-ID manglar | Arrangement-ID mangler |
| 61 | http.Error | ja | nynorsk | ikkje, bilete, frå | Klarte ikkje å hente bilete frå skjemaet | Klarte ikke å hente bildet fra skjemaet |
| 67 | http.Error | ja | nynorsk | ikkje, biletet | Klarte ikkje å erstatte biletet | Klarte ikke å erstatte bildet |
| 75 | http.Error | ja | nynorsk | ikkje, biletet | Klarte ikkje å lagre biletet | Klarte ikke å lagre bildet |
| 82 | http.Error | ja | nynorsk | ikkje, biletet | Klarte ikkje å lagre biletet | Klarte ikke å lagre bildet |
| 87 | http.Error | ja | nynorsk | ikkje, arrangementsbiletet | Klarte ikkje å oppdatere revisjonsspor for arrangementsbiletet | Klarte ikke å oppdatere revisjonsspor for arrangementsbildet |
| 119 | http.Error | nei | nynorsk | manglar | Arrangement-ID manglar | Arrangement-ID mangler |
| 125 | ui |  | nynorsk | bilettype | Ugyldig bilettype | Ugyldig bildetype |
| 132 | http.Error | nei | nynorsk | ikkje, bilete, frå | Klarte ikkje å hente bilete frå skjemaet | Klarte ikke å hente bildet fra skjemaet |
| 141 | http.Error | nei | nynorsk | ikkje, biletet | Klarte ikkje å lagre biletet | Klarte ikke å lagre bildet |
| 148 | http.Error | nei | nynorsk | ikkje, biletet | Klarte ikkje å lagre biletet | Klarte ikke å lagre bildet |
| 153 | http.Error | nei | nynorsk | ikkje, arrangementsbiletet | Klarte ikkje å oppdatere revisjonsspor for arrangementsbiletet | Klarte ikke å oppdatere revisjonsspor for arrangementsbildet |
| 237 | ui | ja | nynorsk | ikkje | Klarte ikkje å hente arrangement: | Klarte ikke å hente arrangement: |
| 242 | ui | ja | nynorsk | Fann, ikkje | Fann ikkje arrangementet | Fant ikke arrangementet |
| 246 | ui | ja | borderline | Lagra | Lagra | Lagret |
| 428 | ui | ja | nynorsk | Heim | Heim | Hjem |
| 431 | ui | ja | nynorsk | bilete | Legg til bilete | Legg til bilde |
| 459 | ui | ja | nynorsk | bilete | Legg til bilete | Legg til bilde |
| 467 | ui | ja | nynorsk | Ver venleg, ikkje, bilete | Ver venleg og ikkje bruk AI/KI-genererte bilete | Vær vennlig og ikke bruk AI/KI-genererte bilder |
| 470 | ui | ja | nynorsk | Sidan, bilete, eit, ønskjer, ikkje | Sidan AI/KI-genererte bilete er eit betent tema, ønskjer vi i RegnCon ikkje | Siden AI/KI-genererte bilder er et betent tema, ønsker vi i RegnCon ikke |
| 471 | ui | ja | nynorsk | arrangørar, bilete, hovudbilete | at arrangørar skal bruke slike bilete som hovudbilete for arrangementet sitt. | at arrangører skal bruke slike bilder som hovedbilde for arrangementet sitt. |
| 474 | ui | ja | nynorsk | bilete, ope | Om du er i tvil om bruk av bilete, anbefaler vi å finne bilete som er ope | Om du er i tvil om bruk av bilder, anbefaler vi å finne bilder som er åpent |
| 475 | ui | ja | nynorsk | tilgjengelege, bilete, elles, brukte | tilgjengelege for bruk, eller bilete som elles blir brukte til å demonstrere | tilgjengelige for bruk, eller bilder som ellers blir brukt til å demonstrere |
| 476 | ui | ja | nynorsk | spelet | og promotere spelet/systemet du skal bruke på arrangementet. | og promotere spillet/systemet du skal bruke på arrangementet. |
| 483 | ui | ja | nynorsk | bilete | Last opp bilete | Last opp bilde |
| 517 | ui | ja | nynorsk | Lastar, bilete, Ver venleg | Lastar bilete. Ver venleg og vent... | Laster bilde. Vær vennlig og vent... |
| 518 | ui | ja | nynorsk | bilete | Klikk eller dra og slipp for å bytte bilete | Klikk eller dra og slipp for å bytte bilde |
| 529 | ui | ja | nynorsk | Lastar, bilete, Ver venleg | $_uploadProgress === 100 ? 'Lastar bilete. Ver venleg og vent...' : '' | $_uploadProgress === 100 ? 'Laster bilde. Vær vennlig og vent...' : '' |
| 549 | ui | ja | nynorsk | Lastar, bilete | Lastar bilete | Laster bilde |
| 562 | ui | ja | nynorsk | Lastar, bilete | Lastar bilete | Laster bilde |

### [components/formsubmission/event_new.templ](../components/formsubmission/event_new.templ) (23)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 86 | http.Error | kanskje | nynorsk | ikkje | Klarte ikkje å lese skjemadata | Klarte ikke å lese skjemadata |
| 104 | http.Error | kanskje | nynorsk | ikkje | Klarte ikkje å lese skjemadata | Klarte ikke å lese skjemadata |
| 266 | ui | kanskje | nynorsk | vere, teikn | Tittel må vere minst 3 teikn lang | Tittel må være minst 3 tegn lang |
| 266 | ui | kanskje | nynorsk | spelmodul | Tittel på spelmodul / arrangement | Tittel på spillmodul / arrangement |
| 267 | ui | kanskje | nynorsk | vere, teikn | Introduksjon må vere minst 3 teikn lang | Introduksjon må være minst 3 tegn lang |
| 268 | ui | kanskje | nynorsk | Skildring, vere, teikn | Skildring må vere minst 3 teikn lang | Beskrivelse må være minst 3 tegn lang |
| 268 | ui | kanskje | nynorsk | Skildring | Skildring | Beskrivelse |
| 269 | ui | kanskje | nynorsk | vere, teikn | System må vere minst 3 teikn langt | System må være minst 3 tegn langt |
| 270 | ui | kanskje | nynorsk | Namn, vere, teikn | Namn på arrangør må vere minst 3 teikn langt | Navn på arrangør må være minst 3 tegn langt |
| 270 | ui | kanskje | nynorsk | Namn | Namn på arrangør | Navn på arrangør |
| 271 | ui | kanskje | nynorsk | E-postadressa, vere | E-postadressa må vere gyldig | E-postadressen må være gyldig |
| 272 | ui | kanskje | nynorsk | vere | Telefonnummeret må vere minst 8 siffer | Telefonnummeret må være minst 8 siffer |
| 273 | ui | kanskje | nynorsk | Vel, ei | Vel ei gyldig aldersgruppe | Velg en gyldig aldersgruppe |
| 274 | ui | kanskje | nynorsk | Vel, ei, varigheit | Vel ei gyldig varigheit | Velg en gyldig varighet |
| 274 | ui | kanskje | nynorsk | Varigheit | Varigheit | Varighet |
| 275 | ui | kanskje | nynorsk | tal, spelarar, vere | Maks tal på spelarar må vere mellom 1 og 18 | Maks antall spillere må være mellom 1 og 18 |
| 275 | ui | kanskje | nynorsk | tal, spelarar | Maks tal på spelarar | Maks antall spillere |
| 276 | ui | kanskje | nynorsk | Nybegynnarvennleg | Nybegynnarvennleg | Nybegynnervennlig |
| 277 | ui | kanskje | nynorsk | haldast | Kan haldast på engelsk | Kan holdes på engelsk |
| 295 | ui | ja | borderline | påmeldinga | Takk for påmeldinga! | Takk for påmeldingen! |
| 301 | ui | ja | nynorsk | ikkje | Klarte ikkje å lagre arrangementet. Prøv igjen. Feil: | Klarte ikke å lagre arrangementet. Prøv igjen. Feil: |
| 581 | ui | kanskje | nynorsk | ikkje | Klarte ikkje å oppdatere arrangementet: %v | Klarte ikke å oppdatere arrangementet: %v |
| 595 | http.Error | kanskje | nynorsk | Fann, ikkje, endringar, vart gjorde | Fann ikkje arrangementet, eller ingen endringar vart gjorde | Fant ikke arrangementet, eller ingen endringer ble gjort |

### [components/formsubmission/feedback_message.go](../components/formsubmission/feedback_message.go) (2)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 4 | ui | ja | nynorsk | ikkje, endringa, held fram | Klarte ikkje å lagre endringa. Prøv igjen. Kontakt styret dersom problemet held fram. | Klarte ikke å lagre endringen. Prøv igjen. Kontakt styret dersom problemet vedvarer. |
| 5 | ui | ja | nynorsk | ikkje, endringa, held fram | Klarte ikkje å lagre endringa. Prøv igjen. Sjekk logger dersom problemet held fram. | Klarte ikke å lagre endringen. Prøv igjen. Sjekk logger dersom problemet vedvarer. |

### [components/formsubmission/other_details.templ](../components/formsubmission/other_details.templ) (49)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 25 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å lese skjemadata | Klarte ikke å lese skjemadata |
| 31 | http.Error | nei | nynorsk | manglar | Arrangement-ID manglar | Arrangement-ID mangler |
| 36 | http.Error | nei | nynorsk | ikkje, aldersgruppa | Klarte ikkje å oppdatere aldersgruppa for arrangementet i databasen | Klarte ikke å oppdatere aldersgruppen for arrangementet i databasen |
| 41 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å sende oppdatering | Klarte ikke å sende oppdatering |
| 55 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å lese skjemadata | Klarte ikke å lese skjemadata |
| 61 | http.Error | nei | nynorsk | manglar | Arrangement-ID manglar | Arrangement-ID mangler |
| 66 | http.Error | nei | nynorsk | ikkje, varigheita | Klarte ikkje å oppdatere varigheita for arrangementet i databasen | Klarte ikke å oppdatere varigheten for arrangementet i databasen |
| 71 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å sende oppdatering | Klarte ikke å sende oppdatering |
| 85 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å lese skjemadata | Klarte ikke å lese skjemadata |
| 91 | http.Error | nei | nynorsk | manglar | Arrangement-ID manglar | Arrangement-ID mangler |
| 96 | http.Error | nei | nynorsk | ikkje, nybegynnarvennleg | Klarte ikkje å oppdatere nybegynnarvennleg-statusen for arrangementet i databasen | Klarte ikke å oppdatere nybegynnervennlig-statusen for arrangementet i databasen |
| 101 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å sende oppdatering | Klarte ikke å sende oppdatering |
| 115 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å lese skjemadata | Klarte ikke å lese skjemadata |
| 121 | http.Error | nei | nynorsk | manglar | Arrangement-ID manglar | Arrangement-ID mangler |
| 126 | http.Error | nei | nynorsk | ikkje, haldast | Klarte ikkje å oppdatere om arrangementet kan haldast på engelsk | Klarte ikke å oppdatere om arrangementet kan holdes på engelsk |
| 131 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å sende oppdatering | Klarte ikke å sende oppdatering |
| 145 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å lese skjemadata | Klarte ikke å lese skjemadata |
| 151 | http.Error | nei | nynorsk | manglar | Arrangement-ID manglar | Arrangement-ID mangler |
| 156 | http.Error | nei | nynorsk | ikkje, tal, spelarar | Klarte ikkje å oppdatere maks tal på spelarar for arrangementet i databasen | Klarte ikke å oppdatere maks antall spillere for arrangementet i databasen |
| 161 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å sende oppdatering | Klarte ikke å sende oppdatering |
| 175 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å lese skjemadata | Klarte ikke å lese skjemadata |
| 181 | http.Error | nei | nynorsk | manglar | Arrangement-ID manglar | Arrangement-ID mangler |
| 188 | http.Error | nei | nynorsk | ikkje, merknadane | Klarte ikkje å oppdatere merknadane for arrangementet i databasen | Klarte ikke å oppdatere merknadene for arrangementet i databasen |
| 193 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å sende oppdatering | Klarte ikke å sende oppdatering |
| 202 | ui | ja | nynorsk | Viss, skremmande, krev, meir, deltakarar, velje, setje, ei | Viss arrangementet ditt er skremmande eller krev litt meir modne deltakarar, kan du velje dette for å setje ei attenårsgrense. | Hvis arrangementet ditt er skremmende eller krever litt mer modne deltakere, kan du velge dette for å sette en attenårsgrense. |
| 204 | ui | ja | nynorsk | tel, dei, passar, sette | På Regncon tel dei som er tolv år eller yngre som barn. Arrangement som passar for barn, blir sette opp i søndagspuljen. | På Regncon regnes de som er tolv år eller yngre som barn. Arrangementer som passer for barn, blir satt opp i søndagspuljen. |
| 206 | ui | ja | nynorsk | Ho, omfattar, dei | Standard aldersgruppe. Ho omfattar dei som kan delta i alle puljene, altså alle som er tretten år eller eldre. | Standard aldersgruppe. Den omfatter de som kan delta i alle puljene, altså alle som er tretten år eller eldre. |
| 224 | ui | ja | nynorsk | Vel, viss, truleg, timar, Då, veit, deltakarane, dei, noko, gjere, medan, ventar | Vel denne viss arrangementet truleg varer mindre enn tre timar. Då veit deltakarane at dei må finne noko å gjere medan dei ventar på neste pulje. | Velg denne hvis arrangementet trolig varer mindre enn tre timer. Da vet deltakerne at de må finne noe å gjøre mens de venter på neste pulje. |
| 226 | ui | ja | nynorsk | Vel, viss, truleg, timar, meir, Då, veit, deltakarane, påverke | Vel denne viss arrangementet truleg varer seks timar eller meir. Då veit deltakarane at det kan påverke nattesøvnen eller den neste puljen. | Velg denne hvis arrangementet trolig varer seks timer eller mer. Da vet deltakerne at det kan påvirke nattesøvnen eller den neste puljen. |
| 228 | ui | ja | nynorsk | Ein, vanleg, timar, varigheita, dei, arrangementa, haldne | Ein vanleg pulje varer 4-5 timar på Regncon. Erfaringsmessig er dette varigheita til dei aller fleste arrangementa som blir haldne på connet. | En vanlig pulje varer 4-5 timer på Regncon. Erfaringsmessig er dette varigheten til de aller fleste arrangementene som blir holdt på connet. |
| 245 | ui | ja | nynorsk | detaljar | Andre detaljar | Andre detaljer |
| 266 | ui | ja | nynorsk | Kva, tilrådd | Kva alder er arrangementet tilrådd for? | Hvilken alder er arrangementet anbefalt for? |
| 270 | ui | ja | nynorsk | Varigheit | Varigheit | Varighet |
| 288 | ui | ja | nynorsk | Kor | Kor lenge varer puljen? | Hvor lenge varer puljen? |
| 294 | ui | ja | nynorsk | Nybegynnarvennleg | Nybegynnarvennleg arrangement | Nybegynnervennlig arrangement |
| 296 | ui | ja | nynorsk | Valfritt | Valfritt | Valgfritt |
| 311 | ui | ja | nynorsk | nybegynnarvennleg | Arrangementet er nybegynnarvennleg | Arrangementet er nybegynnervennlig |
| 316 | ui | ja | nynorsk | nybegynnarvennleg | Arrangementet er nybegynnarvennleg | Arrangementet er nybegynnervennlig |
| 318 | ui | ja | nynorsk | dei, viss, spelet | Erfaringsmessig stemmer dette for dei fleste arrangement – men viss spelet er komplisert, | Erfaringsmessig stemmer dette for de fleste arrangementer – men hvis spillet er komplisert, |
| 319 | ui | ja | nynorsk | ikkje, reglar | og du ikkje vil bruke tid på å forklare reglar, kan du la boksen stå tom. | og du ikke vil bruke tid på å forklare regler, kan du la boksen stå tom. |
| 324 | ui | ja | nynorsk | haldast | Kan haldast på engelsk | Kan holdes på engelsk |
| 325 | ui | ja | nynorsk | Valfritt | Valfritt | Valgfritt |
| 340 | ui | ja | nynorsk | haldast | Arrangementet kan haldast på engelsk | Arrangementet kan holdes på engelsk |
| 345 | ui | ja | nynorsk | haldast | Arrangementet kan haldast på engelsk | Arrangementet kan holdes på engelsk |
| 347 | ui | ja | nynorsk | Vel, viss, halde | Vel denne viss du kan halde arrangementet på engelsk | Velg denne hvis du kan holde arrangementet på engelsk |
| 348 | ui | ja | nynorsk | finn, deltakarane, startar | (så finn du eventuelt ut av det med deltakarane når arrangementet startar) | (så finner du eventuelt ut av det med deltakerne når arrangementet starter) |
| 352 | ui | ja | nynorsk | tal, spelarar, utanom | Maks tal på spelarar (utanom arrangør) | Maks antall spillere (utenom arrangør) |
| 368 | ui | ja | nynorsk | Kor, vere, spele | Kor mange kan vere med og spele? | Hvor mange kan være med og spille? |
| 377 | ui | ja | nynorsk | til dømes, nokre, ikkje, passar, halde | Andre merknader – er det til dømes nokre tidspunkt det ikkje passar å halde arrangementet? | Andre merknader – er det for eksempel noen tidspunkt det ikke passer å holde arrangementet? |

### [components/formsubmission/puljefordeling.templ](../components/formsubmission/puljefordeling.templ) (18)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 24 | http.Error | nei | nynorsk | manglar | Arrangement-ID manglar | Arrangement-ID mangler |
| 33 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å lese skjemadata | Klarte ikke å lese skjemadata |
| 39 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å oppdatere puljefordeling | Klarte ikke å oppdatere puljefordeling |
| 45 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å sende oppdatering | Klarte ikke å sende oppdatering |
| 60 | http.Error | nei | nynorsk | manglar | Arrangement-ID manglar | Arrangement-ID mangler |
| 67 | http.Error | nei | nynorsk | manglar | Pulje-ID manglar | Pulje-ID mangler |
| 81 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å lese skjemadata | Klarte ikke å lese skjemadata |
| 101 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å oppdatere arrangementet i puljen | Klarte ikke å oppdatere arrangementet i puljen |
| 107 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å sende oppdatering | Klarte ikke å sende oppdatering |
| 122 | http.Error | nei | nynorsk | manglar | Arrangement-ID manglar | Arrangement-ID mangler |
| 129 | http.Error | nei | nynorsk | manglar | Pulje-ID manglar | Pulje-ID mangler |
| 140 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å lese skjemadata | Klarte ikke å lese skjemadata |
| 147 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å oppdatere rommet for puljen | Klarte ikke å oppdatere rommet for puljen |
| 153 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å sende oppdatering | Klarte ikke å sende oppdatering |
| 235 | ui | ja | nynorsk | ikkje | Klarte ikkje å hente arrangementspulje: | Klarte ikke å hente arrangementspulje: |
| 240 | ui | ja | nynorsk | ikkje | Klarte ikkje å hente rom | Klarte ikke å hente rom |
| 287 | ui | ja | nynorsk | Ikkje | Ikkje tildelt | Ikke tildelt |
| 326 | ui | ja | nynorsk | ikkje | Klarte ikkje å hente puljer: | Klarte ikke å hente puljer: |

### [components/formsubmission/set_event_in_pulje_test.go](../components/formsubmission/set_event_in_pulje_test.go) (3)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 21 | test | – | nynorsk | Spel (test fixture title) | INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players) VALUES ('e1', 'Spel', '', '', 'Ola', 'ola@x.no', '', 4) | INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players) VALUES ('e1', 'Spill', '', '', 'Ola', 'ola@x.no', '', 4) |
| 65 | test | – | nynorsk | Spel (test fixture title) | INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players) VALUES ('e1', 'Spel', '', '', 'Ola', 'ola@x.no', '', 4) | INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players) VALUES ('e1', 'Spill', '', '', 'Ola', 'ola@x.no', '', 4) |
| 112 | test | – | nynorsk | Spel (test fixture title) | INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players) VALUES ('e1', 'Spel', '', '', 'Ola', 'ola@x.no', '', 4) | INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players) VALUES ('e1', 'Spill', '', '', 'Ola', 'ola@x.no', '', 4) |

### [components/formsubmission/statusCard.templ](../components/formsubmission/statusCard.templ) (5)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 26 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å lese skjemadata | Klarte ikke å lese skjemadata |
| 32 | http.Error | nei | nynorsk | manglar | Arrangement-ID manglar | Arrangement-ID mangler |
| 37 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å oppdatere statusen for arrangementet i databasen | Klarte ikke å oppdatere statusen for arrangementet i databasen |
| 42 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å sende oppdatering | Klarte ikke å sende oppdatering |
| 72 | ui | ja | nynorsk | nedst | (bruk knappen nedst i skjemaet) | (bruk knappen nederst i skjemaet) |

### [components/formsubmission/submit_section.templ](../components/formsubmission/submit_section.templ) (4)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 18 | http.Error | nei | nynorsk | manglar | Arrangement-ID manglar | Arrangement-ID mangler |
| 24 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å oppdatere statusen for arrangementet i databasen | Klarte ikke å oppdatere statusen for arrangementet i databasen |
| 29 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å sende oppdatering | Klarte ikke å sende oppdatering |
| 46 | ui | ja | nynorsk | lagra, treng, berre | Kladden blir lagra automatisk. Når du er klar til å sende arrangementet inn, treng du berre å klikke på knappen her! | Kladden blir lagret automatisk. Når du er klar til å sende arrangementet inn, trenger du bare å klikke på knappen her! |

### [components/header/menu.templ](../components/header/menu.templ) (3)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 154 | ui | ja | borderline | framside | Regncon framside | Regncon forside |
| 626 | ui | ja | nynorsk | Vanlege | Vanlege Spørsmål | Vanlige Spørsmål |
| 881 | ui | ja | nynorsk | Vanlege | Vanlege Spørsmål | Vanlige Spørsmål |

### [components/header/menu_test.go](../components/header/menu_test.go) (4)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 36 | test | – | nynorsk | vanlege | Så skal brukeren bare få navigasjonslenker til forsiden, egen profil, utlogging og vanlege spørsmål. | Så skal brukeren bare få navigasjonslenker til forsiden, egen profil, utlogging og vanlige spørsmål. |
| 55 | test | – | borderline | vanlege (external regncon.no URL slug; only change if the site URL changes) | a[href="https://www.regncon.no/vanlege-sporsmal/"] .inline-icon | a[href="https://www.regncon.no/vanlege-sporsmal/"] .inline-icon |
| 68 | test | – | nynorsk | vanlege | Så skal brukeren få navigasjonslenker til forsiden, egen profil, utlogging, adminområdene og vanlege spørsmål. | Så skal brukeren få navigasjonslenker til forsiden, egen profil, utlogging, adminområdene og vanlige spørsmål. |
| 90 | test | – | borderline | vanlege (external regncon.no URL slug; only change if the site URL changes) | a[href="https://www.regncon.no/vanlege-sporsmal/"] .inline-icon | a[href="https://www.regncon.no/vanlege-sporsmal/"] .inline-icon |

### [components/profile/my_program.templ](../components/profile/my_program.templ) (3)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 147 | ui | ja | nynorsk | ikkje, enno | Programmet for Regncon er ikkje publisert enno | Programmet for Regncon er ikke publisert ennå |
| 149 | ui | ja | nynorsk | framleis, planlegginga, leggje | Vi arbeider framleis med planlegginga og vil leggje ut programmet så snart det er klart. | Vi arbeider fortsatt med planleggingen og vil legge ut programmet så snart det er klart. |
| 150 | ui | ja | nynorsk | nettsida, oppdateringar, nyheiter | Følg med på nettsida for oppdateringar og nyheiter. | Følg med på nettsiden for oppdateringer og nyheter. |

### [components/profile/my_program_render_test.go](../components/profile/my_program_render_test.go) (1)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 125 | test | – | nynorsk | ikkje, enno | Programmet for Regncon er ikkje publisert enno | Programmet for Regncon er ikke publisert ennå |

### [components/profile/my_program_test.go](../components/profile/my_program_test.go) (6)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 224 | test | – | nynorsk | ei, open | Gitt en manuell spillerplassering i ei open pulje. | Gitt en manuell spillerplassering i en åpen pulje. |
| 225 | test | – | borderline | lasta | Når festivalprogrammet blir lasta. | Når festivalprogrammet blir lastet. |
| 226 | test | – | nynorsk | ein gong, gjer | Så viser programmet arrangementet med ein gong, slik varselet i interessedialogen gjer. | Så viser programmet arrangementet med en gang, slik varselet i interessedialogen gjør. |
| 258 | test | – | nynorsk | ei | Gitt en manuell spillerplassering i ei låst pulje. | Gitt en manuell spillerplassering i en låst pulje. |
| 259 | test | – | borderline | lasta | Når festivalprogrammet blir lasta. | Når festivalprogrammet blir lastet. |
| 260 | test | – | nynorsk | sjølv, puljefordelinga, ikkje | Så viser programmet arrangementet sjølv om puljefordelinga ikkje er ferdig. | Så viser programmet arrangementet selv om puljefordelingen ikke er ferdig. |

### [components/profile/my_tickets.templ](../components/profile/my_tickets.templ) (1)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 29 | ui | ja | nynorsk | Billettar | Billettar | Billetter |

### [models/billettholder.go](../models/billettholder.go) (1)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 87 | ui | nei | nynorsk | spelar | spelar | spiller |

### [models/event-model.go](../models/event-model.go) (14)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 20 | ui | ja | nynorsk | Godkjend | Godkjend | Godkjent |
| 22 | ui | ja | borderline | Forkasta | Forkasta | Forkastet |
| 40 | ui | ja | nynorsk | Rollespel | Rollespel | Rollespill |
| 42 | ui | ja | nynorsk | Brettspel | Brettspel | Brettspill |
| 44 | ui | ja | nynorsk | Kortspel | Kortspel | Kortspill |
| 46 | ui | ja | nynorsk | Anna (annet) | Anna | Annet |
| 67 | ui | ja | nynorsk | Vaksne | Vaksne (18+) | Voksne (18+) |
| 76 | ui | ja | nynorsk | Barnevennleg | Barnevennleg | Barnevennlig |
| 78 | ui | ja | nynorsk | Eigna, vaksne | Eigna for vaksne | Egnet for voksne |
| 104 | ui | ja | nynorsk | Vanleg | Vanleg pulje | Vanlig pulje |
| 106 | ui | ja | nynorsk | Kortare, timar | Kortare (2-3 timar) | Kortere (2-3 timer) |
| 108 | ui | ja | nynorsk | timar | Lengre (6+ timar) | Lengre (6+ timer) |
| 117 | ui | ja | nynorsk | timar | Varer under 3 timar | Varer under 3 timer |
| 119 | ui | ja | nynorsk | timar | Varer over 6 timar | Varer over 6 timer |

### [models/interests-model.go](../models/interests-model.go) (1)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 25 | ui | ja | nynorsk | Ikkje | Ikkje interessert | Ikke interessert |

### [pages/admin/admin.go](../pages/admin/admin.go) (1)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 129 | ui | ja | nynorsk | ikkje | Finner ikkje rom | Finner ikke rom |

### [pages/admin/admin_page.templ](../pages/admin/admin_page.templ) (1)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 74 | ui | ja | nynorsk | spel, puljane, sjå, kor, nybegynnarvennlege, kvar, Åtvarar, ein, spelleiar, same | Dra spel inn i puljane og sjå kor mange spel, 18+ og nybegynnarvennlege kvar pulje har. Åtvarar om ein spelleiar har to spel i same pulje. | Dra spill inn i puljene og se hvor mange spill, 18+ og nybegynnervennlige hver pulje har. Advarer om en spilleder har to spill i samme pulje. |

### [pages/admin/approval/editForm/edit_form_index_test.go](../pages/admin/approval/editForm/edit_form_index_test.go) (3)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 15 | test | – | nynorsk | redigeringssida, eit | Gitt redigeringssida for eit arrangement. | Gitt redigeringssiden for et arrangement. |
| 16 | test | – | borderline | sida, rendra | Når sida blir rendra. | Når siden blir rendret. |
| 17 | test | – | nynorsk | inneheld, ho, ikkje, spelartildeling | Så inneheld ho ikkje lenger spelartildeling. | Så inneholder den ikke lenger spillertildeling. |

### [pages/admin/puljefordeling_tab_assignment_test.go](../pages/admin/puljefordeling_tab_assignment_test.go) (13)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 92 | test | – | nynorsk | ei, berre, éin, tilgjengeleg, tilkopling | Gitt ei SQLite-database med berre éin tilgjengeleg tilkopling. | Gitt en SQLite-database med bare én tilgjengelig tilkobling. |
| 93 | test | – | nynorsk | lastar, spelarleiarar | Når dialogen lastar interesser og spelarleiarar. | Når dialogen laster interesser og spilledere. |
| 94 | test | – | nynorsk | kvar, spørjing, tilkoplinga | Så blir kvar spørjing ferdig før den neste bruker tilkoplinga. | Så blir hver spørring ferdig før den neste bruker tilkoblingen. |
| 175 | test | – | nynorsk | interessa | Når en administrator endrer interessa i puljefordeling. | Når en administrator endrer interessen i puljefordeling. |
| 176 | test | – | nynorsk | interessa | Så blir interessa oppdatert. | Så blir interessen oppdatert. |
| 201 | test | – | nynorsk | eitt | Gitt en åpen pulje med interesse for eitt arrangement. | Gitt en åpen pulje med interesse for ett arrangement. |
| 202 | test | – | nynorsk | vel, Ikkje | Når administratoren vel Ikkje interessert. | Når administratoren velger Ikke interessert. |
| 203 | test | – | nynorsk | berre, interessa, fjerna | Så blir berre interessa for det arrangementet fjerna. | Så blir bare interessen for det arrangementet fjernet. |
| 224 | test | – | nynorsk | ei | Gitt en publisert pulje med ei interesse. | Gitt en publisert pulje med en interesse. |
| 225 | test | – | nynorsk | interessa | Når administratoren prøver å fjerne interessa. | Når administratoren prøver å fjerne interessen. |
| 226 | test | – | nynorsk | interessa, endringa | Så blir endringa avvist og interessa står urørt. | Så blir endringen avvist og interessen står urørt. |
| 274 | test | – | nynorsk | Ikkje | Ikkje interessert | Ikke interessert |
| 503 | test | – | nynorsk | Voksenspel | VALUES ('evA','Voksenspel',…) | VALUES ('evA','Voksenspill',…) |

### [pages/admin/puljefordeling_tab_test.go](../pages/admin/puljefordeling_tab_test.go) (7)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 332 | test | – | nynorsk | ein, deltakar, eit | Gitt ein deltakar under 18 som er manuelt plassert i eit 18+-arrangement. | Gitt en deltaker under 18 som er manuelt plassert i et 18+-arrangement. |
| 318 | test | – | nynorsk | Voksenspel | VALUES ('evA','Voksenspel',…) | VALUES ('evA','Voksenspill',…) |
| 334 | test | – | nynorsk | flisa, merkast | Så skal flisa merkast med «Under 18». | Så skal flisen merkes med «Under 18». |
| 379 | test | – | nynorsk | ein, spelleiar, eit | Gitt ein spelleiar under 18 på eit 18+-arrangement. | Gitt en spilleder under 18 på et 18+-arrangement. |
| 381 | test | – | nynorsk | spelleiar, lina, merkast | Så skal spelleiar-lina merkast med «Under 18». | Så skal spilleder-linjen merkes med «Under 18». |
| 387 | test | – | nynorsk | Voksenspel | VALUES ('evA','Voksenspel',…) | VALUES ('evA','Voksenspill',…) |
| 408 | test | – | nynorsk | Voksenspel | VALUES ('evA','Voksenspel',…) | VALUES ('evA','Voksenspill',…) |

### [pages/admin/puljeoppsett.templ](../pages/admin/puljeoppsett.templ) (8)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 43 | ui | ja | nynorsk | ikkje, byggje | Klarte ikkje å byggje brettet: | Klarte ikke å bygge brettet: |
| 47 | ui | ja | nynorsk | spelleiar, kollisjon(ar) | %d spelleiar-kollisjon(ar) | %d spilleder-kollisjon(er) |
| 64 | ui | ja | nynorsk | spel | %d spel | %d spill |
| 84 | ui | ja | nynorsk | spel | %d spel · %d 18+ · %d nyb. | %d spill · %d 18+ · %d nyb. |
| 90 | ui | ja | nynorsk | Spelleiar | ⚠ Spelleiar-kollisjon — | ⚠ Spilleder-kollisjon — |
| 90 | ui | ja | nynorsk | spel | spel i denne puljen | spill i denne puljen |
| 219 | http.Error | nei | nynorsk | manglar | Arrangement-ID manglar | Arrangement-ID mangler |
| 229 | http.Error | nei | nynorsk | ikkje | Klarte ikkje å oppdatere puljen | Klarte ikke å oppdatere puljen |

### [pages/admin/puljeoppsett_render_test.go](../pages/admin/puljeoppsett_render_test.go) (4)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 18 | test | – | nynorsk | spel | Ledig spel | Ledig spill |
| 23 | test | – | nynorsk | spel | Ledig spel | Ledig spill |
| 23 | test | – | nynorsk | spel | 1 spel | 1 spill |
| 89 | test | – | nynorsk | Spelleiar | Spelleiar-kollisjon | Spilleder-kollisjon |

### [pages/admin/puljeoppsett_route_test.go](../pages/admin/puljeoppsett_route_test.go) (2)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 25 | test | – | nynorsk | Spel | Spel | Spill |
| 63 | test | – | nynorsk | Spel | Spel | Spill |

### [pages/admin/rooms/rooms_assignment_page.templ](../pages/admin/rooms/rooms_assignment_page.templ) (1)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 235 | ui | ja | nynorsk | ikkje | Klarte ikkje å tildele rommet. Prøv igjen. | Klarte ikke å tildele rommet. Prøv igjen. |

### [pages/event/event.go](../pages/event/event.go) (7)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 42 | ui | ja | nynorsk | ikkje, interessa, billettheldaren | Du har ikkje tilgang til å endre interessa til denne billettheldaren. Kontakt styret. | Du har ikke tilgang til å endre interessen til denne billettholderen. Kontakt styret. |
| 45 | ui | ja | nynorsk | pulja, ikkje, tilgjengeleg | Denne pulja er ikkje tilgjengeleg for dette arrangementet. | Denne puljen er ikke tilgjengelig for dette arrangementet. |
| 48 | ui | ja | nynorsk | Pulja, ikkje, medan, spelarar | Pulja er låst. Du kan ikkje melde eller endre interesse lenger medan vi fordeler spelarar. | Puljen er låst. Du kan ikke melde eller endre interesse lenger mens vi fordeler spillere. |
| 51 | ui | ja | nynorsk | Puljefordelinga, sjå, kva, fekk | Puljefordelinga er klar. Gå til profilen din for å sjå kva du fekk. | Puljefordelingen er klar. Gå til profilen din for å se hva du fikk. |
| 56 | ui | ja | nynorsk | ein, då, interessa, lagrast, held fram | Det oppstod ein feil då interessa skulle lagrast. Prøv igjen, eller kontakt styret dersom feilen held fram. | Det oppstod en feil da interessen skulle lagres. Prøv igjen, eller kontakt styret dersom feilen fortsetter. |
| 249 | ui | ja | nynorsk | Vel, billetthelder | Vel billetthelder før du melder interesse. | Velg billettholder før du melder interesse. |
| 256 | ui | ja | nynorsk | Vel | Vel pulje før du melder interesse. | Velg pulje før du melder interesse. |

### [pages/event/event_interest_test.go](../pages/event/event_interest_test.go) (6)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 58 | test | – | borderline | vanlege (external regncon.no URL slug; only change if the site URL changes) | a[href="https://www.regncon.no/vanlege-sporsmal/"] .inline-icon | a[href="https://www.regncon.no/vanlege-sporsmal/"] .inline-icon |
| 522 | test | – | nynorsk | ein, eit | Gitt at ein billettholder under 18 år er tildelt eit arrangement i puljen. | Gitt at en billettholder under 18 år er tildelt et arrangement i puljen. |
| 523 | test | – | nynorsk | eit, same, vert | Når eit 18-års arrangement i same pulje vert vist. | Når et 18-års arrangement i samme pulje blir vist. |
| 524 | test | – | nynorsk | berre, ikkje | Så skal berre aldersvarselet vises, ikkje tildelingsvarselet. | Så skal bare aldersvarselet vises, ikke tildelingsvarselet. |
| 592 | test | – | nynorsk | Vel, billetthelder | Vel billetthelder | Velg billettholder |
| 592 | test | – | nynorsk | Vel | Vel pulje | Velg pulje |

### [pages/event/event_page_admin_test.go](../pages/event/event_page_admin_test.go) (3)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 58 | test | – | nynorsk | pulja, ikkje, tilgjengeleg | Denne pulja er ikkje tilgjengeleg for dette arrangementet. | Denne puljen er ikke tilgjengelig for dette arrangementet. |
| 59 | test | – | nynorsk | Pulja, ikkje, medan, spelarar | Pulja er låst. Du kan ikkje melde eller endre interesse lenger medan vi fordeler spelarar. | Puljen er låst. Du kan ikke melde eller endre interesse lenger mens vi fordeler spillere. |
| 60 | test | – | nynorsk | Puljefordelinga, sjå, kva, fekk | Puljefordelinga er klar. Gå til profilen din for å sjå kva du fekk. | Puljefordelingen er klar. Gå til profilen din for å se hva du fikk. |

### [pages/login/login.go](../pages/login/login.go) (1)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 109 | ui | ja | nynorsk | Velkomen | Velkomen tilbake til Regncon 2026! | Velkommen tilbake til Regncon 2026! |

### [pages/login/login.templ](../pages/login/login.templ) (6)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 22 | ui | ja | nynorsk | Velkomen | Velkomen tilbake! | Velkommen tilbake! |
| 23 | ui | ja | nynorsk | allereie, innlogga, treng, ikkje | Det ser ut til at du allereie er innlogga og du treng ikkje logge inn på nytt. | Det ser ut til at du allerede er innlogget, og du trenger ikke logge inn på nytt. |
| 24 | ui | ja | nynorsk | vert, vidare, framsida | Du vert sendt vidare til framsida om | Du blir sendt videre til forsiden om |
| 24 | ui | ja | borderline | sekund (nynorsk plural) | sekund. | sekunder. |
| 25 | ui | ja | nynorsk | vidaresendinga, ikkje, denna, lenkja | Om vidaresendinga ikkje fungerer, kan du bruke denna lenkja til | Om videresendingen ikke fungerer, kan du bruke denne lenken til |
| 25 | ui | ja | borderline | framsida | framsida | forsiden |

### [service/checkIn/assign_users_test.go](../service/checkIn/assign_users_test.go) (12)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 80 | test | – | nynorsk | ein, ei, eksisterande, brukar, same, annan | Gitt at ein billettholder har fått lagt til ei manuell e-postadresse, og ein eksisterande brukar har same e-postadresse med annan casing. | Gitt at en billettholder har fått lagt til en manuell e-postadresse, og en eksisterende bruker har samme e-postadresse med annen casing. |
| 81 | test | – | nynorsk | e-postadressa, forsona, brukarar | Når e-postadressa blir forsona mot brukarar. | Når e-postadressen blir avstemt mot brukere. |
| 82 | test | – | nynorsk | ei, brukar-tilknyting | Så skal billettholderen få ei varig brukar-tilknyting. | Så skal billettholderen få en varig brukertilknytning. |
| 110 | test | – | nynorsk | ein, allereie, knytt, brukar, ei | Gitt at ein billettholder allereie er knytt til ein brukar via ei manuell e-postadresse. | Gitt at en billettholder allerede er knyttet til en bruker via en manuell e-postadresse. |
| 111 | test | – | nynorsk | same, e-postforsoning, køyrer | Når same e-postforsoning køyrer på nytt. | Når samme e-postavstemming kjører på nytt. |
| 112 | test | – | nynorsk | framleis, berre, finnast, ei, brukar-tilknyting | Så skal det framleis berre finnast ei brukar-tilknyting. | Så skal det fortsatt bare finnes én brukertilknytning. |
| 141 | test | – | nynorsk | ei, fjerna, frå, ein, attverande, brukaren | Gitt at ei manuell e-postadresse er fjerna frå ein billettholder, og ingen attverande e-postadresser på billettholderen samsvarer med brukaren. | Gitt at en manuell e-postadresse er fjernet fra en billettholder, og ingen gjenværende e-postadresser på billettholderen samsvarer med brukeren. |
| 142 | test | – | nynorsk | e-postadressa, forsona, brukar-tilknytingar | Når e-postadressa blir forsona mot brukar-tilknytingar. | Når e-postadressen blir avstemt mot brukertilknytninger. |
| 143 | test | – | nynorsk | brukar-tilknytinga, fjernast | Så skal den varige brukar-tilknytinga fjernast. | Så skal den varige brukertilknytningen fjernes. |
| 174 | test | – | nynorsk | ei, fjerna, frå, ein, anna, attverande, same, framleis, brukaren | Gitt at ei manuell e-postadresse er fjerna frå ein billettholder, men ei anna attverande e-postadresse på same billettholder framleis samsvarer med brukaren. | Gitt at en manuell e-postadresse er fjernet fra en billettholder, men en annen gjenværende e-postadresse på samme billettholder fortsatt samsvarer med brukeren. |
| 175 | test | – | nynorsk | e-postadressa, forsona, brukar-tilknytingar | Når e-postadressa blir forsona mot brukar-tilknytingar. | Når e-postadressen blir avstemt mot brukertilknytninger. |
| 176 | test | – | nynorsk | brukar-tilknytinga, behaldast | Så skal den varige brukar-tilknytinga behaldast. | Så skal den varige brukertilknytningen beholdes. |

### [service/puljefordeling/emulate_test.go](../service/puljefordeling/emulate_test.go) (2)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 251 | test | – | borderline | Vaksen (fixture first name; nynorsk for voksen) | Vaksen | Voksen |
| 270 | test | – | borderline | Vaksen (fixture name; must match 5664) | Vaksen Voksdal | Voksen Voksdal |

### [service/puljefordeling/tildelinger.go](../service/puljefordeling/tildelinger.go) (6)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 166 | error | nei | nynorsk | frå | fjern %s-tildeling for billettholder %d frå %s i %s: %w | fjern %s-tildeling for billettholder %d fra %s i %s: %w |
| 170 | error | nei | borderline | fjerna | les fjerna %s-tildeling: %w | les fjernet %s-tildeling: %w |
| 318 | error | nei | nynorsk | berre, ein, spelar | %w: dra-og-slipp kan berre flytte ein spelar | %w: dra-og-slipp kan bare flytte én spiller |
| 367 | ui | ja | nynorsk | ein, spelar | Bruk Legg til for å plassere ein GM som spelar. | Bruk Legg til for å plassere en GM som spiller. |
| 396 | ui | ja | nynorsk | spelar | Legg til som spelar på «%s» | Legg til som spiller på «%s» |
| 403 | ui | ja | nynorsk | spelar | Legg til som spelar på «%s» | Legg til som spiller på «%s» |

### [service/puljefordeling/tildelinger_test.go](../service/puljefordeling/tildelinger_test.go) (3)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 59 | test | – | nynorsk | spelar | Legg til som spelar på «Alpha» | Legg til som spiller på «Alpha» |
| 141 | test | – | nynorsk | ikkje | kan-ikkje-overstyre | kan-ikke-overstyre |
| 485 | test | – | nynorsk | Tilskodar | Tilskodar | Tilskuer |

### [service/rooms/rooms_status_test.go](../service/rooms/rooms_status_test.go) (1)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 24 | test | – | nynorsk | Laurdag | Laurdag | Lørdag |

### [service/rooms/rooms_test.go](../service/rooms/rooms_test.go) (1)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 24 | test | – | nynorsk | eit | Dette er eit gyldig rom | Dette er et gyldig rom |

### [service/rooms/rooms_update_test.go](../service/rooms/rooms_update_test.go) (1)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 25 | test | – | nynorsk | ei | Dette er ei oppdatert note | Dette er en oppdatert note |

### [service/rooms/rooms_validation.go](../service/rooms/rooms_validation.go) (4)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 16 | ui | ja | nynorsk | namn, ikkje, berre, innehalde | Rom namn kan ikkje berre innehalde mellomrom | Romnavn kan ikke bare inneholde mellomrom |
| 22 | ui | ja | nynorsk | Namn, ikkje, vere, teikn | Namn kan ikkje vere lengre enn 50 teikn | Navn kan ikke være lengre enn 50 tegn |
| 27 | ui | ja | nynorsk | påkravd | Romnummer er påkravd | Romnummer er påkrevd |
| 33 | ui | ja | nynorsk | ikkje, vere, teikn | Rom nummer kan ikkje vere lengre enn 10 teikn | Romnummer kan ikke være lengre enn 10 tegn |

### [static/js/error_feedback.js](../static/js/error_feedback.js) (1)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 6 | js | ja | nynorsk | ikkje, endringa | "Klarte ikkje å lagre endringa. Prøv igjen." | Klarte ikke å lagre endringen. Prøv igjen. |

### [static/web_components/banner_cropper.js](../static/web_components/banner_cropper.js) (1)

| Linje | Type | Vises? | Grad | Nynorsk | Tekst i dag | Forslag bokmål |
|---|---|---|---|---|---|---|
| 35 | js | nei | nynorsk | ikkje, endringa, held fram | 'Klarte ikkje å lagre endringa. Prøv igjen. Kontakt styret dersom problemet held fram.' | Klarte ikke å lagre endringen. Prøv igjen. Kontakt styret dersom problemet vedvarer. |

