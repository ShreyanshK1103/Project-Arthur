import { useEffect, useRef, useState } from "react";
import { ArrowDown, Copy } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/Button";
import type { LogLine } from "@/api/types";
import { ansiToHtml } from "@/lib/ansi";
import { cn } from "@/lib/utils";

function logTime(iso: string): string {
  return new Date(iso).toLocaleTimeString("en-US", {
    hour12: false,
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  });
}

interface LogViewerProps {
  logs: LogLine[];
  done: boolean;
  error: string | null;
}

export function LogViewer({ logs, done, error }: LogViewerProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const [stickToBottom, setStickToBottom] = useState(true);

  useEffect(() => {
    if (!stickToBottom || !containerRef.current) return;
    containerRef.current.scrollTop = containerRef.current.scrollHeight;
  }, [logs, stickToBottom]);

  const onScroll = () => {
    const el = containerRef.current;
    if (!el) return;
    const nearBottom = el.scrollHeight - el.scrollTop - el.clientHeight < 48;
    setStickToBottom(nearBottom);
  };

  const copyAll = async () => {
    try {
      await navigator.clipboard.writeText(logs.map((l) => l.message).join("\n"));
      toast.success("Logs copied");
    } catch {
      toast.error("Could not copy logs");
    }
  };

  return (
    <div className="relative rounded-xl border border-[var(--border-primary)] bg-[#0a0a0c] overflow-hidden">
      <div className="flex items-center justify-between px-3 py-2 border-b border-[var(--border-primary)] bg-[var(--surface-primary)]">
        <div className="flex items-center gap-2 text-xs text-[var(--text-secondary)]">
          <span
            className={cn(
              "w-1.5 h-1.5 rounded-full",
              done ? "bg-emerald-400" : "bg-amber-400 animate-pulse"
            )}
          />
          {done ? "Build finished" : "Streaming logs…"}
          {error && <span className="text-red-400">· {error}</span>}
        </div>
        <Button type="button" variant="ghost" size="sm" onClick={copyAll}>
          <Copy size={14} />
          Copy
        </Button>
      </div>

      <div
        ref={containerRef}
        onScroll={onScroll}
        className="h-[480px] overflow-auto font-mono text-xs leading-5 p-3"
      >
        {logs.length === 0 ? (
          <p className="text-[var(--text-tertiary)]">Waiting for log output…</p>
        ) : (
          logs.map((line) => {
            const html = ansiToHtml(line.message);
            return (
              <div key={line.seq} className="flex gap-3 hover:bg-white/5">
                <span className="select-none text-[var(--text-tertiary)] w-10 text-right flex-shrink-0">
                  {line.seq}
                </span>
                <span
                  className="select-none text-[var(--text-tertiary)]/80 w-[68px] flex-shrink-0"
                  title={new Date(line.created_at).toLocaleString()}
                >
                  {logTime(line.created_at)}
                </span>
                <span
                  className="text-[var(--text-secondary)] whitespace-pre-wrap break-all"
                  dangerouslySetInnerHTML={{ __html: html }}
                />
              </div>
            );
          })
        )}
      </div>

      {!stickToBottom && (
        <button
          type="button"
          onClick={() => {
            setStickToBottom(true);
            if (containerRef.current) {
              containerRef.current.scrollTop =
                containerRef.current.scrollHeight;
            }
          }}
          className="absolute bottom-4 right-4 inline-flex items-center gap-1.5 rounded-lg bg-[var(--accent)] text-zinc-950 text-xs px-3 py-1.5 shadow-lg hover:bg-[var(--accent-hover)]"
        >
          <ArrowDown size={14} />
          Jump to latest
        </button>
      )}
    </div>
  );
}
