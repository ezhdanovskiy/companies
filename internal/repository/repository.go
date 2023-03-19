// Package repository contains all the functionality for working with the DB.
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ezhdanovskiy/companies/internal/models"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// Repo performs database operations.
type Repo struct {
	log *zap.SugaredLogger
	db  *sqlx.DB
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
func NewRepo(logger *zap.SugaredLogger, db *sqlx.DB) (*Repo, error) {
	return &Repo{
		log: logger,
		db:  db,
	}, nil
}

// CreateCompany insert a company.
func (r *Repo) CreateCompany(ctx context.Context, c *models.Company) error {
	r.log.With("id", c.ID, "name", c.Name, "descr", c.Description, "amount", c.EmployeesAmount,
		"registered", c.Registered, "type", c.Type).Debug("Repo.CreateCompany")

	const query = `
INSERT INTO companies (id, name, description, employees_amount, registered, type)
VALUES (:id, :name, :description, :employees_amount, :registered, :type)
`
	_, err := r.db.NamedExecContext(ctx, query, newCompany(c))
	if err != nil {
		return fmt.Errorf("insert company: %w", err)
	}

	return nil
}

// UpdateCompany insert a company.
func (r *Repo) UpdateCompany(ctx context.Context, c *models.CompanyPatch) error {
	r.log.With("id", c.ID, "name", c.Name, "descr", c.Description, "amount", c.EmployeesAmount,
		"registered", c.Registered, "type", c.Type).Debug("Repo.CreateCompany")

	args := []any{c.ID}
	var setStr string
	if c.Name != nil {
		args = append(args, *c.Name)
		setStr += fmt.Sprintf("%s = $%v, ", "name", len(args)+1)
	}

	setStr = strings.TrimSuffix(setStr, ", ")

	const query = `
UPDATE companies 
SET %s
WHERE id = $1
`

	s := fmt.Sprintf(query, setStr)
	_, err := r.db.ExecContext(ctx, s, args...)
	if err != nil {
		r.log.Error(err)
		return fmt.Errorf("update company: %w", err)
	}

	return nil
}

// DeleteCompany deletes a company.
func (r *Repo) DeleteCompany(ctx context.Context, uuid string) error {
	r.log.With("uuid", uuid).Debug("Repo.DeleteCompany")

	const query = `
DELETE FROM companies
WHERE id = $1
`

	_, err := r.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return fmt.Errorf("delete company: %w", err)
	}

	return nil
}

// GetCompany selects company by uuid.
func (r *Repo) GetCompany(ctx context.Context, uuid string) (*models.Company, error) {
	r.log.With("uuid", uuid).Debug("Repo.GetCompany")
	const query = `
SELECT * 
FROM companies 
WHERE id = $1 
`

	var company Company
	err := r.db.GetContext(ctx, &company, query, uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("select company: %w", err)
	}

	return company.toDomain(), nil
}
