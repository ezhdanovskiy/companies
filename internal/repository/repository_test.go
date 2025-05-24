package repository

import (
	"testing"

	"github.com/ezhdanovskiy/companies/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestPrepareCompanyPatch_Empty(t *testing.T) {
	c := &models.CompanyPatch{ID: "uuid1"}
	company, fields := prepareCompanyPatch(c)
	assert.Len(t, fields, 1)
	assert.Equal(t, c.ID, company.ID)
	assert.NotEmpty(t, company.UpdatedAt)
}

func TestPrepareCompanyPatch_Name(t *testing.T) {
	c := &models.CompanyPatch{
		ID:   "uuid1",
		Name: newStringPointer("Name1"),
	}
	company, fields := prepareCompanyPatch(c)
	assert.Equal(t, fields, []string{"updated_at", "name"})
	assert.Equal(t, c.ID, company.ID)
	assert.Equal(t, *c.Name, company.Name)
	assert.NotEmpty(t, company.UpdatedAt)
}

func TestPrepareCompanyPatch_EmptyName(t *testing.T) {
	c := &models.CompanyPatch{
		ID:   "uuid1",
		Name: newStringPointer(""),
	}
	company, fields := prepareCompanyPatch(c)
	assert.Equal(t, fields, []string{"updated_at"})
	assert.Equal(t, c.ID, company.ID)
	assert.NotEmpty(t, company.UpdatedAt)
}

func TestPrepareCompanyPatch_Description(t *testing.T) {
	description := "Test description"
	c := &models.CompanyPatch{
		ID:          "uuid1",
		Description: &description,
	}
	company, fields := prepareCompanyPatch(c)
	assert.Equal(t, fields, []string{"updated_at", "description"})
	assert.Equal(t, c.ID, company.ID)
	assert.Equal(t, *c.Description, company.Description)
	assert.NotEmpty(t, company.UpdatedAt)
}

func TestPrepareCompanyPatch_EmployeesAmount(t *testing.T) {
	amount := 100
	c := &models.CompanyPatch{
		ID:              "uuid1",
		EmployeesAmount: &amount,
	}
	company, fields := prepareCompanyPatch(c)
	assert.Equal(t, fields, []string{"updated_at", "employees_amount"})
	assert.Equal(t, c.ID, company.ID)
	assert.Equal(t, *c.EmployeesAmount, company.EmployeesAmount)
	assert.NotEmpty(t, company.UpdatedAt)
}

func TestPrepareCompanyPatch_Registered(t *testing.T) {
	registered := true
	c := &models.CompanyPatch{
		ID:         "uuid1",
		Registered: &registered,
	}
	company, fields := prepareCompanyPatch(c)
	assert.Equal(t, fields, []string{"updated_at", "registered"})
	assert.Equal(t, c.ID, company.ID)
	assert.Equal(t, *c.Registered, company.Registered)
	assert.NotEmpty(t, company.UpdatedAt)
}

func TestPrepareCompanyPatch_Type(t *testing.T) {
	companyType := string(models.Corporations)
	c := &models.CompanyPatch{
		ID:   "uuid1",
		Type: &companyType,
	}
	company, fields := prepareCompanyPatch(c)
	assert.Equal(t, fields, []string{"updated_at", "type"})
	assert.Equal(t, c.ID, company.ID)
	assert.Equal(t, *c.Type, company.Type)
	assert.NotEmpty(t, company.UpdatedAt)
}

func TestPrepareCompanyPatch_AllFields(t *testing.T) {
	name := "Test Company"
	description := "Test description"
	amount := 500
	registered := true
	companyType := string(models.NonProfit)

	c := &models.CompanyPatch{
		ID:              "uuid1",
		Name:            &name,
		Description:     &description,
		EmployeesAmount: &amount,
		Registered:      &registered,
		Type:            &companyType,
	}
	company, fields := prepareCompanyPatch(c)
	assert.Equal(t, fields, []string{"updated_at", "name", "description", "employees_amount", "registered", "type"})
	assert.Equal(t, c.ID, company.ID)
	assert.Equal(t, *c.Name, company.Name)
	assert.Equal(t, *c.Description, company.Description)
	assert.Equal(t, *c.EmployeesAmount, company.EmployeesAmount)
	assert.Equal(t, *c.Registered, company.Registered)
	assert.Equal(t, *c.Type, company.Type)
	assert.NotEmpty(t, company.UpdatedAt)
}

func newStringPointer(str string) *string {
	return &str
}
