// Package http contains the HTTP server and associated endpoint handlers.
package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ezhdanovskiy/companies/internal/middlewares"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	log        *zap.SugaredLogger
	httpPort   int
	httpServer *http.Server
	svc        Service
}

func NewServer(logger *zap.SugaredLogger, httpPort int, svc Service) *Server {
	return &Server{
		log:      logger,
		httpPort: httpPort,
		svc:      svc,
	}
}

func (s *Server) Run() error {
	router := gin.Default()
	apiV1 := router.Group("/api/v1")
	s.SetAPIV1Routes(apiV1)

	s.httpServer = &http.Server{
		Addr:              fmt.Sprintf(":%d", s.httpPort),
		Handler:           router,
		ReadHeaderTimeout: 3 * time.Second,
	}

	err := s.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("start http server: %w", err)
	}

	return nil
}

func (s *Server) SetAPIV1Routes(rg *gin.RouterGroup) {
	rg.GET("/companies/:uuid", s.GetCompany)
	secured := rg.Group("/secured").Use(middlewares.Auth())
	secured.POST("/companies", s.CreateCompany)
	secured.PATCH("/companies/:uuid", s.UpdateCompany)
	secured.DELETE("/companies/:uuid", s.DeleteCompany)
}

func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.log.Errorf("http server shutdown: %s", err)
	}
}
