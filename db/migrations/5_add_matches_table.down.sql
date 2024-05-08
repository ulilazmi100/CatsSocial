DROP TABLE IF EXISTS matches;

DROP INDEX IF EXISTS idx_matches_user_id;
DROP INDEX IF EXISTS idx_matches_match_user_id;
DROP INDEX IF EXISTS idx_matches_match_cat_id;
DROP INDEX IF EXISTS idx_matches_user_cat_id;
DROP INDEX IF EXISTS idx_matches_status_approved;
DROP INDEX IF EXISTS idx_matches_status_removed;
DROP INDEX IF EXISTS idx_matches_created_at;