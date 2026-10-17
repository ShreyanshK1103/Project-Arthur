import { z } from "zod";

/** Mirrors backend validRepoURL: https only, host required, no userinfo. */
function isValidRepoURL(raw: string): boolean {
  try {
    const u = new URL(raw);
    return u.protocol === "https:" && !!u.host && !u.username && !u.password;
  } catch {
    return false;
  }
}

export const createProjectSchema = z.object({
  name: z.string().trim().min(1, "Name is required"),
  repo_url: z
    .string()
    .trim()
    .min(1, "Repository URL is required")
    .refine(isValidRepoURL, "Must be a public https URL (no credentials)"),
  branch: z.string().trim().optional(),
  install_command: z.string().trim().optional(),
  build_command: z.string().trim().optional(),
  output_dir: z.string().trim().optional(),
});

export type CreateProjectFormValues = z.infer<typeof createProjectSchema>;

export const envVarSchema = z.object({
  key: z
    .string()
    .trim()
    .min(1, "Key is required")
    .max(128, "Key must be at most 128 characters")
    .regex(
      /^[A-Za-z_][A-Za-z0-9_]*$/,
      "Key must start with a letter or underscore"
    ),
  value: z
    .string()
    .min(1, "Value is required")
    .max(4096, "Value must be at most 4096 characters"),
});

export type EnvVarFormValues = z.infer<typeof envVarSchema>;
