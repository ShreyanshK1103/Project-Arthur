import { useState } from "react";
import { useNavigate } from "react-router";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Plus, Trash2 } from "lucide-react";
import { toast } from "sonner";
import { useCreateProject } from "@/api/hooks";
import { apiPost, apiPut, ApiRequestError } from "@/api/client";
import type { Deployment } from "@/api/types";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { Label } from "@/components/ui/Label";
import {
  createProjectSchema,
  type CreateProjectFormValues,
} from "@/features/projects/schemas";

interface EnvRow {
  key: string;
  value: string;
}

export function NewProjectPage() {
  const navigate = useNavigate();
  const createProject = useCreateProject();
  const [envRows, setEnvRows] = useState<EnvRow[]>([]);
  const [submitting, setSubmitting] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<CreateProjectFormValues>({
    resolver: zodResolver(createProjectSchema),
    defaultValues: {
      name: "",
      repo_url: "",
      branch: "main",
      install_command: "npm install",
      build_command: "npm run build",
      output_dir: "dist",
    },
  });

  const onSubmit = handleSubmit(async (values) => {
    setSubmitting(true);
    try {
      const project = await createProject.mutateAsync({
        name: values.name,
        repo_url: values.repo_url,
        branch: values.branch || undefined,
        install_command: values.install_command || undefined,
        build_command: values.build_command || undefined,
        output_dir: values.output_dir || undefined,
      });

      const validEnv = envRows.filter((r) => r.key.trim() && r.value);
      for (const row of validEnv) {
        await apiPut(`/v1/projects/${project.project_id}/env`, {
          key: row.key.trim(),
          value: row.value,
        });
      }

      const deployment = await apiPost<Deployment>("/v1/deployments", {
        project_id: project.project_id,
      });

      toast.success("Project created — build queued");
      navigate(
        `/projects/${project.project_id}/deployments/${deployment.id}`,
        { replace: true }
      );
    } catch (err) {
      toast.error(
        err instanceof ApiRequestError
          ? err.message
          : "Could not create project"
      );
    } finally {
      setSubmitting(false);
    }
  });

  return (
    <div className="max-w-2xl">
      <p className="text-xs uppercase tracking-[0.18em] text-[var(--accent)] mb-2">
        Deploy
      </p>
      <h1 className="text-2xl sm:text-3xl font-semibold font-display tracking-tight mb-2">
        New Project
      </h1>
      <p className="text-sm text-[var(--text-secondary)] mb-8">
        Point Arthur at a public GitHub repo. Static sites only — no long-running
        servers.
      </p>

      <form onSubmit={onSubmit} className="panel p-5 sm:p-6 space-y-5" noValidate>
        <div>
          <Label htmlFor="name">Project name</Label>
          <Input id="name" error={!!errors.name} {...register("name")} />
          {errors.name && (
            <p className="mt-1 text-xs text-red-400">{errors.name.message}</p>
          )}
        </div>

        <div>
          <Label htmlFor="repo_url">Repository URL</Label>
          <Input
            id="repo_url"
            placeholder="https://github.com/org/repo"
            error={!!errors.repo_url}
            {...register("repo_url")}
          />
          {errors.repo_url && (
            <p className="mt-1 text-xs text-red-400">{errors.repo_url.message}</p>
          )}
        </div>

        <div className="grid sm:grid-cols-2 gap-4">
          <div>
            <Label htmlFor="branch">Branch</Label>
            <Input id="branch" {...register("branch")} />
          </div>
          <div>
            <Label htmlFor="output_dir">Output directory</Label>
            <Input id="output_dir" {...register("output_dir")} />
          </div>
        </div>

        <div className="grid sm:grid-cols-2 gap-4">
          <div>
            <Label htmlFor="install_command">Install command</Label>
            <Input id="install_command" {...register("install_command")} />
          </div>
          <div>
            <Label htmlFor="build_command">Build command</Label>
            <Input id="build_command" {...register("build_command")} />
          </div>
        </div>

        <div className="rounded-xl border border-[var(--border-primary)] bg-[var(--bg-primary)]/40 p-4">
          <div className="flex items-center justify-between mb-3">
            <div>
              <h2 className="text-sm font-medium text-[var(--text-primary)]">
                Environment variables
              </h2>
              <p className="text-xs text-[var(--text-tertiary)] mt-0.5">
                Optional. Static-site build variables end up in the public
                bundle. Don&apos;t put secrets here.
              </p>
            </div>
            <Button
              type="button"
              variant="ghost"
              size="sm"
              onClick={() => setEnvRows((r) => [...r, { key: "", value: "" }])}
            >
              <Plus size={14} />
              Add
            </Button>
          </div>

          {envRows.length === 0 ? (
            <p className="text-xs text-[var(--text-tertiary)]">None added.</p>
          ) : (
            <div className="space-y-2">
              {envRows.map((row, i) => (
                <div key={i} className="flex gap-2">
                  <Input
                    placeholder="KEY"
                    value={row.key}
                    onChange={(e) => {
                      const next = [...envRows];
                      next[i] = { ...next[i], key: e.target.value };
                      setEnvRows(next);
                    }}
                    className="font-mono"
                  />
                  <Input
                    placeholder="value"
                    value={row.value}
                    onChange={(e) => {
                      const next = [...envRows];
                      next[i] = { ...next[i], value: e.target.value };
                      setEnvRows(next);
                    }}
                    className="font-mono"
                  />
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    onClick={() =>
                      setEnvRows((rows) => rows.filter((_, j) => j !== i))
                    }
                    aria-label="Remove variable"
                  >
                    <Trash2 size={14} />
                  </Button>
                </div>
              ))}
            </div>
          )}
        </div>

        <div className="flex gap-3 pt-2">
          <Button type="submit" loading={submitting}>
            Deploy
          </Button>
          <Button
            type="button"
            variant="ghost"
            onClick={() => navigate("/dashboard")}
          >
            Cancel
          </Button>
        </div>
      </form>
    </div>
  );
}
