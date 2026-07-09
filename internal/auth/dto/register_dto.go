package dto

type RegisterRequestDTO struct {
	FullName        string `json:"fullName" binding:"required,min=2,max=100"`
	Identifier      string `json:"identifier" binding:"required"`
	Password        string `json:"password" binding:"required,min=8,max=128"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=Password"`
}

type RegisterResponseDTO struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
}
