-- 002_user_repo_many_to_many.sql

-- First, remove the user_id column and its unique constraint from repositories
ALTER TABLE repositories DROP CONSTRAINT repositories_user_id_github_repo_id_key;
ALTER TABLE repositories DROP COLUMN user_id;

-- Add a unique constraint on github_repo_id so we can share repository records
ALTER TABLE repositories ADD CONSTRAINT repositories_github_repo_id_key UNIQUE (github_repo_id);

-- Create the join table for many-to-many relationship between users and repositories
CREATE TABLE user_repository (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    repository_id UUID NOT NULL REFERENCES repositories(id) ON DELETE CASCADE,
    added_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, repository_id)
);

-- Cleanup legacy users table fields that are no longer needed
ALTER TABLE users DROP COLUMN github_id;
ALTER TABLE users DROP COLUMN oauth_token;
ALTER TABLE users DROP COLUMN username;
