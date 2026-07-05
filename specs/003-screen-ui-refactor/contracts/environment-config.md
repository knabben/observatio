# Contract: Environment Configuration — *revised per Clarification 2026-07-05*

The SPA is embedded and served by the **same Go binary** that serves the API/WebSocket, so endpoints are
resolved **same-origin at runtime** from `window.location` — not baked in at build time. There is no server
runtime in the exported bundle, but origin derivation is a client computation, so this needs no server. A
`NEXT_PUBLIC_*` override exists **only** for the split development mode (frontend dev server → separate
backend).

## `front/app/lib/config.ts`

```ts
// Same-origin by default; the browser is already on the correct host/port
// because the binary that served this page also serves the API + WS.
function sameOriginWs(path: string): string {
  if (typeof window === 'undefined') return path;            // SSR/export build: relative
  const proto = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  return `${proto}//${window.location.host}${path}`;
}

// Dev override (frontend :3000 → backend :8080); unset in the embedded build.
export const API_URL =
  process.env.NEXT_PUBLIC_API_URL ?? '';                      // '' ⇒ relative to origin
export const WS_URL =
  process.env.NEXT_PUBLIC_WS_URL ?? sameOriginWs('/ws');
export const WS_URL_CHATBOT =
  process.env.NEXT_PUBLIC_WS_URL_CHATBOT ?? sameOriginWs('/chatbot');
```

**Contract**:
- All modules import endpoints from `config.ts`; the hardcoded `URL` const in `data.tsx` is removed. (FR-036)
- **Default is same-origin**: REST uses relative paths (`API_URL === ''`), WebSocket URLs derive from
  `window.location`. The embedded production binary contains **zero hardcoded backend addresses**. (SC-009)
- Relocating the binary to another host/port requires **no rebuild** — addressing follows the serving origin.
- The `NEXT_PUBLIC_*` overrides are honored only for split dev mode; they are unset in `make build`.

## Variables (development override only)

| Variable | Purpose | Embedded prod default |
|----------|---------|-----------------------|
| `NEXT_PUBLIC_API_URL` | REST base URL for split dev | unset ⇒ same-origin (relative) |
| `NEXT_PUBLIC_WS_URL` | live resource WebSocket for split dev | unset ⇒ `ws(s)://<origin>/ws` |
| `NEXT_PUBLIC_WS_URL_CHATBOT` | AI chat WebSocket for split dev | unset ⇒ `ws(s)://<origin>/chatbot` |

`make run-frontend` (dev on :3000) sets these to point at the backend on :8080. Document them in the
README / `.env.example` as **dev-only**. Actual endpoint paths (`/ws`, `/chatbot`, REST routes) must match
the backend router — confirm against `webserver/internal/web/handlers` during implementation.
