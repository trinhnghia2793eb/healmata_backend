-- +goose Up
-- +goose StatementBegin
CREATE TYPE user_status AS ENUM ('active', 'disabled', 'deleted');

CREATE TABLE users (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	full_name VARCHAR(255) NOT NULL,
	email VARCHAR(255),
	phone VARCHAR(255) UNIQUE,
	password_hash VARCHAR(255),
	status user_status DEFAULT 'active',
	first_setup_completed BOOLEAN DEFAULT FALSE,
	created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX users_email_unique_idx ON users (LOWER(email));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
DROP TYPE user_status;
-- +goose StatementEnd
