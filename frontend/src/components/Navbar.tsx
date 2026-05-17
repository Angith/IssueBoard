'use client';

import { useAuth } from './AuthProvider';
import Link from 'next/link';
import { useState } from 'react';
import SettingsModal from './SettingsModal';

export default function Navbar() {
  const { user, signOut } = useAuth();
  const [isSettingsOpen, setIsSettingsOpen] = useState(false);

  if (!user) return null;

  return (
    <>
      <nav className="border-b border-zinc-800 bg-[#09090b] px-6 py-4 flex justify-between items-center">
        <Link href="/inventory" className="text-xl font-semibold text-zinc-100 tracking-tight hover:text-white transition-colors">
          IssueBoard
        </Link>
        <div className="flex items-center gap-6">
          <span className="text-sm font-medium text-zinc-400">{user.email}</span>
          <div className="flex items-center gap-3">
            <button
              onClick={() => setIsSettingsOpen(true)}
              className="text-sm font-medium text-zinc-400 hover:text-zinc-100 transition-colors"
            >
              Settings
            </button>
            <button
              onClick={signOut}
              className="text-sm font-medium text-zinc-400 hover:text-red-400 transition-colors"
            >
              Logout
            </button>
          </div>
        </div>
      </nav>
      <SettingsModal isOpen={isSettingsOpen} onClose={() => setIsSettingsOpen(false)} />
    </>
  );
}
