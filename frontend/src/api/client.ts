import type { ApiError } from "./types";

const API_BASE = import.meta.env.VITE_API_BASE || "";

// ─── Typed error class ───────────────────────────────────────────────────────

export class ApiRequestError extends Error {
  status: number;
  body: ApiError;

  constructor(status: number, body: ApiError) {
    super(body.error);
    this.name = "ApiRequestError";
    this.status = status;
    this.body = body;
  }
}

// ─── Shared in-flight refresh promise ────────────────────────────────────────

let refreshPromise: Promise<boolean> | null = null;

async function refreshTokens(): Promise<boolean> {
  try {
    const res = await fetch(`${API_BASE}/v1/auth/refresh`, {
      method: "POST",
      credentials: "include",
    });
    return res.ok;
  } catch {
    return false;
  }
}

/**
 * Ensure only one refresh request is in-flight at any time.
 * Concurrent 401s will all await the same promise.
 */
function ensureRefresh(): Promise<boolean> {
  if (!refreshPromise) {
    refreshPromise = refreshTokens().finally(() => {
      refreshPromise = null;
    });
  }
  return refreshPromise;
}

// ─── Core fetch wrapper ──────────────────────────────────────────────────────

interface RequestOptions {
  method: string;
  path: string;
  body?: unknown;
  /** Skip the 401 → refresh → retry cycle (used by the refresh call itself). */
  skipRetry?: boolean;
}

async function request<T>(opts: RequestOptions): Promise<T> {
  const { method, path, body, skipRetry = false } = opts;

  const url = `${API_BASE}${path}`;

  const headers: Record<string, string> = {};
  if (body !== undefined) {
    headers["Content-Type"] = "application/json";
  }

  const init: RequestInit = {
    method,
    headers,
    credentials: "include",
    body: body !== undefined ? JSON.stringify(body) : undefined,
  };

  let res = await fetch(url, init);

  // ── 401 → refresh → retry once ──
  if (res.status === 401 && !skipRetry) {
    const refreshed = await ensureRefresh();
    if (refreshed) {
      res = await fetch(url, init);
    }
  }

  // ── Parse response ──
  if (!res.ok) {
    let errorBody: ApiError;
    try {
      errorBody = (await res.json()) as ApiError;
    } catch {
      errorBody = { error: res.statusText || "Request failed" };
    }
    throw new ApiRequestError(res.status, errorBody);
  }

  // Some endpoints return 204 No Content
  if (res.status === 204) {
    return undefined as T;
  }

  return (await res.json()) as T;
}

// ─── Public helpers ──────────────────────────────────────────────────────────

export function apiGet<T>(path: string): Promise<T> {
  return request<T>({ method: "GET", path });
}

export function apiPost<T>(path: string, body?: unknown): Promise<T> {
  return request<T>({ method: "POST", path, body });
}

export function apiPut<T>(path: string, body?: unknown): Promise<T> {
  return request<T>({ method: "PUT", path, body });
}

export function apiDelete<T>(path: string): Promise<T> {
  return request<T>({ method: "DELETE", path });
}
