-- +goose Up
CREATE TABLE "chirps" (
    id uuid DEFAULT gen_random_uuid(),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    body VARCHAR(140) NOT NULL,
    user_id uuid REFERENCES users (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE chirps;
