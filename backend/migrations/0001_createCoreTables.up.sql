CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    handle text NOT NULL UNIQUE,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE profiles (
    user_id uuid PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    accessibility_settings jsonb NOT NULL DEFAULT '{}'::jsonb,
    statistics jsonb NOT NULL DEFAULT '{}'::jsonb,
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE rooms (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    host_user_id uuid NOT NULL REFERENCES users(id),
    status text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    closed_at timestamptz
);

CREATE TABLE matches (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id uuid NOT NULL REFERENCES rooms(id),
    scenario_id text NOT NULL,
    status text NOT NULL,
    seed_commitment text NOT NULL,
    encrypted_seed bytea,
    outcome jsonb,
    started_at timestamptz,
    ended_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE match_players (
    match_id uuid NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES users(id),
    seat smallint NOT NULL CHECK (seat BETWEEN 1 AND 10),
    encrypted_role bytea,
    result jsonb,
    PRIMARY KEY (match_id, user_id),
    UNIQUE (match_id, seat)
);

CREATE TABLE match_events (
    match_id uuid NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    sequence bigint NOT NULL CHECK (sequence > 0),
    event_type text NOT NULL,
    phase text NOT NULL,
    public_payload jsonb NOT NULL DEFAULT '{}'::jsonb,
    encrypted_private_payload bytea,
    server_time timestamptz NOT NULL,
    PRIMARY KEY (match_id, sequence)
);

CREATE INDEX match_events_time_idx ON match_events (match_id, server_time);

CREATE FUNCTION reject_match_event_mutation() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    RAISE EXCEPTION 'match_events is append-only';
END;
$$;

CREATE TRIGGER match_events_immutable
BEFORE UPDATE OR DELETE ON match_events
FOR EACH ROW EXECUTE FUNCTION reject_match_event_mutation();

CREATE TABLE match_snapshots (
    match_id uuid NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    event_sequence bigint NOT NULL,
	match_revision bigint NOT NULL CHECK (match_revision >= 0),
    phase text NOT NULL,
    compressed_state bytea NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (match_id, event_sequence),
    FOREIGN KEY (match_id, event_sequence) REFERENCES match_events(match_id, sequence)
);

CREATE TABLE role_stats (
    role_id text NOT NULL,
    season_id text NOT NULL,
    aggregate jsonb NOT NULL DEFAULT '{}'::jsonb,
    updated_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (role_id, season_id)
);

CREATE TABLE reports (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    reporter_user_id uuid NOT NULL REFERENCES users(id),
    reported_user_id uuid REFERENCES users(id),
    match_id uuid REFERENCES matches(id),
    category text NOT NULL,
    evidence_reference text,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE asset_manifests (
    version text PRIMARY KEY,
    content_hash text NOT NULL UNIQUE,
    manifest jsonb NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);
