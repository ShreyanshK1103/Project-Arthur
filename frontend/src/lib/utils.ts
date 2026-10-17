import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";
import type { DeploymentStatus } from "@/api/types";

/** Merge Tailwind classes without conflicts. */
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

// ─── Relative time ───────────────────────────────────────────────────────────

const DIVISIONS: { amount: number; unit: Intl.RelativeTimeFormatUnit }[] = [
  { amount: 60, unit: "seconds" },
  { amount: 60, unit: "minutes" },
  { amount: 24, unit: "hours" },
  { amount: 7, unit: "days" },
  { amount: 4.345, unit: "weeks" },
  { amount: 12, unit: "months" },
  { amount: Number.POSITIVE_INFINITY, unit: "years" },
];

const rtf = new Intl.RelativeTimeFormat("en", { numeric: "auto" });

/** "2 min ago", "yesterday", etc. */
export function formatRelativeTime(dateStr: string): string {
  let seconds = (new Date(dateStr).getTime() - Date.now()) / 1000;

  for (const division of DIVISIONS) {
    if (Math.abs(seconds) < division.amount) {
      return rtf.format(Math.round(seconds), division.unit);
    }
    seconds /= division.amount;
  }

  return new Date(dateStr).toLocaleDateString();
}

/** Full ISO timestamp for tooltips. */
export function formatExactTime(dateStr: string): string {
  return new Date(dateStr).toLocaleString("en-US", {
    dateStyle: "medium",
    timeStyle: "short",
  });
}

/** Elapsed time between two backend timestamps, e.g. "1m 12s". */
export function formatDuration(startStr: string, endStr: string): string {
  const ms = Math.max(
    0,
    new Date(endStr).getTime() - new Date(startStr).getTime()
  );
  const seconds = Math.round(ms / 1000);
  if (seconds < 60) return `${seconds}s`;
  const minutes = Math.floor(seconds / 60);
  const rem = seconds % 60;
  if (minutes < 60) return rem ? `${minutes}m ${rem}s` : `${minutes}m`;
  const hours = Math.floor(minutes / 60);
  const minRem = minutes % 60;
  return minRem ? `${hours}h ${minRem}m` : `${hours}h`;
}

export function firstName(name: string): string {
  return name.trim().split(/\s+/)[0] || name;
}

export function formatRepo(url: string): string {
  return url.replace(/^https:\/\/github\.com\//i, "").replace(/\.git$/i, "");
}

// ─── Status config ───────────────────────────────────────────────────────────

export interface StatusConfig {
  label: string;
  color: string;        // Tailwind text color class
  bgColor: string;      // Tailwind bg color class
  borderColor: string;  // Tailwind border color class
  pulse: boolean;       // Whether the badge should animate
}

const STATUS_MAP: Record<DeploymentStatus, StatusConfig> = {
  queued: {
    label: "Queued",
    color: "text-zinc-400",
    bgColor: "bg-zinc-400/10",
    borderColor: "border-zinc-400/30",
    pulse: true,
  },
  building: {
    label: "Building",
    color: "text-amber-400",
    bgColor: "bg-amber-400/10",
    borderColor: "border-amber-400/30",
    pulse: true,
  },
  success: {
    label: "Ready",
    color: "text-emerald-400",
    bgColor: "bg-emerald-400/10",
    borderColor: "border-emerald-400/30",
    pulse: false,
  },
  failed: {
    label: "Failed",
    color: "text-red-400",
    bgColor: "bg-red-400/10",
    borderColor: "border-red-400/30",
    pulse: false,
  },
};

/** Returns display config for a deployment status; unknown values get a neutral fallback. */
export function getStatusConfig(status: string): StatusConfig {
  return (
    STATUS_MAP[status as DeploymentStatus] ?? {
      label: status,
      color: "text-zinc-500",
      bgColor: "bg-zinc-500/10",
      borderColor: "border-zinc-500/30",
      pulse: false,
    }
  );
}
