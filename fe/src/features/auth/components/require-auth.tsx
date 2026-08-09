import type { ReactNode } from "react";
import { useTranslation } from "react-i18next";
import { Navigate, useLocation } from "react-router-dom";
import { useAuth } from "../hooks/use-auth";

export function RequireAuth({ children }: { children: ReactNode }) {
	const { t } = useTranslation();
	const { status } = useAuth();
	const location = useLocation();
	if (status === "loading") {
		return <main className="grid min-h-screen place-items-center bg-background text-muted-foreground">{t("auth.login.pending")}</main>;
	}
	if (status === "unauthenticated") {
		return <Navigate replace state={{ from: location.pathname }} to="/login" />;
	}
	return children;
}
