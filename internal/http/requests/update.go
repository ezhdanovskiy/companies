package requests

import "github.com/ezhdanovskiy/companies/internal/models"

type UpdateCompany struct {
	Name            *string `json:"name" binding:"omitempty,max=15"`
	Description     *string `json:"description" binding:"omitempty,max=3000"`
	EmployeesAmount *int    `json:"employees_amount" binding:"omitempty"`
	Registered      *bool   `json:"registered" binding:"omitempty"`
	Type            *string `json:"type" binding:"omitempty,oneof=Corporations NonProfit Cooperative 'Sole Proprietorship'"`
}

func (c *UpdateCompany) ToDomain(uuid string) *models.CompanyPatch {
	return &models.CompanyPatch{
		ID:              uuid,
		Name:            c.Name,
		Description:     c.Description,
		EmployeesAmount: c.EmployeesAmount,
		Registered:      c.Registered,
		Type:            c.Type,
	}
}
