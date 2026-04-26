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
    <form onSubmit={handleSubmit} className="rounded-lg border bg-gray-50 p-6">
      <h2 className="mb-4 text-xl font-semibold">Add New Repository</h2>
      <div className="flex gap-2">
        <input
          type="text"
          placeholder="https://github.com/owner/repo"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          required
          className="flex-1 rounded border px-4 py-2 text-black"
        />
        <button
          type="submit"
          disabled={submitting}
          className="rounded bg-black px-6 py-2 text-white hover:bg-gray-800 disabled:opacity-50"
        >
          {submitting ? 'Adding...' : 'Add'}
        </button>
      </div>
      {error && <p className="mt-2 text-sm text-red-600">{error}</p>}
    </form>
  );
}
