package gorel

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAssertWhichFails(t *testing.T) {
	assert.Panics(t, func() {
		Assert(false, "expected %v but received %v", true, false)
	})
}

func TestAssertWhichPasses(t *testing.T) {
	assert.NotPanics(t, func() {
		Assert(true, "expected %v but received %v", true, false)
	})
}
