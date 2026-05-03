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
            <Link
              key={repo.id}
              href={`/repos/${repo.id}`}
              className="p-4 border rounded hover:bg-gray-50 flex justify-between items-center"
            >
              <div>
                <h2 className="text-xl font-semibold">{repo.full_name}</h2>
                <p className="text-sm text-gray-600">{repo.url}</p>
              </div>
              <span className="text-blue-600">View Board &rarr;</span>
            </Link>
          ))
        )}
      </div>
    </div>
  );
}
