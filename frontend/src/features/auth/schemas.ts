import { z } from "zod";

/** Mirrors backend/internal/validators ValidateEmail + ValidatePassword. */
export const loginSchema = z.object({
  email: z
    .string()
    .trim()
    .min(1, "Email is required")
    .email("Invalid email address"),
  password: z.string().min(1, "Password is required"),
});

export const registerSchema = z.object({
  name: z.string().trim().min(1, "Name is required"),
  email: z
    .string()
    .trim()
    .min(1, "Email is required")
    .email("Invalid email address"),
  password: z
    .string()
    .min(8, "Password must be at least 8 characters long")
    .refine((p) => /[A-Z]/.test(p), "Password must contain an uppercase letter")
    .refine((p) => /[a-z]/.test(p), "Password must contain a lowercase letter")
    .refine((p) => /\d/.test(p), "Password must contain a digit")
    .refine(
      (p) => /[^A-Za-z0-9]/.test(p),
      "Password must contain a special character"
    ),
});

export type LoginFormValues = z.infer<typeof loginSchema>;
export type RegisterFormValues = z.infer<typeof registerSchema>;
