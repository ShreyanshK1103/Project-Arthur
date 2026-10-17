import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Eye, EyeOff, Plus, Trash2 } from "lucide-react";
import { toast } from "sonner";
import {
  useDeleteEnvVar,
  useEnvVars,
  useSetEnvVar,
} from "@/api/hooks";
import { ApiRequestError } from "@/api/client";
import { Button } from "@/components/ui/Button";
import { Dialog } from "@/components/ui/Dialog";
import { Input } from "@/components/ui/Input";
import { Label } from "@/components/ui/Label";
import { Skeleton } from "@/components/ui/Skeleton";
import {
  envVarSchema,
  type EnvVarFormValues,
} from "@/features/projects/schemas";
import { formatExactTime, formatRelativeTime } from "@/lib/utils";

interface EnvVarsPanelProps {
  projectId: string;
}

export function EnvVarsPanel({ projectId }: EnvVarsPanelProps) {
  const { data, isLoading, isError, error } = useEnvVars(projectId);
  const setEnv = useSetEnvVar(projectId);
  const deleteEnv = useDeleteEnvVar(projectId);
  const [showForm, setShowForm] = useState(false);
  const [deleteKey, setDeleteKey] = useState<string | null>(null);
  const [revealed, setRevealed] = useState(false);

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<EnvVarFormValues>({
    resolver: zodResolver(envVarSchema),
    defaultValues: { key: "", value: "" },
  });

  const onSubmit = handleSubmit(async (values) => {
    try {
      await setEnv.mutateAsync(values);
      toast.success(`Saved ${values.key}`);
      reset();
      setShowForm(false);
      setRevealed(false);
    } catch (err) {
      toast.error(
        err instanceof ApiRequestError ? err.message : "Could not save variable"
      );
    }
  });

  return (
    <div>
      <div className="rounded-lg border border-amber-500/30 bg-amber-500/5 px-4 py-3 mb-6 text-sm text-amber-200/90">
        Static-site build variables end up in the public bundle. Don&apos;t put
        secrets here.
      </div>

      <div className="flex items-center justify-between mb-4">
        <h2 className="text-sm font-medium text-[var(--text-primary)]">
          Variables
        </h2>
        <Button
          type="button"
          size="sm"
          variant="secondary"
          onClick={() => setShowForm((v) => !v)}
        >
          <Plus size={14} />
          Add variable
        </Button>
      </div>

      {showForm && (
        <form
          onSubmit={onSubmit}
          className="mb-6 rounded-xl border border-[var(--border-primary)] p-4 space-y-3"
          noValidate
        >
          <div className="grid sm:grid-cols-2 gap-3">
            <div>
              <Label htmlFor="env-key">Key</Label>
              <Input
                id="env-key"
                className="font-mono"
                error={!!errors.key}
                {...register("key")}
              />
              {errors.key && (
                <p className="mt-1 text-xs text-red-400">{errors.key.message}</p>
              )}
            </div>
            <div>
              <Label htmlFor="env-value">Value</Label>
              <div className="relative">
                <Input
                  id="env-value"
                  type={revealed ? "text" : "password"}
                  className="font-mono pr-10"
                  error={!!errors.value}
                  {...register("value")}
                />
                <button
                  type="button"
                  onClick={() => setRevealed((v) => !v)}
                  className="absolute right-2 top-1/2 -translate-y-1/2 text-[var(--text-tertiary)] hover:text-[var(--text-primary)]"
                  aria-label={revealed ? "Hide value" : "Reveal value"}
                >
                  {revealed ? <EyeOff size={16} /> : <Eye size={16} />}
                </button>
              </div>
              {errors.value && (
                <p className="mt-1 text-xs text-red-400">{errors.value.message}</p>
              )}
            </div>
          </div>
          <div className="flex gap-2">
            <Button type="submit" size="sm" loading={setEnv.isPending}>
              Save
            </Button>
            <Button
              type="button"
              size="sm"
              variant="ghost"
              onClick={() => {
                setShowForm(false);
                reset();
              }}
            >
              Cancel
            </Button>
          </div>
        </form>
      )}

      {isLoading && (
        <div className="space-y-2">
          <Skeleton className="h-12" />
          <Skeleton className="h-12" />
        </div>
      )}

      {isError && (
        <p className="text-sm text-red-400">
          {error instanceof Error ? error.message : "Failed to load variables"}
        </p>
      )}

      {!isLoading && data && data.length === 0 && (
        <p className="text-sm text-[var(--text-tertiary)]">
          No environment variables yet.
        </p>
      )}

      {!isLoading && data && data.length > 0 && (
        <ul className="divide-y divide-[var(--border-primary)] rounded-xl border border-[var(--border-primary)] overflow-hidden">
          {data.map((row) => (
            <li
              key={row.key}
              className="flex items-center justify-between gap-3 px-4 py-3 bg-[var(--surface-primary)]"
            >
              <div className="min-w-0">
                <p className="font-mono text-sm text-[var(--text-primary)] truncate">
                  {row.key}
                </p>
                <p
                  className="text-xs text-[var(--text-tertiary)]"
                  title={formatExactTime(row.updated_at)}
                >
                  Updated {formatRelativeTime(row.updated_at)} · value masked
                  (API never returns values)
                </p>
              </div>
              <Button
                type="button"
                variant="ghost"
                size="sm"
                onClick={() => setDeleteKey(row.key)}
                aria-label={`Delete ${row.key}`}
              >
                <Trash2 size={14} className="text-red-400" />
              </Button>
            </li>
          ))}
        </ul>
      )}

      <Dialog
        open={!!deleteKey}
        onClose={() => setDeleteKey(null)}
        title="Delete variable"
        description={
          deleteKey
            ? `Remove ${deleteKey}? Redeploy to apply the change.`
            : undefined
        }
        confirmLabel="Delete"
        confirmVariant="danger"
        loading={deleteEnv.isPending}
        onConfirm={async () => {
          if (!deleteKey) return;
          try {
            await deleteEnv.mutateAsync(deleteKey);
            toast.success("Variable deleted");
            setDeleteKey(null);
          } catch (err) {
            toast.error(
              err instanceof ApiRequestError
                ? err.message
                : "Could not delete variable"
            );
          }
        }}
      />
    </div>
  );
}
