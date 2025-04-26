import Link from 'next/link';
import Image from 'next/image'
import NavLinks from '@/app/ui/dashboard/nav-links';

export default function SideNav() {
  return (
    <div className="flex h-full flex-col px-3 py-4 md:px-2">
        <Image
          className="mb-3 flex items-end justify-center rounded-md p-0 md:h-55"
          src="/logo.png"
          alt={`logo`}
          width={250}
          height={250}
        />
      <div className="flex grow flex-row justify-between space-x-2 md:flex-col md:space-x-0 md:space-y-2">
        <NavLinks/>
        <div className="hidden h-auto w-full grow rounded-md md:block"></div>
      </div>
    </div>
  );
}
