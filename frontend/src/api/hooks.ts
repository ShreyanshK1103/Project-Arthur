import {
  useQuery,
  useMutation,
  useQueryClient,
  keepPreviousData,
} from "@tanstack/react-query";
import { useEffect, useRef, useState } from "react";
import { apiGet, apiPost, apiPut, apiDelete, ApiRequestError } from "./client";
import type {
  User,
  LoginRequest,
  RegisterRequest,
  Project,
  CreateProjectRequest,
  Deployment,
  DeploymentLogsResponse,
  EnvVar,
  SetEnvVarRequest,
} from "./types";

// ─── Query keys ──────────────────────────────────────────────────────────────

export const queryKeys = {
  me: ["auth", "me"] as const,
  projects: ["projects"] as const,
  project: (id: string) => ["projects", id] as const,
  deployments: (projectId: string) =>
    ["projects", projectId, "deployments"] as const,
  deployment: (id: string) => ["deployments", id] as const,
  env: (projectId: string) => ["projects", projectId, "env"] as const,
} as const;

// ─── Auth ────────────────────────────────────────────────────────────────────

/**
 * Fetches the current user via `GET /v1/auth/me`.
 * Returns `null` when the user is not authenticated (401).
 */
export function useMe() {
  return useQuery<User | null>({
    queryKey: queryKeys.me,
    queryFn: async () => {
      try {
        return await apiGet<User>("/v1/auth/me");
      } catch (err) {
        if (err instanceof ApiRequestError && err.status === 401) {
          return null;
        }
        // API down / network error: treat as logged out so auth pages still render.
        if (!(err instanceof ApiRequestError)) {
          return null;
        }
        throw err;
      }
    },
    staleTime: 5 * 60 * 1000,
    retry: false,
  });
}

export function useLogin() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (body: LoginRequest) =>
      apiPost<{ message: string }>("/v1/auth/login", body),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: queryKeys.me });
    },
  });
}

export function useRegister() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (body: RegisterRequest) =>
      apiPost<User>("/v1/auth/register", body),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: queryKeys.me });
    },
  });
}

export function useLogout() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => apiPost("/v1/auth/logout"),
    onSuccess: () => {
      queryClient.setQueryData(queryKeys.me, null);
      queryClient.invalidateQueries();
    },
  });
}

// ─── Projects ────────────────────────────────────────────────────────────────

export function useProjects() {
  return useQuery({
    queryKey: queryKeys.projects,
    queryFn: () => apiGet<Project[]>("/v1/projects"),
  });
}

export function useProject(id: string | undefined) {
  return useQuery({
    queryKey: queryKeys.project(id ?? ""),
    queryFn: () => apiGet<Project>(`/v1/projects/${id}`),
    enabled: !!id,
  });
}

export function useCreateProject() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (body: CreateProjectRequest) =>
      apiPost<Project>("/v1/projects", body),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.projects });
    },
  });
}

