// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ezhdanovskiy/companies/internal/http (interfaces: Service)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	models "github.com/ezhdanovskiy/companies/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// CreateCompany mocks base method.
func (m *MockService) CreateCompany(arg0 context.Context, arg1 *models.Company) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCompany", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateCompany indicates an expected call of CreateCompany.
func (mr *MockServiceMockRecorder) CreateCompany(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCompany", reflect.TypeOf((*MockService)(nil).CreateCompany), arg0, arg1)
}

// DeleteCompany mocks base method.
func (m *MockService) DeleteCompany(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCompany", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCompany indicates an expected call of DeleteCompany.
func (mr *MockServiceMockRecorder) DeleteCompany(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCompany", reflect.TypeOf((*MockService)(nil).DeleteCompany), arg0, arg1)
}

// GetCompany mocks base method.
func (m *MockService) GetCompany(arg0 context.Context, arg1 string) (*models.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCompany", arg0, arg1)
	ret0, _ := ret[0].(*models.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCompany indicates an expected call of GetCompany.
func (mr *MockServiceMockRecorder) GetCompany(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCompany", reflect.TypeOf((*MockService)(nil).GetCompany), arg0, arg1)
}

// UpdateCompany mocks base method.
func (m *MockService) UpdateCompany(arg0 context.Context, arg1 *models.CompanyPatch) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCompany", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCompany indicates an expected call of UpdateCompany.
func (mr *MockServiceMockRecorder) UpdateCompany(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCompany", reflect.TypeOf((*MockService)(nil).UpdateCompany), arg0, arg1)
}
