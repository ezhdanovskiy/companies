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

// todo: write more tests

func newStringPointer(str string) *string {
	return &str
}
