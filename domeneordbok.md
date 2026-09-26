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

Ordet har to betydninger, avhengig av arrangementet:

* **Interessenivå**: på vanlige arrangementer velger billettholderen hvor interessert hen er: `Veldig interessert`, `Middels interessert` eller `Litt interessert`. Interessenivået brukes når spillere blir fordelt på arrangementer i puljefordelingen, og `Veldig interessert` avgjør om billettholderen får [førstevalg](#førstevalg).
* **Påmelding**: på arrangementer uten interessevelger (for eksempel cosplay) melder billettholderen seg bare på. Det finnes ikke noe nivå å velge. Se [Påmelding](#påmelding).

I koden heter dette `interests` (tabell) og `interest_level` (interessenivå), der én rad gjelder én billettholder, ett arrangement og én pulje.

## Førstevalg
En billettholder har fått førstevalg når hen er tildelt som spiller på et arrangement hen har gitt `Veldig interessert`.

Interesse gjelder én pulje, men førstevalg gjelder for alle puljene samlet. Målet er at alle billettholdere får minst ett førstevalg i løpet av festivalen. Har man fått førstevalg i én pulje, har man fått førstevalget sitt, også i de neste puljene.

I puljefordelingen går billettholdere som ikke har fått førstevalg ennå, foran på arrangementene de har gitt `Veldig interessert`. Jo flere puljer de har hatt `Veldig interessert` på et arrangement uten å få plass, jo høyere blir de prioritert. Å være GM på et arrangement gir ikke førstevalg.

I koden heter dette `first_choice` og `Forstevalg`, og i fordelingen (`solver`) `satisfied` og `top choice`.

## Påmelding
Påmelding er noe en spiller kan gjøre på arrangementer som er langvarige arrangementer, som alle som vil, kan melde seg på. For eksempel "Blood on the clock tower" eller "Cosplay"

## Kode
* _index betyr at filen skal sette opp NATS-integrasjon for domenet
* _page betyr at dette er siden som skal bruke entry #id som blei satt opp i _index til å vise html
