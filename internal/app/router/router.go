package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"healmata_backend/internal/app/config"
	"healmata_backend/internal/app/middleware"
	dbpkg "healmata_backend/pkg/db"
	"healmata_backend/pkg/email"
	"healmata_backend/pkg/jwt"
)

// dependencies for router
type Dependencies struct {
	DB          *pgxpool.Pool
	Config      *config.Config
	Transactor  *dbpkg.SQLTxManager
	EmailSender *email.Sender
	JWTManager  *jwt.JWTManager
}

func RegisterRoutes(r *gin.Engine, deps Dependencies) {
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS(deps.Config.CorsAllowedOrigins))

	registerAuthRoutes(r, deps.DB, deps.Transactor, deps.EmailSender, deps.JWTManager)

}
