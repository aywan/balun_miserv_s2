package core_test

import (
	"testing"

	"github.com/aywan/balun_miserv_s2/chat-server/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestDoTrue(t *testing.T) {
	assert.True(t, core.DoTrue())
}
