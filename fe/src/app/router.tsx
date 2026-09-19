import { lazy, Suspense, type ReactNode } from "react";
import { Navigate, createBrowserRouter } from "react-router-dom";
import { AppLayout } from "../layouts/app-layout";
import { AuthLayout } from "../layouts/auth-layout";
import { NotFoundPage } from "../pages/not-found-page";
import { RequireAuth } from "../features/auth/components/require-auth";
import { LoadingState } from "../components/feedback/loading-state";

const LoginPage = lazy(async () => ({ default: (await import("../features/auth/pages/login-page")).LoginPage }));
const RegisterPage = lazy(async () => ({ default: (await import("../features/auth/pages/register-page")).RegisterPage }));
const TasksPage = lazy(async () => ({ default: (await import("../features/tasks/pages/tasks-page")).TasksPage }));
const LearnPage = lazy(async () => ({ default: (await import("../features/linguistics/pages/learn-page")).LearnPage }));
const ReviewSessionPage = lazy(async () => ({ default: (await import("../features/linguistics/pages/review-session-page")).ReviewSessionPage }));
const VocabularyManagePage = lazy(async () => ({ default: (await import("../features/linguistics/pages/vocabulary-manage-page")).VocabularyManagePage }));
const FinancePage = lazy(async () => ({ default: (await import("../features/finance/pages/finance-page")).FinancePage }));
const BudgetsPage = lazy(async () => ({ default: (await import("../features/finance/pages/budgets-page")).BudgetsPage }));
const GoalsPage = lazy(async () => ({ default: (await import("../features/finance/pages/goals-page")).GoalsPage }));
const JournalPage = lazy(async () => ({ default: (await import("../features/journal/pages/journal-page")).JournalPage }));
const JournalDetailPage = lazy(async () => ({ default: (await import("../features/journal/pages/journal-detail-page")).JournalDetailPage }));
const PomodoroPage = lazy(async () => ({ default: (await import("../pages/pomodoro-page")).PomodoroPage }));
const TodayPage = lazy(async () => ({ default: (await import("../features/today/pages/today-page")).TodayPage }));
const SettingsPage = lazy(async () => ({ default: (await import("../features/settings/pages/settings-page")).SettingsPage }));

function deferred(page: ReactNode) {
  return <Suspense fallback={<LoadingState />}>{page}</Suspense>;
}

export const router = createBrowserRouter([
  { path: "/", element: <Navigate to="/app/today" replace /> },
  {
    element: <AuthLayout />,
    children: [
      { path: "/login", element: deferred(<LoginPage />) },
      { path: "/register", element: deferred(<RegisterPage />) },
    ],
  },
  {
    path: "/app",
    element: <RequireAuth><AppLayout /></RequireAuth>,
    children: [
      { index: true, element: <Navigate to="today" replace /> },
      { path: "today", element: deferred(<TodayPage />) },
      { path: "tasks", element: deferred(<TasksPage />) },
      { path: "pomodoro", element: deferred(<PomodoroPage />) },
      { path: "learn", element: deferred(<LearnPage />) },
      { path: "learn/review", element: deferred(<ReviewSessionPage />) },
      { path: "learn/vocabulary", element: deferred(<VocabularyManagePage />) },
      { path: "finance", element: deferred(<FinancePage />) },
      { path: "finance/transactions", element: deferred(<FinancePage />) },
      { path: "finance/accounts", element: deferred(<FinancePage />) },
      { path: "finance/budgets", element: deferred(<BudgetsPage />) },
      { path: "finance/goals", element: deferred(<GoalsPage />) },
      { path: "journal", element: deferred(<JournalPage />) },
      { path: "journal/:journalId", element: deferred(<JournalDetailPage />) },
      { path: "settings", element: deferred(<SettingsPage />) },
    ],
  },
  { path: "*", element: <NotFoundPage /> },
]);
