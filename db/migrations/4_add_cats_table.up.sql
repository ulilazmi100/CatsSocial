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

CREATE INDEX idx_cats_user_id ON cats(user_id);
CREATE INDEX idx_cats_has_matched ON cats(has_matched);
CREATE INDEX idx_cats_race_sex_age ON cats(race, sex, age_in_month);
CREATE INDEX idx_cats_created_at ON cats(created_at);