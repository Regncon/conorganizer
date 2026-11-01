# Live page skeleton

Everything below comes from `pages/profile/tickets`, the cleanest live page in the repo. Sections 1 to 3 match the real files, except that the real `/api` sub-router is called `billettholderAdminRouter`. Section 4 is shortened. Sections 5 and 7 rewrite the real delete and add-e-mail handlers into the shape new code should have. Sections 6, 8 and 9 use the fixture helpers that already exist in `pages/profile/tickets/billettholder_email_routes_test.go`. The Go code in sections 5 to 9 was compiled and run against the repo. Replace `profileTickets`/`profile-tickets` with your page and keep the structure.

Contents:

1. Route setup: normal GET, `/api` stream, mutations
2. The `_index` wrapper: full content plus `data-init`
3. The `_page` renderer: same id, no `data-init`
4. The shared content component
5. A mutation handler in the right order
6. Route test with embedded NATS
7. Field errors the user should read
8. Testing field errors
9. End-to-end: an open page re-renders after a mutation
10. Doc rows to add

## 1. Route setup: normal GET, `/api` stream, mutations

`pages/profile/tickets/tickets_index.templ`:

```go
func ProfileTicketsRoute(router chi.Router, liveManager *live.Manager, db *sql.DB, logger *slog.Logger) error {
	logger = logger.With("component", "profile_tickets")
	router.Route("/tickets", func(profileRouter chi.Router) {
		profileRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			var ctx = r.Context()
			var user = userctx.GetUserRequestInfo(ctx)

			if err := layouts.Base(
				"Mine Billetter",
				user,
				db,
				logger,
				profileTicketsIndex(user.Id, db, logger),
			).Render(ctx, w); err != nil {
				logger.Error(fmt.Errorf("failed to render profile tickets page: %w", err).Error(), "user_id", user.Id)
			}
		})
		profileRouter.Route("/api", func(apiRouter chi.Router) {
			apiRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
				liveManager.Stream(w, r, live.Page{
					Buckets: []live.Bucket{live.BucketBillettholders},
					Render: func(ctx context.Context, r *http.Request) templ.Component {
						user := userctx.GetUserRequestInfo(ctx)
						return profileTicketsPage(user.Id, db, logger)
					},
				})
			})
			addEmailToBilettholderRoute(apiRouter, db, logger, liveManager)
			deleteEmailFromBillettholderRoute(apiRouter, db, logger, liveManager)
			getTicketsRouter(apiRouter, db, logger, liveManager)
		})
	})
	return nil
}
```

Notes:

- The `/api` route sits under the same router as the page, so login middleware applies. Per-resource access, such as ownership of an event, is not inherited. Check it before `Stream` and in every mutation under `/api`, for example with a middleware added by `Use` on the `/api` sub-router, as `newevent.RequireEventOwnerOrAdmin` does in `pages/profile/profile.go`.
- Validate URL parameters before `Stream`. The room assignment stream in `pages/admin/admin.go` parses the pulje and answers `http.Error(..., http.StatusBadRequest)` first, because once `Stream` runs, the headers are flushed.
- Query state the page depends on is read in the route and closed over, or read from `r` inside `Render`: `pages/root/root.go` reads `requestedDate := r.URL.Query().Get(programDateQueryParam)` and passes it to `rootPage`.

## 2. The `_index` wrapper: full content plus `data-init`

```templ
templ profileTicketsIndex(userId string, db *sql.DB, logger *slog.Logger) {
	<div id="profile-tickets-page-wrapper" data-init={ live.DatastarInit("/profile/tickets/api") }>
		@profileTicketsPageContent(userId, db, logger)
	</div>
}
```

When the live URL depends on the browser, use the expression form, as `pages/profile/profile_index.templ` does:

```templ
data-init={ live.DatastarInitExpression("'/profile/api' + window.location.search") }
```

When it depends on the request, build it in Go, as `pages/event/event_index.templ` does:

```templ
data-init={ live.DatastarInit(fmt.Sprintf("/event/api/%s?%s", eventId, r.URL.RawQuery)) }
```

