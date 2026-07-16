# Quick specification statement: controller metrics integration

The companion guide's diagnostics section is built around three endpoints every CAPI controller
already exposes — a Prometheus `/metrics` endpoint (reconciliation latency, work-queue depth, error
counts), a dynamic log-level endpoint (`/debug/flags/v`), and pprof profiling — and recommends scraping
these into Prometheus/Grafana with alerts on reconciliation failures and queue saturation. Observātiō's
006 feature detects "Level 2" provider-controller degradation purely from Pod/Deployment status
(crash-looping, not Ready) — it never looks at what the controller itself is reporting about its own
reconciliation health, so a controller that's alive but silently falling behind (queue depth climbing,
error rate rising, no crash yet) is invisible today.

Add a metrics-ingestion path that scrapes or queries each watched controller's `/metrics` endpoint
(CAPI core in `capi-system`, and each infrastructure/bootstrap provider's controller in its own
namespace — the same controllers 006's `severity.go` and the new Logs view already resolve
Deployment→Pod for) and feeds reconciliation-latency, queue-depth, and error-rate signals into the
existing severity classification as an earlier-warning input, ahead of an outright crash. This is
explicitly a follow-on to 006, not a replacement for its Pod-Ready heuristic — it should sit alongside it
as a second, complementary signal. Planning should decide whether Observātiō queries these endpoints
directly (no new external dependency, consistent with the product's current no-new-infra stance) or
expects an existing Prometheus to be present and only reads from it.
