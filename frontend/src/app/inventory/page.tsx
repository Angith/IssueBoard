'use client';

import { useEffect, useState, useCallback } from 'react';
import { supabase } from '@/lib/supabase';
import { useRouter } from 'next/navigation';
import AddRepoForm from '@/components/AddRepoForm';
import Link from 'next/link';

interface Repository {
  id: string;
  full_name: string;
  owner: string;
  name: string;
  url: string;
}

export default function InventoryPage() {
  const [repos, setRepos] = useState<Repository[]>([]);
  const [loading, setLoading] = useState(true);
  const router = useRouter();

  const fetchRepos = useCallback(async (token: string) => {
    try {
      const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/repos`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      if (res.ok) {
        const data = await res.json();
        setRepos(data || []);
      }
    } catch (err) {
      console.error('Failed to fetch repos', err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    const checkUser = async () => {
      const { data: { session } } = await supabase.auth.getSession();
      if (!session) {
        router.push('/login');
        return;
      }
      fetchRepos(session.access_token);
    };

    checkUser();
  }, [router, fetchRepos]);

  if (loading) return <div className="p-8 text-center">Loading...</div>;

  return (
    <div className="mx-auto max-w-4xl p-8">
      <h1 className="mb-8 text-3xl font-bold">Your Repository Inventory</h1>
      
      <AddRepoForm onRepoAdded={(newRepo: Repository) => setRepos([...repos, newRepo])} />

      <div className="mt-8 grid gap-4">
        {repos.length === 0 ? (
          <p className="text-gray-500 text-center py-8">No repositories added yet.</p>
        ) : (
          repos.map((repo) => (
            <div key={repo.id} className="rounded-lg border p-4 shadow-sm hover:shadow-md transition-shadow flex justify-between items-center">
              <div>
                <h3 className="text-lg font-semibold">{repo.full_name}</h3>
                <a href={repo.url} target="_blank" rel="noopener noreferrer" className="text-sm text-blue-600 hover:underline">
                  View on GitHub
                </a>
              </div>
              <Link href={`/repos/${repo.id}`} className="rounded bg-gray-100 px-4 py-2 text-sm font-medium hover:bg-gray-200">
                View Board
              </Link>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
