import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { Link } from "react-router-dom";
import { z } from "zod";
import { toApiError } from "../../../lib/api-client";
import { AuthFormFields } from "../components/auth-form-fields";
import { useRegister } from "../hooks/use-register";
import type { RegisterInput } from "../types";

const registerSchema = z.object({ email: z.string().trim().email("auth.validation.email"), password: z.string().min(8, "auth.validation.password"), confirmPassword: z.string().min(8, "auth.validation.confirmPassword") }).refine((value) => value.password === value.confirmPassword, { path: ["confirmPassword"], message: "auth.validation.passwordMatch" });

export function RegisterPage() {
  const { t } = useTranslation();
  const registerMutation = useRegister();
  const form = useForm<RegisterInput>({ resolver: zodResolver(registerSchema), defaultValues: { email: "", password: "", confirmPassword: "" } });
  const error = registerMutation.isError ? toApiError(registerMutation.error) : undefined;
  const onSubmit = form.handleSubmit(async (input) => { await registerMutation.mutateAsync(input); form.reset(); });

  return <div><p className="text-sm font-semibold text-primary">GoPA</p><h1 className="mt-2 text-2xl font-semibold tracking-tight">{t("auth.register.title")}</h1><p className="mt-1 text-muted-foreground">{t("auth.register.description")}</p>{registerMutation.isSuccess ? <div className="mt-6 rounded-xl bg-emerald-50 p-3 text-sm text-emerald-800 dark:bg-emerald-950/40 dark:text-emerald-200" role="status">{t("auth.register.success")}</div> : <form className="mt-7 space-y-5" noValidate onSubmit={onSubmit}><AuthFormFields confirmPassword={form.register("confirmPassword")} email={form.register("email")} errors={form.formState.errors} password={form.register("password")} />{error === undefined ? null : <p className="rounded-xl bg-red-50 p-3 text-sm text-red-700 dark:bg-red-950/40 dark:text-red-300" role="alert">{t("auth.register.error")}</p>}<button className="h-11 w-full rounded-xl bg-primary px-4 text-sm font-semibold text-primary-foreground disabled:cursor-not-allowed disabled:opacity-60" disabled={registerMutation.isPending} type="submit">{registerMutation.isPending ? t("auth.register.pending") : t("auth.register.submit")}</button></form>}<p className="mt-6 text-sm text-muted-foreground">{t("auth.register.loginPrompt")} <Link className="font-semibold text-primary hover:underline" to="/login">{t("auth.register.loginAction")}</Link></p></div>;
}
