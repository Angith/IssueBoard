'use client';

import { useState } from 'react';
import Link from 'next/link';
import { usePathname, useRouter } from 'next/navigation';
import { useAuth } from './AuthProvider';
import { useRepos } from './RepoProvider';
import SettingsModal from './SettingsModal';

interface SidebarProps {
  onCloseMobile?: () => void;
}

export default function Sidebar({ onCloseMobile }: SidebarProps) {
  const { user, signOut } = useAuth();
  const { repos, loading, addRepo, error: repoContextError } = useRepos();
  const pathname = usePathname();
  const router = useRouter();

  const [search, setSearch] = useState('');
  const [isSettingsOpen, setIsSettingsOpen] = useState(false);
  
  // Quick Add Form state
  const [showAddForm, setShowAddForm] = useState(false);
  const [newUrl, setNewUrl] = useState('');
  const [adding, setAdding] = useState(false);
  const [addError, setAddError] = useState('');

  const filteredRepos = repos.filter((repo) =>
    repo.full_name.toLowerCase().includes(search.toLowerCase())
  );

  const handleAddRepo = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newUrl.trim()) return;

    setAdding(true);
    setAddError('');
    try {
      const createdRepo = await addRepo(newUrl.trim());
      setNewUrl('');
      setShowAddForm(false);
      // Automatically redirect to the newly added repository board
      router.push(`/repos/${createdRepo.id}`);
      if (onCloseMobile) onCloseMobile();
    } catch (err: any) {
      setAddError(err.message || 'Error adding repository');
    } finally {
      setAdding(false);
    }
  };

  const handleRepoClick = () => {
    if (onCloseMobile) onCloseMobile();
  };

  return (
    <div className="flex flex-col h-full w-full bg-[#0c0c0e] border-r border-zinc-900 text-zinc-300 font-sans select-none">
      {/* Header Logo */}
      <div className="p-6 border-b border-zinc-900/60 flex items-center justify-between">
        <Link 
          href="/inventory" 
          onClick={handleRepoClick}
          className="flex items-center gap-2.5 text-zinc-100 hover:text-white transition-colors"
        >
          <div className="w-7 h-7 rounded-lg bg-gradient-to-tr from-cyan-500 to-blue-500 flex items-center justify-center text-black font-bold text-sm shadow-md shadow-cyan-500/10">
            IB
          </div>
          <span className="font-semibold tracking-tight text-base">IssueBoard</span>
        </Link>
      </div>

      {/* Repository Section */}
      <div className="flex-1 flex flex-col min-h-0 py-4 px-3">
        {/* Search and Add Quick Actions */}
        <div className="mb-4">
          <div className="flex items-center gap-1.5 px-2 mb-2">
            <span className="text-xs font-semibold text-zinc-500 uppercase tracking-wider">Repositories</span>
          </div>
          <div className="flex items-center gap-2">
            <div className="relative flex-1">
              <input
                type="text"
                placeholder="Search..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="w-full pl-8 pr-2 py-1.5 bg-zinc-900/40 border border-zinc-800/80 rounded-lg text-xs text-zinc-200 placeholder-zinc-600 focus:outline-none focus:border-zinc-700 focus:ring-1 focus:ring-zinc-700 transition-colors"
              />
              <svg xmlns="http://www.w3.org/2000/svg" className="h-3.5 w-3.5 text-zinc-600 absolute left-2.5 top-2.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
            <button
              onClick={() => {
                setShowAddForm(!showAddForm);
                setAddError('');
              }}
              className={`p-1.5 border rounded-lg transition-colors ${showAddForm ? 'bg-zinc-800 border-zinc-700 text-zinc-100' : 'border-zinc-800/80 hover:bg-zinc-900/60 hover:text-zinc-100 text-zinc-400'}`}
              title="Add Repository"
            >
              <svg xmlns="http://www.w3.org/2000/svg" className="h-3.5 w-3.5" viewBox="0 0 20 20" fill="currentColor">
                <path fillRule="evenodd" d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z" clipRule="evenodd" />
              </svg>
            </button>
          </div>

          {/* Quick Add Form Drawer */}
          {showAddForm && (
            <form onSubmit={handleAddRepo} className="mt-2.5 p-3 rounded-lg border border-zinc-850 bg-zinc-900/30 space-y-2 animate-fadeIn">
              <input
                type="text"
                placeholder="GitHub URL or owner/repo"
                value={newUrl}
                onChange={(e) => setNewUrl(e.target.value)}
                disabled={adding}
                required
                className="w-full px-2.5 py-1.5 bg-zinc-950 border border-zinc-850 rounded text-xs text-zinc-200 placeholder-zinc-650 focus:outline-none focus:border-zinc-700"
              />
              <div className="flex justify-end gap-1.5">
                <button
                  type="button"
                  onClick={() => setShowAddForm(false)}
                  disabled={adding}
                  className="px-2 py-1 text-[10px] text-zinc-400 hover:text-zinc-200 transition-colors"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={adding}
                  className="px-2.5 py-1 text-[10px] bg-zinc-100 hover:bg-white text-zinc-900 rounded font-medium disabled:opacity-50 transition-colors"
                >
                  {adding ? 'Adding...' : 'Add'}
                </button>
              </div>
              {addError && <p className="text-[10px] text-red-500 leading-normal">{addError}</p>}
            </form>
          )}
        </div>

        {/* Scrollable Repositories List */}
        <div className="flex-1 overflow-y-auto space-y-1 pr-1 custom-scrollbar">
          {loading ? (
            <div className="flex justify-center items-center py-8">
              <div className="animate-spin rounded-full h-4 w-4 border-2 border-zinc-800 border-t-zinc-400"></div>
            </div>
          ) : filteredRepos.length === 0 ? (
            <div className="text-center py-6 text-xs text-zinc-600">
              {search ? 'No repositories match' : 'No tracked repositories'}
            </div>
          ) : (
            filteredRepos.map((repo) => {
              const isActive = pathname === `/repos/${repo.id}`;
              return (
                <Link
                  key={repo.id}
                  href={`/repos/${repo.id}`}
                  onClick={handleRepoClick}
                  className={`group flex items-center justify-between px-3 py-2.5 rounded-lg border text-xs transition-all ${
                    isActive
                      ? 'bg-zinc-900/60 border-zinc-850 text-zinc-100 border-l-2 border-l-cyan-500 shadow-sm'
                      : 'border-transparent hover:bg-zinc-900/30 hover:text-zinc-200 text-zinc-400'
                  }`}
                >
                  <div className="min-w-0 flex-1">
                    <p className="font-medium truncate">{repo.name}</p>
                    <p className="text-[10px] text-zinc-600 truncate mt-0.5 group-hover:text-zinc-500 transition-colors">{repo.owner}</p>
                  </div>
                </Link>
              );
            })
          )}
        </div>
      </div>

      {/* Footer Profile & Settings */}
      {user && (
        <div className="p-4 border-t border-zinc-900 bg-zinc-950/40 flex flex-col gap-3">
          <div className="flex items-center gap-3 min-w-0">
            <div className="w-8 h-8 rounded-full bg-zinc-800 border border-zinc-700/50 flex items-center justify-center text-xs font-semibold text-zinc-300">
              {user.email?.[0].toUpperCase()}
            </div>
            <div className="min-w-0 flex-1">
              <p className="text-xs font-medium text-zinc-350 truncate">{user.email}</p>
            </div>
          </div>

          <div className="flex items-center gap-4 text-xs font-medium text-zinc-500 px-1">
            <button
              onClick={() => setIsSettingsOpen(true)}
              className="hover:text-zinc-300 transition-colors flex items-center gap-1.5"
            >
              <svg xmlns="http://www.w3.org/2000/svg" className="h-4.5 w-4.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
              Settings
            </button>
            <button
              onClick={signOut}
              className="hover:text-red-400 transition-colors ml-auto flex items-center gap-1.5"
            >
              Logout
            </button>
          </div>
        </div>
      )}

      {/* Settings Modal */}
      <SettingsModal isOpen={isSettingsOpen} onClose={() => setIsSettingsOpen(false)} />
    </div>
  );
}
