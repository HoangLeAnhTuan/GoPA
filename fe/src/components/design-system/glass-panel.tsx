import type { HTMLAttributes } from "react";
import { cn } from "../../lib/cn";

export function GlassPanel({ className, ...props }: HTMLAttributes<HTMLDivElement>) {
  return <div className={cn("rounded-2xl border border-white/50 bg-white/75 shadow-lg shadow-slate-950/5 backdrop-blur-xl dark:border-white/10 dark:bg-slate-950/70", className)} {...props} />;
}
