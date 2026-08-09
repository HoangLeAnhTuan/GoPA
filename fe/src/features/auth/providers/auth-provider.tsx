import { useQuery, useQueryClient } from "@tanstack/react-query";
import type { PropsWithChildren } from "react";
import { clearAccessToken, logout, restoreSession } from "../api/auth-api";
import { authKeys } from "../api/auth-keys";
import { AuthContext } from "./auth-context";

export function AuthProvider({ children }: PropsWithChildren) {
  const queryClient = useQueryClient();
  const session = useQuery({ queryKey: authKeys.currentUser, queryFn: restoreSession, retry: false, staleTime: Infinity });
  const status = session.isPending ? "loading" : session.isSuccess ? "authenticated" : "unauthenticated";

  async function endSession() {
    try {
      await logout();
    } finally {
      clearAccessToken();
      queryClient.removeQueries();
    }
  }

  return <AuthContext.Provider value={{ user: session.data ?? null, status, logout: endSession }}>{children}</AuthContext.Provider>;
}
