const base = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080/api';

export async function fetchIdeas(query = '') {
  const res = await fetch(`${base}/ideas${query ? `?${query}` : ''}`, { cache: 'no-store' });
  if (!res.ok) throw new Error('Failed to fetch ideas');
  return res.json();
}

export async function fetchIdea(id) {
  const res = await fetch(`${base}/ideas/${id}`, { cache: 'no-store' });
  if (!res.ok) throw new Error('Failed to fetch idea');
  return res.json();
}
