import { Link } from "react-router";
import { FolderOpen, GitBranch, Plus } from "lucide-react";
import { useMe, useProjects } from "@/api/hooks";
import { Avatar } from "@/components/ui/Avatar";
import { Button } from "@/components/ui/Button";
import { SiteLink } from "@/components/ui/SiteLink";
import { Skeleton } from "@/components/ui/Skeleton";
import { TimeMeta } from "@/components/ui/TimeMeta";

export function DashboardPage() {
  const { data: user } = useMe();
  const { data: projects, isLoading, isError, error, refetch } = useProjects();

  return (
    <div>
      <div className="flex flex-col sm:flex-row sm:items-end sm:justify-between gap-4 mb-8">
        <div className="flex items-center gap-4 min-w-0">
          <Avatar
            src={user?.avatar_url}
            name={user?.name ?? "User"}
            size="lg"
          />
          <div className="min-w-0">
            <p className="text-xs uppercase tracking-[0.18em] text-[var(--accent)] mb-1">
              Dashboard
            </p>
            <h1 className="text-2xl sm:text-3xl font-semibold font-display tracking-tight truncate">
              {user?.name ? `Hey, ${user.name.split(" ")[0]}` : "Your projects"}
            </h1>
            <p className="text-sm text-[var(--text-secondary)] mt-1">
              {projects
                ? `${projects.length} project${projects.length === 1 ? "" : "s"}`
                : "Loading your sites"}
            </p>
          </div>
        </div>
        <Link to="/projects/new">
          <Button>
            <Plus size={16} />
            New Project
          </Button>
        </Link>
      </div>

      {isLoading && (
        <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
          {Array.from({ length: 3 }).map((_, i) => (
            <Skeleton key={i} className="h-48 rounded-2xl" />
          ))}
        </div>
      )}

      {isError && (
        <div className="panel p-8 text-center border-red-500/30">
          <p className="text-sm text-red-400 mb-3">
            {error instanceof Error ? error.message : "Failed to load projects"}
          </p>
          <Button variant="secondary" size="sm" onClick={() => refetch()}>
            Retry
          </Button>
        </div>
      )}

      {!isLoading && !isError && projects && projects.length === 0 && (
        <div className="panel flex flex-col items-center justify-center py-20 px-6 text-center">
          <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-cyan-400/20 to-indigo-500/20 flex items-center justify-center mb-4 border border-cyan-400/20">
            <FolderOpen size={28} className="text-[var(--accent)]" />
          </div>
          <h2 className="text-lg font-medium font-display mb-1">
            No projects yet
          </h2>
          <p className="text-sm text-[var(--text-secondary)] max-w-[340px] mb-6">
            Deploy your first static site. Connect a GitHub repo and Arthur will
            build and host it for you.
          </p>
          <Link to="/projects/new">
            <Button>
              <Plus size={16} />
              Create project
            </Button>
          </Link>
        </div>
      )}

      {!isLoading && !isError && projects && projects.length > 0 && (
        <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
          {projects.map((project) => (
            <Link
              key={project.project_id}
              to={`/projects/${project.project_id}`}
              className="group panel panel-hover p-5 relative overflow-hidden"
            >
              <div
                aria-hidden
                className="absolute inset-x-0 top-0 h-0.5 bg-gradient-to-r from-cyan-400 via-indigo-400 to-transparent opacity-70"
              />
              <h2 className="text-base font-semibold truncate mb-2 group-hover:text-[var(--accent)] transition-colors font-display">
                {project.name}
              </h2>
              <p className="text-xs text-[var(--text-tertiary)] truncate mb-4 font-mono flex items-center gap-1.5">
                <GitBranch size={12} className="flex-shrink-0" />
                {project.repo_url.replace(/^https:\/\/github\.com\//, "")}
              </p>
              <SiteLink
                slug={project.slug}
                compact
                stopPropagation
                className="mb-4"
              />
              <TimeMeta
                createdAt={project.created_at}
                updatedAt={project.updated_at}
                compact
              />
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
