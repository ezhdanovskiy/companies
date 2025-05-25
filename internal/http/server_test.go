package http_test

import (
	"bytes"
	"sync"
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
	var wg sync.WaitGroup
	wg.Add(1)
	errCh := make(chan error, 1)
	go func() {
		defer wg.Done()
		errCh <- server.Run()
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Server should be running, shutdown it
	server.Shutdown()

	// Wait for Run to complete
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Check error
		select {
		case err := <-errCh:
			assert.NoError(t, err)
		default:
			// No error
		}
	case <-time.After(10 * time.Second):
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
	var wg sync.WaitGroup
	wg.Add(1)
	errCh := make(chan error, 1)
	go func() {
		defer wg.Done()
		errCh <- server.Run()
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Call shutdown
	server.Shutdown()

	// Wait for server to stop
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Check error
		select {
		case err := <-errCh:
			assert.NoError(t, err)
		default:
			// No error
		}
	case <-time.After(10 * time.Second):
		t.Fatal("Server did not shut down in time")
	}
}

func TestServer_Routes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	logger := zap.NewNop().Sugar()
	
	server := internalhttp.NewServer(logger, 8080, mockService)

	// Create a test router and register routes
	gin.SetMode(gin.TestMode)
	router := gin.New()
	apiV1 := router.Group("/api/v1")
	server.SetAPIV1Routes(apiV1)

	// Verify routes are registered
	routes := router.Routes()
	expectedPaths := map[string]bool{
		"/api/v1/companies/:uuid": false,
		"/api/v1/secured/companies": false,
		"/api/v1/secured/companies/:uuid": false,
	}

	for _, route := range routes {
		if _, ok := expectedPaths[route.Path]; ok {
			expectedPaths[route.Path] = true
		}
	}

	for path, found := range expectedPaths {
		assert.True(t, found, "Route %s not found", path)
	}
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
	var wg sync.WaitGroup
	wg.Add(1)
	errCh := make(chan error, 1)
	go func() {
		defer wg.Done()
		errCh <- server.Run()
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Call shutdown
	server.Shutdown()

	// Wait for server to stop
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Server stopped
	case <-time.After(10 * time.Second):
		t.Fatal("Server did not shut down in time")
	}
	
	// Check if error was logged
	logOutput := logBuffer.String()
	t.Logf("Log output: %s", logOutput)
}