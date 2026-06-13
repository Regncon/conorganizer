# Hvordan vi skriver manuelle tester

Dette dokumentet beskriver hvordan vi skriver manuelle testsjekklister for Conorganizer. Målgruppen er både utviklere og LLM-baserte kodeassistenter som skal opprette, utvide eller vedlikeholde testfilene i `documentation/testing/`.

Målet er å skrive tester som:

- beskriver observerbar oppførsel
- er nyttige før launch
- er enkle å lese for mennesker
- er lette å videreføre til automatiske tester
- holder seg stabile selv om intern implementasjon endres

## Grunnprinsipper

- Skriv alle tester på bokmål.
- Skriv testene som avkrysningspunkter med kort tittel.
- Skriv `Gitt / Når / Så` på tre egne linjer.
- Test oppførsel, ikke implementasjon.
- Test fra brukerens perspektiv.
- Inkluder både normal oppførsel, edge cases, feilhåndtering og kosmetiske forhold.
- Anta at testeren kan tenke selv. Ikke skriv oppskrifter med detaljert testdata eller trinn-for-trinn-instruksjoner.
- Vær konkret nok til at det er tydelig hva som skal observeres.

## Hva en god test skal beskrive

En god test beskriver:

- en kort tittel som gjør sjekkpunktet lett å skanne
- hvilken rolle som tester
- hvilken situasjon eller forutsetning som gjelder
- hvilken handling som skjer
- hvilket observerbart resultat som forventes

Eksempel:

- [ ] **Min Side viser en helhetlig oversikt**<br>
  **Gitt** at brukeren er innlogget og har tilgang til Min Side.<br>
  **Når** brukeren åpner siden.<br>
  **Så** skal egne arrangementer, billetter og program vises uten at innhold overlapper eller fremstår uferdig.

## Hva vi mener med å teste oppførsel

Vi tester det brukeren kan observere og forholde seg til:

- hva som vises
- hva som kan trykkes på
- hva som lagres
- hva som endres
- hva som skjer ved feil
- hvordan tilgangen oppfører seg
- hvordan siden oppfører seg ved uventede eller tomme data

Vi tester ikke interne tekniske detaljer i den manuelle sjekklisten.

Unngå formuleringer som:

- `SSE-strømmen oppdaterer DOM direkte.`
- `Handleren returnerer 500 når databasen kaster en exception.`

Foretrekk formuleringer som:

- `Oppdatert innhold vises uten ødelagt sidetilstand.`
- `Mislykket lagring gir tydelig beskjed om at innsendingen ikke ble fullført.`

## Format for hver testfil

Hver testfil skal normalt inneholde:

1. En kort tittel.
2. En kort beskrivelse av hva siden eller flyten dekker.
3. En tydelig angivelse av hvilken rolle som skal teste der det er relevant.
4. En sjekkliste gruppert i korte `###`-seksjoner når filen har flere typer oppførsel.
5. Avkrysningspunkter med kort tittel og tre egne `Gitt / Når / Så`-linjer.

## Roller

Roller skal angis i den enkelte testfil der det er relevant. Ikke legg dem bare i en felles hovedfil.

Vanlige roller er:

- ikke-innlogget bruker
- innlogget bruker
- admin

Hvis flere roller er relevante i samme fil, skal det fremgå tydelig i selve teksten eller i korte underoverskrifter.

## Hvordan vi skriver sjekkpunkter

Hvert sjekkpunkt skal være:

- selvstendig
- tydelig
- observerbart
- relevant for launch

Bruk denne formen:

- [ ] **Kort tittel på forventet oppførsel**<br>
  **Gitt** at relevant forutsetning finnes.<br>
  **Når** brukeren gjør handlingen eller siden vises.<br>
  **Så** skal observerbart resultat være tydelig.

Foretrekk:

- ett tydelig fokus per punkt
- konkrete forventninger
- formuleringer som sier noe om utfallet
- titler som beskriver resultatet, ikke bare forutsetningen

