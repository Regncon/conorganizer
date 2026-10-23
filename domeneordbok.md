# Domeneordbok
Domene or følger DDD domain-driven design, og må være det samme i alle deler av systemet, frontend, backend, database, dokumentasjon  etc.

Ikke domeneord bør være på engelsk i kode, database og dokumentasjon, men kan være på norsk i tekster som vises til brukere, for eksempel i feilmeldinger.

Ikke domeneord kan godt være egne services eller tabeller etc. For eksempel ticket og ticket service.

## Domeneord
Interesse
Førstevalg
Pulje
Puljefordeling

## Ikke domeneord
Påmelding
Deltager
Interessenivå
Pujlearrangement
Programarrangement
Arrangement
Arrangementer
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

## Deprecated
Billettholder -> skal bli til deltager.
Billettholdere

## Om domeneord

### Pulje
En pulje er et tidspunkt styret har valgt, og alle arrangementer i puljen skal spilles innenfor dette tidspunktet. For eksempel Fredag kveld: 18 - 23

### Puljefordeling
er en handling styret gjør for og fordele deltakere på arrangementer for en pulje

### Interesse
Interesse er det en billettholder melder inn på et arrangement i en pulje. Den gjelder bare for det arrangementet i den valgte puljen. Samme arrangement kan gå i flere puljer, og da er interessen i hver pulje uavhengig av de andre.

Hva billettholderen kan gjøre, avhenger av arrangementet:

* **Interessenivå**: på arrangementer som er med i puljefordelingen (`events.is_in_puljefordeling = 1`, "Med i puljefordeling" i skjemaet) velger billettholderen hvor interessert hen er: `Veldig interessert`, `Middels interessert` eller `Litt interessert`. Interessenivået brukes når spillere blir fordelt på arrangementer i puljefordelingen, og `Veldig interessert` avgjør om billettholderen får [førstevalg](#førstevalg).

I koden heter dette `interests` (tabell) og `interest_level` (interessenivå), der én rad gjelder én billettholder, ett arrangement og én pulje.

### Førstevalg
En billettholder har fått førstevalg når hen er tildelt med rollen `Player` (spiller) på et arrangement i en pulje, og har gitt `Veldig interessert` på det samme arrangementet i den samme puljen. En `GM`-tildeling alene teller ikke. Er hen både `Player` og `GM` på et arrangement med `Veldig interessert`, teller det, fordi spillerrollen finnes.

I admin vises dette som "Har fått førstevalg" / "Har ikke fått førstevalg". Spillederstatus vises som et eget merke: "Spilleder (GM/DM)" / "Ikke spilleder".

Interesse gjelder én pulje, men førstevalg gjelder for alle puljene samlet. Målet er at alle billettholdere får minst ett førstevalg i løpet av festivalen. Har man fått førstevalg i én pulje, har man fått førstevalget sitt, også i de neste puljene.

I puljefordelingen går billettholdere som ikke har fått førstevalg ennå, foran på arrangementene de har gitt `Veldig interessert`. Jo flere puljer de har hatt `Veldig interessert` på et arrangement uten å få plass, jo høyere blir de prioritert. Å være GM på et arrangement gir ikke førstevalg.

I koden heter dette `first_choice` og `Forstevalg`, og i fordelingen (`solver`) `satisfied` og `top choice`.

## Kode
* _index betyr at filen skal sette opp NATS-integrasjon for domenet
* _page betyr at dette er siden som skal bruke entry #id som blei satt opp i _index til å vise html
