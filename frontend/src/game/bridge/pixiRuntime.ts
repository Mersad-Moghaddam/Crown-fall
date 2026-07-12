import { Application, Graphics, Text } from 'pixi.js';

export type BoardProjection = Readonly<{ phase: string; revision: number }>;
export type PlayerIntention = Readonly<{ type: 'BOARD_SELECTED'; payload: { targetId: string } }>;

export interface PixiRuntime {
  mount(
    host: HTMLDivElement,
    projection: BoardProjection,
    emit: (intention: PlayerIntention) => void,
  ): Promise<void>;
  update(projection: BoardProjection): void;
  destroy(): void;
}

export type PixiRuntimeFactory = () => PixiRuntime;

export function createPixiRuntime(): PixiRuntime {
  const application = new Application();
  let label: Text | undefined;
  let destroyed = false;
  return {
    async mount(host, projection) {
      await application.init({ background: '#171025', resizeTo: host, antialias: true });
      if (destroyed) return;
      host.replaceChildren(application.canvas);
      const table = new Graphics()
        .circle(400, 220, 180)
        .fill({ color: '#34264a' })
        .stroke({ color: '#c7a667', width: 4 });
      label = new Text({
        text: `Authoritative projection · ${projection.phase}`,
        style: { fill: '#f3ead8', fontFamily: 'Georgia', fontSize: 22 },
      });
      label.anchor.set(0.5);
      label.position.set(400, 220);
      application.stage.addChild(table, label);
    },
    update(projection) {
      if (label) label.text = `Authoritative projection · ${projection.phase}`;
    },
    destroy() {
      if (destroyed) return;
      destroyed = true;
      application.destroy(true, { children: true, texture: true, textureSource: true });
    },
  };
}
