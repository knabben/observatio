'use client';

import {
  HomeIcon,
  DocumentDuplicateIcon,
  BanknotesIcon,
  BookOpenIcon,
} from '@heroicons/react/24/outline';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import clsx from 'clsx';

const links = [
  { name: 'Dashboard', href: '/dashboard', icon: HomeIcon },
  {
    name: 'Clusters',
    href: '/dashboard/clusters',
    icon: DocumentDuplicateIcon,
  },
  {
    name: 'Machines Deployment',
    href: '/dashboard/machinedeployments',
    icon: BanknotesIcon,
  },
  {
    name: 'Machines',
    href: '/dashboard/machines',
    icon: BookOpenIcon,
  },
];

export default function NavLinks() {
  const pathname = usePathname();

  return (
    <>
      {links.map((link) => {
        const LinkIcon = link.icon;
        return (
          <Link
            key={link.name}
            href={link.href}
            prefetch={false}
            className={clsx(
                'flex h-[48px] grow items-center justify-center gap-2 rounded-md text-sm md:flex-none md:justify-start md:p-2 md:px-3',
                {
                  'bg-[#aaf16a] text-black font-bold': pathname.startsWith(link.href),
                },
            )}
          >
            <LinkIcon className="w-6" />
            <p className="hidden md:block">{link.name}</p>
          </Link>
        );
      })}
    </>
  );
}
