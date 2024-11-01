-- +goose Up
CREATE TABLE "users" (
    id uuid DEFAULT gen_random_uuid(),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    email VARCHAR(255) NOT NULL
);

-- +goose Down
DROP TABLE users;
