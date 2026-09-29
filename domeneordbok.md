# Domeneordbok

## Ord
Interesse
Interesser
Interessenivå
Førstevalg
Påmelding
Billettholder
Billettholdere
Pulje
Puljer
Puljefordeling
Arrangement
Arrangementer
Programarrangement
Kladd
Innsendt
Godkjent
Annonsert
Forkastet
Publisert
Romfordeling
Aldersgrense
Aldersvarsel
Fredag Kveld
Lordag Morgen
Lordag Kveld
Sondag Morgen

## Billettholder
En billettholder er en CheckIn billett
Den er ikke av typen "Middag"

## Pulje
En pulje er et tidspunkt styret har valgt, og alle arrangementer i puljen skal spilles innenfor dette tidspunktet. For eksempel Fredag kveld: 18 - 23

## Interesse
Interesse er det en billettholder melder inn på et arrangement i en pulje. Den gjelder bare for det arrangementet i den valgte puljen. Samme arrangement kan gå i flere puljer, og da er interessen i hver pulje uavhengig av de andre.

Hva billettholderen kan gjøre, avhenger av arrangementet:

* **Interessenivå**: på arrangementer som er med i puljefordelingen (`events.is_in_puljefordeling = 1`, "Med i puljefordeling" i skjemaet) velger billettholderen hvor interessert hen er: `Veldig interessert`, `Middels interessert` eller `Litt interessert`. Interessenivået brukes når spillere blir fordelt på arrangementer i puljefordelingen, og `Veldig interessert` avgjør om billettholderen får [førstevalg](#førstevalg).
* **Programarrangement**: arrangementer som ikke er med i puljefordelingen (`events.is_in_puljefordeling = 0`) er åpne for alle. Billettholderen trenger ikke melde interesse eller reservere plass, og arrangementssiden viser ingen interessevelger. Se [Påmelding](#påmelding).

I koden heter dette `interests` (tabell) og `interest_level` (interessenivå), der én rad gjelder én billettholder, ett arrangement og én pulje.

## Førstevalg
En billettholder har fått førstevalg når hen er tildelt med rollen `Player` (spiller) på et arrangement i en pulje, og har gitt `Veldig interessert` på det samme arrangementet i den samme puljen. En `GM`-tildeling alene teller ikke. Er hen både `Player` og `GM` på et arrangement med `Veldig interessert`, teller det, fordi spillerrollen finnes.

I admin vises dette som "Har fått førstevalg" / "Har ikke fått førstevalg". Spillederstatus vises som et eget merke: "Spilleder (GM/DM)" / "Ikke spilleder".

Interesse gjelder én pulje, men førstevalg gjelder for alle puljene samlet. Målet er at alle billettholdere får minst ett førstevalg i løpet av festivalen. Har man fått førstevalg i én pulje, har man fått førstevalget sitt, også i de neste puljene.

I puljefordelingen går billettholdere som ikke har fått førstevalg ennå, foran på arrangementene de har gitt `Veldig interessert`. Jo flere puljer de har hatt `Veldig interessert` på et arrangement uten å få plass, jo høyere blir de prioritert. Å være GM på et arrangement gir ikke førstevalg.

I koden heter dette `first_choice` og `Forstevalg`, og i fordelingen (`solver`) `satisfied` og `top choice`.

## Påmelding
Påmelding er ikke en egen funksjon i systemet. Langvarige arrangementer som alle som vil kan være med på, for eksempel "Blood on the clock tower" eller "Cosplay", legges inn som [programarrangementer](#interesse) uten puljefordeling. Der trenger man ikke melde seg på eller reservere plass.

I tekster i løsningen brukes "påmelding" også om å melde interesse, for eksempel "Påmelding lukkes" i tidsplanen.

## Aldersgrense (18+)
Et arrangement med aldersgruppe `AdultsOnly` (`models.AgeGroupAdultsOnly`) har 18-årsgrense og vises som 18+. Om en billettholder er over 18, står i `billettholdere.is_over_18`, som settes fra fødselsdatoen i CheckIn.

Aldersgrensen håndteres manuelt av styret. Puljefordelingen (`solver`) ser ikke på alder, så en billettholder under 18 kan bli satt opp som spiller på et 18+-arrangement. Admin får i stedet beskjed:

* **Aldersvarsel**: puljefordelingssiden i admin viser "Under 18 på 18+-arrangement" under "Varsler om tildelinger", med én oppføring per arrangement og navnet på hver spiller under 18 (`FinnAldersvarsler` i `service/puljefordeling/aldersvarsler.go`). Varselet gjelder både plasser fra fordelingen og manuelle plasser. GM-er tas ikke med.
* Spillere og GM-er under 18 på et 18+-arrangement får merket "Under 18" på arrangementskortet.
* Før admin plasserer en billettholder under 18 manuelt på et 18+-arrangement, må hen bekrefte det i dialogen "Aldersgrense" ("Plasser likevel").

På arrangementssiden hindres ingen i å melde interesse. En billettholder under 18 ser notisen «Anbefalt for voksne (18+)» i interessemodalen for et 18+-arrangement, men kan velge interesse som vanlig. Notisen vises ikke når billettholderen allerede er tildelt et arrangement i puljen (`LoadInterestNoticeState` i `components/event_components/event_interests.templ`).

## Arrangementstatus
Et arrangement har alltid nøyaktig én av fem statuser (`events.status`, `models.EventStatus` i `models/event-model.go`):

| Status | Go-konstant | Betydning |
| --- | --- | --- |
| `Kladd` | `EventStatusDraft` | Arrangøren jobber fortsatt med arrangementet. |
| `Innsendt` | `EventStatusSubmitted` | Sendt inn til styret. |
| `Godkjent` | `EventStatusApproved` | Godkjent av styret, men ikke offentlig. Kan få rom i romfordelingen. |
| `Annonsert` | `EventStatusAnnounced` | Offentlig synlig. Den eneste statusen som vises på forsiden og i forrige/neste-navigasjonen. Kan også få rom. |
| `Forkastet` | `EventStatusArchived` | Avvist eller trukket. |

Det finnes ingen `Avvist`-status: et avvist arrangement er `Forkastet`, og Go-konstanten heter `EventStatusArchived`. `Annonsert` erstattet den gamle statusen `Publisert` (`EventStatusPublished`), som ikke finnes lenger. I statuskortet i arrangementsskjemaet er det bare admin som får valgene `Godkjent` og `Annonsert`.

Se [documentation/pulje-status-and-publishing.md](documentation/pulje-status-and-publishing.md) for pulje-status og publisering av programmet.

## Publisert – flere betydninger
"Publisert" betyr flere forskjellige ting i koden. Hold dem fra hverandre:

1. **Arrangementet er annonsert**: `events.status = 'Annonsert'` (`models.EventStatusAnnounced`). Arrangementet er offentlig synlig. Kall dette "annonsert", ikke "publisert".
2. **Programmet er publisert**: `program_publishing_state.is_published`, én rad med `id = 1`, lest av `program.IsPublished`. Gjelder hele programmet: forsiden med dager og puljer, interessevalg, Mitt festivalprogram og puljer, tider og rom på arrangementssiden. Styres fra "Publiser program" på `/admin`.
3. **Puljefordelingen er publisert**: `puljer.status = 'Completed'`, vist som "Puljefordeling publisert" i admin. Tildelingene i den puljen er publisert og kan ikke endres.
4. **Romfordelingen er publisert**: `puljer.rooms_published`, per pulje, styrt av bryteren "Publiser romfordeling" på `/admin/rooms/assignment/{pulje}` ("Synlig for alle" / "Ikke publisert"). Først da vises rommet, de offentlige romnotatene og romkartet for den puljen på arrangementssiden og i utskriftsvennlig program, også for admin. Arrangementssiden krever i tillegg at programmet er publisert. Flagget henger ikke sammen med `puljer.status`.
5. **`relation_event_puljer.is_published`**: en gammel kolonne per arrangement og pulje. Ingen produksjonskode leser den lenger, og `SetEventInPulje` lar den bevisst være urørt. Ikke bygg ny oppførsel på den.

Disse feltene handler ikke om publisering, men blandes lett sammen med den:

* `relation_event_puljer.is_in_pulje`: om arrangementet er satt opp i puljen. `v_event_puljer_active` filtrerer på `is_in_pulje = 1`, og `v_events_by_pulje_active` filtrerer i tillegg på status `Annonsert`.
* `events.is_in_puljefordeling`: om arrangementet er med i puljefordelingen, eller er et [programarrangement](#interesse).
* `puljer.status` (`Open`, `Locked`, `Completed`) styrer ikke hvilke arrangementer som vises offentlig eller i forrige/neste-navigasjonen.

## Kode
* _index betyr at filen skal sette opp NATS-integrasjon for domenet
* _page betyr at dette er siden som skal bruke entry #id som blei satt opp i _index til å vise html
