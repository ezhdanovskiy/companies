package http

import (
	"context"

	"github.com/ezhdanovskiy/companies/internal/models"
)

// Service describes the service methods required for the server.
type Service interface {
	CreateCompany(ctx context.Context, company *models.Company) error
}
