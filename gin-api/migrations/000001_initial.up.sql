CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL CHECK (char_length(trim(name)) BETWEEN 2 AND 100),
    phone TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS staff (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL CHECK (char_length(trim(name)) BETWEEN 2 AND 100),
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS rewards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    required_stamps INTEGER NOT NULL CHECK (required_stamps > 0),
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS stamp_cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL UNIQUE REFERENCES customers(id) ON DELETE CASCADE,
    stamp_count INTEGER NOT NULL DEFAULT 0 CHECK (stamp_count >= 0),
    required_stamps INTEGER NOT NULL DEFAULT 10 CHECK (required_stamps > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS stamp_qrs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token_hash TEXT NOT NULL UNIQUE,
    staff_id UUID NOT NULL REFERENCES staff(id),
    status TEXT NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'USED', 'EXPIRED', 'CANCELLED')),
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    used_by_customer_id UUID REFERENCES customers(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK ((status = 'USED' AND used_at IS NOT NULL AND used_by_customer_id IS NOT NULL) OR status <> 'USED')
);

CREATE TABLE IF NOT EXISTS stamp_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES customers(id),
    staff_id UUID REFERENCES staff(id),
    reward_id UUID REFERENCES rewards(id),
    stamp_qr_id UUID UNIQUE REFERENCES stamp_qrs(id),
    type TEXT NOT NULL CHECK (type IN ('STAMP_ADDED', 'STAMP_REVERSED', 'REWARD_REDEEMED')),
    stamp_delta INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS stamp_transactions_customer_created_idx ON stamp_transactions (customer_id, created_at DESC);
CREATE INDEX IF NOT EXISTS stamp_qrs_staff_created_idx ON stamp_qrs (staff_id, created_at DESC);
CREATE INDEX IF NOT EXISTS stamp_qrs_active_expiry_idx ON stamp_qrs (expires_at) WHERE status = 'ACTIVE';