Unngå:

- vage formuleringer som `fungerer som forventet`
- punkter som tester flere uavhengige ting samtidig
- formuleringer som er så tekniske at de ikke lenger beskriver brukeropplevelsen

Bra:

- [ ] **Arrangementer ligger under riktig pulje**<br>
  **Gitt** at forsiden inneholder arrangementer i flere puljer.<br>
  **Når** brukeren åpner forsiden.<br>
  **Så** skal hvert arrangement vises under riktig pulje og med lesbar informasjon om tidspunkt og innhold.

Mindre bra:

`Test at forsiden fungerer.`

## Feilhåndtering og edge cases

Hver fil skal inneholde egne punkter for edge cases og feilhåndtering. Disse skal ikke skyves over i en egen generell restliste hvis de hører naturlig hjemme i den konkrete siden eller flyten.

Vi skal blant annet tenke på:

- tomme tilstander
- manglende data
- ugyldige data
- delvis data
- manglende tilgang
- mislykket lagring
- uventet navigasjon
- refresh og tilbakeknapp
- flere raske handlinger etter hverandre
- tilstander som ser teknisk riktige ut, men som er uforståelige for brukeren

## Kosmetiske forhold

Kosmetiske forhold skal være med i sjekklistene, men de skal fortsatt beskrives som oppførsel.

Bra:

- [ ] **Mobilvisning er lesbar og navigerbar**<br>
  **Gitt** at siden vises på mobil.<br>
  **Når** brukeren scroller og navigerer.<br>
  **Så** skal tekst, knapper, ikoner og kort være lesbare og ikke overlappe eller havne utenfor skjermen.

Mindre bra:

`Se etter CSS-feil.`

## Automatisering og rapport

Automatiseringskandidater skal ikke lagres som egne seksjoner i de manuelle testfilene. Slike lister blir raskt utdaterte.

Når et manuelt punkt automatiseres, skal den automatiserte testen ha en tydelig BDD-kommentar øverst i testen. Kjør `task test:report` for å se hvilke automatiserte tester som finnes, hvilke BDD-kommentarer de har, og hvilke tester som mangler BDD-kommentar.

## For utviklere

Når du oppdaterer funksjonalitet, skal du vurdere om en eksisterende testfil må oppdateres.

Spør:

- har siden fått ny observerbar oppførsel?
- har en gammel oppførsel blitt endret eller fjernet?
- finnes det en ny feiltilstand eller edge case?
- bør noe som i dag testes manuelt flyttes til automatiserte tester?

## For LLM-baserte kodeassistenter

Når du skriver eller oppdaterer testfiler:

- følg strukturen i dette dokumentet
- skriv på bokmål
- bruk `Gitt / Når / Så`
- prioriter observerbar oppførsel fremfor implementasjonsdetaljer
- unngå å dikte opp funksjonalitet som ikke finnes i kodebasen
- hold deg til ruter, roller og brukerflyter som faktisk finnes
- sørg for at edge cases og feilhåndtering er med i hver relevant fil

Hvis en flyt er deprecated eller ikke en del av launch, skal den ikke inn i launch-sjekklistene.

## Kvalitetssjekk før en testfil regnes som ferdig

Før en ny eller oppdatert testfil anses som ferdig, skal den kunne bestå denne kontrollen:

- Dekker filen den faktiske siden eller flyten som finnes i appen?
- Er teksten skrevet på bokmål?
- Er punktene skrevet som avkrysningspunkter?
- Har hvert punkt en kort tittel?
- Bruker punktene tre egne `Gitt / Når / Så`-linjer?
- Er punktene brukerorienterte og observerbare?
- Inneholder filen edge cases og feilhåndtering?
- Inneholder filen kosmetiske forhold der det er relevant?
- Er det tydelig hvilken rolle som tester?
- Er automatiserte punkter flyttet til automatiserte tester med BDD-kommentarer?
