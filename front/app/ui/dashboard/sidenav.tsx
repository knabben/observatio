import Link from 'next/link';
import Image from 'next/image'
import NavLinks from '@/app/ui/dashboard/nav-links';

export default function SideNav() {
  return (
    <div className="flex h-full flex-col px-3 py-4 md:px-2">
        <div className="mb-2 flex h-20 items-end justify-center rounded-md bg-[#000] p-1 md:h-55">
          <Image
            src="/logo.png"
            alt={`logo`}
            width={350}
            height={350}
          />
        </div>
      <div className="flex grow flex-row justify-between space-x-2 md:flex-col md:space-x-0 md:space-y-2">
        <NavLinks/>
        <div className="hidden h-auto w-full grow rounded-md bg-gray-50 md:block"></div>
      </div>
    </div>
  );
}
