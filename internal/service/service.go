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
	s.log.Debug("Service.CreateCompany")
	return s.repo.CreateCompany(ctx, company)
}