## 3. The `_page` renderer: same id, no `data-init`

```templ
templ profileTicketsPage(userId string, db *sql.DB, logger *slog.Logger) {
	<div id="profile-tickets-page-wrapper">
		@profileTicketsPageContent(userId, db, logger)
	</div>
}
```

`Stream` sends `sse.PatchElementTempl(page.Render(ctx, r))`. The default patch mode is an outer morph matched by the top-level element's id, so the `_index` wrapper and the `_page` wrapper must share the id. The content component is shared between both so the first HTTP response and every later patch render the same thing.

## 4. The shared content component

`profileTicketsPageContent` queries SQLite each time it renders and logs query failures itself. Shortened:

```templ
templ profileTicketsPageContent(userId string, db *sql.DB, logger *slog.Logger) {
	<section class="flex flex-col page-content-container">
		<h1 class="page-heading">Mine billetter</h1>
		{{ billettholdere, err := billettholderService.GetBillettholdere(userId, db) }}
		if err != nil {
			{{ logger.Error(err.Error(), "user_id", userId) }}
			<p>Feil ved henting av billetter: { err.Error() }</p>
			return
		}
		for _, billettholder := range billettholdere {
			@billettholderCard(billettholder)
		}
	</section>
}
```

It re-renders on every broadcast to `billettholders`, on every tab-visibility change and whenever another tab in the same browser opens a page that touches the shared connection key. Nothing in it may depend on "this is the first render".

## 5. A mutation handler in the right order

The repo's e-mail routes answer a failed check with a 200 SSE that patches `errorMessage-{id}`, and then call `http.Error`, which only logs a superfluous `WriteHeader`. This is the same delete handler in the order this skill asks for: read and validate, check access, write, broadcast, then SSE. Use this shape for new code.

```go
var errEmailNotManual = errors.New("only manual e-mail addresses can be deleted")

func deleteEmailFromBillettholderRoute(router chi.Router, db *sql.DB, logger *slog.Logger, liveManager *live.Manager) {
	logger = logger.With("component", "profile_tickets")
	router.Post("/delete-email/{id}/{emailID}/", func(w http.ResponseWriter, r *http.Request) {
		billettholderID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, "Ugyldig billettholder", http.StatusBadRequest)
			return
		}
		emailID := chi.URLParam(r, "emailID")
		user := userctx.GetUserRequestInfo(r.Context())

		linked, err := billettholderLinkedToUser(db, billettholderID, user.Id)
		if err != nil {
			logger.Error(fmt.Errorf("failed to check billettholder access: %w", err).Error(), "billettholder_id", billettholderID)
			http.Error(w, "Noe gikk galt, prøv igjen", http.StatusInternalServerError)
			return
		}
		if !linked {
			http.Error(w, "Du har ikke tilgang", http.StatusForbidden)
			return
		}

		email, err := deleteManualEmail(db, billettholderID, emailID)
		if errors.Is(err, errEmailNotManual) {
			http.Error(w, "Kun manuelle epostadresser kan slettes", http.StatusConflict)
			return
		}
		if err != nil {
			logger.Error(fmt.Errorf("failed to delete billettholder email: %w", err).Error(), "billettholder_id", billettholderID)
			http.Error(w, "Noe gikk galt, prøv igjen", http.StatusInternalServerError)
			return
		}

		if err := liveManager.Broadcast(r.Context(), live.BucketBillettholders); err != nil {
			logger.Error(fmt.Errorf("failed to broadcast billettholder email delete: %w", err).Error(), "billettholder_id", billettholderID)
			http.Error(w, "Failed to broadcast update", http.StatusInternalServerError)
			return
		}

		sse := datastar.NewSSE(w, r)
		if err := sse.MarshalAndPatchSignals(map[string]string{
			fmt.Sprintf("successMessage-%d", billettholderID): fmt.Sprintf("Epostadressen %s er slettet", email),
			fmt.Sprintf("errorMessage-%d", billettholderID):   "",
		}); err != nil {
			logger.Error(fmt.Errorf("failed to patch billettholder signals: %w", err).Error(), "billettholder_id", billettholderID)
		}
	})
}
```

