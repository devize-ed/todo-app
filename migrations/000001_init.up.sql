--migrations/000001_init.up.sql
CREATE SCHEMA IF NOT EXISTS todoapp;

--Create table for Users
CREATE TABLE IF NOT EXISTS todoapp.users (
    id          UUID                  PRIMARY KEY
,   name        TEXT         NOT NULL CHECK (char_length(name) >= 3)
,   email       TEXT                  CHECK (char_length(email) >= 6)
,   created_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
);

-- Create table for tasks
CREATE TABLE IF NOT EXISTS todoapp.tasks (
    id           UUID                   PRIMARY KEY
,   user_id      UUID         NOT NULL  REFERENCES todoapp.users(id) ON DELETE CASCADE
,   title        TEXT         NOT NULL  CHECK (char_length(title) >= 3)
,   description  TEXT
,   completed    BOOLEAN      NOT NULL  DEFAULT FALSE
,   completed_at TIMESTAMPTZ            
,   created_at   TIMESTAMPTZ  NOT NULL  DEFAULT now()

,   CHECK (
        (NOT completed AND completed_at IS NULL)
        OR (completed = TRUE AND completed_at IS NOT NULL AND completed_at >= created_at)
    )
);



CREATE INDEX IF NOT EXISTS tasks_user_id_idx
    ON todoapp.tasks (user_id);
CREATE INDEX IF NOT EXISTS tasks_user_created_at_idx
    ON todoapp.tasks (user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS tasks_user_open_idx
    ON todoapp.tasks (user_id)
    WHERE completed = FALSE;
CREATE INDEX IF NOT EXISTS tasks_user_completed_at_idx
    ON todoapp.tasks (user_id, completed_at DESC)
    WHERE completed = TRUE;