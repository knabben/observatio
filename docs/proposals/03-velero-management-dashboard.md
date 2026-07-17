# Quick specification statement: Velero management dashboard

Once Velero objects are being watched (proposal 01), a summary rollup on the Day-2 Ops landing page is
enough to answer "are we protected right now," but not enough to actually manage backups and restores.
Add a dedicated "Backups" destination to Observātiō's lateral navigation
(`front/app/ui/dashboard/nav-links.tsx`), following the same route-shell convention 006's new "Logs"
destination uses (`layout.tsx` + `page.tsx` under `front/app/dashboard/`) rather than a one-off page
structure.

The view should list `Backup`, `Schedule`, and `Restore` objects with their health/staleness at a
glance, mirroring the existing Clusters/Machines/MachineDeployments list-page pattern already in this
codebase. Distinguish two levels of scope for planning: a **read-only v1** that only visualizes backup
state (list, describe, health, staleness — consistent with the rest of the product's current read-only
constitution) versus a **v2 with guarded write actions** (triggering an on-demand backup, triggering a
restore with `--existing-resource-policy=update` per the companion guide's recovery workflow, or a
one-click pause/unpause paired with proposal 01's `Cluster.spec.paused` visibility). The write-capable
version is a meaningfully bigger step for this product and should be scoped as an explicit, separate
decision rather than assumed as part of the initial dashboard. Also track whether a backup has been
validated by a recent recovery test (the guide's "an untested backup is not a backup" principle),
even if only as a manually-recorded field in this first iteration.
