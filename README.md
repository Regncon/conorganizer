# Con Organizer

## Description

This is the source code for the Regncon 2024 app. It's a work in progress, but the idea is to have a place to keep track
of your con schedule, your panels, and your expenses.

## referat fra post-mortem: sorteres inn i must-haves og nice-to-haves av neste styre

-   bedre loginsystem? Passwordless login?
-   tydeligere hva pameldingen gjor. tooltip? Bedre ord enn "pamelding"
-   klokkeslett klar for connet
-   bedre grensesnitt og automatisering av opprop? Forslag: liste med deltakere, liste med queue, mest populare forst
-   statistikk pa slutten av connet: hvor mange fikk spilt ting de hadde hvor lyst til? hvor mange dukket ikke opp? osv
-   britisk flagg som funker i alle nettlesere
-   pameldingsstatus tilbake til listviewen - hvordan med flere brukere? Bare prikker eller ikoner?
-   bash-script eller noe annen automatisering for billedbehandling. F.eks: https://www.npmjs.com/package/sharp
-   mulighet for at folk kan legge inn arrangement selv (som skjult, sa kan styret godkjenne og evt endre)
    -   obs: folk ma vare innlogget, sa de ma ha kjopt billett?
    -   obs: hvor mye tekstformatering skal vi ha med? Fet og kursiv, f.eks? Lenker?
    -   obs: fint hvis folk kan laste opp bilder selv og f.eks. bruke https://www.npmjs.com/package/cropperjs
    -   fint hvis folk kan laste opp modulen ogsa, og sjekke av for at den skal vare med i modulkonkurranse
    -   mulighet for a trekke arrangementet?
-   refakturere og/eller lage fra scratch med andre teknologier
-   brukertesting av grensesnittet, gi folk oppdrag som "meld deg pa the One Ring" og se hva som skjer

## Todo

### Must haves for opprop

-   See favorites in participant admin
-   See assigned games in participant admin
-   See assigned games in events
-   Manually assign participants to games
-   Generate assignment suggestions
-   Lock pool
-   Add is gm to event
-   Impersonate participants to assign interests


#### Nice to haves for publishing events

-   Fix cookie lifetime?
-   Filter my signups, events and tags
-   Check preloading of pages
-   Debounce/fix text box, my events
-   Display label on new or unread events

## Nice to haves

-   Fix small card layout
-   Picture upload
-   Add keyboard navigation
-   Add command palatte
-   Add search
-   Tags
-   Add event to calendar
-   storskjermversjon som viser resten av dagen automatisk

## Done

-   Icons | Grethe :D
-   Help text explaining how the algorithm works
-   Bug: Problems with hitting enter while editing event
-   Add picture url | Gerhard
-   Log-in bug, e-mail input loses focus when you input a letter | Torstein
-   Fix events front page | Torstein
-   Desktop view | Torstein
-   Convert my events to use new event system
-   Add none events to list
-   Match header and footer to design
-   Check if firebase token should be private
-   Add media queries mobile first
-   Add event images
-   Add error boundary
-   Save user choices for events
-   Add timeslots
-   Add rooms
-   Add favicon
-   Set page title
-   Authorization
-   Auth
-   query string for event id
-   Add event types
-   Add read more
-   Unlist event
-   Get event from database
-   Add new event
-   Edit event
-   Easy next and previous navigation
-   Make it very clear what day is displayed
-   Get tickets from checkin
-   Room assignment | Gerhard
-   Legge in beta logo
-   Switch default card image
-   Google log-in redirect bug
-   Disable interest function
-   Legge til tags i events og pool-event
-   Participant admin page | Gerhard
-   Create participant in database
-   Check if participant is over 18 |
-   Connect to checkin system
-   filere ut events som ikke er publisert
-   add button to redirect to administration page of event
-   Assign user to participant
-   Fix race condition in assignParticipantByEmail
-   Tildel billett til bruker (Super nesten ferdig, "bank i bordet")
-   Vise hvilken deltaker du er viss du har flere deltakere på samme billett (Nesten ferdig)
-   Gjør det tydelig hvilke pulje arrangementet er i og om det kjøres i andre puljer også.
-   Lagre interessene til deltagere i databasen.- Assign user to participant
-   See favorites in events
