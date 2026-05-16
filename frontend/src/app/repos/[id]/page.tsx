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

  // Label configuration state
  const [showLabelConfig, setShowLabelConfig] = useState(false);
  const [availableLabels, setAvailableLabels] = useState<any[]>([]);
  const [selectedLabels, setSelectedLabels] = useState<Set<string>>(new Set());
  const [configLoading, setConfigLoading] = useState(false);
  const [savingLabels, setSavingLabels] = useState(false);

  const fetchBoard = async () => {
    if (!session?.access_token) return;
    try {
      const data = await issueService.getBoard(id, session.access_token);
      setBoard(data);
      if (data && data.is_tracking_configured === false) {
        openLabelConfig();
      }
      setError('');
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const openLabelConfig = async () => {
    if (!session?.access_token) return;
    setShowLabelConfig(true);
    setConfigLoading(true);
    try {
      const available = await issueService.getAvailableLabels(id, session.access_token);
      const tracked = await issueService.getTrackedLabels(id, session.access_token);
      setAvailableLabels(available || []);
      setSelectedLabels(new Set(tracked || []));
    } catch (err: any) {
      setError(err.message);
    } finally {
      setConfigLoading(false);
    }
  };

  const saveLabelConfig = async () => {
    if (!session?.access_token) return;
    setSavingLabels(true);
    try {
      const labelsArray = Array.from(selectedLabels);
      await issueService.updateTrackedLabels(id, labelsArray, session.access_token);
      setShowLabelConfig(false);
      handleRefresh(); // Fetch issues based on new labels
    } catch (err: any) {
      setError(err.message);
      setSavingLabels(false);
    }
  };

  const toggleLabel = (labelName: string) => {
    const newSelected = new Set(selectedLabels);
    if (newSelected.has(labelName)) {
      newSelected.delete(labelName);
    } else {
      newSelected.add(labelName);
    }
    setSelectedLabels(newSelected);
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
    <div className="p-8 min-h-screen bg-[#09090b] text-zinc-100">
      <div className="flex justify-between items-center mb-8">
        <div>
          <h1 className="text-3xl font-semibold tracking-tight text-zinc-100">{board?.repository || 'Repository'}</h1>
          <p className="text-zinc-500 mt-1">Issues categorized by tracked labels</p>
        </div>
        <div className="flex gap-4">
          <button
            onClick={openLabelConfig}
            disabled={refreshing || loading}
            className="border border-zinc-700/50 text-zinc-300 px-4 py-2 rounded-lg hover:bg-zinc-800/50 hover:text-zinc-100 transition-colors flex items-center gap-2 text-sm font-medium"
          >
            <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M11.49 3.17c-.38-1.56-2.6-1.56-2.98 0a1.532 1.532 0 01-2.286.948c-1.372-.836-2.942.734-2.106 2.106.54.886.061 2.042-.947 2.287-1.561.379-1.561 2.6 0 2.978a1.532 1.532 0 01.947 2.287c-.836 1.372.734 2.942 2.106 2.106a1.532 1.532 0 012.287.947c.379 1.561 2.6 1.561 2.978 0a1.533 1.533 0 012.287-.947c1.372.836 2.942-.734 2.106-2.106a1.533 1.533 0 01.947-2.287c1.561-.379 1.561-2.6 0-2.978a1.532 1.532 0 01-.947-2.287c.836-1.372-.734-2.942-2.106-2.106a1.532 1.532 0 01-2.287-.947zM10 13a3 3 0 100-6 3 3 0 000 6z" clipRule="evenodd" />
            </svg>
            Configure Labels
          </button>
          <button
            onClick={handleRefresh}
            disabled={refreshing}
            className="bg-zinc-100 text-zinc-900 px-5 py-2 rounded-lg hover:bg-white disabled:opacity-50 transition-colors text-sm font-medium shadow-sm"
          >
            {refreshing ? 'Syncing Issues...' : 'Sync from GitHub'}
          </button>
        </div>
      </div>

      {!board?.is_tracking_configured && !showLabelConfig ? (
        <div className="bg-zinc-900/50 border border-zinc-800 text-zinc-300 p-6 rounded-xl text-center">
          <h2 className="text-xl font-semibold text-zinc-100 mb-2">Setup Required</h2>
          <p className="mb-4 text-zinc-400">You haven't configured which labels to track for this repository yet.</p>
          <button onClick={openLabelConfig} className="bg-zinc-100 text-zinc-900 px-6 py-2 rounded-lg hover:bg-white font-medium transition-colors text-sm shadow-sm">
            Configure Labels
          </button>
        </div>
      ) : (
        <IssueBoard categories={board?.categories || []} error={error} />
      )}

      {/* Label Configuration Modal */}
      {showLabelConfig && (
        <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center p-4 z-50">
          <div className="bg-zinc-900 rounded-xl w-full max-w-2xl max-h-[80vh] flex flex-col border border-zinc-800 shadow-2xl">
            <div className="p-6 border-b border-zinc-800/60">
              <h3 className="text-xl font-semibold tracking-tight text-zinc-100">Configure Tracked Labels</h3>
              <p className="text-zinc-400 mt-1">Select which labels you want to sync issues for.</p>
            </div>
            
            <div className="p-6 overflow-y-auto flex-1">
              {configLoading ? (
                <p className="text-zinc-400">Loading labels from GitHub...</p>
              ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                  {availableLabels.map((label) => (
                    <label 
                      key={label.name} 
                      className={`flex items-start gap-3 p-3 border rounded-lg cursor-pointer transition-colors ${selectedLabels.has(label.name) ? 'border-zinc-500 bg-zinc-800/80 ring-1 ring-zinc-500/50' : 'border-zinc-800/60 hover:bg-zinc-800/50 hover:border-zinc-700'}`}
                    >
                      <input 
                        type="checkbox" 
                        className="mt-1 accent-zinc-100"
                        checked={selectedLabels.has(label.name)}
                        onChange={() => toggleLabel(label.name)}
                      />
                      <div>
                        <div className="flex items-center gap-2">
                          <span 
                            className="w-3 h-3 rounded-full inline-block shadow-sm" 
                            style={{ backgroundColor: `#${label.color}` }}
                          ></span>
                          <span className="font-medium text-zinc-200 text-sm">{label.name}</span>
                        </div>
                        {label.description && (
                          <p className="text-xs text-zinc-500 mt-1">{label.description}</p>
                        )}
                      </div>
                    </label>
                  ))}
                  {availableLabels.length === 0 && !configLoading && (
                    <p className="text-zinc-400">No labels found in this repository.</p>
                  )}
                </div>
              )}
            </div>

            <div className="p-6 border-t border-zinc-800/60 flex justify-end gap-3 bg-zinc-900/50 rounded-b-xl">
              <button 
                onClick={() => setShowLabelConfig(false)}
                className="px-4 py-2 text-zinc-400 hover:bg-zinc-800 hover:text-zinc-200 rounded-lg text-sm font-medium transition-colors"
                disabled={savingLabels}
              >
                Cancel
              </button>
              <button 
                onClick={saveLabelConfig}
                className="px-5 py-2 bg-zinc-100 text-zinc-900 rounded-lg hover:bg-white disabled:opacity-50 text-sm font-medium transition-colors shadow-sm"
                disabled={savingLabels || configLoading}
              >
                {savingLabels ? 'Saving...' : 'Save & Sync'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
