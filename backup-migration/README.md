# Backup migration tool
> [!NOTE]
> Hvis du har en `.env` i samme dir som applikasjonen vill den prøve å laste den inn automatisk når du starter

Super enkel gui for å gjøre ned og opplasting av databaser for Regncon 2025.

For best mulig funksjonalitet anbefaler jeg at du bygger programmet. Krever at du har [Go installert](https://go.dev/doc/install)
```console
go build . && /.backup-migration
```


## Manuel migration med Goose
> [!TIP]
> Hvis du ønsker å kjøre Goose via cli kan du finne installasjons veiledning her [link](https://pressly.github.io/goose/installation/)

Som en del av migrasjons prossesen bruker vi [Goose](https://github.com/pressly/goose) som er tilgjengelig via både go functions og som CLI verktøy.
Backup-migration applikasjonen bruker dette programmet for å gjøre nødvendige oppgraderinger av databasen automatisk. Filene for migrasjonen eksisterer i `/migrations`.

Du kan bruke backup-migration applikasjonen til å laste ned nyeste versjon av regncon databasen. Filen vill lagres i samme mappe som binærfilen kjøres fra og få følgende navn struktur `regncon-<date>.db`


## Goose command line
> [!NOTE]
> Hvis du har en `.env` med [goose variabler](https://pressly.github.io/goose/documentation/environment-variables/) vill cli verktøyet prøve å lese disse

### Lag ny migrasjon
Hvis du ønsker å gjøre en endring i databasen må du først lage migrasjonsfilen, følgende kommando lager en `.sql` fil som du kan skrive inn endringene i.
Gjerne bruk et forklarende navn med snake_case, f.eks: `events_updating_status_default` eller `users_adding_age`

```console
goose create <subject> sql
```

Les mer [her](https://pressly.github.io/goose/documentation/annotations/) for veiledning og eksempler i å skrive migrasjonsfiler i sql

### Oppdater database med alle migrasjons steg
Følgende kommando går gradvis igjennom alle migrasjonsfiler og gjør nødvendige endringer i databasen

```console
goose <db file path> up
```

### Hvordan kan jeg gå tilbake etter oppdatering?
Du kan bruke følgende kommando for å returnere databasen sin orginale versjon
```console
goose <db file path> reset
```
