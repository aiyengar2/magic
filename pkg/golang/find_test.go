package golang

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindModule(t *testing.T) {
	module, err := FindModule(".")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "github.com/aiyengar2/magic", module, "found wrong interface")
}