`billettholderLinkedToUser` and `deleteManualEmail` stand for the service code you write. The access check reads `relation_billettholdere_users` joined with `users.external_id = user.Id`, never an e-mail match (AGENTS.md, "Billettholder access"). The real e-mail routes do not check this yet. When the requester needs no signals at all, replace the SSE block with `w.WriteHeader(http.StatusNoContent)`.

The broadcast reaches every page that streams `billettholders`; the live map lists them. A page that renders billettholdere without streaming that bucket, such as `/event/{id}` with its billettholder picker, stays stale. Name such a gap in your notes, or add the bucket to its `Stream` call.

The matching trigger in the card. The handler rejects with `http.Error`, and the client drops that text, so the trigger needs a feedback key, or a rejected delete shows nothing:

```templ
<button
	type="button"
	class="btn btn--ghost"
	data-on:click={ fmt.Sprintf("@post('/profile/tickets/api/delete-email/%d/%d/')", billettholder.ID, email.ID) }
	{ errorfeedback.Attrs(fmt.Sprintf("deleteEmail%d", email.ID), "Klarte ikke å slette epostadressen. Prøv igjen.")... }
>
	Slett
</button>
@errorfeedback.RenderMessage(fmt.Sprintf("deleteEmail%d", email.ID))
<strong
	data-signals={ components.DataSignals(map[string]string{fmt.Sprintf("successMessage-%d", billettholder.ID): ""}) }
	data-text={ fmt.Sprintf("$successMessage-%d", billettholder.ID) }
></strong>
```

The client reads `$successMessage-42` as one signal name, which is why the card can put the id in the name. Each card gets its own signal, so the cards do not overwrite each other.

## 6. Route test with embedded NATS

From `pages/admin/rooms_assignment_route_test.go` and the helper in `components/formsubmission/other_details_test.go`. Use it when the handler treats the broadcast as required. The helper returns the embedded server as well, for section 9.

