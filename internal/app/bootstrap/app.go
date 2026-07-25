package bootstrap

import (
	dbpkg "healmata_backend/pkg/db"
	"healmata_backend/pkg/email"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Config      *Config
	DB          dbpkg.DBEngine
	Transactor  dbpkg.Transactor
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

	// ==========================================
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
	}, nil
}

func (a *App) Close() {
	if a.DB != nil {
		// type assert into pgxpool.Pool to call Close()
		if pool, ok := a.DB.(*pgxpool.Pool); ok {
			pool.Close()
		}
	}
}
