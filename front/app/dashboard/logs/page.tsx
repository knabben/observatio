import React, {Suspense} from "react";

import {LogsView} from '@/app/ui/dashboard/logs/logs-view'

export default async function Logs() {
  return (
    <Suspense>
      <LogsView/>
    </Suspense>
  )
}
