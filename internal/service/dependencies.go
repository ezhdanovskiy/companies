package service

import (
	"context"

	"github.com/ezhdanovskiy/companies/internal/models"
)

// Repository describes the repository methods required for the service.
type Repository interface {
	CreateCompany(ctx context.Context, company *models.Company) error
	UpdateCompany(ctx context.Context, companyPatch *models.CompanyPatch) (affected int64, err error)
	DeleteCompany(ctx context.Context, companyUUID string) (affected int64, err error)
	GetCompany(ctx context.Context, companyUUID string) (*models.Company, error)
}

//go:generate mockgen -destination=./mocks/repository_mock.go -package=mocks . Repository
