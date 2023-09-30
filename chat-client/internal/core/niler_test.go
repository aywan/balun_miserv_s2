package core_test

import (
	"testing"

	"github.com/aywan/balun_miserv_s2/chat-client/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestDoNil(t *testing.T) {
	assert.Nil(t, core.DoNil())
}
