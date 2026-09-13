import { useTranslation } from "react-i18next";
import { cn } from "../../lib/cn";

export interface AmountDisplayProps {
  amount: string | number;
  currency?: string;
  className?: string;
  signed?: "positive" | "negative" | "neutral";
  colorize?: boolean;
}

function numberLocale(language: string | undefined) {
  if (language === "vi") return "vi-VN";
  if (language === "ja") return "ja-JP";
  return "en-US";
}

export function AmountDisplay({
  amount,
  currency = "VND",
  className,
  signed = "neutral",
  colorize = true,
}: AmountDisplayProps) {
  const { i18n } = useTranslation();
  const numeric = typeof amount === "number" ? amount : Number(amount);
  const absolute = Number.isFinite(numeric) ? Math.abs(numeric) : 0;
  const value = signed === "negative" ? -absolute : signed === "positive" ? absolute : numeric;
  
  const formatted = new Intl.NumberFormat(numberLocale(i18n.resolvedLanguage), {
    style: "currency",
    currency,
    currencyDisplay: "narrowSymbol",
    signDisplay: signed === "neutral" ? "auto" : "always",
  }).format(Number.isFinite(value) ? value : 0);

  const isPositive = signed === "positive" || (signed === "neutral" && numeric > 0);
  const isNegative = signed === "negative" || (signed === "neutral" && numeric < 0);

  const colorClass = colorize
    ? isPositive && signed !== "neutral"
      ? "text-emerald-400 dark:text-emerald-400 font-semibold"
      : isNegative && signed !== "neutral"
      ? "text-rose-400 dark:text-rose-400 font-semibold"
      : ""
    : "";

  return (
    <span className={cn("font-mono tabular-nums tracking-tight", colorClass, className)}>
      {formatted}
    </span>
  );
}
