'use client';

import { useEffect, useState, use } from 'react';
import { useAuth } from '@/components/AuthProvider';
import { issueService } from '@/services/api';
import IssueBoard from '@/components/IssueBoard';

export default function RepoDetailsPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const { session, loading: authLoading } = useAuth();
  const [board, setBoard] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [error, setError] = useState('');

  const fetchBoard = async () => {
    if (!session?.access_token) return;
    try {
      const data = await issueService.getBoard(id, session.access_token);
      setBoard(data);
      setError('');
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleRefresh = async () => {
    if (!session?.access_token) return;
    setRefreshing(true);
    try {
      await issueService.refresh(id, session.access_token);
      await fetchBoard();
    } catch (err: any) {
      setError(err.message);
    } finally {
      setRefreshing(false);
    }
  };

  useEffect(() => {
    if (!authLoading && session) {
      fetchBoard();
    }
  }, [session, authLoading, id]);

  if (authLoading || (loading && !board)) return <div className="p-8">Loading board...</div>;

  return (
    <div className="p-8">
      <div className="flex justify-between items-center mb-8">
        <div>
          <h1 className="text-3xl font-bold">{board?.repository || 'Repository'}</h1>
          <p className="text-gray-600">Issues categorized by labels</p>
        </div>
        <button
          onClick={handleRefresh}
          disabled={refreshing}
          className="bg-black text-white px-6 py-2 rounded hover:bg-gray-800 disabled:opacity-50"
        >
          {refreshing ? 'Refreshing...' : 'Refresh from GitHub'}
        </button>
      </div>

      <IssueBoard categories={board?.categories || []} error={error} />
    </div>
  );
}
