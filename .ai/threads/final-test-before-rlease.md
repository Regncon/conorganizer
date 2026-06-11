--- Codex Session ---
ID: 019eb602-a829-78f1-b582-b26b22179d22

# Final test before release

We are nearing the release and we need to test everything to make sure it works as expected.

You can see some of our previous discussions about testing in the following threads:
.ai/threads/create-test-ckecklist.md
.ai/threads/add-more-automated-tests.md

Have a look at the manual test in documentation/testing/

The code has changed since we started writing the manual tests, so we should verify that all the tests are still valid and update them if necessary.

We should add a checklist to documentation/testing/index.md with a check mark for each test file.

Also, when we add any automated tests, we should follow the new Go test structure guidelines as outlined in AGENTS.MD and copied below.
Also find some good examples in the existing tests that use this new structure.

## Go test structure
When writing or cleaning up Go tests, prefer behavior-focused tests that read from intent to verification.

- Add a short BDD-style comment at the top of each test when it improves readability:
  - Given some initial context,
  - When an event occurs,
  - Then ensure the expected outcome.
- Split each test into three visible sections:
  - `// Given`
  - `// When`
  - `// Then`
- Start the `// Given` section with the expected output or expected result. This makes the test's purpose clear before setup details.
- Keep test names specific to one behavior. Prefer splitting creation, duplicate/idempotency, validation, and error cases into separate tests.
- Move noisy setup SQL and repeated assertions into local test helpers when that makes the test read like behavior.

Example shape:

```go
func TestService_EventProducesExpectedResult(t *testing.T) {
	// Given an existing domain object,
	// when the service reconciles related data,
	// then the expected durable relation is stored.

	// Given
	expectedResult := models.SomeResult{ID: 123}
	db := createTestDB(t)
	insertRequiredFixture(t, db, expectedResult.ID)

	// When
	err := ServiceFunction(expectedResult.ID, db)

	// Then
	if err != nil {
		t.Fatalf("expected service call to succeed: %v", err)
	}
	assertExpectedResult(t, db, expectedResult)
}

```

Let's make a plan for how we will approach this testing. 
I think first make sure that all pages are covered.
Then go through each test and verify that it is still valid and update it if necessary.
I'm interested in automating as much as possible, so as we go through the tests, we should identify which ones can be automated and add them to our automated test suite.
I don't trust our previous list of automation candidates, so we should review each one and make sure it still makes sense to automate it.

Let's have a discussion about the design before we start coding. If you have any clarification questions or suggestions, please let me know.

Time: 2026-06-11 11:32:17
--- Codex Run ---

