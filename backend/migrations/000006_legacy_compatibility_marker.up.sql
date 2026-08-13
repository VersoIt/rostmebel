-- No-op compatibility migration.
-- Deployed databases already recorded version 6, but the historical SQL files are no longer in the repository.
-- Keeping version 6 in the migration source lets golang-migrate continue from the live state without forcing
-- manual edits to schema_migrations or risky production resets.
SELECT 1;