```go
func newTestLiveManager(t *testing.T) (*live.Manager, *embeddednats.Server) {
	t.Helper()
	ns, err := embeddednats.New(context.Background(), embeddednats.WithNATSServerOptions(&natsserver.Options{
		Host: "127.0.0.1", Port: -1, JetStream: true, StoreDir: t.TempDir(),
	}))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ns.Close() })
	ns.WaitForServer()
	manager, err := live.NewManager(context.Background(), ns, sessions.NewCookieStore([]byte("profile-tickets-test")))
	if err != nil {
		t.Fatal(err)
	}
	return manager, ns
}

func TestDeleteEmailRoute_RemovesManualEmail(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en billettholder knyttet til brukeren, med en manuell e-postadresse.",
		When:  "Når brukeren sletter adressen.",
		Then:  "Så er raden borte og svaret bekrefter slettingen.",
	})

	// Given
	expectedStatus := http.StatusOK
	db, logger := testutil.CreateTestDBAndLogger(t, "profile_tickets_delete_email")
	insertProfileRouteTestBillettholder(t, db, 42)
	insertProfileRouteTestUser(t, db, 7, "kari.nordmann@example.com")
	insertProfileRouteTestBillettholderUserAssociation(t, db, 42, 7)
	emailID := insertProfileRouteTestEmail(t, db, 42, "kari.manuell@example.com", models.BillettholderEmailKindManual)
	manager, _ := newTestLiveManager(t)
	router := chi.NewRouter()
	deleteEmailFromBillettholderRoute(router, db, logger, manager)

	// When
	request := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/delete-email/42/%d/", emailID), nil)
	request = request.WithContext(authctx.WithUserToken(request.Context(), "profile-route-user-7", "kari.nordmann@example.com"))
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != expectedStatus {
		t.Fatalf("expected status %d, got %d: %s", expectedStatus, recorder.Code, recorder.Body.String())
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_billettholder_emails WHERE id = ?`, emailID); got != 0 {
		t.Fatalf("expected the e-mail row to be deleted, found %d", got)
	}
}
```

Imports: `github.com/delaneyj/toolbelt/embeddednats`, `github.com/gorilla/sessions`, `natsserver "github.com/nats-io/nats-server/v2/server"`, `github.com/Regncon/conorganizer/service/authctx`.

`insertProfileRouteTestUser` gives user 7 the external id `profile-route-user-7`, and `userctx.GetUserRequestInfo(ctx).Id` is that external id, so the token must carry it. `authctx.WithUserToken` and `authctx.WithAdminUserToken` live in `service/authctx/test_support.go`.

This test proves the write and that the required broadcast does not fail. It still passes without the `Broadcast` line, because the status and the row do not depend on it. Section 9 proves an open page re-renders.

`&live.Manager{}` is enough only when the handler treats the broadcast as best effort, the test stops before the broadcast, or the test checks the broadcast-failure path on purpose; its `Broadcast` returns `unknown live bucket: ...` because no buckets were created.

## 7. Field errors the user should read

Section 5 rejects with `http.Error`, so the user only ever sees the trigger's fallback text. When the user can fix the problem and needs to know which one, answer 200 with an `errorfeedback` patch instead. The order stays the same as in section 5, with one branch added: validate before the write, and open the SSE only in the branch that patches. `http.Error` then stays available for everything before that, including a failed broadcast. This is the profile tickets add-e-mail handler in that shape:

```go
func addEmailToBillettholderRoute(router chi.Router, db *sql.DB, logger *slog.Logger, liveManager *live.Manager) {
	logger = logger.With("component", "profile_tickets")
	router.Post("/new-email/{id}/", func(w http.ResponseWriter, r *http.Request) {
		billettholderID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, "Ugyldig billettholder", http.StatusBadRequest)
			return
		}
		signals := map[string]any{}
		if err := datastar.ReadSignals(r, &signals); err != nil {
			http.Error(w, "Klarte ikke å lese skjemadata", http.StatusBadRequest)
			return
		}
		inputSignal := fmt.Sprintf("newEmail-%d", billettholderID)
		email, _ := signals[inputSignal].(string)
		email = strings.TrimSpace(email)
		user := userctx.GetUserRequestInfo(r.Context())

		linked, err := billettholderLinkedToUser(db, billettholderID, user.Id)
		if err != nil {
			logger.Error(fmt.Errorf("failed to check billettholder access: %w", err).Error(), "billettholder_id", billettholderID)
			http.Error(w, "Noe gikk galt, prøv igjen", http.StatusInternalServerError)
			return
		}
		if !linked {
			http.Error(w, "Du har ikke tilgang", http.StatusForbidden)
			return
		}

		feedbackKey := fmt.Sprintf("newEmail%d", billettholderID)
		feedback := errorfeedback.New(feedbackKey)
		if email == "" {
			feedback.Set(feedbackKey, "Skriv inn en e-postadresse.")
		} else {
			onBillettholder, err := billettholderHasEmail(db, billettholderID, email)
			if err != nil {
				logger.Error(fmt.Errorf("failed to check existing billettholder e-mails: %w", err).Error(), "billettholder_id", billettholderID)
				http.Error(w, "Noe gikk galt, prøv igjen", http.StatusInternalServerError)
				return
			}
			if onBillettholder {
				feedback.Set(feedbackKey, fmt.Sprintf("%s står allerede på denne billetten.", email))
			}
		}
		if feedback.HasErrors() {
			sse := datastar.NewSSE(w, r)
			if err := feedback.Patch(sse); err != nil {
				logger.Error(fmt.Errorf("failed to patch e-mail feedback: %w", err).Error(), "billettholder_id", billettholderID)
			}
			return
		}

		if err := addManualEmail(db, billettholderID, email); err != nil {
			logger.Error(fmt.Errorf("failed to add billettholder e-mail: %w", err).Error(), "billettholder_id", billettholderID)
			http.Error(w, "Noe gikk galt, prøv igjen", http.StatusInternalServerError)
			return
		}
		if err := liveManager.Broadcast(r.Context(), live.BucketBillettholders); err != nil {
			logger.Error(fmt.Errorf("failed to broadcast billettholder e-mail add: %w", err).Error(), "billettholder_id", billettholderID)
			http.Error(w, "Failed to broadcast update", http.StatusInternalServerError)
			return
		}

		sse := datastar.NewSSE(w, r)
		response := feedback.Signals()
		response[inputSignal] = ""
		if err := sse.MarshalAndPatchSignals(response); err != nil {
			logger.Error(fmt.Errorf("failed to patch e-mail signals: %w", err).Error(), "billettholder_id", billettholderID)
		}
	})
}
```

Imports beyond section 5: `strings` and `github.com/Regncon/conorganizer/components/errorfeedback`.

- Access is checked before validation. A field error sent before the access check would answer a user without access with 200, and a check that reads the database would tell them what is stored on a billettholder they may not see.
- Read the signals into `map[string]any` or a struct. The browser sends every signal that does not start with `_`, and `RenderMessage` declares `feedbackErrors` as an object, so reading into `map[string]string`, as the real route does, makes `ReadSignals` return an error on every click.
- `billettholderHasEmail` and `addManualEmail` stand for service code, like `billettholderLinkedToUser` in section 5. Check for duplicates in SQL with `email = ?`: the column is `COLLATE NOCASE`, so `Kari.Manuell@example.com` counts as a duplicate of `kari.manuell@example.com`, as in the real route. `addManualEmail` also links users through `checkIn.AssociateUsersWithBillettholderEmail`, as the real route does; adding a manual e-mail is one of the flows AGENTS.md allows to create links.
- `errorfeedback.New(feedbackKey)` registers the one key this response owns. On the error path `Patch` sends the message. On success `feedback.Signals()` sends the same key as `""`, and the extra `inputSignal` entry empties the input; both go in one `MarshalAndPatchSignals` call. `error_feedback.js` already clears the trigger's own key when a request starts, so the server's `""` matters for keys other triggers set and for a trigger without `errorfeedback.Attrs`.
- Nothing is written on the error path, so there is nothing to broadcast. On success the new e-mail reaches the card through the broadcast, on every open tab, not through this response.
- The feedback key (`newEmail42`) and the input signal (`newEmail-42`) are different things: the key lives inside `$feedbackErrors`, the input value is its own per-card signal.
- Validation messages are expected outcomes, so the handler does not log them; the conorganizer-logging skill has the rule.
- Moving the real route to this shape changes what users see. The real card confirms with "Epostadressen … er lagt til" through `successMessage-{id}` and uses other wording for its errors, and `documentation/testing/profile-tickets.md`, `documentation/billettholdere.md` and `billettholder_email_routes_test.go` quote those texts. To keep the confirmation, add `successMessage-{id}` to the response map on success and set it to `""` on the error path. If the texts change, update those files and add negative `!strings.Contains` checks, as AGENTS.md asks. The real function is spelled `addEmailToBilettholderRoute`; use one name in the route setup, the handler and the tests.

The matching markup in the card. Enter in the input and the button both start the request, so both carry the same feedback key:

```templ
<input
	class="input"
	type="email"
	data-bind={ fmt.Sprintf("newEmail-%d", billettholder.ID) }
	data-on:keydown={ fmt.Sprintf("if (evt.key === 'Enter') { evt.preventDefault(); @post('/profile/tickets/api/new-email/%d/') }", billettholder.ID) }
	{ errorfeedback.Attrs(fmt.Sprintf("newEmail%d", billettholder.ID), "Klarte ikke å legge til e-postadressen. Prøv igjen.")... }
