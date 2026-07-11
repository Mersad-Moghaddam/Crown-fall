import { describe, expect, it } from 'vitest';
import { serverEventEnvelopeSchema } from './envelope';

describe('server event envelope', () => {
  it('accepts a public authoritative event', () => {
    const event = serverEventEnvelopeSchema.parse({
      protocol_version: '1.0.0',
      match_id: 'match-1',
      seq: 1,
      type: 'match.stateUpdated',
      phase: 'LOBBY',
      public: {},
      server_time: '2026-07-11T10:00:00Z',
      accepted_revision: 0,
    });
    expect(event.phase).toBe('LOBBY');
  });

  it('rejects non-normative phases', () => {
    expect(() =>
      serverEventEnvelopeSchema.parse({
        protocol_version: '1.0.0',
        match_id: 'match-1',
        seq: 1,
        type: 'match.stateUpdated',
        phase: 'EPILOGUE_REMATCH',
        public: {},
        server_time: '2026-07-11T10:00:00Z',
        accepted_revision: 0,
      }),
    ).toThrow();
  });
});
