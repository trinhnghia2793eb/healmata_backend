-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_sessions (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id uuid NOT NULL,
	refresh_token_id uuid NOT NULL,
	device_id VARCHAR(255) NOT NULL,
	platform VARCHAR(255) NOT NULL,
	ip_address VARCHAR(255) NOT NULL,
	user_agent VARCHAR(255) NOT NULL,
	last_active_at timestamptz,
	revoked_at timestamptz,
	created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	
	CONSTRAINT fk_session_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	CONSTRAINT fk_session_tokens FOREIGN KEY (refresh_token_id) REFERENCES refresh_tokens(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_sessions;
-- +goose StatementEnd
