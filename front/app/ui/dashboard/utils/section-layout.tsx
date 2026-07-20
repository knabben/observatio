import React from 'react';

/**
 * Shared no-op layout for dashboard route segments that add nothing beyond what
 * front/app/dashboard/layout.tsx already provides (SideNav, AI panel, theme font). Next.js App
 * Router requires a layout.tsx per segment for it to compose cleanly with its children — this
 * exists so that requirement doesn't mean re-typing `<div>{children}</div>` in every segment.
 */
export default function SectionLayout({children}: {children: React.ReactNode}) {
  return <div>{children}</div>;
}
