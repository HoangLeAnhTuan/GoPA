CREATE TABLE accounts (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL CHECK (char_length(trim(name)) BETWEEN 1 AND 120),
    type TEXT NOT NULL CHECK (type IN ('CASH', 'BANK', 'CREDIT_CARD', 'SAVINGS', 'INVESTMENT', 'CRYPTO')),
    currency CHAR(3) NOT NULL CHECK (currency ~ '^[A-Z]{3}$'),
    initial_balance NUMERIC(14, 2) NOT NULL CHECK (initial_balance >= 0),
    current_balance NUMERIC(14, 2) NOT NULL,
    color TEXT NOT NULL DEFAULT '#10B981' CHECK (color ~ '^#[0-9A-Fa-f]{6}$'),
    icon TEXT NOT NULL DEFAULT 'wallet' CHECK (char_length(trim(icon)) BETWEEN 1 AND 50),
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX accounts_user_active_idx ON accounts (user_id, is_archived, created_at DESC);

CREATE TABLE categories (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL CHECK (char_length(trim(name)) BETWEEN 1 AND 100),
    type TEXT NOT NULL CHECK (type IN ('INCOME', 'EXPENSE')),
    icon TEXT NOT NULL DEFAULT 'tag' CHECK (char_length(trim(icon)) BETWEEN 1 AND 50),
    color TEXT NOT NULL DEFAULT '#64748B' CHECK (color ~ '^#[0-9A-Fa-f]{6}$'),
    parent_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    monthly_budget NUMERIC(14, 2) CHECK (monthly_budget > 0),
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    UNIQUE (user_id, type, name),
    CHECK (parent_id IS NULL OR parent_id <> id)
);

CREATE INDEX categories_user_type_idx ON categories (user_id, type, name);

CREATE TABLE transactions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
    to_account_id UUID REFERENCES accounts(id) ON DELETE RESTRICT,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    type TEXT NOT NULL CHECK (type IN ('INCOME', 'EXPENSE', 'TRANSFER')),
    amount NUMERIC(14, 2) NOT NULL CHECK (amount > 0),
    currency CHAR(3) NOT NULL CHECK (currency ~ '^[A-Z]{3}$'),
    exchange_rate NUMERIC(10, 6) NOT NULL DEFAULT 1 CHECK (exchange_rate > 0),
    description TEXT NOT NULL DEFAULT '' CHECK (char_length(description) <= 500),
    merchant TEXT CHECK (char_length(merchant) <= 150),
    tags TEXT[] NOT NULL DEFAULT '{}',
    occurred_at TIMESTAMPTZ NOT NULL,
    is_recurring BOOLEAN NOT NULL DEFAULT FALSE,
    recurring_rule JSONB,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    CHECK ((type = 'TRANSFER' AND to_account_id IS NOT NULL AND category_id IS NULL AND to_account_id <> account_id) OR (type <> 'TRANSFER' AND to_account_id IS NULL))
);

CREATE INDEX idx_transactions_user_occurred ON transactions (user_id, occurred_at DESC);
CREATE INDEX idx_transactions_account ON transactions (account_id, occurred_at DESC);
CREATE INDEX transactions_category_idx ON transactions (category_id, occurred_at DESC);

CREATE TABLE budgets (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    name TEXT NOT NULL CHECK (char_length(trim(name)) BETWEEN 1 AND 150),
    amount NUMERIC(14, 2) NOT NULL CHECK (amount > 0),
    period TEXT NOT NULL CHECK (period IN ('WEEKLY', 'MONTHLY', 'QUARTERLY', 'YEARLY', 'CUSTOM')),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL CHECK (end_date >= start_date),
    alert_threshold NUMERIC(4, 2) NOT NULL DEFAULT 0.80 CHECK (alert_threshold > 0 AND alert_threshold <= 1),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX budgets_user_active_idx ON budgets (user_id, is_active, start_date, end_date);

CREATE TABLE savings_goals (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL CHECK (char_length(trim(name)) BETWEEN 1 AND 150),
    target_amount NUMERIC(14, 2) NOT NULL CHECK (target_amount > 0),
    current_amount NUMERIC(14, 2) NOT NULL DEFAULT 0 CHECK (current_amount >= 0),
    linked_account_id UUID REFERENCES accounts(id) ON DELETE SET NULL,
    target_date DATE,
    color TEXT NOT NULL DEFAULT '#10B981' CHECK (color ~ '^#[0-9A-Fa-f]{6}$'),
    icon TEXT NOT NULL DEFAULT 'target' CHECK (char_length(trim(icon)) BETWEEN 1 AND 50),
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX savings_goals_user_idx ON savings_goals (user_id, created_at DESC);
