import { zodResolver } from "@hookform/resolvers/zod";
import { useQueryClient } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { Link, useNavigate } from "react-router-dom";
import { z } from "zod";
import { toApiError } from "../../../lib/api-client";
import { authKeys } from "../api/auth-keys";
import { AuthFormFields } from "../components/auth-form-fields";
import { useLogin } from "../hooks/use-login";
import type { LoginInput } from "../types";

const loginSchema = z.object({ email: z.string().trim().email("auth.validation.email"), password: z.string().min(8, "auth.validation.password") });

export function LoginPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const loginMutation = useLogin();
  const form = useForm<LoginInput>({ resolver: zodResolver(loginSchema), defaultValues: { email: "", password: "" } });
  const error = loginMutation.isError ? toApiError(loginMutation.error) : undefined;

  const onSubmit = form.handleSubmit(async (input) => {
    const result = await loginMutation.mutateAsync(input);
    queryClient.setQueryData(authKeys.currentUser, result.user);
    await navigate("/app/today", { replace: true });
  });

  return <div><p className="text-sm font-semibold text-primary">GoPA</p><h1 className="mt-2 text-2xl font-semibold tracking-tight">{t("auth.login.title")}</h1><p className="mt-1 text-muted-foreground">{t("auth.login.description")}</p><form className="mt-7 space-y-5" noValidate onSubmit={onSubmit}><AuthFormFields email={form.register("email")} errors={form.formState.errors} password={form.register("password")} />{error === undefined ? null : <p className="rounded-xl bg-red-50 p-3 text-sm text-red-700 dark:bg-red-950/40 dark:text-red-300" role="alert">{t("auth.login.error")}</p>}<button className="h-11 w-full rounded-xl bg-primary px-4 text-sm font-semibold text-primary-foreground disabled:cursor-not-allowed disabled:opacity-60" disabled={loginMutation.isPending} type="submit">{loginMutation.isPending ? t("auth.login.pending") : t("auth.login.submit")}</button></form><p className="mt-6 text-sm text-muted-foreground">{t("auth.login.registerPrompt")} <Link className="font-semibold text-primary hover:underline" to="/register">{t("auth.login.registerAction")}</Link></p></div>;
}
