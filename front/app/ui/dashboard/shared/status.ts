/**
 * Normalized health used by every status indicator.
 *
 * - `healthy`  : readiness flag present and true
 * - `degraded` : not failed, but has an open risk warning or an implicated/inconclusive debugging
 *   layer (Day-2 Ops dashboard only — never produced by `toStatusState`/`allReady` below)
 * - `notready` : readiness flag present and false
 * - `unknown`  : readiness flag absent/undefined
 *
 * Strict comparisons only — an absent field is NEVER treated as `false`/failed.
 */
export type StatusState = 'healthy' | 'degraded' | 'notready' | 'unknown';

/**
 * Derives a tri-state from a readiness value. `undefined`/`null` → `unknown`,
 * so a resource with unknown availability is never rendered as failed.
 */
export function toStatusState(ready: boolean | undefined | null): StatusState {
  if (ready === true) return 'healthy';
  if (ready === false) return 'notready';
  return 'unknown';
}

/** True only when every provided readiness flag is explicitly `true`. */
export function allReady(...flags: Array<boolean | undefined | null>): StatusState {
  if (flags.some(f => f == null)) return 'unknown';
  return flags.every(f => f === true) ? 'healthy' : 'notready';
}
