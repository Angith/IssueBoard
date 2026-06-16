'use client';

import { createContext, useContext, useState, useEffect, useCallback } from 'react';
import { useAuth } from './AuthProvider';
import { repositoryService, userService } from '@/services/api';

export interface Repository {
  id: string;
  full_name: string;
  owner: string;
  name: string;
  url: string;
}

interface RepoContextType {
  repos: Repository[];
  loading: boolean;
  hasToken: boolean | null;
  error: string;
  fetchRepos: () => Promise<void>;
  addRepo: (url: string) => Promise<Repository>;
  removeRepo: (repoId: string) => Promise<void>;
}

const RepoContext = createContext<RepoContextType | undefined>(undefined);

export function RepoProvider({ children }: { children: React.ReactNode }) {
  const { session, loading: authLoading } = useAuth();
  const [repos, setRepos] = useState<Repository[]>([]);
  const [loading, setLoading] = useState(true);
  const [hasToken, setHasToken] = useState<boolean | null>(null);
  const [error, setError] = useState('');

  const fetchRepos = useCallback(async () => {
    if (!session?.access_token) return;
    try {
      const settings = await userService.getSettings(session.access_token);
      setHasToken(settings.has_github_token);
      
      if (settings.has_github_token) {
        const data = await repositoryService.list(session.access_token);
        setRepos(data || []);
      } else {
        setRepos([]);
      }
      setError('');
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
  }, [session]);

  const addRepo = useCallback(async (url: string) => {
    if (!session?.access_token) throw new Error('Not authenticated');
    try {
      const newRepo = await repositoryService.add(url, session.access_token);
      await fetchRepos();
      return newRepo;
    } catch (err: any) {
      throw new Error(err.message || 'Failed to add repository');
    }
  }, [session, fetchRepos]);

  const removeRepo = useCallback(async (repoId: string) => {
    if (!session?.access_token) throw new Error('Not authenticated');
    try {
      await repositoryService.remove(repoId, session.access_token);
      await fetchRepos();
    } catch (err: any) {
      throw new Error(err.message || 'Failed to remove repository');
    }
  }, [session, fetchRepos]);

  useEffect(() => {
    if (!authLoading && session) {
      fetchRepos();
    } else if (!authLoading && !session) {
      setRepos([]);
      setLoading(false);
    }
  }, [session, authLoading, fetchRepos]);

  return (
    <RepoContext.Provider value={{ repos, loading, hasToken, error, fetchRepos, addRepo, removeRepo }}>
      {children}
    </RepoContext.Provider>
  );
}

export const useRepos = () => {
  const context = useContext(RepoContext);
  if (context === undefined) {
    throw new Error('useRepos must be used within a RepoProvider');
  }
  return context;
};
