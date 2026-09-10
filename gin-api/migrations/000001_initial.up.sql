CREATE TABLE IF NOT EXISTS customers (
    id CHAR(36) CHARACTER SET ascii PRIMARY KEY DEFAULT (UUID()),
    name VARCHAR(100) NOT NULL CHECK (CHAR_LENGTH(TRIM(name)) BETWEEN 2 AND 100),
    phone VARCHAR(32) NOT NULL UNIQUE,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
);

CREATE TABLE IF NOT EXISTS staff (
    id CHAR(36) CHARACTER SET ascii PRIMARY KEY DEFAULT (UUID()),
    name VARCHAR(100) NOT NULL CHECK (CHAR_LENGTH(TRIM(name)) BETWEEN 2 AND 100),
    email VARCHAR(254) COLLATE utf8mb4_unicode_ci NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
);

CREATE TABLE IF NOT EXISTS rewards (
    id CHAR(36) CHARACTER SET ascii PRIMARY KEY DEFAULT (UUID()),
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    required_stamps INT NOT NULL CHECK (required_stamps > 0),
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
);

CREATE TABLE IF NOT EXISTS stamp_cards (
    id CHAR(36) CHARACTER SET ascii PRIMARY KEY DEFAULT (UUID()),
    customer_id CHAR(36) CHARACTER SET ascii NOT NULL UNIQUE,
    stamp_count INT NOT NULL DEFAULT 0 CHECK (stamp_count >= 0),
    required_stamps INT NOT NULL DEFAULT 10 CHECK (required_stamps > 0),
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS stamp_qrs (
    id CHAR(36) CHARACTER SET ascii PRIMARY KEY DEFAULT (UUID()),
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    staff_id CHAR(36) CHARACTER SET ascii NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'USED', 'EXPIRED', 'CANCELLED')),
    expires_at DATETIME(6) NOT NULL,
    used_at DATETIME(6),
    used_by_customer_id CHAR(36) CHARACTER SET ascii,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    FOREIGN KEY (staff_id) REFERENCES staff(id),
    FOREIGN KEY (used_by_customer_id) REFERENCES customers(id),
    INDEX stamp_qrs_staff_created_idx (staff_id, created_at DESC),
    INDEX stamp_qrs_status_expiry_idx (status, expires_at),
    CONSTRAINT stamp_qrs_usage_matches_status_check CHECK (
        (status = 'USED' AND used_at IS NOT NULL AND used_by_customer_id IS NOT NULL)
        OR
        (status <> 'USED' AND used_at IS NULL AND used_by_customer_id IS NULL)
    )
);

CREATE TABLE IF NOT EXISTS stamp_transactions (
    id CHAR(36) CHARACTER SET ascii PRIMARY KEY DEFAULT (UUID()),
    customer_id CHAR(36) CHARACTER SET ascii NOT NULL,
    staff_id CHAR(36) CHARACTER SET ascii,
    reward_id CHAR(36) CHARACTER SET ascii,
    stamp_qr_id CHAR(36) CHARACTER SET ascii UNIQUE,
    type VARCHAR(32) NOT NULL CHECK (type IN ('STAMP_ADDED', 'STAMP_REVERSED', 'REWARD_REDEEMED')),
    stamp_delta INT NOT NULL,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    FOREIGN KEY (customer_id) REFERENCES customers(id),
    FOREIGN KEY (staff_id) REFERENCES staff(id),
    FOREIGN KEY (reward_id) REFERENCES rewards(id),
    FOREIGN KEY (stamp_qr_id) REFERENCES stamp_qrs(id),
    INDEX stamp_transactions_customer_created_idx (customer_id, created_at DESC),
    CONSTRAINT stamp_transactions_type_relationships_check CHECK (
        (
            type = 'STAMP_ADDED'
            AND stamp_delta > 0
            AND reward_id IS NULL
        )
        OR
        (
            type = 'STAMP_REVERSED'
            AND stamp_delta < 0
            AND reward_id IS NULL
            AND stamp_qr_id IS NULL
        )
        OR
        (
            type = 'REWARD_REDEEMED'
            AND stamp_delta < 0
            AND reward_id IS NOT NULL
            AND stamp_qr_id IS NULL
        )
    )
);
