'use client';

import { useEffect, useState } from 'react';
import { useAuth } from '@/components/AuthProvider';
import { repositoryService, userService } from '@/services/api';
import Link from 'next/link';
import AddRepoForm from '@/components/AddRepoForm';
import SettingsModal from '@/components/SettingsModal';

export default function InventoryPage() {
  const { session, loading: authLoading } = useAuth();
  const [repos, setRepos] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [repoToDelete, setRepoToDelete] = useState<any>(null);
  const [isDeleting, setIsDeleting] = useState(false);
  const [hasToken, setHasToken] = useState<boolean | null>(null);
  const [isSettingsOpen, setIsSettingsOpen] = useState(false);

  const fetchRepos = async () => {
    if (!session?.access_token) return;
    try {
      const settings = await userService.getSettings(session.access_token);
      setHasToken(settings.has_github_token);
      
      if (settings.has_github_token) {
        const data = await repositoryService.list(session.access_token);
        setRepos(data);
      }
    } catch (err: any) {
      if (err.message?.includes('token') || err.message?.includes('Unauthorized')) {
        setError("Your GitHub session expired or the token was revoked. Please authenticate again.");
        setHasToken(false);
      } else {
        setError(err.message);
      }
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async () => {
    if (!session?.access_token || !repoToDelete) return;
    setIsDeleting(true);
    try {
      await repositoryService.remove(repoToDelete.id, session.access_token);
      setRepoToDelete(null);
      fetchRepos();
    } catch (err: any) {
      setError(err.message);
    } finally {
      setIsDeleting(false);
    }
  };

  useEffect(() => {
    if (!authLoading && session) {
      fetchRepos();
    }
  }, [session, authLoading]);

  if (authLoading || loading) return <div className="min-h-screen bg-[#09090b] p-8 text-zinc-400 font-sans flex items-center justify-center">Loading...</div>;

  return (
    <div className="min-h-screen bg-[#09090b] text-zinc-100 font-sans selection:bg-zinc-800">
      <div className="p-8 max-w-4xl mx-auto">
        <h1 className="text-2xl font-semibold mb-8 tracking-tight">My Repositories</h1>
        
        <div className="mb-8">
          {hasToken === false ? (
            <div className="border border-zinc-800 border-dashed rounded-lg p-8 text-center bg-zinc-900/30">
              <h3 className="text-lg font-medium text-zinc-200 mb-2">Connect GitHub</h3>
              <p className="text-sm text-zinc-400 mb-4 max-w-md mx-auto">
                You need to add a GitHub Personal Access Token to sync repositories and issues.
              </p>
              <button
                onClick={() => setIsSettingsOpen(true)}
                className="bg-zinc-100 text-zinc-900 hover:bg-white px-4 py-2 rounded-md text-sm font-medium transition-colors"
              >
                Add a Token
              </button>
            </div>
          ) : (
            <AddRepoForm onRepoAdded={fetchRepos} />
          )}
        </div>

        {error && <p className="text-red-500 text-sm mb-4">{error}</p>}

        <div className="flex flex-col border border-zinc-800/60 rounded-md overflow-hidden bg-[#09090b]">
          {repos.length === 0 ? (
            <p className="p-4 text-sm text-zinc-400">No repositories tracked yet. Add one above.</p>
          ) : (
            repos.map((repo, index) => (
              <div 
                key={repo.id} 
                className={`flex justify-between items-center group px-4 py-3 hover:bg-zinc-800/20 transition-colors ${
                  index !== repos.length - 1 ? 'border-b border-zinc-800/60' : ''
                }`}
              >
                <Link href={`/repos/${repo.id}`} className="flex-1 min-w-0">
                  <div className="flex items-baseline gap-2">
                    <h2 className="text-sm font-medium text-zinc-100 truncate">{repo.full_name}</h2>
                    <p className="text-xs text-zinc-500 truncate hidden sm:block">{repo.url}</p>
                  </div>
                </Link>
                <div className="flex items-center gap-3 ml-4">
                  <Link 
                    href={`/repos/${repo.id}`} 
                    className="text-xs font-medium text-zinc-400 hover:text-zinc-100 transition-colors"
                  >
                    View Board &rarr;
                  </Link>
                  <button 
                    onClick={() => setRepoToDelete(repo)}
                    className="text-zinc-600 hover:text-red-500 p-1.5 opacity-0 group-hover:opacity-100 transition-all focus:opacity-100"
                    title="Remove Repository"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                      <path fillRule="evenodd" d="M9 2a1 1 0 00-.894.553L7.382 4H4a1 1 0 000 2v10a2 2 0 002 2h8a2 2 0 002-2V6a1 1 0 100-2h-3.382l-.724-1.447A1 1 0 0011 2H9zM7 8a1 1 0 012 0v6a1 1 0 11-2 0V8zm5-1a1 1 0 00-1 1v6a1 1 0 102 0V8a1 1 0 00-1-1z" clipRule="evenodd" />
                    </svg>
                  </button>
                </div>
              </div>
            ))
          )}
        </div>

        {/* Delete Confirmation Modal */}
        {repoToDelete && (
          <div className="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center p-4 z-50">
            <div className="bg-[#121214] border border-zinc-800 rounded-lg p-6 max-w-sm w-full shadow-xl">
              <h3 className="text-lg font-medium text-zinc-100 mb-2">Remove Repository</h3>
              <p className="mb-6 text-sm text-zinc-400">Are you sure you want to remove <strong className="text-zinc-200">{repoToDelete.full_name}</strong> from your inventory?</p>
              <div className="flex justify-end gap-3">
                <button 
                  onClick={() => setRepoToDelete(null)}
                  className="px-3 py-1.5 text-sm text-zinc-400 hover:text-zinc-100 hover:bg-zinc-800 rounded transition-colors"
                  disabled={isDeleting}
                >
                  Cancel
                </button>
                <button 
                  onClick={handleDelete}
                  className="px-3 py-1.5 text-sm bg-red-500/10 text-red-500 border border-red-500/20 rounded hover:bg-red-500/20 disabled:opacity-50 transition-colors"
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
        onClose={() => {
          setIsSettingsOpen(false);
          fetchRepos();
        }} 
      />
    </div>
  );
}
