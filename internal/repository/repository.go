// Package repository contains all the functionality for working with the DB.
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ezhdanovskiy/companies/internal/models"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"go.uber.org/zap"
)

// Repo performs database operations.
type Repo struct {
	log *zap.SugaredLogger
	db  *bun.DB
}

// MigrateUp applies migrations to DB.
func MigrateUp(logger *zap.SugaredLogger, db *sql.DB, path string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(path, "postgres", driver)
	if err != nil {
		return fmt.Errorf("migrate NewWithDatabaseInstance: %w", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate Up: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate version: %w", err)
	}
	logger.With("version", version, "dirty", dirty).Info("Migrations applied")

	return nil
}

// NewRepo creates instance of repository using existing DB.
func NewRepo(logger *zap.SugaredLogger, db *bun.DB) (*Repo, error) {
	return &Repo{
		log: logger,
		db:  db,
	}, nil
}

// CreateCompany insert a company.
func (r *Repo) CreateCompany(ctx context.Context, c *models.Company) error {
	r.log.With("id", c.ID, "name", c.Name, "descr", c.Description, "amount", c.EmployeesAmount,
		"registered", c.Registered, "type", c.Type).Debug("Repo.CreateCompany")

	_, err := r.db.NewInsert().Model(newCompany(c)).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

// UpdateCompany insert a company.
func (r *Repo) UpdateCompany(ctx context.Context, c *models.CompanyPatch) (affected int64, err error) {
	r.log.With("id", c.ID, "name", c.Name, "descr", c.Description, "amount", c.EmployeesAmount,
		"registered", c.Registered, "type", c.Type).Debug("Repo.CreateCompany")

	t := time.Now()
	company := &Company{
		ID:        c.ID,
		UpdatedAt: &t,
	}
	fields := []string{"updated_at"}

	if c.Name != nil && *c.Name != "" {
		company.Name = *c.Name
		fields = append(fields, "name")
	}

	if c.Description != nil {
		company.Description = *c.Description
		fields = append(fields, "description")
	}

	if c.EmployeesAmount != nil {
		company.EmployeesAmount = *c.EmployeesAmount
		fields = append(fields, "employees_amount")
	}

	if c.Registered != nil {
		company.Registered = *c.Registered
		fields = append(fields, "registered")
	}

	if c.Type != nil {
		company.Type = *c.Type
		fields = append(fields, "type")
	}

	res, err := r.db.NewUpdate().Model(company).Column(fields...).WherePK().Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("update company: %w", err)
	}

	affected, err = res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("delete company rows affected: %w", err)
	}

	return affected, nil
}

// DeleteCompany deletes a company.
func (r *Repo) DeleteCompany(ctx context.Context, uuid string) (affected int64, err error) {
	r.log.With("uuid", uuid).Debug("Repo.DeleteCompany")

	res, err := r.db.NewDelete().Model(&Company{ID: uuid}).WherePK().Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("delete company: %w", err)
	}

	affected, err = res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("delete company rows affected: %w", err)
	}

	return affected, nil
}

// GetCompany selects company by uuid.
func (r *Repo) GetCompany(ctx context.Context, uuid string) (*models.Company, error) {
	r.log.With("uuid", uuid).Debug("Repo.GetCompany")

	company := new(Company)
	err := r.db.NewSelect().Model(company).Where("id = ?", uuid).Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("select company: %w", err)
	}

	return company.toDomain(), nil
}
