package models

type Company struct {
	ID              string
	Name            string
	Description     string
	EmployeesAmount int
	Registered      bool
	Type            string // Corporations | NonProfit | Cooperative | Sole Proprietorship
}
