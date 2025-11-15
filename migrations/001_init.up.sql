CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    username VARCHAR(256) UNIQUE NOT NULL,
    team_id UUID NOT NULL,
    is_active BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS teams (
    id UUID PRIMARY KEY,
    name VARCHAR(256) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS statuses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(32) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS pull_requests (
    id UUID PRIMARY KEY,
    title VARCHAR(256) NOT NULL,
    author_id UUID NOT NULL,
    status_id INT NOT NULL REFERENCES statuses(id),
    need_more_reviewers BOOLEAN NOT NULL,
    merged_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS pull_requests_users (
    pr_id UUID NOT NULL REFERENCES pull_requests(id),
    user_id UUID NOT NULL REFERENCES users(id),
    PRIMARY KEY (pr_id, user_id)
);

CREATE INDEX IF NOT EXISTS pull_requests_users_pr_id_idx ON pull_requests_users (pr_id);
CREATE INDEX IF NOT EXISTS pull_requests_users_user_id_idx ON pull_requests_users (user_id);
CREATE INDEX IF NOT EXISTS users_team_id_idx ON users (team_id);

INSERT INTO statuses (name) VALUES ('OPEN'), ('MERGED');