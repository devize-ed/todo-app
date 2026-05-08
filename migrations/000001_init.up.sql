--migrations/000001_init.up.sql
CREATE SCHEMA IF NOT EXISTS todoapp;

--Create table for Users
CREATE TABLE IF NOT EXISTS todoapp.users (
    id          UUID                  PRIMARY KEY
,   name        TEXT         NOT NULL CHECK (char_length(trim(name)) >= 3 AND char_length(trim(name)) <= 100)
,   email       TEXT                  UNIQUE CHECK (char_length(trim(email)) >= 6 AND char_length(trim(email)) <= 255)
,   created_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
,   updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
);

-- Create table for tasks
CREATE TABLE IF NOT EXISTS todoapp.tasks (
    id           UUID                   PRIMARY KEY
,   user_id      UUID         NOT NULL  REFERENCES todoapp.users(id) ON DELETE CASCADE
,   title        TEXT         NOT NULL  CHECK (char_length(trim(title)) >= 1 AND char_length(trim(title)) <= 100)
,   description  TEXT                   CHECK (char_length(trim(description)) >= 1 AND char_length(trim(description)) <= 1000)
,   completed    BOOLEAN      NOT NULL  DEFAULT FALSE
,   completed_at TIMESTAMPTZ            
,   created_at   TIMESTAMPTZ  NOT NULL  DEFAULT now()
,   updated_at   TIMESTAMPTZ  NOT NULL  DEFAULT now()

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