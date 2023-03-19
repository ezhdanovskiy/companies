package repository

import (
	"time"

	"github.com/ezhdanovskiy/companies/internal/models"
	"github.com/uptrace/bun"
)

type Company struct {
	bun.BaseModel `bun:"table:companies,alias:c"`

	ID              string     `bun:"id,pk"`
	Name            string     `bun:"name"`
	Description     string     `bun:"description"`
	EmployeesAmount int        `bun:"employees_amount"`
	Registered      bool       `bun:"registered"`
	Type            string     `bun:"type"` // Corporations | NonProfit | Cooperative | Sole Proprietorship
	CreatedAt       time.Time  `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt       *time.Time `bun:"updated_at"`
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
