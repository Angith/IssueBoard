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
    const errorData = await response.json().catch(() => ({}));
    throw new Error(errorData.error || `API Error: ${response.statusText}`);
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
