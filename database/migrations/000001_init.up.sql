CREATE SCHEMA IF NOT EXISTS schema_subscriptions;

CREATE TABLE IF NOT EXISTS schema_subscriptions.subscriptions
(
    service_name TEXT,
    price INT NOT NULL,
    user_id TEXT,
    start_time TEXT NOT NULL,
    end_time TEXT
);