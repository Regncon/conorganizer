#!/usr/bin/env bash
# Builds a small conorganizer-like git checkout for the commit message evals.
#
# Usage: bash make-fixture-repo.sh <dest> <uncommitted-retry|bad-messages>
#
# uncommitted-retry: branch login-sync-retry with an unstaged retry change.
# bad-messages: branch pulje-closing-warning with four commits, three of
# them with messages that break the house rules.
set -eu

dest="${1:?usage: make-fixture-repo.sh <dest> <uncommitted-retry|bad-messages>}"
scenario="${2:?usage: make-fixture-repo.sh <dest> <uncommitted-retry|bad-messages>}"

rm -rf "$dest"
mkdir -p "$dest"
cd "$dest"
git init -q -b main
git config user.name "Kari Nordmann"
git config user.email "kari@example.com"
git config core.autocrlf false

stamp=1760000000
commit() {
	stamp=$((stamp + 3600))
	GIT_AUTHOR_DATE="@$stamp +0100" GIT_COMMITTER_DATE="@$stamp +0100" git commit -q -F -
}

mkdir -p service/checkIn service/puljer pages/profile/tickets static/css documentation

cat > AGENTS.md <<'EOF'
## Names
- Code is in English, except the domain words in domeneordbok.md, which stay Norwegian: `getBillettholder`, not `getTicketHolder`.
- Write domain words without æøå in identifiers, as domeneordbok.md does: `Forstevalg`, `Lordag`.

## Text and names
- User-facing text is bokmål, never nynorsk.
- An event is "annonsert", never "publisert". "Publisert" is for the program, the puljefordeling and the romfordeling.

## Git
- Commit subjects are imperative English and describe the user-visible outcome, keeping Norwegian domain words. No Co-Authored-By trailer.
- Bring a branch up to date by merging main into it, not by rebasing.
EOF

cat > domeneordbok.md <<'EOF'
# Domeneordbok

- **Billettholder**: a person with a ticket to Regncon. One user can be linked to several billettholdere.
- **Pulje**: a time slot in the festival programme. Arrangementer run inside a pulje.
- **Puljefordeling**: the assignment of billettholdere to arrangementer in a pulje.
- **Interesse**: a billettholder's wish to join an arrangement in a pulje. Never "påmelding": nobody signs up, they show interesse and the puljefordeling decides.
- **Førstevalg**: the one interesse a billettholder marks as their first choice in a pulje.
EOF

cat > README.md <<'EOF'
# Conorganizer

Conorganizer runs the programme for Regncon: billetholdere show interesse in arrangementer, and the puljefordeling places them. The programme is publisert when the board says so.

## Dev setup

Run `go tool task start` and open http://localhost:8080. A fresh checkout has no database, so download data first with `go tool task download:db`.
EOF

cat > static/css/index.css <<'EOF'
:root {
	--color-text: #1b1b1b;
	--color-accent-yellow: #f4c542;
	--banner-bg: #2a2a2a;
}
EOF

cat > service/checkIn/import.go <<'EOF'
package checkIn

import (
	"context"
	"fmt"
)

type Client interface {
	TicketsByEmail(ctx context.Context, email string) ([]Ticket, error)
}

type Ticket struct {
	ID    string
	Email string
	Type  string
}

func ImportUserTickets(ctx context.Context, client Client, email string) ([]Ticket, error) {
	tickets, err := client.TicketsByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("fetch tickets for login sync: %w", err)
	}
	return tickets, nil
}
EOF

cat > service/puljer/closing.go <<'EOF'
package puljer

import (
	"fmt"
	"time"
)

type Pulje struct {
	Name     string
	StartsAt time.Time
	ClosesAt time.Time
}

func ClosingWarning(p Pulje, now time.Time) string {
	return fmt.Sprintf("Påmeldingen til %s lukkes %s", p.Name, p.ClosesAt.Format("02.01.2006 15:04"))
}
EOF

cat > service/puljer/status.go <<'EOF'
package puljer

import (
	"context"
	"fmt"
)

type Broadcaster interface {
	Broadcast(ctx context.Context, bucket string) error
}

const BucketEvents = "events"

