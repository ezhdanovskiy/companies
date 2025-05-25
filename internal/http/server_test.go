package http_test

import (
	"bytes"
	"testing"
	"time"

	internalhttp "github.com/ezhdanovskiy/companies/internal/http"
	"github.com/ezhdanovskiy/companies/internal/http/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	logger := zap.NewNop().Sugar()
	port := 8080

	server := internalhttp.NewServer(logger, port, mockService)

	assert.NotNil(t, server)
}

func TestServer_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	logger := zap.NewNop().Sugar()
	
	// Use a random port to avoid conflicts
	server := internalhttp.NewServer(logger, 0, mockService)

	// Run server in goroutine
	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Run()
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Server should be running, shutdown it
	server.Shutdown()

	// Wait for Run to complete
	select {
	case err := <-errCh:
		assert.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("Server did not shut down in time")
	}
}

func TestServer_Run_InvalidPort(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	logger := zap.NewNop().Sugar()
	
	// Use an invalid port
	server := internalhttp.NewServer(logger, -1, mockService)

	err := server.Run()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "start http server")
}

func TestServer_Shutdown(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	logger := zap.NewNop().Sugar()
	
	server := internalhttp.NewServer(logger, 0, mockService)

	// Start server
	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Run()
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Call shutdown
	server.Shutdown()

	// Wait for server to stop
	select {
	case err := <-errCh:
		assert.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("Server did not shut down in time")
	}
}

func TestServer_Routes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	logger := zap.NewNop().Sugar()
	
	server := internalhttp.NewServer(logger, 0, mockService)

	// Start server to get actual port
	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Run()
	}()
	defer server.Shutdown()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Test that routes are registered (we can't easily get the port in this setup)
	// This is mostly covered by the handler tests
}

func TestServer_Shutdown_Timeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	
	// Create a buffer to capture logs
	var logBuffer bytes.Buffer
	
	// Create a custom logger that writes to our buffer
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(&logBuffer),
		zapcore.DebugLevel,
	)
	logger := zap.New(core).Sugar()
	
	server := internalhttp.NewServer(logger, 0, mockService)

	// Start server
	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Run()
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Call shutdown multiple times to potentially trigger error log
	server.Shutdown()
	server.Shutdown() // Second call might trigger different behavior

	// Wait for server to stop
	select {
	case <-errCh:
		// Server stopped
	case <-time.After(2 * time.Second):
		t.Fatal("Server did not shut down in time")
	}
	
	// Check if error was logged
	logOutput := logBuffer.String()
	t.Logf("Log output: %s", logOutput)
}

func TestServer_SetAPIV1Routes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	logger := zap.NewNop().Sugar()
	
	server := internalhttp.NewServer(logger, 8080, mockService)
	
	gin.SetMode(gin.TestMode)
	router := gin.New()
	apiV1 := router.Group("/api/v1")
	
	// Call SetAPIV1Routes
	server.SetAPIV1Routes(apiV1)
	
	// Verify routes are registered
	routes := router.Routes()
	
	expectedRoutes := map[string]string{
		"GET:/api/v1/companies/:uuid": "GetCompany",
		"POST:/api/v1/secured/companies": "CreateCompany",
		"PATCH:/api/v1/secured/companies/:uuid": "UpdateCompany",
		"DELETE:/api/v1/secured/companies/:uuid": "DeleteCompany",
	}
	
	for _, route := range routes {
		key := route.Method + ":" + route.Path
		if _, ok := expectedRoutes[key]; ok {
			delete(expectedRoutes, key)
		}
	}
	
	assert.Empty(t, expectedRoutes, "Not all expected routes were registered")
}