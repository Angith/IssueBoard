const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export async function apiFetch(path: string, options: RequestInit = {}, token?: string) {
  const headers = new Headers(options.headers);
  if (token) {
    headers.set('Authorization', `Bearer ${token}`);
  }
  headers.set('Content-Type', 'application/json');

  const response = await fetch(`${API_URL}${path}`, {
    ...options,
    headers,
  });

  if (!response.ok) {
    let errorMessage = `API Error: ${response.statusText}`;
    try {
      const errorData = await response.json();
      errorMessage = errorData.error || errorData.message || errorMessage;
    } catch {
      const text = await response.text().catch(() => '');
      if (text) errorMessage = text.trim();
    }
    
    // Convert generic GitHub API errors
    if (response.status === 401 || errorMessage.includes('401') || errorMessage.includes('Bad credentials')) {
      errorMessage = "Your GitHub session expired or the token was revoked. Please authenticate again.";
    } else if (response.status === 403 || errorMessage.includes('403') || errorMessage.includes('API rate limit')) {
      errorMessage = "GitHub API rate limit exceeded. Please wait a moment and try again.";
    }

    throw new Error(errorMessage);
  }

  if (response.status === 204) {
    return null;
  }

  return response.json();
}

export const repositoryService = {
  list: (token: string) => apiFetch('/api/repos', {}, token),
  add: (url: string, token: string) => apiFetch('/api/repos', {
    method: 'POST',
    body: JSON.stringify({ url }),
  }, token),
  remove: (repoId: string, token: string) => apiFetch(`/api/repos/${repoId}`, {
    method: 'DELETE',
  }, token),
};

export const issueService = {
  getBoard: (repoId: string, token: string) => apiFetch(`/api/repos/${repoId}/issues`, {}, token),
  refresh: (repoId: string, token: string) => apiFetch(`/api/repos/${repoId}/refresh`, { method: 'POST' }, token),
  getAvailableLabels: (repoId: string, token: string) => apiFetch(`/api/repos/${repoId}/labels/available`, {}, token),
  getTrackedLabels: (repoId: string, token: string) => apiFetch(`/api/repos/${repoId}/labels/tracked`, {}, token),
  updateTrackedLabels: (repoId: string, labels: string[], token: string) => apiFetch(`/api/repos/${repoId}/labels/tracked`, {
    method: 'PUT',
    body: JSON.stringify({ labels }),
  }, token),
};

export const userService = {
  getSettings: (token: string) => apiFetch('/api/user/settings', {}, token),
  setGitHubToken: (githubToken: string, token: string) => apiFetch('/api/user/settings/github-token', {
    method: 'PUT',
    body: JSON.stringify({ token: githubToken }),
  }, token),
  deleteGitHubToken: (token: string) => apiFetch('/api/user/settings/github-token', {
    method: 'DELETE',
  }, token),
};

