import { Link } from 'react-router-dom';

export default function LobbyPage() {
  return (
    <main className="page">
      <p className="eyebrow">Pre-production foundation</p>
      <h1>Crownfall</h1>
      <p>A voice-first fantasy social-deduction game for six to ten players.</p>
      <section className="panel">
        <h2>Lobby</h2>
        <p>Matchmaking is not implemented. Enter the architecture demonstration room.</p>
        <Link className="button" to="/rooms/demo-room">
          Enter room
        </Link>
      </section>
    </main>
  );
}
