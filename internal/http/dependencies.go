package http

import (
	"context"

	"github.com/ezhdanovskiy/companies/internal/models"
	_ "github.com/golang/mock/mockgen/model"
)

// Service describes the service methods required for the server.
type Service interface {
	CreateCompany(ctx context.Context, company *models.Company) error
	UpdateCompany(ctx context.Context, companyPatch *models.CompanyPatch) error
	DeleteCompany(ctx context.Context, companyUUID string) error
	GetCompany(ctx context.Context, companyUUID string) (*models.Company, error)
}

//go:generate mockgen -destination=./mocks/service_mock.go -package=mocks . Service
