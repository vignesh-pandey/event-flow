CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    email_address VARCHAR(255) NOT NULL,
    created_at VARCHAR(255),
    merged_at VARCHAR(255),
    deleted_at VARCHAR(255),
    parent_user_id FLOAT
);
