package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetCompany(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

func (s *Server) CreateCompany(c *gin.Context) {
	s.log.Debug("CreateCompany")

	var company Company
	if err := c.ShouldBindJSON(&company); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	err := s.svc.CreateCompany(c.Request.Context(), company.toDomain())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, nil)
}
