package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ezhdanovskiy/companies/internal/models"
	"github.com/ezhdanovskiy/companies/internal/service/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const logsEnabled = false

var ctx context.Context

func TestNewService_CreateCompany(t *testing.T) {
	ts := newTestService(t)
	defer ts.Finish()

	company := &models.Company{}

	ts.mockRepo.EXPECT().CreateCompany(ctx, company).
		Return(nil)

	ts.mockProducer.EXPECT().Publish(ctx, gomock.Any()).
		Return(nil)

	err := ts.svc.CreateCompany(ctx, company)
	require.NoError(t, err)
}

func TestNewService_CreateCompany_Error(t *testing.T) {
	ts := newTestService(t)
	defer ts.Finish()

	company := &models.Company{}
	expectedErr := errors.New("CreateCompanyError")

	ts.mockRepo.EXPECT().CreateCompany(ctx, company).
		Return(expectedErr)

	err := ts.svc.CreateCompany(ctx, &models.Company{})
	require.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestNewService_UpdateCompany(t *testing.T) {
	ts := newTestService(t)
	defer ts.Finish()

	company := &models.CompanyPatch{}
	affected := int64(1)

	ts.mockRepo.EXPECT().UpdateCompany(ctx, company).
		Return(affected, nil)

	ts.mockProducer.EXPECT().Publish(ctx, gomock.Any()).
		Return(nil)

	err := ts.svc.UpdateCompany(ctx, company)
	require.NoError(t, err)
}

func TestNewService_UpdateCompany_Error(t *testing.T) {
	ts := newTestService(t)
	defer ts.Finish()

	company := &models.CompanyPatch{}
	affected := int64(0)
	expectedErr := errors.New("CreateCompanyError")

	ts.mockRepo.EXPECT().UpdateCompany(ctx, company).
		Return(affected, expectedErr)

	err := ts.svc.UpdateCompany(ctx, company)
	require.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestNewService_UpdateCompany_NotFound(t *testing.T) {
	ts := newTestService(t)
	defer ts.Finish()

	company := &models.CompanyPatch{}
	affected := int64(0)

	ts.mockRepo.EXPECT().UpdateCompany(ctx, company).
		Return(affected, nil)

	err := ts.svc.UpdateCompany(ctx, company)
	require.Error(t, err)
	assert.Equal(t, models.ErrCompanyNotFound, err)
}

// todo: write more tests

// TestService ---------------------------------------------------------------------------------------------------------
type TestService struct {
	t            *testing.T
	log          *zap.SugaredLogger
	mockCtrl     *gomock.Controller
	mockRepo     *mocks.MockRepository
	mockProducer *mocks.MockProducer
	svc          *Service
}

func newTestService(t *testing.T) TestService {
	t.Parallel()
	mockCtrl := gomock.NewController(t)
	ts := TestService{
		t:            t,
		mockCtrl:     mockCtrl,
		mockRepo:     mocks.NewMockRepository(mockCtrl),
		mockProducer: mocks.NewMockProducer(mockCtrl),
	}

	if logsEnabled {
		logger, _ := zap.NewDevelopment()
		ts.log = logger.Sugar()
	} else {
		ts.log = zap.NewNop().Sugar()
	}

	ts.svc = NewService(ts.log, ts.mockRepo, ts.mockProducer)

	return ts
}

func (ts *TestService) Finish() {
	ts.mockCtrl.Finish()
}