Codex output:
I’ll read the existing testing notes and project guidance first, then I’ll turn that into a proposed test strategy without changing files yet.
I found the testing docs split by page/flow, plus the two prior thread notes. Next I’m reading the actual checklist content and current AGENTS guidance so the plan is based on what is in the repo now, not the earlier assumptions.
Done reading AGENTS.MD file.
I’ve got the prior decisions: bokmål, BDD-style manual checklist, page/flow files, and Go tests that start from expected outcomes. I’m now checking the current checklist contents against the current route/component surface so we can decide how to audit and automate without duplicating work.
The current docs already show a few drift signals: `index.md` links “Mine arrangementer” to `profile.md`, `general.md` still has some broad/manual nav checks that may overlap with automated header tests, and some checklist points are phrased against flows that may depend on the auth provider. I’m reading the remaining page-specific files before proposing the audit order.
There is already a meaningful automated layer beyond the first menu tests: root program rendering, event visibility/navigation/interest behavior, profile program behavior, admin billettholder counts, live/session checks, and service-level rules. I’m doing one focused route scan now so the plan can include “pages covered” as an explicit first pass.
The route scan confirms at least one current gap in the manual package: admin rooms routes exist under `/admin/rooms`, but there is no corresponding manual checklist file yet. That should be handled in the first coverage pass, not discovered late while testing.
I also see room service tests already exist, but they look older and less aligned with the new BDD structure than the newer root/event tests. For this release pass I’d treat rooms as both a documentation gap and a candidate for a small amount of higher-level UI/component coverage, while leaving deeper room service cleanup as a separate decision.
The newer tests give us the examples to copy: explicit behavior names, BDD comment, `// Given` starts with expected values, then setup, render/action, and assertions. I would not use the older table-heavy/service tests as the style model unless we intentionally clean them up.
I’m going to run the current automated suite once as a baseline. That gives us a concrete starting point for the plan: either “tests pass before the release audit” or “stabilize existing failures first.”
Baseline is green: `go test ./...` passes locally with `GOCACHE=/tmp/go-build-cache`. I’m doing a current `.templ` logic scan next so our automation candidates are based on today’s code, not the older list.
One more concrete drift item: `my-events.md` claims to cover `/profile`, which is the same page as `profile.md`. That may be acceptable if it is a flow section rather than a separate route, but the index link is currently wrong and we should decide whether to keep it as a separate flow file or fold it into `profile.md`.
The current “Send inn arrangement” entry is worth correcting during the audit: the front-page CTA links to `/profile`, while the actual create action on Min Side posts to `/profile/api/create` and redirects to `/profile/new/{id}`. The manual test should describe that two-step behavior instead of implying the front-page button opens the form directly.
So far I’d structure the work as a release audit rather than “write all missing tests blindly”: first route/file coverage, then validity corrections, then automation. The most valuable automation candidates appear to be deterministic server-rendered behavior and route/service rules; browser E2E should stay limited to true interaction/layout cases.
A key design choice: I would not remove manual tests merely because component tests exist unless the automated test covers the same user-visible contract. Example: header href sets are automated, but keyboard focus, mobile menu usability, and “clear path after access denied” still belong in manual testing unless we add browser tests.
Admin access is another place to verify carefully. Some admin API routes use explicit `RequireAdmin`, but the page-level `/admin` route is mounted under the app router with auth middleware and may rely on layout/page checks elsewhere; part of the audit should confirm whether non-admin access is actually blocked at every admin page, not just hidden in navigation.
There is a concrete risk to investigate when we move from planning to fixes: `/admin/`, `/admin/api/`, and `/admin/approval/` render without an explicit `RequireAdmin` wrapper in the route definitions I just checked, while some billettholder APIs do have it. I would put admin route protection near the top of the release test plan.
I agree with your proposed order, with one addition: first build a route-to-test-file coverage map from the current code, not from the existing docs.

Baseline check: `GOCACHE=/tmp/go-build-cache go test ./...` passes right now.

**Proposed Plan**
1. **Coverage inventory**
   Compare active routes/pages/components against `documentation/testing/`. I already see likely gaps/drift:
   `admin-rooms.md` is missing for `/admin/rooms`.
   `index.md` links “Mine arrangementer” to `profile.md` instead of `my-events.md`.
   `my-events.md` covers a section of `/profile`, not a separate page.
   `event-form.md` should reflect that the front-page CTA goes to `/profile`, then the form is created from Min Side.

2. **Update `documentation/testing/index.md`**
   I suggest turning the overview into a checklist like:
   `- [ ] [Generelle tester](./general.md)`
   One checkbox per manual test file. I’d leave them unchecked so they can be used for the release run.

3. **Validate each manual test file**
   For each file:
   verify route/page still exists,
   verify behavior still matches code,
   remove or rewrite stale points,
   add missing launch-critical behavior,
   keep manual checks for layout, accessibility-ish keyboard use, visual polish, and third-party auth flows.

4. **Automation pass per file**
   For each checklist item, decide:
   already automated,
   automate now with Go component/service/route test,
   keep manual,
   or defer as browser E2E only if interaction/layout truly needs a browser.

   I would add small automated tests as we go when the behavior is deterministic and cheap to test.

