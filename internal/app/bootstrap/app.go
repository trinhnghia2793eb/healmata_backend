package bootstrap

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"healmata_backend/pkg/email"
)

type App struct {
	Config      *Config
	DB          *pgxpool.Pool
	EmailSender email.EmailSender
}

func NewApp() (*App, error) {

	// load .env --> return cfg
	cfg, err := LoadEnv()
	if err != nil {
		return nil, err
	}

	// use cfg to create db connect
	db, err := NewDatabase(cfg)
	if err != nil {
		return nil, err
	}

	// init email sender
	emailSender := email.NewEmailSender(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUser,
		cfg.SMTPPassword,
		cfg.MailFromAddress,
		cfg.MailFromName,
	)

	// return
	return &App{
		Config:      cfg,
		DB:          db,
		EmailSender: emailSender,
	}, nil
}

func (a *App) Close() {
	if a.DB != nil {
		a.DB.Close()
	}
}