func broadcastStatus(ctx context.Context, b Broadcaster) error {
	if err := b.Broadcast(ctx, BucketEvents); err != nil {
		return fmt.Errorf("broadcast pulje status: %w", err)
	}
	return nil
}
EOF

cat > documentation/billettholdere.md <<'EOF'
# Billettholdere

The login sync (checkIn.ImportUserTickets) fetches the tickets CheckIn has for the user's e-mail and links the user to each billettholder it finds.
EOF

git add -A
commit <<'EOF'
Add the billettholder import from CheckIn
EOF

cat >> documentation/billettholdere.md <<'EOF'

A user is linked to a billettholder only through relation_billettholdere_users. Every read, picker and access check uses that link; an e-mail is only used to create it.
EOF
git add -A
commit <<'EOF'
Link billettholdere to users by id everywhere except ticket import

An e-mail is text a user typed, not a key, so matching it against the
logged-in user let anyone see the billettholdere behind an address they
could type. Every read and access check now goes through
relation_billettholdere_users; only the CheckIn import still uses the
e-mail, because that is where the link is created.
EOF

cat > pages/profile/tickets/tickets_page.templ <<'EOF'
package tickets

templ NoSeatMessage(puljeName string) {
	<p class="no-seat">Du fikk ikke plass i { puljeName } denne gangen.</p>
}
EOF
git add -A
commit <<'EOF'
Show no-seat message on profile for published puljer
EOF

sed -i 's/The programme is publisert when the board says so./The programme is publisert when the board says so; an arrangement is annonsert./' README.md
git add -A
commit <<'EOF'
Stop calling events publisert in user-facing text

domeneordbok.md reserves "publisert" for the program, the puljefordeling
and the romfordeling. An event is "annonsert".
EOF

case "$scenario" in
uncommitted-retry)
	git checkout -q -b login-sync-retry
	cat > service/checkIn/import.go <<'EOF'
package checkIn

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Client interface {
	TicketsByEmail(ctx context.Context, email string) ([]Ticket, error)
}

type Ticket struct {
	ID    string
	Email string
	Type  string
}

var ErrTransient = errors.New("checkin unavailable")

const fetchAttempts = 3
const fetchRetryDelay = 500 * time.Millisecond

func ImportUserTickets(ctx context.Context, client Client, email string) ([]Ticket, error) {
	var lastErr error
	for attempt := 1; attempt <= fetchAttempts; attempt++ {
		tickets, err := client.TicketsByEmail(ctx, email)
		if err == nil {
			return tickets, nil
		}
		lastErr = err
		if !errors.Is(err, ErrTransient) {
			break
		}
		time.Sleep(fetchRetryDelay)
	}
	return nil, fmt.Errorf("fetch tickets for login sync: %w", lastErr)
}
EOF
	cat > service/checkIn/import_test.go <<'EOF'
package checkIn

import (
	"context"
	"testing"
)

type flakyClient struct {
	failures int
	calls    int
}

func (c *flakyClient) TicketsByEmail(ctx context.Context, email string) ([]Ticket, error) {
	c.calls++
	if c.calls <= c.failures {
		return nil, ErrTransient
	}
	return []Ticket{{ID: "t-1", Email: email, Type: "Festivalpass"}}, nil
}

func TestImportUserTickets_RetriesAfterTransientCheckInError(t *testing.T) {
	// Given
	client := &flakyClient{failures: 2}

	// When
	tickets, err := ImportUserTickets(context.Background(), client, "kari@example.com")

	// Then
	if err != nil {
		t.Fatalf("expected the sync to recover, got %v", err)
	}
	if len(tickets) != 1 || client.calls != 3 {
		t.Fatalf("expected one ticket after three calls, got %d tickets after %d calls", len(tickets), client.calls)
	}
}
EOF
	;;
bad-messages)
	git checkout -q -b pulje-closing-warning
	cat > service/puljer/closing.go <<'EOF'
package puljer

import (
	"fmt"
	"time"
)

type Pulje struct {
	Name     string
	StartsAt time.Time
	ClosesAt time.Time
}

func ClosingWarning(p Pulje, now time.Time) string {
	if now.After(p.StartsAt.Add(time.Hour)) {
		return ""
	}
	return fmt.Sprintf("Påmeldingen til %s lukkes %s", p.Name, p.ClosesAt.Format("02.01.2006 15:04"))
}
EOF
	git add -A
	commit <<'EOF'
