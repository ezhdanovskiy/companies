package service

import (
	"context"

	"github.com/ezhdanovskiy/companies/internal/models"
	"go.uber.org/zap"
)

type Service struct {
	log  *zap.SugaredLogger
	repo Repository
}

func NewService(log *zap.SugaredLogger, repo Repository) *Service {
	return &Service{
		log:  log,
		repo: repo,
	}
}

func (s *Service) CreateCompany(ctx context.Context, company *models.Company) error {
	s.log.With("id", company.ID).Debug("Service.CreateCompany")
	return s.repo.CreateCompany(ctx, company)
}

func (s *Service) GetCompany(ctx context.Context, uuid string) (*models.Company, error) {
	s.log.With("uuid", uuid).Debug("Service.GetCompany")
	return s.repo.GetCompany(ctx, uuid)
}
