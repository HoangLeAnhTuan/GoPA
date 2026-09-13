import React from "react";
import { motion, type HTMLMotionProps } from "framer-motion";
import type { LucideIcon } from "lucide-react";
import { cn } from "../../lib/cn";

export type ButtonVariant = "primary" | "secondary" | "emerald" | "violet" | "amber" | "rose" | "ghost";

export interface ButtonInButtonProps extends Omit<HTMLMotionProps<"button">, "children"> {
  icon?: LucideIcon;
  variant?: ButtonVariant;
  children: React.ReactNode;
  iconPosition?: "right" | "left";
  iconClassName?: string;
  size?: "sm" | "md" | "lg";
}

const variantStyles: Record<ButtonVariant, string> = {
  primary: "bg-blue-600 hover:bg-blue-500 text-white shadow-md shadow-blue-600/25 border border-blue-400/30",
  secondary: "bg-slate-100 hover:bg-slate-200 text-foreground border border-slate-200/80 shadow-sm dark:bg-white/10 dark:hover:bg-white/15 dark:text-foreground dark:border-white/10",
  emerald: "bg-emerald-600 hover:bg-emerald-500 text-white shadow-md shadow-emerald-600/25 border border-emerald-400/30",
  violet: "bg-violet-600 hover:bg-violet-500 text-white shadow-md shadow-violet-600/25 border border-violet-400/30",
  amber: "bg-amber-600 hover:bg-amber-500 text-white shadow-md shadow-amber-600/25 border border-amber-400/30",
  rose: "bg-rose-600 hover:bg-rose-500 text-white shadow-md shadow-rose-600/25 border border-rose-400/30",
  ghost: "bg-transparent hover:bg-slate-100 text-muted-foreground hover:text-foreground border border-transparent dark:hover:bg-white/10",
};

// Padding when an icon is on the right
const sizeStylesRightIcon = {
  sm: "pl-4 pr-1.5 py-1.5 text-xs min-h-[36px]",
  md: "pl-5 pr-2 py-2 text-sm min-h-[42px]",
  lg: "pl-6 pr-2.5 py-2.5 text-base min-h-[48px]",
};

// Padding when an icon is on the left
const sizeStylesLeftIcon = {
  sm: "pl-1.5 pr-4 py-1.5 text-xs min-h-[36px] flex-row-reverse",
  md: "pl-2 pr-5 py-2 text-sm min-h-[42px] flex-row-reverse",
  lg: "pl-2.5 pr-6 py-2.5 text-base min-h-[48px] flex-row-reverse",
};

// Padding when there is no icon
const sizeStylesNoIcon = {
  sm: "px-4 py-1.5 text-xs min-h-[36px]",
  md: "px-5 py-2 text-sm min-h-[42px]",
  lg: "px-6 py-2.5 text-base min-h-[48px]",
};

const iconWrapperSizes = {
  sm: "size-6",
  md: "size-7",
  lg: "size-8",
};

const iconSizes = {
  sm: 12,
  md: 14,
  lg: 16,
};

export const ButtonInButton = React.forwardRef<HTMLButtonElement, ButtonInButtonProps>(
  (
    {
      children,
      icon: Icon,
      variant = "primary",
      size = "md",
      iconPosition = "right",
      iconClassName,
      className,
      disabled,
      type = "button",
      ...props
    },
    ref
  ) => {
    const hasIcon = Boolean(Icon);
    const sizeCls = hasIcon
      ? iconPosition === "left"
        ? sizeStylesLeftIcon[size]
        : sizeStylesRightIcon[size]
      : sizeStylesNoIcon[size];

    return (
      <motion.button
        ref={ref}
        type={type}
        disabled={disabled}
        whileHover={disabled ? undefined : { scale: 1.015 }}
        whileTap={disabled ? undefined : { scale: 0.98 }}
        transition={{ type: "spring", stiffness: 450, damping: 28 }}
        className={cn(
          "group relative inline-flex items-center justify-center gap-3 rounded-full font-semibold",
          "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2 focus-visible:ring-offset-background",
          "transition-colors duration-200 disabled:opacity-50 disabled:pointer-events-none cursor-pointer",
          variantStyles[variant],
          sizeCls,
          className
        )}
        {...props}
      >
        <span className="tracking-tight select-none leading-none inline-flex items-center text-center">
          {children}
        </span>

        {Icon && (
          <div
            className={cn(
              "flex shrink-0 items-center justify-center rounded-full",
              variant === "secondary" || variant === "ghost"
                ? "bg-slate-900/10 text-foreground dark:bg-white/15 dark:text-white"
                : "bg-white/25 text-white dark:bg-white/20",
              "transition-transform duration-300 ease-[cubic-bezier(0.32,0.72,0,1)]",
              "group-hover:translate-x-0.5 group-hover:-translate-y-0.5 group-hover:scale-105",
              iconWrapperSizes[size],
              iconClassName
            )}
          >
            <Icon size={iconSizes[size]} aria-hidden="true" />
          </div>
        )}
      </motion.button>
    );
  }
);

ButtonInButton.displayName = "ButtonInButton";
