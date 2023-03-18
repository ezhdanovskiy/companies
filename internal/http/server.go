// Package http contains the HTTP server and associated endpoint handlers.
package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

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
	apiV1.GET("/companies", s.GetCompany)
	//apiV1.Use(middlewares.Auth()).
	apiV1.POST("/companies", s.CreateCompany)
	//{
	//	api.POST("/token", controllers.GenerateToken)
	//	api.POST("/user/register", controllers.RegisterUser)
	//	secured := api.Group("/secured").Use(middlewares.Auth())
	//	{
	//		secured.GET("/ping", controllers.Ping)
	//	}
	//}

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.httpPort),
		Handler: router,
	}

	err := s.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("start http server: %w", err)
	}

	return nil
}

func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.log.Errorf("http server shutdown: %s", err)
	}
}
