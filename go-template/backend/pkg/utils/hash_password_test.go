package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const pass = "password"

func TestHashPassword(t *testing.T) {
	p, err := HashPassword(pass)
	assert.Nil(t, err)
	s := CheckPasswordHash(pass, p)
	assert.True(t, s)
}

func TestCheckPasswordHash(t *testing.T) {
	p, err := HashPassword(pass)
	assert.Nil(t, err)

	s := CheckPasswordHash(pass, p)
	assert.True(t, s)

	s2 := CheckPasswordHash("asdf", p)
	assert.False(t, s2)
}
