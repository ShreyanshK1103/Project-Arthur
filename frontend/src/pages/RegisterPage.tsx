import { Link, useNavigate } from "react-router";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { toast } from "sonner";
import { AuthShell } from "@/components/ui/AuthShell";
import { Button } from "@/components/ui/Button";
import { GitHubIcon } from "@/components/ui/GitHubIcon";
import { Input } from "@/components/ui/Input";
import { Label } from "@/components/ui/Label";
import { useLogin, useRegister } from "@/api/hooks";
import { ApiRequestError } from "@/api/client";
import {
  registerSchema,
  type RegisterFormValues,
} from "@/features/auth/schemas";

export function RegisterPage() {
  const navigate = useNavigate();
  const registerMutation = useRegister();
  const login = useLogin();

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<RegisterFormValues>({
    resolver: zodResolver(registerSchema),
    defaultValues: { name: "", email: "", password: "" },
  });

  const onSubmit = handleSubmit(async (values) => {
    try {
      await registerMutation.mutateAsync(values);
      await login.mutateAsync({
        email: values.email,
        password: values.password,
      });
      toast.success("Account created");
      navigate("/dashboard", { replace: true });
    } catch (err) {
      const message =
        err instanceof ApiRequestError
          ? err.message
          : "Could not create account. Try again.";
      toast.error(message);
    }
  });

  return (
    <AuthShell
      title="Create your account"
      subtitle="Deploy static sites from GitHub in minutes."
      footer={
        <>
          Already have an account?{" "}
          <Link
            to="/login"
            className="text-[var(--accent)] hover:underline font-medium"
          >
            Sign in
          </Link>
        </>
      }
    >
      <Button
        type="button"
        variant="outline"
        className="w-full h-11"
        onClick={() => {
          window.location.href = "/v1/auth/github/login";
        }}
      >
        <GitHubIcon size={18} />
        Continue with GitHub
      </Button>
      <p className="mt-2 mb-5 text-center text-[11px] text-[var(--text-tertiary)]">
        GitHub sign-in uses your GitHub name and avatar in the dashboard.
      </p>

      <div className="relative mb-5">
        <div className="absolute inset-0 flex items-center">
          <div className="w-full border-t border-[var(--border-primary)]" />
        </div>
        <div className="relative flex justify-center text-xs">
          <span className="bg-[var(--surface-secondary)] px-2 text-[var(--text-tertiary)] rounded">
            or email
          </span>
        </div>
      </div>

      <form onSubmit={onSubmit} className="space-y-4" noValidate>
        <div>
          <Label htmlFor="name">Name</Label>
          <Input
            id="name"
            autoComplete="name"
            error={!!errors.name}
            {...register("name")}
          />
          {errors.name && (
            <p className="mt-1 text-xs text-red-400">{errors.name.message}</p>
          )}
        </div>

        <div>
          <Label htmlFor="email">Email</Label>
          <Input
            id="email"
            type="email"
            autoComplete="email"
            error={!!errors.email}
            {...register("email")}
          />
          {errors.email && (
            <p className="mt-1 text-xs text-red-400">{errors.email.message}</p>
          )}
        </div>

        <div>
          <Label htmlFor="password">Password</Label>
          <Input
            id="password"
            type="password"
            autoComplete="new-password"
            error={!!errors.password}
            {...register("password")}
          />
          {errors.password && (
            <p className="mt-1 text-xs text-red-400">
              {errors.password.message}
            </p>
          )}
          <p className="mt-1.5 text-xs text-[var(--text-tertiary)]">
            8+ chars with upper, lower, digit, and special character.
          </p>
        </div>

        <Button
          type="submit"
          className="w-full"
          loading={registerMutation.isPending || login.isPending}
        >
          Create account
        </Button>
      </form>
    </AuthShell>
  );
}
