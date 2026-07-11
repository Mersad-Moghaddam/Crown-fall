import { Link, useParams } from 'react-router-dom';

export default function RoomPage() {
  const { roomId } = useParams();
  return (
    <main className="page">
      <p className="eyebrow">Room {roomId}</p>
      <h1>The Empty Throne</h1>
      <div className="connection">
        <span aria-hidden="true" /> Mock connection: ready
      </div>
      <section className="panel">
        <h2>Players</h2>
        <p>1 / 10 architecture test seats occupied</p>
      </section>
      <Link className="button" to="/matches/demo-match">
        Open match table
      </Link>
    </main>
  );
}
