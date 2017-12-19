package main

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestData(t *testing.T) {
	input := `{"content":"中文测试","owner":"root","group":"root","mode":"400"}`
	parsed := ParseInput("/lain/app", input)
	assert.Equal(t, parsed.Content, "中文测试")
}
