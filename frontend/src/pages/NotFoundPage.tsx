import { Link } from "react-router";
import { ArrowLeft } from "lucide-react";

export function NotFoundPage() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-[var(--bg-primary)] px-4">
      <div className="text-center">
        <p className="text-8xl font-bold text-[var(--accent)] mb-4 tabular-nums">
          404
        </p>
        <h1 className="text-xl font-semibold text-[var(--text-primary)] mb-2">
          Page not found
        </h1>
        <p className="text-sm text-[var(--text-secondary)] mb-8 max-w-[320px] mx-auto">
          The page you're looking for doesn't exist or has been moved.
        </p>
        <Link
          to="/"
          className="inline-flex items-center gap-2 rounded-lg bg-[var(--accent)] px-4 py-2.5 text-sm font-medium text-zinc-950 hover:brightness-110 transition focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[var(--accent)]"
        >
          <ArrowLeft size={16} />
          Back home
        </Link>
      </div>
    </div>
  );
}
