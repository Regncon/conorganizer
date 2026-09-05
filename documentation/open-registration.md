# Open registration

Open registration lets a billettholder take a confirmed player place in a specific event and pulje without waiting for the solver. It is configured per event and applies to every published pulje occurrence of that event.

## User stories

- As a billettholder, I can register directly for an open-registration event in an open pulje so that my place is confirmed immediately.
- As a billettholder, I can keep `Interessert` or `Litt interessert` as a fallback for an open-registration event instead of registering directly.
- As a billettholder, I can register for several open-registration events in the same pulje.
- As a billettholder, I can register for the same event in several puljer independently.
- As a billettholder, I can opt out of my own registration or a manual player assignment while the program is published and the pulje is open.
- As a billettholder with a confirmed assignment, I can see which linked events prevent me from changing ordinary interests in that pulje.
- As an admin, I can mark an event as open registration and manually assign or remove a GM, player, or first-choice player from Puljefordeling.
- As an attendee, I can see confirmed registrations, manual assignments, and GM assignments on my profile without waiting for Puljefordeling to be completed.

## Rules

`Meld deg på` replaces `Veldig interessert` for an open-registration event. `Interessert`, `Litt interessert`, and no interest remain ordinary solver preferences. A confirmed manual or registration player sees `Meld deg av`, including on an event without open registration.

A registration change is allowed only when:

- the program is published;
- the event is announced and its occurrence is published in the selected pulje;
- the pulje is `Open`;
- the user owns the selected billettholder;
- the billettholder is not the GM of an event in the pulje; and
- an underage billettholder is not registering for an 18+ event.

Manual gate check-in remains the source of truth for whether somebody is physically attending. Conorganizer does not derive pulje access from the externally controlled ticket type ID or name and does not add another attendance check to this flow.

## Assignment and interest behavior

Self-registration creates a player assignment with source `registration`. Admin player assignment uses source `manual`; the two workflows remain distinct even though both are confirmed assignments from the billettholder's perspective.

Registering, opting out, or being manually assigned removes only the interest for the same billettholder, event, and pulje. Other interests and assignments remain untouched.

A confirmed manual or registration player assignment in a pulje blocks changes to ordinary interests in that pulje. It does not block registration for another open-registration event. A GM assignment blocks all registration and interest changes in that pulje.

Puljefordeling shows registration assignments as confirmed attendees and excludes those billettholdere from ordinary solver placement in the same pulje. A registration in one pulje does not affect solver eligibility in another pulje. Registration rows keep their `registration` source when an admin commits the distribution.

## Profile visibility

Assignments are shown only after the program is published. Once published:

| Assignment | `Open` | `Locked` | `Completed` |
| --- | --- | --- | --- |
| Manual player | Shown | Shown | Shown |
| Registration player | Shown | Shown | Shown |
| GM | Shown | Shown | Shown |
| Solver player | Hidden | Hidden | Shown |

When an assignment is shown for a pulje, ordinary interest levels for that pulje are hidden.

## Deliberate non-goals

- Open registration is configured per event, not per event/pulje occurrence.
- `max_players` behavior and solver capacity are unchanged by this feature.
- Changing the open-registration flag does not reconcile existing attendees.
- The admin interface does not prevent configuration changes after registration has started.
- Ticket types are not mapped to puljer and ticket names are not parsed.

See [the manual open-registration checklist](./testing/open-registration.md) for end-to-end acceptance checks.
