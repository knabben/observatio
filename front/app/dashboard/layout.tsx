import React from "react";
import SideNav from '@/app/ui/dashboard/sidenav';
import {openSans} from '@/fonts';

export default function Layout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className={openSans.className}>
      <div className="flex h-screen flex-col overflow-hidden md:flex-row">
        <div className="w-full flex-none md:w-64">
          <SideNav/>
        </div>
        <div className="flex-grow overflow-y-auto p-6 md:p-12">{children}</div>
      </div>
    </div>
  )
}
