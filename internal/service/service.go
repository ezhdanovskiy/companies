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

func (s *Service) UpdateCompany(ctx context.Context, companyPatch *models.CompanyPatch) error {
	s.log.With("id", companyPatch.ID).Debug("Service.UpdateCompany")
	affected, err := s.repo.UpdateCompany(ctx, companyPatch)
	if err != nil {
		return err
	}
	if affected == 0 {
		return models.ErrCompanyNotFound
	}
	return nil
}

func (s *Service) DeleteCompany(ctx context.Context, uuid string) error {
	s.log.With("uuid", uuid).Debug("Service.DeleteCompany")
	affected, err := s.repo.DeleteCompany(ctx, uuid)
	if err != nil {
		return err
	}
	if affected == 0 {
		return models.ErrCompanyNotFound
	}
	return nil
}

func (s *Service) GetCompany(ctx context.Context, uuid string) (*models.Company, error) {
	s.log.With("uuid", uuid).Debug("Service.GetCompany")
	return s.repo.GetCompany(ctx, uuid)
}
