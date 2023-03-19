package service

import (
	"context"
	"encoding/json"

	"github.com/ezhdanovskiy/companies/internal/models"
	"go.uber.org/zap"
)

type Service struct {
	log      *zap.SugaredLogger
	repo     Repository
	producer Producer
}

func NewService(log *zap.SugaredLogger, repo Repository, producer Producer) *Service {
	return &Service{
		log:      log,
		repo:     repo,
		producer: producer,
	}
}

type Event struct {
	Message string
	Body    interface{}
}

func (s *Service) publish(ctx context.Context, ev *Event) error {
	message, err := json.Marshal(ev)
	if err != nil {
		return err
	}

	if err := s.producer.Publish(ctx, message); err != nil {
		return err
	}
	s.log.With("message", string(message)).Debug("Event published")

	return nil
}

func (s *Service) CreateCompany(ctx context.Context, company *models.Company) error {
	s.log.With("id", company.ID).Debug("Service.CreateCompany")
	err := s.repo.CreateCompany(ctx, company)
	if err != nil {
		return err
	}

	err = s.publish(ctx, &Event{
		Message: "Company created",
		Body:    company,
	})
	if err != nil {
		s.log.With("error", err).Warn("Failed to publish message")
	}

	return nil
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

	err = s.publish(ctx, &Event{
		Message: "Company updated",
		Body:    companyPatch,
	})
	if err != nil {
		s.log.With("error", err).Warn("Failed to publish message")
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

	err = s.publish(ctx, &Event{
		Message: "Company deleted",
		Body:    uuid,
	})
	if err != nil {
		s.log.With("error", err).Warn("Failed to publish message")
	}

	return nil
}

func (s *Service) GetCompany(ctx context.Context, uuid string) (*models.Company, error) {
	s.log.With("uuid", uuid).Debug("Service.GetCompany")
	return s.repo.GetCompany(ctx, uuid)
}
