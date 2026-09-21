import { cloneElement, type ReactElement, type ReactNode } from "react";
import type { FieldError, UseFormRegisterReturn } from "react-hook-form";
import { useTranslation } from "react-i18next";

interface AuthFormFieldsProps {
  email: UseFormRegisterReturn;
  password: UseFormRegisterReturn;
  confirmPassword?: UseFormRegisterReturn;
  errors: { email?: FieldError; password?: FieldError; confirmPassword?: FieldError };
}

export function AuthFormFields({ email, password, confirmPassword, errors }: AuthFormFieldsProps) {
  const { t } = useTranslation();
  return (
    <div className="space-y-4">
      <Field error={errors.email?.message} label={t("auth.fields.email")}><input autoComplete="email" {...email} id="email" inputMode="email" type="email" /></Field>
      <Field error={errors.password?.message} label={t("auth.fields.password")}><input autoComplete={confirmPassword === undefined ? "current-password" : "new-password"} {...password} id="password" type="password" /></Field>
      {confirmPassword === undefined ? null : <Field error={errors.confirmPassword?.message} label={t("auth.fields.confirmPassword")}><input autoComplete="new-password" {...confirmPassword} id="confirmPassword" type="password" /></Field>}
    </div>
  );
}

interface FieldProps {
  label: string;
  error?: string;
  children: ReactNode;
}

function Field({ label, error, children }: FieldProps) {
  const { t } = useTranslation();
  const input = children as ReactElement<{ id: string; className?: string; "aria-describedby"?: string; "aria-invalid"?: boolean }>;
  const describedBy = error === undefined ? undefined : `${input.props.id}-error`;
  return <div><label className="mb-1.5 block text-sm font-medium text-foreground" htmlFor={input.props.id}>{label}</label>{cloneElement(input, { "aria-describedby": describedBy, "aria-invalid": error !== undefined, className: "h-11 w-full rounded-xl border bg-white/70 px-3 text-foreground outline-none transition placeholder:text-muted-foreground focus:border-primary dark:bg-slate-900/70" })}{error === undefined ? null : <p className="mt-1 text-sm text-red-600 dark:text-red-400" id={describedBy} role="alert">{t(error)}</p>}</div>;
}