Hide pulje alerts one hour after the pulje starts

An alert for a pulje that is already running told people to hurry to
something they could no longer join. The bar drops the alert an hour
after the pulje starts, when the last late arrivals have been seated.
EOF

	cat > service/puljer/closing.go <<'EOF'
package puljer

import (
	"fmt"
	"time"
)

type Pulje struct {
	Name     string
	StartsAt time.Time
	ClosesAt time.Time
}

func ClosingWarning(p Pulje, now time.Time) string {
	if p.ClosesAt.IsZero() {
		return ""
	}
	if now.After(p.StartsAt.Add(time.Hour)) {
		return ""
	}
	return fmt.Sprintf("Påmeldingen til %s lukkes %s", p.Name, p.ClosesAt.Format("02.01.2006 15:04"))
}
EOF
	cat > service/puljer/closing_test.go <<'EOF'
package puljer

import (
	"testing"
	"time"
)

func TestClosingWarning_EmptyWhenPuljeHasNoClosingTime(t *testing.T) {
	// Given
	now := time.Date(2026, 11, 3, 12, 0, 0, 0, time.UTC)
	p := Pulje{Name: "Fredag kveld", StartsAt: now.Add(48 * time.Hour)}

	// When
	got := ClosingWarning(p, now)

	// Then
	if got != "" {
		t.Fatalf("expected no warning for a pulje without a closing time, got %q", got)
	}
}
EOF
	git add -A
	commit <<'EOF'
fix the thing.
EOF

	sed -i 's/"Påmeldingen til %s lukkes %s"/"Du kan melde interesse for %s fram til %s"/' service/puljer/closing.go
	sed -i 's/--banner-bg: #2a2a2a;/--banner-bg: color-mix(in srgb, var(--color-accent-yellow) 20%, #2a2a2a);/' static/css/index.css
	sed -i 's/billetholdere/billettholdere/' README.md
	git add -A
	commit <<'EOF'
651 | Updated the pulje closing warning text and also changed the banner color and fixed a typo in the readme

- changed text in closing.go
- changed the banner color in index.css
- fixed typo in README

Co-Authored-By: Claude <noreply@anthropic.com>
EOF

	cat > service/puljer/status.go <<'EOF'
package puljer

import (
	"context"
	"fmt"
	"time"
)

type Broadcaster interface {
	Broadcast(ctx context.Context, bucket string) error
}

const BucketEvents = "events"

const broadcastAttempts = 3
const broadcastRetryDelay = 500 * time.Millisecond

func broadcastStatus(ctx context.Context, b Broadcaster) error {
	var err error
	for attempt := 1; attempt <= broadcastAttempts; attempt++ {
		err = b.Broadcast(ctx, BucketEvents)
		if err == nil {
			return nil
		}
		time.Sleep(broadcastRetryDelay)
	}
	return fmt.Errorf("broadcast pulje status: %w", err)
}
EOF
	cat > service/puljer/status_test.go <<'EOF'
package puljer

import (
	"context"
	"errors"
	"testing"
)

type reconnectingBroadcaster struct {
	failures int
	calls    int
}

func (b *reconnectingBroadcaster) Broadcast(ctx context.Context, bucket string) error {
	b.calls++
	if b.calls <= b.failures {
		return errors.New("nats: connection reconnecting")
	}
	return nil
}

func TestBroadcastStatus_RetriesWhileNATSReconnects(t *testing.T) {
	// Given
	b := &reconnectingBroadcaster{failures: 2}

	// When
	err := broadcastStatus(context.Background(), b)

	// Then
	if err != nil {
		t.Fatalf("expected the broadcast to go through once NATS is back, got %v", err)
	}
	if b.calls != 3 {
		t.Fatalf("expected three attempts, got %d", b.calls)
	}
}
EOF
	git add -A
	commit <<'EOF'
Add retry to the pulje status broadcast

Added a for loop in broadcastStatus that retries three times with a
500ms sleep between attempts. Changed the function to take the number
of attempts. Updated the test.
EOF
	;;
*)
	echo "unknown scenario: $scenario" >&2
	exit 1
	;;
esac

echo "fixture repo ready in $dest ($scenario)"
