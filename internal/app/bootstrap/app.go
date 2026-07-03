package bootstrap

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Config *Config
	DB     *pgxpool.Pool
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

	// return 
	return &App{
		Config: cfg,
		DB:     db,
	}, nil
}

func (a *App) Close() {
	if a.DB != nil {
		a.DB.Close()
	}
}