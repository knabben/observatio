/**
 * Same-origin endpoint configuration.
 *
 * The SPA is embedded in and served by the same Go binary that serves the API and
 * WebSocket, so by default endpoints are derived from the page origin — one binary
 * works on any host/port with no rebuild. The NEXT_PUBLIC_* overrides exist ONLY for
 * the split development mode (frontend dev server on :3000 → backend on :8080) and are
 * unset in the embedded production build.
 *
 * See specs/003-screen-ui-refactor/contracts/environment-config.md
 */

function wsFromOrigin(path: string): string {
  if (typeof window === 'undefined') return path;
  const proto = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  return `${proto}//${window.location.host}${path}`;
}

/** REST base URL. Empty string ⇒ relative to the serving origin (production). */
export const API_URL = process.env.NEXT_PUBLIC_API_URL ?? '';

/** Live resource watcher WebSocket. */
export const WS_URL_WATCHER =
  process.env.NEXT_PUBLIC_WS_URL ?? wsFromOrigin('/ws/watcher');

/** AI troubleshooting chat WebSocket. */
export const WS_URL_CHATBOT =
  process.env.NEXT_PUBLIC_WS_URL_CHATBOT ?? wsFromOrigin('/ws/analysis');
