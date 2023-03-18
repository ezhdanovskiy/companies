package http

import "github.com/ezhdanovskiy/companies/internal/models"

type Company struct {
	ID              string `json:"id" binding:"required,uuid"`
	Name            string `json:"name" binding:"required,max=3000"`
	Description     string `json:"description" binding:"omitempty,max=3000"`
	EmployeesAmount int    `json:"employees_amount" binding:"required"`
	Registered      bool   `json:"registered" binding:"required"`
	Type            string `json:"type" binding:"required,oneof=Corporations NonProfit Cooperative 'Sole Proprietorship'"` // Corporations | NonProfit | Cooperative | Sole Proprietorship
}

func (c *Company) toDomain() *models.Company {
	return &models.Company{
		ID:              c.ID,
		Name:            c.Name,
		Description:     c.Description,
		EmployeesAmount: c.EmployeesAmount,
		Registered:      c.Registered,
		Type:            c.Type,
	}
}