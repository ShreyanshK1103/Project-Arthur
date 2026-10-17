import type { ReactNode } from "react";
import { Link } from "react-router";
import { Rocket } from "lucide-react";

interface AuthShellProps {
  title: string;
  subtitle: string;
  children: ReactNode;
  footer: ReactNode;
}

export function AuthShell({ title, subtitle, children, footer }: AuthShellProps) {
  return (
    <div className="relative min-h-screen overflow-hidden flex items-center justify-center px-4 py-12">
      <div
        aria-hidden
        className="pointer-events-none absolute inset-0"
        style={{
          backgroundImage: `
            radial-gradient(ellipse 80% 50% at 20% -10%, rgba(34,211,238,0.18), transparent 50%),
            radial-gradient(ellipse 60% 40% at 90% 10%, rgba(99,102,241,0.16), transparent 45%),
            radial-gradient(ellipse 50% 40% at 50% 110%, rgba(34,211,238,0.08), transparent 50%)
          `,
        }}
      />
      <div
        aria-hidden
        className="pointer-events-none absolute inset-0 opacity-[0.07]"
        style={{
          backgroundImage:
            "linear-gradient(var(--border-primary) 1px, transparent 1px), linear-gradient(90deg, var(--border-primary) 1px, transparent 1px)",
          backgroundSize: "48px 48px",
        }}
      />

      <div className="relative w-full max-w-md">
        <Link
          to="/"
          className="flex items-center justify-center gap-2.5 mb-8 hover:opacity-90"
        >
          <div className="w-11 h-11 rounded-xl bg-gradient-to-br from-cyan-400 to-indigo-500 flex items-center justify-center shadow-lg shadow-cyan-500/20">
            <Rocket size={22} className="text-zinc-950" />
          </div>
          <span className="text-2xl font-semibold tracking-tight font-display">
            Arthur
          </span>
        </Link>

        <h1 className="text-2xl font-semibold text-center mb-2 font-display tracking-tight">
          {title}
        </h1>
        <p className="text-sm text-[var(--text-secondary)] text-center mb-8">
          {subtitle}
        </p>

        <div className="panel p-6 sm:p-7">{children}</div>
        <div className="mt-6 text-center text-sm text-[var(--text-secondary)]">
          {footer}
        </div>
      </div>
    </div>
  );
}
