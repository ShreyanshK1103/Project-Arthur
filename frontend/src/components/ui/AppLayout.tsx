import { Outlet } from "react-router";
import { Topbar } from "./Topbar";
import { Sidebar } from "./Sidebar";

export function AppLayout() {
  return (
    <div className="min-h-screen bg-[var(--bg-primary)]">
      <Topbar />
      <div className="flex">
        <Sidebar />
        <main className="flex-1 min-w-0 p-6 lg:p-10">
          <div className="mx-auto w-full max-w-6xl">
            <Outlet />
          </div>
        </main>
      </div>
    </div>
  );
}
