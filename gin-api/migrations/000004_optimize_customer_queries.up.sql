CREATE INDEX stamp_transactions_customer_page_idx
    ON stamp_transactions (customer_id, created_at DESC, id DESC);

DROP INDEX stamp_transactions_customer_created_idx
    ON stamp_transactions;
