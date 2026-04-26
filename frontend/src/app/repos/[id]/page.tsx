'use client';

import { useEffect, useState, useCallback } from 'react';
import { supabase } from '@/lib/supabase';
import { useRouter, useParams } from 'next/navigation';
import IssueBoard from '@/components/IssueBoard';

interface Label {
  name: string;
  color: string;
}

interface Issue {
  id: string;
  number: number;
  title: string;
  url: string;
  state: string;
}

interface IssueCategory {
  label: Label;
  issues: Issue[];
}

interface Board {
  repository: string;
  categories: IssueCategory[];
}

export default function RepoBoardPage() {
  const { id } = useParams();
  const [board, setBoard] = useState<Board | null>(null);
  const [loading, setLoading] = useState(true);
  const router = useRouter();

  const fetchBoard = useCallback(async (token: string) => {
    try {
      const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/repos/${id}/issues`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      if (res.ok) {
        const data = await res.json();
        setBoard(data);
      }
    } catch (err) {
      console.error('Failed to fetch board', err);
    } finally {
      setLoading(false);
    }
  }, [id]);

  useEffect(() => {
    const checkUser = async () => {
      const { data: { session } } = await supabase.auth.getSession();
      if (!session) {
        router.push('/login');
        return;
      }
      fetchBoard(session.access_token);
    };

    checkUser();
  }, [id, router, fetchBoard]);

  const handleRefresh = async () => {
    setLoading(true);
    try {
      const { data: { session } } = await supabase.auth.getSession();
      if (!session) return;

      const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/repos/${id}/refresh`, {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${session.access_token}`,
        },
      });
      if (res.ok) {
        fetchBoard(session.access_token);
      }
    } catch (err) {
      console.error('Failed to refresh board', err);
      setLoading(false);
    }
  };

  if (loading) return <div className="p-8 text-center">Loading Board...</div>;
  if (!board) return <div className="p-8 text-center">Board not found.</div>;

  return (
    <div className="p-8">
      <div className="mb-8 flex items-center justify-between">
        <h1 className="text-3xl font-bold">Issue Board</h1>
        <button
          onClick={handleRefresh}
          className="rounded bg-blue-600 px-4 py-2 text-white hover:bg-blue-700"
        >
          Refresh from GitHub
        </button>
      </div>
      
      <IssueBoard categories={board.categories} />
    </div>
  );
}