export function useDeleteProject() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) =>
      apiDelete<{ message: string }>(`/v1/projects/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.projects });
    },
  });
}

// ─── Deployments ─────────────────────────────────────────────────────────────

export function useProjectDeployments(projectId: string | undefined) {
  return useQuery({
    queryKey: queryKeys.deployments(projectId ?? ""),
    queryFn: () =>
      apiGet<Deployment[]>(`/v1/projects/${projectId}/deployments`),
    enabled: !!projectId,
    refetchInterval: (query) => {
      const data = query.state.data;
      if (!data) return false;
      const active = data.some(
        (d) => d.status === "queued" || d.status === "building"
      );
      return active ? 2000 : false;
    },
  });
}

export function useDeployment(id: string | undefined) {
  return useQuery({
    queryKey: queryKeys.deployment(id ?? ""),
    queryFn: () => apiGet<Deployment>(`/v1/deployments/${id}`),
    enabled: !!id,
    refetchInterval: (query) => {
      const status = query.state.data?.status;
      if (status === "queued" || status === "building") return 2000;
      return false;
    },
    placeholderData: keepPreviousData,
  });
}

export function useRedeploy() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (projectId: string) =>
      apiPost<Deployment>(`/v1/projects/${projectId}/redeploy`),
    onSuccess: (_data, projectId) => {
      queryClient.invalidateQueries({
        queryKey: queryKeys.deployments(projectId),
      });
      queryClient.invalidateQueries({ queryKey: queryKeys.project(projectId) });
    },
  });
}

export function usePromote() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      projectId,
      deploymentId,
    }: {
      projectId: string;
      deploymentId: string;
    }) =>
      apiPost<Project>(
        `/v1/projects/${projectId}/deployments/${deploymentId}/promote`
      ),
    onSuccess: (_data, { projectId }) => {
      queryClient.invalidateQueries({
        queryKey: queryKeys.deployments(projectId),
      });
      queryClient.invalidateQueries({ queryKey: queryKeys.project(projectId) });
      queryClient.invalidateQueries({ queryKey: queryKeys.projects });
    },
  });
}

/** Polls deployment logs until `done` is true. */
export function useDeploymentLogs(deploymentId: string | undefined) {
  const [snapshot, setSnapshot] = useState<{
    id: string | undefined;
    logs: DeploymentLogsResponse["logs"];
    status: string | null;
    done: boolean;
    error: string | null;
  }>({ id: undefined, logs: [], status: null, done: false, error: null });

  const afterRef = useRef(0);
  const backoffRef = useRef(1500);

  useEffect(() => {
    if (!deploymentId) return;

    afterRef.current = 0;
    backoffRef.current = 1500;

    let cancelled = false;
    let timer: ReturnType<typeof setTimeout>;

    const poll = async () => {
      try {
        const res = await apiGet<DeploymentLogsResponse>(
          `/v1/deployments/${deploymentId}/logs?after=${afterRef.current}`
        );
        if (cancelled) return;

        if (res.logs.length > 0) {
          afterRef.current = res.next_after;
        }
        backoffRef.current = 1500;

        setSnapshot((prev) => ({
          id: deploymentId,
          logs:
            prev.id === deploymentId && res.logs.length > 0
              ? [...prev.logs, ...res.logs]
              : prev.id === deploymentId
                ? prev.logs
                : res.logs,
          status: res.status,
          done: res.done,
          error: null,
        }));

        if (!res.done) {
          timer = setTimeout(poll, backoffRef.current);
        }
      } catch (err) {
        if (cancelled) return;
        setSnapshot((prev) => ({
          id: deploymentId,
          logs: prev.id === deploymentId ? prev.logs : [],
          status: prev.id === deploymentId ? prev.status : null,
          done: false,
          error:
            err instanceof ApiRequestError
              ? err.message
              : "Failed to load logs",
        }));
        backoffRef.current = Math.min(backoffRef.current * 2, 10000);
        timer = setTimeout(poll, backoffRef.current);
      }
    };

    poll();

    return () => {
      cancelled = true;
      clearTimeout(timer);
    };
  }, [deploymentId]);

  const active = snapshot.id === deploymentId;
  return {
    logs: active ? snapshot.logs : [],
    status: active ? snapshot.status : null,
    done: active ? snapshot.done : false,
    error: active ? snapshot.error : null,
  };
}

// ─── Env vars ────────────────────────────────────────────────────────────────

export function useEnvVars(projectId: string | undefined) {
  return useQuery({
    queryKey: queryKeys.env(projectId ?? ""),
    queryFn: () => apiGet<EnvVar[]>(`/v1/projects/${projectId}/env`),
    enabled: !!projectId,
  });
}

export function useSetEnvVar(projectId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (body: SetEnvVarRequest) =>
      apiPut<{ key: string }>(`/v1/projects/${projectId}/env`, body),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.env(projectId) });
    },
  });
}

export function useDeleteEnvVar(projectId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (key: string) =>
      apiDelete<{ message: string }>(
        `/v1/projects/${projectId}/env/${encodeURIComponent(key)}`
      ),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.env(projectId) });
    },
  });
}
