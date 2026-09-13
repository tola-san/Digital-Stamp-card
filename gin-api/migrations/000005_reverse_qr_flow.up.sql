RENAME TABLE stamp_qrs TO customer_qr_tokens;

ALTER TABLE customer_qr_tokens
    CHANGE COLUMN staff_id customer_id CHAR(36) CHARACTER SET ascii NOT NULL,
    CHANGE COLUMN used_by_customer_id used_by_staff_id CHAR(36) CHARACTER SET ascii,
    DROP INDEX stamp_qrs_staff_created_idx,
    DROP INDEX stamp_qrs_status_expiry_idx,
    ADD INDEX customer_qr_tokens_customer_created_idx (customer_id, created_at DESC),
    ADD INDEX customer_qr_tokens_status_expiry_idx (status, expires_at),
    ADD CONSTRAINT customer_qr_tokens_customer_fk
        FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    ADD CONSTRAINT customer_qr_tokens_used_by_staff_fk
        FOREIGN KEY (used_by_staff_id) REFERENCES staff(id),
    ADD CONSTRAINT customer_qr_tokens_usage_matches_status_check CHECK (
        (status = 'USED' AND used_at IS NOT NULL AND used_by_staff_id IS NOT NULL)
        OR
        (status <> 'USED' AND used_at IS NULL AND used_by_staff_id IS NULL)
    );

ALTER TABLE stamp_transactions
    CHANGE COLUMN stamp_qr_id customer_qr_token_id CHAR(36) CHARACTER SET ascii,
    ADD CONSTRAINT stamp_transactions_customer_qr_token_fk
        FOREIGN KEY (customer_qr_token_id) REFERENCES customer_qr_tokens(id);
