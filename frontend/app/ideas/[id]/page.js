import { fetchIdea } from '../../../lib/api';

export default async function IdeaDetail({ params }) {
  const { id } = params;
  const data = await fetchIdea(id);
  return (
    <main style={{ padding: 24 }}>
      <h1>Idea Detail</h1>
      <p><strong>Idea:</strong> {data.idea.idea_text}</p>
      <p><strong>Source product:</strong> {data.idea.product_name}</p>
      <p><strong>Demand score:</strong> {data.idea.demand_score}</p>
      <p><strong>Competition:</strong> {data.idea.competition_level}</p>
      <h3>Related ideas</h3>
      <ul>
        {data.related_ideas.map((i) => <li key={i.id}>{i.idea_text}</li>)}
      </ul>
    </main>
  );
}
