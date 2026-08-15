package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// declare dependencies object
	deps := router.Dependencies{
		Config:      app.Config,
		DB:          app.DB,
		Transactor:  app.Transactor,
		EmailSender: app.EmailSender,
		JWTManager:  app.JWTManager,
	}

	router.RegisterRoutes(r, deps)

	// =======================================================================
	// RUN SERVER: GRACEFUL SHUTDOWN
	// declare http.Server instead of r.Run()
	srv := &http.Server{
		Addr:    ":" + app.Config.AppPort,
		Handler: r,
	}
	// run server in a seperate Groutine
	go func() {
		logger.Log.Info("Server is running on port: " + app.Config.AppPort)
		// ErrServerClosed is returned when called srv.Shutdown(), not error
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Error("Server interrupted", "error", err)
			os.Exit(1)
		}
	}()

	// create a channel to listen for signals from the operating system (OS Signals)
	quit := make(chan os.Signal, 1)

	// listen for SIGINT (Ctrl+C) and SIGTERM (Docker/K8s kill command) signals
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// The main() thread will be blocked here until the channel receives a signal
	<-quit
	logger.Log.Info("Shutting down server...")

	// create a Context with a countdown
	// --> the server has a time to finish processing pending requests
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// shutdown: stop accepting new requests and perform cleanup
	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Error("Error when shutting down", "error", err)
	}

	// after srv.Shutdown() completes / timeout expires
	// --> `main()` function terminates --> defer app.Close() triggered
	logger.Log.Info("Server shutted down.")
	// =====================================================================
}
