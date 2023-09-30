package core_test

import (
	"testing"

	"github.com/aywan/balun_miserv_s2/auth-server/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestDoFalse(t *testing.T) {
	assert.False(t, core.DoFalse())
}