5. **High-priority automation candidates**
   First candidates I’d prioritize:
   admin route protection, especially `/admin`, `/admin/api`, `/admin/approval`, `/admin/rooms`;
   missing rooms page behavior;
   `EventCard` href/status/fallback rendering;
   `MyEvents` status-based links and empty title fallback;
   front-page CTA href behavior;
   admin program publishing and pulje status route behavior.

6. **Verification**
   After each batch: `gofmt`, focused `go test`, then `go test ./...`.
   If `.templ` files change, also run templ generation before tests.

**Clarifications**
Do you want `my-events.md` to stay as a separate flow checklist for the “Mine arrangementer” section on `/profile`, or should we merge it into `profile.md`?

I also recommend adding `documentation/testing/admin-rooms.md` since rooms are now exposed from the admin page.

Time: 2026-06-11 11:32:17
--- Next Prompt ---
Let's merge `my-events.md` into `profile.md` 
Yes, let's add `admin-rooms.md` 

Go ahead and start with the coverage inventory, and update `index.md` with the checklist. Then we can go through each file and verify the tests one at a time. Don't start the next file until we have verified the previous one.

Let's start each file with a discussion about the file before we start coding. If you have any clarification questions or suggestions, please let me know.

Time: 2026-06-11 11:48:55
--- Codex Run ---

Codex output:
I’ll make the inventory/index pass only now: update the file list, remove the standalone “Mine arrangementer” entry, add admin rooms, and leave deeper validation of each checklist for the one-file-at-a-time pass. For the merge, I’ll keep it mechanical and avoid rewriting profile behavior until we discuss `profile.md` specifically.
The active route inventory confirms `/print` still exists, but prior launch scope excludes it, so I’ll list it as deliberately excluded rather than adding it back. `/my-profile*` does not appear in the active route setup.
I’m going to edit three docs now: `index.md` for the inventory/checklist, `profile.md` for the mechanical merge of Mine arrangementer, and a new `admin-rooms.md` skeleton so the coverage list has a real target. I’ll delete the old standalone `my-events.md` after the merge to avoid duplicate launch files.
The inventory/index edits are in place. I’m checking the resulting docs for broken references and making sure the deleted `my-events.md` no longer appears in the testing file list.
The merge left `profile.md` with a clean section for Mine arrangementer, and `index.md` now has one checkbox per launch test file. I’m checking the diff and repo status before I report back and ask which file we should discuss first.
There is an unrelated untracked thread file under `.ai/threads/` that I did not touch. I’ll leave it alone and just run a whitespace check on the documentation edits.
Updated the coverage inventory and index checklist.

Changes:
- [index.md](/home/cinmay/Documents/conorganizer/documentation/testing/index.md) now has one checkbox per launch test file and a route coverage inventory.
- [profile.md](/home/cinmay/Documents/conorganizer/documentation/testing/profile.md) now includes the merged “Mine arrangementer” section.
- [my-events.md](/home/cinmay/Documents/conorganizer/documentation/testing/my-events.md) was deleted to avoid duplicate launch files.
- [admin-rooms.md](/home/cinmay/Documents/conorganizer/documentation/testing/admin-rooms.md) was added as a skeleton only; we’ll fill it when we get to that file.

Coverage notes added:
- `/print` is explicitly excluded from launch testing.
- `/auth/test` is excluded as diagnostic.
- API/live routes are covered through the owning page or flow.

Verification:
- `go test ./...` passed before these doc edits.
- `git diff --check` passed for the documentation changes.

Suggested next step: start with `general.md`. Before editing it, I’d like to discuss whether we should keep only cross-cutting manual behavior there: navigation usability, protected-page behavior, mobile/keyboard checks, logout/access-state consistency, and remove anything already fully covered by the header component tests.

