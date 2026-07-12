import { useEffect, useRef } from 'react';
import {
  createPixiRuntime,
  type BoardProjection,
  type PixiRuntimeFactory,
  type PlayerIntention,
} from './pixiRuntime';

export type PixiBoardProps = {
  projection: BoardProjection;
  onIntention: (intention: PlayerIntention) => void;
  runtimeFactory?: PixiRuntimeFactory;
};

export function PixiBoard({
  projection,
  onIntention,
  runtimeFactory = createPixiRuntime,
}: PixiBoardProps) {
  const host = useRef<HTMLDivElement>(null);
  const runtime = useRef<ReturnType<PixiRuntimeFactory> | null>(null);
  const projectionRef = useRef(projection);
  const intentionRef = useRef(onIntention);

  useEffect(() => {
    projectionRef.current = projection;
  }, [projection]);
  useEffect(() => {
    intentionRef.current = onIntention;
  }, [onIntention]);

  useEffect(() => {
    const instance = runtimeFactory();
    const hostElement = host.current!;
    runtime.current = instance;
    void instance.mount(hostElement, projectionRef.current, (intention) =>
      intentionRef.current(intention),
    );
    return () => {
      instance.destroy();
      runtime.current = null;
      hostElement.replaceChildren();
    };
  }, [runtimeFactory]);

  useEffect(() => {
    runtime.current?.update(projection);
  }, [projection]);

  return (
    <div
      className="board"
      ref={host}
      data-testid="pixi-board"
      aria-label={`Animated Crownfall table in ${projection.phase} phase`}
    />
  );
}
