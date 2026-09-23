# Domeneordbok

## Ord
Interesse
Interesser
Interessenivå
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
En pulje er tidspunkt der alle arrangementer som skal spiller innen for tidspunktet som styret har valg For eksempel Fredag kveld: 18 - 23

## Interesse
Interesse er det en billettholder melder inn på et arrangement i en pulje. Den gjelder bare for det arrangementet i den valgte puljen. Samme arrangement kan gå i flere puljer, og da er interessen i hver pulje uavhengig av de andre.

Ordet har to betydninger, avhengig av arrangementet:

* **Interessenivå**: på vanlige arrangementer velger billettholderen hvor interessert hen er: `Veldig interessert`, `Middels interessert` eller `Litt interessert`. Interessenivået brukes når spillere blir fordelt på arrangementer i puljefordelingen.
* **Påmelding**: på arrangementer uten interessevelger (for eksempel cosplay) melder billettholderen seg bare på. Det finnes ikke noe nivå å velge. Se [Påmelding](#påmelding).

I koden heter dette `interests` (tabell) og `interest_level` (interessenivå), der én rad gjelder én billettholder, ett arrangement og én pulje.

## Påmelding
Påmelding er noe en spiller kan gjøre på arrangementer som er langvarige arrangementer, som alle så vil kan melde seg på.Som for eksempel "Blood on the clock tower" eller "Cosplay"

## Kode
* _index betyr at filen skal sette opp NATS-integrasjon for domenet
* _page betyr at dette er siden som skal bruke entry #id som blei satt opp i _index til og vise html
