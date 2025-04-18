import Link from 'next/link';
import Image from 'next/image'
import NavLinks from '@/app/ui/dashboard/nav-links';
import { Select } from '@mantine/core';

export default function SideNav() {
  return (
    <div className="flex h-full flex-col px-3 py-4 md:px-2">
      <Link
        className="mb-2 flex h-20 items-end justify-center rounded-md bg-[#A4A89C] p-4 md:h-40"
        href="/dashboard"
      >
        <div className="w-32 text-white md:w-40">
          <Image
            src="/logo.png"
            alt={`logo`}
            width={200}
            height={200}
          />
        </div>
      </Link>
      <div className="flex grow flex-row justify-between space-x-2 md:flex-col md:space-x-0 md:space-y-2">
        <form>
          <Select
            placeholder="Kubernetes Context"
            data={['Default']}
          />
        </form>
        <NavLinks/>
        <div className="hidden h-auto w-full grow rounded-md bg-gray-50 md:block"></div>
      </div>
    </div>
  );
}
