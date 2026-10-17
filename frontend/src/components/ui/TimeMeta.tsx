import { cn, formatExactTime, formatRelativeTime } from "@/lib/utils";

interface TimeMetaProps {
  createdAt: string;
  updatedAt?: string;
  className?: string;
  compact?: boolean;
}

export function TimeMeta({
  createdAt,
  updatedAt,
  className,
  compact = false,
}: TimeMetaProps) {
  const showUpdated =
    updatedAt && new Date(updatedAt).getTime() !== new Date(createdAt).getTime();

  return (
    <div
      className={cn(
        "flex flex-wrap items-center gap-x-3 gap-y-0.5 text-[var(--text-tertiary)]",
        compact ? "text-[11px]" : "text-xs",
        className
      )}
    >
      <time dateTime={createdAt} title={formatExactTime(createdAt)}>
        Created {formatRelativeTime(createdAt)}
      </time>
      {showUpdated && (
        <>
          <span aria-hidden className="text-[var(--border-secondary)]">
            ·
          </span>
          <time dateTime={updatedAt} title={formatExactTime(updatedAt)}>
            Updated {formatRelativeTime(updatedAt)}
          </time>
        </>
      )}
    </div>
  );
}
