import { lazy, Suspense } from "react";
import { Link } from "react-router";
import {
  GitBranch,
  Hammer,
  Rocket,
  RotateCcw,
  Shield,
  Terminal,
  ArrowRight,
} from "lucide-react";
import { Button } from "@/components/ui/Button";
import { useMe } from "@/api/hooks";

const PixelGrid = lazy(() =>
  import("@/components/vendor/PixelGrid").then((m) => ({ default: m.PixelGrid }))
);

const STEPS = [
  {
    icon: GitBranch,
    title: "Connect",
    body: "Point Arthur at a public GitHub repo over HTTPS.",
  },
  {
    icon: Hammer,
    title: "Build",
    body: "Sandboxed builds run with your install and build commands.",
  },
  {
    icon: Rocket,
    title: "Deploy",
    body: "Get a stable slug URL like project-arthur-abc12.localhost:8080, plus a permalink for every deploy.",
  },
  {
    icon: RotateCcw,
    title: "Rollback",
    body: "Promote any successful deployment instantly — no rebuild.",
  },
] as const;

const FEATURES = [
  {
    icon: Shield,
    title: "Encrypted env vars",
    body: "Build-time variables stored encrypted. Remember: static bundles are public.",
  },
  {
    icon: Terminal,
    title: "Live build logs",
    body: "Stream ANSI-colored logs while a deployment is queued or building.",
  },
  {
    icon: Rocket,
    title: "Stable URLs",
    body: "Projects live at slug.localhost:8080. Each deployment is the first 12 characters of its id at that same host.",
  },
] as const;

export function LandingPage() {
  const { data: user } = useMe();

  return (
    <div className="min-h-screen bg-[var(--bg-primary)] text-[var(--text-primary)]">
      <header className="relative z-10 flex items-center justify-between px-6 lg:px-10 h-16">
        <div className="flex items-center gap-2.5">
          <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-cyan-400 to-indigo-500 flex items-center justify-center">
            <Rocket size={18} className="text-zinc-950" />
          </div>
          <span className="text-lg font-semibold tracking-tight font-display">Arthur</span>
        </div>
        <div className="flex items-center gap-2">
          {user ? (
            <Link to="/dashboard">
              <Button size="sm">Dashboard</Button>
            </Link>
          ) : (
            <>
              <Link to="/login">
                <Button variant="ghost" size="sm">
                  Sign in
                </Button>
              </Link>
              <Link to="/register">
                <Button size="sm">Get started</Button>
              </Link>
            </>
          )}
        </div>
      </header>

      <section className="relative overflow-hidden border-b border-[var(--border-primary)] min-h-[520px] flex flex-col">
        <div className="absolute inset-0 opacity-25 pointer-events-none">
          <Suspense fallback={null}>
            <PixelGrid />
          </Suspense>
        </div>
        <div className="absolute inset-0 pointer-events-none bg-gradient-to-b from-[var(--bg-primary)]/20 via-transparent to-[var(--bg-primary)]" />

        <div className="relative z-10 max-w-4xl mx-auto px-6 pt-24 pb-28 text-center flex-1 flex flex-col items-center justify-center">
          <p className="text-sm font-medium text-cyan-300 mb-4 tracking-wide">
            Static site deployments
          </p>
          <h1 className="text-5xl sm:text-6xl lg:text-7xl font-semibold tracking-tight mb-5 font-display text-white">
            Arthur
          </h1>
          <p className="text-lg text-zinc-300 max-w-xl mx-auto mb-10">
            Push a repo. Watch the build. Ship a live URL — with one-click
            rollback when you need it.
          </p>
          <div className="flex flex-wrap items-center justify-center gap-3">
            <Link to={user ? "/dashboard" : "/register"}>
              <Button size="lg">
                {user ? "Open dashboard" : "Start deploying"}
                <ArrowRight size={18} />
              </Button>
            </Link>
            {!user && (
              <Link to="/login">
                <Button size="lg" variant="outline">
                  Sign in
                </Button>
              </Link>
            )}
          </div>
        </div>
      </section>

      <section className="max-w-5xl mx-auto px-6 py-20">
        <h2 className="text-2xl font-semibold text-center mb-3">How it works</h2>
        <p className="text-center text-[var(--text-secondary)] mb-12 max-w-lg mx-auto">
          Four steps from repo to live site.
        </p>
        <ol className="grid sm:grid-cols-2 lg:grid-cols-4 gap-4">
          {STEPS.map((step, i) => (
            <li key={step.title} className="panel p-5">
              <div className="w-10 h-10 rounded-xl bg-[var(--accent)]/10 text-[var(--accent)] flex items-center justify-center mb-4">
                <step.icon size={20} />
              </div>
              <p className="text-xs text-[var(--text-tertiary)] mb-1">
                Step {i + 1}
              </p>
              <h3 className="font-semibold mb-2 font-display">{step.title}</h3>
              <p className="text-sm text-[var(--text-secondary)]">{step.body}</p>
            </li>
          ))}
        </ol>
      </section>

      <section className="border-t border-[var(--border-primary)] bg-[var(--bg-secondary)]">
        <div className="max-w-5xl mx-auto px-6 py-20">
          <h2 className="text-2xl font-semibold text-center mb-3">
            Built for real deploys
          </h2>
          <p className="text-center text-[var(--text-secondary)] mb-12 max-w-lg mx-auto">
            Features that match what the Arthur backend actually ships today.
          </p>
          <div className="grid md:grid-cols-3 gap-4">
            {FEATURES.map((f) => (
              <div key={f.title} className="panel p-5">
                <div className="w-10 h-10 rounded-xl bg-[var(--accent)]/10 flex items-center justify-center mb-4 text-[var(--accent)]">
                  <f.icon size={20} />
                </div>
                <h3 className="font-semibold mb-2 font-display">{f.title}</h3>
                <p className="text-sm text-[var(--text-secondary)]">{f.body}</p>
              </div>
            ))}
          </div>
          <p className="mt-12 text-center text-xs text-[var(--text-tertiary)] max-w-2xl mx-auto">
            Static sites only. No custom domains or PR preview deployments yet.
            Builds may queue on a single worker under load.
          </p>
        </div>
      </section>

      <footer className="border-t border-[var(--border-primary)] px-6 py-8 flex flex-col sm:flex-row items-center justify-between gap-4 text-sm text-[var(--text-tertiary)]">
        <div className="flex items-center gap-2">
          <Rocket size={16} className="text-[var(--accent)]" />
          Arthur
        </div>
        <p>Deploy static sites from GitHub.</p>
      </footer>
    </div>
  );
}
