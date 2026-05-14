-- Migration from old schema to new schema
-- This script will be run against the new database after creating it from schema.sql

-- Tables that need special handling:
-- 1. users: old has user_id (TEXT), new has external_id (TEXT)
-- 2. billettholdere: old has inserted_time, new has created_at, updated_at, created_by_id, updated_by_id
-- 3. relation_billettholder_emails: old is billettholder_emails, new has created_by_id, updated_by_id
-- 4. events: old has host/image_url/pulje_name, new has created_by_id, updated_by_id, status_changed_by_id, etc.
-- 5. relation_event_puljer: old is event_pujer, column names changed
-- 6. relation_events_players: old is events_players, structure changed
-- 7. interests: old has just interest_level column, new adds created_by_id, updated_by_id, pulje_id
-- 8. puljer: old has start_time/end_time, new has start_at/end_at, status

-- This will be a manual migration since we need to handle old database separately
