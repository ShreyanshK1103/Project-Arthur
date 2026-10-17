import { useState } from "react";
import { Link, useNavigate, useParams } from "react-router";
import {
  RefreshCw,
  Trash2,
  ArrowLeft,
  GitBranch,
} from "lucide-react";
import { toast } from "sonner";
import {
  useDeleteProject,
  useProject,
  useProjectDeployments,
  usePromote,
  useRedeploy,
} from "@/api/hooks";
import { ApiRequestError } from "@/api/client";
import { Button } from "@/components/ui/Button";
import { Dialog } from "@/components/ui/Dialog";
import { SiteLink } from "@/components/ui/SiteLink";
import { Skeleton } from "@/components/ui/Skeleton";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { TimeMeta } from "@/components/ui/TimeMeta";
import { EnvVarsPanel } from "@/features/env/EnvVarsPanel";
import { deploymentSlug } from "@/lib/siteUrl";
import { cn, formatExactTime, formatRelativeTime } from "@/lib/utils";

type Tab = "deployments" | "environment" | "settings";

export function ProjectDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [tab, setTab] = useState<Tab>("deployments");
  const [promoteId, setPromoteId] = useState<string | null>(null);
  const [confirmDelete, setConfirmDelete] = useState(false);

  const projectQuery = useProject(id);
  const deploymentsQuery = useProjectDeployments(id);
  const redeploy = useRedeploy();
  const promote = usePromote();
  const deleteProject = useDeleteProject();

  if (projectQuery.isLoading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-8 w-64" />
        <Skeleton className="h-4 w-96" />
        <Skeleton className="h-64" />
      </div>
    );
  }

  if (projectQuery.isError || !projectQuery.data) {
    return (
      <div className="text-center py-16">
        <p className="text-red-400 mb-4">
          {projectQuery.error instanceof Error
            ? projectQuery.error.message
            : "Project not found"}
        </p>
        <Link to="/dashboard">
          <Button variant="secondary">Back to dashboard</Button>
        </Link>
      </div>
    );
  }

  const project = projectQuery.data;

  return (
    <div>
      <Link
        to="/dashboard"
        className="inline-flex items-center gap-1.5 text-sm text-[var(--text-secondary)] hover:text-[var(--text-primary)] mb-4"
      >
        <ArrowLeft size={16} />
        Projects
      </Link>

      <div className="panel p-5 sm:p-6 mb-6">
        <div className="flex flex-col lg:flex-row lg:items-start lg:justify-between gap-4">
          <div className="min-w-0">
            <h1 className="text-2xl sm:text-3xl font-semibold font-display tracking-tight truncate">
              {project.name}
            </h1>
            <p className="text-sm text-[var(--text-tertiary)] font-mono truncate mt-2 flex items-center gap-1.5">
              <GitBranch size={14} className="flex-shrink-0" />
              {project.repo_url}
            </p>
            <TimeMeta
              createdAt={project.created_at}
              updatedAt={project.updated_at}
              className="mt-3"
            />
            <div className="flex flex-col gap-1.5 mt-4">
              <SiteLink slug={project.slug} />
              {project.active_deployment_id && (
                <p className="text-xs text-[var(--text-tertiary)]">
                  Active permalink:{" "}
                  <SiteLink
                    slug={deploymentSlug(project.active_deployment_id)}
                    compact
                  />
                </p>
              )}
            </div>
          </div>
          <Button
            onClick={async () => {
              try {
                const d = await redeploy.mutateAsync(project.project_id);
                toast.success("Redeploy queued");
                navigate(
                  `/projects/${project.project_id}/deployments/${d.id}`
                );
              } catch (err) {
                toast.error(
                  err instanceof ApiRequestError
                    ? err.message
                    : "Redeploy failed"
                );
              }
            }}
            loading={redeploy.isPending}
          >
            <RefreshCw size={16} />
            Redeploy
          </Button>
        </div>
      </div>

      <div className="flex gap-1 border-b border-[var(--border-primary)] mb-6 overflow-x-auto">
        {(
          [
            ["deployments", "Deployments"],
            ["environment", "Environment"],
            ["settings", "Settings"],
          ] as const
        ).map(([key, label]) => (
          <button
            key={key}
            type="button"
            onClick={() => setTab(key)}
            className={cn(
              "px-4 py-2.5 text-sm font-medium border-b-2 -mb-px transition-colors whitespace-nowrap",
              tab === key
                ? "border-[var(--accent)] text-[var(--accent)]"
                : "border-transparent text-[var(--text-secondary)] hover:text-[var(--text-primary)]"
            )}
          >
            {label}
          </button>
        ))}
      </div>

      {tab === "deployments" && (
        <div>
          {deploymentsQuery.isLoading && (
            <div className="space-y-2">
              <Skeleton className="h-12" />
              <Skeleton className="h-12" />
            </div>
          )}
          {deploymentsQuery.isError && (
            <p className="text-sm text-red-400">
              {deploymentsQuery.error instanceof Error
                ? deploymentsQuery.error.message
                : "Failed to load deployments"}
            </p>
          )}
          {deploymentsQuery.data && deploymentsQuery.data.length === 0 && (
            <p className="text-sm text-[var(--text-tertiary)]">
              No deployments yet. Click Redeploy to start a build.
            </p>
          )}
          {deploymentsQuery.data && deploymentsQuery.data.length > 0 && (
            <div className="overflow-x-auto rounded-xl border border-[var(--border-primary)]">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-[var(--border-primary)] bg-[var(--surface-secondary)] text-left text-[var(--text-tertiary)]">
                    <th className="px-4 py-3 font-medium">Status</th>
                    <th className="px-4 py-3 font-medium">Created</th>
                    <th className="px-4 py-3 font-medium">Updated</th>
                    <th className="px-4 py-3 font-medium">URL</th>
                    <th className="px-4 py-3 font-medium text-right">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {deploymentsQuery.data.map((d) => (
                    <tr
                      key={d.id}
                      className="border-b border-[var(--border-primary)] last:border-0 hover:bg-[var(--surface-hover)]/50"
                    >
                      <td className="px-4 py-3">
                        <div className="flex items-center gap-2">
                          <StatusBadge status={d.status} />
                          {d.is_active && (
                            <span className="text-[10px] uppercase tracking-wide text-emerald-400 border border-emerald-400/30 rounded px-1.5 py-0.5">
                              Active
                            </span>
                          )}
                        </div>
                      </td>
                      <td className="px-4 py-3 text-[var(--text-secondary)]">
                        <span title={formatExactTime(d.created_at)}>
                          {formatRelativeTime(d.created_at)}
                        </span>
                      </td>
                      <td className="px-4 py-3 text-[var(--text-secondary)]">
                        <span title={formatExactTime(d.updated_at)}>
                          {formatRelativeTime(d.updated_at)}
                        </span>
                      </td>
                      <td className="px-4 py-3">
                        {d.status === "success" ? (
                          <SiteLink
                            slug={deploymentSlug(d.id)}
                            compact
                            className="max-w-[280px]"
                          />
                        ) : (
                          <span className="text-[var(--text-tertiary)]">—</span>
                        )}
                      </td>
                      <td className="px-4 py-3 text-right space-x-2">
                        <Link
                          to={`/projects/${project.project_id}/deployments/${d.id}`}
                          className="text-[var(--accent)] hover:underline text-xs"
                        >
                          Logs
                        </Link>
                        {d.status === "success" && !d.is_active && (
                          <button
                            type="button"
                            onClick={() => setPromoteId(d.id)}
                            className="text-xs text-[var(--text-secondary)] hover:text-[var(--text-primary)]"
                          >
                            Promote
                          </button>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      )}

      {tab === "environment" && <EnvVarsPanel projectId={project.project_id} />}

      {tab === "settings" && (
        <div className="max-w-xl space-y-6">
          <div className="panel p-5 space-y-3">
            <h2 className="text-sm font-medium text-[var(--text-primary)]">
              Build settings
            </h2>
            <dl className="grid grid-cols-[120px_1fr] gap-y-2 text-sm">
              <dt className="text-[var(--text-tertiary)]">Branch</dt>
              <dd className="font-mono text-[var(--text-secondary)]">
                {project.branch}
              </dd>
              <dt className="text-[var(--text-tertiary)]">Install</dt>
              <dd className="font-mono text-[var(--text-secondary)]">
                {project.install_command}
              </dd>
              <dt className="text-[var(--text-tertiary)]">Build</dt>
              <dd className="font-mono text-[var(--text-secondary)]">
                {project.build_command}
              </dd>
              <dt className="text-[var(--text-tertiary)]">Output</dt>
              <dd className="font-mono text-[var(--text-secondary)]">
                {project.output_dir}
              </dd>
              <dt className="text-[var(--text-tertiary)]">Live URL</dt>
              <dd>
                <SiteLink slug={project.slug} className="text-sm" />
              </dd>
              <dt className="text-[var(--text-tertiary)]">Created</dt>
              <dd
                className="text-[var(--text-secondary)]"
                title={formatExactTime(project.created_at)}
              >
                {formatRelativeTime(project.created_at)}
              </dd>
              <dt className="text-[var(--text-tertiary)]">Updated</dt>
              <dd
                className="text-[var(--text-secondary)]"
                title={formatExactTime(project.updated_at)}
              >
                {formatRelativeTime(project.updated_at)}
              </dd>
            </dl>
          </div>

          <div className="rounded-xl border border-red-500/30 p-5">
            <h2 className="text-sm font-medium text-red-400 mb-2">
              Danger zone
            </h2>
            <p className="text-sm text-[var(--text-secondary)] mb-4">
              Delete this project and all of its deployments. This cannot be
              undone.
            </p>
            <Button
              variant="danger"
              onClick={() => setConfirmDelete(true)}
            >
              <Trash2 size={16} />
              Delete project
            </Button>
          </div>
        </div>
      )}

      <Dialog
        open={!!promoteId}
        onClose={() => setPromoteId(null)}
        title="Promote deployment"
        description="Make this successful deployment the live site. No rebuild — files already in storage."
        confirmLabel="Promote"
        loading={promote.isPending}
        onConfirm={async () => {
          if (!promoteId) return;
          try {
            await promote.mutateAsync({
              projectId: project.project_id,
              deploymentId: promoteId,
            });
            toast.success("Deployment promoted");
            setPromoteId(null);
          } catch (err) {
            toast.error(
              err instanceof ApiRequestError
                ? err.message
                : "Could not promote"
            );
          }
        }}
      />

      <Dialog
        open={confirmDelete}
        onClose={() => setConfirmDelete(false)}
        title="Delete project"
        description={`Permanently delete “${project.name}” and all deployments?`}
        confirmLabel="Delete"
        confirmVariant="danger"
        loading={deleteProject.isPending}
        onConfirm={async () => {
          try {
            await deleteProject.mutateAsync(project.project_id);
            toast.success("Project deleted");
            navigate("/dashboard", { replace: true });
          } catch (err) {
            toast.error(
              err instanceof ApiRequestError
                ? err.message
                : "Could not delete project"
            );
          }
        }}
      />
    </div>
  );
}
