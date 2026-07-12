import { describe, expect, it } from 'vitest';
import { clientEnvelopeSchema, serverEventEnvelopeSchema } from './envelope';

describe('realtime protocol envelopes', () => {
  it('accepts the implemented command contract', () => {
    const message = clientEnvelopeSchema.parse({
      version: '1.0.0',
      messageId: 'message-1',
      command: {
        commandId: 'command-1',
        matchId: 'match-1',
        playerId: 'player-1',
        expectedRevision: 0,
        commandType: 'JOIN_ROOM',
        payload: {},
        clientTimestamp: '2026-07-11T10:00:00Z',
        clientSequence: 1,
      },
    });
    expect(message.command.commandType).toBe('JOIN_ROOM');
  });

  it('accepts a server projection and rejects unknown versions', () => {
    const data = {
      version: '1.0.0',
      messageId: 'message-1',
      matchId: 'match-1',
      sequence: 1,
      serverTime: '2026-07-11T10:00:00Z',
      type: 'match.publicState',
      payload: {},
    };
    expect(serverEventEnvelopeSchema.parse(data).type).toBe('match.publicState');
    expect(() => serverEventEnvelopeSchema.parse({ ...data, version: '2.0.0' })).toThrow();
  });
});
