package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) CreateCompany(c *gin.Context) {
	s.log.Debug("Server.CreateCompany")

	var req CreateCompany
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	err := s.svc.CreateCompany(c.Request.Context(), req.toDomain())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (s *Server) GetCompany(c *gin.Context) {
	uid := strings.ToLower(c.Param("uuid"))

	if _, err := uuid.Parse(uid); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	s.log.With("uuid", uid).Debug("Server.GetCompany")

	company, err := s.svc.GetCompany(c.Request.Context(), uid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	if company == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Company not found",
		})
		return
	}

	c.JSON(http.StatusOK, company)
}
