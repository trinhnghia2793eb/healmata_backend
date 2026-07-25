package router

import (
	"github.com/gin-gonic/gin"

	"healmata_backend/internal/app/middleware"
	dbpkg "healmata_backend/pkg/db"
	"healmata_backend/pkg/email"
)

// dependencies for router
type Dependencies struct {
	DB          dbpkg.DBEngine
	Transactor  dbpkg.Transactor
	EmailSender email.EmailSender
}

func RegisterRoutes(r *gin.Engine, deps Dependencies) {
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())

	registerAuthRoutes(r, deps.DB, deps.Transactor, deps.EmailSender)

}
