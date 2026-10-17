import { useState } from "react";
import { Link, useLocation } from "react-router";
import { motion, AnimatePresence } from "motion/react";
import {
  LayoutDashboard,
  ChevronLeft,
  ChevronRight,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { useMe } from "@/api/hooks";
import { Avatar } from "./Avatar";

interface NavItem {
  label: string;
  path: string;
  icon: React.ReactNode;
}

const NAV_ITEMS: NavItem[] = [
  {
    label: "Dashboard",
    path: "/dashboard",
    icon: <LayoutDashboard size={20} />,
  },
];

export function Sidebar() {
  const [collapsed, setCollapsed] = useState(false);
  const location = useLocation();
  const { data: user } = useMe();

  return (
    <motion.aside
      initial={false}
      animate={{ width: collapsed ? 64 : 256 }}
      transition={{ type: "spring", stiffness: 300, damping: 30 }}
      className={cn(
        "hidden lg:flex flex-col flex-shrink-0",
        "h-[calc(100vh-56px)] sticky top-14",
        "bg-[var(--surface-primary)] border-r border-[var(--border-primary)]",
        "overflow-hidden z-20"
      )}
    >
      {/* ── Scrollable nav ── */}
      <nav className="flex-1 overflow-y-auto overflow-x-hidden py-3 px-2 scrollbar-thin">
        <ul className="space-y-1">
          {NAV_ITEMS.map((item) => {
            const active = location.pathname.startsWith(item.path);
            return (
              <li key={item.path}>
                <Link
                  to={item.path}
                  className={cn(
                    "flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors",
                    active
                      ? "bg-[var(--accent)]/10 text-[var(--accent)]"
                      : "text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--surface-hover)]"
                  )}
                >
                  <span className="flex-shrink-0">{item.icon}</span>
                  <AnimatePresence mode="wait">
                    {!collapsed && (
                      <motion.span
                        initial={{ opacity: 0, width: 0 }}
                        animate={{ opacity: 1, width: "auto" }}
                        exit={{ opacity: 0, width: 0 }}
                        transition={{ duration: 0.15 }}
                        className="truncate"
                      >
                        {item.label}
                      </motion.span>
                    )}
                  </AnimatePresence>
                </Link>
              </li>
            );
          })}
        </ul>
      </nav>

      {/* ── User info (bottom) ── */}
      <div className="border-t border-[var(--border-primary)] p-2">
        <div
          className={cn(
            "flex items-center gap-3 rounded-lg px-3 py-2",
            "text-[var(--text-secondary)]"
          )}
        >
          <Avatar
            src={user?.avatar_url}
            name={user?.name ?? "User"}
            size="sm"
          />
          <AnimatePresence mode="wait">
            {!collapsed && (
              <motion.span
                initial={{ opacity: 0, width: 0 }}
                animate={{ opacity: 1, width: "auto" }}
                exit={{ opacity: 0, width: 0 }}
                transition={{ duration: 0.15 }}
                className="truncate text-sm font-medium text-[var(--text-primary)]"
              >
                {user?.name ?? "User"}
              </motion.span>
            )}
          </AnimatePresence>
        </div>
      </div>

      {/* ── Collapse toggle ── */}
      <button
        type="button"
        onClick={() => setCollapsed(!collapsed)}
        className={cn(
          "flex items-center justify-center h-10 border-t border-[var(--border-primary)]",
          "text-[var(--text-tertiary)] hover:text-[var(--text-primary)] hover:bg-[var(--surface-hover)]",
          "transition-colors"
        )}
        aria-label={collapsed ? "Expand sidebar" : "Collapse sidebar"}
      >
        {collapsed ? <ChevronRight size={16} /> : <ChevronLeft size={16} />}
      </button>
    </motion.aside>
  );
}

// ─── Mobile sidebar (overlay) ────────────────────────────────────────────────

interface MobileSidebarProps {
  open: boolean;
  onClose: () => void;
}

export function MobileSidebar({ open, onClose }: MobileSidebarProps) {
  const location = useLocation();
  const { data: user } = useMe();

  return (
    <AnimatePresence>
      {open && (
        <>
          {/* Backdrop */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onClick={onClose}
            className="fixed inset-0 bg-black/50 z-40 lg:hidden"
          />

          {/* Panel */}
          <motion.aside
            initial={{ x: -280 }}
            animate={{ x: 0 }}
            exit={{ x: -280 }}
            transition={{ type: "spring", stiffness: 300, damping: 30 }}
            className={cn(
              "fixed left-0 top-14 bottom-0 w-[280px] z-50 lg:hidden",
              "bg-[var(--surface-primary)] border-r border-[var(--border-primary)]",
              "flex flex-col overflow-hidden"
            )}
          >
            <nav className="flex-1 overflow-y-auto py-3 px-2">
              <ul className="space-y-1">
                {NAV_ITEMS.map((item) => {
                  const active = location.pathname.startsWith(item.path);
                  return (
                    <li key={item.path}>
                      <Link
                        to={item.path}
                        onClick={onClose}
                        className={cn(
                          "flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors",
                          active
                            ? "bg-[var(--accent)]/10 text-[var(--accent)]"
                            : "text-[var(--text-secondary)] hover:text-[var(--text-primary)] hover:bg-[var(--surface-hover)]"
                        )}
                      >
                        <span className="flex-shrink-0">{item.icon}</span>
                        <span>{item.label}</span>
                      </Link>
                    </li>
                  );
                })}
              </ul>
            </nav>

            {/* User info */}
            <div className="border-t border-[var(--border-primary)] p-3">
              <div className="flex items-center gap-3 px-2">
                <Avatar
                  src={user?.avatar_url}
                  name={user?.name ?? "User"}
                  size="sm"
                />
                <span className="text-sm font-medium text-[var(--text-primary)] truncate">
                  {user?.name ?? "User"}
                </span>
              </div>
            </div>
          </motion.aside>
        </>
      )}
    </AnimatePresence>
  );
}