Time: 2026-06-11 11:48:55
--- Next Prompt ---
Yes, let's remove anything that is already fully covered by the header component tests.

Time: 2026-06-11 12:00:32
--- Codex Run ---

Codex output:
For `general.md`, the header component tests fully cover the rendered menu href sets by role: anonymous, logged-in non-admin, and admin. They do not cover clicking, access-denied pages, logout effects, keyboard focus, mobile layout, or whether external links are understandable, so I’ll keep those manual checks.

The only item I see as close enough to remove is the broad admin point about “all links work as expected” from the main nav, because it overlaps with menu-link coverage and is too vague. I’ll also move the external-link check out of the “Ikke innlogget bruker” subsection, since those external FAQ links only render for logged-in users/admins.
I’m updating only `general.md` now, keeping the change scoped to the header-test overlap and the misplaced external-link check.
The `general.md` cleanup is done. I’m checking the file and diff for the exact movement/removal before summarizing.
Updated [general.md](/home/cinmay/Documents/conorganizer/documentation/testing/general.md).

I removed the vague admin navigation item that overlapped with the header component href tests, and moved the external-link check from “ikke innlogget bruker” to “innlogget bruker”, since the FAQ link is only rendered for logged-in users/admins.

I kept the remaining items because the header tests do not cover clicking, logout effects, access-denied behavior, keyboard use, mobile layout, or visual clarity.

Validation: `git diff --check` passed for `general.md`.

My suggestion: treat `general.md` as verified unless you want to tighten wording further. Next file in the index is `root.md`; before editing it, we should discuss which current root tests are already automated and what should stay manual.

Time: 2026-06-11 12:00:32
--- Next Prompt ---
Before we move on, see if you can clean up the text in the `general.md` file to make it more concise and easier to read. 
Also, double check that the test coverage is complete for all behaviors. 
Then, double check if there is anything that can be automated. 

Time: 2026-06-11 12:06:50
--- Codex Run ---

Codex output:
I’ll re-audit `general.md` against the current header/auth route behavior and existing tests, then tighten the checklist text and add the missing automation section if needed.
The current `general.md` is missing a `Kan automatiseres` section, and the remaining checklist is complete for the cross-cutting behavior I can see: login entry, protected-page handling, logout, admin access, responsive/menu usability, keyboard/focus, and external links. I’m going to tighten the wording and add automation notes without changing the intended coverage.
I’ll replace the checklist with a shorter structure: shared navigation checks once, then role-specific checks. That reduces repetition while keeping the same coverage surface.
The text is tightened and the automation section is now present. I’m doing one pass over the final file and checking the diff for formatting issues.
The checklist now has a complete cross-cutting coverage shape and explicitly calls out what is already automated versus still manual. I’m running the header component tests once because those are the automated coverage this file now relies on.
Header tests pass. I’m checking the working diff/status now so the final note is precise.
Cleaned up [general.md](/home/cinmay/Documents/conorganizer/documentation/testing/general.md).

I rewrote it into a tighter structure:
- `Alle roller`
- `Ikke innlogget bruker`
- `Innlogget bruker`
- `Admin`
- `Kan automatiseres`

I double-checked it against the current header tests, menu template, route setup, `UserMiddleware`, `RequireAdmin`, and the unauthorized page. Coverage now includes the cross-cutting behavior that belongs in `general.md`: navigation stability, responsive/menu usability, keyboard focus, protected-page access, logout state, external links, and admin access.

I added automation notes for:
- already-covered header role link sets
- protected-page route tests
- admin route protection tests
- logout cookie-clearing handler test
- what should remain manual/browser-based

Verification:
- `git diff --check -- documentation/testing/general.md` passed.
- `GOCACHE=/tmp/go-build-cache go test ./components/header` passed.

`general.md` is ready to treat as verified. Next up is `root.md`, and we should discuss its scope before editing.

Time: 2026-06-11 12:06:50
--- Next Prompt ---

