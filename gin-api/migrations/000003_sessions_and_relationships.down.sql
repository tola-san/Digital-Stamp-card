ALTER TABLE stamp_transactions
    DROP CONSTRAINT IF EXISTS stamp_transactions_type_relationships_check;

ALTER TABLE stamp_qrs
    DROP CONSTRAINT IF EXISTS stamp_qrs_usage_matches_status_check;

DROP INDEX IF EXISTS staff_email_case_insensitive_uidx;
DROP TABLE IF EXISTS staff_sessions;
DROP TABLE IF EXISTS customer_sessions;
