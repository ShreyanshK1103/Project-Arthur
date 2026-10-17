import { Navigate, Outlet } from "react-router";
import { useMe } from "@/api/hooks";
import { Loader2 } from "lucide-react";

/**
 * Route guard: renders children only when authenticated.
 * While loading → spinner. If unauthenticated → redirect to /login.
 */
export function ProtectedRoute() {
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

  if (!user) {
    return <Navigate to="/login" replace />;
  }

  return <Outlet />;
}
