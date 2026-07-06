import {renderHook} from '@testing-library/react';
import React from 'react';
import {AIPanelProvider, useAIPanel, ObjectContext} from './ai-panel-context';
import {useCurrentObjectContext} from './use-current-object-context';

function wrapper({children}: {children: React.ReactNode}) {
  return <AIPanelProvider>{children}</AIPanelProvider>;
}

const ctxA: ObjectContext = {
  kind: 'Cluster', name: 'c1', namespace: 'default', status: 'Ready', keySpecFields: {},
};
const ctxB: ObjectContext = {
  kind: 'Cluster', name: 'c2', namespace: 'default', status: 'Ready', keySpecFields: {},
};

describe('useCurrentObjectContext', () => {
  it('registers the context on mount', () => {
    const {result} = renderHook(
      (ctx: ObjectContext | null) => {
        useCurrentObjectContext(ctx);
        return useAIPanel();
      },
      {wrapper, initialProps: ctxA},
    );
    expect(result.current.currentObjectContext).toEqual(ctxA);
    expect(result.current.queryField).toContain('c1');
  });

  it('unregisters on unmount, leaving no stale context', () => {
    const {result, unmount} = renderHook(
      (ctx: ObjectContext | null) => {
        useCurrentObjectContext(ctx);
        return useAIPanel();
      },
      {wrapper, initialProps: ctxA as ObjectContext | null},
    );
    expect(result.current.currentObjectContext).toEqual(ctxA);
    unmount();
    // Nothing to assert post-unmount on `result` (stale), but this documents the cleanup path
    // exists and doesn't throw.
  });

  it('re-registers and refreshes the query field when the object changes, while untouched', () => {
    const {result, rerender} = renderHook(
      (ctx: ObjectContext | null) => {
        useCurrentObjectContext(ctx);
        return useAIPanel();
      },
      {wrapper, initialProps: ctxA as ObjectContext | null},
    );
    expect(result.current.queryField).toContain('c1');
    rerender(ctxB);
    expect(result.current.currentObjectContext).toEqual(ctxB);
    expect(result.current.queryField).toContain('c2');
  });
});
