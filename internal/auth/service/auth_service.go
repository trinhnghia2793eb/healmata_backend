package service

import (
	"context"
	"healmata_backend/internal/auth/repository"
)

type AuthService interface {
	Register(ctx context.Context, req )
}

type authService struct {
	repo repository.AuthRepository
}

func NewAuthService() AuthService {

	return &authService{}
}
