/**
 * Lightweight pixel-grid background inspired by Componentry Pixel Canvas.
 * Source pattern: https://componentry.dev/docs/components/pixel-canvas
 * Copied/adapted locally so auth/landing stay dependency-light.
 */
import { useEffect, useRef } from "react";
import { cn } from "@/lib/utils";

interface PixelGridProps {
  className?: string;
  colors?: string[];
  gap?: number;
}

export function PixelGrid({
  className,
  colors = ["#06b6d4", "#22d3ee", "#67e8f9", "#0891b2"],
  gap = 8,
}: PixelGridProps) {
  const canvasRef = useRef<HTMLCanvasElement>(null);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    const reduced = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
    const ctx = canvas.getContext("2d");
    if (!ctx) return;

    let raf = 0;
    let width = 0;
    let height = 0;
    let cols = 0;
    let rows = 0;
    const heat: number[] = [];
    let mx = -999;
    let my = -999;

    const resize = () => {
      const parent = canvas.parentElement;
      if (!parent) return;
      width = parent.clientWidth;
      height = parent.clientHeight;
      canvas.width = width * devicePixelRatio;
      canvas.height = height * devicePixelRatio;
      canvas.style.width = `${width}px`;
      canvas.style.height = `${height}px`;
      ctx.setTransform(devicePixelRatio, 0, 0, devicePixelRatio, 0, 0);
      cols = Math.ceil(width / gap);
      rows = Math.ceil(height / gap);
      heat.length = cols * rows;
      heat.fill(0);
    };

    const onMove = (e: PointerEvent) => {
      const rect = canvas.getBoundingClientRect();
      mx = e.clientX - rect.left;
      my = e.clientY - rect.top;
    };

    const draw = () => {
      ctx.clearRect(0, 0, width, height);
      for (let y = 0; y < rows; y++) {
        for (let x = 0; x < cols; x++) {
          const i = y * cols + x;
          const cx = x * gap + gap / 2;
          const cy = y * gap + gap / 2;
          if (!reduced) {
            const dx = cx - mx;
            const dy = cy - my;
            const dist = Math.hypot(dx, dy);
            if (dist < 80) {
              heat[i] = Math.min(1, heat[i] + (1 - dist / 80) * 0.35);
            } else {
              heat[i] *= 0.92;
            }
          } else {
            heat[i] = ((x + y) % 7 === 0) ? 0.15 : 0.04;
          }
          const h = heat[i];
          if (h < 0.02) continue;
          const color = colors[Math.floor(h * (colors.length - 1))];
          ctx.globalAlpha = h;
          ctx.fillStyle = color;
          ctx.fillRect(x * gap + 1, y * gap + 1, gap - 2, gap - 2);
        }
      }
      ctx.globalAlpha = 1;
      if (!reduced) raf = requestAnimationFrame(draw);
    };

    resize();
    draw();
    window.addEventListener("resize", resize);
    canvas.addEventListener("pointermove", onMove);

    return () => {
      cancelAnimationFrame(raf);
      window.removeEventListener("resize", resize);
      canvas.removeEventListener("pointermove", onMove);
    };
  }, [colors, gap]);

  return (
    <canvas
      ref={canvasRef}
      className={cn("absolute inset-0 w-full h-full", className)}
      aria-hidden
    />
  );
}
