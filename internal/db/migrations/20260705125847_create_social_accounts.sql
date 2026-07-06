-- +goose Up
-- +goose StatementBegin
CREATE TABLE social_accounts (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id uuid NOT NULL,
	provider VARCHAR(255) NOT NULL,
	provider_user_id VARCHAR(255) NOT NULL,
	provider_email VARCHAR(255) NOT NULL,
	created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	
	CONSTRAINT fk_social_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	CONSTRAINT social_accounts_provider_user_unique UNIQUE (provider, provider_user_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE social_accounts;
-- +goose StatementEnd
