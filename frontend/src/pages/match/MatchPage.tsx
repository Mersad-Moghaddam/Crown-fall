import { useParams } from 'react-router-dom';
import { PixiBoard } from '../../game/bridge/PixiBoard';

export default function MatchPage() {
  const { matchId } = useParams();
  return (
    <main className="matchPage">
      <header>
        <div>
          <p className="eyebrow">Match {matchId}</p>
          <h1>Council Table</h1>
        </div>
        <div className="connection">
          <span aria-hidden="true" /> Mock connection: synchronized
        </div>
      </header>
      <PixiBoard projection={{ phase: 'LOBBY', revision: 0 }} onIntention={() => undefined} />
      <aside className="panel">
        <h2>Authoritative phase</h2>
        <strong>LOBBY</strong>
        <p>The board renders a projection. It cannot validate or advance the match.</p>
      </aside>
    </main>
  );
}
