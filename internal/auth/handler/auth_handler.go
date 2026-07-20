package handler

import (
	authErrors "healmata_backend/internal/auth/errors"
	"healmata_backend/internal/auth/service"
)

var validationErr = authErrors.Validation

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(s service.AuthService) *AuthHandler {

	return &AuthHandler{
		service: s,
	}
}
