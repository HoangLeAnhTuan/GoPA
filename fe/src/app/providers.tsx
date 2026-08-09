import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import type { PropsWithChildren } from "react";
import { ThemeController } from "../components/design-system/theme-controller";
import { AuthProvider } from "../features/auth/providers/auth-provider";
import "../i18n/config";

const queryClient = new QueryClient({
  defaultOptions: {
    queries: { retry: 1, refetchOnWindowFocus: false },
  },
});

export function AppProviders({ children }: PropsWithChildren) {
  return (
    <QueryClientProvider client={queryClient}>
      <ThemeController><AuthProvider>{children}</AuthProvider></ThemeController>
    </QueryClientProvider>
  );
}
