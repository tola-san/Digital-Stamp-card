CREATE TABLE IF NOT EXISTS customer_sessions (
    id CHAR(36) CHARACTER SET ascii PRIMARY KEY DEFAULT (UUID()),
    customer_id CHAR(36) CHARACTER SET ascii NOT NULL,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    expires_at DATETIME(6) NOT NULL,
    revoked_at DATETIME(6),
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    INDEX customer_sessions_customer_idx (customer_id),
    INDEX customer_sessions_active_expiry_idx (revoked_at, expires_at),
    CHECK (expires_at > created_at),
    CHECK (revoked_at IS NULL OR revoked_at >= created_at)
);

CREATE TABLE IF NOT EXISTS staff_sessions (
    id CHAR(36) CHARACTER SET ascii PRIMARY KEY DEFAULT (UUID()),
    staff_id CHAR(36) CHARACTER SET ascii NOT NULL,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    expires_at DATETIME(6) NOT NULL,
    revoked_at DATETIME(6),
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    FOREIGN KEY (staff_id) REFERENCES staff(id) ON DELETE CASCADE,
    INDEX staff_sessions_staff_idx (staff_id),
    INDEX staff_sessions_active_expiry_idx (revoked_at, expires_at),
    CHECK (expires_at > created_at),
    CHECK (revoked_at IS NULL OR revoked_at >= created_at)
);
