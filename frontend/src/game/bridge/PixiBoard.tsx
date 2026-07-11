import { Application, Graphics, Text } from 'pixi.js';
import { useEffect, useRef } from 'react';

export function PixiBoard({ phase }: { phase: string }) {
  const host = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const application = new Application();
    let disposed = false;
    void application
      .init({ background: '#171025', resizeTo: host.current!, antialias: true })
      .then(() => {
        if (disposed) return;
        host.current?.appendChild(application.canvas);
        const table = new Graphics()
          .circle(400, 220, 180)
          .fill({ color: '#34264a' })
          .stroke({ color: '#c7a667', width: 4 });
        const label = new Text({
          text: `Authoritative projection · ${phase}`,
          style: { fill: '#f3ead8', fontFamily: 'Georgia', fontSize: 22 },
        });
        label.anchor.set(0.5);
        label.position.set(400, 220);
        application.stage.addChild(table, label);
      });
    return () => {
      disposed = true;
      application.destroy(true, { children: true });
    };
  }, [phase]);

  return (
    <div className="board" ref={host} aria-label={`Animated Crownfall table in ${phase} phase`} />
  );
}
