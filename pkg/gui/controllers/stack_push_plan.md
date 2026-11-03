# Stack push menu preferences

## Status and scope

Implementation complete and locally validated on 2026-11-03. Steps 1–3 are
complete. External CI has not been run for these changes. The three requested
capabilities are accepted: disable the stack push feature, add direct menu
shortcuts, and configure the initially selected choice. Exact configuration
names and behavior below are adopted from the proposed plan.
One tracker owns this work: this file.

Baseline inspected: `master`, `0002c0354770ff8ec6abe4d5ab4d025ea40e1619`.
The checkout was clean before this plan. No separate reference implementation
or reference branch was supplied; the baseline implementation is the reference.
Implementation changes and validation are tracked below.

## Baseline evidence (before Steps 1–3)

- `sync_controller.go`: `HandlePush` reaches `push`, which unconditionally
  opens a two-choice menu when `unpushedBranchesBelow` returns any branches.
  The first choice pushes the current branch plus eligible branches below;
  the second calls `pushCurrentBranch`. Force-push checks and upstream
  resolution happen inside those paths, after the choice.
- `helpers/branch_stack.go`: stack membership comes from loaded unmerged
  commits, excluding configured main branches and commits reachable from local
  main branches. `unpushedBranchesBelow` additionally requires a push remote,
  a remote branch that has not disappeared, and commits ahead for push.
  Preserve this detection logic.
- `types/common.go`: `MenuItem.Keys` already supports direct execution and
  displays the first shortcut. Existing `menuKey` callers in controllers and
  `helpers/refs_helper.go` demonstrate owner-local shortcut declarations.
  Reuse these declarations and callbacks; no new keybinding framework is needed.
- `../menu_panel.go`: `createMenu` always calls `SetSelection(0)` and strips
  shortcuts conflicting with primary confirm/cancel/navigation bindings.
  `CreateMenuOptions` currently has no initial-selection field.
- `../../config/user_config.go` owns defaults and schema annotations;
  `user_config_validation.go` owns runtime enum validation. Both must agree.
- `../../integration/tests/sync/push_stacked_branches*.go` already exercises
  both push choices, force pushing, and missing upstreams. `MenuDriver.Confirm`
  can exercise Enter without navigation; view matchers can assert selection.
- `.github/workflows/ci.yml` runs short unit tests on Linux and Windows,
  integration tests through `scripts/run_integration_tests.sh` on Linux with
  several Git versions and a race job, generated-file freshness, and lint.
  These are inspected CI callers, not CI runs performed for this plan.

## Accepted behavior and decisions

Configuration defaults, preserving baseline behavior:

```yaml
git:
  # False bypasses this menu and pushes only the current branch.
  # Existing force-push and upstream prompts still apply.
  enableStackedPush: true
  # Allowed: allBranches, currentBranch. Used only when the menu opens.
  stackedPushDefault: allBranches
```

For the user's preferred initial selection:

```yaml
git:
  stackedPushDefault: currentBranch
```

- Disabling means push only the current branch, independently
  of the menu default. Do not silently push other branches when disabled.
  Adopted for implementation following the user's request to implement the plan.
- `a` directly executes the current-plus-below choice, and `o`
  directly executes only-current. Adopted for implementation. Shortcuts invoke the existing callbacks, including subsequent
  force confirmations; they do not merely highlight a row.
- Keep item order stable. Changing the default changes the highlighted row;
  Enter runs that row. Reordering rows would also change the visual layout and
  the relationship between choices and shortcuts, so prefer explicit selection.
- Add `InitialSelection int` to `CreateMenuOptions`, with zero retaining the
  existing behavior for every other menu. Apply it after installing menu items;
  clamp invalid indices to zero. The push controller chooses index 0 or 1 from
  the validated enum. Cancel remains available and is never a configured default.
- Preserve primary-key conflict handling. A conflicting menu shortcut is
  suppressed as in other menus; navigation and Enter remain available.
- Do not change the analogous pull menu, stack detection, Git push arguments,
  force-push protections, or translation files. Reuse current English labels.
- This is a feature change, so the bug-demonstration EXPECTED/ACTUAL pattern
  does not apply. No speculative refactor or shared push/pull policy is needed.