/>
<button
	type="button"
	class="btn btn--secondary"
	data-on:click={ fmt.Sprintf("@post('/profile/tickets/api/new-email/%d/')", billettholder.ID) }
	{ errorfeedback.Attrs(fmt.Sprintf("newEmail%d", billettholder.ID), "Klarte ikke å legge til e-postadressen. Prøv igjen.")... }
>
	Legg til
</button>
@errorfeedback.RenderMessage(fmt.Sprintf("newEmail%d", billettholder.ID))
```

The fallback text in `Attrs` covers the `http.Error` answers and a failed request; the handler's own messages cover what the user can fix.

## 8. Testing field errors

The feedback arrives in the `data: signals` lines of the SSE body. `feedbackErrors` is nested, so read each payload as `map[string]any`; the flat `profilePatchedSignals` helper in `billettholder_email_routes_test.go` cannot read it. `mergeSignalPatch` merges nested objects the way the browser does, so two patches that each carry `feedbackErrors` keep both keys. Send what the browser sends: every signal that does not start with `_`, including the `feedbackErrors` object. Test both paths: the error path shows the message and writes nothing, and the success path patches every key the response owns to `""`.

```go
func patchedSignals(t *testing.T, body string) map[string]any {
	t.Helper()
	signals := map[string]any{}
	for line := range strings.SplitSeq(body, "\n") {
		payload, ok := strings.CutPrefix(line, "data: signals ")
		if !ok {
			continue
		}
		patch := map[string]any{}
		if err := json.Unmarshal([]byte(payload), &patch); err != nil {
			t.Fatalf("failed to unmarshal Datastar signal patch %q: %v", payload, err)
		}
		mergeSignalPatch(signals, patch)
	}
	return signals
}

