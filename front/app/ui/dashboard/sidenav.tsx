'use client';

import Image from 'next/image'
import {Burger, Drawer} from '@mantine/core';
import {useDisclosure} from '@mantine/hooks';
import NavLinks from '@/app/ui/dashboard/nav-links';

export default function SideNav() {
  const [opened, {toggle, close}] = useDisclosure(false);

  const logo = (
    <Image
      className="mb-3 flex items-end justify-center rounded-md p-0 md:h-55"
      src="/logo.png"
      alt="Observatio logo"
      sizes="(max-width: 768px) 40vw, 250px"
      width={250}
      height={250}
    />
  );

  return (
    <div className="flex h-full flex-col px-3 py-4 md:px-2">
      <div className="flex items-center justify-between md:block">
        {logo}
        <Burger opened={opened} onClick={toggle} aria-label="Toggle navigation menu" hiddenFrom="md"/>
      </div>
      <div className="hidden grow flex-col space-y-2 md:flex">
        <NavLinks/>
        <div className="hidden h-auto w-full grow rounded-md md:block"></div>
      </div>
      <Drawer opened={opened} onClose={close} title="Navigation" hiddenFrom="md" size="xs">
        <div className="flex flex-col space-y-2" onClick={close}>
          <NavLinks/>
        </div>
      </Drawer>
    </div>
  );
}
