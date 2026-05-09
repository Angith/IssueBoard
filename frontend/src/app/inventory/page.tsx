'use client';

import { useEffect, useState } from 'react';
import { useAuth } from '@/components/AuthProvider';
import { repositoryService } from '@/services/api';
import Link from 'next/link';
import AddRepoForm from '@/components/AddRepoForm';

export default function InventoryPage() {
  const { session, loading: authLoading } = useAuth();
  const [repos, setRepos] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [repoToDelete, setRepoToDelete] = useState<any>(null);
  const [isDeleting, setIsDeleting] = useState(false);

  const fetchRepos = async () => {
    if (!session?.access_token) return;
    try {
      const data = await repositoryService.list(session.access_token);
      setRepos(data);
    } catch (err: any) {
      setError(err.message);
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

  if (authLoading || loading) return <div className="p-8">Loading...</div>;

  return (
    <div className="p-8 max-w-4xl mx-auto">
      <h1 className="text-3xl font-bold mb-8">My Repositories</h1>
      
      <div className="mb-8">
        <AddRepoForm onRepoAdded={fetchRepos} />
      </div>

      {error && <p className="text-red-600 mb-4">{error}</p>}

      <div className="grid gap-4">
        {repos.length === 0 ? (
          <p>No repositories tracked yet. Add one above!</p>
        ) : (
          repos.map((repo) => (
            <div key={repo.id} className="p-4 border rounded hover:bg-gray-50 flex justify-between items-center group">
              <Link href={`/repos/${repo.id}`} className="flex-1">
                <div>
                  <h2 className="text-xl font-semibold">{repo.full_name}</h2>
                  <p className="text-sm text-gray-600">{repo.url}</p>
                </div>
              </Link>
              <div className="flex items-center gap-4">
                <Link href={`/repos/${repo.id}`} className="text-blue-600">
                  View Board &rarr;
                </Link>
                <button 
                  onClick={() => setRepoToDelete(repo)}
                  className="text-red-500 hover:text-red-700 p-2"
                  title="Remove Repository"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
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
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-lg p-6 max-w-sm w-full">
            <h3 className="text-xl font-bold mb-4">Remove Repository</h3>
            <p className="mb-6">Are you sure you want to remove <strong>{repoToDelete.full_name}</strong> from your inventory?</p>
            <div className="flex justify-end gap-4">
              <button 
                onClick={() => setRepoToDelete(null)}
                className="px-4 py-2 text-gray-600 hover:bg-gray-100 rounded"
                disabled={isDeleting}
              >
                Cancel
              </button>
              <button 
                onClick={handleDelete}
                className="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700 disabled:opacity-50"
                disabled={isDeleting}
              >
                {isDeleting ? 'Removing...' : 'Remove'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