func mergeSignalPatch(signals map[string]any, patch map[string]any) {
	for key, value := range patch {
		nestedPatch, patchIsObject := value.(map[string]any)
		nested, signalIsObject := signals[key].(map[string]any)
		if patchIsObject && signalIsObject {
			mergeSignalPatch(nested, nestedPatch)
			continue
		}
		signals[key] = value
	}
}

func patchedFeedbackError(t *testing.T, body string, key string) (string, bool) {
	t.Helper()
	feedbackErrors, _ := patchedSignals(t, body)["feedbackErrors"].(map[string]any)
	message, ok := feedbackErrors[key].(string)
	return message, ok
}

func postSignalsAsKari(t *testing.T, router http.Handler, path string, signalsJSON string) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(http.MethodPost, path, strings.NewReader(signalsJSON))
	request.Header.Set("Content-Type", "application/json")
	request = request.WithContext(authctx.WithUserToken(request.Context(), "profile-route-user-7", "kari.nordmann@example.com"))
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	return recorder
}

func TestAddEmailRoute_WhenEmailIsAlreadyOnBillettholder_ShowsFieldError(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en billettholder knyttet til brukeren, der e-postadressen allerede står.",
		When:  "Når brukeren legger til den samme adressen igjen.",
		Then:  "Så forteller feltet hvorfor, og ingenting lagres.",
	})

	// Given
	expectedMessage := "Kari.Manuell@example.com står allerede på denne billetten."
	db, logger := testutil.CreateTestDBAndLogger(t, "profile_tickets_add_email_duplicate")
	insertProfileRouteTestBillettholder(t, db, 42)
	insertProfileRouteTestUser(t, db, 7, "kari.nordmann@example.com")
	insertProfileRouteTestBillettholderUserAssociation(t, db, 42, 7)
	insertProfileRouteTestEmail(t, db, 42, "kari.manuell@example.com", models.BillettholderEmailKindManual)
	manager, _ := newTestLiveManager(t)
	router := chi.NewRouter()
	addEmailToBillettholderRoute(router, db, logger, manager)

	// When
	recorder := postSignalsAsKari(t, router, "/new-email/42/", `{"newEmail-42": "Kari.Manuell@example.com", "feedbackErrors": {"newEmail42": ""}}`)

	// Then
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}
	if message, _ := patchedFeedbackError(t, recorder.Body.String(), "newEmail42"); message != expectedMessage {
		t.Fatalf("field error mismatch\nexpected: %q\nactual:   %q", expectedMessage, message)
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_billettholder_emails WHERE billettholder_id = 42`); got != 1 {
		t.Fatalf("expected no new e-mail row, found %d rows", got)
	}
}

func TestAddEmailRoute_WhenEmailIsNew_ClearsFieldErrorAndInput(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en billettholder knyttet til brukeren.",
		When:  "Når brukeren legger til en ny e-postadresse.",
		Then:  "Så lagres adressen, feltfeilen tømmes og feltet blir tomt.",
	})

	// Given
	expectedEmail := "kari.ny@example.com"
	db, logger := testutil.CreateTestDBAndLogger(t, "profile_tickets_add_email_new")
	insertProfileRouteTestBillettholder(t, db, 42)
	insertProfileRouteTestUser(t, db, 7, "kari.nordmann@example.com")
	insertProfileRouteTestBillettholderUserAssociation(t, db, 42, 7)
	manager, _ := newTestLiveManager(t)
	router := chi.NewRouter()
	addEmailToBillettholderRoute(router, db, logger, manager)

	// When
	recorder := postSignalsAsKari(t, router, "/new-email/42/", `{"newEmail-42": "kari.ny@example.com", "feedbackErrors": {"newEmail42": ""}}`)

	// Then
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}
	if message, patched := patchedFeedbackError(t, recorder.Body.String(), "newEmail42"); !patched || message != "" {
		t.Fatalf("expected the field error to be cleared, got %q (patched: %v)", message, patched)
	}
	if input := patchedSignals(t, recorder.Body.String())["newEmail-42"]; input != "" {
		t.Fatalf("expected the input to be emptied, got %v", input)
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_billettholder_emails WHERE billettholder_id = 42 AND email = ?`, expectedEmail); got != 1 {
		t.Fatalf("expected the new e-mail to be saved once, found %d rows", got)
	}
}
```

Imports beyond section 6: `encoding/json` and `strings`.

A recorder is enough for these tests: the handler reads its signals before any SSE is opened, so the body is intact. A handler that reads after `NewSSE` would still pass a recorder test and fail on a real server, because Go's HTTP/1 server discards the unread body once the headers are flushed. To catch that, post real signals to the route over the section 9 server, read the whole response, and check the row:

```go
response, err := server.Client().Post(server.URL+"/profile/tickets/api/new-email/42/", "application/json", strings.NewReader(`{"newEmail-42": "kari.ny@example.com", "feedbackErrors": {"newEmail42": ""}}`))
if err != nil {
	t.Fatal(err)
}
defer response.Body.Close()
if _, err := io.ReadAll(response.Body); err != nil {
	t.Fatal(err)
}
if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_billettholder_emails WHERE billettholder_id = 42 AND email = ?`, "kari.ny@example.com"); got != 1 {
	t.Fatalf("expected the e-mail to be saved over a real server, found %d rows", got)
}
```

