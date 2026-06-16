'use client';

import { useState } from 'react';
import { useAuth } from './AuthProvider';
import { RepoProvider } from './RepoProvider';
import Sidebar from './Sidebar';
import { usePathname } from 'next/navigation';

export default function AppLayoutWrapper({ children }: { children: React.ReactNode }) {
  const { user, loading } = useAuth();
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const pathname = usePathname();

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-[#09090b] text-zinc-400 font-sans">
        <div className="flex flex-col items-center gap-3">
          <div className="animate-spin rounded-full h-8 w-8 border-2 border-zinc-800 border-t-cyan-500"></div>
          <span className="text-sm font-medium tracking-wide">Loading workspace...</span>
        </div>
      </div>
    );
  }

  // If not logged in, render layout without sidebar (e.g. login screen)
  if (!user) {
    const isPublicPath = pathname === '/login' || pathname === '/';
    if (!isPublicPath) {
      return (
        <div className="flex items-center justify-center min-h-screen bg-[#09090b] text-zinc-400 font-sans">
          <div className="flex flex-col items-center gap-3">
            <div className="animate-spin rounded-full h-8 w-8 border-2 border-zinc-800 border-t-cyan-500"></div>
            <span className="text-sm font-medium tracking-wide">Redirecting...</span>
          </div>
        </div>
      );
    }
    return <>{children}</>;
  }

  return (
    <RepoProvider>
      <div className="flex h-screen w-screen overflow-hidden bg-[#09090b] text-zinc-100 font-sans">
        {/* Desktop Sidebar (Fixed Left) */}
        <aside className="hidden md:block w-72 h-full flex-shrink-0">
          <Sidebar />
        </aside>

        {/* Mobile Sidebar Slide-out Drawer */}
        {mobileMenuOpen && (
          <div className="fixed inset-0 z-50 flex md:hidden animate-fadeIn">
            {/* Backdrop */}
            <div 
              className="fixed inset-0 bg-black/60 backdrop-blur-sm transition-opacity"
              onClick={() => setMobileMenuOpen(false)}
            />
            {/* Drawer Content */}
            <div className="relative flex-1 flex flex-col max-w-[280px] w-full bg-[#0c0c0e] h-full shadow-2xl animate-slideRight">
              <Sidebar onCloseMobile={() => setMobileMenuOpen(false)} />
              <button
                onClick={() => setMobileMenuOpen(false)}
                className="absolute top-5 right-[-45px] p-2 text-zinc-400 hover:text-white focus:outline-none"
              >
                <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
          </div>
        )}

        {/* Main Area */}
        <div className="flex-1 flex flex-col min-w-0 h-full relative">
          {/* Mobile Top Navigation Header */}
          <header className="flex md:hidden items-center justify-between px-6 py-4 bg-[#0c0c0e] border-b border-zinc-900 flex-shrink-0">
            <div className="flex items-center gap-2.5">
              <div className="w-6.5 h-6.5 rounded bg-gradient-to-tr from-cyan-500 to-blue-500 flex items-center justify-center text-black font-bold text-xs shadow-md shadow-cyan-500/10">
                IB
              </div>
              <span className="font-semibold text-sm text-zinc-100">IssueBoard</span>
            </div>
            <button
              onClick={() => setMobileMenuOpen(true)}
              className="p-1 text-zinc-400 hover:text-zinc-150 focus:outline-none"
            >
              <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
              </svg>
            </button>
          </header>

          {/* Page content view */}
          <main className="flex-grow overflow-auto min-w-0">
            {children}
          </main>
        </div>
      </div>
    </RepoProvider>
  );
}
