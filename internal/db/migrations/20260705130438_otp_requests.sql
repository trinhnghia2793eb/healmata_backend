-- +goose Up
-- +goose StatementBegin
CREATE TABLE otp_requests (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	identifier VARCHAR(255) NOT NULL,
	otp_hash VARCHAR(255) NOT NULL,
	purpose VARCHAR(255) NOT NULL,
	attempts INT DEFAULT 0,
	expires_at timestamptz NOT NULL,
	verified_at timestamptz,
	created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX otp_requests_identifier ON otp_requests (identifier);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE otp_requests;
-- +goose StatementEnd
