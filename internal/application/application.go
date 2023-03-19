// Package application runs the required components depending on the parameters.
package application

import (
	"database/sql"
	"fmt"

	"github.com/ezhdanovskiy/companies/internal/auth"
	"github.com/ezhdanovskiy/companies/internal/kafka"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
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

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		a.cfg.DB.User, a.cfg.DB.Password, a.cfg.DB.Host, a.cfg.DB.Port, a.cfg.DB.DBName)
	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(pgdb, pgdialect.New())

	// Print all queries to stdout.
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	if err := db.Ping(); err != nil {
		return fmt.Errorf("db pinf: %w", err)
	}

	if err := repository.MigrateUp(a.log, db.DB, "file://"+a.cfg.DB.MigrationsPath); err != nil {
		return fmt.Errorf("migrate up: %w", err)
	}

	repo, err := repository.NewRepo(a.log, db)
	if err != nil {
		return fmt.Errorf("new repo: %w", err)
	}

	producer := kafka.NewAsyncProducer(&kafka.ProducerConfig{
		Brokers:      []string{a.cfg.Kafka.Addr},
		Topic:        a.cfg.Kafka.Topic,
		BatchSize:    a.cfg.Kafka.BatchSize,
		BatchTimeout: a.cfg.Kafka.BatchTimeout,
	})

	a.svc = service.NewService(a.log, repo, producer)

	auth.SetJWTKey(a.cfg.JWTKey)
	a.httpServer = http.NewServer(a.log, a.cfg.HTTPPort, a.svc)

	a.log.Infof("Run HTTP server on port %v", a.cfg.HTTPPort)

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
