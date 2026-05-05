DROP INDEX IF EXISTS tasks_user_id_idx;
DROP INDEX IF EXISTS tasks_user_created_at_idx;
DROP INDEX IF EXISTS tasks_user_open_idx;
DROP INDEX IF EXISTS tasks_user_completed_at_idx;

DROP TABLE IF EXISTS todoapp.tasks;
DROP TABLE IF EXISTS todoapp.users;
DROP SCHEMA IF EXISTS todoapp;