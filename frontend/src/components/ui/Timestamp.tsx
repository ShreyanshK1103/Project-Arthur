import { cn, formatExactTime, formatRelativeTime } from "@/lib/utils";

interface TimestampProps {
  at: string;
  prefix?: string;
  className?: string;
}

/** Relative time with the exact `created_at` / `updated_at` in a tooltip. */
export function Timestamp({ at, prefix, className }: TimestampProps) {
  return (
    <time
      dateTime={at}
      title={formatExactTime(at)}
      className={cn("tabular-nums", className)}
    >
      {prefix ? `${prefix} ` : null}
      {formatRelativeTime(at)}
    </time>
  );
}
