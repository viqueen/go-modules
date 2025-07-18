package pet_test

import (
	"testing"

	"github.com/ovechkin-dm/mockio/v2/mock"
	"github.com/viqueen/go-template/pkg/pet"

	"github.com/stretchr/testify/assert"
)

func TestProject(t *testing.T) {
	t.Parallel()
	ctrl := mock.NewMockController(t)
	project := mock.Mock[pet.Project](ctrl)
	mock.WhenSingle(project.Get()).ThenReturn("project")
	assert.Equal(t, "project", project.Get())
}