## Delivery steps

### Step 1: Allow opting out of stacked pushes

Problem: `push` always performs stack discovery and opens the menu when eligible
branches exist. There is no opt-out in configuration.

- [x] Adopt the proposed disabling behavior and configuration name.
- [x] Add `git.enableStackedPush`, default true, with an unwrapped field comment.
  When false, return `pushCurrentBranch` before discovering branches below.
- [x] Add integration coverage with an eligible rebased stack: disabling pushes
  only the current branch, retains force-push confirmation, and leaves lower
  remote refs unchanged. Preserve existing enabled and no-eligible-branch cases.
- [x] Run `just generate`; include generated config/schema and test registration
  changes. Never edit `docs/` or generated config documentation by hand.
- [x] Run the new focused e2e case and existing stacked push cases, then the
  implementation commit gates below; update this tracker and commit Step 1.

Exit: opt-out bypasses the menu while retaining current-branch push safeguards;
omitted configuration retains baseline behavior.

### Step 2: Add direct menu shortcuts

Depends on Step 1 for serial delivery, but changes a separate capability.

- [x] Adopt the proposed `a` / `o` shortcuts.
- [x] Add `menuKey` declarations to the two existing push menu items. Keep
  their labels, ordering, and callbacks; rely on the existing shortcut renderer.
- [x] Test both shortcuts through actual menu keypresses without navigation or
  Enter. Verify remote refs distinguish pushing only current from pushing all
  eligible branches, and required force confirmations remain in effect.
- [x] Cover primary-key conflict behavior using the existing menu coverage if
  available; otherwise add one focused regression through a customized binding.
  Keep cancellation and ordinary navigation covered by existing cases.
- [x] Run `just generate` for added integration tests and affected cheatsheets,
  inspect generated changes, run focused e2e cases and commit gates, update the
  tracker, and commit Step 2.

Exit: both actions are available via visible shortcuts and preserve the same
behavior as selecting and confirming their corresponding menu items.

### Step 3: Configure the initial selection

Depends on Steps 1 and 2. The shared API has a real consumer in this same step;
adding an optional selection is new behavior, not a preparatory pure refactor.

- [x] Add the zero-default `InitialSelection` option and consume it in
  `createMenu`, replacing the hard-coded selection with a bounded selection.
- [x] Add `git.stackedPushDefault` with `allBranches` and `currentBranch` values,
  default `allBranches`. Add schema enum tags and runtime `validateEnum` coverage
  for both accepted values and invalid values, following existing config tests.
- [x] Map the setting to initial selection in the push menu. Keep shortcut
  actions stable under either selection. When Step 1 disables the feature,
  the setting has no effect.
- [x] Test opening the menu and pressing Enter without any selection call for
  each default; assert selected row and resulting remote refs. Reopen after
  choosing a different row to prove selection resets to the configured default.
- [x] Cover disabled plus `allBranches`, disabled plus `currentBranch`, no
  eligible branches, and an unrelated menu retaining its first-row default.
  There is no existing isolated `createMenu` unit-test seam; actual menu
  integration tests cover selection and reset, with an explicit invalid-index
  fallback retained in `createMenu`.
- [x] Run `just generate`, focused config/menu tests and stacked-push e2e tests;
  run commit gates, update this tracker, and commit Step 3.

Exit: Enter follows the configured default on every opening, omitted settings
preserve baseline behavior, and other menus retain their existing defaults.

## Orchestration, validation, and handoff

- User override: implement directly without sub-agents. The primary agent owns
  implementation, checks, focused commits, and tracker updates, in order 1–3.
  This supersedes the original serial delegation workflow. No agents are active.
- No model or effort was specified by the user. Inherit the session defaults
  for implementation; use fresh context if requested. No review checkpoint is
  currently selected. If selected later, the skill prescribes a fresh-context
  `gpt-6.1-sol` reviewer at high effort; report unavailability, never substitute
  silently. Reviewers edit the tracker and identify corrections, not fixes.
- Keep scope, baseline/current revisions, dependencies, and open gates in this
  tracker. If a later review is requested, pass these and repository instructions
  to the reviewer. Required in-scope corrections get stable suffixes such as 2a,
  run before dependent work, and receive their own validation; optional advice
  stays separate. Do not renumber existing steps.
