package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionParsingSimple(t *testing.T) {
	v := Parse("1.0.0")
	assert.Equal(t, uint8(3), v.numParts)
	assert.Equal(t, []byte{}, v.suffix)
	assert.Equal(t, [][]byte{[]byte("1"), []byte("0"), []byte("0")}, v.parts)
}

func TestVersionParsingWithSuffix(t *testing.T) {
	v := Parse("1.0.0-rc0")
	assert.Equal(t, uint8(3), v.numParts)
	assert.Equal(t, []byte("rc0"), v.suffix)
	assert.Equal(t, [][]byte{[]byte("1"), []byte("0"), []byte("0")}, v.parts)
}

func TestVersionIncrementBuild(t *testing.T) {
	v := NewVersion(0, 0, 0, 1)
	if err := v.Increment(3, 1); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, "0.0.0.2", v.String())
	}
}

func TestVersionIncrementPatch(t *testing.T) {
	v := NewVersion(0, 0, 1, 1)
	if err := v.Increment(2, 1); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, "0.0.2.0", v.String())
	}
}

func TestVersionIncrementMinor(t *testing.T) {
	v := NewVersion(0, 1, 1, 1)
	if err := v.Increment(1, 1); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, "0.2.0.0", v.String())
	}
}

func TestVersionIncrementMajor(t *testing.T) {
	v := NewVersion(1, 1, 1, 1)
	if err := v.Increment(0, 1); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, "2.0.0.0", v.String())
	}
}

func TestVersionIncrementSuffix(t *testing.T) {
	v := NewVersion(1, 0, 0, 0)
	v.suffix = []byte("rc0")
	if err := v.Increment(4, 1); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, "1.0.0.0-rc1", v.String())
	}
}

func TestVersionDecrementBuild(t *testing.T) {
	v := NewVersion(0, 0, 0, 2)
	if err := v.Decrement(3, 1); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, "0.0.0.1", v.String())
	}
}

func TestVersionDecrementPatch(t *testing.T) {
	v := NewVersion(0, 0, 2, 1)
	if err := v.Decrement(2, 1); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, "0.0.1.0", v.String())
	}
}

func TestVersionDecrementMinor(t *testing.T) {
	v := NewVersion(0, 2, 1, 1)
	if err := v.Decrement(1, 1); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, "0.1.0.0", v.String())
	}
}

func TestVersionDecrementMajor(t *testing.T) {
	v := NewVersion(2, 1, 1, 1)
	if err := v.Decrement(0, 1); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, "1.0.0.0", v.String())
	}
}

func TestVersionDecrementSuffix(t *testing.T) {
	v := NewVersion(1, 0, 0, 0)
	v.suffix = []byte("rc1")
	if err := v.Decrement(4, 1); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, "1.0.0.0-rc0", v.String())
	}
}
