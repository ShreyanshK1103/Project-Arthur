import type { MouseEvent } from "react";
import { ExternalLink } from "lucide-react";
import { CopyButton } from "@/components/ui/CopyButton";
import { siteHost, siteHref } from "@/lib/siteUrl";
import { cn } from "@/lib/utils";

interface SiteLinkProps {
  slug: string;
  className?: string;
  /** Hide extra padding on tight table cells. */
  compact?: boolean;
  /** Use when nested inside a parent <Link> / card. */
  stopPropagation?: boolean;
}

export function SiteLink({
  slug,
  className,
  compact = false,
  stopPropagation = false,
}: SiteLinkProps) {
  const host = siteHost(slug);
  const href = siteHref(slug);

  const isolate = stopPropagation
    ? (e: MouseEvent) => e.stopPropagation()
    : undefined;

  const isolateCopy = stopPropagation
    ? (e: MouseEvent) => {
        e.preventDefault();
        e.stopPropagation();
      }
    : undefined;

  return (
    <span className={cn("inline-flex items-center gap-1 min-w-0", className)}>
      <a
        href={href}
        target="_blank"
        rel="noreferrer"
        onClick={isolate}
        className={cn(
          "inline-flex items-center gap-1.5 min-w-0 text-[var(--accent)] hover:underline font-mono",
          compact ? "text-xs" : "text-sm"
        )}
      >
        <span className="truncate">{host}</span>
        <ExternalLink size={14} className="flex-shrink-0" />
      </a>
      <span onClick={isolateCopy} className="flex-shrink-0">
        <CopyButton
          value={href}
          className={compact ? "px-1" : undefined}
        />
      </span>
    </span>
  );
}
