package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	
	"healmata_backend/pkg/email"
)

func RegisterRoutes(r *gin.Engine, db *pgxpool.Pool, emailSender email.EmailSender) {
	registerAuthRoutes(r, db, emailSender)

	// registerUserRoutes(r)
	// registerPaymentRoutes(r)
}
