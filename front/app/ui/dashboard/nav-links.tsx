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

/** Nested routes (e.g. `/dashboard/clusters/x`) should highlight their parent link. */
function isActiveLink(pathname: string, href: string): boolean {
  if (href === '/dashboard') return pathname === '/dashboard';
  return pathname === href || pathname.startsWith(`${href}/`);
}

export default function NavLinks() {
  const pathname = usePathname();

  return (
    <>
      {links.map((link) => {
        const LinkIcon = link.icon;
        const active = isActiveLink(pathname, link.href);
        return (
          <Link
            key={link.name}
            href={link.href}
            prefetch={false}
            aria-label={link.name}
            aria-current={active ? 'page' : undefined}
            className={clsx(
                'flex h-[48px] grow items-center justify-center gap-2 rounded-md text-sm md:flex-none md:justify-start md:p-2 md:px-3',
                {
                  'bg-[var(--mantine-color-brand-4)] text-black font-bold': active,
                },
            )}
          >
            <LinkIcon className="w-6" aria-hidden="true" />
            <p className="hidden md:block">{link.name}</p>
          </Link>
        );
      })}
    </>
  );
}
