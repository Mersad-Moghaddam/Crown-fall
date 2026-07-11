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

export const serverEventEnvelopeSchema = z.object({
  protocol_version: z.literal('1.0.0'),
  match_id: z.string().min(1),
  seq: z.number().int().positive(),
  type: z.string().min(1),
  phase: phaseSchema,
  public: z.record(z.string(), z.unknown()),
  private: z.record(z.string(), z.record(z.string(), z.unknown())).optional(),
  server_time: z.iso.datetime(),
  accepted_revision: z.number().int().nonnegative(),
});

export type ServerEventEnvelope = z.infer<typeof serverEventEnvelopeSchema>;