- Implementation commit gates: `just format`, `just build`, `just unit-test`,
  `just lint`, and relevant `just e2e sync/<test_name>` cases. Generate whenever
  configuration, tests, or keybindings require it. Run focused tests per step;
  after the shared menu change also run existing menu integration coverage.
  Record unavailable gates instead of reporting success. `just` was not found
  on PATH during planning; installed at `/tmp/lazygit-tools/bin/just`.
  Run unit tests with `env -u NO_COLOR` because the session sets `NO_COLOR=1`
  and presentation tests assert ANSI output.
- Planning-only validation is final diff review and `git diff --check`; no
  expensive application suites are needed for this prose-only commit, per the
  invoked skill. Implementation results are recorded per step below.
- Follow repository commit conventions: small coherent commits, `[codex]`
  titles under 60 characters including the stable implementation step number,
  explanatory body wrapped within 72 columns (also within the global 80-column
  limit), and a final model co-author trailer using the harness identity.
  Refinements use separate `fixup!`/`amend!` commits; never amend or autosquash.
  Stop for genuine design forks or newly discovered bugs as AGENTS requires.
- Routine scoped commits are authorized by repository policy. Branch rewriting,
  publishing, and external writes are not authorized. No PRs may be created.
- Before handoff or context exhaustion, persist active step, HEAD and relevant
  commits, dirty-file ownership, commands/results, artifacts, processes/agents,
  open gates, decisions, and next action here. Successors reconcile actual
  checkout/process state first, preserve user changes, and finish/account for
  live processes before dispatching another agent. Partial checks stay partial.
- Final handoff: Step 1 is `55a4e37bd`; Step 2 is `3c350ea39`; Step 3
  completes in this commit. No processes or agents are running. Logs are under
  `/tmp/lazygit-step3-{unit,e2e}.log`. All local implementation gates passed;
  external CI remains unexecuted. No further implementation is pending.

## Worklog

- 2026-11-03 — Planning at baseline `0002c0354770ff8ec6abe4d5ab4d025ea40e1619`:
  traced push/menu/config/test/CI callers and drafted Steps 1–3. No application
  checks executed. Final plan diff reviewed; `git diff --check` passed.
  All implementation and execution gates remain open.

- 2026-11-03 — Step 1 (`55a4e37bd`): added `enableStackedPush` and a disabled
  integration case. `just generate`, `just format`, `just build`,
  `env -u NO_COLOR just unit-test`, `just lint` (0 issues), and
  `just e2e 'sync/push_stacked_branches.*'` passed. The initial unit run
  failed solely on ANSI expectations with `NO_COLOR=1`; the clean-environment
  rerun passed. Steps 2–3 and their gates remain open.

- 2026-11-03 — Step 2 (`3c350ea39`), after `55a4e37bd`: added visible `a`/`o`
  shortcuts and exercised both existing push paths through keypresses. Added
  a navigation-conflict/cancel integration case. `just generate`, `just format`,
  `just build`, `env -u NO_COLOR just unit-test`, `just lint` (0 issues), and
  `just e2e 'sync/push_stacked_branches.*'` passed. Initial e2e assertions
  were corrected for prompt rows and adjacency before the passing rerun.
  Step 3 and its gates remain open.

- 2026-11-03 — Step 3 (this commit), after `3c350ea39`: added
  `stackedPushDefault`, runtime/schema enum validation, and an initial-selection
  menu option with first-item fallback for invalid indices. Tests cover Enter
  for both defaults, reopening/reset, shortcuts independent of selection, both
  disabled/default combinations, ordinary push, and an unrelated menu.
  `just generate`, `just format`, `just build`,
  `env -u NO_COLOR just unit-test`, and `just lint` (0 issues) passed.
  `just e2e 'sync/.*' 'ui/menu_scroll_position_is_reset' 'ui/empty_menu'
  'filter_and_search/filter_menu.*' 'custom_commands/.*menu.*'` passed.
  Final diff reviewed and `git diff --check` passed. No implementation gates
  remain open locally; external CI was not executed.
