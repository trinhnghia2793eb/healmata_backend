package handler

import (
	"healmata_backend/internal/auth/service"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(s service.AuthService) *AuthHandler {

	return &AuthHandler{
		service: s,
	}
}
