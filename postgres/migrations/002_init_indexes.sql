CREATE INDEX idx_team ON teams(name);

CREATE INDEX idx_users_team_active ON users(team_name, is_active);

CREATE INDEX idx_users_pull_requests_user ON users_pull_requests(users_id);