## 9. End-to-end: an open page re-renders after a mutation

Open the stream for real, wait until it watches, mutate, and read the next patch. `httptest.NewRecorder` cannot do this, because `Stream` blocks until the request ends; use `httptest.NewServer`.

```go
const liveTestTimeout = 5 * time.Second

func newProfileTicketsLiveServer(t *testing.T, db *sql.DB, logger *slog.Logger, manager *live.Manager, externalID string) *httptest.Server {
	t.Helper()
	router := chi.NewRouter()
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(authctx.WithUserToken(r.Context(), externalID, "kari.nordmann@example.com")))
		})
	})
	router.Route("/profile", func(profileRouter chi.Router) {
		if err := ProfileTicketsRoute(profileRouter, manager, db, logger); err != nil {
			t.Fatal(err)
		}
	})
	server := httptest.NewServer(router)
	t.Cleanup(server.Close)
	return server
}

func openLiveStream(t *testing.T, server *httptest.Server, path string) <-chan string {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL+path, nil)
	if err != nil {
		t.Fatal(err)
	}
	response, err := server.Client().Do(request)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { response.Body.Close() })
	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected the live stream to open with 200, got %d", response.StatusCode)
	}
	patches := make(chan string)
	go func() {
		defer close(patches)
		scanner := bufio.NewScanner(response.Body)
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		var elements []string
		for scanner.Scan() {
			line := scanner.Text()
			if payload, ok := strings.CutPrefix(line, "data: elements "); ok {
				elements = append(elements, payload)
				continue
			}
			if line != "" || len(elements) == 0 {
				continue
			}
			select {
			case patches <- strings.Join(elements, "\n"):
			case <-ctx.Done():
				return
			}
			elements = nil
		}
	}()
	return patches
}

func nextLivePatch(t *testing.T, patches <-chan string) string {
	t.Helper()
	select {
	case patch, ok := <-patches:
		if !ok {
			t.Fatal("expected another live patch, but the stream closed")
		}
		return patch
	case <-time.After(liveTestTimeout):
		t.Fatal("timed out waiting for a live patch")
	}
	return ""
}

func waitUntilWatching(t *testing.T, ns *embeddednats.Server, bucket live.Bucket) {
	t.Helper()
	nc, err := ns.Client()
	if err != nil {
		t.Fatal(err)
	}
	defer nc.Close()
	js, err := jetstream.New(nc)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), liveTestTimeout)
	defer cancel()
	stream, err := js.Stream(ctx, "KV_"+string(bucket))
	if err != nil {
		t.Fatal(err)
	}
	for {
		info, err := stream.Info(ctx)
		if err != nil {
			t.Fatalf("timed out waiting for a watcher on %s: %v", bucket, err)
		}
		if info.State.Consumers > 0 {
			return
		}
		select {
		case <-ctx.Done():
			t.Fatalf("timed out waiting for a watcher on %s", bucket)
		case <-time.After(10 * time.Millisecond):
		}
	}
}

func postDatastar(t *testing.T, server *httptest.Server, path string) {
	t.Helper()
	response, err := server.Client().Post(server.URL+path, "application/json", strings.NewReader("{}"))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode/100 != 2 {
		t.Fatalf("expected %s to succeed, got %d", path, response.StatusCode)
	}
}

func TestProfileTicketsLive_DeletedEmailDisappearsFromOpenPage(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at billettsiden er åpen med en manuell e-postadresse på kortet.",
		When:  "Når adressen slettes.",
		Then:  "Så får den åpne siden en ny patch uten adressen.",
	})

	// Given
	expectedEmail := "kari.manuell@example.com"
	db, logger := testutil.CreateTestDBAndLogger(t, "profile_tickets_live")
	insertProfileRouteTestBillettholder(t, db, 42)
	insertProfileRouteTestUser(t, db, 7, "kari.nordmann@example.com")
	insertProfileRouteTestBillettholderUserAssociation(t, db, 42, 7)
	emailID := insertProfileRouteTestEmail(t, db, 42, expectedEmail, models.BillettholderEmailKindManual)
	manager, ns := newTestLiveManager(t)
	server := newProfileTicketsLiveServer(t, db, logger, manager, "profile-route-user-7")
	patches := openLiveStream(t, server, "/profile/tickets/api")
	if first := nextLivePatch(t, patches); !strings.Contains(first, expectedEmail) {
		t.Fatalf("expected the first patch to show %s", expectedEmail)
	}
	waitUntilWatching(t, ns, live.BucketBillettholders)

	// When
	postDatastar(t, server, fmt.Sprintf("/profile/tickets/api/delete-email/42/%d/", emailID))

	// Then
	if patch := nextLivePatch(t, patches); strings.Contains(patch, expectedEmail) {
		t.Fatalf("expected the open page to re-render without %s", expectedEmail)
	}
}
```

Extra imports: `bufio`, `database/sql`, `log/slog`, `strings`, `time`, `github.com/nats-io/nats.go/jetstream`.

Over `httptest.NewServer` the test's own request context never reaches the handler, so the user goes in through the router middleware. Without a token the renderer gets an empty user id, and `GetBillettholdere("")` then returns every billettholder, which hides access bugs. `waitUntilWatching` counts consumers on the bucket, so open only one stream per bucket in a test, or wait for the right number. Delete the `Broadcast` line once by hand and watch this test fail; that is the proof the handler broadcasts.

## 10. Doc rows to add

In `documentation/live-update-lifecycle.md`:

- "Page subscriptions": `| /profile/tickets | /profile/tickets/api | billettholders | |`
- "Buckets", column "Typical broadcasters" of the bucket you broadcast: add the new mutation.
- "Cross-bucket broadcasts": only when the handler broadcasts more than one bucket.
