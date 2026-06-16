'use client';

import { useState } from 'react';
import { useRepos } from '@/components/RepoProvider';
import Link from 'next/link';
import AddRepoForm from '@/components/AddRepoForm';
import SettingsModal from '@/components/SettingsModal';

export default function InventoryPage() {
  const { repos, loading, hasToken, error, fetchRepos, removeRepo } = useRepos();
  const [repoToDelete, setRepoToDelete] = useState<any>(null);
  const [isDeleting, setIsDeleting] = useState(false);
  const [isSettingsOpen, setIsSettingsOpen] = useState(false);

  const handleDelete = async () => {
    if (!repoToDelete) return;
    setIsDeleting(true);
    try {
      await removeRepo(repoToDelete.id);
      setRepoToDelete(null);
    } catch (err: any) {
      // Error handled by repo context or set locally
    } finally {
      setIsDeleting(false);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-[#09090b] p-8 text-zinc-400 font-sans flex items-center justify-center">
        <div className="flex flex-col items-center gap-3">
          <div className="animate-spin rounded-full h-8 w-8 border-2 border-zinc-800 border-t-cyan-500"></div>
          <span className="text-sm">Loading dashboard...</span>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-[#09090b] text-zinc-100 font-sans selection:bg-zinc-800">
      <div className="p-8 max-w-5xl mx-auto space-y-10">
        
        {/* Welcome Section */}
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 border-b border-zinc-900/60 pb-6">
          <div>
            <h1 className="text-3xl font-semibold tracking-tight text-zinc-100 font-sans">Workspace Dashboard</h1>
            <p className="text-zinc-500 mt-1.5 text-sm">
              Manage your repositories and view your issue boards.
            </p>
          </div>
          {hasToken !== false && (
            <div className="flex items-center gap-3">
              <button
                onClick={() => setIsSettingsOpen(true)}
                className="px-4 py-2 border border-zinc-800 hover:border-zinc-700 bg-zinc-900/30 text-zinc-300 hover:text-zinc-100 rounded-lg text-sm font-medium transition-colors"
              >
                GitHub Connection Settings
              </button>
            </div>
          )}
        </div>

        {/* Global Connection Alerts */}
        {hasToken === false && (
          <div className="border border-yellow-500/20 rounded-xl p-6 bg-yellow-500/5 flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
            <div className="space-y-1">
              <h3 className="text-base font-semibold text-yellow-550 flex items-center gap-2">
                <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                  <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
                </svg>
                GitHub Token Required
              </h3>
              <p className="text-sm text-zinc-400 font-sans leading-relaxed">
                You need to configure a GitHub Personal Access Token to sync and view your tracked repositories and issues.
              </p>
            </div>
            <button
              onClick={() => setIsSettingsOpen(true)}
              className="bg-zinc-100 text-zinc-900 hover:bg-white px-4 py-2 rounded-lg text-sm font-medium transition-colors shadow-sm"
            >
              Add GitHub Token
            </button>
          </div>
        )}

        {/* Form and Overview Cards */}
        {hasToken !== false && (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            
            {/* Quick Add Form widget */}
            <div className="md:col-span-2 border border-zinc-850/60 rounded-xl p-6 bg-zinc-900/10 space-y-4">
              <div>
                <h3 className="text-sm font-semibold text-zinc-200">Track New Repository</h3>
                <p className="text-xs text-zinc-500 mt-1">Add a GitHub repository URL to sync its issues immediately.</p>
              </div>
              <AddRepoForm onRepoAdded={fetchRepos} />
            </div>

            {/* Quick Stats Summary widget */}
            <div className="border border-zinc-850/60 rounded-xl p-6 bg-zinc-900/10 flex flex-col justify-between">
              <div>
                <h3 className="text-sm font-semibold text-zinc-200">Workspace Summary</h3>
                <p className="text-xs text-zinc-500 mt-1">Snapshot of your currently tracked project boards.</p>
              </div>
              <div className="flex items-baseline gap-2 mt-4">
                <span className="text-3xl font-bold text-zinc-100">{repos.length}</span>
                <span className="text-xs text-zinc-500">Tracked Repositories</span>
              </div>
            </div>

          </div>
        )}

        {error && <p className="text-red-500 text-sm">{error}</p>}

        {/* Repository Grid Cards */}
        {hasToken !== false && (
          <div className="space-y-4">
            <h2 className="text-lg font-semibold text-zinc-250">Connected Project Boards</h2>
            {repos.length === 0 ? (
              <div className="border border-dashed border-zinc-850 rounded-xl p-12 text-center bg-zinc-900/10">
                <p className="text-sm text-zinc-500 mb-2 font-sans">No repositories added to this workspace yet.</p>
                <p className="text-xs text-zinc-600">Paste a GitHub URL above to get started.</p>
              </div>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {repos.map((repo) => (
                  <div 
                    key={repo.id} 
                    className="group relative flex flex-col justify-between p-5 border border-zinc-850 bg-zinc-900/20 hover:border-zinc-750 hover:bg-zinc-900/40 rounded-xl transition-all duration-200"
                  >
                    <div className="space-y-1.5 pr-8">
                      <div className="flex items-center gap-2">
                        <svg className="h-4 w-4 text-zinc-500" viewBox="0 0 16 16" fill="currentColor">
                          <path d="M2 2.5A2.5 2.5 0 014.5 0h8.75a.75.75 0 01.75.75v12.5a.75.75 0 01-.75.75h-2.5a.75.75 0 110-1.5h1.75v-2h-8a1 1 0 00-.714 1.7.75.75 0 01-1.072 1.05A2.495 2.495 0 012 11.5v-9zm10.5-1V9h-8c-.356 0-.694.074-1 .208V2.5a1 1 0 011-1h8z"></path>
                        </svg>
                        <h3 className="font-semibold text-zinc-200 text-sm truncate">{repo.name}</h3>
                      </div>
                      <p className="text-xs text-zinc-500 truncate">{repo.owner}</p>
                      <p className="text-[10px] text-zinc-600 truncate mt-1">{repo.url}</p>
                    </div>

                    <div className="flex items-center justify-between mt-6 pt-3 border-t border-zinc-900">
                      <Link 
                        href={`/repos/${repo.id}`} 
                        className="text-xs font-semibold text-cyan-400 hover:text-cyan-300 transition-colors flex items-center gap-1"
                      >
                        View Issue Board 
                        <span className="transform translate-x-0 group-hover:translate-x-1 transition-transform">&rarr;</span>
                      </Link>
                      
                      <button 
                        onClick={() => setRepoToDelete(repo)}
                        className="text-zinc-650 hover:text-red-500 p-1.5 transition-colors focus:opacity-100"
                        title="Remove Repository"
                      >
                        <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                          <path fillRule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clipRule="evenodd" />
                        </svg>
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}

        {/* Delete Confirmation Modal */}
        {repoToDelete && (
          <div className="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center p-4 z-50">
            <div className="bg-[#121214] border border-zinc-800 rounded-xl p-6 max-w-sm w-full shadow-2xl animate-scaleIn">
              <h3 className="text-lg font-medium text-zinc-100 mb-2">Remove Repository</h3>
              <p className="mb-6 text-sm text-zinc-400 font-sans leading-relaxed">
                Are you sure you want to remove <strong className="text-zinc-200">{repoToDelete.full_name}</strong> from your inventory?
              </p>
              <div className="flex justify-end gap-3">
                <button 
                  onClick={() => setRepoToDelete(null)}
                  className="px-3.5 py-1.5 text-sm text-zinc-455 hover:text-zinc-100 hover:bg-zinc-805 rounded-lg transition-colors"
                  disabled={isDeleting}
                >
                  Cancel
                </button>
                <button 
                  onClick={handleDelete}
                  className="px-4 py-1.5 text-sm bg-red-500/10 text-red-500 border border-red-500/20 rounded-lg hover:bg-red-500/20 disabled:opacity-50 transition-colors"
                  disabled={isDeleting}
                >
                  {isDeleting ? 'Removing...' : 'Remove'}
                </button>
              </div>
            </div>
          </div>
        )}
      </div>

      <SettingsModal 
        isOpen={isSettingsOpen} 
        onClose={() => setIsSettingsOpen(false)} 
      />
    </div>
  );
}
