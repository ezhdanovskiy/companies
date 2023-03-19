package repository

import (
	"time"

	"github.com/ezhdanovskiy/companies/internal/models"
)

type Company struct {
	ID              string     `db:"id"`
	Name            string     `db:"name"`
	Description     string     `db:"description"`
	EmployeesAmount int        `db:"employees_amount"`
	Registered      bool       `db:"registered"`
	Type            string     `db:"type"` // Corporations | NonProfit | Cooperative | Sole Proprietorship
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       *time.Time `db:"updated_at"`
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

func newCompany(m *models.Company) *Company {
	return &Company{
		ID:              m.ID,
		Name:            m.Name,
		Description:     m.Description,
		EmployeesAmount: m.EmployeesAmount,
		Registered:      m.Registered,
		Type:            m.Type,
	}
}
