package errorsx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const msg = "test"

func TestNewErrorx(t *testing.T) {
	err := NewErrorx(ErrorTypeNoType, msg)
	s := err.Error()
	assert.Equal(t, s, msg)
}

func TestNew(t *testing.T) {
	e := New(msg)
	assert.Equal(t, e.Error(), msg)
}

func TestIs(t *testing.T) {
	err1 := New(msg)
	r := Is(err1, err1)
	assert.True(t, r)

	err2 := New("asdf")
	r = Is(err1, err2)
	assert.False(t, r)
}
