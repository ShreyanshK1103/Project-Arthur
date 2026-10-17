/**
 * TypeScript types derived 1:1 from the Go backend models.
 * Source: backend/internal/models/models.go, backend/internal/models/auth.go
 *
 * NEVER invent fields — every type here maps to a real JSON tag in the Go structs.
 */

// ─── Auth ────────────────────────────────────────────────────────────────────

export interface RegisterRequest {
  name: string;
  email: string;
  password: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

// ─── User ────────────────────────────────────────────────────────────────────

export interface User {
  id: string;
  created_at: string;
  updated_at: string;
  name: string;
  email: string;
  avatar_url?: string;
  /** Present on the User model; `/auth/me` may omit it. */
  provider?: string;
  email_verified: boolean;
}

// ─── Project ─────────────────────────────────────────────────────────────────

export interface Project {
  project_id: string;
  name: string;
  user_id: string;
  repo_url: string;
  install_command: string;
  build_command: string;
  output_dir: string;
  created_at: string;
  updated_at: string;
  github_repo_id?: number;
  branch: string;
  auto_deploy: boolean;
  slug: string;
  url: string;
  active_deployment_id?: string;
}

export interface CreateProjectRequest {
  name: string;
  repo_url: string;
  install_command?: string;
  build_command?: string;
  output_dir?: string;
  branch?: string;
  github_repo_id?: number;
  auto_deploy?: boolean;
}

// ─── Deployment ──────────────────────────────────────────────────────────────

/** Status values the backend actually writes. */
export type DeploymentStatus = "queued" | "building" | "success" | "failed";

export interface Deployment {
  id: string;
  project_id: string;
  status: DeploymentStatus;
  url?: string;
  created_at: string;
  updated_at: string;
  /** Present only in the project-deployments list endpoint. */
  project_url?: string;
  /** Present only in the project-deployments list endpoint. */
  is_active?: boolean;
}

// ─── Deployment Logs ─────────────────────────────────────────────────────────

export interface LogLine {
  seq: number;
  message: string;
  created_at: string;
}

export interface DeploymentLogsResponse {
  logs: LogLine[];
  status: DeploymentStatus;
  next_after: number;
  done: boolean;
}

// ─── Environment Variables ───────────────────────────────────────────────────

export interface EnvVar {
  key: string;
  updated_at: string;
}

export interface SetEnvVarRequest {
  key: string;
  value: string;
}

// ─── Error ───────────────────────────────────────────────────────────────────

/** Every backend error uses this shape: { "error": "message" } */
export interface ApiError {
  error: string;
}
