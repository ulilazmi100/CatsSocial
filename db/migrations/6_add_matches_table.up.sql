CREATE TABLE matches (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    match_user_id INT REFERENCES users(id),
    match_cat_id INT REFERENCES cats(id),
    user_cat_id INT REFERENCES cats(id),
    message TEXT,
    status VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);