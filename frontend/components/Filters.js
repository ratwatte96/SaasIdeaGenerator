'use client';

import { useRouter, useSearchParams } from 'next/navigation';

export default function Filters() {
  const router = useRouter();
  const params = useSearchParams();

  const update = (key, value) => {
    const next = new URLSearchParams(params.toString());
    if (!value) next.delete(key); else next.set(key, value);
    router.push(`/?${next.toString()}`);
  };

  return (
    <div style={{ display: 'grid', gap: 8 }}>
      <label>Category <input defaultValue={params.get('category') || ''} onBlur={(e) => update('category', e.target.value)} /></label>
      <label>Min Demand <input type='number' defaultValue={params.get('min_demand_score') || ''} onBlur={(e) => update('min_demand_score', e.target.value)} /></label>
      <label>Competition
        <select defaultValue={params.get('competition_level') || ''} onChange={(e) => update('competition_level', e.target.value)}>
          <option value=''>All</option><option value='low'>low</option><option value='medium'>medium</option><option value='high'>high</option>
        </select>
      </label>
    </div>
  );
}
