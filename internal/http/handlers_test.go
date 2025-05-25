package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	internalhttp "github.com/ezhdanovskiy/companies/internal/http"
	"github.com/ezhdanovskiy/companies/internal/http/mocks"
	"github.com/ezhdanovskiy/companies/internal/http/requests"
	"github.com/ezhdanovskiy/companies/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func setupTestServer(t *testing.T) (*internalhttp.Server, *mocks.MockService, *gin.Engine) {
	t.Helper()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockService := mocks.NewMockService(ctrl)
	logger := zap.NewNop().Sugar()
	server := internalhttp.NewServer(logger, 8080, mockService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Register routes without auth middleware for testing
	router.GET("/api/v1/companies/:uuid", server.GetCompany)
	router.POST("/api/v1/secured/companies", server.CreateCompany)
	router.PATCH("/api/v1/secured/companies/:uuid", server.UpdateCompany)
	router.DELETE("/api/v1/secured/companies/:uuid", server.DeleteCompany)

	return server, mockService, router
}

func TestCreateCompany(t *testing.T) {
	tests := []struct {
		name           string
		body           interface{}
		setupMock      func(m *mocks.MockService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success",
			body: requests.CreateCompany{
				ID:              "550e8400-e29b-41d4-a716-446655440000",
				Name:            "Test Company",
				Description:     "Test Description",
				EmployeesAmount: 100,
				Registered:      true,
				Type:            "Corporations",
			},
			setupMock: func(m *mocks.MockService) {
				m.EXPECT().CreateCompany(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid json",
			body: `{"invalid json"`,
			setupMock: func(m *mocks.MockService) {
				// No expectations, should fail before service call
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "validation error - empty name",
			body: requests.CreateCompany{
				ID:              "550e8400-e29b-41d4-a716-446655440000",
				Name:            "",
				Description:     "Test Description",
				EmployeesAmount: 100,
				Registered:      true,
				Type:            "Corporations",
			},
			setupMock: func(m *mocks.MockService) {
				// No expectations, should fail validation
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "validation error - invalid type",
			body: map[string]interface{}{
				"id":               "550e8400-e29b-41d4-a716-446655440000",
				"name":             "Test Company",
				"description":      "Test Description",
				"employees_amount": 100,
				"registered":       true,
				"type":             "InvalidType",
			},
			setupMock: func(m *mocks.MockService) {
				// No expectations, should fail validation
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service error",
			body: requests.CreateCompany{
				ID:              "550e8400-e29b-41d4-a716-446655440000",
				Name:            "Test Company",
				Description:     "Test Description",
				EmployeesAmount: 100,
				Registered:      true,
				Type:            "Corporations",
			},
			setupMock: func(m *mocks.MockService) {
				m.EXPECT().CreateCompany(gomock.Any(), gomock.Any()).
					Return(errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"database error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, mockService, router := setupTestServer(t)
			tt.setupMock(mockService)

			var body []byte
			var err error
			if str, ok := tt.body.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.body)
				require.NoError(t, err)
			}

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/secured/companies", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestUpdateCompany(t *testing.T) {
	tests := []struct {
		name           string
		uuid           string
		body           interface{}
		setupMock      func(m *mocks.MockService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success",
			uuid: "550e8400-e29b-41d4-a716-446655440000",
			body: requests.UpdateCompany{
				Name:            stringPtr("Updated Company"),
				Description:     stringPtr("Updated Description"),
				EmployeesAmount: intPtr(200),
				Registered:      boolPtr(false),
				Type:            stringPtr("Sole Proprietorship"),
			},
			setupMock: func(m *mocks.MockService) {
				m.EXPECT().UpdateCompany(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid uuid",
			uuid: "invalid-uuid",
			body: requests.UpdateCompany{},
			setupMock: func(m *mocks.MockService) {
				// No expectations, should fail UUID validation
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "company not found",
			uuid: "550e8400-e29b-41d4-a716-446655440000",
			body: requests.UpdateCompany{
				Name: stringPtr("Updated Company"),
			},
			setupMock: func(m *mocks.MockService) {
				m.EXPECT().UpdateCompany(gomock.Any(), gomock.Any()).
					Return(models.ErrCompanyNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"message":"Company not found"}`,
		},
		{
			name: "invalid json",
			uuid: "550e8400-e29b-41d4-a716-446655440000",
			body: `{"invalid json"`,
			setupMock: func(m *mocks.MockService) {
				// No expectations, should fail before service call
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "validation error - invalid type",
			uuid: "550e8400-e29b-41d4-a716-446655440000",
			body: map[string]interface{}{
				"type": "InvalidType",
			},
			setupMock: func(m *mocks.MockService) {
				// No expectations, should fail validation
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service error",
			uuid: "550e8400-e29b-41d4-a716-446655440000",
			body: requests.UpdateCompany{
				Name: stringPtr("Updated Company"),
			},
			setupMock: func(m *mocks.MockService) {
				m.EXPECT().UpdateCompany(gomock.Any(), gomock.Any()).
					Return(errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"database error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, mockService, router := setupTestServer(t)
			tt.setupMock(mockService)

			var body []byte
			var err error
			if str, ok := tt.body.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.body)
				require.NoError(t, err)
			}

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPatch, "/api/v1/secured/companies/"+tt.uuid, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestDeleteCompany(t *testing.T) {
	tests := []struct {
		name           string
		uuid           string
		setupMock      func(m *mocks.MockService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success",
			uuid: "550e8400-e29b-41d4-a716-446655440000",
			setupMock: func(m *mocks.MockService) {
				m.EXPECT().DeleteCompany(gomock.Any(), "550e8400-e29b-41d4-a716-446655440000").
					Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid uuid",
			uuid: "invalid-uuid",
			setupMock: func(m *mocks.MockService) {
				// No expectations, should fail UUID validation
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "company not found",
			uuid: "550e8400-e29b-41d4-a716-446655440000",
			setupMock: func(m *mocks.MockService) {
				m.EXPECT().DeleteCompany(gomock.Any(), "550e8400-e29b-41d4-a716-446655440000").
					Return(models.ErrCompanyNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"message":"Company not found"}`,
		},
		{
			name: "service error",
			uuid: "550e8400-e29b-41d4-a716-446655440000",
			setupMock: func(m *mocks.MockService) {
				m.EXPECT().DeleteCompany(gomock.Any(), "550e8400-e29b-41d4-a716-446655440000").
					Return(errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"database error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, mockService, router := setupTestServer(t)
			tt.setupMock(mockService)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodDelete, "/api/v1/secured/companies/"+tt.uuid, nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestGetCompany(t *testing.T) {
	tests := []struct {
		name           string
		uuid           string
		setupMock      func(m *mocks.MockService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success",
			uuid: "550e8400-e29b-41d4-a716-446655440000",
			setupMock: func(m *mocks.MockService) {
				m.EXPECT().GetCompany(gomock.Any(), "550e8400-e29b-41d4-a716-446655440000").
					Return(&models.Company{
						ID:              "550e8400-e29b-41d4-a716-446655440000",
						Name:            "Test Company",
						Description:     "Test Description",
						EmployeesAmount: 100,
						Registered:      true,
						Type:            "Corporations",
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: `{
				"ID": "550e8400-e29b-41d4-a716-446655440000",
				"Name": "Test Company",
				"Description": "Test Description",
				"EmployeesAmount": 100,
				"Registered": true,
				"Type": "Corporations"
			}`,
		},
		{
			name: "invalid uuid",
			uuid: "invalid-uuid",
			setupMock: func(m *mocks.MockService) {
				// No expectations, should fail UUID validation
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "company not found",
			uuid: "550e8400-e29b-41d4-a716-446655440000",
			setupMock: func(m *mocks.MockService) {
				m.EXPECT().GetCompany(gomock.Any(), "550e8400-e29b-41d4-a716-446655440000").
					Return(nil, nil)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"message":"Company not found"}`,
		},
		{
			name: "service error",
			uuid: "550e8400-e29b-41d4-a716-446655440000",
			setupMock: func(m *mocks.MockService) {
				m.EXPECT().GetCompany(gomock.Any(), "550e8400-e29b-41d4-a716-446655440000").
					Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"database error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, mockService, router := setupTestServer(t)
			tt.setupMock(mockService)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/companies/"+tt.uuid, nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestGetCompany_CaseInsensitiveUUID(t *testing.T) {
	_, mockService, router := setupTestServer(t)

	expectedUUID := "550e8400-e29b-41d4-a716-446655440000"
	mockService.EXPECT().GetCompany(gomock.Any(), expectedUUID).
		Return(&models.Company{
			ID:   expectedUUID,
			Name: "Test Company",
		}, nil)

	// Test with uppercase UUID
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/companies/550E8400-E29B-41D4-A716-446655440000", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}

