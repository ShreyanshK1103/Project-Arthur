import { cn, getStatusConfig } from "@/lib/utils";

interface StatusBadgeProps {
  status: string;
  className?: string;
}

export function StatusBadge({ status, className }: StatusBadgeProps) {
  const config = getStatusConfig(status);

  return (
    <span
      className={cn(
        "inline-flex items-center gap-1.5 rounded-full border px-2.5 py-0.5 text-xs font-medium",
        config.color,
        config.bgColor,
        config.borderColor,
        className
      )}
    >
      <span
        className={cn(
          "w-1.5 h-1.5 rounded-full bg-current",
          config.pulse && "animate-pulse"
        )}
      />
      {config.label}
    </span>
  );
}
