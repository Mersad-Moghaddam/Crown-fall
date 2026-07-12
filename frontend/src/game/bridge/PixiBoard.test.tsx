import { render } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { PixiBoard } from './PixiBoard';
import type { PixiRuntime } from './pixiRuntime';

describe('PixiBoard lifecycle', () => {
  it('mounts once, updates on rerender, and disposes once', () => {
    const runtime: PixiRuntime = {
      mount: vi.fn(async (host) => {
        host.append(document.createElement('canvas'));
      }),
      update: vi.fn(),
      destroy: vi.fn(),
    };
    const factory = vi.fn(() => runtime);
    const view = render(
      <PixiBoard
        projection={{ phase: 'LOBBY', revision: 0 }}
        onIntention={vi.fn()}
        runtimeFactory={factory}
      />,
    );
    expect(factory).toHaveBeenCalledTimes(1);
    expect(view.getByTestId('pixi-board').querySelectorAll('canvas')).toHaveLength(1);
    view.rerender(
      <PixiBoard
        projection={{ phase: 'ROLE_DEAL', revision: 1 }}
        onIntention={vi.fn()}
        runtimeFactory={factory}
      />,
    );
    expect(factory).toHaveBeenCalledTimes(1);
    expect(runtime.update).toHaveBeenLastCalledWith({ phase: 'ROLE_DEAL', revision: 1 });
    view.unmount();
    expect(runtime.destroy).toHaveBeenCalledTimes(1);
  });

  it('creates a fresh runtime after route remount', () => {
    const instances: PixiRuntime[] = [];
    const factory = vi.fn(() => {
      const runtime: PixiRuntime = {
        mount: vi.fn(async () => undefined),
        update: vi.fn(),
        destroy: vi.fn(),
      };
      instances.push(runtime);
      return runtime;
    });
    const first = render(
      <PixiBoard
        projection={{ phase: 'LOBBY', revision: 0 }}
        onIntention={vi.fn()}
        runtimeFactory={factory}
      />,
    );
    first.unmount();
    const second = render(
      <PixiBoard
        projection={{ phase: 'LOBBY', revision: 0 }}
        onIntention={vi.fn()}
        runtimeFactory={factory}
      />,
    );
    second.unmount();
    expect(factory).toHaveBeenCalledTimes(2);
    expect(instances[0].destroy).toHaveBeenCalledOnce();
    expect(instances[1].destroy).toHaveBeenCalledOnce();
  });
});
