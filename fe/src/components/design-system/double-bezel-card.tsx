import React from "react";
import { cn } from "../../lib/cn";

export interface DoubleBezelCardProps extends React.HTMLAttributes<HTMLDivElement> {
  outerClassName?: string;
  glowColor?: string;
  innerClassName?: string;
}

export const DoubleBezelCard = React.forwardRef<HTMLDivElement, DoubleBezelCardProps>(
  ({ children, className, outerClassName, innerClassName, glowColor, style, ...props }, ref) => {
    const dynamicStyle = glowColor
      ? {
          boxShadow: `0 0 35px -8px ${glowColor}`,
          ...style,
        }
      : style;

    return (
      <div
        ref={ref}
        style={dynamicStyle}
        className={cn(
          "group/bezel relative rounded-[1.5rem] p-1.5",
          "border border-slate-200/80 bg-slate-100/75 shadow-sm shadow-slate-200/50 backdrop-blur-xl",
          "dark:border-white/10 dark:bg-white/[0.04] dark:shadow-none",
          "transition-all duration-500 ease-[cubic-bezier(0.32,0.72,0,1)]",
          "hover:border-slate-300 dark:hover:border-white/20 hover:bg-slate-200/50 dark:hover:bg-white/[0.07]",
          outerClassName
        )}
        {...props}
      >
        <div
          className={cn(
            "relative h-full w-full rounded-[calc(1.5rem-0.375rem)] p-5 sm:p-6",
            "border border-white/80 bg-white/95 text-foreground shadow-sm shadow-slate-900/5",
            "dark:border-white/5 dark:bg-slate-950/80 dark:text-foreground dark:shadow-[inset_0_1px_1px_rgba(255,255,255,0.12)]",
            "transition-colors duration-300",
            innerClassName,
            className
          )}
        >
          {children}
        </div>
      </div>
    );
  }
);

DoubleBezelCard.displayName = "DoubleBezelCard";
