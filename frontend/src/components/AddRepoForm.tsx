'use client';

import { useState } from 'react';
import { supabase } from '@/lib/supabase';

interface Repository {
  id: string;
  full_name: string;
  owner: string;
  name: string;
  url: string;
}

interface Props {
  onRepoAdded: (repo: Repository) => void;
}

export default function AddRepoForm({ onRepoAdded }: Props) {
  const [url, setUrl] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    setError('');

    try {
      const { data: { session } } = await supabase.auth.getSession();
      if (!session) throw new Error('Not authenticated');

      const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/repos`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${session.access_token}`,
        },
        body: JSON.stringify({ url }),
      });

      if (!res.ok) {
        const msg = await res.text();
        throw new Error(msg || 'Failed to add repository');
      }

      const newRepo = await res.json();
      onRepoAdded(newRepo);
      setUrl('');
    } catch (err: unknown) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError('An unexpected error occurred');
      }
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="w-full">
      <div className="flex items-center gap-3 rounded-md border border-zinc-800/60 bg-[#09090b] px-3 py-1.5 focus-within:border-zinc-700 focus-within:ring-1 focus-within:ring-zinc-700 transition-all">
        <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 text-zinc-500" viewBox="0 0 20 20" fill="currentColor">
          <path fillRule="evenodd" d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z" clipRule="evenodd" />
        </svg>
        <input
          type="text"
          placeholder="Add new repository (e.g., https://github.com/owner/repo)"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          required
          className="flex-1 bg-transparent border-none text-sm text-zinc-100 placeholder-zinc-500 focus:outline-none focus:ring-0 p-0"
        />
        <button
          type="submit"
          disabled={submitting}
          className="rounded bg-zinc-100 px-3 py-1 text-xs font-medium text-black hover:bg-white disabled:opacity-50 transition-colors"
        >
          {submitting ? 'Adding...' : 'Add'}
        </button>
      </div>
      {error && <p className="mt-2 text-sm text-red-500">{error}</p>}
    </form>
  );
}
