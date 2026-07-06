import {useEffect} from 'react';
import {ObjectContext, useAIPanel} from '@/app/ui/dashboard/ai-panel/ai-panel-context';

/**
 * Registers a detail screen's current object with the global AI panel (FR-006, FR-007), so
 * opening the panel from anywhere pre-fills a description of whatever is actually in view.
 * Unregisters on unmount so navigating to a list/overview screen leaves no stale context behind.
 */
export function useCurrentObjectContext(ctx: ObjectContext | null) {
  const {setCurrentObjectContext} = useAIPanel();
  const specFieldsKey = ctx ? JSON.stringify(ctx.keySpecFields) : '';

  useEffect(() => {
    setCurrentObjectContext(ctx);
    return () => setCurrentObjectContext(null);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [ctx?.kind, ctx?.name, ctx?.namespace, ctx?.status, specFieldsKey, setCurrentObjectContext]);
}
