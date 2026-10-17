import { useState } from "react";
import { cn } from "@/lib/utils";

interface AvatarProps {
  src?: string | null;
  name: string;
  size?: "sm" | "md" | "lg";
  className?: string;
}

const SIZE = {
  sm: "w-7 h-7 text-xs",
  md: "w-9 h-9 text-sm",
  lg: "w-12 h-12 text-base",
} as const;

function initials(name: string): string {
  const parts = name.trim().split(/\s+/).filter(Boolean);
  if (parts.length === 0) return "?";
  if (parts.length === 1) return parts[0].slice(0, 2).toUpperCase();
  return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase();
}

/** GitHub avatars 403 without `referrerPolicy="no-referrer"`. */
export function Avatar({ src, name, size = "sm", className }: AvatarProps) {
  const [failed, setFailed] = useState(false);
  const showImage = Boolean(src) && !failed;

  return (
    <span
      className={cn(
        "relative inline-flex flex-shrink-0 items-center justify-center overflow-hidden rounded-full",
        size === "sm"
          ? "ring-1 ring-[var(--accent)]/30"
          : "ring-2 ring-[var(--accent)]/25 ring-offset-2 ring-offset-[var(--bg-primary)]",
        SIZE[size],
        className
      )}
    >
      {showImage ? (
        <img
          src={src!}
          alt={name}
          referrerPolicy="no-referrer"
          onError={() => setFailed(true)}
          className="h-full w-full object-cover"
        />
      ) : (
        <span className="flex h-full w-full items-center justify-center bg-gradient-to-br from-cyan-400 to-indigo-500 font-semibold text-zinc-950">
          {initials(name)}
        </span>
      )}
    </span>
  );
}
