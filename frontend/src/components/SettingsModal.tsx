import { useState, useEffect } from 'react';
import { userService } from '../services/api';
import { useAuth } from './AuthProvider';
import { useRepos } from './RepoProvider';

interface SettingsModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export default function SettingsModal({ isOpen, onClose }: SettingsModalProps) {
  const { session } = useAuth();
  const { fetchRepos } = useRepos();
  const [tokenInput, setTokenInput] = useState('');
  const [hasToken, setHasToken] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  useEffect(() => {
    if (isOpen && session?.access_token) {
      loadSettings();
    }
  }, [isOpen, session]);

  const loadSettings = async () => {
    try {
      const data = await userService.getSettings(session!.access_token);
      setHasToken(data.has_github_token);
    } catch (err) {
      console.error('Failed to load settings', err);
    }
  };

  const handleUpdate = async () => {
    if (!tokenInput.trim()) {
      setError('Please enter a valid GitHub token');
      return;
    }

    setIsLoading(true);
    setError('');
    setSuccess('');

    try {
      await userService.setGitHubToken(tokenInput, session!.access_token);
      setSuccess('Token updated successfully.');
      setHasToken(true);
      setTokenInput('');
      await fetchRepos();
    } catch (err: any) {
      setError(err.message || 'Failed to update token.');
    } finally {
      setIsLoading(false);
    }
  };

  const handleRemove = async () => {
    setIsLoading(true);
    setError('');
    setSuccess('');

    try {
      await userService.deleteGitHubToken(session!.access_token);
      setSuccess('Token removed successfully.');
      setHasToken(false);
      setTokenInput('');
      await fetchRepos();
    } catch (err: any) {
      setError(err.message || 'Failed to remove token.');
    } finally {
      setIsLoading(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm">
      <div className="w-full max-w-md rounded-xl border border-zinc-800 bg-[#09090b] p-6 shadow-2xl">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-xl font-semibold text-zinc-100">Settings</h2>
          <button
            onClick={onClose}
            className="text-zinc-400 hover:text-zinc-100 transition-colors"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg>
          </button>
        </div>

        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-zinc-300 mb-1">
              GitHub Personal Access Token
            </label>
            <p className="text-xs text-zinc-500 mb-3">
              Required to fetch repositories and issues from GitHub.
              {hasToken && <span className="ml-2 text-emerald-500 font-medium">✓ Token is active</span>}
            </p>
            
            <input
              type="password"
              placeholder="ghp_..."
              value={tokenInput}
              onChange={(e) => setTokenInput(e.target.value)}
              className="w-full rounded-md border border-zinc-800 bg-zinc-900 px-3 py-2 text-sm text-zinc-100 placeholder-zinc-500 focus:border-zinc-700 focus:outline-none focus:ring-1 focus:ring-zinc-700"
            />
          </div>

          {error && <p className="text-sm text-red-500">{error}</p>}
          {success && <p className="text-sm text-emerald-500">{success}</p>}

          <div className="flex gap-3 pt-2">
            <button
              onClick={handleUpdate}
              disabled={isLoading || !tokenInput}
              className="flex-1 rounded-md bg-zinc-100 px-4 py-2 text-sm font-medium text-zinc-900 transition-colors hover:bg-white disabled:opacity-50 disabled:cursor-not-allowed flex justify-center items-center h-10"
            >
              {isLoading && tokenInput ? (
                <div className="w-4 h-4 border-2 border-zinc-900 border-t-transparent rounded-full animate-spin"></div>
              ) : (
                'Update Token'
              )}
            </button>
            {hasToken && (
              <button
                onClick={handleRemove}
                disabled={isLoading}
                className="flex-1 rounded-md border border-zinc-800 bg-transparent px-4 py-2 text-sm font-medium text-zinc-300 transition-colors hover:bg-zinc-900 hover:text-white disabled:opacity-50 disabled:cursor-not-allowed flex justify-center items-center h-10"
              >
                {isLoading && !tokenInput ? (
                  <div className="w-4 h-4 border-2 border-zinc-300 border-t-transparent rounded-full animate-spin"></div>
                ) : (
                  'Remove Token'
                )}
              </button>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
