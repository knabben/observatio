import { sourceCodePro400 } from '@/fonts'

export default function Home() {
  return (
    <main>
      <h1 className={`${sourceCodePro400.className} mb-4 text-xl md:text-2xl`}>
        Dashboard
      </h1>
      <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
        the middle
        cluster summary and info
      </div>
      <div className="mt-6 grid grid-cols-1 gap-6 md:grid-cols-4 lg:grid-cols-8">
        the other side
      </div>
    </main>
  );
}
