ALTER TABLE stamp_transactions
    DROP FOREIGN KEY stamp_transactions_customer_qr_token_fk;

ALTER TABLE customer_qr_tokens
    DROP FOREIGN KEY customer_qr_tokens_customer_fk,
    DROP FOREIGN KEY customer_qr_tokens_used_by_staff_fk;

RENAME TABLE customer_qr_tokens TO stamp_qrs;

ALTER TABLE stamp_qrs
    CHANGE COLUMN customer_id staff_id CHAR(36) CHARACTER SET ascii NOT NULL,
    CHANGE COLUMN used_by_staff_id used_by_customer_id CHAR(36) CHARACTER SET ascii,
    DROP INDEX customer_qr_tokens_customer_created_idx,
    DROP INDEX customer_qr_tokens_status_expiry_idx,
    ADD INDEX stamp_qrs_staff_created_idx (staff_id, created_at DESC),
    ADD INDEX stamp_qrs_status_expiry_idx (status, expires_at),
    ADD CONSTRAINT stamp_qrs_staff_fk
        FOREIGN KEY (staff_id) REFERENCES staff(id),
    ADD CONSTRAINT stamp_qrs_used_by_customer_fk
        FOREIGN KEY (used_by_customer_id) REFERENCES customers(id);

ALTER TABLE stamp_transactions
    CHANGE COLUMN customer_qr_token_id stamp_qr_id CHAR(36) CHARACTER SET ascii,
    ADD CONSTRAINT stamp_transactions_stamp_qr_fk
        FOREIGN KEY (stamp_qr_id) REFERENCES stamp_qrs(id);
