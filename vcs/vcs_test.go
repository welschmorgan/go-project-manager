package vcs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVCSShouldRetrieveInstanceByName(t *testing.T) {
	if len(All) == 0 {
		t.Fatal("no vcs configured")
	}
	expected := All[0]
	actual := Get(expected.Name())

	assert.Equal(t, expected.Name(), actual.Name())
}
