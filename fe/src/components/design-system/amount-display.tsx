import { useTranslation } from "react-i18next";
import { cn } from "../../lib/cn";

interface AmountDisplayProps {
  amount: string;
  currency: string;
  className?: string;
  signed?: "positive" | "negative" | "neutral";
}

function numberLocale(language: string | undefined) {
  if (language === "vi") return "vi-VN";
  if (language === "ja") return "ja-JP";
  return "en-US";
}

export function AmountDisplay({ amount, currency, className, signed = "neutral" }: AmountDisplayProps) {
  const { i18n } = useTranslation();
  const numeric = Number(amount);
  const absolute = Number.isFinite(numeric) ? Math.abs(numeric) : 0;
  const value = signed === "negative" ? -absolute : signed === "positive" ? absolute : numeric;
  const formatted = new Intl.NumberFormat(numberLocale(i18n.resolvedLanguage), {
    style: "currency",
    currency,
    currencyDisplay: "narrowSymbol",
    signDisplay: signed === "neutral" ? "auto" : "always",
  }).format(Number.isFinite(value) ? value : 0);

  return <span className={cn("tabular-nums", signed === "positive" && "text-emerald-500", signed === "negative" && "text-rose-500", className)}>{formatted}</span>;
}
