import { describe, expect, it } from 'vitest';
import { loadEnvironment } from './environment';

describe('runtime environment', () => {
  it('validates HTTP and WebSocket endpoints', () => {
    expect(
      loadEnvironment({
        VITE_API_URL: 'https://api.example.test',
        VITE_WEBSOCKET_URL: 'wss://api.example.test/ws',
      }),
    ).toBeTruthy();
    expect(() =>
      loadEnvironment({ VITE_API_URL: 'invalid', VITE_WEBSOCKET_URL: 'https://example.test/ws' }),
    ).toThrow();
  });
});
