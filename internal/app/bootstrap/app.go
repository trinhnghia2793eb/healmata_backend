package bootstrap

import (
	"healmata_backend/internal/app/config"
	dbpkg "healmata_backend/pkg/db"
	"healmata_backend/pkg/email"
	"healmata_backend/pkg/jwt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Config      *config.Config
	DB          *pgxpool.Pool
	Transactor  *dbpkg.SQLTxManager
	EmailSender *email.Sender
	JWTManager  *jwt.JWTManager
}

func NewApp() (*App, error) {

	// load config from config folder --> return cfg
	cfg, err := config.LoadConfig(".")
	if err != nil {
		return nil, err
	}

	// use cfg to create db connect
	db, err := NewDatabase(cfg)
	if err != nil {
		return nil, err
	}

	// ==========================================
	// init jwtManager
	jwtManager := jwt.NewJWTManager(cfg.JWTSecret, cfg.JWTAccessExpiry, cfg.JWTRefreshExpiry)

	// init transactor using txManager
	txManager := dbpkg.NewSQLTxManager(db)

	// init email sender
	emailSender := email.NewEmailSender(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUser,
		cfg.SMTPPassword,
		cfg.MailFromAddress,
		cfg.MailFromName,
	)

	// ==========================================

	// return
	return &App{
		Config:      cfg,
		DB:          db,
		Transactor:  txManager,
		EmailSender: emailSender,
		JWTManager:  jwtManager,
	}, nil
}

func (a *App) Close() {
	if a.DB != nil {
		a.DB.Close()
	}
}
