CREATE TABLE matches (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    match_user_id INT REFERENCES users(id),
    match_cat_id INT REFERENCES cats(id),
    user_cat_id INT REFERENCES cats(id),
    message TEXT,
    status VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_matches_user_id ON matches(user_id);
CREATE INDEX idx_matches_match_user_id ON matches(match_user_id);
CREATE INDEX idx_matches_match_cat_id ON matches(match_cat_id);
CREATE INDEX idx_matches_user_cat_id ON matches(user_cat_id);
CREATE INDEX idx_matches_status_approved ON matches(status) WHERE status = 'approved';
CREATE INDEX idx_matches_status_removed ON matches(status) WHERE status = 'removed';
CREATE INDEX idx_matches_created_at ON matches(created_at);