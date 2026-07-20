package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"healmata_backend/pkg/email"
	"healmata_backend/internal/app/middleware"
)

func RegisterRoutes(r *gin.Engine, db *pgxpool.Pool, emailSender email.EmailSender) {
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())

	registerAuthRoutes(r, db, emailSender)

	// registerUserRoutes(r)
	// registerPaymentRoutes(r)
}
