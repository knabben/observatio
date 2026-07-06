import {act, renderHook} from '@testing-library/react';
import React from 'react';
import {AIPanelProvider, useAIPanel, formatObjectContext} from './ai-panel-context';

function wrapper({children}: {children: React.ReactNode}) {
  return <AIPanelProvider>{children}</AIPanelProvider>;
}

describe('AIPanelProvider / useAIPanel', () => {
  it('starts closed with no messages', () => {
    const {result} = renderHook(() => useAIPanel(), {wrapper});
    expect(result.current.isOpen).toBe(false);
    expect(result.current.messages).toEqual([]);
  });

  it('open() opens the panel, close() closes it', () => {
    const {result} = renderHook(() => useAIPanel(), {wrapper});
    act(() => result.current.open());
    expect(result.current.isOpen).toBe(true);
    act(() => result.current.close());
    expect(result.current.isOpen).toBe(false);
  });

  it('preserves messages across close and reopen (FR-003)', () => {
    const {result} = renderHook(() => useAIPanel(), {wrapper});
    act(() => {
      result.current.open();
      result.current.setMessages([
        {id: '1', type: 'chatbot', agent_id: 'a', content: 'hello', timestamp: 'now', actor: 'user'},
      ]);
    });
    act(() => result.current.close());
    act(() => result.current.open());
    expect(result.current.messages).toHaveLength(1);
    expect(result.current.messages[0].content).toBe('hello');
  });

  it('open(prefill) sets the query field and clears the touched flag', () => {
    const {result} = renderHook(() => useAIPanel(), {wrapper});
    act(() => result.current.open('some prefill'));
    expect(result.current.queryField).toBe('some prefill');
    expect(result.current.queryFieldTouched).toBe(false);
  });

  it('setQueryField marks the field as touched', () => {
    const {result} = renderHook(() => useAIPanel(), {wrapper});
    act(() => result.current.setQueryField('operator typed this'));
    expect(result.current.queryField).toBe('operator typed this');
    expect(result.current.queryFieldTouched).toBe(true);
  });

  it('refreshes the query field from currentObjectContext only while untouched', () => {
    const {result} = renderHook(() => useAIPanel(), {wrapper});
    act(() => {
      result.current.setCurrentObjectContext({
        kind: 'Cluster', name: 'c1', namespace: 'default', status: 'Ready', keySpecFields: {},
      });
    });
    expect(result.current.queryField).toContain('Cluster "c1"');

    act(() => result.current.setQueryField('operator edit'));
    act(() => {
      result.current.setCurrentObjectContext({
        kind: 'Cluster', name: 'c2', namespace: 'default', status: 'Ready', keySpecFields: {},
      });
    });
    expect(result.current.queryField).toBe('operator edit');
  });
});

describe('formatObjectContext', () => {
  it('includes identity, status, and key spec fields', () => {
    const text = formatObjectContext({
      kind: 'Machine',
      name: 'm1',
      namespace: 'default',
      status: 'NotReady: InfrastructureNotReady',
      keySpecFields: {version: 'v1.30.0', providerID: 'docker:///m1'},
    });
    expect(text).toContain('Machine "m1"');
    expect(text).toContain('default');
    expect(text).toContain('NotReady: InfrastructureNotReady');
    expect(text).toContain('version=v1.30.0');
    expect(text).toContain('providerID=docker:///m1');
  });

  it('omits the key-spec suffix when there are no fields', () => {
    const text = formatObjectContext({
      kind: 'Cluster', name: 'c1', namespace: 'default', status: 'Ready', keySpecFields: {},
    });
    expect(text).not.toContain('Key spec fields');
  });
});
