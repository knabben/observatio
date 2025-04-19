'use client';

import {
  HomeIcon,
  DocumentDuplicateIcon,
  BanknotesIcon,
  BookOpenIcon,
  MagnifyingGlassPlusIcon,
} from '@heroicons/react/24/outline';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

import clsx from 'clsx';

// Map of links to display in the side navigation.
// Depending on the size of the application, this would be stored in a database.
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
            className={clsx(
                'flex h-[48px] grow items-center justify-center gap-2 rounded-md bg-gray-50 p-3 text-sm font-medium hover:bg-[#8BB94B] hover:text-black-600 md:flex-none md:justify-start md:p-2 md:px-3',
                {
                  'bg-[#8BB94B] text-black-600': pathname === link.href,
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
