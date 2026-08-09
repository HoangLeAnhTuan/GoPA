import { AnimatePresence, motion } from "framer-motion";
import { Check, ChevronDown } from "lucide-react";
import type { ComponentType } from "react";
import { useEffect, useRef, useState } from "react";
import { createPortal } from "react-dom";
import { cn } from "../../lib/cn";

export interface MacSelectOption {
  value: string;
  label: string;
  icon?: ComponentType<{ size?: number | string; className?: string }>;
  description?: string;
}

export interface MacSelectProps {
  options: MacSelectOption[];
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  className?: string;
  triggerClassName?: string;
  "aria-label"?: string;
  disabled?: boolean;
}

export function MacSelect({
  options,
  value,
  onChange,
  placeholder = "Select...",
  className,
  triggerClassName,
  "aria-label": ariaLabel,
  disabled = false,
}: MacSelectProps) {
  const [open, setOpen] = useState(false);
  const triggerRef = useRef<HTMLButtonElement>(null);
  const menuRef = useRef<HTMLDivElement>(null);
  const [coords, setCoords] = useState<{ top: number; left: number; width: number; flipUp: boolean }>({
    top: 0,
    left: 0,
    width: 160,
    flipUp: false,
  });

  const selectedOption = options.find((opt) => opt.value === value);

  const updatePosition = () => {
    if (!triggerRef.current) return;
    const rect = triggerRef.current.getBoundingClientRect();
    const spaceBelow = window.innerHeight - rect.bottom;
    const flipUp = spaceBelow < 240 && rect.top > 240;
    const width = Math.min(Math.max(160, rect.width), window.innerWidth - 16);

    setCoords({
      top: flipUp ? rect.top - 6 : rect.bottom + 6,
      left: Math.max(8, Math.min(rect.left, window.innerWidth - width - 8)),
      width,
      flipUp,
    });
  };

  const handleToggle = () => {
    if (disabled) return;
    if (!open) {
      updatePosition();
    }
    setOpen((prev) => !prev);
  };

  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      const target = event.target as Node;
      if (
        triggerRef.current &&
        !triggerRef.current.contains(target) &&
        menuRef.current &&
        !menuRef.current.contains(target)
      ) {
        setOpen(false);
      }
    }

    function handleScrollOrResize() {
      if (open) {
        updatePosition();
      }
    }

    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === "Escape") {
        setOpen(false);
      }
    }

    if (open) {
      document.addEventListener("mousedown", handleClickOutside);
      document.addEventListener("keydown", handleKeyDown);
      window.addEventListener("scroll", handleScrollOrResize, true);
      window.addEventListener("resize", handleScrollOrResize);
    }

    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
      document.removeEventListener("keydown", handleKeyDown);
      window.removeEventListener("scroll", handleScrollOrResize, true);
      window.removeEventListener("resize", handleScrollOrResize);
    };
  }, [open]);

  return (
    <div className={cn("relative inline-block text-left", className)}>
      <button
        aria-expanded={open}
        aria-haspopup="listbox"
        aria-label={ariaLabel || placeholder}
        className={cn(
          "group flex min-h-10 w-full items-center justify-between gap-3 rounded-xl border border-white/40 bg-white/70 px-3.5 text-sm font-medium text-foreground shadow-sm backdrop-blur-xl transition hover:border-white/60 hover:bg-white/80 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary dark:border-white/10 dark:bg-slate-900/70 dark:hover:border-white/20 dark:hover:bg-slate-900/90",
          disabled && "cursor-not-allowed opacity-50",
          triggerClassName
        )}
        disabled={disabled}
        onClick={handleToggle}
        ref={triggerRef}
        type="button"
      >
        <span className="flex items-center gap-2 truncate">
          {selectedOption?.icon && (
            <selectedOption.icon className="size-4 shrink-0 text-muted-foreground transition group-hover:text-foreground" />
          )}
          <span className="truncate">{selectedOption ? selectedOption.label : placeholder}</span>
        </span>
        <ChevronDown
          className={cn("size-4 shrink-0 text-muted-foreground transition-transform duration-200", open && "rotate-180")}
        />
      </button>

      {typeof document !== "undefined" &&
        createPortal(
          <AnimatePresence>
            {open && (
              <motion.div
                animate={{ opacity: 1, scale: 1, y: 0 }}
                className="fixed z-[9999] overflow-hidden rounded-2xl border border-white/50 bg-white/95 p-1.5 shadow-2xl shadow-slate-950/30 backdrop-blur-2xl dark:border-white/15 dark:bg-slate-900/95"
                exit={{ opacity: 0, scale: 0.95, y: coords.flipUp ? 6 : -6 }}
                initial={{ opacity: 0, scale: 0.95, y: coords.flipUp ? 6 : -6 }}
                ref={menuRef}
                role="listbox"
                style={{
                  top: coords.flipUp ? undefined : `${coords.top}px`,
                  bottom: coords.flipUp ? `${window.innerHeight - coords.top}px` : undefined,
                  left: `${coords.left}px`,
                  width: `${coords.width}px`,
                }}
                transition={{ type: "spring", stiffness: 450, damping: 30 }}
              >
                <div className="max-h-60 overflow-y-auto space-y-0.5 scrollbar-thin">
                  {options.map((option) => {
                    const isSelected = option.value === value;
                    const IconComponent = option.icon;

                    return (
                      <button
                        aria-selected={isSelected}
                        className={cn(
                          "flex w-full items-center justify-between rounded-xl px-3 py-2 text-left text-sm transition-colors",
                          isSelected
                            ? "bg-primary text-primary-foreground font-semibold"
                            : "text-foreground hover:bg-slate-200/70 dark:hover:bg-slate-800/80"
                        )}
                        key={option.value}
                        onClick={() => {
                          onChange(option.value);
                          setOpen(false);
                        }}
                        role="option"
                        type="button"
                      >
                        <div className="flex items-center gap-2.5 truncate">
                          {IconComponent && (
                            <IconComponent
                              className={cn(
                                "size-4 shrink-0",
                                isSelected ? "text-primary-foreground" : "text-muted-foreground"
                              )}
                            />
                          )}
                          <div className="truncate">
                            <div>{option.label}</div>
                            {option.description && (
                              <div
                                className={cn(
                                  "text-xs font-normal",
                                  isSelected ? "text-primary-foreground/80" : "text-muted-foreground"
                                )}
                              >
                                {option.description}
                              </div>
                            )}
                          </div>
                        </div>
                        {isSelected && <Check className="size-4 shrink-0 text-primary-foreground" />}
                      </button>
                    );
                  })}
                </div>
              </motion.div>
            )}
          </AnimatePresence>,
          document.body
        )}
    </div>
  );
}
