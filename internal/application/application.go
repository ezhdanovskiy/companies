// Package application runs the required components depending on the parameters.
package application

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/ezhdanovskiy/companies/internal/config"
	"github.com/ezhdanovskiy/companies/internal/http"
	"github.com/ezhdanovskiy/companies/internal/repository"
	"github.com/ezhdanovskiy/companies/internal/service"
)

// Application contains all components of application.
type Application struct {
	log *zap.SugaredLogger
	cfg *config.Config
	svc *service.Service

	httpServer *http.Server
}

// NewApplication creates and connects instances of all components required to run Application.
func NewApplication() (*Application, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("new config: %w", err)
	}

	log, err := newLogger(cfg.LogLevel, cfg.LogEncoding)
	if err != nil {
		return nil, fmt.Errorf("new logger: %w", err)
	}
	log.Debugf("cfg: %+v", cfg)

	return &Application{
		log: log,
		cfg: cfg,
	}, nil
}

// Run runs configured components.
func (a *Application) Run() error {
	a.log.Info("Run application")

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		a.cfg.DB.Host, a.cfg.DB.Port, a.cfg.DB.User, a.cfg.DB.Password, a.cfg.DB.DBName)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}

	repo, err := repository.NewRepo(a.log, db)
	if err != nil {
		return fmt.Errorf("new repo: %w", err)
	}

	a.svc = service.NewService(a.log, repo)

	a.httpServer = http.NewServer(a.log, a.cfg.HttpPort, a.svc)

	a.log.Infof("Run HTTP server on port %v", a.cfg.HttpPort)

	if err := a.httpServer.Run(); err != nil {
		return fmt.Errorf("HTTP server run: %w", err)
	}
	a.log.Info("HTTP server stopped")

	a.log.Info("Application stopped")
	return nil
}

// Stop terminates configured components.
func (a *Application) Stop() {
	if a.httpServer != nil {
		a.log.Info("Stopping HTTP server")
		a.httpServer.Shutdown()
	}
}
