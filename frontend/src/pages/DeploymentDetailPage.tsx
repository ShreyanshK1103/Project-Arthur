import { Link, useParams } from "react-router";
import { ArrowLeft } from "lucide-react";
import { useDeployment, useDeploymentLogs } from "@/api/hooks";
import { Button } from "@/components/ui/Button";
import { SiteLink } from "@/components/ui/SiteLink";
import { Skeleton } from "@/components/ui/Skeleton";
import { StatusBadge } from "@/components/ui/StatusBadge";
import { TimeMeta } from "@/components/ui/TimeMeta";
import { LogViewer } from "@/features/deployments/LogViewer";
import { deploymentSlug } from "@/lib/siteUrl";

export function DeploymentDetailPage() {
  const { projectId, deploymentId } = useParams<{
    projectId: string;
    deploymentId: string;
  }>();

  const deploymentQuery = useDeployment(deploymentId);
  const { logs, status, done, error } = useDeploymentLogs(deploymentId);

  if (deploymentQuery.isLoading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-8 w-48" />
        <Skeleton className="h-96" />
      </div>
    );
  }

  if (deploymentQuery.isError || !deploymentQuery.data) {
    return (
      <div className="text-center py-16">
        <p className="text-red-400 mb-4">
          {deploymentQuery.error instanceof Error
            ? deploymentQuery.error.message
            : "Deployment not found"}
        </p>
        <Link to={`/projects/${projectId}`}>
          <Button variant="secondary">Back to project</Button>
        </Link>
      </div>
    );
  }

  const deployment = deploymentQuery.data;
  const displayStatus = status ?? deployment.status;

  return (
    <div>
      <Link
        to={`/projects/${projectId}`}
        className="inline-flex items-center gap-1.5 text-sm text-[var(--text-secondary)] hover:text-[var(--text-primary)] mb-4"
      >
        <ArrowLeft size={16} />
        Project
      </Link>

      <div className="panel p-5 sm:p-6 mb-6">
        <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4">
          <div>
            <h1 className="text-2xl font-semibold font-display tracking-tight">
              Deployment
            </h1>
            <p className="text-sm text-[var(--text-tertiary)] font-mono mt-1 break-all">
              {deployment.id}
            </p>
            <TimeMeta
              createdAt={deployment.created_at}
              updatedAt={deployment.updated_at}
              className="mt-3"
            />
          </div>
          <div className="flex flex-wrap items-center gap-3">
            <StatusBadge status={displayStatus} />
            {(displayStatus === "success" || deployment.url) && (
              <SiteLink slug={deploymentSlug(deployment.id)} />
            )}
          </div>
        </div>
      </div>

      <LogViewer logs={logs} done={done} error={error} />
    </div>
  );
}
