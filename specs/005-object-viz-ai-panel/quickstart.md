# Quickstart: Verifying the YAML Tab & Global AI Panel

Manual verification steps once the feature is implemented, covering each acceptance scenario in
`spec.md`.

## 1. YAML tree tab on every detail screen

1. Open a Cluster's detail screen → select the new "YAML" tab → confirm every field the backend
   returns (including fields absent from "Specification", e.g. `spec.clusterNetwork` sub-fields not
   surfaced today) appears as an expandable tree.
2. Collapse a nested node (e.g. `status`) → confirm its children hide and re-expanding restores them.
3. Repeat for Machine, Machine Deployment, and both Cluster/Machine Infrastructure variants (Docker
   and vSphere) → confirm the tab appears and works identically on all six.
4. While the tab is open, cause the underlying object to change (e.g. wait for a status condition to
   flip) → confirm the tree updates rather than staying stale.
5. Confirm the "Specification" tab still shows the header status indicator + conditions table exactly
   as before — no separate "Status" tab was introduced, and nothing that was visible before is now
   missing.

## 2. Global AI panel

1. From the Clusters list, the Dashboard overview, and a Machine detail screen, confirm the same
   persistent AI panel trigger is visible and opens the same collapsible panel each time.
2. Open the panel from a Machine's detail screen → confirm the query field is pre-filled with that
   Machine's identity, status/conditions, and key spec fields — not just a bare condition-reason
   string.
3. Open the panel from the Clusters list (no single object in focus) → confirm the field is empty or
   shows a general prompt, not a broken partial string.
4. Edit the pre-filled text, then send → confirm the edited text (not the original auto-fill) is
   what's sent.
5. With the panel open and an object's context freshly auto-filled but not yet edited/sent, navigate
   to a different object's detail screen → confirm the field updates to the new object's context.
6. Send a message, collapse the panel, reopen it → confirm the conversation (including that message
   and its response once it arrives) is still there.
7. From an object's detail screen, click that screen's "Ask AI about this" quick-action → confirm the
   panel opens already pre-filled with that object's context in one click.
8. Visually compare the AI panel's colors/contrast against the rest of the dashboard in both light
   and dark theme → confirm no hardcoded colors that clash with the active theme.
9. At a narrow (mobile) viewport, open the panel → confirm it remains usable and closable, not
   permanently obscuring the screen.

## Automated coverage

- Backend: `make run-tests-backend` — tests for the raw-object handler's GVR validation (missing
  param → 400, unknown GVR → 400/404).
- Frontend: `make run-tests-frontend` — Jest tests for `to-tree-data.ts` (nested objects/arrays,
  scalars, empty object), the AI panel context (auto-context refresh, edit-locks-prefill, per-object
  quick-action), and each detail screen's updated tab set (no AI Troubleshooting tab, new YAML tab
  present).
