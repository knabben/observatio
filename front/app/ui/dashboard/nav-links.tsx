'use client';

import {
  HomeIcon,
  DocumentDuplicateIcon,
  BanknotesIcon,
  BookOpenIcon,
  DocumentTextIcon,
  HeartIcon,
  CpuChipIcon,
  Squares2X2Icon,
  RectangleGroupIcon,
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
  {
    name: 'Machine Health Checks',
    href: '/dashboard/machinehealthchecks',
    icon: HeartIcon,
  },
  {
    name: 'Kubeadm Control Planes',
    href: '/dashboard/kubeadmcontrolplanes',
    icon: CpuChipIcon,
  },
  {
    name: 'Machine Sets',
    href: '/dashboard/machinesets',
    icon: Squares2X2Icon,
  },
  {
    name: 'Cluster Classes',
    href: '/dashboard/clusterclasses',
    icon: RectangleGroupIcon,
  },
  {
    name: 'Logs',
    href: '/dashboard/logs',
    icon: DocumentTextIcon,
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
