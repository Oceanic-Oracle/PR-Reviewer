CREATE TABLE teams (
    name TEXT PRIMARY KEY
);

CREATE TABLE users (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    team_name TEXT NOT NULL REFERENCES teams(name) ON DELETE CASCADE,
    is_active BOOLEAN NOT NULL DEFAULT true
);

CREATE TABLE pull_requests (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    author_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status TEXT NOT NULL CHECK (status IN ('OPEN', 'MERGED')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    merged_at TIMESTAMPTZ
);

CREATE TABLE users_pull_requests (
    pull_requests_id TEXT NOT NULL REFERENCES pull_requests(id) ON DELETE CASCADE,
    users_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE
);