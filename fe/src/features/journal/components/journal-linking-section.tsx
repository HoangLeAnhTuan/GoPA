import { Link2, Plus } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { MacSelect, type MacSelectOption } from "../../../components/design-system/mac-select";
import type { Journal } from "../types";

interface LinkingSectionProps {
  availableToLink: Journal[];
  onLink: (linkedId: string) => void;
}

export function JournalLinkingSection({ availableToLink, onLink }: LinkingSectionProps) {
  const { t } = useTranslation("journal");
  const [linkTargetId, setLinkTargetId] = useState("");

  const handleSubmit = () => {
    if (!linkTargetId) return;
    onLink(linkTargetId);
    setLinkTargetId("");
  };

  const options: MacSelectOption[] = [
    { value: "", label: t("linking.select") },
    ...availableToLink.map((item) => ({
      value: item.id,
      label: `${item.title} (${item.published_date ? item.published_date.slice(0, 10) : item.created_at.slice(0, 10)})`,
    })),
  ];

  return (
    <div className="pt-3 border-t border-slate-200/80 dark:border-white/10 space-y-2 shrink-0">
      <div className="flex items-center justify-between">
        <h4 className="text-xs font-bold text-muted-foreground uppercase tracking-wider flex items-center gap-1.5">
          <Link2 size={13} className="text-amber-500" />
          {t("linking.title")}
        </h4>
      </div>

      <div className="flex items-center gap-2">
        <div className="flex-1 min-w-0">
          <MacSelect
            className="w-full"
            triggerClassName="h-9 text-xs rounded-xl px-3 border-slate-200 dark:border-white/15 bg-white/90 dark:bg-white/5 shadow-2xs font-medium"
            options={options}
            value={linkTargetId}
            onChange={(val) => setLinkTargetId(val)}
            placeholder={t("linking.select")}
          />
        </div>

        <button
          type="button"
          onClick={handleSubmit}
          disabled={!linkTargetId}
          className="h-9 rounded-xl bg-amber-500 px-3.5 text-xs font-bold text-slate-950 hover:bg-amber-400 disabled:opacity-40 transition-all flex items-center gap-1 shrink-0 shadow-2xs"
        >
          <Plus size={14} />
          <span>{t("linking.add")}</span>
        </button>
      </div>
    </div>
  );
}

