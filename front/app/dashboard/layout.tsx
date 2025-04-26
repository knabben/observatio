import React from "react";
import SideNav from '@/app/ui/dashboard/sidenav';
import {Roboto} from 'next/font/google';

const roboto = Roboto({
  weight: '200',
  subsets: ['latin'],
})

export default function Layout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className={roboto.className}>
      <div className="flex h-screen flex-col md:flex-row md:overflow-hidden">
        <div className="w-full flex-none md:w-64">
          <SideNav/>
        </div>
        <div className="flex-grow p-6 md:overflow-y-auto md:p-12">{children}</div>
      </div>
    </div>
  )
}
