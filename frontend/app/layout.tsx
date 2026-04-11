import './globals.css';
import { Inter } from 'next/font/google';

const inter = Inter({ subsets: ['latin'] });

export const metadata = {
  title: 'Telemetry Engine | Real-time Dashboard',
  description: 'Scalable system monitoring powered by Go and Next.js',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className="dark">
      <body className={`${inter.className} bg-[#050505] text-slate-200 antialiased overflow-hidden`}>
        <div className="relative h-screen w-screen flex flex-col">
          <header className="absolute top-0 left-0 right-0 z-10 p-6 pointer-events-none">
            <div className="flex items-center justify-between">
              <div>
                <h1 className="text-xl font-black tracking-tighter uppercase text-white">
                  Telemetry <span className="text-blue-500">Engine</span>
                </h1>
                <p className="text-[10px] font-mono text-slate-500 uppercase tracking-[0.2em]">
                  Real-time Data Stream v1.0
                </p>
              </div>
            </div>
          </header>

          <main className="flex-1 h-full w-full">
            {children}
          </main>

          <footer className="absolute bottom-4 left-4 z-10 pointer-events-none">
            <div className="flex items-center gap-2 px-3 py-1 rounded-full bg-slate-900/80 border border-slate-800 backdrop-blur-md">
              <div className="h-2 w-2 rounded-full bg-emerald-500 animate-pulse" />
              <span className="text-[10px] font-mono text-slate-300 uppercase tracking-widest">
                System Status: Active
              </span>
            </div>
          </footer>
        </div>
      </body>
    </html>
  );
}