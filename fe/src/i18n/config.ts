import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import enCommon from "../locales/en/common.json";
import viCommon from "../locales/vi/common.json";
import jaCommon from "../locales/ja/common.json";
import enTasks from "../locales/en/tasks.json";
import viTasks from "../locales/vi/tasks.json";
import jaTasks from "../locales/ja/tasks.json";
import enLearning from "../locales/en/learning.json";
import viLearning from "../locales/vi/learning.json";
import jaLearning from "../locales/ja/learning.json";
import enAssets from "../locales/en/assets.json";
import viAssets from "../locales/vi/assets.json";
import jaAssets from "../locales/ja/assets.json";
import enFinance from "../locales/en/finance.json";
import viFinance from "../locales/vi/finance.json";
import jaFinance from "../locales/ja/finance.json";
import enJournal from "../locales/en/journal.json";
import viJournal from "../locales/vi/journal.json";
import jaJournal from "../locales/ja/journal.json";
import enFeedback from "../locales/en/feedback.json";
import viFeedback from "../locales/vi/feedback.json";
import jaFeedback from "../locales/ja/feedback.json";

void i18n.use(initReactI18next).init({
  resources: {
    en: { common: enCommon, tasks: enTasks, learning: enLearning, assets: enAssets, finance: enFinance, journal: enJournal, feedback: enFeedback },
    vi: { common: viCommon, tasks: viTasks, learning: viLearning, assets: viAssets, finance: viFinance, journal: viJournal, feedback: viFeedback },
    ja: { common: jaCommon, tasks: jaTasks, learning: jaLearning, assets: jaAssets, finance: jaFinance, journal: jaJournal, feedback: jaFeedback },
  },
  lng: "vi",
  fallbackLng: "en",
  defaultNS: "common",
  interpolation: { escapeValue: false },
});

export default i18n;
