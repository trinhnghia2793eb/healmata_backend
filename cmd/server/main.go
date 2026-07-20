package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"healmata_backend/internal/app/bootstrap"
	"healmata_backend/internal/app/logger"
	"healmata_backend/internal/app/router"
)

func main() {
	// bootstrap files
	app, err := bootstrap.NewApp()
	if err != nil {
		log.Fatal("Lỗi xảy ra khi khởi chạy app: ", err)
	}
	defer app.Close()

	logger.InitLogger(app.Config.AppEnv)
	logger.Log.Info("Initialization...", "env", app.Config.AppEnv)

	// init gin
	gin.SetMode(app.Config.GinMode)

	r := gin.New()

	if err := r.SetTrustedProxies(nil); err != nil {
		logger.Log.Error("Cannot config trusted proxies:", "error", err)
		log.Fatal(err)
	}

	router.RegisterRoutes(r, app.DB, app.EmailSender)

	// run server
	if err := r.Run(":" + app.Config.AppPort); err != nil {
		logger.Log.Error("Server interrupted", "error", err)
		log.Fatal(err)
	}
}
