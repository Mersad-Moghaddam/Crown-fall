import { z } from 'zod';

const environmentSchema = z.object({
  VITE_API_URL: z.string().url(),
  VITE_WEBSOCKET_URL: z
    .string()
    .url()
    .refine((value) => /^wss?:/.test(value), 'must use ws or wss'),
});

export type Environment = z.infer<typeof environmentSchema>;

export function loadEnvironment(source: Record<string, unknown> = import.meta.env): Environment {
  return environmentSchema.parse({
    VITE_API_URL: source.VITE_API_URL ?? 'http://localhost:8080',
    VITE_WEBSOCKET_URL: source.VITE_WEBSOCKET_URL ?? 'ws://localhost:8080/ws',
  });
}

export const environment = loadEnvironment();
