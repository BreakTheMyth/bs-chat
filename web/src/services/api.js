export async function fetchRecommendedSessions() {
  const resp = await fetch('/sessions/list');
  if (!resp.ok) throw new Error('Failed to fetch sessions');
  return resp.json();
}

export async function fetchSessionDetail(id) {
  const resp = await fetch(`/sessions/${id}`);
  if (!resp.ok) throw new Error('Failed to fetch session');
  return resp.json();
}

export async function fetchUser(id) {
  const resp = await fetch(`/users/${id}`);
  if (!resp.ok) throw new Error('Failed to fetch user');
  return resp.json();
}
