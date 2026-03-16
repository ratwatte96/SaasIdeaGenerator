import Link from 'next/link';
import Filters from '../components/Filters';
import { fetchIdeas } from '../lib/api';

export default async function Home({ searchParams }) {
  const q = new URLSearchParams(searchParams).toString();
  let data = { data: [] };
  let error = null;
  try { data = await fetchIdeas(q); } catch (e) { error = e.message; }

  return (
    <main style={{ display: 'grid', gridTemplateColumns: '280px 1fr', minHeight: '100vh' }}>
      <aside style={{ background: '#fff', padding: 16, borderRight: '1px solid #ddd' }}>
        <h3>Filters</h3>
        <Filters />
      </aside>
      <section style={{ padding: 16 }}>
        <h1>Idea List</h1>
        {error && <p>{error}</p>}
        {!error && data.data.length === 0 && <p>No ideas yet.</p>}
        {!error && data.data.length > 0 && (
          <table border='1' cellPadding='8' style={{ borderCollapse: 'collapse', width: '100%', background: '#fff' }}>
            <thead><tr><th>Idea</th><th>Demand Score</th><th>Competition</th><th>Source Product</th></tr></thead>
            <tbody>
              {data.data.map((row) => (
                <tr key={row.id}>
                  <td><Link href={`/ideas/${row.id}`}>{row.idea_text}</Link></td>
                  <td>{row.demand_score}</td>
                  <td>{row.competition_level}</td>
                  <td>{row.product_name}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </section>
    </main>
  );
}
