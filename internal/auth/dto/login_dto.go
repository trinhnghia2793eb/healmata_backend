package dto

type LoginRequestDTO struct {
	Identifier string `json:"identifier" binding:"required,is_identifier"`
	Password   string `json:"password" binding:"required,min=8"`
}

type LoginResponseDTO struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
}
