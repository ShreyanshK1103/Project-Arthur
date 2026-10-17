import { forwardRef, type ButtonHTMLAttributes } from "react";
import { cn } from "@/lib/utils";

export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: "primary" | "secondary" | "ghost" | "danger" | "outline";
  size?: "sm" | "md" | "lg";
  loading?: boolean;
}

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  (
    {
      className,
      variant = "primary",
      size = "md",
      loading = false,
      disabled,
      children,
      ...props
    },
    ref
  ) => {
    return (
      <button
        ref={ref}
        disabled={disabled || loading}
        className={cn(
          "inline-flex items-center justify-center gap-2 font-medium rounded-lg transition-all",
          "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--accent)] focus-visible:ring-offset-2 focus-visible:ring-offset-[var(--bg-primary)]",
          "disabled:opacity-50 disabled:pointer-events-none",
          {
            "bg-[var(--accent)] text-zinc-950 hover:bg-[var(--accent-hover)] shadow-md shadow-cyan-500/20":
              variant === "primary",
            "bg-[var(--surface-secondary)] text-[var(--text-primary)] hover:bg-[var(--surface-hover)] border border-[var(--border-primary)]":
              variant === "secondary",
            "bg-transparent text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--surface-hover)]":
              variant === "ghost",
            "bg-red-500/10 text-red-400 hover:bg-red-500/20 border border-red-500/30":
              variant === "danger",
            "bg-transparent text-[var(--text-primary)] border border-[var(--border-secondary)] hover:bg-[var(--surface-hover)]":
              variant === "outline",
          },
          {
            "h-8 px-3 text-xs": size === "sm",
            "h-10 px-4 text-sm": size === "md",
            "h-12 px-6 text-base": size === "lg",
          },
          className
        )}
        {...props}
      >
        {loading && (
          <span
            className="w-4 h-4 border-2 border-current border-t-transparent rounded-full animate-spin"
            aria-hidden
          />
        )}
        {children}
      </button>
    );
  }
);

Button.displayName = "Button";
