import {Conditions, Meta} from "@/app/ui/dashboard/base/types";

/**
 * Represents the type definition for a ClusterClass — reuses the existing flat
 * `models.ClusterClass` shape (name/namespace/generation/conditions) already served by the
 * main-dashboard widget, plus a `metadata` mirror added by the watcher so this first-class page
 * can reuse the same `BaseLister`/`ObjectTable` row-key and search conventions as every other kind.
 */
export type ClusterClassType = {
  metadata?: Meta,
  name?: string,
  namespace?: string,
  generation?: number,
  conditions?: Conditions[],
}
