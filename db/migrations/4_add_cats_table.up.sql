CREATE TABLE cats (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    name VARCHAR(255) NOT NULL,
    race VARCHAR(255),
    sex VARCHAR(10),
    age_in_month INT,
    description TEXT,
    image_urls TEXT[],
    has_matched BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
