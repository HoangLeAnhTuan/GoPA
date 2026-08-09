import { Navigate, createBrowserRouter } from "react-router-dom";
import { AppLayout } from "../layouts/app-layout";
import { AuthLayout } from "../layouts/auth-layout";
import { PlaceholderPage } from "../pages/placeholder-page";
import { NotFoundPage } from "../pages/not-found-page";
import { RequireAuth } from "../features/auth/components/require-auth";
import { LoginPage } from "../features/auth/pages/login-page";
import { RegisterPage } from "../features/auth/pages/register-page";
import { TasksPage } from "../features/tasks/pages/tasks-page";
import { LearnPage } from "../features/linguistics/pages/learn-page";
import { FinancePage } from "../features/finance/pages/finance-page";
import { JournalPage } from "../features/journal/pages/journal-page";

import { PomodoroPage } from "../pages/pomodoro-page";

const appPage = (titleKey: string, descriptionKey: string) => <PlaceholderPage titleKey={titleKey} descriptionKey={descriptionKey} />;

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
      { path: "today", element: appPage("navigation.today", "pages.today.description") },
      { path: "tasks", element: <TasksPage /> },
      { path: "pomodoro", element: <PomodoroPage /> },
      { path: "learn", element: <LearnPage /> },
      { path: "learn/review", element: <LearnPage /> },
      { path: "learn/vocabulary", element: <LearnPage /> },
      { path: "finance", element: <FinancePage /> },
      { path: "finance/transactions", element: <FinancePage /> },
      { path: "finance/accounts", element: <FinancePage /> },
      { path: "finance/budgets", element: <FinancePage /> },
      { path: "finance/goals", element: <FinancePage /> },
      { path: "journal", element: <JournalPage /> },
      { path: "journal/:journalId", element: appPage("pages.journalDetail.title", "pages.journalDetail.description") },
      { path: "settings", element: appPage("navigation.settings", "pages.settings.description") },
    ],
  },
  { path: "*", element: <NotFoundPage /> },
]);
