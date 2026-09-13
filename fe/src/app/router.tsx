import { Navigate, createBrowserRouter } from "react-router-dom";
import { AppLayout } from "../layouts/app-layout";
import { AuthLayout } from "../layouts/auth-layout";
import { NotFoundPage } from "../pages/not-found-page";
import { RequireAuth } from "../features/auth/components/require-auth";
import { LoginPage } from "../features/auth/pages/login-page";
import { RegisterPage } from "../features/auth/pages/register-page";
import { TasksPage } from "../features/tasks/pages/tasks-page";
import { LearnPage } from "../features/linguistics/pages/learn-page";
import { ReviewSessionPage } from "../features/linguistics/pages/review-session-page";
import { VocabularyManagePage } from "../features/linguistics/pages/vocabulary-manage-page";
import { FinancePage } from "../features/finance/pages/finance-page";
import { BudgetsPage } from "../features/finance/pages/budgets-page";
import { GoalsPage } from "../features/finance/pages/goals-page";
import { JournalPage } from "../features/journal/pages/journal-page";
import { JournalDetailPage } from "../features/journal/pages/journal-detail-page";

import { PomodoroPage } from "../pages/pomodoro-page";
import { TodayPage } from "../features/today/pages/today-page";
import { SettingsPage } from "../features/settings/pages/settings-page";

export const router = createBrowserRouter([
  { path: "/", element: <Navigate to="/app/today" replace /> },
  {
    element: <AuthLayout />,
    children: [
      { path: "/login", element: <LoginPage /> },
      { path: "/register", element: <RegisterPage /> },
    ],
  },
  {
    path: "/app",
    element: <RequireAuth><AppLayout /></RequireAuth>,
    children: [
      { index: true, element: <Navigate to="today" replace /> },
      { path: "today", element: <TodayPage /> },
      { path: "tasks", element: <TasksPage /> },
      { path: "pomodoro", element: <PomodoroPage /> },
      { path: "learn", element: <LearnPage /> },
      { path: "learn/review", element: <ReviewSessionPage /> },
      { path: "learn/vocabulary", element: <VocabularyManagePage /> },
      { path: "finance", element: <FinancePage /> },
      { path: "finance/transactions", element: <FinancePage /> },
      { path: "finance/accounts", element: <FinancePage /> },
      { path: "finance/budgets", element: <BudgetsPage /> },
      { path: "finance/goals", element: <GoalsPage /> },
      { path: "journal", element: <JournalPage /> },
      { path: "journal/:journalId", element: <JournalDetailPage /> },
      { path: "settings", element: <SettingsPage /> },
    ],
  },
  { path: "*", element: <NotFoundPage /> },
]);
