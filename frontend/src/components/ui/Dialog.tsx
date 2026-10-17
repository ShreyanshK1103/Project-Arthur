import { useEffect, type ReactNode } from "react";
import { X } from "lucide-react";
import { cn } from "@/lib/utils";
import { Button } from "./Button";

interface DialogProps {
  open: boolean;
  onClose: () => void;
  title: string;
  description?: string;
  children?: ReactNode;
  confirmLabel?: string;
  confirmVariant?: "primary" | "danger";
  onConfirm?: () => void;
  loading?: boolean;
}

export function Dialog({
  open,
  onClose,
  title,
  description,
  children,
  confirmLabel = "Confirm",
  confirmVariant = "primary",
  onConfirm,
  loading,
}: DialogProps) {
  useEffect(() => {
    if (!open) return;
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") onClose();
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [open, onClose]);

  if (!open) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div
        className="absolute inset-0 bg-black/60"
        onClick={onClose}
        aria-hidden
      />
      <div
        role="dialog"
        aria-modal="true"
        aria-labelledby="dialog-title"
        className={cn(
          "relative z-10 w-full max-w-md rounded-xl",
          "bg-[var(--surface-primary)] border border-[var(--border-primary)]",
          "shadow-2xl shadow-black/40 p-6"
        )}
      >
        <div className="flex items-start justify-between gap-4 mb-4">
          <div>
            <h2
              id="dialog-title"
              className="text-lg font-semibold text-[var(--text-primary)]"
            >
              {title}
            </h2>
            {description && (
              <p className="mt-1 text-sm text-[var(--text-secondary)]">
                {description}
              </p>
            )}
          </div>
          <button
            type="button"
            onClick={onClose}
            className="p-1 rounded-md text-[var(--text-tertiary)] hover:text-[var(--text-primary)] hover:bg-[var(--surface-hover)]"
            aria-label="Close"
          >
            <X size={18} />
          </button>
        </div>

        {children && <div className="mb-6">{children}</div>}

        {onConfirm && (
          <div className="flex justify-end gap-2">
            <Button variant="ghost" onClick={onClose} disabled={loading}>
              Cancel
            </Button>
            <Button
              variant={confirmVariant}
              onClick={onConfirm}
              loading={loading}
            >
              {confirmLabel}
            </Button>
          </div>
        )}
      </div>
    </div>
  );
}
