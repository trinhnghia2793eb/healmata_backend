package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"healmata_backend/internal/app/bootstrap"
	"healmata_backend/internal/app/router"
)

func main() {
	//
	app, err := bootstrap.NewApp()
	if err != nil {
		log.Fatal("Lỗi xảy ra khi khởi chạy app: ", err)
	}
	defer app.Close()

	// init gin
	gin.SetMode(app.Config.GinMode)

	r := gin.Default()

	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatal("Cannot config trusted proxies:", err)
	}

	router.RegisterRoutes(r, app.DB, app.EmailSender)

	//
	log.Printf(
		"%s is running on :%s",
		app.Config.AppName,
		app.Config.AppPort,
	)

	if err := r.Run(":" + app.Config.AppPort); err != nil {
		log.Fatal(err)
	}
}
