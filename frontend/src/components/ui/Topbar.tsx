import { useState } from "react";
import { Link } from "react-router";
import { Menu, LogOut, Rocket } from "lucide-react";
import { toast } from "sonner";
import { cn } from "@/lib/utils";
import { useMe, useLogout } from "@/api/hooks";
import { Avatar } from "./Avatar";
import { MobileSidebar } from "./Sidebar";

export function Topbar() {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const { data: user } = useMe();
  const logout = useLogout();
  const [userMenuOpen, setUserMenuOpen] = useState(false);

  return (
    <>
      <header
        className={cn(
          "sticky top-0 z-30 h-14",
          "flex items-center justify-between px-4 lg:px-6",
          "bg-[var(--surface-primary)]/80 backdrop-blur-xl",
          "border-b border-[var(--border-primary)]"
        )}
      >
        {/* ── Left: hamburger + logo ── */}
        <div className="flex items-center gap-3">
          <button
            type="button"
            onClick={() => setMobileMenuOpen(true)}
            className="lg:hidden p-1.5 rounded-md text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--surface-hover)] transition-colors"
            aria-label="Open menu"
          >
            <Menu size={20} />
          </button>

          <Link
            to="/dashboard"
            className="flex items-center gap-2.5 text-[var(--text-primary)] hover:opacity-80 transition-opacity"
          >
            <div className="w-7 h-7 rounded-lg bg-gradient-to-br from-cyan-400 to-indigo-500 flex items-center justify-center">
              <Rocket size={16} className="text-zinc-950" />
            </div>
            <span className="text-base font-semibold tracking-tight font-display">
              Arthur
            </span>
          </Link>
        </div>

        {/* ── Right: user menu ── */}
        <div className="relative">
          <button
            type="button"
            onClick={() => setUserMenuOpen(!userMenuOpen)}
            className={cn(
              "flex items-center gap-2 rounded-full p-1 pr-3",
              "hover:bg-[var(--surface-hover)] transition-colors"
            )}
          >
            <Avatar
              src={user?.avatar_url}
              name={user?.name ?? "User"}
              size="sm"
            />
            <span className="hidden sm:block text-sm text-[var(--text-primary)] font-medium truncate max-w-[120px]">
              {user?.name ?? "User"}
            </span>
          </button>

          {/* Dropdown */}
          {userMenuOpen && (
            <>
              <div
                className="fixed inset-0 z-40"
                onClick={() => setUserMenuOpen(false)}
              />
              <div
                className={cn(
                  "absolute right-0 top-full mt-2 z-50",
                  "w-52 rounded-xl overflow-hidden",
                  "bg-[var(--surface-secondary)] border border-[var(--border-primary)]",
                  "shadow-xl shadow-black/20"
                )}
              >
                <div className="px-4 py-3 border-b border-[var(--border-primary)] flex items-center gap-3">
                  <Avatar
                    src={user?.avatar_url}
                    name={user?.name ?? "User"}
                    size="md"
                  />
                  <div className="min-w-0">
                    <p className="text-sm font-medium text-[var(--text-primary)] truncate">
                      {user?.name}
                    </p>
                    <p className="text-xs text-[var(--text-tertiary)] truncate">
                      {user?.email}
                    </p>
                  </div>
                </div>
                <div className="p-1.5">
                  <button
                    type="button"
                    onClick={() => {
                      setUserMenuOpen(false);
                      logout.mutate(undefined, {
                        onSuccess: () => toast.success("Signed out"),
                        onError: () => toast.error("Could not sign out"),
                      });
                    }}
                    className={cn(
                      "w-full flex items-center gap-2 rounded-lg px-3 py-2 text-sm",
                      "text-red-400 hover:bg-red-400/10 transition-colors"
                    )}
                  >
                    <LogOut size={16} />
                    Sign out
                  </button>
                </div>
              </div>
            </>
          )}
        </div>
      </header>

      {/* Mobile sidebar */}
      <MobileSidebar
        open={mobileMenuOpen}
        onClose={() => setMobileMenuOpen(false)}
      />
    </>
  );
}
