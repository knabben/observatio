import React from "react";
import MachineHealthCheckLister from '@/app/ui/dashboard/components/machinehealthchecks/lister'

export default async function MachineHealthChecks() {
  return (
    <div>
      <main>
        <MachineHealthCheckLister />
      </main>
    </div>
  )
}
