import { Navigate, Outlet } from "react-router";
import { useMe } from "@/api/hooks";
import { Loader2 } from "lucide-react";

/**
 * Route guard for public-only pages (login, register).
 * If already logged in → redirect to /dashboard.
 */
export function PublicRoute() {
  const { data: user, isLoading, isError } = useMe();

  if (isLoading && !isError) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-[var(--bg-primary)]">
        <Loader2
          size={32}
          className="animate-spin text-[var(--accent)]"
        />
      </div>
    );
  }

  if (user) {
    return <Navigate to="/dashboard" replace />;
  }

  return <Outlet />;
}
