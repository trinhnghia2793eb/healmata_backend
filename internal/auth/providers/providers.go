package providers

import (
	"context"
	authError "healmata_backend/internal/auth/errors"
	"os"
	"strings"

	"github.com/Timothylock/go-signin-with-apple/apple"
	"google.golang.org/api/idtoken"
)

type ProviderUser struct {
	ProviderUserID string
	Email          string
	FullName       string
}

func VerifyGoogleToken(ctx context.Context, providerToken string, clientIDs []string) (*ProviderUser, error) {
	if os.Getenv("DEV_SKIP_VERIFY_TOKEN") == "true" {
		return mockProviderUser(providerToken, "Google")
	}

	var payload *idtoken.Payload
	var err error

	// Verify token against all allowed Client IDs to handle multiple platforms
	for _, clientID := range clientIDs {
		if clientID == "" {
			continue
		}
		payload, err = idtoken.Validate(ctx, providerToken, clientID)
		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, authError.AUTH_SOCIAL_002
	}

	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)

	if email == "" {
		return nil, authError.AUTH_SOCIAL_005
	}

	return &ProviderUser{
		ProviderUserID: payload.Subject,
		Email:          email,
		FullName:       name,
	}, nil
}

func VerifyAppleToken(ctx context.Context, providerToken string, clientIDs []string) (*ProviderUser, error) {
	if os.Getenv("DEV_SKIP_VERIFY_TOKEN") == "true" {
		return mockProviderUser(providerToken, "Apple")
	}

	client := apple.New()
	var claims *apple.IDTokenClaims
	var err error

	// Verify token against all allowed Client IDs to handle multiple platforms
	for _, clientID := range clientIDs {
		if clientID == "" {
			continue
		}
		claims, err = client.VerifyIDToken(ctx, providerToken, clientID)
		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, authError.AUTH_SOCIAL_002
	}

	// Apple only provides user name on the first authorization request.
	// Fallback to name extracted from email prefix on consecutive logins.
	fullName := strings.Split(claims.Email, "@")[0]

	return &ProviderUser{
		ProviderUserID: claims.Subject,
		Email:          claims.Email,
		FullName:       fullName,
	}, nil
}

func mockProviderUser(token, provider string) (*ProviderUser, error) {
	email := token
	if !strings.Contains(email, "@") {
		email = "mock_" + strings.ToLower(provider) + "@example.com"
	}
	name := strings.Split(email, "@")[0]

	return &ProviderUser{
		ProviderUserID: "mock_" + strings.ToLower(provider) + "_id_123",
		Email:          email,
		FullName:       name,
	}, nil
}
