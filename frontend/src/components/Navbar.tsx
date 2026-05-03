'use client';

import { useAuth } from './AuthProvider';
import Link from 'next/link';

export default function Navbar() {
  const { user, signOut } = useAuth();

  if (!user) return null;

  return (
    <nav className="bg-gray-800 text-white p-4 flex justify-between items-center">
      <Link href="/inventory" className="text-xl font-bold">IssueBoard</Link>
      <div className="flex items-center gap-4">
        <span>{user.email}</span>
        <button
          onClick={signOut}
          className="bg-red-600 px-3 py-1 rounded hover:bg-red-700"
        >
          Logout
        </button>
      </div>
    </nav>
  );
}
