import { z } from 'zod';

export const phaseSchema = z.enum([
  'LOBBY',
  'ROLE_DEAL',
  'CHAPTER_START',
  'CRISIS_REVEAL',
  'NOMINATION',
  'COUNCIL_VOTE',
  'PARTY_BUILD',
  'QUEST_COMMIT',
  'QUEST_REVEAL',
  'EXECUTIVE_POWER',
  'TRIAL_WINDOW',
  'ROUND_END',
  'FINAL_RECKONING',
  'EPILOGUE',
  'PAUSED_RECONNECT',
]);

export const commandSchema = z.object({
  commandId: z.string().min(1),
  matchId: z.string().min(1),
  playerId: z.string().min(1),
  expectedRevision: z.number().int().nonnegative(),
  commandType: z.enum(['JOIN_ROOM', 'SET_READY', 'START_MATCH', 'ACKNOWLEDGE_ROLE']),
  payload: z.record(z.string(), z.unknown()).default({}),
  clientTimestamp: z.iso.datetime(),
  clientSequence: z.number().int().positive(),
});

export const clientEnvelopeSchema = z.object({
  version: z.literal('1.0.0'),
  messageId: z.string().min(1),
  command: commandSchema,
});

export const serverEventEnvelopeSchema = z.object({
  version: z.literal('1.0.0'),
  messageId: z.string().min(1),
  matchId: z.string().min(1),
  sequence: z.number().int().nonnegative(),
  serverTime: z.iso.datetime(),
  type: z.enum([
    'connection.accepted',
    'connection.resynced',
    'match.publicState',
    'match.privateState',
    'match.privateEvents',
    'protocol.error',
  ]),
  payload: z.unknown(),
});

export type ClientEnvelope = z.infer<typeof clientEnvelopeSchema>;
export type ServerEventEnvelope = z.infer<typeof serverEventEnvelopeSchema>;
