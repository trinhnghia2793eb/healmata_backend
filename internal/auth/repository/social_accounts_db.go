package repository

import (
	"context"
	"healmata_backend/internal/auth/model"

	"github.com/jackc/pgx/v5"
)

func (r *authRepository) GetSocialAccount(ctx context.Context, provider string, providerUserID string) (*model.SocialAccounts, error) {
	var socialAccount model.SocialAccounts
	query := `SELECT id, user_id, provider, provider_user_id, provider_email, created_at FROM social_accounts WHERE provider = $1 AND provider_user_id = $2`
	err := r.db.QueryRow(ctx, query, provider, providerUserID).Scan(
		&socialAccount.ID,
		&socialAccount.UserID,
		&socialAccount.Provider,
		&socialAccount.ProviderUserID,
		&socialAccount.ProviderEmail,
		&socialAccount.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &socialAccount, nil
}

func (r *authRepository) CreateSocialAccount(ctx context.Context, tx pgx.Tx, socialAcc *CreateSocialAccountPayload) (*model.SocialAccounts, error) {
	query := `INSERT INTO social_accounts (user_id, provider, provider_user_id, provider_email) 
			  VALUES ($1, $2, $3, $4) 
			  RETURNING id, created_at`
	var socialAccount model.SocialAccounts
	err := tx.QueryRow(ctx, query,
		socialAcc.UserID,
		socialAcc.Provider,
		socialAcc.ProviderUserID,
		socialAcc.ProviderEmail,
	).Scan(
		&socialAccount.ID,
		&socialAccount.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	socialAccount.UserID = socialAcc.UserID
	socialAccount.Provider = socialAcc.Provider
	socialAccount.ProviderUserID = socialAcc.ProviderUserID
	socialAccount.ProviderEmail = socialAcc.ProviderEmail
	return &socialAccount, nil
}
