package main

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestSlug(t *testing.T) {
	assert.Equal(t, true, slugPattern.MatchString("app-hello"))
	assert.Equal(t, false, slugPattern.MatchString("app-_hello"))
	assert.Equal(t, true, slugPattern.MatchString("app_hello"))
	assert.Equal(t, false, slugPattern.MatchString("app--hello"))
	assert.Equal(t, true, slugPattern.MatchString("app&hello"))
}
