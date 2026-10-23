# Age limit (18+)

This document describes how the 18+ age limit (aldersgrense) is modelled and where the application warns about billettholdere under 18. The board handles the age limit by hand: nothing in the application stops a billettholder under 18 from showing interest in, or being seated on, an 18+ event.

## The two flags

| Flag | Where | Meaning |
| --- | --- | --- |
| `events.age_group = 'AdultsOnly'` (`models.AgeGroupAdultsOnly`) | `models/event-model.go` | The event has an 18+ age limit. The event form's select shows it as "Voksne (18+)" (`AgeGroup.Label()`). Event cards, the event page and the event tags on the admin puljefordeling page show `icons.EventAdultOnly` labelled "Egnet for voksne" (`AgeGroup.BadgeLabel()`): as visible text next to the icon on the event page, as a tippy tooltip on event cards, and as the icon's title and aria-label on the puljefordeling page. |
| `billettholdere.is_over_18` | `service/checkIn/cache.go` | Set from the CheckIn ticket's birth date by `isOver18`: true when the billettholder is 18 or older on `regnConDate` (2026-10-04). A missing or unparseable birth date counts as under 18. Billettholder cards in admin and on the profile show "Alder: Over 18" or "Alder: Under 18". See [billettholdere.md](billettholdere.md). |

`regnConDate` changes every festival (AGENTS.md, "Festival year").

## The solver ignores age

The puljefordeling solver does not look at `is_over_18` or `age_group`, so it can seat a billettholder under 18 as a player on an 18+ event. This is a board decision (#618 reversed #524). Do not add age back into the solver. The application warns the admin instead.

## Warnings on the admin puljefordeling page

`/admin/puljefordeling/{pulje}` (`pages/admin/puljefordeling_tab.templ`) shows three things:

- **Aldersvarsel** under "Varsler om tildelinger" (`puljeFordelingsvarsler` in `pages/admin/puljefordeling_varsler.templ`): the heading "Under 18 på 18+-arrangement", the text "Disse spillerne er under 18 år og er satt opp på et arrangement med 18-årsgrense. Aldersgrensen håndteres manuelt, så sjekk dem før puljen publiseres.", and one collapsible entry per event with the count ("N under 18") and the name of every player under 18. `FinnAldersvarsler` (`service/puljefordeling/aldersvarsler.go`) builds the list from the displayed pulje, so it covers solver seats and manual pins alike. GMs are not included. Events are sorted by title, players by name.
- **The "Under 18" badge** (`pulje-badge pulje-badge--error`) next to every player and GM under 18 on an 18+ event card.
- **The "Under 18" chip** in the assignment dialog's candidate list (`pages/admin/puljefordeling_assignment_dialog.templ`), shown for every candidate under 18 whatever the event's age group.

## Confirming a manual pin

Before an admin pins a billettholder under 18 on an 18+ event, the pin must be confirmed:

- `buildTildelingsvarsel` (`service/puljefordeling/tildelinger.go`) fills `Tildelingsvarsel.Aldersvarsel` with "Navn er under 18 år, og «Tittel» er 18+." (the name without guillemets, the title in «») for the chosen event, and adds one such sentence, space-separated, for any other 18+ event in the pulje where the same pin turns the billettholder's solver seat into a manual seat.
- The assignment dialog and drag-and-drop both post to `/admin/api/puljefordeling/assign/preview` first. The preview always opens the assignment preview dialog (`pages/admin/tildelingsdialog.templ`), which shows the warning as `.tildeling-age-warning`, and the admin confirms the whole assignment there. Its confirm button sends the matching confirmation token (`Bekreftelse`, signal `assignmentConfirmation`) together with `assignmentAgeConfirmed`.
- `POST /admin/api/puljefordeling/assign` saves the pin when the request carries a matching token, or, without a token, when `AlderBekreftet` (signal `assignmentAgeConfirmed`) is set. Without either, and when the billettholder has no other assignments in the pulje and there is no capacity warning, `sendTildelingsvarsel` (`pages/admin/tildelinger.go`) patches `$ageWarningText`, which opens the fallback "Aldersgrense" dialog (`puljeAgeWarningDialog` in `puljefordeling_tab.templ`) with the buttons "Avbryt" and "Plasser likevel". "Plasser likevel" posts again with `assignmentAgeConfirmed` set.

## The event page

Nobody is stopped from showing interest. `LoadInterestNoticeState` (`components/event_components/event_interests.templ`) sets `ShowUnder18` when the billettholder is under 18, the event is `AdultsOnly` and the billettholder has no manually pinned Player seat and no GM assignment in that pulje (a solver seat does not count). The interest dialog then shows the notice "Anbefalt for voksne (18+)", and the billettholder can still choose and save an interest level. Manual test: [testing/event-details.md](testing/event-details.md).
