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
    <html lang="en">
      <body className={inter.className}>
        <div className="min-h-screen flex flex-col">
          <header className="p-4 border-b border-slate-700 bg-slate-900/50">
            <h1 className="text-xl font-bold tracking-tight">Telemetry Engine</h1>
          </header>
          
          <main className="flex-1 container mx-auto p-4">
            {children}
          </main>

          <footer className="p-4 text-center text-slate-500 text-sm border-t border-slate-800">
            System Status: Connected
          </footer>
        </div>
      </body>
    </html>
  );
}